# jsond
jsond is a dynamic JSON handling library for Go, designed to provide flexible and dynamic access to JSON data structures. 

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
	"log"

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
		Get("artifacts").
		Get(1).
		Get("name").
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