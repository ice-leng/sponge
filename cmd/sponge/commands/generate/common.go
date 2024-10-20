// Package generate is to generate code, including model, cache, dao, handler, http, service, grpc, grpc-gw, grpc-cli code.
package generate

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/huandu/xstrings"

	"github.com/zhufuyi/sponge/pkg/gobash"
	"github.com/zhufuyi/sponge/pkg/gofile"
	"github.com/zhufuyi/sponge/pkg/replacer"
)

const (
	defaultGoModVersion = "go 1.20"

	// TplNameSponge name of the template
	TplNameSponge = "sponge"

	// DBDriverMysql mysql driver
	DBDriverMysql = "mysql"
	// DBDriverPostgresql postgresql driver
	DBDriverPostgresql = "postgresql"
	// DBDriverTidb tidb driver
	DBDriverTidb = "tidb"
	// DBDriverSqlite sqlite driver
	DBDriverSqlite = "sqlite"
	// DBDriverMongodb mongodb driver
	DBDriverMongodb = "mongodb"

	undeterminedDBDriver = "undetermined" // used in services created based on protobuf.

	// code name
	codeNameHTTP        = "http"
	codeNameGRPC        = "grpc"
	codeNameHTTPPb      = "http-pb"
	codeNameGRPCPb      = "grpc-pb"
	codeNameGRPCGW      = "grpc-gw-pb"
	codeNameGRPCHTTP    = "grpc-http-pb"
	codeNameHandler     = "handler"
	codeNameHandlerPb   = "handler-pb"
	codeNameService     = "service"
	codeNameServiceHTTP = "service-handler"
	codeNameDao         = "dao"
	codeNameProtobuf    = "protobuf"
	codeNameModel       = "model"
	codeNameGRPCConn    = "grpc-conn"
	codeNameCache       = "cache"

	wellPrefix    = "## "
	mgoSuffix     = ".mgo"
	pkgPathSuffix = "/pkg"
	expSuffix     = ".exp"
	apiDocsSuffix = " api docs"
)

var (
	modelFile     = "model/userExample.go"
	modelFileMark = "// todo generate model code to here"

	modelInitDBFile     = "model/init.go"
	modelInitDBFileMark = "// todo generate initialisation database code here"

	cacheFile = "cache/cacheNameExample.go"

	daoFile     = "dao/userExample.go"
	daoMgoFile  = "dao/userExample.go.mgo"
	daoFileMark = "// todo generate the update fields code to here"
	daoTestFile = "dao/userExample_test.go"

	handlerFile       = "types/userExample_types.go"
	handlerMgoFile    = "types/userExample_types.go.mgo"
	handlerFileMark   = "// todo generate the request and response struct to here"
	handlerTestFile   = "handler/userExample_test.go"
	handlerPbTestFile = "handler/userExample_logic_test.go"

	handlerLogicFile = "handler/userExample_logic.go"
	serviceLogicFile = "service/userExample.go"
	embedTimeMark    = "// todo generate the conversion createdAt and updatedAt code here"

	httpFile = "server/http.go"

	protoFile     = "v1/userExample.proto"
	protoFileMark = "// todo generate the protobuf code here"

	serviceTestFile      = "service/userExample_test.go"
	serviceClientFile    = "service/userExample_client_test.go"
	serviceClientMgoFile = "service/userExample_client_test.go.mgo"
	serviceFileMark      = "// todo generate the service struct code here"

	dockerFile     = "scripts/build/Dockerfile"
	dockerFileMark = "# todo generate dockerfile code for http or grpc here"

	dockerFileBuild     = "scripts/build/Dockerfile_build"
	dockerFileBuildMark = "# todo generate dockerfile_build code for http or grpc here"

	imageBuildFile     = "scripts/image-build.sh"
	imageBuildFileMark = "# todo generate image-build code for http or grpc here"

	imageBuildLocalFile     = "scripts/image-build-local.sh"
	imageBuildLocalFileMark = "# todo generate image-build-local code for http or grpc here"

	dockerComposeFile     = "deployments/docker-compose/docker-compose.yml"
	dockerComposeFileMark = "# todo generate docker-compose.yml code for http or grpc here"

	k8sDeploymentFile     = "deployments/kubernetes/serverNameExample-deployment.yml"
	k8sDeploymentFileMark = "# todo generate k8s-deployment.yml code for http or grpc here"

	k8sServiceFile     = "deployments/kubernetes/serverNameExample-svc.yml"
	k8sServiceFileMark = "# todo generate k8s-svc.yml code for http or grpc here"

	protoShellFile         = "scripts/protoc.sh"
	protoShellFileGRPCMark = "# todo generate grpc files here"
	protoShellFileMark     = "# todo generate api template code command here"

	appConfigFile      = "configs/serverNameExample.yml"
	appConfigFileMark  = "# todo generate http or rpc server configuration here"
	appConfigFileMark2 = "# todo generate the database configuration here"

	expectedSQLForDeletion = "expectedSQLForDeletion := \"UPDATE .*\""

	//deploymentConfigFile     = "kubernetes/serverNameExample-configmap.yml"
	//deploymentConfigFileMark = "# todo generate the database configuration for deployment here"

	spongeTemplateVersionMark = "// todo generate the local sponge template code version here"

	configmapFileMark = "# todo generate server configuration code here"

	readmeFile    = "sponge/README.md"
	makeFile      = "sponge/Makefile"
	gitIgnoreFile = "sponge/.gitignore"

	startMarkStr  = "// delete the templates code start"
	endMarkStr    = "// delete the templates code end"
	startMark     = []byte(startMarkStr)
	endMark       = []byte(endMarkStr)
	wellStartMark = symbolConvert(startMarkStr)
	wellEndMark   = symbolConvert(endMarkStr)

	// embed FS template file when using
	selfPackageName = "github.com/zhufuyi/sponge"
)

