package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/gargath/buildconfig/pkg/buildconfig"
)

func main() {

	flag.Bool("help", false, "print this help and exit")
	flag.Bool("version", false, "print version and exit")
	flag.StringP("buildconfig", "f", "buildconfig.yaml", "path to the buildconfig file to use")
	flag.BoolP("makefile", "m", false, "also write Makefile template")

	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	if viper.GetBool("help") {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, flag.CommandLine.FlagUsages())
		os.Exit(0)
	}

	if viper.GetBool("version") {
		fmt.Fprintf(os.Stderr, "buildconfig %s:\n", version())
		os.Exit(0)
	}

	err := buildconfig.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
		os.Exit(1)
	}
}
