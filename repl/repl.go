package repl

// Package repl implements the Read-Eval-Print-Loop or interactive console
// by lexing, parsing and evaluating the input in the interpreter

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/g2d/eval"
	"github.com/lucasepe/g2d/lexer"
	"github.com/lucasepe/g2d/object"
	"github.com/lucasepe/g2d/parser"
)

// prompt is the REPL prompt displayed for each input
const prompt = "▶▶ "

const banner = `Hello {{USER}}!
              ___  ___
Welcome __ _ |_  )|   \
    to / _' | / / | |) |
       \__, |/___||___/    
       |___/  Programming language.`

// Options are the REPL parameters.
type Options struct {
	Debug       bool
	Interactive bool
}

// REPL is the read-evaluate-print loop.
type REPL struct {
	user string
	args []string
	opts *Options
}

// New create a REPL
func New(user string, args []string, opts *Options) *REPL {
	object.StandardInput = os.Stdin
	object.StandardOutput = os.Stdout
	object.ExitFunction = os.Exit

	return &REPL{user, args, opts}
}

// Eval parses and evalulates the program given by f and returns the resulting
// environment, any errors are printed to stderr
func (r *REPL) Eval(f io.Reader) (env *object.Environment) {
	env = object.NewEnvironment()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source file: %s", err)
		return
	}

	l := lexer.New(string(b))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		return
	}

	// if obj := eval.Eval(program, env); obj.Type() == object.ERROR {//
	if obj := eval.BeginEval(program, env, l); (obj != nil) && (obj.Type() == object.ERROR) {
		printParserErrors(os.Stderr, []string{obj.String()})
	}
	return
}

// StartEvalLoop starts the REPL in a continious eval loop
func (r *REPL) StartEvalLoop(in io.Reader, out io.Writer, env *object.Environment) {
	scanner := bufio.NewScanner(in)

	if env == nil {
		env = object.NewEnvironment()
	}

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		obj := eval.Eval(program, env)
		if obj != nil {
			if _, ok := obj.(*object.Null); !ok {
				io.WriteString(out, obj.Inspect())
				io.WriteString(out, "\n")
			}
		}
	}
}

// Run execute the read-eval loop.
func (r *REPL) Run() {
	object.Arguments = make([]string, len(r.args))
	copy(object.Arguments, r.args)

	if len(r.args) == 0 {
		welcome(r.user)
		r.StartEvalLoop(os.Stdin, os.Stdout, nil)
		return
	}

	if len(r.args) > 0 {
		workDir, err := filepath.Abs(filepath.Dir(r.args[0]))
		if err != nil {
			workDir = filepath.Dir(r.args[0])
		}

		sourceFile := filepath.Base(r.args[0])

		object.WorkDir = workDir
		object.SourceFile = strings.TrimSuffix(sourceFile, filepath.Ext(sourceFile))

		f, err := os.Open(r.args[0])
		if err != nil {
			log.Fatalf("could not open source file %s: %s", r.args[0], err)
		}

		// Remove program argument (zero)
		r.args = r.args[1:]
		object.Arguments = object.Arguments[1:]

		env := r.Eval(f)
		if r.opts.Interactive {
			r.StartEvalLoop(os.Stdin, os.Stdout, env)
		}

	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, fmt.Sprintf(" - %s\n", msg))
	}
}

func welcome(user string) {
	str := strings.Replace(banner, "{{USER}}", user, 1)
	fmt.Print(str, "\n\n")
	fmt.Println("Feel free to type in commands.")
}
