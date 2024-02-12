/*
pphack - The Most Advanced Client-Side Prototype Pollution Scanner

This repository is under MIT License https://github.com/edoardottt/pphack/blob/main/LICENSE
*/

package output_test

import (
	"reflect"
	"testing"

	"github.com/edoardottt/pphack/pkg/output"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name  string
		input output.JSONData
		want  string
	}{
		{name: "no error",
			input: output.JSONData{
				URL:          "https://edoardottt.github.io/pp-test?constructor.prototype.oigczu=oigczu\u0026__proto__[oigczu]=oigczu\u0026constructor[prototype][oigczu]=oigczu\u0026__proto__.oigczu=oigczu\u0026__proto__.oigczu=1|2|3\u0026__proto__[oigczu]={\"json\":\"value\"}#__proto__[oigczu]=oigczu",
				JSEvaluation: "oigczu",
				Error:        "",
			},
			want: `{"URL":"https://edoardottt.github.io/pp-test?constructor.prototype.oigczu=oigczu\u0026__proto__[oigczu]=oigczu\u0026constructor[prototype][oigczu]=oigczu\u0026__proto__.oigczu=oigczu\u0026__proto__.oigczu=1|2|3\u0026__proto__[oigczu]={\"json\":\"value\"}#__proto__[oigczu]=oigczu","JSEvaluation":"oigczu"}`, //nolint:lll
		},
		{name: "error",
			input: output.JSONData{
				URL:          "https://edoardottt.github.io/pp-tes?constructor.prototype.sqtiwx=sqtiwx\u0026__proto__[sqtiwx]=sqtiwx\u0026constructor[prototype][sqtiwx]=sqtiwx\u0026__proto__.sqtiwx=sqtiwx\u0026__proto__.sqtiwx=1|2|3\u0026__proto__[sqtiwx]={\"json\":\"value\"}#__proto__[sqtiwx]=sqtiwx",
				JSEvaluation: "",
				Error:        "encountered an undefined value",
			},
			want: `{"URL":"https://edoardottt.github.io/pp-tes?constructor.prototype.sqtiwx=sqtiwx\u0026__proto__[sqtiwx]=sqtiwx\u0026constructor[prototype][sqtiwx]=sqtiwx\u0026__proto__.sqtiwx=sqtiwx\u0026__proto__.sqtiwx=1|2|3\u0026__proto__[sqtiwx]={\"json\":\"value\"}#__proto__[sqtiwx]=sqtiwx","Error":"encountered an undefined value"}`, //nolint:lll
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := output.FormatJSON(tt.input.URL, tt.input.JSEvaluation, tt.input.Error); !reflect.DeepEqual(string(got), tt.want) { //nolint:lll
				t.Errorf("GetJSONString\n%v", string(got))
				t.Errorf("want\n%v", tt.want)
			}
		})
	}
}
