/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"context"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/edoardottt/pphack/pkg/exploit"
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
func Scan(pctx context.Context, headers map[string]interface{}, detection bool,
	js, targetURL string) (string, []string, []string, error, error) {
	var (
		resScan                string
		resDetection           []string
		result                 []string
		chromedpTasksScan      chromedp.Tasks
		chromedpTasksDetection chromedp.Tasks
		errScan                error
		errDetection           error
	)

	if headers != nil {
		chromedpTasksScan = append(chromedpTasksScan, network.SetExtraHTTPHeaders(network.Headers(headers)))
	}

	chromedpTasksScan = append(
		chromedpTasksScan, chromedp.Navigate(targetURL),
		chromedp.EvaluateAsDevTools(js, &resScan),
	)

	ctx, cancel := context.WithTimeout(pctx, time.Second*time.Duration(10))
	ctx, _ = chromedp.NewContext(ctx)
	defer cancel()

	errScan = chromedp.Run(ctx, chromedpTasksScan)

	// if I have to detect the exploit, no errors and it's vulnerable.
	if detection && errScan == nil {
		if resTrimmed := strings.TrimSpace(resScan); resTrimmed != "" {
			chromedpTasksScan = append(chromedpTasksScan, chromedp.EvaluateAsDevTools(exploit.Fingerprint, &resDetection))

			errDetection := chromedp.Run(ctx, chromedpTasksScan)
			if errDetection != nil {
				return resScan, resDetection, result, errScan, errDetection
			}

			if headers != nil {
				chromedpTasksDetection = append(chromedpTasksDetection, network.SetExtraHTTPHeaders(network.Headers(headers)))
			}

			result, _ = exploit.CheckExploit(pctx, chromedpTasksDetection, resDetection, targetURL)

			// if err != nil {
			// verbose output
			// }
		}
	}

	return resScan, resDetection, result, errScan, errDetection
}
