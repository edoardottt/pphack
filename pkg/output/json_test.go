/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package output_test

import (
	"testing"

	"github.com/edoardottt/pphack/pkg/output"
	"github.com/stretchr/testify/require"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name  string
		input output.ResultData
		want  string
	}{
		{name: "no error",
			input: output.ResultData{
				TargetURL: "https://edoardottt.github.io/pp-test",
				ScanURL: "https://edoardottt.github.io/pp-test?constructor.prototype.oigczu=oigczu" +
					"&__proto__[oigczu]=oigczu&constructor[prototype][oigczu]=oigczu&__proto__.oigczu=oigczu" +
					"&__proto__.oigczu=1|2|3&__proto__[oigczu]={\"json\":\"value\"}#__proto__[oigczu]=oigczu",
				JSEvaluation: "oigczu",
				ScanError:    "",
			},
			want: `{"TargetURL":"https://edoardottt.github.io/pp-test",` +
				`"ScanURL":"https://edoardottt.github.io/pp-test?constructor.prototype.oigczu=oigczu&__proto__` +
				`[oigczu]=oigczu&constructor[prototype][oigczu]=oigczu&__proto__.oigczu=oigczu&__proto__.` +
				`oigczu=1|2|3&__proto__[oigczu]={\"json\":\"value\"}#__proto__[oigczu]=oigczu",` +
				`"JSEvaluation":"oigczu"}`,
		},
		{name: "error",
			input: output.ResultData{
				TargetURL: "https://edoardottt.github.io/pp-tes",
				ScanURL: "https://edoardottt.github.io/pp-tes?constructor.prototype.sqtiwx=sqtiwx" +
					"&__proto__[sqtiwx]=sqtiwx&constructor[prototype][sqtiwx]=sqtiwx&__proto__.sqtiwx=sqtiwx" +
					"&__proto__.sqtiwx=1|2|3&__proto__[sqtiwx]={\"json\":\"value\"}#__proto__[sqtiwx]=sqtiwx",
				JSEvaluation: "",
				ScanError:    "encountered an undefined value",
			},
			want: `{"TargetURL":"https://edoardottt.github.io/pp-tes",` +
				`"ScanURL":"https://edoardottt.github.io/pp-tes?constructor.prototype.sqtiwx=sqtiwx&__proto__[sqtiwx]` +
				`=sqtiwx&constructor[prototype][sqtiwx]=sqtiwx&__proto__.sqtiwx=sqtiwx&__proto__.` +
				`sqtiwx=1|2|3&__proto__[sqtiwx]={\"json\":\"value\"}#__proto__[sqtiwx]=` +
				`sqtiwx","ScanError":"encountered an undefined value"}`,
		},
		{name: "exploit",
			input: output.ResultData{
				TargetURL: "https://edoardottt.github.io/pp-test",
				ScanURL: "https://edoardottt.github.io/pp-test?constructor.prototype.lfhfqn=lfhfqn" +
					"&__proto__[lfhfqn]=lfhfqn&constructor[prototype][lfhfqn]=lfhfqn&__proto__.lfhfqn=lfhfqn" +
					"&__proto__.lfhfqn=1|2|3&__proto__[lfhfqn]={\"json\":\"value\"}#__proto__[lfhfqn]=lfhfqn",
				JSEvaluation: "lfhfqn",
				Fingerprint:  []string{"jQuery"},
				ExploitURLs: []string{"https://edoardottt.github.io/pp-test/?__proto__[url][]=data:,alert(1337)//" +
					"&__proto__[dataType]=script",
					"https://edoardottt.github.io/pp-test/?__proto__[context]=%3Cimg/src/onerror%3dalert(1337)" +
						"%3E&__proto__[jquery]=x"},
				ScanError: "",
			},
			want: `{"TargetURL":"https://edoardottt.github.io/pp-test",` +
				`"ScanURL":"https://edoardottt.github.io/pp-test?constructor.prototype.lfhfqn=lfhfqn` +
				`&__proto__[lfhfqn]=lfhfqn&constructor[prototype][lfhfqn]=lfhfqn&__proto__.lfhfqn=lfhfqn` +
				`&__proto__.lfhfqn=1|2|3&__proto__[lfhfqn]={\"json\":\"value\"}#__proto__[lfhfqn]=lfhfqn",` +
				`"JSEvaluation":"lfhfqn",` +
				`"Fingerprint":["jQuery"],` +
				`"ExploitURLs":["https://edoardottt.github.io/pp-test/?__proto__[url][]=data:,alert(1337)//` +
				`&__proto__[dataType]=script","https://edoardottt.github.io/pp-test/?__proto__[context]=` +
				`%3Cimg/src/onerror%3dalert(1337)%3E&__proto__[jquery]=x"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := output.FormatJSON(&tt.input)
			require.JSONEq(t, string(got), tt.want)
		})
	}
}
