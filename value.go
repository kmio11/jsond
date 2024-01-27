package jsond

import (
	"encoding/json"
	"fmt"
)

// jsonvalue represents a json.Unmarshal result.
// (https://pkg.go.dev/encoding/json#Unmarshal)
// It stores one of these in the any value:
// - bool, for JSON booleans
// - float64, for JSON numbers
// - string, for JSON strings
// - []any, for JSON arrays
// - map[string]any, for JSON objects
// - nil for JSON null
type jsonvalue any

func getJSONValue(v any) (jsonvalue, error) {
	if v == nil {
		return v, nil
	}

	switch t := v.(type) {
	case *Node:
		return t.value, t.err
	case bool, float64, string:
		return v, nil
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		jv := *new(jsonvalue)
		err = json.Unmarshal(data, &jv)
		return jv, err
	}
}

func getTypeString(v jsonvalue) string {
	switch v.(type) {
	case bool:
		return "bool"
	case float64:
		return "number"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	case nil:
		return "nil"
	default:
		panic(fmt.Sprintf("invalid jsonvalue. v=%v", v))
	}
}
