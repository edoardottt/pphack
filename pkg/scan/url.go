/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"net/url"
	"strings"

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
