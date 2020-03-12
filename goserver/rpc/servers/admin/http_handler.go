package admin

import (
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
	id, _ := strconv.Atoi(c.Query(requestContext.UserContextKey))
	var user models.User
	err := rpc.sqlStore.First(&user, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			views.NotFound(c, fmt.Sprintf("user, specified as <%s>, not found", requestContext.UserQueryKey))
			return
		}

		views.InternalError(c)
		return
	}

	c.Request = c.Request.WithContext(common.NewRequestContext(c.Request.Context(), &user, nil, c.Request))
	rpc.ServeHTTP(c.Writer, c.Request)
}
