package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/Laughs-In-Flowers/flip"
	"github.com/Laughs-In-Flowers/log"
	"github.com/Laughs-In-Flowers/valhalla/lib/viewer"
)

type Options struct {
	Path, Formatter string
	log.Logger
}

func defaultOptions() *Options {
	currentDir, _ := os.Getwd()
	currentPath = filepath.Join(currentDir, "image.jpg")
	l := log.New(os.Stdout, log.LInfo, log.DefaultNullFormatter())
	return &Options{
		currentPath, "null", l,
	}
}

func tFlags(fs *flip.FlagSet, o *Options) *flip.FlagSet {
	fs.StringVar(&o.Formatter, "formatter", o.Formatter, "Specify the log formatter.")
	return fs
}

type Execute func(*Options) error

type Executing interface {
	Run(*Options) error
}

type executing struct {
	has []Execute
}

func NewExecuting(e ...Execute) *executing {
	return &executing{
		e,
	}
}

func (e *executing) Run(o *Options) error {
	for _, v := range e.has {
		err := v(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func xFormatter(o *Options) error {
	if o.Formatter != "null" {
		switch o.Formatter {
		case "raw":
			o.SwapFormatter(log.GetFormatter("raw"))
		case "text", "stdout":
			o.SwapFormatter(log.GetFormatter("valhalla_text"))
		}
	}
	return nil
}

var topExecute = NewExecuting(xFormatter)

func TopCommand(o *Options) flip.Command {
	fs := flip.NewFlagSet("top", flip.ContinueOnError)
	fs = tFlags(fs, o)
	return flip.NewCommand(
		"",
		"valhalla",
		"image viewer top level flag use",
		1,
		func(c context.Context, a []string) flip.ExitStatus {
			topExecute.Run(o)
			return flip.ExitNo
		},
		fs,
	)
}

func vFlags(fs *flip.FlagSet, o *Options) *flip.FlagSet {
	fs.StringVar(&o.Path, "path", o.Path, "The path of the image to view.")
	return fs
}

func basicErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", e)
		os.Exit(-1)
	}
}

func ViewCommand(o *Options) flip.Command {
	fs := flip.NewFlagSet("view", flip.ContinueOnError)
	fs = vFlags(fs, o)

	return flip.NewCommand(
		"",
		"view",
		"image viewer view command.",
		1,
		func(c context.Context, a []string) flip.ExitStatus {
			v, err := viewer.New(o.Path, o.Logger)
			if err != nil {
				basicErr(err)
				return flip.ExitFailure
			}
			v.Run()
			return flip.ExitSuccess
		},
		fs,
	)
}

var (
	versionPackage string = path.Base(os.Args[0])
	versionTag     string = "No Tag"
	versionHash    string = "No Hash"
	versionDate    string = "No Date"
)

var (
	options     *Options
	currentPath string
	C           *flip.Commander
)

func init() {
	options = defaultOptions()
	log.SetFormatter("valhalla_text", log.MakeTextFormatter(versionPackage))
	C = flip.BaseWithVersion(versionPackage, versionTag, versionHash, versionDate)
	C.RegisterGroup("top", 1, TopCommand(options))
	C.RegisterGroup("view", 10, ViewCommand(options))
}

func main() {
	ctx := context.Background()
	C.Execute(ctx, os.Args)
	os.Exit(0)
}
