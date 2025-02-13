/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package input

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	fileutil "github.com/projectdiscovery/utils/file"
)

var (
	ErrMutexFlags       = errors.New("incompatible flags specified")
	ErrNoInput          = errors.New("no input specified")
	ErrNegativeValue    = errors.New("must be positive")
	ErrMalformedURL     = errors.New("malformed input URL")
	ErrMalformedPayload = errors.New("malformed input payload (follow JavaScript variables naming rules)")
)

func (options *Options) validateOptions() error {
	if options.Silent && options.Verbose {
		return fmt.Errorf("%w: %s and %s", ErrMutexFlags, "silent", "verbose")
	}

	if options.Input == "" && options.FileInput == "" && !fileutil.HasStdin() {
		return fmt.Errorf("%w", ErrNoInput)
	}

	if options.Concurrency <= 0 {
		return fmt.Errorf("concurrency: %w", ErrNegativeValue)
	}

	if options.RateLimit != 0 && options.RateLimit <= 0 {
		return fmt.Errorf("rate limit: %w", ErrNegativeValue)
	}

	if options.JS != "" && options.JSFile != "" {
		return fmt.Errorf("%w: %s and %s", ErrMutexFlags, "javascript", "javascript-file")
	}

	if len(options.Headers) != 0 && options.HeadersFile != "" {
		return fmt.Errorf("%w: %s and %s", ErrMutexFlags, "headers", "headers-file")
	}

	if !payloadOk(options.Payload) {
		return fmt.Errorf("%w", ErrMalformedPayload)
	}

	return nil
}

func payloadOk(payload string) bool {
	if payload == "" {
		return true
	}

	if !(payload[0] == '_' || payload[0] == '$' || isLetter(payload[0])) {
		return false
	}

	if strings.ContainsAny(payload, " ") {
		return false
	}

	if !isAlphanumeric(payload[1:]) {
		return false
	}

	return true
}

func isLetter(r byte) bool {
	if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
		return false
	}

	return true
}

func isAlphanumeric(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(s)
}
