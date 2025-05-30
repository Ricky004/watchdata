package factory

// Configurable defines an interface that configs must implement.
type Configurable interface {
	// Check if config is valid.
	Validate() error
}

// Config defines something that expose its uique name.
type Factory interface {
	Idx
	NewConfig() Configurable
}

// CreatorFunc is a function that returns a new Configurable.
type CreatorFunc func() Configurable

// Factory implements NamedFactory and can create new configs.
type factory struct {
	id      Id
	creator CreatorFunc
}

// Id returns the unique identifier of the factory.
func (f *factory) Id() Id {
	return f.id
}

// NewConfig creates a new instance of the configuration.
func (f *factory) NewConfig() Configurable {
	return f.creator()
}

// NewFactory returns a new Factory instance.
func NewFactory(id Id, cf CreatorFunc) Factory {
	return &factory{
		id:      id,
		creator: cf,
	}
}
