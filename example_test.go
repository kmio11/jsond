package jsond_test

import (
	"fmt"

	"github.com/kmio11/jsond"
)

func ExampleNode_Get_methodChain() {
	src := []byte(`
	{
		"total_count": 2,
		"artifacts": [
		  {
			"id": 11,
			"name": "Rails"
		  },
		  {
			"id": 13,
			"name": "Test output"
		  }
		]
	  }
	`)

	var name string
	_ = jsond.Parse(src).
		Get("artifacts").
		Get(1).
		Get("name").
		Unmarshal(&name)

	fmt.Println(name)

	// Output:
	// Test output
}

func ExampleNode_Get() {
	src := []byte(`
	{
		"total_count": 2,
		"artifacts": [
		  {
			"id": 11,
			"name": "Rails"
		  },
		  {
			"id": 13,
			"name": "Test output"
		  }
		]
	  }
	`)

	var name string
	_ = jsond.Parse(src).
		Get("artifacts", 1, "name").
		Unmarshal(&name)

	fmt.Println(name)

	// Output:
	// Test output
}

func ExampleNode_Set() {
	src := []byte(`
	{
		"total_count": 2,
		"artifacts": [
		  {
			"id": 11,
			"name": "Rails"
		  },
		  {
			"id": 13,
			"name": "Test output"
		  }
		]
	  }
	`)

	b, _ := jsond.Parse(src).
		Set(
			"Golang",               // new value
			"artifacts", 0, "name", // path to set new value
		).
		Marshal()

	fmt.Println(string(b))

	// Output:
	// {"artifacts":[{"id":11,"name":"Golang"},{"id":13,"name":"Test output"}],"total_count":2}
}

func ExampleNode_AsArray() {
	src := []byte(`
	{
		"total_count": 2,
		"artifacts": [
		  {
			"id": 11,
			"name": "Rails"
		  },
		  {
			"id": 13,
			"name": "Test output"
		  }
		]
	  }
	`)

	artifacts, _ := jsond.Parse(src).
		Get("artifacts").
		AsArray()

	for i, node := range artifacts {
		var name string
		_ = node.
			Get("name").
			Unmarshal(&name)

		fmt.Printf("%d : %s\n", i, name)
	}

	// Output:
	// 0 : Rails
	// 1 : Test output
}

func ExampleNode_AsObject() {
	src := []byte(`
	{
		"total_count": 2,
		"artifacts": [
		  {
			"id": 11,
			"name": "Rails"
		  },
		  {
			"id": 13,
			"name": "Test output"
		  }
		]
	  }
	`)

	artifact, _ := jsond.Parse(src).
		Get("artifacts").
		Get(0).
		AsObject()

	for key, node := range artifact {
		var value any
		_ = node.
			Unmarshal(&value)

		fmt.Printf("%s : %v\n", key, value)
	}

	// Unordered Output:
	// id : 11
	// name : Rails
}

func ExampleNode_Unmarshal() {
	src := []byte(`
	{
		"total_count": 2,
		"artifacts": [
		  {
			"id": 11,
			"name": "Rails"
		  },
		  {
			"id": 13,
			"name": "Test output"
		  }
		]
	  }
	`)

	artifact := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{}

	_ = jsond.Parse(src).
		Get("artifacts", 1).
		Unmarshal(&artifact)

	fmt.Printf("id  : %d\n", artifact.ID)
	fmt.Printf("name: %s", artifact.Name)

	// Output:
	// id  : 13
	// name: Test output
}

func ExampleNode_Unmarshal_error() {
	src := []byte(`
	{
		"key1": {
			"key2" : null
		}
	}
	`)

	var v string
	err := jsond.Parse(src).
		Get("key1").
		Get("key2"). // null
		Get("xxx").  // cannot read properties of null
		Get("yyy").
		Unmarshal(&v)

	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// cannot read properties of null (reading 'xxx') at $['key1']['key2']['xxx']
}

func ExampleUnmarshalNode() {
	src := []byte(`
	{
		"total_count": 2,
		"artifacts": [
		  {
			"id": 11,
			"name": "Rails"
		  },
		  {
			"id": 13,
			"name": "Test output"
		  }
		]
	  }
	`)

	artifact, _ := jsond.UnmarshalNode[struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}](
		jsond.Parse(src).Get("artifacts", 1),
	)

	fmt.Printf("id  : %d\n", artifact.ID)
	fmt.Printf("name: %s", artifact.Name)

	// Output:
	// id  : 13
	// name: Test output
}

func ExampleUndefined() {
	src := []byte(`
	{
		"key1" : "value1"
	}
	`)

	var v any
	err := jsond.Parse(src).
		Get("invalid_key").
		Unmarshal(&v)

	if jsond.IsUndefined(err) {
		fmt.Println(err)
	}

	// Output:
	// undefined
}

func ExampleNode_Marshal() {
	src := []byte(`
	{
		"total_count": 2,
		"artifacts": [
		  {
			"id": 11,
			"name": "Rails"
		  },
		  {
			"id": 13,
			"name": "Test output"
		  }
		]
	  }
	`)

	artifact, _ := jsond.Parse(src).
		Get("artifacts").
		Get(1).
		Marshal()

	fmt.Println(string(artifact))

	// Output:
	// {"id":13,"name":"Test output"}
}
