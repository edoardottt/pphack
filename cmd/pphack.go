/*
pphack - The Most Advanced Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package main

import (
	"github.com/edoardottt/pphack/pkg/input"
	"github.com/edoardottt/pphack/pkg/scan"
)

func main() {
	options := input.ParseOptions()
	runner := scan.New(options)
	runner.Run()
}
