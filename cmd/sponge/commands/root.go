// Package commands are subcommands of the sponge command.
package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/cmd/sponge/commands/generate"
)

var (
	version     = "v0.0.0"
	versionFile = GetSpongeDir() + "/.sponge/.github/version"
)

// NewRootCMD command entry
func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sponge",
		Long: fmt.Sprintf(`
A powerful and easy-to-use Go development framework that enables you to effortlessly 
build stable, reliable, and high-performance backend services with a "low-code" approach.
Repo: %s
Docs: %s`,
			color.HiCyanString("https://github.com/go-dev-frame/sponge"),
			color.HiCyanString("https://go-sponge.com")),
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       getVersion(),
	}

	cmd.AddCommand(
		InitCommand(),
		UpgradeCommand(),
		PluginsCommand(),
		GenWebCommand(),
		GenMicroCommand(),
		generate.ConfigCommand(),
		OpenUICommand(),
		MergeCommand(),
		PatchCommand(),
		GenGraphCommand(),
		TemplateCommand(),
		AssistantCommand(),
	)

	return cmd
}

func getVersion() string {
	data, _ := os.ReadFile(versionFile)
	v := string(data)
	if v != "" {
		return v
	}
	return "unknown, execute command \"sponge init\" to get version"
}

// GetSpongeDir get sponge home directory
func GetSpongeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("can't get home directory'")
		return ""
	}

	return dir
}
