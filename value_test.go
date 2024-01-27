package jsond

import (
	"encoding/json"
	"testing"
)

func TestGetTypeString(t *testing.T) {

	tests := []struct {
		src  string
		want string
	}{
		{src: `true`, want: "bool"},
		{src: `1`, want: "number"},
		{src: `"aa"`, want: "string"},
		{src: `[1,2,3]`, want: "array"},
		{src: `{"key":"value"}`, want: "object"},
		{src: `null`, want: "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.src, func(t *testing.T) {
			var v = *new(any)

			err := json.Unmarshal([]byte(tt.src), &v)
			if err != nil {
				t.Fatal(err)
			}

			got := getTypeString(v)
			if got != tt.want {
				t.Errorf("\ngot  %s\nwant %s", got, tt.want)
			}
		})
	}
}
