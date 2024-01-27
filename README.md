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
### Getting JSON value
```go
package main

import (
	"fmt"

	"github.com/kmio11/jsond"
)

func main() {
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
```

### Setting JSON value
```go
package main

import (
	"fmt"

	"github.com/kmio11/jsond"
)

func main() {
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
```