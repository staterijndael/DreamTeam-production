package sms

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"math/rand"
	"time"
)

var (
	key      = "kg1DOpEC6mG7Y8nwVP8x9WxnN8rfPejA" //TODO: byteHand api key
	template = "Authorization code: %d"
)

type Manager struct {
	random *rand.Rand
	queue  *queue.Queue
	codes  map[uint]int
}

type generatedCode struct {
	uid    uint
	expire time.Time
}

func (m *Manager) GenerateCode(uid uint) int {
	code := m.random.Intn(1000000)
	for ; code < 100000; code = m.random.Intn(1000000) {
	}
	generated := &generatedCode{
		uid:    uid,
		expire: time.Now().Add(time.Minute * 5),
	}

	m.queue.Enqueue(generated)
	m.codes[uid] = code
	return code
}

func (m *Manager) cleaner() {
	for {
		if m.queue.Len() == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		codeInfo := m.queue.Dequeue().(*generatedCode)
		if _, ok := m.codes[codeInfo.uid]; !ok {
			continue
		}

		time.Sleep(time.Until(codeInfo.expire))
		if _, ok := m.codes[codeInfo.uid]; !ok {
			continue
		}

		delete(m.codes, codeInfo.uid)
	}
}

func (m *Manager) Send(uid uint, phone string) (int, error) {
	code := m.GenerateCode(uid)
	_, err := sendSms(key, phone, fmt.Sprintf(template, code))
	return code, err
}

func (m *Manager) Get(uid uint) (code int, ok bool) {
	code, ok = m.codes[uid]
	return
}

func (m *Manager) Delete(uid uint) {
	delete(m.codes, uid)
}

func NewManager() *Manager {
	m := &Manager{
		queue:  queue.New(),
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		codes:  make(map[uint]int),
	}

	go m.cleaner()
	return m
}
