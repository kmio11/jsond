package main

import (
	"fmt"

	"github.com/kmio11/jsond"
)

func main() {
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

	//
	// Parsing JSON Data
	//
	rootNode := jsond.Parse(data)

	//
	// Retrieving Values
	//
	firstCommentNode := rootNode.Get("post", "comments", 0)

	b, _ := firstCommentNode.Marshal()
	fmt.Printf("firstCommentNode's value %s\n", string(b))
	// Output: firstCommentNode's value: {"author":"Alice","content":"Nice post!","id":1}

	//
	// Setting Values
	//
	newFirstCommentNode := firstCommentNode.Set("Very nice post!", "content")

	b, _ = newFirstCommentNode.Marshal()
	fmt.Printf("newFirstCommentNode's value: %s\n", string(b))
	// Output: newFirstCommentNode's value: {"author":"Alice","content":"Very nice post!","id":1}

	//
	// Unmarshalling JSON data
	//

	// Getting author, and unmarshal to string.
	var author string
	_ = firstCommentNode.
		Get("author").
		Unmarshal(&author)
	fmt.Printf("author: %s\n", author)
	// Output: author: Alice

	// Unmarshal a comment node to the struct Comment
	type Comment struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	firstComment := Comment{}
	_ = firstCommentNode.Unmarshal(&firstComment)
	fmt.Printf("first comment: %+v\n", firstComment)
	// Output: first comment: {ID:1 Content:Nice post! Author:Alice}

	//
	// Marshalling JSON data
	//
	jsonData, _ := firstCommentNode.Marshal()
	fmt.Printf("jsonData: %s\n", string(jsonData))
	// Output: jsonData: {"author":"Alice","content":"Nice post!","id":1}

	//
	// Error Handling and Undefined Values
	//
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

	errorNode := undefinedNode.Get("key")
	if err := errorNode.Error(); err != nil {
		fmt.Println(err) // cannot read properties of undefined (reading 'key') at $['undefined_key']['key']
	}
}
