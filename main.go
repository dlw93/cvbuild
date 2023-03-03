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
	flag.StringVar(&args.InputFile, "input-file", "./index.html", "input file")
	flag.StringVar(&args.OutputDirectory, "output-directory", "./dist", "output directory")
	flag.StringVar(&args.ProjectRootAbsolute, "project-root", cwd, "absolute path to the project root directory")
}

func main() {
	flag.Parse()

	if info, err := os.Stat(args.InputFile); err != nil || info.Mode().IsDir() {
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
		options := BuildOptions{args.OutputDirectory, args.ProjectRootAbsolute}
		if result, err := Build(path, options); err != nil {
			return "", err
		} else {
			return result, nil
		}
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
