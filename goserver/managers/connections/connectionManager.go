package connections

import (
	"container/list"
	"dt/models"
	"github.com/gorilla/websocket"
	"sync"
)

func NewManager() *UserConnectionsManager {
	uMap := make(map[uint]*ConnList)
	return &UserConnectionsManager{
		Mutex: sync.Mutex{},
		users: uMap,
	}
}

type UserConnectionsManager struct {
	sync.Mutex
	users map[uint]*ConnList
}

func (m *UserConnectionsManager) Add(u *models.User, con *websocket.Conn) {
	m.Lock()
	defer m.Unlock()

	if old, ok := m.users[u.ID]; ok {
		old.Lock()
		old.PushBack(con)
		old.Unlock()
	} else {
		l := list.New()
		l.PushBack(con)

		newList := &ConnList{
			Mutex: sync.Mutex{},
			List:  *l,
		}

		m.users[u.ID] = newList
	}
}

func (m *UserConnectionsManager) Send(uid uint, msg interface{}, exceptConn *websocket.Conn) error {
	cons, ok := m.users[uid]
	if !ok {
		return UnavailableErr
	}

	cons.Lock()
	defer cons.Unlock()

	for element := cons.Front(); element != nil; {
		if con, ok := element.Value.(*websocket.Conn); con != exceptConn {
			if !ok {
				elementForRemove := element
				element = element.Next()
				cons.Remove(elementForRemove)
				continue
			}

			if err := con.WriteJSON(msg); err != nil {
				elementForRemove := element
				element = element.Next()
				cons.Remove(elementForRemove)
				continue
			}
		}

		element = element.Next()
	}

	return nil
}
