/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package output

import (
	"bytes"
	"encoding/json"
)

// JSONData struct holds the JSON output data.
type ResultData struct {
	// TargetURL is the original target URL.
	TargetURL string `json:"TargetURL,omitempty"`
	// ScanURL is the target URL + payload used for prototype pollution scan.
	ScanURL string `json:"ScanURL,omitempty"`
	// JSEvaluation is the JS result after prototype pollution scan
	// in the browser console.
	JSEvaluation string `json:"JSEvaluation,omitempty"`
	// ScanError is the error after prototype pollution scan (if present).
	ScanError string `json:"ScanError,omitempty"`
	// Fingerprint is the JS result after fingerprint scan
	// in the browser console.
	Fingerprint []string `json:"Fingerprint,omitempty"`
	// FingerprintError is the error after fingerprint scan (if present).
	FingerprintError string `json:"FingerprintError,omitempty"`
	// ExploitURLs are the URLs crafted to exploit the target URL (if present).
	ExploitURLs []string `json:"ExploitURLs,omitempty"`
	// ExploitError is the error after exploit scan (if present).
	ExploitError string `json:"ExploitError,omitempty"`
	// References are the links to read more on the vulnerable target exploitation.
	References []string `json:"References,omitempty"`
}

// FormatJSON returns the input as JSON string.
func FormatJSON(input *ResultData) ([]byte, error) {
	var jsonOutput bytes.Buffer

	enc := json.NewEncoder(&jsonOutput)
	enc.SetEscapeHTML(false)

	err := enc.Encode(input)
	if err != nil {
		return nil, err
	}

	return jsonOutput.Bytes(), nil
}
