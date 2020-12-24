package buildconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

// Run executes the script builder
func Run() error {

	opts, err := parseInfo(viper.GetString("buildconfig"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid buildconfig: %v\n", err)
		os.Exit(1)
	}

	// TODO: make this experience nicer
	if _, err := os.Stat("configure"); err == nil {
		fmt.Fprintf(os.Stderr, "./configure: file exists. Overwrite? y/n\n")
		if !checkOverwrite() {
			return fmt.Errorf("./configure: file exists")
		}

		if _, err := os.Stat("Makefile"); err == nil && viper.GetBool("makefile") {
			fmt.Fprintf(os.Stderr, "Makefile: file exists. Overwrite? y/n\n")
			if !checkOverwrite() {
				return fmt.Errorf("Makefile: file exists")
			}
		}
		fmt.Println()
	}

	out, err := build(opts)
	if err != nil {
		return fmt.Errorf("error while templating ./configure: %v", err)
	}

	if viper.GetBool("makefile") {
		makefile, err := buildMakefile(opts)
		if err != nil {
			return fmt.Errorf("error while templating Makefile: %v", err)
		}
		err = ioutil.WriteFile("Makefile", makefile, 0644)
		if err != nil {
			return fmt.Errorf("error writing Makefile: %v", err)
		}
		fmt.Printf("Makefile created successfully\n")
	}

	err = ioutil.WriteFile("configure", out, 0755)
	if err != nil {
		return fmt.Errorf("error writing ./configure: %v", err)
	}
	fmt.Printf("./configure created successfully\n")

	return nil

}

func build(opts *options) ([]byte, error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	tmpl, err := template.New("script").Funcs(funcMap).Parse(scriptTmpl)

	if err != nil {
		return nil, err
	}
	var result bytes.Buffer

	tmpl.Execute(&result, opts.info)

	return result.Bytes(), nil
}

func buildMakefile(opts *options) ([]byte, error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	tmpl, err := template.New("makefile").Funcs(funcMap).Parse(makefileTmpl)

	if err != nil {
		return nil, err
	}
	var result bytes.Buffer

	tmpl.Execute(&result, opts.info)

	return result.Bytes(), nil
}

func checkOverwrite() bool {
	for {
		response := make([]byte, 1)
		os.Stdin.Read(response)
		if response[0] == 'y' {
			return true
		} else if response[0] == 'n' {
			return false
		}
	}
}
