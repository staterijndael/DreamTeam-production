package connections

import (
	"container/list"
	"errors"
	"sync"
)

var (
	UnavailableErr = errors.New("unavailable user")
)

type ConnList struct {
	sync.Mutex
	list.List
}

func (m *ConnList) remove(el *list.Element) {

}
