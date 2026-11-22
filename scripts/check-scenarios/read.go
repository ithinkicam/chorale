package main

import "os"

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path) //nolint:gosec // paths found under the OpenSpec folder
	return string(data), err
}
