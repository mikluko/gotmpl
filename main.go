package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func stdin() (data interface{}, err error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return
	}
	err = json.NewDecoder(os.Stdin).Decode(&data)
	return
}

func main() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		_, _ = fmt.Fprintln(flag.CommandLine.Output(),
			"Usage: gotmpl template [template ...] < input.json > output.txt")
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	data, err := stdin()
	if err != nil {
		log.Fatal(err)
	}

	name := filepath.Base(flag.Arg(0))
	tmpl, err := template.New(name).Funcs(sprig.TxtFuncMap()).ParseFiles(flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		log.Fatal(err)
	}
}
