/*
pphack - The Most Advanced Prototype Pollution Scanner

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
	Input     []string
	InputChan chan string
	Result    output.Result
	UserAgent string
	Options   input.Options
	OutMutex  *sync.Mutex
}

func New(options *input.Options) Runner {
	if options.FileOutput != "" {
		_, err := os.Create(options.FileOutput)
		if err != nil {
			gologger.Error().Msgf("%s", err)
		}
	}

	return Runner{
		Input:     []string{},
		InputChan: make(chan string, options.Concurrency),
		Result:    output.New(),
		UserAgent: golazy.GenerateRandomUserAgent(),
		Options:   *options,
		OutMutex:  &sync.Mutex{},
	}
}

func (r *Runner) Run() {
	copts := getChromeOptions(r)

	ectx, ecancel := chromedp.NewExecAllocator(context.Background(), copts...)
	defer ecancel()

	pctx, pcancel := chromedp.NewContext(ectx)
	defer pcancel()

	if err := chromedp.Run(pctx); err != nil {
		gologger.Fatal().Msgf("error starting browser: %s", err.Error())
	}

	var wg sync.WaitGroup

	for i := 0; i < r.Options.Concurrency; i++ {
		wg.Add(1)

		go func() {
			for value := range r.InputChan {
				targetURL, payload, err := PrepareURL(value, r.Options.Payload)
				if err != nil {
					if r.Options.Verbose {
						gologger.Error().Msg(err.Error())
					}
				}

				ctx, cancel := context.WithTimeout(pctx, time.Second*time.Duration(r.Options.Timeout))

				ctx, _ = chromedp.NewContext(ctx)

				var res string

				err = chromedp.Run(ctx,
					chromedp.Navigate(targetURL),
					chromedp.Evaluate("window."+payload, &res),
				)
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
