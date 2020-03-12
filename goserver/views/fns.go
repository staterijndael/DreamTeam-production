package views

import "dt/models"

type FNSKey struct {
	Key        string `json:"key"`
	ExpireDate int64  `json:"expireDate"`
}

func FNSKeyFromModel(key *models.FNSKey) *FNSKey {
	return &FNSKey{
		Key:        key.Key,
		ExpireDate: key.ExpireDate.Unix(),
	}
}
