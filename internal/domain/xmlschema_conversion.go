package domain

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type XMLType struct {
	Type string `xml:"type"`
}

type XMLProperty struct {
	Type     string            `xml:"type"`
	Attrs    map[string]string `xml:",any,attr"`
	InnerXML string            `xml:",innerxml"`
}

type XMLRoot struct {
	XMLName    xml.Name      `xml:"root"`
	Type       XMLType       `xml:"type"`
	Properties []XMLProperty `xml:"properties"`
	Required   []string      `xml:"required"`
}

type JSONField struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Name     string                 `json:"name"`
	Nullable bool                   `json:"nullable"`
	Type     string                 `json:"type"`
}

type JSONObject struct {
	Fields []JSONField `json:"fields"`
	Type   string      `json:"type"`
}

func ConvertXMLSchemaToSparkSchema(xmlSchema []byte) ([]byte, error) {
	var root XMLRoot
	err := xml.Unmarshal(xmlSchema, &root)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return nil, err
	}
	fmt.Println(string(xmlSchema))
	fmt.Println(root.Properties)

	// Create a map to store required fields
	requiredFields := make(map[string]bool)
	for _, required := range root.Required {
		requiredFields[required] = true
	}
	fmt.Println(len(root.Properties))

	jsonFields := make([]JSONField, 0)

	for _, prop := range root.Properties {
		metadata := make(map[string]interface{})
		for k, v := range prop.Attrs {
			metadata[k] = v
		}
		jsonFields = append(jsonFields, JSONField{
			Metadata: metadata,
			Name:     prop.Type,
			Nullable: requiredFields[prop.Type], // Set nullable based on if the field is required
			Type:     prop.Type,
		})
	}

	// Create JSON object
	jsonObj := JSONObject{
		Fields: jsonFields,
		Type:   root.Type.Type,
	}

	// Convert JSON object to JSON string
	jsonBytes, err := json.MarshalIndent(jsonObj, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}
	return jsonBytes, nil
}

//// Root represents the root element of the XML schema.
//type Root struct {
//	XMLName    xml.Name   `xml:"root"`
//	Type       string     `xml:"type"`
//	Properties []Property `xml:"properties>property"`
//	Required   []string   `xml:"required"`
//}
//
//// Property represents a property in the XML schema.
//type Property struct {
//	Name    string `xml:"name"`
//	Type    Type   `xml:"type"`
//	Minimum int    `xml:"minimum,omitempty"`
//}
//
//// Type represents property in the XML schema.
//type Type struct {
//	Type    string `xml:"type"`
//	Minimum int    `xml:"minimum,omitempty"`
//}
//
//// ConvertXMLSchemaToSparkSchema recursively converts XML schema to Spark schema.
//func ConvertXMLSchemaToSparkSchema(xmlSchema []byte) (map[string]interface{}, error) {
//	var root Root
//	err := xml.Unmarshal(xmlSchema, &root)
//	if err != nil {
//		return nil, fmt.Errorf("failed to unmarshal XML schema: %w", err)
//	}
//
//	fmt.Println(string(xmlSchema))
//	sparkSchema := make(map[string]interface{})
//	sparkSchema["type"] = "struct"
//
//	var fields []map[string]interface{}
//
//	for _, prop := range root.Properties {
//		field := make(map[string]interface{})
//		field["name"] = prop.Name
//
//		field["nullable"] = !contain(root.Required, prop.Name)
//		field["metadata"] = make(map[string]interface{})
//
//		// Set all metadata fields dynamically
//		setMetadataFields(field["metadata"].(map[string]interface{}), prop)
//
//		fieldType := xmlTypeToSparkType(prop.Type.Type)
//		field["type"] = fieldType
//
//		fields = append(fields, field)
//	}
//
//	sparkSchema["fields"] = fields
//	return sparkSchema, nil
//}
//
//// xmlTypeToSparkType maps XML types to Spark data types.
//func xmlTypeToSparkType(xmlType string) string {
//	switch xmlType {
//	case "string":
//		return "StringType"
//	case "integer":
//		return "IntegerType"
//	// Add more mappings as needed
//	default:
//		return "UnknownType"
//	}
//}
//
//// contain checks if a value is present in a slice.
//func contain(slice []string, item string) bool {
//	for _, i := range slice {
//		if i == item {
//			return true
//		}
//	}
//	return false
//}
