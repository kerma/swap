package main

import (
	"fmt"
	"io"

	"github.com/joho/godotenv"
)

type DotEnv struct{}

func (de DotEnv) read(r io.Reader) (map[string]string, error) {
	return godotenv.Parse(r)
}

func (de DotEnv) write(w io.Writer, m map[string]string) error {
	content, err := godotenv.Marshal(m)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n", content)
	return err
}
