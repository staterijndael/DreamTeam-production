package stores

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Nickname{}, &models.User{}, &models.File{}, &models.Organization{}, &models.Community{},
		&models.MembershipOfCommunity{}, &models.GroupCreationRequest{}, &models.Group{}, &models.GroupJoinRequest{},
		&models.OrgJoinRequest{}, &models.Chat{}, &models.Message{}, &models.FNSKey{}, &models.RatingEvent{},
		&models.RatingEventAverageEstimate{}, &models.Estimate{}, &models.Notification{}, &models.OrgWall{},
		&models.AdminOrgWallSeen{}, &models.UserNotificationSeen{}, &models.Timeline{}, &models.RatingOrgConfig{},
		&models.OrgExMember{},
	)
}

func initModels(db *gorm.DB) {
	db.Model(&models.Message{}).
		AddForeignKey("sender", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("chat", "chats(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.Chat{}).
		AddForeignKey("community", "communities(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.OrgJoinRequest{}).
		AddForeignKey("initiator", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("acceptor", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("organization", "organizations(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("group", "groups(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.GroupJoinRequest{}).
		AddForeignKey("initiator", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("acceptor", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("group", "groups(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.User{}).
		AddForeignKey("nickname", "nicknames(id)", "RESTRICT", "NO ACTION").
		AddForeignKey("avatar", "files(id)", "RESTRICT", "NO ACTION")
	db.Model(&models.Organization{}).
		AddForeignKey("nickname", "nicknames(id)", "RESTRICT", "NO ACTION").
		AddForeignKey("avatar", "files(id)", "RESTRICT", "NO ACTION").
		AddForeignKey("community", "communities(id)", "CASCADE", "NO ACTION")
	db.Model(&models.MembershipOfCommunity{}).
		AddForeignKey("community", "communities(id)", "CASCADE", "NO ACTION").
		AddForeignKey("user", "users(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.GroupCreationRequest{}).
		AddForeignKey("initiator", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("acceptor", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("hm", "groups(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("organization", "organizations(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.Group{}).
		AddForeignKey("creator", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("admin", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("parent", "groups(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("community", "communities(id)", "CASCADE", "NO ACTION").
		AddForeignKey("avatar", "files(id)", "RESTRICT", "NO ACTION").
		AddForeignKey("nickname", "nicknames(id)", "RESTRICT", "NO ACTION").
		AddForeignKey("organization", "organizations(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("chat", "chats(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.RatingEvent{}).
		AddForeignKey("organization", "organizations(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.RatingEventAverageEstimate{}).
		AddForeignKey("user", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("group", "groups(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("event", "rating_events(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("organization", "organizations(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.Estimate{}).
		AddForeignKey("estimator", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("estimated", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("group", "groups(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("event", "rating_events(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("organization", "organizations(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.UserNotificationSeen{}).
		AddForeignKey("notification", "notifications(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("user", "users(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.OrgWall{}).
		AddForeignKey("notification", "notifications(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("organization", "organizations(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.AdminOrgWallSeen{}).
		AddForeignKey("user", "users(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("org_wall", "org_walls(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.Timeline{}).
		AddForeignKey("group", "groups(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("event", "rating_events(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("notification", "notifications(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.RatingOrgConfig{}).
		AddForeignKey("org_id", "organizations(id)", "NO ACTION", "NO ACTION")
	db.Model(&models.OrgExMember{}).
		AddForeignKey("organization", "organizations(id)", "NO ACTION", "NO ACTION").
		AddForeignKey("community", "communities(id)", "NO ACTION", "NO ACTION")

}
