/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"context"

	"github.com/chromedp/chromedp"
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

func GetChromeBrowser(copts []func(*chromedp.ExecAllocator)) (context.CancelFunc,
	context.Context, context.CancelFunc) {
	ectx, ecancel := chromedp.NewExecAllocator(context.Background(), copts...)
	pctx, pcancel := chromedp.NewContext(ectx)

	if err := chromedp.Run(pctx); err != nil {
		gologger.Fatal().Msgf("error starting browser: %s", err.Error())
	}

	return ecancel, pctx, pcancel
}
