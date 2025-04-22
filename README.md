# DC Package

A lightweight dependency container for Go applications.

## Overview

The `dc` package provides a simple and type-safe dependency container.
It does not inject dependencies automatically, but allows you to:
- Create providers for dependencies using factory functions
- Get dependencies on demand using the `Use()` method
- Mock dependencies for testing
- Detect circular dependencies
- Reset dependencies and mocks

## Features

- **Type Safety**: Uses Go generics for type-safe dependency injection
- **Lazy Loading**: Dependencies are created only when first used via `Use()`
- **Circular Dependency Detection**: Panics if circular dependencies are detected
- **Mocking Support**: Easy mocking for testing using the `Mock()` method
- **Reset Capability**: Reset dependencies and mocks using the `Reset()` function

## Design Principles

- **Configuration as Code**: Component composition is defined in code rather than configuration files
- **Debugging Friendly**: Set breakpoints in factory functions - it's just regular Go code
- **Modular Design**: 
  - Factory functions act as module constructors
  - Internal dependencies are local variables
  - External dependencies are providers
  - Create providers for reusable components, use local variables for one-time use
- **No Magic**:
  - No reflection or code generation
  - No metadata storage about connections
  - No automatic connection computation

## API Overview

The entire API consists of just four functions:
- `dc.Provider(factoryFunction)` - Creates a provider
- `Provider.Use()` - Executes the factory function or returns a cached instance
- `Provider.Mock(&MockInstance)` - Substitutes the service with a mock object
- `dc.Reset()` - Clears the container


## Examples

### Basic Usage

This is local package `<project_path>/dc/search.go`

```go
package dc 

import ....

import "github.com/abratko/dc"

// Define  Reset as current package function.
// It is necessary  for  resetting  container from external scope without import core package. 
// This will come in handy in tests.
var Reset = dc.Reset


type Controller interface {
    ...
} 

// Create new provider with using  factory function 
var SearchGrpcController = dc.Provider(func() Controller {
	...
	// call `Use()` to get or to create instance
    var someExternalDependency = ExternalDependency.Use()

    // inject external dependency manually 
    var internalDependency = app.NewInternalDepenedency(someExternalDependency.Use())
	...

    // Create service instance and inject dependency manually.
	return app.NewController(internalDependency)
})
```

This is `<project_path>/main.go`

```go

package main 

import "<project_path>/dc"

func main() {
    ...
    controller = dc.SearchGrpcController.Use()
    ...
}
```

### Mocking in tests

```go

import "<project_path>/dc"

type Mock struct {
	mock.Mock
}


func TestSearchServer_Exec(t *testing.T) {

	t.Run("your test name", func(t *testing.T) {
		...
		mock := &Mock{}
		
		// create mock for external dependency 
		dc.ExternalDependency.Mock(mock)
		
		// here is we have controller with mocker external dependency 
		service := dc.SearchGrpcController.Use()
		
		// don't forget to dump the container after each test.
		dc.Reset()	
	}
}
```

