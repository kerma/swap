package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type reader func(io.Reader) (map[string]string, error)
type writer func(io.Writer, map[string]string) error

func main() {
	src := flag.String("src", "", "input filename (default: stdin)")
	dst := flag.String("dst", "", "output filename (default: stdout)")
	in := flag.String("in", "", "input type")
	out := flag.String("out", "debug", "output type")
	name := flag.String("name", "", "ouput name (if applicable)")
	flag.Parse()

	var err error
	r := os.Stdin
	if *src != "" {
		r, err = os.Open(*src)
		if err != nil {
			log.Fatal(err)
		}
	}

	w := os.Stdout
	if *dst != "" {
		w, err = os.Create(*dst)
		if err != nil {
			log.Fatal(err)
		}
	}

	read := getReader(*in)
	data, err := read(r)
	if err != nil {
		log.Fatal(err)
	}

	write := getWriter(getName(*src, *dst, *name), *out)
	err = write(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func getName(src, dst, name string) string {
	if name != "" {
		return name
	}

	if dst != "" {
		return strings.TrimSuffix(filepath.Base(dst), filepath.Ext(dst))
	}

	if src != "" {
		return strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))
	}

	return "swap-generated-name"
}

func getReader(kind string) (r reader) {
	switch kind {
	case "dotenv":
		r = new(DotEnv).read
	case "configmap":
		r = new(ConfigMap).read
	case "secret":
		r = new(Secret).read
	default:
		fatal("unknown input type", kind)
	}
	return
}

func getWriter(name, kind string) (w writer) {
	switch kind {
	case "dotenv":
		w = new(DotEnv).write
	case "configmap":
		cm := ConfigMap{name: name}
		w = cm.write
	case "secret":
		s := Secret{name: name}
		w = s.write
	default:
		w = writeDebug
	}
	return
}

func writeDebug(w io.Writer, m map[string]string) error {
	_, err := fmt.Fprintf(w, "%#v\n", m)
	return err
}

func fatal(a ...interface{}) {
	fmt.Print("Error: ")
	fmt.Println(a...)
	os.Exit(1)
}
