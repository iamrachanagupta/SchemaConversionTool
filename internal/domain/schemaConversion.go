package domain

import (
	"SchemaConversionTool/internal/adapters"
	"encoding/json"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"os"
	"strings"
)

func ConvertInputSchemaToSparkSchema(SchemaType, path string) (SchemaResponse, error) {
	var schema, sparkSchema map[string]interface{}
	var originalSchema interface{}
	var response SchemaResponse
	originalSchemaBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading XML file:", err)
		return SchemaResponse{}, err
	}

	if SchemaType == "XML" {
		xml := strings.NewReader(string(originalSchemaBytes))

		jsonBuff, err := xj.Convert(xml)
		if err != nil {
			return SchemaResponse{}, err
		}
		originalSchema = string(originalSchemaBytes)
		originalSchemaBytes = jsonBuff.Bytes()
	}
	// Unmarshal JSON data into the map
	err = json.Unmarshal(originalSchemaBytes, &schema)
	if err != nil {
		fmt.Println("Error:", err)
		return SchemaResponse{}, err
	}

	if SchemaType == "XML" {
		schema = schema["root"].(map[string]interface{})
	} else {
		originalSchema = schema
	}

	// Simulated Spark schema data
	sparkSchema, err = adapters.ConvertJSONSchemaToSparkSchema(schema)
	if err != nil {
		fmt.Println("Error:", err)
		return SchemaResponse{}, err
	}

	// Create a response struct
	response = SchemaResponse{
		SchemaType:     SchemaType,
		OriginalSchema: originalSchema,
		SparkSchema:    sparkSchema,
	}
	return response, err

}
