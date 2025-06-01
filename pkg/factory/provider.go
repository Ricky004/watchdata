package factory

import (
	"context"
)

// Provider is a constraint for any type
type Provider = any

// NewProviderFunc defines a function that creates a new provider
type NewProviderFunc[P Provider, C Configurable] func(ctx context.Context, config C) (P, error)

// ProviderFactory defines an interface for a typed provider factory
type ProviderFactory[P Provider, C Configurable] interface {
	Idx
	New(context.Context, C) (P, error)
}

// providerFactory is the default implementation
type providerFactory[P Provider, C Configurable] struct {
	id              Id
	newProviderFunc NewProviderFunc[P, C]
}

func (factory *providerFactory[P, C]) Id() Id {
	return factory.id
}

func (factory *providerFactory[P, C]) New(ctx context.Context, config C) (p P, err error) {
	provider, err := factory.newProviderFunc(ctx, config)
	if err != nil {
		return
	}

	p = provider
	return
}

// NewProviderFactory returns a new generic provider factory
func NewProviderFactory[P Provider, C Configurable](id Id, fn NewProviderFunc[P, C]) ProviderFactory[P, C] {
	return &providerFactory[P, C]{
		id:              id,
		newProviderFunc: fn,
	}
}

// NewProviderFromNamedMap creates a provider instance from a factory map using the given key
func NewProviderFromIdxMap[P Provider, C Configurable](
	ctx context.Context,
	config C,
	factories IdxMap[ProviderFactory[P, C]],
	key string,
) (p P, err error) {
	pFactory, err := factories.Get(key)
	if err != nil {
		return
	}

	provider, err := pFactory.New(ctx, config)
	if err != nil {
		return
	}

	p = provider 
	return 
}
