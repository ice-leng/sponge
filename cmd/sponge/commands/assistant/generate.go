package assistant

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/pkg/aicli"
	"github.com/go-dev-frame/sponge/pkg/goast"
	"github.com/go-dev-frame/sponge/pkg/gofile"
)

// GenerateCommand  command
func GenerateCommand() *cobra.Command {
	var (
		assistantType string
		apiKey        string
		model         string
		roleDesc      string
		maxToken      int
		temperature   float32
		enableContext bool

		onlyPrintPrompt bool // for test only
		maxAssistantNum int
		dir             string
		files           []string // specified Go files
	)

	//nolint
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate code using AI assistant",
		Long:  "Generate code using AI assistant. Automatically locate the positions in Go files that require code implementation, and let the AI assistant generate the corresponding business logic based on the context.",
		Example: color.HiBlackString(`  # Generate code using deepseek, default model is deepseek-chat, you can specify deepseek-reasoner by --model.
  sponge assistant generate --type=deepseek --api-key=your-api-key --dir=your-project-dir
  
  # Generate code using gemini, default model is gemini-2.5-pro-exp-03-25
  sponge assistant generate --type=gemini --api-key=your-api-key --dir=your-project-dir

  # Generate code using chatgpt, default model is gpt-4o
  sponge assistant generate --type=chatgpt --api-key=your-api-key --dir=your-project-dir`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := checkDirAndFile(dir, files)
			if err != nil {
				return err
			}

			isUseChinese, fileCodeMap, err := parseFiles(dir, files)
			if err != nil {
				return err
			}
			total := len(fileCodeMap)
			if total == 0 {
				fmt.Println(ErrnoAssistantMarker)
				return nil
			}

			if maxAssistantNum > total {
				maxAssistantNum = total
			}

			if roleDesc == "" {
				if isUseChinese {
					roleDesc = gopherRoleDescCN
				} else {
					roleDesc = gopherRoleDescEN
				}
			}

			assistantType = strings.ToLower(assistantType)
			asst := &assistantParams{
				Type:          assistantType,
				apiKey:        apiKey,
				model:         defaultModelMap[assistantType],
				enableContext: true,
				roleDesc:      roleDesc,
				maxToken:      maxToken,
				temperature:   temperature,
			}

			g := &assistantGenerator{
				maxAssistantNum: maxAssistantNum,
				asst:            asst,

				fileCodeMap: fileCodeMap,

				dir:               dir,
				files:             files,
				isOnlyPrintPrompt: onlyPrintPrompt,
			}
			err = g.generateCode()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&assistantType, "type", "t", "", "assistant type, supported types: chatgpt, deepseek, gemini")
	_ = cmd.MarkFlagRequired("type")
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "assistant api key")
	_ = cmd.MarkFlagRequired("api-key")
	cmd.Flags().StringVarP(&model, "model", "m", "", "assistant model, corresponding assistant type.")
	cmd.Flags().StringVarP(&roleDesc, "role", "r", "", "role description, for example, you are a psychologist.")
	cmd.Flags().IntVarP(&maxToken, "max-token", "s", 0, "maximum number of tokens")
	cmd.Flags().Float32VarP(&temperature, "temperature", "e", 0, "temperature of the model")
	cmd.Flags().BoolVarP(&enableContext, "enable-context", "c", false, "whether the assistant supports context")
	cmd.Flags().BoolVarP(&onlyPrintPrompt, "only-print-prompt", "p", false, "skip AI assistant request, only print prompt")
	cmd.Flags().IntVarP(&maxAssistantNum, "max-assistant-num", "n", 20, "maximum number of assistant running simultaneously")
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Go project directory")
	cmd.Flags().StringSliceVarP(&files, "file", "f", nil, "specified Go files")

	return cmd
}

type assistantGenerator struct {
	maxAssistantNum int
	asst            *assistantParams

	fileCodeMap map[string]*codeInfo

	dir   string
	files []string

	isOnlyPrintPrompt bool
}

