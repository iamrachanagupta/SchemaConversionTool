## Steps to use the protoc-gen-jsonschema library

- Follow https://github.com/chrusty/protoc-gen-jsonschema
  - go install github.com/chrusty/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema@latest
  - sudo cp ~/go/bin/protoc-gen-jsonschema /usr/local/bin/protoc-gen-jsonschema
  - protoc-gen-jsonschema --version
  - Check if this command works `protoc --jsonschema_out=./inputSchemas/generatedByProto/ --proto_path=inputSchemas/ inputSchemas/proto_schema.proto`
- go run cmd/main.go
  * Go to http://localhost:8080/spark-schema/avro
  * Go to http://localhost:8080/spark-schema/xml
  * Go to http://localhost:8080/spark-schema/json
  * Go to http://localhost:8080/spark-schema/proto