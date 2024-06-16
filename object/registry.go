// registry.go
package object

import (
	"log"
	"sync"
)

var (
	functionRegistry = make(map[string]Extended)
	registryMutex    = sync.RWMutex{}
)

// RegisterFunction registers a function in the global registry
func RegisterFunction(name string, fn Extended) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	log.Printf("Registering function: %s", name)
	functionRegistry[name] = fn
}

// GetFunction retrieves a function from the global registry and casts to Extended
func GetExtendedFunction(name string) (Extended, bool) {
	fn, exists := GetFunction(name)
	return Extended(fn), exists
}

func GetFunction(name string) (Extended, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()
	fn, exists := functionRegistry[name]
	return fn, exists
}
