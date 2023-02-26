package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var args struct {
	InputFile           string
	OutputDirectory     string
	ProjectRootAbsolute string
}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	strvar := func(p *string, name string, short rune, value string, usage string) {
		flag.StringVar(p, name, value, usage)
		flag.StringVar(p, string(short), value, usage+" (shorthand)")
	}
	strvar(&args.InputFile, "input-file", 'i', "./index.html", "input file")
	strvar(&args.OutputDirectory, "output-directory", 'o', "./dist", "output directory")
	strvar(&args.ProjectRootAbsolute, "project-root", 'p', cwd, "absolute path to the project root directory")
}

func main() {
	flag.Parse()

	if info, err := os.Stat(args.InputFile); err != nil {
		log.Fatalf("entry point %s does not exist", args.InputFile)
	} else if info.IsDir() {
		log.Fatalf("entry point %s is a directory", args.InputFile)
	}

	if info, err := os.Stat(args.ProjectRootAbsolute); err != nil {
		log.Fatalf("project root %s does not exist", args.ProjectRootAbsolute)
	} else if !info.IsDir() {
		log.Fatalf("project root %s is not a directory", args.ProjectRootAbsolute)
	}

	file, err := os.Open(args.InputFile)
	if err != nil {
		log.Fatalf("failed to open entry point %s: %s", args.InputFile, err)
	}
	defer file.Close()

	doc, err := NewDocument(file)
	if err != nil {
		log.Fatalf("failed to parse %s: %s", args.InputFile, err)
	}

	err = doc.Walk(func(path string) (string, error) {
		result, err := Build(path, args.OutputDirectory)
		if err != nil {
			return "", err
		}
		return result[0], nil
	})
	if err != nil {
		log.Fatalf("failed to process dependency in %s: %s", args.InputFile, err)
	}

	outpath := filepath.Join(args.OutputDirectory, filepath.Base(args.InputFile))
	outfile, err := os.Create(outpath)
	if err != nil {
		log.Fatalf("failed to create %s: %s", outpath, err)
	}
	defer outfile.Close()

	err = doc.WriteTo(outfile)
	if err != nil {
		log.Fatalf("failed to write to %s: %s", outpath, err)
	}
}
