package factory

import "fmt"

type Idx interface {
	Id() Id
}

type IdxMap[T Idx] struct {
	factories        map[Id]T
	factoriesInOrder []T
}

func NewIdxMap[T Idx](factories ...T) (IdxMap[T], error) {
	fmap := make(map[Id]T)
	for _, factory := range factories {
		if _, ok := fmap[factory.Id()]; ok {
			return IdxMap[T]{}, fmt.Errorf("cannot build factory map, duplicate name %q found", factory.Id())
		}

		fmap[factory.Id()] = factory
	}

	return IdxMap[T]{factories: fmap, factoriesInOrder: factories}, nil
}

func MustNewIdxMap[T Idx](factories ...T) IdxMap[T] {
	nm, err := NewIdxMap(factories...)
	if err != nil {
		panic(err)
	}
	return nm
}

func (i *IdxMap[T]) Get(idstr string) (t T, err error) {
	id, err := NewId(idstr)
	if err != nil {
		return
	}

	factory, ok := i.factories[id]
	if !ok {
		err = fmt.Errorf("factory %q not found or not registered", id)
	}

	t = factory
	return
}

func (i *IdxMap[T]) Add(factory T) (err error) {
	id := factory.Id()
	if _, ok := i.factories[id]; ok {
		return fmt.Errorf("factory %q already exists", id)
	}
	
	i.factories[id] = factory
	i.factoriesInOrder = append(i.factoriesInOrder, factory)
	return nil
}

func (i *IdxMap[T]) GetInOrder() []T {
	return i.factoriesInOrder
}