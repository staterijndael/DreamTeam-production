//go:generate zenrpc
package report

import "dt/config"

type Service struct {
	conf *config.Config
} //zenrpc

func New(conf *config.Config) *Service {
	return &Service{
		conf: conf,
	}
}
