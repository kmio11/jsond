# jsond

jsond is a dynamic JSON handling library for Go. It provides a simple and flexible way to parse and manipulate JSON data.

## Features

- **Dynamic JSON parsing**: Parse JSON data into a dynamic structure that can be easily navigated and manipulated.
- **Flexible Data Access**: Access data in the JSON structure using a simple Get method. You can retrieve data by index for arrays or by key for objects.
- **Data Modification**: Modify data in the JSON structure using the Set method. You can set data by index for arrays or by key for objects.
- **Type Conversion**: Easily convert a value into a specified type using Unmarshal / Marshal functions.

## Installation

```bash
go get github.com/kmio11/jsond
```

## Usage

The core concepts of `jsond` is a `Node`, which represents a node in the JSON data structure.
Each `Node` contains a value, and an error. This structure allows for efficient and easy traversal and manipulation of JSON data.

### Parsing JSON Data

To parse JSON data, use the `Parse` function. This function returns a `Node` representing the parsed structure.

The following code shows how to parse a JSON string into a `Node`.

```go
data := []byte(`
{
	"post": {
		"id": 1,
		"title": "Hello World!",
		"content": "This is my first post.",
		"author": "John Doe",
		"comments": [
			{
				"id": 1,
				"content": "Nice post!",
				"author": "Alice"
			},
			{
				"id": 2,
				"content": "Thanks for sharing.",
				"author": "Bob"
			}
		]
	},
	"status": "published",
	"views": 1500
}
`)

rootNode := jsond.Parse(data)
```

### Retrieving Values

To retrieve a value from the JSON data, use the `Get` method on a `Node`.
This method takes a series of properties (either indexes for an array or keys for an object) and returns the new `Node` at the specified path.

The following code shows how to retrieve the first comment from a blog post.

```go
firstCommentNode := rootNode.Get("post", "comments", 0)
```

This `firstCommentNode` is a `Node` representing the following JSON data:

```json
{
  "id": 1,
  "content": "Nice post!",
  "author": "Alice"
}
```

### Setting Values

To set a value in the JSON data, use the `Set` method on a `Node`.
This method takes a value and a series of properties (either indexes for an array or keys for an object) and returns a new `Node` with the value set at the specified path.

The following code shows how to change the content of a first comment.

```go
newFirstCommentNode := firstCommentNode.Set("Very nice post!", "content")
```

This `newFirstCommentNode` is a `Node` representing a following JSON data:

```json
{
  "id": 1,
  "content": "Very nice post!",
  "author": "Alice"
}
```

### Unmarshalling and Marshalling JSON data

The `Unmarshal` method allows you to unmarshal a `Node`'s value into a specified variable.
This is particularly useful when you want to convert a `Node`'s value into a custom data type.

The following code shows how to unmarshal the author of a first comment.

```go
var author string
err = firstCommentNode.
	Get("author").
	Unmarshal(&author)

fmt.Printf("author: %s\n", author)
// author: Alice
```

You can also unmarshal to a struct data type.  
The following code shows how to unmarshal the first comment.

```go
type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

firstComment := Comment{}
err = firstCommentNode.Unmarshal(&firstComment)

fmt.Printf("first comment: %+v\n", firstComment)
// first comment: {ID:1 Content:Nice post! Author:Alice}
```

On the other hand, the `Marshal` method allows you to marshal a `Node`'s value into JSON format.
This is useful when you want to convert a `Node`'s value back into JSON data.

The following code shows how to marshal a `firstCommentNode` back into a JSON string.

```go
jsonData, err := firstCommentNode.Marshal()
fmt.Printf("jsonData: %s\n", string(jsonData))
// jsonData: {"author":"Alice","content":"Nice post!","id":1}
```

### Error Handling and Undefined Values

If an error occurs during any operation, the resulting `Node`'s `Error` method will return an error.
In addition to regular errors, it also provides a way to represent undefined values through the `Undefined` type.
If a `Node` represents an undefined value, the `IsUndefined` method can be used to check this.
This is particularly useful when dealing with optional fields in JSON data.

The following code shows how to handle errors and undefined values.

```go
undefinedNode := rootNode.Get("undefined_key")
if err := undefinedNode.Error(); err != nil {
	if undefinedNode.IsUndefined() {
		// handle undefined value
		fmt.Println(err) // undefined
	} else {
		// handle other errors
		fmt.Println(err)
	}
}
```

And if you try to read an undefined node, for example, the error is set to node.

```go
errorNode := undefinedNode.Get("key")
if err := errorNode.Error(); err != nil {
	fmt.Println(err) // cannot read properties of undefined (reading 'key') at $['undefined_key']['key']
}
```

This allows you to distinguish between regular errors and undefined values, providing more control over your error handling logic.
