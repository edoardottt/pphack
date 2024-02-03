/*
pphack - The Most Advanced Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"net/url"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/edoardottt/pphack/pkg/input"
)

const (
	minURLLength = 4
)

// PrepareURL takes as input a string (URL) and prepares
// the input to be scanned.
func PrepareURL(inputURL, payloadInput string) (string, string, error) {
	if len(inputURL) < minURLLength {
		return "", "", input.ErrMalformedURL
	}

	if !strings.Contains(inputURL, "://") {
		inputURL = "http://" + inputURL
	}

	u, err := url.Parse(inputURL)
	if err != nil {
		return "", "", err
	}

	var (
		payload     string
		testPayload string
	)

	if payloadInput != "" {
		payload, testPayload = GenCustomPayload(payloadInput)
	} else {
		payload, testPayload = GenPayload()
	}

	return u.Scheme + "://" + u.Host + u.Path + "?" + payload, testPayload, nil
}

func getChromeOptions(r *Runner) []func(*chromedp.ExecAllocator) {
	copts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent(r.UserAgent),
	)

	if r.Options.Proxy != "" {
		copts = append(copts, chromedp.ProxyServer(r.Options.Proxy))
	}

	return copts
}
