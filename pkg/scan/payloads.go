/*
pphack - The Most Advanced Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"fmt"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// GenPayload returns a ready to use HTTP GET query and the payload
// used in the query.
func GenPayload() (string, string) {
	testPayload := randStringBytes(15)

	var payload = fmt.Sprintf("constructor.prototype." + testPayload + "=" + testPayload +
		"&__proto__[" + testPayload + "]=" + testPayload +
		"&constructor[prototype][" + testPayload + "]=" + testPayload +
		"&__proto__." + testPayload + "=" + testPayload +
		"&__proto__." + testPayload + "=1|2|3" +
		"&__proto__[" + testPayload + "]={\"json\":\"value\"}" +
		"#__proto__[" + testPayload + "]=" + testPayload)

	return payload, testPayload
}
