package assistant

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"

	"github.com/go-dev-frame/sponge/pkg/aicli"
	"github.com/go-dev-frame/sponge/pkg/aicli/chatgpt"
	"github.com/go-dev-frame/sponge/pkg/aicli/deepseek"
	"github.com/go-dev-frame/sponge/pkg/aicli/gemini"
	"github.com/go-dev-frame/sponge/pkg/gobash"
	"github.com/go-dev-frame/sponge/pkg/gofile"
	"github.com/go-dev-frame/sponge/pkg/utils"
)

const (
	typeChatGPT  = "chatgpt"
	typeDeepSeek = "deepseek"
	typeGemini   = "gemini"

	gopherRoleDescCN = aicli.GopherRoleDescCN
	gopherRoleDescEN = aicli.GopherRoleDescEN

	successSymbol = "✔ "
	failureSymbol = "❌ "
)

var assistantTypeMap = map[string]string{
	typeChatGPT:  "ChatGPT",
	typeDeepSeek: "DeepSeek",
	typeGemini:   "Gemini",
}

var defaultModelMap = map[string]string{
	typeChatGPT:  chatgpt.DefaultModel,
	typeDeepSeek: deepseek.DefaultModel,
	typeGemini:   gemini.DefaultModel,
}

type assistantParams struct {
	Type string

	apiKey        string
	model         string
	enableContext bool

	// only for chatgpt and deepseek
	roleDesc    string
	maxToken    int
	temperature float32
}

func (a *assistantParams) newClient() (aicli.Assistanter, error) {
	asstType := strings.ToLower(a.Type)
	switch asstType {
	case typeChatGPT, typeDeepSeek:
		var opts []chatgpt.ClientOption
		if a.model != "" {
			opts = append(opts, chatgpt.WithModel(a.model))
		}
		if a.enableContext {
			opts = append(opts, chatgpt.WithEnableContext())
		}
		if a.roleDesc != "" {
			opts = append(opts, chatgpt.WithInitialRole(a.roleDesc))
		}
		if a.maxToken > 0 {
			opts = append(opts, chatgpt.WithMaxTokens(a.maxToken))
		}
		if a.temperature > 0 {
			opts = append(opts, chatgpt.WithTemperature(a.temperature))
		}

		if asstType == typeChatGPT {
			return chatgpt.NewClient(a.apiKey, opts...)
		} else if asstType == typeDeepSeek {
			return deepseek.NewClient(a.apiKey, opts...)
		}

	case typeGemini:
		var opts []gemini.ClientOption
		if a.model != "" {
			opts = append(opts, gemini.WithModel(a.model))
		}
		if a.enableContext {
			opts = append(opts, gemini.WithEnableContext())
		}
		return gemini.NewClient(a.apiKey, opts...)
	}

	return nil, fmt.Errorf("unsupported assistant type: %s", a.Type)
}

// --------------------------------------------------------------------------

const (
	gormDao = `package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"github.com/go-dev-frame/sponge/internal/cache"
	"github.com/go-dev-frame/sponge/internal/database"
	"github.com/go-dev-frame/sponge/internal/model"
)

var _ UserExampleDao = (*userExampleDao)(nil)

// UserExampleDao defining the dao interface
type UserExampleDao interface {
	Create(ctx context.Context, table *model.UserExample) error
	UpdateByID(ctx context.Context, table *model.UserExample) error
	GetByID(ctx context.Context, id uint64) (*model.UserExample, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.UserExample) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.UserExample) error
}

type userExampleDao struct {
	db    *gorm.DB
	cache cache.UserExampleCache // if nil, the cache is not used.
	sfg   *singleflight.Group    // if cache is nil, the sfg is not used.
}

// NewUserExampleDao creating the dao interface
func NewUserExampleDao(db *gorm.DB, xCache cache.UserExampleCache) UserExampleDao {
	if xCache == nil {
		return &userExampleDao{db: db}
	}
	return &userExampleDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

`

	mongoDao = `package dao

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/singleflight"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/mgo"
	"github.com/go-dev-frame/sponge/pkg/mgo/query"

	"github.com/go-dev-frame/sponge/internal/cache"
	"github.com/go-dev-frame/sponge/internal/database"
	"github.com/go-dev-frame/sponge/internal/model"
)

var _ UserExampleDao = (*userExampleDao)(nil)

// UserExampleDao defining the dao interface
type UserExampleDao interface {
	Create(ctx context.Context, record *model.UserExample) error
	UpdateByID(ctx context.Context, record *model.UserExample) error
	GetByID(ctx context.Context, id string) (*model.UserExample, error)
}

type userExampleDao struct {
	collection *mongo.Collection
	cache      cache.UserExampleCache // if nil, the cache is not used.
	sfg        *singleflight.Group    // if cache is nil, the sfg is not used.
}

// NewUserExampleDao creating the dao interface
func NewUserExampleDao(collection *mongo.Collection, xCache cache.UserExampleCache) UserExampleDao {
	if xCache == nil {
		return &userExampleDao{collection: collection}
	}
	return &userExampleDao{
		collection: collection,
		cache:      xCache,
		sfg:        new(singleflight.Group),
	}
}

`

	codeDelimiterMarker = "/**code-delimiter**/"
)

