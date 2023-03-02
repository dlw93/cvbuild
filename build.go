package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type BuildOptions struct {
	OutputDirectory     string
	ProjectRootAbsolute string
}

type messageError api.Message

func ImportMetaUrlPlugin() api.Plugin {
	jsStringLiteralRE := func(qmarks ...rune) string {
		return strings.Join(Map(qmarks, func(q rune) string { return fmt.Sprintf(`%c%c%c`, q, q, q) }), "|")
	}
	importMetaUrlRE := regexp.MustCompile(`(?m)\bnew\s+URL\s*\(\s*(` + jsStringLiteralRE('\'', '"', '`') + `)\s*,\s*import\.meta\.url\s*(?:,\s*)?\)`)
	deps := func(text string) []string {
		matches := importMetaUrlRE.FindAllStringSubmatchIndex(text, -1)
		return Map(matches, func(match []int) string {
			return text[match[2]+1 : match[3]-1] // +1/-1 to remove quotes
		})
	}

	return api.Plugin{
		Name: "import-meta-url",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(api.OnLoadOptions{Filter: ".js$"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				b, err := os.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}
				s := string(b)
				return api.OnLoadResult{
					Contents:   &s,
					ResolveDir: filepath.Dir(args.Path),
					Loader:     api.LoaderJS,
					WatchFiles: deps(s),
				}, nil
			})
		},
	}
}

func AbsolutePathPlugin() api.Plugin {
	return api.Plugin{
		Name: "absolute-path",
		Setup: func(build api.PluginBuild) {
			root := filepath.Clean(build.InitialOptions.AbsWorkingDir)

			build.OnResolve(api.OnResolveOptions{Filter: "^/"},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					rel, err := filepath.Rel(args.ResolveDir, filepath.Join(root, args.Path))
					if err != nil {
						err := fmt.Errorf("failed to resolve %s in %s: %w", args.Path, root, err)
						return api.OnResolveResult{}, err
					}

					path := "./" + filepath.ToSlash(rel)
					options := api.ResolveOptions{
						ResolveDir: args.ResolveDir,
						Kind:       args.Kind,
						Importer:   args.Importer,
						Namespace:  args.Namespace,
					}

					result := build.Resolve(path, options)
					errs := errors.Join(Map(result.Errors, newMessageError)...)

					return api.OnResolveResult{
						Path:      result.Path,
						Namespace: args.Namespace,
					}, errs
				})
		},
	}
}

func Build(input string, options BuildOptions) ([]string, error) {
	opts := api.BuildOptions{
		EntryPoints:   []string{input},
		AbsWorkingDir: options.ProjectRootAbsolute,

		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,

		Outdir:     options.OutputDirectory,
		Write:      true,
		Format:     api.FormatESModule,
		EntryNames: "[name]",
		AssetNames: "[name]",

		Plugins: []api.Plugin{
			AbsolutePathPlugin(),
			ImportMetaUrlPlugin(),
		},
	}

	ctx, err := api.Context(opts)
	if err != nil {
		return nil, errors.Join(Map(err.Errors, newMessageError)...)
	}
	defer ctx.Dispose()

	result := ctx.Rebuild()

	errs := errors.Join(Map(result.Errors, newMessageError)...)
	outfiles := Map(result.OutputFiles, func(f api.OutputFile) string { return f.Path })

	return outfiles, errs
}

func newMessageError(msg api.Message) error {
	return messageError(msg)
}

func (msg messageError) Error() string {
	if msg.Location != nil {
		return fmt.Sprintf("%s [Ln %d, Col %d]: %s", msg.Location.File, msg.Location.Line, msg.Location.Column, msg.Text)
	} else {
		return msg.Text
	}
}
