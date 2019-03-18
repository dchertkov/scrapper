package types

type Stat struct {
	RequestTotal     int `json:"request_total"`
	RequestPerMinute int `json:"request_per_minute"`
}

type StatStore interface {
	Find(host string) (*Stat, error)
	Add(host string)
}
