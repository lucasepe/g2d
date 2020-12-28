package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	banner = `       ___  ___
 __ _ |_  )|   \
/ _' | / / | |) |
\__, |/___||___/    
|___/  Scripting Interpreter`

	appSummary = "Create beautiful drawings using a simple scripting language."
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Use:                   fmt.Sprintf("%s <COMMAND>", appName()),
		Short:                 appSummary,
		Long:                  banner,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}} - Luca Sepe <luca.sepe@gmail.com>
`)
}

func appName() string {
	return filepath.Base(os.Args[0])
}
