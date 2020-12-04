package main

// Package main implements the main process which executes a program if
// a filename is supplied as an argument or invokes the interpreter's
// REPL and waits for user input before lexing, parsing nad evaulating.

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/lucasepe/g2d/repl"
)

var (
	// Version release version
	Version = "0.0.1"

	// GitCommit will be overwritten automatically by the build system
	GitCommit = "HEAD"
)

var (
	interactive bool
	version     bool
	debug       bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&version, "v", false, "display version information")
	flag.BoolVar(&debug, "d", false, "enable debug mode")

	flag.BoolVar(&interactive, "i", false, "enable interactive mode")
}

// Indent indents a block of text with an indent string
func Indent(text, indent string) string {
	if text[len(text)-1:] == "\n" {
		result := ""
		for _, j := range strings.Split(text[:len(text)-1], "\n") {
			result += indent + j + "\n"
		}
		return result
	}
	result := ""
	for _, j := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		result += indent + j + "\n"
	}
	return result[:len(result)-1]
}

// FullVersion returns the full version and commit hash
func FullVersion() string {
	return fmt.Sprintf("%s@%s", Version, GitCommit)
}

func main() {
	flag.Parse()

	if version {
		fmt.Printf("%s %s", path.Base(os.Args[0]), FullVersion())
		os.Exit(0)
	}

	user, err := user.Current()
	if err != nil {
		log.Fatalf("could not determine current user: %s", err)
	}

	args := flag.Args()

	opts := &repl.Options{
		Debug:       debug,
		Interactive: interactive,
	}
	repl := repl.New(user.Username, args, opts)
	repl.Run()
}
