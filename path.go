package jsond

import "fmt"

// jsonpath represents json jsonpath.
type jsonpath []property

type (
	// property represents path's element.
	// - int, for array index
	// - string, for object key
	property any

	arrayIndex int
	objectKey  string
)

func (p jsonpath) String() string {
	if len(p) == 0 {
		return ""
	}
	joined := "$"
	for _, prop := range p {
		switch t := prop.(type) {
		case arrayIndex:
			joined = fmt.Sprintf("%s[%d]", joined, t)
		case objectKey:
			joined = fmt.Sprintf("%s['%s']", joined, t)
		default:
			panic(fmt.Sprintf("invalid property. prop=%v, type=%t", prop, prop))
		}
	}
	return joined
}

func (p jsonpath) append(prop property) jsonpath {
	switch t := prop.(type) {
	case arrayIndex:
		return append(p, t)
	case objectKey:
		return append(p, t)
	default:
		panic(fmt.Sprintf("invalid property. prop=%v, type=%t", prop, prop))
	}
}

func getProperty(v any) (property, error) {
	switch t := v.(type) {
	case
		int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64:

		idx, err := getInt(t)
		if err != nil {
			return nil, fmt.Errorf("invalid index : %v", idx)
		}
		return arrayIndex(idx), nil

	case string:
		return objectKey(t), nil

	default:
		return nil, fmt.Errorf("invalid property : %v", v)
	}
}

func getInt(n any) (int, error) {
	switch i := n.(type) {
	case int:
		return i, nil
	case int64:
		return int(i), nil
	case int32:
		return int(i), nil
	case int16:
		return int(i), nil
	case int8:
		return int(i), nil
	case uint:
		return int(i), nil
	case uint64:
		return int(i), nil
	case uint32:
		return int(i), nil
	case uint16:
		return int(i), nil
	case uint8:
		return int(i), nil
	case float64:
		return int(i), nil
	case float32:
		return int(i), nil
	default:
		return 0, fmt.Errorf("invalid integer. n=%v, type=%t", n, n)
	}
}
