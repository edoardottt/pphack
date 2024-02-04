/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"os"

	"github.com/projectdiscovery/gologger"
)

// GetJavascript returns the Javascript code must be run on
// the target to verify the vulnerability.
func GetJavascript(r *Runner, testPayload string) string {
	if r.Options.JS != "" {
		return r.Options.JS
	}

	if r.Options.JSFile != "" {
		js, err := os.ReadFile(r.Options.JSFile)
		if err != nil {
			gologger.Fatal().Msg(err.Error())
		}

		return string(js)
	}

	return "window." + testPayload
}
