package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/dlw93/cvbuild/util"
)

type BuildOptions struct {
	OutputDirectory     string
	ProjectRootAbsolute string
}

type BuildError api.Message

func ImportMetaUrlPlugin() api.Plugin {
	return api.Plugin{
		Name: "import-meta-url",
		Setup: func(build api.PluginBuild) {
			build.OnLoad(api.OnLoadOptions{Filter: ".js$"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				b, err := os.ReadFile(args.Path)
				if err != nil {
					return api.OnLoadResult{}, err
				}
				s := NewDependencyScanner(b)
				contents := s.String()
				deps := s.Scan()
				return api.OnLoadResult{
					Contents:   &contents,
					ResolveDir: filepath.Dir(args.Path),
					Loader:     api.LoaderJS,
					WatchFiles: deps,
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

			build.OnResolve(api.OnResolveOptions{Filter: "^/"}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
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
				errs := newBuildError(result.Errors)

				return api.OnResolveResult{
					Path:      result.Path,
					Namespace: args.Namespace,
				}, errs
			})
		},
	}
}

func Build(input string, options BuildOptions) (string, error) {
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
	result := api.Build(opts)

	if len(result.Errors) > 0 {
		return "", newBuildError(result.Errors)
	} else if len(result.OutputFiles) == 0 {
		return "", fmt.Errorf("no output generated for %s", input)
	}

	return result.OutputFiles[0].Path, nil
}

func newBuildError[T api.Message | []api.Message](msg T) error {
	switch msg := any(msg).(type) {
	case api.Message:
		return BuildError(msg)
	case []api.Message:
		return errors.Join(util.Map(msg, newBuildError[api.Message])...)
	default:
		panic("unreachable")
	}
}

func (msg BuildError) Error() string {
	if msg.Location != nil {
		return fmt.Sprintf("%s [Ln %d, Col %d]: %s", msg.Location.File, msg.Location.Line, msg.Location.Column, msg.Text)
	} else {
		return msg.Text
	}
}
