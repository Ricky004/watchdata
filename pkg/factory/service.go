package factory

import "context"

type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type IdxService interface {
	Idx
	Service
}

type idxService struct {
	id Id
	Service
}

func (i *idxService) Id() Id {
	return i.id
}

func NewIdxService(id Id, service Service) IdxService {
	return &idxService{
		id:      id,
		Service: service,
	}
}