var (
	ModelInitDBFile     = modelInitDBFile
	ModelInitDBFileMark = modelInitDBFileMark
	AppConfigFileDBMark = appConfigFileMark2
	StartMark           = startMark
	EndMark             = endMark
)

func symbolConvert(str string, additionalChar ...string) []byte {
	char := ""
	if len(additionalChar) > 0 {
		char = additionalChar[0]
	}

	return []byte(strings.Replace(str, "//", "#", 1) + char)
}

func convertServerName(serverName string) string {
	return strings.ReplaceAll(serverName, "-", "_")
}

func convertProjectAndServerName(projectName, serverName string) (pn string, sn string, err error) {
	if strings.HasSuffix(serverName, "-test") {
		err = fmt.Errorf(`the server name (%s) suffix "-test" is not supported for code generation, please delete suffix "-test" or change it to another name. `, serverName)
	}
	if strings.HasSuffix(serverName, "_test") {
		err = fmt.Errorf(`the server name (%s) suffix "_test" is not supported for code generation, please delete suffix "_test" or change it to another name. `, serverName)
	}

	sn = strings.ReplaceAll(serverName, "-", "_")
	pn = xstrings.ToKebabCase(projectName)
	return pn, sn, err
}

func adjustmentOfIDType(handlerCodes string, dbDriver string) string {
	if dbDriver == DBDriverMongodb {
		return idTypeToStr(handlerCodes)
	}
	return idTypeToUint64(idTypeFixToUint64(handlerCodes))
}

func idTypeFixToUint64(handlerCodes string) string {
	subStart := "ByIDRequest struct {"
	subEnd := "`" + `json:"id" binding:""` + "`"
	if subBytes := gofile.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID uint64 " + subEnd + " // uint64 id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func idTypeToUint64(handlerCodes string) string {
	subStart := "ObjDetail struct {"
	subEnd := "`" + `json:"id"` + "`"
	if subBytes := gofile.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID uint64 " + subEnd + " // convert to uint64 id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func idTypeToStr(handlerCodes string) string {
	subStart := "ObjDetail struct {"
	subEnd := "`" + `json:"id"` + "`"
	if subBytes := gofile.FindSubBytesNotIn([]byte(handlerCodes), []byte(subStart), []byte(subEnd)); len(subBytes) > 0 {
		old := subStart + string(subBytes) + subEnd
		newStr := subStart + "\n\tID string " + subEnd + " // convert to string id\n"
		handlerCodes = strings.ReplaceAll(handlerCodes, old, newStr)
	}

	return handlerCodes
}

