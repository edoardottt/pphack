/*
pphack - The Most Advanced Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"net/url"
	"strings"

	"github.com/edoardottt/pphack/pkg/input"
)

const (
	MinURLLength = 4
)

// PrepareURL takes as input a string and prepares
// the input URL in order to get the favicon icon.
func PrepareURL(inputURL string) (string, string, error) {
	if len(inputURL) < MinURLLength {
		return "", "", input.ErrMalformedURL
	}

	if !strings.Contains(inputURL, "://") {
		inputURL = "http://" + inputURL
	}

	u, err := url.Parse(inputURL)
	if err != nil {
		return "", "", err
	}

	payload, testPayload := GenPayload()

	return u.Scheme + "://" + u.Host + u.Path + "?" + payload, testPayload, nil
}
