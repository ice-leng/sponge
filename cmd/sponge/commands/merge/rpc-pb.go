package merge

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// GRPCServiceCode merge the grpc service code
func GRPCServiceCode() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "rpc-pb",
		Short: "Merge the generated grpc related code into the template file",
		Long:  "Merge the generated grpc related code into the template file.",
		Example: color.HiBlackString(`  # Merge go template file in local server directory
  sponge merge rpc-pb

  # Merge go template file in specified directory
  sponge merge rpc-pb --dir=/path/to/server/directory`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			codeTypes := []mergeType{errCodeType, serviceGRPCTmplType, serviceGRPCClientType}
			for _, codeType := range codeTypes {
				if err := newMergeParams(dir, codeType).runMerge(); err != nil {
					fmt.Printf("[Warning] %v\n", err)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "input directory")

	return cmd
}
