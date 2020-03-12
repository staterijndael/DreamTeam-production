package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

//requires address, not pointer
func IsMemberOfAnyGroupOfOrg(p *bool, oid, uid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		membership := models.MembershipOfCommunity{UserID: uid}
		group := models.Group{OrganizationID: oid}
		_db := db.
			Model(&membership).
			Where(&membership).
			Where(
				"community in (?)",
				db.
					Model(&group).
					Where(&group).
					Select("community").
					SubQuery(),
			).
			First(&membership)
		if gorm.IsRecordNotFoundError(_db.Error) {
			_db.Error = nil
			*p = false
		} else {
			*p = _db.Error == nil
		}

		return _db
	}
}
