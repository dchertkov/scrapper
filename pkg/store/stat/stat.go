package stat

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dchertkov/scrapper/pkg/types"

	"github.com/paulbellamy/ratecounter"
)

var ErrNotFound = errors.New("Record is not found")

const ratecounterInterval = time.Minute

type statStore struct {
	mu   sync.RWMutex
	data map[string]*stat
}

type stat struct {
	Total     int
	PerMinute *ratecounter.RateCounter
}

func (s *statStore) Find(host string) (*types.Stat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stat, ok := s.data[host]
	if !ok {
		return nil, ErrNotFound
	}

	fmt.Println(stat.PerMinute.Rate())

	return &types.Stat{
		RequestTotal:     stat.Total,
		RequestPerMinute: int(stat.PerMinute.Rate()),
	}, nil
}

func (s *statStore) Add(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.data[host]
	if !ok {
		v = &stat{
			PerMinute: ratecounter.NewRateCounter(ratecounterInterval),
		}
		s.data[host] = v
	}

	v.Total++
	v.PerMinute.Incr(1)
}

func NewStore() types.StatStore {
	return &statStore{
		data: make(map[string]*stat),
	}
}
