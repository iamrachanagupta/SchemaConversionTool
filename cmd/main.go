package main

import (
	"SchemaConversionTool/internal/domain"
	"encoding/json"
	_ "encoding/xml"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// SchemaResponse holds both the original JSON schema and the converted Spark schema
type SchemaResponse struct {
	OriginalSchema interface{} `json:"original_schema"`
	SparkSchema    interface{} `json:"spark_schema"`
}

func getSchemaFromPath(path string) ([]byte, error) {
	var Schema, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Schema, nil
}

//var Schema = map[string]interface{}{
//	"schema": "http://json-schema.org/draft-07/schema#",
//	"type":   "object",
//	"properties": map[string]interface{}{
//		"type":             map[string]interface{}{"type": "string"},
//		"id":               map[string]interface{}{"type": "string"},
//		"customerId":       map[string]interface{}{"type": "string"},
//		"createdAt":        map[string]interface{}{"type": "string"},
//		"updatedAt":        map[string]interface{}{"type": "string"},
//		"name":             map[string]interface{}{"type": "string"},
//		"associationCount": map[string]interface{}{"type": "integer", "minimum": 0},
//	},
//	"required": []string{"type", "id", "customerId", "createdAt", "updatedAt", "name"},
//}

func main() {

	http.HandleFunc("/json_to_spark", jsonsparkSchemaHandler)
	http.HandleFunc("/xml_to_spark", xmlSparkSchemaHandler)
	http.Handle("/", http.FileServer(http.Dir("."))) // Serve static files (like the HTML file)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

// Handler for Spark schema endpoint
func jsonsparkSchemaHandler(w http.ResponseWriter, r *http.Request) {
	schemaBytes, err := getSchemaFromPath("inputSchemas/json_schema.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var schema map[string]interface{}

	// Unmarshal JSON data into the map
	err = json.Unmarshal(schemaBytes, &schema)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Simulated Spark schema data
	sparkSchema, _ := domain.ConvertJSONSchemaToSparkSchema(schema)
	// Create a response struct
	response := SchemaResponse{
		OriginalSchema: schema,
		SparkSchema:    sparkSchema,
	}
	//	// Marshal the schema to JSON with indentation
	//	schemaJSON, err := json.MarshalIndent(response, "", "    ")
	//	if err != nil {
	//		fmt.Println("Error marshaling schema to JSON:", err)
	//		return
	//	}
	//
	//	// Print the pretty formatted JSON schema
	//	//fmt.Println(string(schemaJSON))
	//	// Set the content type header
	//	w.Header().Set("Content-Type", "application/json")
	//	// Write the Spark schema JSON response
	//	w.Write(schemaJSON)
	//
	//}

	// Load HTML template
	tmpl, err := template.New("schema").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) (string, error) {
			jsonData, err := json.MarshalIndent(v, "", "    ")
			if err != nil {
				return "", err
			}
			return string(jsonData), nil
		},
	}).Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Schema Conversion</title>
    <style>
        .container {
            display: flex;
        }
        .schema {
            flex: 1;
            margin: 10px;
            border: 1px solid #ccc;
            padding: 10px;
            overflow: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="schema">
            <h2>Original Schema</h2>
            <pre>{{ toJSON .OriginalSchema }}</pre>
        </div>
        <div class="schema">
            <h2>Spark Schema</h2>
            <pre>{{ toJSON .SparkSchema }}</pre>
        </div>
    </div>
</body>
</html>
`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the response data
	err = tmpl.Execute(w, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler for Spark schema endpoint
func xmlSparkSchemaHandler(w http.ResponseWriter, r *http.Request) {
	originalSchema, schemaBytes := convertXMLToJSON()
	var schema map[string]interface{}

	// Unmarshal JSON data into the map
	err := json.Unmarshal(schemaBytes, &schema)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Simulated Spark schema data
	sparkSchema, _ := domain.ConvertJSONSchemaToSparkSchema(schema["root"].(map[string]interface{}))
	// Create a response struct
	response := SchemaResponse{
		OriginalSchema: string(originalSchema),
		SparkSchema:    sparkSchema,
	}
	// Load HTML template
	tmpl, err := template.New("schema").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) (string, error) {
			jsonData, err := json.MarshalIndent(v, "", "    ")
			if err != nil {
				return "", err
			}
			return string(jsonData), nil
		},
	}).Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Schema Conversion</title>
    <style>
        .container {
            display: flex;
        }
        .schema {
            flex: 1;
            margin: 10px;
            border: 1px solid #ccc;
            padding: 10px;
            overflow: auto;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="schema">
            <h2>Original Schema</h2>
            <pre>{{ .OriginalSchema }}</pre>
        </div>
        <div class="schema">
            <h2>Spark Schema</h2>
            <pre>{{ toJSON .SparkSchema }}</pre>
        </div>
    </div>
</body>
</html>
`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the response data
	err = tmpl.Execute(w, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func convertXMLToJSON() ([]byte, []byte) {

	// Read XML content
	xmlData, err := os.ReadFile("inputSchemas/xml_schema.xml")
	if err != nil {
		fmt.Println("Error reading XML file:", err)
		return nil, nil
	}
	xml := strings.NewReader(string(xmlData))

	jsonBytes, err := xj.Convert(xml)
	if err != nil {
		return xmlData, nil
	}
	return xmlData, jsonBytes.Bytes()
}
