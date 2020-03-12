package fns

import (
	"dt/logwrap"
	"dt/models"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type FNSManager struct {
	db     *gorm.DB
	cache  sync.Map
	curKey *FNSKey
}

func NewManager(db *gorm.DB) (*FNSManager, error) {
	var key models.FNSKey
	if err := db.Where(`current_timestamp < expire_date`).First(&key).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &FNSManager{
				db:    db,
				cache: sync.Map{},
				curKey: &FNSKey{
					RWMutex: sync.RWMutex{},
					key: &models.FNSKey{
						Key:        "",
						ExpireDate: time.Time{},
					},
				},
			}, nil
		}

		return nil, err
	}

	return &FNSManager{
		db:    db,
		cache: sync.Map{},
		curKey: &FNSKey{
			RWMutex: sync.RWMutex{},
			key:     &key,
		},
	}, nil
}

func (self *FNSManager) Search(text string) (*OrganizationInfo, error) {
	self.curKey.RLock()
	res, err := http.Get(fmt.Sprintf("https://api-fns.ru/api/search?key=%s&q=%s", self.curKey.key, text))
	self.curKey.RUnlock()
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if len(b) < 150 {
		stringMessage := string(b)
		if strings.Contains(stringMessage, "Ошибка") {
			logwrap.Info("FNS Error: %s", stringMessage)
			return nil, NewError(stringMessage)
		}
	}

	var msg OrganizationInfo

	err = json.Unmarshal(b, &msg)
	if err != nil {
		return nil, err
	}

	for _, org := range msg.Items {
		if org.UL != nil {
			self.cache.Store(*org.UL.INN, org)
		} else {
			self.cache.Store(*org.IP.INN, org)
		}
	}

	return &msg, nil
}

func (self *FNSManager) GetItem(inn string) *Item {
	o, ok := self.cache.Load(inn)
	if !ok {
		return nil
	}

	res, _ := o.(*Item)
	return res
}

func (self *FNSManager) GetCurrentKey() *models.FNSKey {
	return self.curKey.key
}

func (self *FNSManager) GetAllKeys() ([]*models.FNSKey, error) {
	var keys []*models.FNSKey
	if err := self.db.Find(&keys).Error; err != nil {
		logwrap.Error(`FNSManager: "` + err.Error() + `"`)
		return nil, err
	}

	return keys, nil
}

func (self *FNSManager) SetKey(key models.FNSKey) error {
	var err error = nil
	err = self.db.Create(&key).Error
	if err == nil {
		self.curKey.Lock()
		defer self.curKey.Unlock()

		self.curKey.key = &key
	}

	return err
}

type FNSKey struct {
	sync.RWMutex
	key *models.FNSKey
}
