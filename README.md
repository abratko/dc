# DC Package

A lightweight dependency injection container for Go applications.

## Overview

The `dc` package provides a simple and type-safe dependency injection system. It allows you to:
- Register dependencies with their creation functions
- Mock dependencies for testing
- Detect circular dependencies
- Reset dependencies and mocks

## Features

- **Type Safety**: Uses Go generics for type-safe dependency injection
- **Lazy Loading**: Dependencies are created only when first used
- **Circular Dependency Detection**: Panics if circular dependencies are detected
- **Mocking Support**: Easy mocking for testing
- **Reset Capability**: Reset dependencies and mocks as needed


## Usage

### Registering Dependencies

```go
// Create a provider for a dependency
userRepoProvider := dc.Provider(func() UserRepository {
    return NewUserRepository()
})

// Use the dependency
userRepo := userRepoProvider.Use()
```

### Mocking Dependencies

```go

// In tests
func TestUserService(t *testing.T) {
    // Create a mock
    mockRepo := &MockUserRepository{}
    // Set the mock
    userRepo.Mock(mockRepo)

    
    service := userRepo.Use
    // Test with mock...
    
    userRepoProvider.ResetMock()
}


// Use the mock
userRepo := userRepoProvider.Use() // Returns mockRepo

// Reset the mock when done
userRepoProvider.ResetMock()
```

### Resetting Dependencies

```go
// Reset a single provider
userRepoProvider.Reset()

// Reset all providers
dc.Reset()
```

### Resetting Mocks

```go
// Reset all mocks
dc.ResetMocks()
```


## Example

```go
// Define your dependencies
type UserRepository interface {
    GetUser(id string) (*User, error)
}

// Create providers
userRepoProvider := dc.Provider(func() UserRepository {
    return NewUserRepository()
})

// Use in your code
func GetUserService() *UserService {
    repo := userRepoProvider.Use()
    return NewUserService(repo)
}

```

## Best Practices

1. Create providers at the package level or in your dependency injection setup
2. Use the `Use()` method to access dependencies
3. Reset mocks after tests
4. Use the global `Reset()` and `ResetMocks()` functions to clean up all providers
5. Keep provider creation functions simple and focused

## Notes

- The package uses a global registry of providers
- Circular dependencies are detected and will cause a panic
- Dependencies are created lazily (only when first used)
- Mocking is thread-safe but not concurrent-safe
