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
		name             string
		preset           map[string]string
		overrideExisting bool
		payload          []string
		expect           map[string]string
		wantErr          bool
	}{
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
		{
			name:             "should load data from file and populate env and override existing",
			overrideExisting: true,
			preset: map[string]string{
				"existing": "foo",
			},
			payload: []string{
				"DEBUG=1",
				"  #doesnt = exist",
				"foo=bar",
				"existing=baz",
			},
			expect: map[string]string{
				"foo":      "bar",
				"existing": "baz",
			},
		},
		{
			name: "should load data from file and populate env",
			preset: map[string]string{
				"existing": "foo",
			},
			payload: []string{
				"DEBUG=1",
				"  #doesnt = exist",
				"foo=bar",
				"fooSpace= bar",
				"fooQuote=\" baz \"",
				"baz=hello world",
				"bazSpace=\"hello world   \"",
				"bazEqual=hello=world",
				"existing=baz",
			},
			expect: map[string]string{
				"foo":      "bar",
				"fooSpace": "bar",
				"fooQuote": " baz ",
				"baz":      "hello world",
				"bazSpace": "hello world   ",
				"bazEqual": "hello=world",
				"#doesnt":  "",
				"existing": "foo", // not overriden
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := bytes.NewReader([]byte(strings.Join(tt.payload, "\n")))

			for k, v := range tt.preset {
				t.Setenv(k, v)
			}

			if tt.overrideExisting {
				dotenv.Options.OverrideExisting = true
				defer func() {
					dotenv.Options.OverrideExisting = false
				}()
			}

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
				os.Unsetenv(key)
			}
		})
	}
}
