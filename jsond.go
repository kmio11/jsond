package jsond

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Node represents a node in the JSON data structure.
type Node struct {
	parent *Node
	value  jsonvalue
	path   jsonpath
	err    error
}

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

// Parse parses the given JSON data and returns a Node representing the parsed structure.
func Parse(data []byte) *Node {
	value := *new(jsonvalue)
	path := []property{}

	err := unmarshal(path, data, &value)
	return &Node{
		parent: nil,
		value:  value,
		path:   path,
		err:    err,
	}
}

// newChild creates a new child node with the given arguments.
func (n *Node) newChild(value jsonvalue, path jsonpath, err error) *Node {
	return &Node{
		parent: n,
		value:  value,
		path:   path,
		err:    err,
	}
}

// IsUndefined checks if the Node represents an undefined value.
func (n *Node) IsUndefined() bool {
	return IsUndefined(n.err)
}

// Error returns Node's error.
func (n *Node) Error() error {
	return n.err
}

func (n *Node) getArrayElement(idx int) *Node {
	path := n.path.append(idx)

	array, ok := n.value.([]any)
	if !ok || idx > len(array) {
		return n.newChild(nil, path, newUndefined(path))
	}

	return &Node{
		parent: n,
		value:  array[idx],
		path:   path,
		err:    nil,
	}
}

func (n *Node) getObjectValue(key string) *Node {
	path := n.path.append(key)

	object, ok := n.value.(map[string]any)
	if !ok {
		return n.newChild(nil, path, newUndefined(path))
	}

	v, ok := object[key]
	if !ok {
		return n.newChild(nil, path, newUndefined(path))
	}

	return n.newChild(v, path, nil)
}

// Get retrieves a child node based on the specified property (index or key).
func (n *Node) Get(prop any) *Node {
	validProp, err := getProperty(prop)
	if err != nil {
		panic(fmt.Sprintf("invalid property. prop=%v, type=%t", prop, prop))
	}

	path := n.path.append(validProp)

	if n.IsUndefined() {
		return n.newChild(nil, path, newReadUndefinedError(path))
	}

	if n.err != nil {
		return n
	}

	if n.value == nil {
		return n.newChild(nil, path, newReadNullError(path))
	}

	switch typedProp := validProp.(type) {
	case int:
		return n.getArrayElement(typedProp)

	case string:
		return n.getObjectValue(typedProp)

	default:
		panic(fmt.Sprintf("invalid property. prop=%v, type=%t", validProp, validProp))
	}
}

// AsArray returns the Node's value as an array of child nodes.
func (n *Node) AsArray() ([]*Node, error) {
	if n.err != nil {
		return nil, n.err
	}
	if n.value == nil {
		return nil, errors.New("node value is nil")
	}
	if a, ok := n.value.([]any); ok {
		nodeArray := []*Node{}
		for i, v := range a {
			nodeArray = append(nodeArray,
				n.newChild(v, n.path.append(i), nil),
			)
		}
		return nodeArray, nil
	}
	return nil, errors.New("node is not an array")
}

// AsObject returns the Node's value as a map of string keys to child nodes.
func (n *Node) AsObject() (map[string]*Node, error) {
	if n.err != nil {
		return nil, n.err
	}
	if n.value == nil {
		return nil, errors.New("node value is nil")
	}
	if m, ok := n.value.(map[string]any); ok {
		nodeMap := map[string]*Node{}
		for k, v := range m {
			nodeMap[k] = n.newChild(v, n.path.append(k), nil)
		}
		return nodeMap, nil
	}
	return nil, errors.New("node is not a object")
}

// Unmarshal unmarshals the Node's value into the specified variable.
func (n *Node) Unmarshal(v any) error {
	if n.err != nil {
		return n.err
	}

	data, err := marshal(n.path, n.value)
	if err != nil {
		return err
	}

	return unmarshal(n.path, data, &v)
}

// Marshal marshals the Node's value into JSON format.
func (n *Node) Marshal() ([]byte, error) {
	if n.err != nil {
		return nil, n.err
	}

	return marshal(n.path, n.value)
}

// Typed is a helper function to unmarshal a Node's value into a specified type.
func Typed[T any](node *Node) (T, error) {
	var v = *new(T)
	err := node.Unmarshal(v)
	return v, err
}

func unmarshal(path jsonpath, data []byte, v jsonvalue) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return &NodeError{
			code: codeUnmarshalError,
			path: path,
			err:  err,
		}
	}
	return nil
}

func marshal(path jsonpath, v jsonvalue) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, &NodeError{
			code: codeMarshalError,
			path: path,
			err:  err,
		}
	}
	return data, nil
}

func getInt(n any) (int, error) {
	switch i := any(n).(type) {
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
