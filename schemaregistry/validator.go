package schemaregistry

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

var schemas = map[string]string{
	"tasks.created":   "task/created",
	"tasks.shuffled":  "task/shuffled",
	"tasks.completed": "task/completed",
}

func Validate(eventName string, data []byte, version int) (bool, error) {
	schema, ok := schemas[eventName]
	if !ok {
		return false, fmt.Errorf("schema for event \"%s\" was not found", eventName)
	}

	// ToDo: don't want to spend time for having an absolute path.
	schemaPath := fmt.Sprintf(
		"file:///Users/v.eizele/Documents/Projects/async-architecture-course/schemaregistry/schema/%s/%v.json",
		schema, version,
	)
	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)
	dataLoader := gojsonschema.NewBytesLoader(data)
	res, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return false, err
	}
	if !res.Valid() {
		fmt.Println(res.Errors())
		return false, nil
	}

	return true, nil
}
