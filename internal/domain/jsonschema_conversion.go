// (C) Copyright 2023-2024 Hewlett Packard Enterprise Development LP

package domain

import "fmt"

const (
	ObjectType = "object"
	ArrayType  = "array"
)

// setMetadataFields sets metadata fields dynamically based on propMap
func setMetadataFields(metadata, propMap map[string]interface{}) {
	for key, value := range propMap {
		if key != "type" && key != "items" && key != "properties" && key != "required" {
			metadata[key] = value
		}
	}
}

// ConvertJSONSchemaToSparkSchema recursively converts JSON schema to Spark schema.
func ConvertJSONSchemaToSparkSchema(jsonSchema map[string]interface{}) (map[string]interface{}, error) {
	sparkSchema := make(map[string]interface{})
	sparkSchema["type"] = "struct"

	var fields []map[string]interface{}

	properties, ok := jsonSchema["properties"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to map properties in JSON schema")
	}

	for propName, propData := range properties {
		field := make(map[string]interface{})
		field["name"] = propName

		propMap, ok := propData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to map data fields in JSON schema")
		}

		fieldType, ok := propMap["type"].(string)
		if !ok {
			return nil, fmt.Errorf("failed to map type field in JSON schema")
		}

		field["type"] = fieldType
		field["nullable"] = !contains(jsonSchema["required"], propName)
		field["metadata"] = make(map[string]interface{})

		// Set all metadata fields dynamically
		setMetadataFields(field["metadata"].(map[string]interface{}), propMap)

		if fieldType == ArrayType {
			if items, ok := propMap["items"].([]map[string]interface{}); ok {
				elementTypes := make([]interface{}, 0)
				for _, item := range items {
					itemMap := item
					if itemMap["type"].(string) == ObjectType {
						var err error
						field["type"], err = ConvertJSONSchemaToSparkSchema(itemMap)
						if err != nil {
							return nil, fmt.Errorf("converting the %s property: %w", propName, err)
						}
						elementTypes = append(elementTypes, field["type"])
					} else if itemType, isString := itemMap["type"].(string); isString {
						elementTypes = append(elementTypes, itemType)
					}
				}
				if len(elementTypes) > 0 {
					field["type"] = map[string]interface{}{
						"type":         ArrayType,
						"elementType":  elementTypes,
						"containsNull": true,
					}
				}
			}
		}

		if fieldType == ObjectType {
			var err error
			field["type"], err = ConvertJSONSchemaToSparkSchema(propMap)
			if err != nil {
				return nil, fmt.Errorf("converting the %s property: %w", propName, err)
			}
		}

		fields = append(fields, field)
	}

	sparkSchema["fields"] = fields
	return sparkSchema, nil
}

// contains checks if a value is present in a slice.
func contains(slice, item interface{}) bool {
	if slice, ok := slice.([]string); ok {
		for _, i := range slice {
			if i == item {
				return true
			}
		}
	}
	return false
}