func (g *assistantGenerator) generateCode() error {
	fileCount := len(g.fileCodeMap)

	// initialize worker pool
	workerPool, err := NewWorkerPool(context.Background(), g.maxAssistantNum, fileCount)
	if err != nil {
		return err
	}
	workerPool.Start()

	fmt.Printf("\n[INFO] Detected %s files for code generation. Processing concurrently with %s AI assistants (%s).\n\n",
		color.HiCyanString(strconv.Itoa(fileCount)), color.HiCyanString(strconv.Itoa(g.maxAssistantNum)), g.asst.model)

	jobID := 0

	// submit tasks to worker pool
	for file, info := range g.fileCodeMap {
		jobID++
		dependentFile, prompt := getPrompt(file, info)
		client, err := g.asst.newClient()
		if err != nil {
			return err
		}
		task := &assistantTask{
			jobID:         jobID,
			file:          file,
			dependentFile: dependentFile,
			funcNames:     info.getFuncNames(),
			prompt:        prompt,
			client:        client,
			Type:          g.asst.Type,

			isOnlyPrintPrompt: g.isOnlyPrintPrompt,
		}
		err = workerPool.Submit(Job{
			ID:   jobID,
			Task: task,
		}, time.Millisecond*20)
		if err != nil {
			return err
		}
	}

	// wait for all tasks to complete
	go func() {
		workerPool.Wait()
	}()

	var (
		outputFiles  []string
		successCount int
		failedCount  int
		p            = newPrintLog(time.Millisecond * 200)
	)

	// handle results from worker pool
	for result := range workerPool.Results() {
		reply := result.Value.(*Reply)
		if g.isOnlyPrintPrompt {
			l := fmt.Sprintf("File: [%s]\nPrompt: %s\n\n%s\n\n", cutFilePath(reply.SrcFile), reply.Prompt, strings.Repeat("-", 80))
			p.StopPrint(l)
			p = newPrintLog()
			continue
		}
		if result.Err != nil {
			failedCount++
			l := fmt.Sprintf("\n[ERROR] Job %s - File: [%s] | Functions: [%s] → Code generation failed! Error: [%s]\n",
				color.HiCyanString(strconv.Itoa(reply.JobID)),
				color.HiCyanString(cutFilePath(reply.SrcFile)),
				color.HiCyanString(strings.Join(reply.Functions, ", ")),
				color.HiRedString(reply.ErrMsg),
			)
			p.StopPrint(l)
			p = newPrintLog()
		} else {
			successCount++
			var newFiles []string
			for newFile := range reply.Contents {
				newFiles = append(newFiles, cutFilePath(newFile))
				outputFiles = append(outputFiles, newFile)
			}
			l := fmt.Sprintf("\n[SUCCESS] Job %s - File: [%s] | Functions: [%s] | Output: [%s] | Time: %s\n",
				color.HiCyanString(strconv.Itoa(reply.JobID)),
				color.HiCyanString(cutFilePath(reply.SrcFile)),
				color.HiCyanString(strings.Join(reply.Functions, ", ")),
				color.HiGreenString(strings.Join(newFiles, ", ")),
				color.HiCyanString("%.2fs", result.EndTime.Sub(result.StartTime).Seconds()),
			)
			p.StopPrint(l)
			p = newPrintLog()
		}
	}

	// stop worker pool
	workerPool.Stop()

	time.Sleep(time.Millisecond * 220)
	p.StopPrint("")

	total := successCount + failedCount
	if total > 0 {
		fmt.Printf("\nJobs Summary:\n    → Total Jobs: %d\n    → Successful Jobs: %d\n    → Failed Jobs: %d\n\n", total, successCount, failedCount)
	}

	if len(outputFiles) > 0 {
		fmt.Println("Successful Jobs output files:")
		for _, file := range outputFiles {
			fmt.Printf("    %s %s\n", successSymbol, color.HiGreenString(cutFilePath(file)))
		}
	}

	return nil
}

