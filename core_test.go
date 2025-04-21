package dc

import (
	"testing"
)

type testService struct {
	value string
}

// Helper function to reset the global state
func resetGlobalState() {
	providers = make([]providerAbstract, 0)
}

func TestProvider_Use(t *testing.T) {
	resetGlobalState()

	t.Run("should create instance using factory function", func(t *testing.T) {
		descriptor := Provider(func() testService {
			return testService{value: "test"}
		})
		instance := descriptor.Use()

		if instance.value != "test" {
			t.Errorf("Use() = %v, want %v", instance.value, "test")
		}
	})
	t.Run("should create function using factory function", func(t *testing.T) {
		descriptor := Provider(func() func(string) string {
			return func(s string) string {
				return s
			}
		})
		instance := descriptor.Use()

		result := instance("test")
		if result != "test" {
			t.Errorf("Use() = %v, want %v", result, "test")
		}
	})

	t.Run("should return mock function when set", func(t *testing.T) {
		descriptor := Provider(func() func(string) string {
			return func(s string) string {
				return s
			}
		})
		descriptor.Mock(func(s string) string {
			return "mocked"
		})

		instance := descriptor.Use()

		result := instance("test")
		if result != "mocked" {
			t.Errorf("Use() = %v, want %v", result, "mocked")
		}
	})

	t.Run("should return mock instance when set", func(t *testing.T) {
		descriptor := Provider(func() testService {
			return testService{value: "original"}
		})
		descriptor.Mock(testService{value: "mocked"})
		instance := descriptor.Use()

		if instance.value != "mocked" {
			t.Errorf("Use() = %v, want %v", instance.value, "mocked")
		}
	})

	t.Run("should create instance only once", func(t *testing.T) {
		counter := 0
		descriptor := Provider(func() testService {
			counter++
			return testService{value: "test"}
		})

		// Call Use() multiple times
		descriptor.Use()
		descriptor.Use()
		descriptor.Use()

		if counter != 1 {
			t.Errorf("Use() created multiple instances, counter = %d", counter)
		}
	})

	t.Run("should panic on circular dependency", func(t *testing.T) {
		resetGlobalState()

		var circularDescriptor *provider[testService]
		circularDescriptor = Provider(func() testService {
			// This will cause a circular dependency
			circularDescriptor.Use()
			return testService{value: "test"}
		})

		defer func() {
			if r := recover(); r == nil {
				t.Error("Use() did not panic on circular dependency")
			}
		}()

		circularDescriptor.Use()
	})
}

func TestProvider_Registration(t *testing.T) {
	resetGlobalState()

	t.Run("should create and register descriptor", func(t *testing.T) {
		descriptor := Provider(func() testService {
			return testService{value: "test"}
		})
		instance := descriptor.Use()

		if instance.value != "test" {
			t.Errorf("Provider() = %v, want %v", instance.value, "test")
		}

		if len(providers) != 1 {
			t.Errorf("Provider() did not register descriptor, got %d descriptors", len(providers))
		}
	})
}

func TestReset(t *testing.T) {
	resetGlobalState()

	t.Run("should call factory function after reset", func(t *testing.T) {
		counter := 0
		descriptor := Provider(func() testService {
			counter++
			return testService{value: "test"}
		})

		// First use
		instance1 := descriptor.Use()
		if instance1.value != "test" {
			t.Errorf("First Use() = %v, want %v", instance1.value, "test")
		}
		if counter != 1 {
			t.Errorf("Factory function called %d times, want 1", counter)
		}

		// Reset
		Reset()

		// Second use should call factory again
		instance2 := descriptor.Use()
		if instance2.value != "test" {
			t.Errorf("Second Use() = %v, want %v", instance2.value, "test")
		}
		if counter != 2 {
			t.Errorf("Factory function called %d times after reset, want 2", counter)
		}
	})
}
