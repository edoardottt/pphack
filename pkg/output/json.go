/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package output

import "encoding/json"

// JSONData.
type JSONData struct {
	URL          string `json:"URL,omitempty"`
	JSEvaluation string `json:"JSEvaluation,omitempty"`
	Error        string `json:"Error,omitempty"`
}

// FormatJSON returns the input as JSON string.
func FormatJSON(url, jsEval, e string) ([]byte, error) {
	input := &JSONData{
		URL:          url,
		JSEvaluation: jsEval,
		Error:        e,
	}

	jsonOutput, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return jsonOutput, nil
}
