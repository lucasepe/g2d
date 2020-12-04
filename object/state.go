package object

import (
	"io"
)

var (
	WorkDir        string
	SourceFile     string
	SaveCounter    int
	Arguments      []string
	StandardInput  io.Reader
	StandardOutput io.Writer
	ExitFunction   func(int)
)
