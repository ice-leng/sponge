package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/gobash"
)

var pluginNames = []string{
	"go",
	"protoc",
	"protoc-gen-go",
	"protoc-gen-go-grpc",
	"protoc-gen-validate",
	"protoc-gen-gotag",
	"protoc-gen-go-gin",
	"protoc-gen-go-rpc-tmpl",
	"protoc-gen-openapiv2",
	"protoc-gen-doc",
	"swag",
	//"golangci-lint",
	//"go-callvis",
}

var installPluginCommands = map[string]string{
	"go":                     "go: please install manually yourself, download url is https://go.dev/dl/ or https://golang.google.cn/dl/",
	"protoc":                 "protoc: please install manually yourself, download url is https://github.com/protocolbuffers/protobuf/releases/tag/v25.2",
	"protoc-gen-go":          "google.golang.org/protobuf/cmd/protoc-gen-go@latest",
	"protoc-gen-go-grpc":     "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
	"protoc-gen-validate":    "github.com/envoyproxy/protoc-gen-validate@latest",
	"protoc-gen-gotag":       "github.com/srikrsna/protoc-gen-gotag@latest",
	"protoc-gen-go-gin":      "github.com/zhufuyi/sponge/cmd/protoc-gen-go-gin@latest",
	"protoc-gen-go-rpc-tmpl": "github.com/zhufuyi/sponge/cmd/protoc-gen-go-rpc-tmpl@latest",
	"protoc-gen-openapiv2":   "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest",
	"protoc-gen-doc":         "github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest",
	"swag":                   "github.com/swaggo/swag/cmd/swag@v1.8.12",
	//"golangci-lint":          "github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	//"go-callvis":             "github.com/ofabry/go-callvis@latest",
}

const (
	installedSymbol = "✔ "
	lackSymbol      = "❌ "
	warnSymbol      = "⚠ "
)

// PluginsCommand plugins management
func PluginsCommand() *cobra.Command {
	var installFlag bool
	var skipPluginName string

	cmd := &cobra.Command{
		Use:   "plugins",
		Short: "Managing sponge dependency plugins",
		Long: color.HiBlackString(`managing sponge dependency plugins.

Examples:
  # show all dependency plugins.
  sponge plugins

  # install all dependency plugins.
  sponge plugins --install

  # skip installing dependency plugins, multiple plugin names separated by commas
  sponge plugins --install --skip=go-callvis
`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			installedNames, lackNames := checkInstallPlugins()
			lackNames = filterLackNames(lackNames, skipPluginName)
			if installFlag {
				installPlugins(lackNames)
			} else {
				showDependencyPlugins(installedNames, lackNames)
			}

			return nil
		},
	}
	cmd.Flags().BoolVarP(&installFlag, "install", "i", false, "install dependency plugins")
	cmd.Flags().StringVarP(&skipPluginName, "skip", "s", "", "skip installing dependency plugins")

	return cmd
}

func checkInstallPlugins() ([]string, []string) {
	var installedNames, lackNames = []string{}, []string{}
	for _, name := range pluginNames {
		_, err := gobash.Exec("which", name)
		if err != nil {
			lackNames = append(lackNames, name)
			continue
		}
		installedNames = append(installedNames, name)
	}

	data, _ := os.ReadFile(versionFile)
	v := string(data)
	if v != "" {
		version = v
	}

	return installedNames, lackNames
}

func showDependencyPlugins(installedNames []string, lackNames []string) {
	var content string

	if len(installedNames) > 0 {
		content = "Installed dependency plugins:\n"
		for _, name := range installedNames {
			content += "    " + installedSymbol + " " + name + "\n"
		}
	}

	if len(lackNames) > 0 {
		content += "\nUninstalled dependency plugins:\n"
		for _, name := range lackNames {
			content += "    " + lackSymbol + " " + name + "\n"
		}
		content += "\nInstalling dependency plugins using the command: sponge plugins --install\n"
	} else {
		content += "\nAll dependency plugins installed.\n"
	}

	fmt.Println(content)
}

func installPlugins(lackNames []string) {
	if len(lackNames) == 0 {
		fmt.Printf("\n    All dependency plugins installed.\n\n")
		return
	}
	fmt.Printf("install a total of %d dependency plugins, need to wait a little time.\n\n", len(lackNames))

	var wg = &sync.WaitGroup{}
	var manuallyNames []string
	for _, name := range lackNames {
		if name == "go" || name == "protoc" {
			manuallyNames = append(manuallyNames, name)
			continue
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			ctx, _ := context.WithTimeout(context.Background(), time.Minute*3) //nolint
			pkgAddr, ok := installPluginCommands[name]
			if !ok {
				return
			}
			pkgAddr = adaptInternalCommand(name, pkgAddr)
			result := gobash.Run(ctx, "go", "install", pkgAddr)
			for v := range result.StdOut {
				_ = v
			}
			if result.Err != nil {
				fmt.Printf("%s %s, %v\n", lackSymbol, name, result.Err)
			} else {
				fmt.Printf("%s %s\n", installedSymbol, name)
			}
		}(name)
	}

	wg.Wait()

	for _, name := range manuallyNames {
		fmt.Println(warnSymbol + " " + installPluginCommands[name])
	}
	fmt.Println()
}

func adaptInternalCommand(name string, pkgAddr string) string {
	if name == "protoc-gen-go-gin" || name == "protoc-gen-go-rpc-tmpl" {
		if version != "v0.0.0" {
			return strings.ReplaceAll(pkgAddr, "@latest", "@"+version)
		}
	}

	return pkgAddr
}

func filterLackNames(lackNames []string, skipPluginName string) []string {
	if skipPluginName == "" {
		return lackNames
	}
	skipPluginNames := strings.Split(skipPluginName, ",")

	names := []string{}
	for _, name := range lackNames {
		isMatch := false
		for _, pluginName := range skipPluginNames {
			if name == pluginName {
				isMatch = true
				continue
			}
		}
		if !isMatch {
			names = append(names, name)
		}
	}
	return names
}
