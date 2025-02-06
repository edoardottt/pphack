/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/edoardottt/pphack/pkg/exploit"
	"github.com/edoardottt/pphack/pkg/output"
	"github.com/projectdiscovery/gologger"
)

// GetChromeOptions takes as input the runner settings and returns
// the chrome options.
func GetChromeOptions(r *Runner) []func(*chromedp.ExecAllocator) {
	copts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent(r.UserAgent),
	)

	if r.Options.Proxy != "" {
		copts = append(copts, chromedp.ProxyServer(r.Options.Proxy))
	}

	return copts
}

// GetChromeBrowser takes as input the chrome options and returns
// the contexts with the associated cancel functions to use the
// headless chrome browser it creates.
func GetChromeBrowser(copts []func(*chromedp.ExecAllocator)) (context.CancelFunc,
	context.Context, context.CancelFunc) {
	ectx, ecancel := chromedp.NewExecAllocator(context.Background(), copts...)
	pctx, pcancel := chromedp.NewContext(ectx)

	if err := chromedp.Run(pctx); err != nil {
		gologger.Fatal().Msgf("error starting browser: %s", err.Error())
	}

	return ecancel, pctx, pcancel
}

// Scan is the actual function that takes as input a browser context, other info
// and performs the scan.
func Scan(pctx context.Context, r *Runner, headers map[string]interface{},
	js, value, targetURL string) (output.ResultData, error) {
	var (
		resScan                string
		resDetection           []string
		chromedpTasksScan      chromedp.Tasks
		chromedpTasksDetection chromedp.Tasks
	)

	resultData := output.ResultData{
		TargetURL: value,
		ScanURL:   targetURL,
	}

	if headers != nil {
		chromedpTasksScan = append(chromedpTasksScan, network.SetExtraHTTPHeaders(network.Headers(headers)))
	}

	chromedpTasksScan = append(
		chromedpTasksScan, chromedp.Navigate(targetURL),
		chromedp.EvaluateAsDevTools(js, &resScan),
	)

	ctx, cancel := context.WithTimeout(pctx, time.Second*time.Duration(r.Options.Timeout))
	ctx, _ = chromedp.NewContext(ctx)
	defer cancel()

	errScan := chromedp.Run(ctx, chromedpTasksScan)

	if errScan != nil {
		resultData.ScanError = errScan.Error()
	}

	resultData.JSEvaluation = strings.TrimSpace(resScan)

	// if I have to detect the exploit, no errors and it's vulnerable.
	if r.Options.Exploit && errScan == nil {
		if resTrimmed := strings.TrimSpace(resScan); resTrimmed != "" {
			if r.Options.Verbose {
				gologger.Info().Label("VULN").Msg(fmt.Sprintf("Target is Vulnerable %s", targetURL))
			}

			chromedpTasksScan = append(chromedpTasksScan, chromedp.EvaluateAsDevTools(exploit.Fingerprint, &resDetection))

			errDetection := chromedp.Run(ctx, chromedpTasksScan)
			if errDetection != nil && r.Options.Verbose {
				gologger.Error().Msg(errDetection.Error())
			}

			resultData.Fingerprint = resDetection

			if errDetection != nil {
				resultData.FingerprintError = errDetection.Error()
			}

			if headers != nil {
				chromedpTasksDetection = append(chromedpTasksDetection, network.SetExtraHTTPHeaders(network.Headers(headers)))
			}

			if r.Options.Verbose {
				gologger.Info().Msg(fmt.Sprintf("Trying to exploit %s", value))
			}

			result, errExploit := exploit.CheckExploit(pctx, chromedpTasksDetection, resDetection, targetURL,
				r.Options.Verbose, r.Options.Timeout)

			resultData.ExploitURLs = result

			if errExploit != nil {
				resultData.ExploitError = errDetection.Error()
			}

			if errExploit != nil && !r.Options.Verbose {
				gologger.Error().Msg(errExploit.Error())
			}
		}
	}

	if len(resultData.ExploitURLs) == 0 {
		resultData.References = exploit.GetReferences(resDetection)
	}

	return resultData, nil
}