func deleteFieldsMark(r replacer.Replacer, filename string, startMark []byte, endMark []byte) []replacer.Field {
	var fields []replacer.Field

	data, err := r.ReadFile(filename)
	if err != nil {
		//fmt.Printf("readFile error: %v\n", err)
		return fields
	}
	if subBytes := gofile.FindSubBytes(data, startMark, endMark); len(subBytes) > 0 {
		fields = append(fields,
			replacer.Field{ // clear marked template code
				Old: string(subBytes),
				New: "",
			},
		)
	}

	return fields
}

// DeleteCodeMark delete code mark fragment
func DeleteCodeMark(r replacer.Replacer, filename string, startMark []byte, endMark []byte) []replacer.Field {
	return deleteFieldsMark(r, filename, startMark, endMark)
}

func deleteAllFieldsMark(r replacer.Replacer, filename string, startMark []byte, endMark []byte) []replacer.Field {
	var fields []replacer.Field

	data, err := r.ReadFile(filename)
	if err != nil {
		//fmt.Printf("readFile error: %v\n", err)
		return fields
	}
	allSubBytes := gofile.FindAllSubBytes(data, startMark, endMark)
	for _, subBytes := range allSubBytes {
		fields = append(fields,
			replacer.Field{ // clear marked template code
				Old: string(subBytes),
				New: "",
			},
		)
	}

	return fields
}

func replaceFileContentMark(r replacer.Replacer, filename string, newContent string) []replacer.Field {
	var fields []replacer.Field

	data, err := r.ReadFile(filename)
	if err != nil {
		fmt.Printf("read the file \"%s\" error: %v\n", filename, err)
		return fields
	}

	fields = append(fields, replacer.Field{
		Old: string(data),
		New: newContent,
	})

	return fields
}

// resolving mirror repository host and name
func parseImageRepoAddr(addr string) (host string, name string) {
	splits := strings.Split(addr, "/")

	// default docker hub official repo address
	if len(splits) == 1 {
		return "https://index.docker.io/v1", addr
	}

	// unofficial repo address
	l := len(splits)
	return strings.Join(splits[:l-1], "/"), splits[l-1]
}

// ------------------------------------------------------------------------------------------

func parseProtobufFiles(protobufFile string) ([]string, bool, error) {
	if filepath.Ext(protobufFile) != ".proto" {
		return nil, false, fmt.Errorf("%v is not a protobuf file", protobufFile)
	}

	protobufFiles := gofile.FuzzyMatchFiles(protobufFile)
	countService, countImportTypes := 0, 0
	for _, file := range protobufFiles {
		protoData, err := os.ReadFile(file)
		if err != nil {
			return nil, false, err
		}
		if isExistServiceName(protoData) {
			countService++
		}
		if isDependImport(protoData, "api/types/types.proto") {
			countImportTypes++
		}
	}

	if countService == 0 {
		return nil, false, errors.New("not found service name, protobuf file requires at least one service")
	}

	return protobufFiles, countImportTypes > 0, nil
}

// save the moduleName and serverName to the specified file for external use
func saveGenInfo(moduleName string, serverName string, suitedMonoRepo bool, outputDir string) error {
	genInfo := moduleName + "," + serverName + "," + strconv.FormatBool(suitedMonoRepo)
	dir := outputDir + "/docs"
	_ = os.MkdirAll(dir, 0766)
	file := dir + "/gen.info"
	err := os.WriteFile(file, []byte(genInfo), 0666)
	if err != nil {
		return fmt.Errorf("save file %s error, %v", file, err)
	}
	return nil
}

