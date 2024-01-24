package jsond

import "fmt"

// jsonpath represents json jsonpath.
type jsonpath []property

// property represents path's element.
// - int, for array index
// - string, for object key
type property any

func (p jsonpath) String() string {
	if len(p) == 0 {
		return ""
	}
	joined := "$"
	for _, prop := range p {
		switch t := prop.(type) {
		case int:
			joined = fmt.Sprintf("%s[%d]", joined, t)
		case string:
			joined = fmt.Sprintf("%s['%s']", joined, t)
		default:
			panic(fmt.Sprintf("invalid property. prop=%v, type=%t", prop, prop))
		}
	}
	return joined
}

func (p jsonpath) append(prop property) jsonpath {
	switch t := prop.(type) {
	case int:
		return append(p, t)
	case string:
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

		index, err := getInt(t)
		if err != nil {
			return nil, fmt.Errorf("invalid index : %v", index)
		}
		return index, nil

	case string:
		return t, nil

	default:
		return nil, fmt.Errorf("invalid property : %v", v)
	}
}
