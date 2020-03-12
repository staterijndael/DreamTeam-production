package models

type OrgExMember struct {
	OrganizationID uint         `gorm:"column:organization"`
	CommunityID    uint         `gorm:"column:community"`
	Organization   Organization `gorm:"foreignkey:organization"`
	Community      Community    `gorm:"foreignkey:community"`
}