func saveEmptySwaggerJSON(outputDir string) error {
	dir := outputDir + "/docs"
	_ = os.MkdirAll(dir, 0766)
	file := dir + "/apis.swagger.json"
	err := os.WriteFile(file, []byte(`{"swagger":"2.0","info":{"version":"version not set"}}`), 0666)
	if err != nil {
		return fmt.Errorf("save file %s error, %v", file, err)
	}
	return nil
}

// get moduleName and serverName from directory
func getNamesFromOutDir(dir string) (moduleName string, serverName string, suitedMonoRepo bool) {
	if dir == "" {
		return "", "", false
	}
	data, err := os.ReadFile(dir + "/docs/gen.info")
	if err != nil {
		return "", "", false
	}

	ms := strings.Split(string(data), ",")
	if len(ms) == 2 {
		return ms[0], ms[1], false
	} else if len(ms) >= 3 {
		return ms[0], ms[1], ms[2] == "true"
	}

	return "", "", false
}

func saveProtobufFiles(moduleName string, serverName string, suitedMonoRepo bool, outputDir string, protobufFiles []string) error {
	if suitedMonoRepo {
		outputDir = strings.TrimSuffix(outputDir, serverName)
		outputDir = strings.TrimSuffix(outputDir, gofile.GetPathDelimiter())
	}

	for _, pbFile := range protobufFiles {
		pbContent, err := os.ReadFile(pbFile)
		if err != nil {
			fmt.Printf("read file %s error, %v\n", pbFile, err)
			continue
		}
		pbContent = replacePackage(pbContent, moduleName, serverName)

		dir := outputDir + "/api/" + serverName + "/v1"
		_ = os.MkdirAll(dir, 0766)

		_, name := filepath.Split(pbFile)
		file := dir + "/" + name
		if gofile.IsExists(file) {
			return fmt.Errorf("file %s already exists", file)
		}
		err = os.WriteFile(file, pbContent, 0666)
		if err != nil {
			return fmt.Errorf("save file %s error, %v", file, err)
		}
	}

	return nil
}

func isExistServiceName(data []byte) bool {
	servicePattern := `\nservice (\w+)`
	re := regexp.MustCompile(servicePattern)
	matchArr := re.FindStringSubmatch(string(data))
	return len(matchArr) >= 2
}

func isDependImport(protoData []byte, pkgName string) bool {
	return bytes.Contains(protoData, []byte(pkgName))
}

func replacePackage(data []byte, moduleName string, serverName string) []byte {
	if bytes.Contains(data, []byte("\r\n")) {
		data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	}

	regStr := `\npackage [\w\W]*?;`
	reg := regexp.MustCompile(regStr)
	packageName := reg.Find(data)

	regStr2 := `go_package [\w\W]*?;\n`
	reg2 := regexp.MustCompile(regStr2)
	goPackageName := reg2.Find(data)

	if len(packageName) > 0 {
		newPackage := fmt.Sprintf("\npackage api.%s.v1;", serverName)
		data = bytes.Replace(data, packageName, []byte(newPackage), 1)
	}

	if len(goPackageName) > 0 {
		newGoPackage := fmt.Sprintf("go_package = \"%s/api/%s/v1;v1\";\n", moduleName, serverName)
		data = bytes.Replace(data, goPackageName, []byte(newGoPackage), 1)
	}

	return data
}

func getDBConfigCode(dbDriver string) string {
	dbConfigCode := ""
	switch strings.ToLower(dbDriver) {
	case DBDriverMysql, DBDriverTidb:
		dbConfigCode = mysqlConfigCode
	case DBDriverPostgresql:
		dbConfigCode = postgresqlConfigCode
	case DBDriverSqlite:
		dbConfigCode = sqliteConfigCode
	case DBDriverMongodb:
		dbConfigCode = mongodbConfigCode
	case undeterminedDBDriver:
		dbConfigCode = undeterminedDatabaseConfigCode
	default:
		panic("getDBConfigCode error, unsupported database driver: " + dbDriver)
	}
	return dbConfigCode
}