// nolint
var ErrnoAssistantMarker = fmt.Errorf("\n%s%s\n\n",
	`No code requiring AI assistant generation has been detected. To enable the AI assistant to automatically generate code, the following two conditions must be met:  
1. The function body must include the panic("implement me") marker.  
2. A comment describing the function's purpose must be added above the function.  

Example:`, fmt.Sprintf(`

    %s
    func MyFunc() {
        %s
    }`, color.HiCyanString("// Describe the specific functionality of the function"), color.HiCyanString(`panic("implement me")`)))

func getDaoCode(isMongo bool, objName string, isChinese bool) string {
	var daoCode string
	if isMongo {
		daoCode = mongoDao
	} else {
		daoCode = gormDao
	}
	daoCode = strings.Replace(daoCode, "UserExample", capitalize(objName), -1)
	daoCode = strings.Replace(daoCode, "userExample", objName, -1)

	decs := ""
	if isChinese {
		decs = fmt.Sprintf("// 在这里实现业务逻辑代码, 方法函数的接收者是 %sDao， "+
			"已知 %sDao 结构体实现了 %sDao 接口，其中 %s 方法函数已经实现了，不需要在代码上额外定义和实现这些方法函数。"+
			"如果实现业务逻辑代码过程中需要用到这些方法函数，可以直接调用。",
			objName, objName, capitalize(objName), getFuncNames(isMongo, isChinese))
	} else {
		decs = fmt.Sprintf("// Implement the business logic code here. The method function receiver is %sDao."+
			" It is known that the %sDao struct implements the %sDao interface. The %s method functions"+
			" have already been implemented and do not need to be defined and implemented in code. If you need to use"+
			" these method functions while implementing the business logic, you can call them directly.",
			objName, objName, capitalize(objName), getFuncNames(isMongo, isChinese))
	}

	return daoCode + decs + "\n\n"
}

func getFuncNames(isMongo bool, isChinese bool) string {
	baseFuncNames := "Create, UpdateByID, GetByID"
	if !isMongo {
		baseFuncNames += ", CreateByTx, DeleteByTx, UpdateByTx"
	}
	if isChinese {
		baseFuncNames = strings.ReplaceAll(baseFuncNames, ", ", "、")
	}
	return baseFuncNames
}

