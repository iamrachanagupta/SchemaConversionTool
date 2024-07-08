// (C) Copyright 2023-2024 Hewlett Packard Enterprise Development LP
package adapters

import (
	"fmt"
)

// AvroToSparkType converts an Avro type to a Spark type
func AvroToSparkType(avroType interface{}) (string, error) {
	switch avroType := avroType.(type) {
	case string:
		switch avroType {
		case "string":
			return "StringType", nil
		case "int":
			return "IntegerType", nil
		case "long":
			return "LongType", nil
		case "float":
			return "FloatType", nil
		case "double":
			return "DoubleType", nil
		case "boolean":
			return "BooleanType", nil
		case "bytes":
			return "BinaryType", nil
		default:
			return "", fmt.Errorf("unsupported Avro type: %s", avroType)
		}
	case map[string]interface{}:
		if avroType["type"] == "array" {
			itemType, err := AvroToSparkType(avroType["items"])
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("ArrayType(%s)", itemType), nil
		}
		if avroType["type"] == "record" {
			return "StructType", nil // Will be handled in the struct conversion
		}
		return "", fmt.Errorf("unsupported Avro complex type: %v", avroType)
	default:
		return "", fmt.Errorf("unsupported Avro type format: %v", avroType)
	}
}

// ConvertRecordSchema converts an Avro record schema to a Spark schema
func ConvertRecordSchema(avroSchema map[string]interface{}) (map[string]interface{}, error) {
	fields := avroSchema["fields"].([]interface{})
	sparkFields := make([]map[string]interface{}, len(fields))

	for i, field := range fields {
		fieldMap := field.(map[string]interface{})
		fieldName := fieldMap["name"].(string)
		fieldType, err := AvroToSparkType(fieldMap["type"])
		if err != nil {
			return nil, err
		}
		sparkFields[i] = map[string]interface{}{
			"name": fieldName,
			"type": fieldType,
		}
	}

	return map[string]interface{}{
		"type":   "struct",
		"fields": sparkFields,
	}, nil
}
