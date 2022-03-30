package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

const doubleQuoteSpecialChars = "\\\n\r\"!$`"

type DotEnv struct{}

func (de DotEnv) read(r io.Reader) (map[string]string, error) {
	return godotenv.Parse(r)
}

func (de DotEnv) write(w io.Writer, m map[string]string) error {
	content := make([]string, 0, len(m))
	for k, v := range m {
		content = append(content, fmt.Sprintf(`%s=%s`, k, doubleQuoteEscape(v)))
	}
	sort.Strings(content)
	for _, v := range content {
		_, err := fmt.Fprintf(w, "%s\n", v)
		if err != nil {
			return err
		}
	}
	return nil
}

func doubleQuoteEscape(line string) string {
	for _, c := range doubleQuoteSpecialChars {
		toReplace := "\\" + string(c)
		if c == '\n' {
			toReplace = `\n`
		}
		if c == '\r' {
			toReplace = `\r`
		}
		line = strings.Replace(line, string(c), toReplace, -1)
	}
	return line
}
