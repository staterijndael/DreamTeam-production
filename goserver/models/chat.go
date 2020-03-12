package models

import "github.com/jinzhu/gorm"

type Chat struct {
	gorm.Model
	CommunityID uint      `gorm:"column:community"`
	Community   Community `gorm:"foreignkey:community"`
}

func (*Chat) TableName() string {
	return "chats"
}

type Message struct {
	gorm.Model
	Text	 string `gorm:"column:text_of_message"`
	SenderID uint `gorm:"column:sender"`
	ChatID   uint `gorm:"column:chat"`
	Sender   User `gorm:"foreignkey:sender"`
	Chat     Chat `gorm:"foreignkey:chat;PRELOAD:false"`
}

func (*Message) TableName() string {
	return "messages"
}