// GetDBConfigurationCode get db config code
func GetDBConfigurationCode(dbDriver string) string {
	return getDBConfigCode(dbDriver)
}

func getInitDBCode(dbDriver string) string {
	initDBCode := ""
	switch strings.ToLower(dbDriver) {
	case DBDriverMysql, DBDriverTidb:
		initDBCode = modelInitDBFileMysqlCode
	case DBDriverPostgresql:
		initDBCode = modelInitDBFilePostgresqlCode
	case DBDriverSqlite:
		initDBCode = modelInitDBFileSqliteCode
	case DBDriverMongodb:
		initDBCode = "" // do nothing
	default:
		panic("getInitDBCode error, unsupported database driver: " + dbDriver)
	}
	return initDBCode
}

// GetInitDataBaseCode get init db code
func GetInitDataBaseCode(dbDriver string) string {
	return getInitDBCode(dbDriver)
}

func getLocalSpongeTemplateVersion() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("os.UserHomeDir error:", err)
		return ""
	}

	versionFile := dir + "/.sponge/.github/version"
	data, err := os.ReadFile(versionFile)
	if err != nil {
		fmt.Printf("read file %s error: %v\n", versionFile, err)
		return ""
	}

	v := string(data)
	if v == "" {
		return ""
	}
	return fmt.Sprintf("github.com/zhufuyi/sponge %s", v)
}

func getEmbedTimeCode(isEmbed bool) string {
	if isEmbed {
		return embedTimeCode
	}
	return ""
}

func getExpectedSQLForDeletion(isEmbed bool) string {
	if !isEmbed {
		return strings.ReplaceAll(expectedSQLForDeletion, "UPDATE", "DELETE")
	}

	return expectedSQLForDeletion
}

func getExpectedSQLForDeletionField(isEmbed bool) []replacer.Field {
	var fields []replacer.Field
	esql := getExpectedSQLForDeletion(isEmbed)
	if esql != expectedSQLForDeletion {
		fields = append(fields, []replacer.Field{
			{
				Old: expectedSQLForDeletion,
				New: getExpectedSQLForDeletion(isEmbed),
			},
			{
				Old: "expectedArgsForDeletionTime := d.AnyTime",
				New: "",
			},
			{
				Old: "expectedArgsForDeletionTime := h.MockDao.AnyTime",
				New: "",
			},
			{
				Old: "WithArgs(expectedArgsForDeletionTime, testData.ID)",
				New: "WithArgs(testData.ID)",
			},
		}...)
	}
	return fields
}

func convertYamlConfig(configFile string) (string, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return "", err
	}
	defer f.Close() //nolint

	scanner := bufio.NewScanner(f)
	modifiedLines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		modifiedLine := "    " + line
		modifiedLines = append(modifiedLines, modifiedLine)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(modifiedLines, "\n"), nil
}

func generateConfigmap(serverName string, outPath string) error {
	configFile := fmt.Sprintf(outPath+"/configs/%s.yml", serverName)
	configmapFile := fmt.Sprintf(outPath+"/deployments/kubernetes/%s-configmap.yml", serverName)
	configFileData, err := convertYamlConfig(configFile)
	if err != nil {
		return err
	}
	configmapFileData, err := os.ReadFile(configmapFile)
	if err != nil {
		return err
	}
	data := strings.ReplaceAll(string(configmapFileData), configmapFileMark, configFileData)
	return os.WriteFile(configmapFile, []byte(data), 0666)
}

func sqliteDSNAdaptation(dbDriver string, dsn string) string {
	if dbDriver == DBDriverSqlite && gofile.IsWindows() {
		dsn = strings.Replace(dsn, "\\", "\\\\", -1)
		dsn = strings.Replace(dsn, "/", "\\\\", -1)
	}
	return dsn
}

func removeElements(slice []string, elements ...string) []string {
	if len(elements) == 0 {
		return slice
	}
	filters := make(map[string]struct{})
	for _, element := range elements {
		filters[element] = struct{}{}
	}
	result := make([]string, 0, len(slice)-1)
	for _, s := range slice {
		if _, ok := filters[s]; !ok {
			result = append(result, s)
		}
	}
	return result
}

