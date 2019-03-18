package checker

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/dchertkov/scrapper/pkg/checker/dialer"
	"github.com/dchertkov/scrapper/pkg/config"
	"github.com/dchertkov/scrapper/pkg/types"
)

type Checker struct {
	ctx context.Context

	workers  int
	interval time.Duration
	dialer   *dialer.Dialer

	service types.ServiceStore

	hosts []string
}

func NewChecker(
	conf config.Checker,
	service types.ServiceStore,
	hosts []string,
) *Checker {
	return &Checker{
		workers:  conf.Workers,
		interval: conf.Interval,
		dialer:   dialer.NewDialer(conf.Timeout),
		service:  service,
		hosts:    hosts,
	}
}

func (ch *Checker) Start(ctx context.Context) error {
	err := ch.do()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(ch.interval)
	for {
		select {
		case <-ticker.C:
			if err = ch.do(); err != nil {
				return err
			}
		case <-ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

func (ch *Checker) do() error {
	logrus.Info("start checking")

	hosts := make(chan string)

	g := errgroup.Group{}
	for i := 0; i < ch.workers; i++ {
		g.Go(func() error {
			return ch.worker(hosts)
		})
	}

	go func() {
		for i := range ch.hosts {
			hosts <- ch.hosts[i]
		}
		close(hosts)
	}()

	return g.Wait()
}

func (ch *Checker) worker(hosts <-chan string) error {
	for host := range hosts {
		if err := ch.check(host); err != nil {
			return err
		}
	}
	return nil
}

func (ch *Checker) check(host string) (err error) {
	s := &types.Service{Host: host}

	logger := logrus.WithFields(logrus.Fields{
		"host": host,
	})

	s.AvailabilityTime, err = ch.dialer.Dial(host)
	if err != nil {
		logger.WithError(err).Error("host check")
	} else {
		s.Available = true
	}

	return ch.service.Update(s)
}