type assistantTask struct {
	jobID         int
	file          string
	dependentFile string
	funcNames     []string
	prompt        string
	client        aicli.Assistanter
	Type          string

	isOnlyPrintPrompt bool
}

// Reply execute assistant task result
type Reply struct {
	JobID     int      `json:"jobID"`
	ErrMsg    string   `json:"errMsg"`
	SrcFile   string   `json:"srcFile"`
	Functions []string `json:"functions"`
	Prompt    string   `json:"prompt"`

	Contents map[string]string `json:"contents"` // reply contents
}

// Execute execute assistant task
func (t *assistantTask) Execute(ctx context.Context) (interface{}, error) {
	taskReply := &Reply{
		JobID:     t.jobID,
		SrcFile:   t.file,
		Functions: t.funcNames,
		Prompt:    t.prompt,
	}

	if t.isOnlyPrintPrompt {
		return taskReply, nil
	}

	fmt.Printf("[START] Job %s - File: [%s] | Functions: [%s]\n\n",
		color.HiCyanString(strconv.Itoa(t.jobID)),
		color.HiCyanString(cutFilePath(t.file)),
		color.HiCyanString(strings.Join(t.funcNames, ", ")),
	)

	streamReply := t.client.SendStream(ctx, t.prompt)
	assistantReply := ""
	for content := range streamReply.Content {
		assistantReply += content
	}
	if streamReply.Err != nil {
		taskReply.ErrMsg = streamReply.Err.Error()
		return taskReply, streamReply.Err
	}

	codes := parseCode(assistantReply)
	newFile, err := saveAssistantCode(t.file, codes[0], t.Type)
	if err != nil {
		taskReply.ErrMsg = err.Error()
		return taskReply, err
	}

	contents := make(map[string]string, len(codes))
	contents[newFile] = codes[0]

	if t.dependentFile != "" && len(codes) > 1 {
		newDependentFile, err := saveAssistantCode(t.dependentFile, codes[1], t.Type)
		if err != nil {
			taskReply.ErrMsg = err.Error()
			return taskReply, err
		}
		contents[newDependentFile] = codes[1]
	}
	taskReply.Contents = contents

	return taskReply, nil
}

func parseFiles(dir string, specifiedFiles []string) (bool, map[string]*codeInfo, error) {
	var isChinese bool
	var fileCodeMap = make(map[string]*codeInfo)

	for _, file := range specifiedFiles {
		if info := extractFuncCodeBlock(file); info != nil {
			fileCodeMap[file] = info
			isChinese = isChinese || info.isUseChinesePrompt()
		}
	}

	if dir != "" {
		files, err := gofile.ListFiles(dir, gofile.WithSuffix(".go")) //nolint
		if err != nil {
			return false, nil, err
		}
		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") ||
				strings.HasSuffix(file, ".pb.go") ||
				strings.HasSuffix(file, ".validate.go") {
				continue
			}
			if info := extractFuncCodeBlock(file); info != nil {
				fileCodeMap[file] = info
				isChinese = isChinese || info.isUseChinesePrompt()
			}
		}
	}

	return isChinese, fileCodeMap, nil
}

func saveAssistantCode(file string, code string, asstType string) (string, error) {
	dirPath := gofile.GetDir(file)
	if !gofile.IsExists(dirPath) {
		err := os.MkdirAll(dirPath, 0666)
		if err != nil {
			return "", fmt.Errorf("failed to create directory %s: %v", dirPath, err)
		}
	}

	newFile := fmt.Sprintf("%s.%s.md", file, asstType)
	err := os.WriteFile(newFile, []byte(code), 0666)
	if err != nil {
		return "", fmt.Errorf("failed to save assistant code to file %s: %v", newFile, err)
	}
	return newFile, nil
}

type codeInfo struct {
	funcInfos []goast.FuncInfo
	code      []byte
}

func extractFuncCodeBlock(file string) *codeInfo {
	if file == "" {
		return nil
	}

	code, infos, err := goast.FilterFuncCodeByFile(file)
	if err != nil {
		return nil
	}

	return &codeInfo{
		funcInfos: infos,
		code:      code,
	}
}

