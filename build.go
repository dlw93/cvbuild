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
	InputFile           string
	OutputDirectory     string
	ProjectRootAbsolute string
}

func ImportMetaPlugin() api.Plugin {
	quoteRE := func(q rune) string {
		return fmt.Sprintf(`%c[^%c]+%c`, q, q, q)
	}
	literalRE := strings.Join(Map([]rune{'\'', '"', '`'}, quoteRE), "|")
	importMetaRE := regexp.MustCompile(`(?m)\bnew\s+URL\s*\(\s*(` + literalRE + `)\s*,\s*import\.meta\.url\s*(?:,\s*)?\)`)

	return api.Plugin{
		Name: "import-meta",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(api.OnLoadOptions{Filter: ".js$"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				b, err := os.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				s := string(b)

				for _, match := range importMetaRE.FindAllStringSubmatch(s, -1) {
					str := match[1]
					path := str[1 : len(str)-1]
					fmt.Println("import-meta", path)
				}

				return api.OnLoadResult{
					Contents:   &s,
					ResolveDir: filepath.Dir(args.Path),
					Loader:     api.LoaderJS,
				}, nil
			})
		},
	}
}

func AbsolutePathPlugin(root string) api.Plugin {
	return api.Plugin{
		Name: "absolute-path",
		Setup: func(build api.PluginBuild) {
			build.OnResolve(api.OnResolveOptions{Filter: "^/"},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					// TODO filepath.Rel()
					rel := fmt.Sprintf(".%s", args.Path)
					options := api.ResolveOptions{
						ResolveDir: root,
						Kind:       args.Kind,
						Importer:   args.Importer,
						Namespace:  args.Namespace,
					}
					result := build.Resolve(rel, options)
					if len(result.Errors) > 0 {
						err := fmt.Errorf("failed to resolve %s in %s: %s", args.Path, root, result.Errors[0].Text)
						return api.OnResolveResult{}, err
					}
					return api.OnResolveResult{
						Path:      result.Path,
						Namespace: args.Namespace,
					}, nil
				})
		},
	}
}

func Build(entryFile string, outputDir string) ([]string, error) {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entryFile},

		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,

		Outdir:     outputDir,
		Write:      true,
		Format:     api.FormatESModule,
		EntryNames: "[name]",
		AssetNames: "[name]",

		Plugins: []api.Plugin{
			AbsolutePathPlugin("./"),
			ImportMetaPlugin(),
		},
	})

	if len(result.Errors) > 0 {
		errs := Map(result.Errors, func(err api.Message) error {
			return fmt.Errorf("%s:%d: %s", err.Location.File, err.Location.Line, err.Text)
		})
		return nil, errors.Join(errs...)
	}

	outfiles := Map(result.OutputFiles, func(f api.OutputFile) string { return f.Path })

	return outfiles, nil
}
