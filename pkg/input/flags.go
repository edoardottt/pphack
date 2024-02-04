/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package input

import (
	"io"
	"os"
	"strings"

	"github.com/edoardottt/pphack/pkg/output"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

const (
	DefaultTimeout     = 10
	DefaultConcurrency = 50
	DefaultRateLimit   = 0
)

// Options struct holds all the configuration settings.
type Options struct {
	Input       string
	FileInput   string
	FileOutput  string
	Payload     string
	Output      io.Writer
	Concurrency int
	Timeout     int
	Proxy       string
	JS          string
	JSFile      string
	RateLimit   int
	Silent      bool
	Verbose     bool
}

// configureOutput configures the output on the screen.
func (options *Options) configureOutput() {
	if options.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	} else if options.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	}
}

// ParseOptions parses the command line options for application.
func ParseOptions() *Options {
	options := &Options{}

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`The Most Advanced Client-Side Prototype Pollution Scanner.`)

	// Input.
	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&options.Input, "url", "u", "", `Input URL`),
		flagSet.StringVarP(&options.FileInput, "list", "l", "", `File containing input URLs`),
	)

	// Config.
	flagSet.CreateGroup("config", "Configuration",
		flagSet.IntVarP(&options.Concurrency, "concurrency", "c", DefaultConcurrency, `Concurrency level`),
		flagSet.IntVarP(&options.Timeout, "timeout", "t", DefaultTimeout, `Connection timeout in seconds`),
		flagSet.StringVarP(&options.Proxy, "proxy", "px", "", `Set a proxy server (URL)`),
		flagSet.IntVarP(&options.RateLimit, "rate-limit", "rl", DefaultRateLimit, `Set a rate limit (per second)`),
	)

	// Scan.
	flagSet.CreateGroup("scan", "Scan",
		flagSet.StringVarP(&options.Payload, "payload", "p", "", `Custom payload`),
		flagSet.StringVarP(&options.JS, "javascript", "js", "", `Run custom Javascript on target`),
		flagSet.StringVarP(&options.JSFile, "javascript-file", "jsf", "",
			`File containing custom Javascript to run on target`),
	)

	// Output.
	flagSet.CreateGroup("output", "Output",
		flagSet.StringVarP(&options.FileOutput, "output", "o", "", `File to write output results`),
		flagSet.BoolVarP(&options.Verbose, "verbose", "v", false, `Verbose output`),
		flagSet.BoolVarP(&options.Silent, "silent", "s", false, `Silent output. Print only results`),
	)

	if help() || noArgs() {
		output.ShowBanner()
	}

	if err := flagSet.Parse(); err != nil {
		output.ShowBanner()
		gologger.Fatal().Msgf("%s\n", err)
	}

	// Read the inputs and configure the logging.
	options.configureOutput()

	if err := options.validateOptions(); err != nil {
		output.ShowBanner()
		gologger.Fatal().Msgf("%s\n", err)
	}

	output.ShowBanner()

	return options
}

func help() bool {
	// help usage asked by user.
	for _, arg := range os.Args {
		argStripped := strings.Trim(arg, "-")
		if argStripped == "h" || argStripped == "help" {
			return true
		}
	}

	return false
}

func noArgs() bool {
	// User passed no flag.
	return len(os.Args) < 2
}
