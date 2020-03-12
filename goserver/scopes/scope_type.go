package scopes

import "github.com/jinzhu/gorm"

type Scope func(*gorm.DB) *gorm.DB