func moveProtoFileToAPIDir(moduleName string, serverName string, suitedMonoRepo bool, outputDir string) error {
	apiDir := outputDir + gofile.GetPathDelimiter() + "api"
	protoFiles, _ := gofile.ListFiles(apiDir, gofile.WithNoAbsolutePath(), gofile.WithSuffix(".proto"))
	if err := saveProtobufFiles(moduleName, serverName, suitedMonoRepo, outputDir, protoFiles); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 100)
	_ = os.RemoveAll(apiDir)
	return nil
}

var (
	// for protoc.sh and protoc-doc.sh
	monoRepoAPIPath = `bash scripts/patch-mono.sh
cd ..

protoBasePath="api"`

	// for patch-mono.sh
	monoRepoHTTPPatch = `bash scripts/patch-mono.sh

HOST_ADDR=$1`

	// for patch.sh
	typePbShellCode = `
    if [ ! -d "../api/types" ]; then
        sponge patch gen-types-pb --out=./
        checkResult $?
        mv -f api/types ../api
        rmdir api
    fi`

	dupCodeMark = "--dir=internal/ecode"

	adaptDupCode = func(serverType string, serverName string) string {
		if serverType == codeNameHTTP {
			return dupCodeMark
		}
		return fmt.Sprintf("--dir=%s/internal/ecode", serverName)
	}
)

func serverCodeFields(serverType string, moduleName string, serverName string) []replacer.Field {
	return []replacer.Field{
		{
			Old: fmt.Sprintf("\"%s/internal/", moduleName),
			New: fmt.Sprintf("\"%s/internal/", moduleName+"/"+serverName),
		},
		{
			Old: "=$(cat docs/gen.info",
			New: fmt.Sprintf("=$(cat %s/docs/gen.info", serverName),
		},
		{
			Old: dupCodeMark,
			New: adaptDupCode(serverType, serverName),
		},
		{
			Old: fmt.Sprintf("\"%s/cmd/", moduleName),
			New: fmt.Sprintf("\"%s/cmd/", moduleName+"/"+serverName),
		},
		{
			Old: fmt.Sprintf("\"%s/configs", moduleName),
			New: fmt.Sprintf("\"%s/configs", moduleName+"/"+serverName),
		},
		{
			Old: fmt.Sprintf("\"%s/docs", moduleName),
			New: fmt.Sprintf("\"%s/docs", moduleName+"/"+serverName),
		},
		{
			Old: fmt.Sprintf("\"%s/api", moduleName),
			New: fmt.Sprintf("\"%s/api", moduleName+"/"+serverName),
		},
		{
			Old: "merge_file_name=docs/apis.json",
			New: fmt.Sprintf("merge_file_name=%s/docs/apis.json", serverName),
		},
		{
			Old: "--file=docs/apis.swagger.json",
			New: fmt.Sprintf("--file=%s/docs/apis.swagger.json", serverName),
		},
		{
			Old: "sponge merge http-pb",
			New: fmt.Sprintf("sponge merge http-pb --dir=%s", serverName),
		},
		{
			Old: "sponge merge rpc-pb",
			New: fmt.Sprintf("sponge merge rpc-pb --dir=%s", serverName),
		},
		{
			Old: "sponge merge rpc-gw-pb",
			New: fmt.Sprintf("sponge merge rpc-gw-pb --dir=%s", serverName),
		},
		{
			Old: "docs/apis.html",
			New: fmt.Sprintf("%s/docs/apis.html", serverName),
		},
		{
			Old: `sponge patch gen-types-pb --out=./`,
			New: typePbShellCode,
		},
		{
			Old: `protoBasePath="api"`,
			New: monoRepoAPIPath,
		},
		{
			Old: `HOST_ADDR=$1`,
			New: monoRepoHTTPPatch,
		},
		{
			Old: `genServerType=$1`,
			New: fmt.Sprintf(`genServerType="%s"`, serverType),
		},
		{
			Old: fmt.Sprintf("go get %s@", moduleName),
			New: fmt.Sprintf("go get %s@", "github.com/zhufuyi/sponge"),
		},
	}
}

