package controller

import (
	"dt/models"
	"dt/utils"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type authUserScheme struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func payload(data interface{}) jwt.MapClaims {
	if user, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			jwtIdentityKey: user.ID,
		}
	}
	return jwt.MapClaims{}
}

func identityHandler(c *gin.Context) interface{} {
	id := uint(jwt.ExtractClaims(c)[jwtIdentityKey].(float64))

	//if err != nil {
	//	return nil
	//}

	return &models.User{
		Model: gorm.Model{
			ID: id,
		},
	}
}

func authorizator(data interface{}, c *gin.Context) bool {
	user, ok := data.(*models.User)
	if !ok || user == nil {
		return false
	}

	err := sqlDB.First(user, user.ID).Error
	if err != nil && err != models.EmptyNicknameModelErr {
		return false
	}

	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys["currentUser"] = user

	//token := jwt.GetToken(c)
	//if _, ok := existingTokens.Load(token); !ok {
	//	return false
	//}

	//existingTokens.Delete(jwt.GetToken(c))

	return true
}

func authenticator(c *gin.Context) (interface{}, error) {
	var authData authUserScheme
	if err := c.ShouldBind(&authData); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	var user models.User
	err := sqlDB.Where(&models.User{Phone: utils.FormatPhone(authData.Phone),
		Password: authData.Password}).First(&user).Error
	if err != nil && err != models.EmptyNicknameModelErr {
		return nil, jwt.ErrFailedAuthentication
	}

	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys["currentUser"] = &user

	return &user, nil
}

func loginOnSuccess(c *gin.Context, code int, token string, expired time.Time) {
	user := c.Keys["currentUser"].(*models.User)
	existingTokens.Store(token, true)
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"token": token, "id": user.ID})
}

func unauthorized(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}

func LoginHandler() gin.HandlerFunc {
	return authMiddleware.LoginHandler
}

func CheckIsAuthMiddleware() gin.HandlerFunc {
	return authMiddleware.MiddlewareFunc()
}