func getModelCode(file string, dirName string, fileName string, isChinese bool) string {
	modelFile := strings.TrimSuffix(gofile.GetDir(file), dirName) + "model" + string(filepath.Separator) + fileName
	data, err := os.ReadFile(modelFile)
	if err != nil {
		return ""
	}

	modelCode := ""
	modelFile = "internal" + string(filepath.Separator) + "model" + string(filepath.Separator) + fileName
	if isChinese {
		modelCode = fmt.Sprintf("\n\n原始文件 %s 代码如下：", modelFile)
	} else {
		modelCode = fmt.Sprintf("\n\nThe original file %s code is as follows:", modelFile)
	}

	return modelCode + fmt.Sprintf("\n```go\n%s\n```", string(data))
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// nolint
var (
	defaultPromptFormatEN = `
Please implement the business logic for the functions %s based on the descriptions provided in the comments within the following Go code. Specific implementation requirements are as follows:

1. **Function Implementation**: Ensure that the code functionality aligns with the descriptions in the function comments. Add comments to key steps during the implementation of the code logic to facilitate understanding and maintenance.
2. **Dependency Management**: If you need to import other packages or third-party libraries, ensure these packages already exist and are usable.
3. **Code Quality**: Adhere to Go language's coding style, naming conventions, and commenting standards. Ensure the code is clear, readable, and well-commented.
4. **Compilability**: Ensure that the implemented code can be compiled successfully.
5. **Returned Code**: Please provide the complete Go language code, ensuring that the code is in pure Go format, just use markdown tag go language code block, no additional text description code block.
`
	defaultPromptFormatCN = `
请根据以下Go语言代码中的函数注释描述，实现函数 %s 的业务逻辑。具体实现要求如下： 

1. **功能实现**：确保代码功能与函数注释描述一致。在实现代码逻辑时，对关键步骤添加注释，以便于理解和维护。
2. **依赖管理**：如果需要import其他包或第三方库，请确保这些包已存在并可用。
3. **代码质量**：遵循Go语言的代码风格、命名规范和注释规范，确保代码清晰易读，注释完整。
4. **可编译性**：确保实现的代码能够顺利编译通过。
5. **返回代码**：请提供完整的Go语言代码，确保代码为纯Go语言格式，只需使用Markdown标记Go语言代码块，无需额外文字描述说明代码块。
`

	promptFormatEN = `
Below are the codes for the Go files. Please implement the business logic for the function %s according to the following requirements:

1. **Function Description**:
   - The function comments in the %s file describe the functionality of %s, and the function body contains example reference code.
   - The specific business logic is not implemented in %s but rather in %s.
   - The function in %s needs to call the business logic function implemented in %s.
   - Add comments to key steps during the implementation of the code logic to facilitate understanding and maintenance.

2. **Dependency Management**:
   - Given that the import package in the code already exists, it can be used directly.
   - If you need to import additional packages, make sure they already exist and are available.

3. **Code Quality**:
   - Adhere to Go language's coding style, naming conventions, and commenting standards.
   - Ensure the code is clear, readable, and well-commented.

4. **Compilability**:
   - Ensure that the implemented code can be compiled successfully.

5. **Returned Code**:
   - Please provide the complete Go language code, ensuring that the code is in pure Go format, just use markdown tag go language code block, no additional text description code block. The code for the two files should be separated by /**code-delimiter**/, with the code for the first %s file preceding the code for the second %s file, to facilitate easy splitting of the code based on the /**code-delimiter**/ marker.


The original file %s code is as follows:
%s


The original file %s code is as follows:
%s
`
	promptFormatCN = `
下面是 Go 文件的代码。请根据以下要求实现函数 %s 的业务逻辑：

1. **功能描述**：
   - %s 文件中的函数注释描述了 %s 的功能，函数体里包含了示例参考代码。
   - 具体的业务逻辑不在 %s 中实现，而是在 %s 中实现。
   - %s 中的函数需要调用 %s 中实现的业务逻辑函数。
   - 在实现代码逻辑时，对关键步骤添加注释，以便于理解和维护。

2. **依赖管理**：
   - 已知代码中 import 的包已存在，可直接使用。
   - 如果需要 import 其他包，请确保这些包已存在并可用。

3. **代码质量**：
   - 遵循Go语言的代码风格、命名规范和注释规范。
   - 确保代码清晰易读，注释完整。

4. **可编译性**：
   - 确保实现的代码能够顺利编译通过。

5. **返回代码**：
   - 请提供完整的Go语言代码，确保代码为纯Go语言格式，只需使用Markdown标记Go语言代码块，无需额外文字描述说明代码块。两个文件的代码应以 /**code-delimiter**/ 分隔，其中第一个 %s 文件的代码在前，第二个 %s 文件的代码在后，以便后续通过 /**code-delimiter**/ 标记轻松分割代码。


原始文件 %s 代码如下：
%s


原始文件 %s 代码如下：
%s
`
)

// extractGoCode extracts the Go code blocks from the given markdown string.
func extractGoCode(markdown string) []string {
	var goCodeBlocks []string
	scanner := bufio.NewScanner(strings.NewReader(markdown))
	inCodeBlock := false
	var currentBlock strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// dealing with new ```go
		if strings.HasPrefix(line, "```go") {
			// if it is already in the code block, it means that the previous code block is missing the close identifier and stores it first.
			if inCodeBlock {
				goCodeBlocks = append(goCodeBlocks, currentBlock.String())
				currentBlock.Reset()
			}
			inCodeBlock = true
			continue
		}

		// processing ``` close code block
		if inCodeBlock && strings.HasPrefix(line, "```") {
			inCodeBlock = false
			goCodeBlocks = append(goCodeBlocks, currentBlock.String())
			currentBlock.Reset()
			continue
		}

		// record code content
		if inCodeBlock {
			currentBlock.WriteString(line)
			currentBlock.WriteString("\n")
		}
	}

	// prevents the last code block from not closing
	if inCodeBlock && currentBlock.Len() > 0 {
		goCodeBlocks = append(goCodeBlocks, currentBlock.String())
	}

	return goCodeBlocks
}

