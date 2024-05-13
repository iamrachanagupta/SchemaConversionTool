package domain

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os/exec"
)

// SchemaResponse holds both the original JSON schema and the converted Spark schema
type SchemaResponse struct {
	SchemaType     string      `json:"Schema_type"`
	OriginalSchema interface{} `json:"original_schema"`
	SparkSchema    interface{} `json:"spark_schema"`
}

func generateJSONFromProto() {
	cmd := exec.Command("protoc", "--jsonschema_out=./inputSchemas/generatedByProto/",
		"--proto_path=inputSchemas/", "inputSchemas/proto_schema.proto")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func SparkSchemaHandler(w http.ResponseWriter, r *http.Request) {
	var schemaType, path string
	var response SchemaResponse
	var err error
	// Check the URL path to determine the requested format
	switch r.URL.Path {
	case "/spark-schema/json":
		schemaType = "JSON"
		path = "inputSchemas/json_schema.json"
	case "/spark-schema/xml":
		schemaType = "XML"
		path = "inputSchemas/xml_schema.xml"
	case "/spark-schema/proto":
		schemaType = "PROTO"
		path = "inputSchemas/proto_schema.proto"
	case "/spark-schema/custom":
		schemaType = "CUSTOM"
		path = "inputSchemas/custom_json_schema.json"
	default:
		http.Error(w, "Unsupported format", http.StatusNotFound)
	}
	response, err = ConvertInputSchemaToSparkSchema(schemaType, path)
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
			{{ if eq .SchemaType "JSON" }}
            <pre>{{ toJSON .OriginalSchema }}</pre>
			{{ else }}
			<pre>{{ .OriginalSchema }}</pre>
			{{ end }}
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
