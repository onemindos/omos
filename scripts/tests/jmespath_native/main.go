// Exercise the exact JMESPath dependency pinned by Kyverno 1.19.1. Kubernetes
// dynamic clients retain native int64 values; JSON-loaded CLI resources use float64.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"

	jmespath "github.com/kyverno/go-jmespath"
)

type input struct {
	Queries  []string                 `json:"queries"`
	Contexts []map[string]interface{} `json:"contexts"`
}

type result struct {
	Value interface{} `json:"value"`
	Error string      `json:"error,omitempty"`
}

func nativeNumbers(value interface{}) interface{} {
	switch value := value.(type) {
	case float64:
		if value == math.Trunc(value) {
			return int64(value)
		}
		return value
	case map[string]interface{}:
		converted := make(map[string]interface{}, len(value))
		for key, child := range value {
			converted[key] = nativeNumbers(child)
		}
		return converted
	case []interface{}:
		converted := make([]interface{}, len(value))
		for index, child := range value {
			converted[index] = nativeNumbers(child)
		}
		return converted
	default:
		return value
	}
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var request input
	if err := json.Unmarshal(data, &request); err != nil {
		panic(err)
	}
	output := make(map[string][][]result)
	for _, mode := range []string{"json-float64", "kubernetes-int64"} {
		rows := make([][]result, len(request.Contexts))
		for index, context := range request.Contexts {
			var value interface{} = context
			if mode == "kubernetes-int64" {
				value = nativeNumbers(context)
			}
			for _, query := range request.Queries {
				evaluated, err := jmespath.Search(query, value)
				r := result{Value: evaluated}
				if err != nil {
					r.Error = err.Error()
				}
				rows[index] = append(rows[index], r)
			}
		}
		output[mode] = rows
	}
	encoded, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(encoded))
}
