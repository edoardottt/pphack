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

func GetChromeBrowser(copts []func(*chromedp.ExecAllocator)) (context.CancelFunc, context.Context, context.CancelFunc) {
	ectx, ecancel := chromedp.NewExecAllocator(context.Background(), copts...)
	pctx, pcancel := chromedp.NewContext(ectx)

	if err := chromedp.Run(pctx); err != nil {
		ecancel()
		gologger.Fatal().Msgf("error starting browser: %s", err.Error())
	}

	return ecancel, pctx, pcancel
}

func buildHeaders(headers map[string]interface{}) chromedp.Tasks {
	if headers == nil {
		return nil
	}

	return chromedp.Tasks{network.SetExtraHTTPHeaders(network.Headers(headers))}
}

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

	ctx, ctxCancel := context.WithTimeout(pctx, time.Second*time.Duration(r.Options.Timeout))
	defer ctxCancel()

	tabCtx, tabCancel := chromedp.NewContext(ctx)
	defer tabCancel()

	scanTasks := buildHeaders(headers)
	scanTasks = append(
		scanTasks,
		chromedp.Navigate(targetURL),
		chromedp.EvaluateAsDevTools(js, &resScan),
	)

	errScan := chromedp.Run(tabCtx, scanTasks)
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

	fingerprintTasks := chromedp.Tasks{
		chromedp.EvaluateAsDevTools(exploit.Fingerprint, &resDetection),
	}

	errDetection := chromedp.Run(tabCtx, fingerprintTasks)
	if errDetection != nil {
		gologger.Error().Msg(errDetection.Error())
		resultData.FingerprintError = errDetection.Error()
	}

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
