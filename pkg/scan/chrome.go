/*
pphack - Advanced Client-Side Prototype Pollution Scanner

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
// the chrome options used to configure the headless browser instance.
// It always disables certificate errors and sets a custom user agent.
// If a proxy is configured in the runner options, it is appended as well.
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
// the context with the associated cancel functions to use the
// headless chrome browser it creates.
// Returns ecancel (exec allocator cancel), pctx (parent browser context),
// and pcancel (parent context cancel).
// Callers must invoke pcancel before ecancel to ensure correct cleanup order.
// ecancel is also called internally on fatal browser startup failure
// to avoid leaking the exec allocator before the process exits.
func GetChromeBrowser(copts []func(*chromedp.ExecAllocator)) (context.CancelFunc, context.Context, context.CancelFunc) {
	ectx, ecancel := chromedp.NewExecAllocator(context.Background(), copts...)
	pctx, pcancel := chromedp.NewContext(ectx)

	// Run an empty chromedp step list to verify the browser starts successfully.
	if err := chromedp.Do(pctx); err != nil {
		ecancel()
		gologger.Fatal().Msgf("error starting browser: %s", err.Error())
	}

	return ecancel, pctx, pcancel
}

// buildHeaders is a helper that converts a headers map into a slice of
// chromedp actions containing the SetExtraHTTPHeaders command.
// Returns nil if headers is nil, making it safe to append directly onto any
// existing action slice without an extra nil check at the call site.
func buildHeaders(headers map[string]interface{}) []chromedp.Action[chromedp.Void] {
	if headers == nil {
		return nil
	}

	return []chromedp.Action[chromedp.Void]{
		chromedp.Func(func(ctx context.Context, _ *chromedp.Target) error {
			_, err := chromedp.Call(ctx, network.SetExtraHTTPHeaders,
				network.SetExtraHTTPHeadersParams{Headers: network.Headers(headers)})
			return err
		}),
	}
}

// Scan is the core function that performs the prototype pollution scan.
// It takes a parent browser context (pctx), runner config (r), optional HTTP
// headers, the JavaScript payload (js), the original input value and the
// fully constructed target URL.
//
// Flow:
//  1. Creates a timeout-scoped context and a dedicated Chrome tab context.
//  2. Navigates to targetURL and evaluates the JS pollution payload.
//  3. If exploit mode is enabled and the payload returned a non-empty result,
//     it runs fingerprinting to identify the affected library/sink.
//  4. Attempts exploitation using the fingerprint results.
//  5. Populates and returns a ResultData struct with all findings and errors.
func Scan(
	pctx context.Context,
	r *Runner,
	headers map[string]interface{},
	js, value, targetURL string,
) (output.ResultData, error) {
	var (
		resScan      string
		resDetection []string
	)

	resultData := output.ResultData{
		TargetURL: value,
		ScanURL:   targetURL,
	}

	// Wrap the parent context with a per-scan timeout to avoid blocking.
	ctx, ctxCancel := context.WithTimeout(pctx, time.Second*time.Duration(r.Options.Timeout))
	defer ctxCancel()

	// Open a new Chrome tab scoped to the timeout context.
	// tabCancel explicitly closes the tab when Scan returns
	// preventing tab accumulation across concurrent scans.
	tabCtx, tabCancel := chromedp.NewContext(ctx)
	defer tabCancel()

	// Build the scan step list: optionally inject custom HTTP headers, then
	// navigate to the target.
	scanSteps := buildHeaders(headers)
	scanSteps = append(scanSteps, chromedp.Navigate(targetURL))

	// Execute the scan steps inside the dedicated tab context, then evaluate
	// the prototype pollution JS payload and capture its result.
	errScan := chromedp.Do(tabCtx, scanSteps...)
	if errScan == nil {
		resScan, errScan = chromedp.Run(tabCtx, chromedp.EvaluateAsDevTools[string](js))
	}

	if errScan != nil {
		resultData.ScanError = errScan.Error()
	}

	resultData.JSEvaluation = strings.TrimSpace(resScan)

	if !r.Options.Exploit || errScan != nil || resultData.JSEvaluation == "" {
		return resultData, nil
	}

	if r.Options.Verbose {
		gologger.Info().Label("VULN").Msg(fmt.Sprintf("Target is Vulnerable %s", targetURL))
	}

	// Run fingerprinting as a separate, isolated evaluation.
	resDetection, errDetection := chromedp.Run(tabCtx, chromedp.EvaluateAsDevTools[[]string](exploit.Fingerprint))
	if errDetection != nil {
		gologger.Error().Msg(errDetection.Error())
		resultData.FingerprintError = errDetection.Error()
	}

	// Store fingerprint results and cross-reference known exploit references.
	resultData.Fingerprint = resDetection
	resultData.References = exploit.GetReferences(resDetection)

	if r.Options.Verbose {
		gologger.Info().Msg(fmt.Sprintf("Trying to exploit %s", value))
	}

	exploitTasks := buildHeaders(headers)

	result, errExploit := exploit.CheckExploit(
		pctx,
		exploitTasks,
		resDetection,
		targetURL,
		r.Options.Verbose,
		r.Options.Timeout,
	)

	resultData.ExploitURLs = result

	if errExploit != nil {
		resultData.ExploitError = errExploit.Error()
		gologger.Error().Msg(errExploit.Error())
	}

	return resultData, nil
}
