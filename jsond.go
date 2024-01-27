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
func (n *Node) newChild(value jsonvalue, prop property, err error) *Node {
	return &Node{
		parent: n,
		value:  value,
		path:   n.path.append(prop),
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

func (n *Node) getArrayElement(idx arrayIndex) *Node {
	path := n.path.append(idx)

	array, ok := n.value.([]any)
	if !ok || int(idx) > len(array) {
		return n.newChild(nil, idx, newUndefined(path))
	}

	return &Node{
		parent: n,
		value:  array[idx],
		path:   path,
		err:    nil,
	}
}

func (n *Node) getObjectValue(key objectKey) *Node {
	path := n.path.append(key)

	object, ok := n.value.(map[string]any)
	if !ok {
		return n.newChild(nil, key, newUndefined(path))
	}

	v, ok := object[string(key)]
	if !ok {
		return n.newChild(nil, key, newUndefined(path))
	}

	return n.newChild(v, key, nil)
}

// Get retrieves a child node based on the specified property (index or key).
// If no properties are provided, the current node is returned.
// It supports nested property access using variadic parameters.
// If multiple properties are provided, it recursively calls Get on each property.
// If an error occurs during the operation, it returns a new Node with the error.
func (n *Node) Get(props ...any) *Node {
	if len(props) == 0 {
		return n
	}

	if len(props) >= 2 {
		return n.Get(props[0]).Get(props[1:]...)
	}

	prop := props[0]
	validProp, err := getProperty(prop)
	if err != nil {
		panic(fmt.Sprintf("invalid property. prop=%v, type=%t", prop, prop))
	}

	path := n.path.append(validProp)

	if n.IsUndefined() {
		return n.newChild(nil, validProp, newReadUndefinedError(path))
	}

	if n.err != nil {
		return n
	}

	if n.value == nil {
		return n.newChild(nil, validProp, newReadNullError(path))
	}

	switch typedProp := validProp.(type) {
	case arrayIndex:
		return n.getArrayElement(typedProp)

	case objectKey:
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
				n.newChild(v, arrayIndex(i), nil),
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
			nodeMap[k] = n.newChild(v, objectKey(k), nil)
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

	return unmarshal(n.path, data, v)
}

// Marshal marshals the Node's value into JSON format.
func (n *Node) Marshal() ([]byte, error) {
	if n.err != nil {
		return nil, n.err
	}

	return marshal(n.path, n.value)
}

// UnmarshalNode is a helper function to unmarshal a Node's value into a specified type.
func UnmarshalNode[T any](node *Node) (T, error) {
	var v = *new(T)
	err := node.Unmarshal(&v)
	return v, err
}

// Set sets the specified value at the given property path within the JSON structure.
// If no properties are provided, it returns a new Node with the provided value.
// It supports nested property access using variadic parameters.
// If an error occurs during the operation, it returns a new Node with the error.
func (n *Node) Set(value any, props ...any) *Node {
	if n.err != nil && !n.IsUndefined() {
		return n
	}

	if len(props) == 0 {
		return n.replaceValue(value)
	}

	targetNode := n.
		Get(props...). // get a Node which is replaced by the given value
		Set(value)     // call Set with no props

	if targetNode.err != nil {
		if nodeErr, ok := targetNode.err.(*NodeError); ok {
			isErrInTargetNode := len(n.path)+len(props) == len(targetNode.path)

			// codeReadUndefinedError may occured in Get().
			if nodeErr.code == codeReadUndefinedError && isErrInTargetNode {
				targetNode.err = newSetUndefinedError(targetNode.path)
				return targetNode
			}

			// codeReadNullError may occured in Get().
			if nodeErr.code == codeReadNullError && isErrInTargetNode {
				targetNode.err = newSetNullError(targetNode.path)
				return targetNode
			}

		}
		return targetNode
	}

	return targetNode.newParent(len(props))
}

func (n *Node) replaceValue(value any) *Node {
	jvalue, err := getJSONValue(value)
	if err != nil {
		return &Node{
			parent: n.parent,
			value:  n.value,
			path:   n.path,
			err:    newInternalError(n.path, err),
		}
	}

	return &Node{
		parent: n.parent,
		value:  jvalue,
		path:   n.path,
		err:    nil,
	}
}

func (n *Node) setArrayElement(v jsonvalue, idx arrayIndex) *Node {
	array, ok := n.value.([]any)
	if !ok {
		return n.newChild(nil, idx,
			newCreatePopertyError(n.path.append(idx), n.value),
		)
	}

	newLen := len(array)
	if int(idx) > len(array) {
		newLen = int(idx) + 1
	}

	newValue := make([]any, newLen)
	copy(newValue, array)
	newValue[int(idx)] = v

	return &Node{
		parent: n.parent,
		value:  newValue,
		path:   n.path,
		err:    nil,
	}
}

func (n *Node) setObjectValue(v jsonvalue, key objectKey) *Node {

	object, ok := n.value.(map[string]any)
	if !ok {
		return n.newChild(
			nil, key,
			newCreatePopertyError(n.path.append(key), n.value),
		)
	}

	newValue := map[string]any{}
	for k, v := range object {
		newValue[k] = v
	}
	newValue[string(key)] = v

	return &Node{
		parent: n.parent,
		value:  newValue,
		path:   n.path,
		err:    nil,
	}
}

func (n *Node) setValue(v jsonvalue, prop property) *Node {
	path := n.path.append(prop)

	if n.IsUndefined() {
		return n.newChild(nil, prop, newSetUndefinedError(path))
	}

	if n.err != nil {
		return n
	}

	if n.value == nil {
		return n.newChild(nil, prop, newSetNullError(path))
	}

	switch typedProp := prop.(type) {
	case arrayIndex:
		return n.setArrayElement(v, typedProp)
	case objectKey:
		return n.setObjectValue(v, typedProp)
	default:
		panic(fmt.Sprintf("invalid property. prop=%v, type=%t", prop, prop))
	}
}

// newParent creates a new parent node by traversing up the hierarchy by the specified depth.
// It attempts to set the value at the current node in the parent node, preserving the path up to the last property.
// If the depth is 0 or an error is present in the current node, it returns the current node.
func (n *Node) newParent(depth int) *Node {
	if depth <= 0 || n.err != nil {
		return n
	}

	return n.parent.
		setValue(n.value, n.path[len(n.path)-1]).
		newParent(depth - 1)
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
