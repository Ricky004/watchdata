package factory

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Registry struct {
	services IdxMap[IdxService]
	startch  chan error
	stopch   chan error
}

func NewRegistry(services ...IdxService) (*Registry, error) {
	if len(services) == 0 {
		return nil, fmt.Errorf("cannot build, at least one service required")
	}

	m, err := NewIdxMap(services...)
	if err != nil {
		return nil, err
	}

	return &Registry{
		services: m,
		startch: make(chan error, 1),
		stopch: make(chan error, len(services)),
	}, nil
}

func (r *Registry) Start(ctx context.Context) {
	for _, s := range r.services.GetInOrder() {
		go func(s IdxService) {
			err := s.Start(ctx)
			r.startch <- err
		}(s)
	}
}

func (r *Registry) Wait(ctx context.Context) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Printf("caught context error, exiting : %v", ctx.Err())
	case s := <-interrupt:
		log.Printf("caught interrupt signal, exiting : %v", s)
	case err := <-r.startch:
		log.Printf("caught service error, exiting : %v", err)
		return err
	}

	return nil
}

func (r *Registry) Stop(ctx context.Context) error {
	for _, s := range r.services.GetInOrder() {
		go func(s IdxService) {
			log.Printf("stopping service : %v", s.Id())
			err := s.Stop(ctx)
			r.stopch <- err
		}(s)
	}

	errs := make([]error, len(r.services.GetInOrder()))
	for range r.services.GetInOrder() {
		err := <-r.stopch
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}