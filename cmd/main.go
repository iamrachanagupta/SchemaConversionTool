package main

import (
	"SchemaConversionTool/internal/domain"
	_ "encoding/xml"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/spark-schema/json", domain.SparkSchemaHandler)
	http.HandleFunc("/spark-schema/xml", domain.SparkSchemaHandler)
	http.Handle("/", http.FileServer(http.Dir("."))) // Serve static files (like the HTML file)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

//// Handler for Spark schema endpoint
//func jsonsparkSchemaHandler(w http.ResponseWriter, r *http.Request) {
//	schemaBytes, err := getSchemaFromPath("inputSchemas/json_schema.json")
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	var schema map[string]interface{}
//
//	// Unmarshal JSON data into the map
//	err = json.Unmarshal(schemaBytes, &schema)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	// Simulated Spark schema data
//	sparkSchema, _ := domain.ConvertJSONSchemaToSparkSchema(schema)
//	// Create a response struct
//	response := SchemaResponse{
//		OriginalSchema: schema,
//		SparkSchema:    sparkSchema,
//	}
//	//	// Marshal the schema to JSON with indentation
//	//	schemaJSON, err := json.MarshalIndent(response, "", "    ")
//	//	if err != nil {
//	//		fmt.Println("Error marshaling schema to JSON:", err)
//	//		return
//	//	}
//	//
//	//	// Print the pretty formatted JSON schema
//	//	//fmt.Println(string(schemaJSON))
//	//	// Set the content type header
//	//	w.Header().Set("Content-Type", "application/json")
//	//	// Write the Spark schema JSON response
//	//	w.Write(schemaJSON)
//	//
//	//}
//
//	// Load HTML template
//	tmpl, err := template.New("schema").Funcs(template.FuncMap{
//		"toJSON": func(v interface{}) (string, error) {
//			jsonData, err := json.MarshalIndent(v, "", "    ")
//			if err != nil {
//				return "", err
//			}
//			return string(jsonData), nil
//		},
//	}).Parse(`
//<!DOCTYPE html>
//<html lang="en">
//<head>
//    <meta charset="UTF-8">
//    <meta name="viewport" content="width=device-width, initial-scale=1.0">
//    <title>Schema Conversion</title>
//    <style>
//        .container {
//            display: flex;
//        }
//        .schema {
//            flex: 1;
//            margin: 10px;
//            border: 1px solid #ccc;
//            padding: 10px;
//            overflow: auto;
//        }
//    </style>
//</head>
//<body>
//    <div class="container">
//        <div class="schema">
//            <h2>Original Schema</h2>
//            <pre>{{ toJSON .OriginalSchema }}</pre>
//        </div>
//        <div class="schema">
//            <h2>Spark Schema</h2>
//            <pre>{{ toJSON .SparkSchema }}</pre>
//        </div>
//    </div>
//</body>
//</html>
//`)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Execute the template with the response data
//	err = tmpl.Execute(w, response)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//// Handler for Spark schema endpoint
//func xmlSparkSchemaHandler(w http.ResponseWriter, r *http.Request) {
//	originalSchema, schemaBytes := convertXMLToJSON()
//	var schema map[string]interface{}
//
//	// Unmarshal JSON data into the map
//	err := json.Unmarshal(schemaBytes, &schema)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	// Simulated Spark schema data
//	sparkSchema, _ := domain.ConvertJSONSchemaToSparkSchema(schema["root"].(map[string]interface{}))
//	// Create a response struct
//	response := SchemaResponse{
//		OriginalSchema: string(originalSchema),
//		SparkSchema:    sparkSchema,
//	}
//	// Load HTML template
//	tmpl, err := template.New("schema").Funcs(template.FuncMap{
//		"toJSON": func(v interface{}) (string, error) {
//			jsonData, err := json.MarshalIndent(v, "", "    ")
//			if err != nil {
//				return "", err
//			}
//			return string(jsonData), nil
//		},
//	}).Parse(`
//<!DOCTYPE html>
//<html lang="en">
//<head>
//    <meta charset="UTF-8">
//    <meta name="viewport" content="width=device-width, initial-scale=1.0">
//    <title>Schema Conversion</title>
//    <style>
//        .container {
//            display: flex;
//        }
//        .schema {
//            flex: 1;
//            margin: 10px;
//            border: 1px solid #ccc;
//            padding: 10px;
//            overflow: auto;
//        }
//    </style>
//</head>
//<body>
//    <div class="container">
//        <div class="schema">
//            <h2>Original Schema</h2>
//            <pre>{{ .OriginalSchema }}</pre>
//        </div>
//        <div class="schema">
//            <h2>Spark Schema</h2>
//            <pre>{{ toJSON .SparkSchema }}</pre>
//        </div>
//    </div>
//</body>
//</html>
//`)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Execute the template with the response data
//	err = tmpl.Execute(w, response)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
