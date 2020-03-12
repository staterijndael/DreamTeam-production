package views

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InternalError(c *gin.Context, additionalInfos ...gin.H) {
	ErrorResponse(c, http.StatusInternalServerError, errorResponseBuilder("internal error", additionalInfos...))
}

func BadRequest(c *gin.Context, err string, additionalInfos ...gin.H) {
	ErrorResponse(c, http.StatusBadRequest, errorResponseBuilder(err, additionalInfos...))
}

func NotFound(c *gin.Context, err string, additionalInfos ...gin.H) {
	ErrorResponse(c, http.StatusNotFound, errorResponseBuilder(err, additionalInfos...))
}

func Forbidden(c *gin.Context, err string, additionalInfos ...gin.H) {
	ErrorResponse(c, http.StatusForbidden, errorResponseBuilder(err, additionalInfos...))
}

func Conflict(c *gin.Context, err string, additionalInfos ...gin.H) {
	ErrorResponse(c, http.StatusConflict, errorResponseBuilder(err, additionalInfos...))
}

func ID(c *gin.Context, id interface{}) {
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func OK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": "ok"})
}

func TokenAndID(c *gin.Context, id ,token, expired interface{}) {
	c.JSON(http.StatusOK, gin.H{"expired": expired, "token": token, "id": id})
}

func errorResponseBuilder(err string, additionalInfos ...gin.H) gin.H {
	response := gin.H{"error": err}
	if len(additionalInfos) > 0 {
		infos := make([]gin.H, len(additionalInfos))
		for i := range additionalInfos {
			infos[i] = additionalInfos[i]
		}
		response["info"] = infos
	}

	return response
}

func ErrorResponse(c *gin.Context, code int, msg gin.H) {
	c.AbortWithStatusJSON(code, msg)
}

func stringOrNil(base sql.NullString) *string {
	if !base.Valid {
		return nil
	}

	return &base.String
}