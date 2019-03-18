package types

type Service struct {
	Host             string `json:"host"`
	Available        bool   `json:"available"`
	AvailabilityTime int    `json:"availability_time,omitempty"`
}

type ServiceStore interface {
	Find(host string) (*Service, error)
	Update(*Service) error
	MinTime() (*Service, error)
	MaxTime() (*Service, error)
}
