package stores

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func initIndices(db *gorm.DB) {
	db.Model(&models.Message{}).
		AddIndex("idx_messages_sender", "sender").
		AddIndex("idx_messages_chat", "chat")
	db.Model(&models.Chat{}).
		AddIndex("idx_chat_community", "community")
	db.Model(&models.OrgJoinRequest{}).
		AddIndex("idx_org_join_req_initiator", "initiator").
		AddIndex("idx_org_join_req_acceptor", "acceptor").
		AddIndex("idx_org_join_req_organization", "organization").
		AddIndex("idx_org_join_req_group", "group")
	db.Model(&models.GroupJoinRequest{}).
		AddIndex("idx_group_join_req_initiator", "initiator").
		AddIndex("idx_group_join_req_acceptor", "acceptor").
		AddIndex("idx_group_join_req_group", "group")
	db.Model(&models.Organization{}).
		AddIndex("idx_organization_community", "community").
		AddIndex("idx_organization_director", "director")
	db.Model(&models.MembershipOfCommunity{}).
		AddIndex("idx_membership_of_community_community", "community").
		AddIndex("idx_membership_of_community_user", "user")
	db.Model(&models.GroupCreationRequest{}).
		AddIndex("idx_group_creation_req_initiator", "initiator").
		AddIndex("idx_group_creation_req_acceptor", "acceptor").
		AddIndex("idx_group_creation_req_hm", "hm").
		AddIndex("idx_group_creation_req_organization", "organization")
	db.Model(&models.Group{}).
		AddIndex("idx_group_admin", "admin").
		AddIndex("idx_group_organization", "organization").
		AddIndex("idx_group_community", "community")
	db.Model(&models.RatingEvent{}).
		AddIndex("idx_rating_event_organization", "organization")
	db.Model(&models.RatingEventAverageEstimate{}).
		AddIndex("idx_rating_event_avg_estimate_user", "user").
		AddIndex("idx_rating_event_avg_estimate_group", "group").
		AddIndex("idx_rating_event_avg_estimate_event", "event").
		AddIndex("idx_rating_event_avg_estimate_organization", "organization")
	db.Model(&models.Estimate{}).
		AddIndex("idx_estimate_group", "group").
		AddIndex("idx_estimate_event", "event")
	db.Model(&models.UserNotificationSeen{}).
		AddIndex("idx_user_notification_seen_user", "user").
		AddIndex("idx_user_notification_seen_notification", "notification")
	db.Model(&models.AdminOrgWallSeen{}).
		AddIndex("idx_admin_org_wall_seen_user", "user").
		AddIndex("idx_admin_org_wall_seen_org_wall", "org_wall")
	db.Model(&models.Timeline{}).
		AddIndex("idx_timeline_group", "group").
		AddIndex("idx_timeline_notification", "notification")
	db.Model(&models.OrgExMember{}).
		AddIndex("idx_org_ex_member_organization", "organization").
		AddIndex("idx_org_ex_member_community", "community")
}
