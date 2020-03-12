package plebs

import (
	"dt/logwrap"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/servers/common"
	"dt/views"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func (rpc *Server) HTTPHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query(requestContext.UserQueryKey))
	var user models.User
	err := rpc.sqlStore.First(&user, id).Error
	if err != nil && err != models.EmptyNicknameModelErr {
		if gorm.IsRecordNotFoundError(err) {
			views.NotFound(c, fmt.Sprintf("user, specified as <%s>, not found", requestContext.UserQueryKey))
			return
		}

		logwrap.Error("[plebs.HTTPHandler]: %s", err.Error())
		views.InternalError(c)
		return
	}

	c.Request = c.Request.WithContext(common.NewRequestContext(c.Request.Context(), &user, nil, c.Request))
	rpc.ServeHTTP(c.Writer, c.Request)
}
