//go:generate zenrpc
package fns

import "dt/managers/fns"

type Service struct {
	fnsMgr *fns.FNSManager
} //zenrpc

func New(fns *fns.FNSManager) *Service {
	return &Service{fnsMgr: fns}
}
