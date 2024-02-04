/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/edoardottt/golazy"
	"github.com/edoardottt/pphack/pkg/input"
	"github.com/edoardottt/pphack/pkg/output"
	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
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

	/*

		var headers map[string]interface{}

		if len(r.Options.Headers) != 0 || r.Options.HeadersFile != "" {
			h, err := GetHeaders(r)
			if err != nil {
				gologger.Fatal().Msg(err.Error())
			}

			headers = h
		}
	*/

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
				if err != nil {
					if r.Options.Verbose {
						gologger.Error().Msg(err.Error())
					}
				}

				ctx, cancel := context.WithTimeout(pctx, time.Second*time.Duration(r.Options.Timeout))
				ctx, _ = chromedp.NewContext(ctx)

				rl.Take()

				res, err := Scan(ctx, js, targetURL)
				if err != nil {
					if r.Options.Verbose {
						gologger.Error().Msg(err.Error())
					}
				}

				if resTrimmed := strings.TrimSpace(res); resTrimmed != "" {
					writeOutput(r, targetURL)
				}

				cancel()
			}

			wg.Done()
		}()
	}

	pushInput(r)

	wg.Wait()
}

func writeOutput(r *Runner, targetURL string) {
	if !r.Result.Printed(targetURL) {
		write(r.OutMutex, &r.Options, targetURL)
	}
}

func write(m *sync.Mutex, options *input.Options, o string) {
	if options.FileOutput != "" && options.Output == nil {
		file, err := os.OpenFile(options.FileOutput, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			gologger.Fatal().Msg(err.Error())
		}

		options.Output = file
	}

	m.Lock()

	if options.Output != nil {
		if _, err := options.Output.Write([]byte(o + "\n")); err != nil && options.Verbose {
			gologger.Fatal().Msg(err.Error())
		}
	}

	m.Unlock()

	fmt.Println(o)
}