// SubServerCodeFields sub server code fields
func SubServerCodeFields(moduleName string, serverName string) []replacer.Field {
	return []replacer.Field{
		{
			Old: fmt.Sprintf("\"%s/internal/", moduleName),
			New: fmt.Sprintf("\"%s/internal/", moduleName+"/"+serverName),
		},
		{
			Old: fmt.Sprintf("\"%s/configs", moduleName),
			New: fmt.Sprintf("\"%s/configs", moduleName+"/"+serverName),
		},
		{
			Old: fmt.Sprintf("\"%s/api", moduleName),
			New: fmt.Sprintf("\"%s/api", moduleName+"/"+serverName),
		},
	}
}

func changeOutPath(outPath string, serverName string) string {
	switch outPath {
	case "", ".", "./", ".\\", serverName, "./" + serverName, ".\\" + serverName:
		return serverName
	}
	return outPath + gofile.GetPathDelimiter() + serverName
}

func getSubFiles(selectFiles map[string][]string, replaceFiles map[string][]string) []string {
	files := []string{}
	for dir, filenames := range selectFiles {
		if v, ok := replaceFiles[dir]; ok {
			filenames = v
		}
		for _, filename := range filenames {
			files = append(files, dir+"/"+filename)
		}
	}
	return files
}

func getLocalGoVersion() string {
	result, err := gobash.Exec("go", "version")
	if err != nil {
		return defaultGoModVersion
	}

	pattern := `go(\d+\.\d+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(string(result))
	if len(matches) < 2 {
		return defaultGoModVersion
	}

	localGoVersion := "go " + matches[1]
	if localGoVersion < defaultGoModVersion {
		return defaultGoModVersion
	}

	if len(localGoVersion) != 6 && len(localGoVersion) != 7 {
		return defaultGoModVersion
	}

	return localGoVersion
}

func dbDriverErr(driver string) error {
	return errors.New("unsupported db driver: " + driver)
}

func flagTip(name ...string) string {
	if len(name) == 2 {
		return fmt.Sprintf("if you specify the directory where the web or microservice generated by sponge, the %s and %s flag can be ignored", name[0], name[1])
	}
	return fmt.Sprintf("if you specify the directory where the web or microservice generated by sponge, the %s flag can be ignored", name[0])
}

func cutPath(srcFilePath string) string {
	dirPath, _ := filepath.Abs(".")
	srcFilePath = strings.ReplaceAll(srcFilePath, dirPath, ".")
	return strings.ReplaceAll(srcFilePath, "\\", "/")
}

func wrapPoint(s string) string {
	return "`" + s + "`"
}

func setReadmeTitle(moduleName string, serverName string, serverType string, suitedMonoRepo bool) string {
	var repoType string
	if suitedMonoRepo {
		repoType = "mono-repo"
	} else {
		if serverType == codeNameHTTP {
			repoType = "monolith"
		} else {
			repoType = "multi-repo"
		}
	}

	return wellPrefix + serverName + fmt.Sprintf(`

| Feature             | Value          |
| :----------------: | :-----------: |
| Server name      |  %s   |
| Server type        |  %s   |
| Go module name |  %s  |
| Repository type   |  %s  |

`, wrapPoint(serverName), wrapPoint(serverType), wrapPoint(moduleName), wrapPoint(repoType))
}

// GetGoModFields get go mod fields
func GetGoModFields(moduleName string) []replacer.Field {
	return []replacer.Field{
		{
			Old: "github.com/zhufuyi/sponge",
			New: moduleName,
		},
		{
			Old: defaultGoModVersion,
			New: getLocalGoVersion(),
		},
		{
			Old: spongeTemplateVersionMark,
			New: getLocalSpongeTemplateVersion(),
		},
	}
}
