package registry

import (
	"github.com/ABT-Tech-Limited/beosin-go"
	"github.com/ABT-Tech-Limited/kytunified/kyt"
	beosinprovider "github.com/ABT-Tech-Limited/kytunified/provider/beosin"
)

// RegisterBeosin registers a Beosin provider with the given client.
func RegisterBeosin(client beosin.Client, opts ...beosinprovider.Option) error {
	provider := beosinprovider.New(client, opts...)
	return Register("beosin", func(_ map[string]interface{}) (kyt.Provider, error) {
		return provider, nil
	})
}

// MustRegisterBeosin registers a Beosin provider and panics on error.
func MustRegisterBeosin(client beosin.Client, opts ...beosinprovider.Option) {
	if err := RegisterBeosin(client, opts...); err != nil {
		panic(err)
	}
}

// GetBeosin returns the registered Beosin provider.
func GetBeosin() (kyt.Provider, error) {
	return Create("beosin", nil)
}