func isAllHaveExampleCode(info *codeInfo) bool {
	return bytes.Count(info.code, []byte("// fill in the business logic code here")) == len(info.funcInfos)
}

func getDefaultPrompt(info *codeInfo) string {
	format := ""
	if info.isUseChinesePrompt() {
		format = defaultPromptFormatCN
	} else {
		format = defaultPromptFormatEN
	}
	return fmt.Sprintf(format, info.joinFuncNames()) +
		fmt.Sprintf("\n```go\n%s\n```", string(info.code))
}

// nolint
func getPrompt(file string, info *codeInfo) (dependentFileFullPath string, prompt string) {
	//var prompt string
	dirName := getLastDirName(file)
	isChinese := info.isUseChinesePrompt()

	switch dirName {
	case "handler", "service", "biz":
		internalDirName := strings.TrimSuffix(gofile.GetDir(file), dirName) // end with filepath.Separator
		if getLastDirName(internalDirName) != "internal" {
			return "", getDefaultPrompt(info)
		} else {
			if !isAllHaveExampleCode(info) {
				return "", getDefaultPrompt(info)
			}
			internalDirName = "internal" + string(filepath.Separator)
		}

		fileName := gofile.GetFilename(file)
		srcFile := internalDirName + dirName + string(filepath.Separator) + fileName
		daoFile := internalDirName + "dao" + string(filepath.Separator) + fileName
		dbFile := internalDirName + "database" + string(filepath.Separator) + "init.go"
		objName := strings.TrimSuffix(fileName, ".go")
		isMongo := isMongoOrmType(dbFile)
		daoCode := getDaoCode(isMongo, objName, isChinese)
		funcNames := info.joinFuncNames()

		format := ""
		if isChinese {
			format = promptFormatCN
		} else {
			format = promptFormatEN
		}
		prompt = fmt.Sprintf(format,
			funcNames,
			srcFile, funcNames,
			srcFile, daoFile,
			srcFile, daoFile,
			srcFile, daoFile,
			srcFile, fmt.Sprintf("\n```go\n%s\n```", string(info.code)),
			daoFile, fmt.Sprintf("\n```go\n%s\n```", daoCode),
		)
		modelCode := getModelCode(file, dirName, fileName, isChinese)
		if modelCode != "" {
			prompt += modelCode
		}
		dependentFileFullPath = strings.TrimSuffix(gofile.GetDir(file), dirName) + "dao" + string(filepath.Separator) + fileName
	default:
		prompt = getDefaultPrompt(info)
	}

	return dependentFileFullPath, prompt
}

func (c *codeInfo) isUseChinesePrompt() bool {
	var hasChinese bool
	for _, info := range c.funcInfos {
		if containsChinese(info.Comment) {
			hasChinese = true
			break
		}
	}
	return hasChinese
}

func (c *codeInfo) getFuncNames() []string {
	var funcNames []string
	for _, funcInfo := range c.funcInfos {
		funcNames = append(funcNames, funcInfo.Name)
	}
	return funcNames
}

func (c *codeInfo) joinFuncNames() string {
	delimiter := "、"
	if !c.isUseChinesePrompt() {
		delimiter = ", "
	}
	var funcNames []string
	for _, funcInfo := range c.funcInfos {
		funcNames = append(funcNames, funcInfo.Name)
	}
	return strings.Join(funcNames, delimiter)
}

func containsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func getLastDirName(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	if info.IsDir() {
		return filepath.Base(path)
	}
	return filepath.Base(filepath.Dir(path))
}

func isMongoOrmType(dbFile string) bool {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		return false
	}
	if bytes.Contains(data, []byte(`"github.com/go-dev-frame/sponge/pkg/mgo"`)) {
		return true
	}
	return false
}

func checkDirAndFile(dir string, files []string) error {
	if dir == "" && len(files) == 0 {
		return fmt.Errorf("please specify flag --dir or --file")
	}
	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("directory %s does not exist", dir)
		}
	}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", file)
		}
	}
	return nil
}
