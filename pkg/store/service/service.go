package service

import (
	"errors"
	"sync"

	"github.com/dchertkov/scrapper/pkg/types"
)

var ErrNotFound = errors.New("Service is not found")

type serviceStore struct {
	mu   sync.RWMutex
	data map[string]*types.Service
	min  *types.Service
	max  *types.Service
}

func (s *serviceStore) Find(host string) (*types.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	service, ok := s.data[host]
	if !ok {
		return nil, ErrNotFound
	}

	return service, nil
}

func (s *serviceStore) updateMinMax(service *types.Service) {
	if !service.Available {
		return
	}

	if s.min == nil || s.min.AvailabilityTime > service.AvailabilityTime {
		s.min = service
	}

	if s.max == nil || s.max.AvailabilityTime < service.AvailabilityTime {
		s.max = service
	}
}

func (s *serviceStore) Update(service *types.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[service.Host] = service
	s.updateMinMax(service)

	return nil
}

func (s *serviceStore) MinTime() (*types.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.min == nil {
		return nil, ErrNotFound
	}

	return s.min, nil
}

func (s *serviceStore) MaxTime() (*types.Service, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.max == nil {
		return nil, ErrNotFound
	}

	return s.max, nil
}

func NewStore() types.ServiceStore {
	return &serviceStore{
		data: make(map[string]*types.Service),
	}
}