func parseCode(code string) []string {
	n := strings.Count(code, "```go")
	switch n {
	case 0:
		if strings.Contains(code, codeDelimiterMarker) {
			ss := strings.Split(code, codeDelimiterMarker)
			var codes []string
			for _, s := range ss {
				if strings.Contains(s, "package ") {
					codes = append(codes, s)
				}
			}
			if len(codes) > 0 {
				return reassembleGoMarkdown(codes)
			}
		}
		if strings.Contains(code, "\npackage ") {
			return reassembleGoMarkdown([]string{code})
		}
		return []string{code}
	case 1:
		goCodes := extractGoCode(code)
		var codes []string
		for _, goCode := range goCodes {
			codes = append(codes, strings.Split(goCode, codeDelimiterMarker)...)
		}
		return reassembleGoMarkdown(codes)
	}

	code = strings.ReplaceAll(code, codeDelimiterMarker, "\n\n")
	return reassembleGoMarkdown(extractGoCode(code))
}

func reassembleGoMarkdown(codes []string) []string {
	for i, c := range codes {
		codes[i] = "```go\n" + strings.TrimSpace(c) + "\n```\n"
	}
	return codes
}

func cutFilePath(fullPath string) string {
	cwd, _ := os.Getwd()
	if strings.HasPrefix(fullPath, cwd) {
		return strings.TrimPrefix(fullPath, cwd+string(filepath.Separator))
	}
	return fullPath
}

func newPrintLog(t ...time.Duration) *utils.WaitPrinter {
	if len(t) > 0 {
		time.Sleep(t[0])
	}
	p := utils.NewWaitPrinter(time.Millisecond * 200)
	p.LoopPrint("Waiting for assistant responses ")
	return p
}

func getAssistantSuffixed(assistantType string) string {
	return ".go." + assistantType + ".md"
}

func deleteGenFiles(files []string, assistantType string) {
	var doneFiles []string
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			continue
		}
		doneFiles = append(doneFiles, cutFilePath(file))
	}

	if len(doneFiles) > 0 {
		fmt.Printf("Clean up code files generated by [%s] assistant:\n", assistantType)
		for _, file := range doneFiles {
			fmt.Printf("    %s %s\n", successSymbol, color.HiGreenString(cutFilePath(file)))
		}
		fmt.Println()
	}
}

func getRelativeDirAndFile(srcFile string) (dir string, file string) {
	dirPath, _ := filepath.Abs(".")
	filePath := strings.TrimPrefix(srcFile, dirPath+gofile.GetPathDelimiter())

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return srcFile, ""
	}

	if !fileInfo.IsDir() {
		dir = gofile.GetDir(filePath)
		file = gofile.GetFilename(filePath)
		return dir, file
	}
	return filePath, ""
}

func getBackupDir() string {
	var backupDir = os.TempDir() + string(os.PathSeparator) + "sponge_merge_backup_code"
	return backupDir + string(os.PathSeparator) + time.Now().Format("20060102T150405")
}

func backupFile(file string, backupDir string) {
	relPath, _ := getRelativeDirAndFile(file)
	bkDir := backupDir + string(os.PathSeparator) + relPath + string(os.PathSeparator)
	_ = os.MkdirAll(bkDir, 0744)
	_, _ = gobash.Exec("cp", file, bkDir)
}
