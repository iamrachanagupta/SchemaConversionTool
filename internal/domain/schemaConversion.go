package domain

import (
	"SchemaConversionTool/internal/adapters"
	"encoding/json"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"log"
	"os"
	"strings"
)

func ConvertInputSchemaToSparkSchema(SchemaType, path string) (SchemaResponse, error) {
	var schema, sparkSchema map[string]interface{}
	var originalSchema interface{}
	var response SchemaResponse

	if SchemaType == "PROTO" {
		protoSchema, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Error reading XML file:", err)
			return SchemaResponse{}, err
		}
		originalSchema = string(protoSchema)

		generateJSONFromProto()
		path = "inputSchemas/generatedByProto/schema.json"
	}
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
	} else if SchemaType == "PROTO" {
		schema = schema["definitions"].(map[string]interface{})["schema"].(map[string]interface{})
	} else if SchemaType == "AVRO" {
		originalSchema = string(originalSchemaBytes)
	} else {
		originalSchema = schema
	}

	if SchemaType == "AVRO" {
		sparkSchema, err = adapters.ConvertRecordSchema(schema)
		if err != nil {
			log.Fatalf("Failed to convert Avro schema to Spark schema: %v", err)
		}
	} else {
		// Simulated Spark schema data
		sparkSchema, err = adapters.ConvertJSONSchemaToSparkSchema(schema)
		if err != nil {
			fmt.Println("Error:", err)
			return SchemaResponse{}, err
		}
	}

	// Create a response struct
	response = SchemaResponse{
		SchemaType:     SchemaType,
		OriginalSchema: originalSchema,
		SparkSchema:    sparkSchema,
	}
	return response, err

}
