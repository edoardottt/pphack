/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/edoardottt/golazy"
	"github.com/edoardottt/pphack/pkg/input"
	"github.com/edoardottt/pphack/pkg/output"
	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
)

const (
	DefaultFilePerm = 0644
)

type Runner struct {
	InputChan chan string
	Result    output.Result
	UserAgent string
	Options   input.Options
	OutMutex  *sync.Mutex
}

func New(options *input.Options) (Runner, error) {
	r := Runner{}

	if options.FileOutput != "" {
		_, err := os.Create(options.FileOutput)
		if err != nil {
			return r, err
		}
	}

	r = Runner{
		InputChan: make(chan string, options.Concurrency),
		Result:    output.New(),
		Options:   *options,
		OutMutex:  &sync.Mutex{},
	}

	if options.UserAgent != "" {
		r.UserAgent = options.UserAgent
	} else {
		r.UserAgent = golazy.GenerateRandomUserAgent()
	}

	return r, nil
}

func pushInput(r *Runner) {
	if fileutil.HasStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			r.InputChan <- scanner.Text()
		}
	}

	if r.Options.FileInput != "" {
		for _, line := range golazy.RemoveDuplicateValues(golazy.ReadFileLineByLine(r.Options.FileInput)) {
			r.InputChan <- line
		}
	}

	if r.Options.Input != "" {
		r.InputChan <- r.Options.Input
	}

	close(r.InputChan)
}

func (r *Runner) Run() {
	copts := GetChromeOptions(r)
	ecancel, pctx, pcancel := GetChromeBrowser(copts)
	testPayload := GetTestPayload(r, payloadLength)
	js := GetJavascript(r, testPayload)

	var headers map[string]interface{}

	if len(r.Options.Headers) != 0 || r.Options.HeadersFile != "" {
		h, err := GetHeaders(r)
		if err != nil {
			gologger.Fatal().Msg(err.Error())
		}

		headers = h
	}

	defer ecancel()
	defer pcancel()

	var (
		rl = rateLimiter(r)
		wg sync.WaitGroup
	)

	for i := 0; i < r.Options.Concurrency; i++ {
		wg.Add(1)

		go func() {
			for value := range r.InputChan {
				targetURL, err := PrepareURL(value, testPayload)
				if err != nil && !r.Options.Silent {
					gologger.Error().Msg(fmt.Sprintf("%s: %s", value, err.Error()))
				}

				rl.Take()

				if r.Options.Verbose {
					gologger.Info().Msg(fmt.Sprintf("Targeting %s", value))
				}

				resultData, err := Scan(pctx, r, headers, js, value, targetURL)
				if err != nil && !r.Options.Silent {
					gologger.Error().Msg(err.Error())
				}

				if resTrimmed := strings.TrimSpace(resultData.JSEvaluation); resTrimmed != "" {
					writeOutput(r, resultData)
				}
			}

			wg.Done()
		}()
	}

	pushInput(r)

	wg.Wait()
}

func writeOutput(r *Runner, resultData output.ResultData) {
	if !r.Result.Printed(resultData.TargetURL) {
		write(r.OutMutex, &r.Options, resultData)
	}
}

func write(m *sync.Mutex, options *input.Options, resultData output.ResultData) {
	if options.FileOutput != "" && options.Output == nil {
		file, err := os.OpenFile(options.FileOutput, os.O_CREATE|os.O_RDWR|os.O_APPEND, DefaultFilePerm)
		if err != nil {
			gologger.Fatal().Msg(err.Error())
		}

		options.Output = file
	}

	var (
		o   []byte
		err error
	)

	if options.JSON {
		o, err = output.FormatJSON(&resultData)
		if err != nil {
			gologger.Fatal().Msg(err.Error())
		}

		fmt.Println(string(o))
	} else {
		if options.Exploit {
			if len(resultData.ExploitURLs) != 0 {
				var str string
				for _, e := range resultData.ExploitURLs {
					str += fmt.Sprintf("[EXPLOIT] %s\n", e)
					gologger.Info().Label("EXPLOIT").Msg(e)
				}

				o = []byte(strings.TrimSuffix(str, "\n"))
			} else {
				gologger.Info().Label("VULN").Msg(resultData.ScanURL)

				str := fmt.Sprintf("[VULN] %s\n", resultData.ScanURL)
				for _, e := range resultData.References {
					str += fmt.Sprintf("[REFERENCE] Target is vulnerable but cannot reproduce an exploit, see %s\n", e)
					gologger.Info().Label("REFERENCE").Msg(
						fmt.Sprintf("Target is vulnerable but cannot reproduce an exploit, see %s", e))
				}

				o = []byte(strings.TrimSuffix(str, "\n"))
			}
		} else {
			gologger.Info().Label("VULN").Msg(resultData.ScanURL)
			o = []byte(fmt.Sprintf("[VULN] %s", resultData.ScanURL))
		}
	}

	m.Lock()

	if options.Output != nil {
		if _, err := options.Output.Write([]byte(string(o) + "\n")); err != nil && options.Verbose {
			gologger.Error().Msg(err.Error())
		}
	}

	m.Unlock()
}
