/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package output

import (
	"fmt"
	"sync"

	"github.com/projectdiscovery/gologger"
)

// Result is the struct handling the output.
type Result struct {
	Map   map[string]struct{}
	Mutex *sync.RWMutex
}

// New returns a new Result object.
func New() Result {
	return Result{
		Map:   map[string]struct{}{},
		Mutex: &sync.RWMutex{},
	}
}

// Printed checks if a string was already printed.
func (o *Result) Printed(result string) bool {
	o.Mutex.RLock()

	if _, ok := o.Map[result]; !ok {
		o.Mutex.RUnlock()
		o.Mutex.Lock()
		o.Map[result] = struct{}{}
		o.Mutex.Unlock()

		return false
	} else {
		o.Mutex.RUnlock()
	}

	return true
}

// JSONOutput marshals and prints the JSON output.
func JSONOutput(json *ResultData) {
	data, err := FormatJSON(json)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	fmt.Println(string(data))
}
