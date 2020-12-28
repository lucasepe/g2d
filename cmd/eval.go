package cmd

import (
	"fmt"
	"image"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/lucasepe/g2d/data"
	"github.com/lucasepe/g2d/eval"
	"github.com/lucasepe/g2d/gg/img"
	"github.com/lucasepe/g2d/lexer"
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/parser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	limit = 512 * 1000

	optDirectory = "directory"
	optPrefix    = "prefix"
)

// renderCmd represents the render command
var evalCmd = &cobra.Command{
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Use:                   "eval <script URL or PATH>",
	Short:                 "Evaluate a g2d script",
	Example:               evalCmdExample(),
	Run: func(cmd *cobra.Command, args []string) {
		src, err := data.Fetch(args[0], limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}

		directory, err := cmd.Flags().GetString(optDirectory)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}

		prefix, err := lastPathSegment(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}

		if err := doEval(src, directory, prefix); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}
	},
}

func init() {

	evalCmd.Flags().StringP(optDirectory, "d", "", "snapshots destination folder (note that must exist)")
	//evalCmd.MarkFlagRequired(optDirectory)

	rootCmd.AddCommand(evalCmd)
}

func evalCmdExample() string {
	tpl := `  {{APP}} eval https://github.com/lucasepe/g2d/_examples/circles.g2d
  {{APP}} eval /path/to/my_script.g2d`

	return strings.Replace(tpl, "{{APP}}", appName(), -1)
}

// Eval parses and evalulates the program given by f and returns the resulting
// environment, any errors are printed to stderr
func doEval(src []byte, directory, prefix string) error {
	ctx := img.NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, 1024, 1024)))
	env := object.NewEnvironment(ctx,
		object.WithOutputDir(directory),
		object.WithSnapshotPrefix(prefix))

	l := lexer.New(string(src))
	p := parser.New(l)

	program := p.ParseProgram()

	n := len(p.Errors())

	if n != 0 {
		err := errors.New(p.Errors()[0])
		if n > 1 {
			for i := 1; i < n; i++ {
				err = errors.Wrapf(err, "\n%s", p.Errors()[i])
			}
		}
		return err
	}

	// if obj := eval.Eval(program, env); obj.Type() == object.ERROR {//
	if obj := eval.BeginEval(program, env, l); (obj != nil) && (obj.Type() == object.ERROR) {
		return errors.New(obj.String())
	}

	return nil
}

func lastPathSegment(uri string) (string, error) {
	var res string
	if strings.HasPrefix(uri, "http") {
		u, err := url.Parse(uri)
		if err != nil {
			return "", err
		}
		res = path.Base(u.Path)
	} else {
		res = filepath.Base(uri)
	}

	return strings.TrimSuffix(res, filepath.Ext(res)), nil
}
