// Package dotenv populates env variables from file.
package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const defaultFile = ".env"

// Options contain the options and defaults for the Env loader.
var Options = struct {
	OverrideExisting bool
}{}

// LoadEnvFromFile loads the specified file with env variables (same syntax as in shell) and
// populates the env variables for the running process.
func LoadEnvFromFile(path string) (err error) {
	var fh *os.File

	fh, err = os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer func() {
		err = fh.Close()
	}()

	return LoadEnvFromReader(fh)
}

// LoadEnv loads the default file (.env) with env variables (same syntax as in shell) and
// populates the env variables for the running process.
func LoadEnv() error {
	return LoadEnvFromFile(defaultFile)
}

// LoadEnvFromReader loads the env variables pairs (same syntax as in shell) from the provided
// reader and populates the env variables for the running process.
func LoadEnvFromReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line[0] == '#' {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid line: %s", line)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if value[0] == '"' {
			if value[len(value)-1] != '"' {
				return fmt.Errorf("unbalanced quotes: '%s'", value)
			}

			value = value[1 : len(value)-1]
		}

		if os.Getenv("DEBUG") != "" {
			fmt.Printf("dotenv: %s=%s\n", key, value)
		}

		if !Options.OverrideExisting && os.Getenv(key) != "" {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
