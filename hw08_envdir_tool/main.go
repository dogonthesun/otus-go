package main

import (
	"fmt"
	"os"
)

var progname string

func init() {
	progname = os.Args[0]
}

func printUsage() {
	fmt.Printf("%[1]s: usage: %[1]s dir child\n", progname)
}

const FailExitCode = 111

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(FailExitCode)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("%s: %v", progname, err)
		os.Exit(FailExitCode)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
