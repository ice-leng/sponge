package assistant

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	chatgpt "github.com/go-dev-frame/sponge/pkg/aicli/chatgpt"
	"github.com/go-dev-frame/sponge/pkg/aicli/deepseek"
	"github.com/go-dev-frame/sponge/pkg/gofile"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

// GenerateCommand  command
func GenerateCommand() *cobra.Command {
	var (
		assistantType string
		apiKey        string
		model         string
		role          string
		maxToken      int
		temperature   float32
		useContext    bool

		dir string
	)

	//nolint
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code for project",
		Long:  "Generate code for project using assistant.",
		Example: color.HiBlackString(`  # Generate code using ChatGPT assistant.
  sponge assistant generate --type=chatgpt --api-key=your-api-key --dir=./your-project

  # Generate code using ChatGPT assistant.
  sponge assistant generate --type=deepseek --api-key=your-api-key --dir=./your-project`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var opts []chatgpt.ClientOption
			if model != "" {
				opts = append(opts, chatgpt.WithModel(model))
			}
			if maxToken > 0 {
				opts = append(opts, chatgpt.WithMaxTokens(maxToken))
			}
			if temperature > 0 {
				opts = append(opts, chatgpt.WithTemperature(temperature))
			}
			if role != "" {
				opts = append(opts, chatgpt.WithRole(role))
			}
			opts = append(opts, chatgpt.WithUseContext(useContext))

			var client *chatgpt.Client
			var err error
			switch strings.ToLower(assistantType) {
			case typeChatGPT:
				client, err = chatgpt.NewClient(apiKey, opts...)
				if err != nil {
					return err
				}
			case typeDeepSeek:
				client, err = deepseek.NewClient(apiKey, opts...)
				if err != nil {
					return err
				}
			default:
				return cmd.Usage()
			}

			err = generateCode(client, dir)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&assistantType, "type", "t", "", "assistant type, e.g. chatgpt, deepseek")
	_ = cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "assistant api key")
	_ = cmd.MarkFlagRequired("api-key")
	cmd.Flags().StringVarP(&model, "model", "m", "", "assistant model, corresponding assistant type.")
	cmd.Flags().StringVarP(&role, "role", "r", "", "role of the model")
	cmd.Flags().IntVarP(&maxToken, "max-token", "s", 0, "maximum number of tokens")
	cmd.Flags().Float32VarP(&temperature, "temperature", "e", 0, "temperature of the model")
	cmd.Flags().BoolVarP(&useContext, "use-context", "c", false, "whether the assistant supports context")

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "project directory")
	_ = cmd.MarkFlagRequired("dir")

	return cmd
}

func generateCode(client *chatgpt.Client, dir string) error {
	files, err := getFiles(dir)
	if err != nil {
		return err
	}

	count := 0
	total := len(files)
	for file, data := range files {
		filename := gofile.GetFilename(file)
		filenameColor := color.HiCyanString(filename)
		count++

		p := utils.NewWaitPrinter(time.Millisecond * 200)
		fmt.Println()
		tip := fmt.Sprintf("[%d/%d] %s, assistant is analyzing and writing code ", count, total, filenameColor)
		p.LoopPrint(tip)

		reply, err := client.Send(context.Background(), getPrompt(data))
		if err != nil {
			p.StopPrint(filename + ", " + err.Error())
			return err
		}

		newFile := file + ".assistant"
		err = os.WriteFile(newFile, []byte(reply), 0666)
		if err != nil {
			p.StopPrint(newFile + ", " + err.Error())
			return err
		}
		tip2 := fmt.Sprintf("[%d/%d] %s, assistant has completed coding and saved it in %s", count, total, filenameColor, color.HiGreenString(newFile))
		p.StopPrint(tip2)
	}
	return nil
}

func getFiles(dir string) (map[string][]byte, error) {
	var promptFiles = make(map[string][]byte)
	handlerDirs, err := gofile.ListSubDirs(dir, "handler")
	if err != nil {
		return nil, err
	}
	for _, handlerDir := range handlerDirs {
		files, err := gofile.ListFiles(handlerDir, gofile.WithSuffix(".go")) //nolint
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if data := checkPrompt(file); data != nil {
				promptFiles[file] = data
			}
		}
	}

	serviceDirs, err := gofile.ListSubDirs(dir, "service")
	if err != nil {
		return nil, err
	}
	for _, serviceDir := range serviceDirs {
		files, err := gofile.ListFiles(serviceDir, gofile.WithSuffix(".go")) //nolint
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if data := checkPrompt(file); data != nil {
				promptFiles[file] = data
			}
		}
	}
	return promptFiles, nil
}

func checkPrompt(file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil
	}
	if bytes.Contains(data, []byte(`panic("prompt:`)) {
		return data
	}
	return nil
}

// nolint
func getPrompt(data []byte) string {
	cnPrompt := `下面go语言代码中，请根据 panic("prompt: 后面的提示要求，实现该方法函数的完整业务逻辑代码。如果有多个prompt提示表示多个方法函数需要实现，请按照提示请按顺序回答。` + "\n\n```go\n"
	enPrompt := `In the following go language code, please implement the complete business logic code of this method function according to the prompt requirements after panic("prompt: .  If there are multiple prompt prompts indicating that more than one method function needs to be implemented, please follow the prompts and answer in order.` + "\n\n```go\n"
	var prompt string

	localTime := time.Now()
	_, offset := localTime.Zone()
	east8Offset := 8 * 60 * 60
	if offset == east8Offset {
		prompt += cnPrompt + string(data) + "\n```"
	} else {
		prompt += enPrompt + string(data) + "\n```"
	}
	return prompt
}
