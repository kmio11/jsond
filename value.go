package jsond

import "encoding/json"

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
