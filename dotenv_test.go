package dotenv_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snobb/go-dotenv"
)

func TestLoadEnvFromReader(t *testing.T) {
	tests := []struct {
		name    string
		payload []string
		expect  map[string]string
		wantErr bool
	}{
		{
			name: "should load data from file and populate env",
			payload: []string{
				"  #doesnt = exist",
				"foo=bar",
				"fooSpace= bar",
				"fooQuote=\" baz \"",
				"baz=hello world",
				"bazSpace=\"hello world   \"",
				"bazEqual=hello=world",
			},
			expect: map[string]string{
				"foo":      "bar",
				"fooSpace": "bar",
				"fooQuote": " baz ",
				"baz":      "hello world",
				"bazSpace": "hello world   ",
				"bazEqual": "hello=world",
				"#doesnt":  "",
			},
		},
		{
			name: "should return an error on invalid records",
			payload: []string{
				"foo",
			},
			wantErr: true,
		},
		{
			name: "should return an error on unbalanced quotes",
			payload: []string{
				"foo=\"bar",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := bytes.NewReader([]byte(strings.Join(tt.payload, "\n")))

			err := dotenv.LoadEnvFromReader(rr)
			if tt.wantErr {
				assert.Error(t, err)
				t.Log(err)
				return
			}
			assert.NoError(t, err)

			for key, value := range tt.expect {
				varValue := os.Getenv(key)
				assert.Equal(t, value, varValue)
			}
		})
	}
}
