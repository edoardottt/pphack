/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package scan

import (
	"fmt"
	"math/rand"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyz"
	payloadLength = 6
)

// GetTestPayload returns the payload specified as input
// or a random payload with a specified length.
func GetTestPayload(r *Runner, length int) string {
	if r.Options.Payload != "" {
		return r.Options.Payload
	}

	return randStringBytes(length)
}

// GenQueryPayload returns a ready to use HTTP GET query with a random generated payload
// and the payload used in the query.
func GenQueryPayload() (string, string) {
	testPayload := randStringBytes(payloadLength)
	payload := genQueryPayload(testPayload)

	return payload, testPayload
}

// GenCustomPayload returns a ready to use HTTP GET query with the payload
// supplied as input.
func GenCustomQueryPayload(testPayload string) string {
	var payload = genQueryPayload(testPayload)
	return payload
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func genQueryPayload(testPayload string) string {
	return fmt.Sprintf("constructor.prototype." + testPayload + "=" + testPayload +
		"&__proto__[" + testPayload + "]=" + testPayload +
		"&constructor[prototype][" + testPayload + "]=" + testPayload +
		"&__proto__." + testPayload + "=" + testPayload +
		"&__proto__." + testPayload + "=1|2|3" +
		"&__proto__[" + testPayload + "]={\"json\":\"value\"}" +
		"#__proto__[" + testPayload + "]=" + testPayload)
}
