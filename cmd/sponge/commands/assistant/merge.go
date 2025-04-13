package assistant

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/pkg/goast"
	"github.com/go-dev-frame/sponge/pkg/gofile"
)

// MergeAssistantCode merge AI assistant generated code into source Go file
func MergeAssistantCode() *cobra.Command {
	var (
		dir           string
		files         []string // specified Go files
		assistantType string
		isClean       bool
	)

	cmd := &cobra.Command{
		Use:   "merge",
		Short: "Merge AI assistant generated code into source Go file",
		Long:  "Merge AI assistant generated code into source Go file.",
		Example: color.HiBlackString(`  # Merge deepseek generated code into source Go file in specified directory
  sponge assistant merge --type=deepseek --dir=/path/to/directory

  # Merge chatgpt generated code into source Go file in specified directory
  sponge assistant merge --type=chatgpt --dir=/path/to/directory

  # Merge gemini generated code into source Go file in specified directory
  sponge assistant merge --type=gemini --dir=/path/to/directory`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := checkDirAndFile(dir, files)
			if err != nil {
				return err
			}

			m := &mergeParams{
				assistantType:  strings.ToLower(assistantType),
				dir:            dir,
				specifiedFiles: files,
				isClean:        isClean,
			}

			fileMap, err := m.parseAssistantFiles()
			if err != nil {
				return err
			}
			if len(fileMap) == 0 {
				fmt.Printf("\nNo %s assistant generated code found, nothing to merge, please generate code before merging.\n", m.assistantType)
				return nil
			}

			var mergeCodes = make(map[string]string)
			for srcFile, genFile := range fileMap {
				var newCode string
				newCode, err = m.mergeGoFile(srcFile, genFile)
				if err != nil {
					return fmt.Errorf("merge %s into %s failed: %v", cutFilePath(genFile), cutFilePath(srcFile), err)
				}
				mergeCodes[srcFile] = newCode
			}

			var deleteFiles []string
			backupDir := getBackupDir()
			if len(mergeCodes) > 0 {
				fmt.Printf("Merged to Go files:\n")
				for srcFile, code := range mergeCodes {
					backupFile(srcFile, backupDir)
					if err = os.WriteFile(srcFile, []byte(code), 0666); err != nil {
						return err
					}
					fmt.Printf("    %s %s  →  %s\n", successSymbol,
						color.HiGreenString(cutFilePath(fileMap[srcFile])), color.HiGreenString(cutFilePath(srcFile)))
					deleteFiles = append(deleteFiles, fileMap[srcFile])
				}
				fmt.Println()
			}

			if m.isClean {
				deleteGenFiles(deleteFiles, m.assistantType)
			}

			if len(mergeCodes) > 0 {
				fmt.Printf("\n[Tip] You can view the pre-merge Go code files here:\n    %s\n\n", backupDir)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&assistantType, "type", "t", "", "assistant type, supported types: chatgpt, deepseek, gemini")
	_ = cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "input directory")
	cmd.Flags().StringSliceVarP(&files, "file", "f", nil, "specified Go files")
	cmd.Flags().BoolVarP(&isClean, "is-clean", "c", false, "clean up assistant generated code after merge")

	return cmd
}

type mergeParams struct {
	assistantType  string
	dir            string
	specifiedFiles []string
	isClean        bool
}

func (m *mergeParams) parseAssistantFiles() (map[string]string, error) {
	fileMap := make(map[string]string) // srcFile -> genFile

	for _, genFile := range m.specifiedFiles {
		srcFile, ok := m.getSourceFile(genFile)
		if ok {
			fileMap[srcFile] = genFile
		}
	}

	if m.dir != "" {
		files, err := gofile.ListFiles(m.dir, gofile.WithSuffix(getAssistantSuffixed(m.assistantType)))
		if err != nil {
			return nil, err
		}
		for _, genFile := range files {
			srcFile, ok := m.getSourceFile(genFile)
			if ok {
				fileMap[srcFile] = genFile
			}
		}
	}

	return fileMap, nil
}

func (m *mergeParams) getSourceFile(file string) (string, bool) {
	switch m.assistantType {
	case typeDeepSeek, typeChatGPT, typeGemini:
		if strings.HasSuffix(file, getAssistantSuffixed(m.assistantType)) {
			srcGoFile := strings.TrimSuffix(file, getAssistantSuffixed(m.assistantType)) + ".go"
			if gofile.IsExists(srcGoFile) {
				return srcGoFile, true
			}
		}
	}
	return "", false
}

func (m *mergeParams) mergeGoFile(srcFile string, genFile string) (string, error) {
	srcCode, err := os.ReadFile(srcFile)
	if err != nil {
		return "", err
	}

	genCode, err := os.ReadFile(genFile)
	if err != nil {
		return "", err
	}
	genCodes := extractGoCode(string(genCode))
	genCode = []byte(strings.Join(genCodes, "\n\n"))
	opts := checkPackageName(genFile, genCode)

	codeAst, err := goast.MergeGoCode(srcCode, genCode, opts...)
	if err != nil {
		return "", err
	}
	codeAst.FilePath = srcFile

	return codeAst.Code, nil
}

func checkPackageName(file string, data []byte) []goast.CodeAstOption {
	var opts = []goast.CodeAstOption{goast.WithCoverSameFunc()}
	fp := filepath.ToSlash(filepath.Clean(file))
	if strings.Contains(fp, "/dao/") {
		if bytes.Contains(data, []byte("package dao")) {
			opts = append(opts, goast.WithIgnoreMergeFunc(
				"Create",
				"GetByID",
				"UpdateByID",
				"GetByID",
				"GetByColumns",
				"CreateByTx",
				"DeleteByTx",
				"UpdateByTx",
			))
		}
	}

	return opts
}
