package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/go-dev-frame/sponge/pkg/gofile"
)

// ConvertSwagger2ToOpenAPI3Command convert swagger2 to openapi3
func ConvertSwagger2ToOpenAPI3Command() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "swagger2-to-openapi3",
		Short: "Convert swagger2.0 to openapi3",
		Long:  "Convert swagger2.0 to openapi3.",
		Example: color.HiBlackString(`  # Convert swagger2.0 files to openapi3
  sponge web swagger2-to-openapi3 --file=swagger.json`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return convertOpenAPI3(file)
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "input json file path")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func convertOpenAPI3(inputFile string) error {
	if gofile.GetFileSuffixName(inputFile) != ".json" {
		return fmt.Errorf("input file must be a json file")
	}

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	outputYAML, outputJSON := getOutputFile(inputFile)

	var swaggerDoc openapi2.T
	if err = json.Unmarshal(data, &swaggerDoc); err != nil {
		return fmt.Errorf("parse swagger json file failed: %v", err)
	}

	openapi3Doc, err := openapi2conv.ToV3(&swaggerDoc)
	if err != nil {
		return fmt.Errorf("convert to openapi3 failed: %v", err)
	}

	jsonData, err := json.MarshalIndent(openapi3Doc, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize to json failed: %v", err)
	}
	if err = os.WriteFile(outputJSON, jsonData, 0644); err != nil {
		return fmt.Errorf("write json file failed: %v", err)
	}

	yamlData, err := yaml.Marshal(openapi3Doc)
	if err != nil {
		return fmt.Errorf("serialize to yaml failed: %v", err)
	}
	if err = os.WriteFile(outputYAML, yamlData, 0644); err != nil {
		return fmt.Errorf("write yaml file failed: %v", err)
	}

	fmt.Printf("conversion successful, output files: %s, %s\n", outputJSON, outputYAML)

	return nil
}

func getOutputFile(filePath string) (yamlFile string, jsonFile string) {
	var suffix string
	if strings.HasSuffix(filePath, "swagger.json") {
		suffix = "swagger.json"

	} else {
		suffix = gofile.GetFilename(filePath)
	}
	yamlFile = strings.TrimSuffix(filePath, suffix) + "openapi3.yaml"
	jsonFile = strings.TrimSuffix(filePath, suffix) + "openapi3.json"
	return yamlFile, jsonFile
}
