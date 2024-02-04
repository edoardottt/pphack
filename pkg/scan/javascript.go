/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"os"

	"github.com/projectdiscovery/gologger"
)

func getJavascript(r *Runner, testPayload string) string {
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
