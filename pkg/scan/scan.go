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
	Input     chan string
	Output    chan string
	Result    output.Result
	UserAgent string
	InWg      *sync.WaitGroup
	OutWg     *sync.WaitGroup
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
		Input:     make(chan string, options.Concurrency),
		Output:    make(chan string, options.Concurrency),
		Result:    output.New(),
		UserAgent: golazy.GenerateRandomUserAgent(),
		InWg:      &sync.WaitGroup{},
		OutWg:     &sync.WaitGroup{},
		Options:   *options,
		OutMutex:  &sync.Mutex{},
	}
}

func (r *Runner) Run() {
	r.InWg.Add(1)

	go pushInput(r)
	r.InWg.Add(1)

	go execute(r)
	r.OutWg.Add(1)

	go pullOutput(r)
	r.InWg.Wait()

	close(r.Output)
	r.OutWg.Wait()
}

func pushInput(r *Runner) {
	defer r.InWg.Done()

	if fileutil.HasStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			r.Input <- scanner.Text()
		}
	}

	if r.Options.FileInput != "" {
		for _, line := range golazy.RemoveDuplicateValues(golazy.ReadFileLineByLine(r.Options.FileInput)) {
			r.Input <- line
		}
	}

	if r.Options.Input != "" {
		r.Input <- r.Options.Input
	}

	close(r.Input)
}

func execute(r *Runner) {
	defer r.InWg.Done()

	for i := 0; i < r.Options.Concurrency; i++ {
		r.InWg.Add(1)

		go func() {
			defer r.InWg.Done()

			for value := range r.Input {
				targetURL, payload, err := PrepareURL(value)
				if err != nil {
					if r.Options.Verbose {
						gologger.Error().Msgf("%s", err)
					}

					return
				}

				// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! this should go out of the loop above
				copts := append(chromedp.DefaultExecAllocatorOptions[:],
					chromedp.Flag("ignore-certificate-errors", true),
					chromedp.Flag("disable-extensions", true),
					chromedp.Flag("disable-client-side-phishing-detection", true),
					chromedp.Flag("disable-popup-blocking", true),
					chromedp.UserAgent(r.UserAgent),
				)

				ectx, ecancel := chromedp.NewExecAllocator(context.Background(), copts...)
				defer ecancel()

				pctx, pcancel := chromedp.NewContext(ectx)
				defer pcancel()

				if err := chromedp.Run(pctx); err != nil {
					gologger.Fatal().Msgf("error starting browser: %s", err)
				}
				// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! this should go out of the loop above

				ctx, cancel := context.WithTimeout(pctx, time.Second*time.Duration(r.Options.Timeout))
				ctx, _ = chromedp.NewContext(ctx)

				var res string

				err = chromedp.Run(ctx,
					chromedp.Navigate(targetURL),
					chromedp.Evaluate("window."+payload, &res),
				)
				if err != nil {
					if r.Options.Verbose {
						gologger.Error().Msgf("%s", err)
					}

					cancel()
					return
				}

				if resTrimmed := strings.TrimSpace(res); resTrimmed != "" {
					r.Output <- targetURL
				}

				cancel()
			}
		}()
	}
}

func pullOutput(r *Runner) {
	defer r.OutWg.Done()

	for o := range r.Output {
		if !r.Result.Printed(o) {
			r.OutWg.Add(1)

			go writeOutput(r.OutWg, r.OutMutex, &r.Options, o)
		}
	}
}

func writeOutput(wg *sync.WaitGroup, m *sync.Mutex, options *input.Options, o string) {
	defer wg.Done()

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
