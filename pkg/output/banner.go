/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package output

import "github.com/projectdiscovery/gologger"

var printed = false

const (
	Version = "v0.1.0"
	banner  = `                 __               __  
    ____  ____  / /_  ____ ______/ /__
   / __ \/ __ \/ __ \/ __ ` + "`" + `/ ___/ //_/
  / /_/ / /_/ / / / / /_/ / /__/ ,<   
 / .___/ .___/_/ /_/\__,_/\___/_/|_|  
/_/   /_/                             `
)

// ShowBanner prints the banner.
func ShowBanner() {
	if !printed {
		gologger.Print().Msgf("%s%s\n\n", banner, Version)
		gologger.Print().Msgf("\t\t@edoardottt, https://edoardottt.com/\n")
		gologger.Print().Msgf("\t\t             https://github.com/edoardottt/\n\n")

		printed = true
	}
}
