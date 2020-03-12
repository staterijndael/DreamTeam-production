//go:generate zenrpc
package wall

import "github.com/jinzhu/gorm"

type Service struct {
	db *gorm.DB
} //zenrpc

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}
