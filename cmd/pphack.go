/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package main

import (
	"github.com/edoardottt/pphack/pkg/input"
	"github.com/edoardottt/pphack/pkg/scan"
	"github.com/projectdiscovery/gologger"
)

func main() {
	options := input.ParseOptions()

	runner, err := scan.New(options)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	runner.Run()
}
