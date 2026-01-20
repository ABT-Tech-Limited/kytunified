package registry

import (
	"fmt"
	"sync"

	"github.com/ABT-Tech-Limited/kytunified/kyt"
)

// ProviderFactory is a function that creates a Provider from configuration.
type ProviderFactory func(config map[string]interface{}) (kyt.Provider, error)

// Registry manages KYT provider registration and instantiation.
// It provides a centralized way to register and create providers,
// making it easy to switch between different KYT services.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]ProviderFactory
}

// Global registry instance
var globalRegistry = NewRegistry()

// NewRegistry creates a new provider registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]ProviderFactory),
	}
}

// Register registers a provider factory with the given name.
// Returns an error if a provider with that name is already registered.
func (r *Registry) Register(name string, factory ProviderFactory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[name]; exists {
		return fmt.Errorf("provider %q already registered", name)
	}

	r.factories[name] = factory
	return nil
}

// MustRegister registers a provider factory and panics on error.
// This is useful for init() functions.
func (r *Registry) MustRegister(name string, factory ProviderFactory) {
	if err := r.Register(name, factory); err != nil {
		panic(err)
	}
}

// Unregister removes a provider factory from the registry.
// Returns an error if the provider is not found.
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[name]; !exists {
		return fmt.Errorf("%w: %s", kyt.ErrProviderNotFound, name)
	}

	delete(r.factories, name)
	return nil
}

// Create creates a provider instance by name with the given configuration.
// Returns an error if the provider is not found or creation fails.
func (r *Registry) Create(name string, config map[string]interface{}) (kyt.Provider, error) {
	r.mu.RLock()
	factory, exists := r.factories[name]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %s", kyt.ErrProviderNotFound, name)
	}

	return factory(config)
}

// Has returns true if a provider with the given name is registered.
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.factories[name]
	return exists
}

// List returns all registered provider names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

// Clear removes all registered providers.
// This is mainly useful for testing.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories = make(map[string]ProviderFactory)
}

// Global registry functions

// Register registers a provider in the global registry.
func Register(name string, factory ProviderFactory) error {
	return globalRegistry.Register(name, factory)
}

// MustRegister registers a provider in the global registry and panics on error.
func MustRegister(name string, factory ProviderFactory) {
	globalRegistry.MustRegister(name, factory)
}

// Unregister removes a provider from the global registry.
func Unregister(name string) error {
	return globalRegistry.Unregister(name)
}

// Create creates a provider from the global registry.
func Create(name string, config map[string]interface{}) (kyt.Provider, error) {
	return globalRegistry.Create(name, config)
}

// Has returns true if the provider exists in the global registry.
func Has(name string) bool {
	return globalRegistry.Has(name)
}

// List lists all providers in the global registry.
func List() []string {
	return globalRegistry.List()
}

// Clear clears the global registry.
func Clear() {
	globalRegistry.Clear()
}
