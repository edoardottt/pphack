/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import "go.uber.org/ratelimit"

func rateLimiter(r *Runner) ratelimit.Limiter {
	var ratelimiter ratelimit.Limiter
	if r.Options.RateLimit > 0 {
		ratelimiter = ratelimit.New(r.Options.RateLimit)
	} else {
		ratelimiter = ratelimit.NewUnlimited()
	}

	return ratelimiter
}
