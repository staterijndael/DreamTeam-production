package controller

import (
	"dt/config"
	"dt/events"
	"dt/managers/eventEmitter"
	"dt/managers/sms"
	"dt/rpc/servers/admin"
	"dt/rpc/servers/auth"
	"dt/rpc/servers/bug"
	"dt/rpc/servers/plebs"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	authMiddleware *jwt.GinJWTMiddleware
	existingTokens *sync.Map
	sqlDB          *gorm.DB
	jwtIdentityKey string
)

func NewServer(
	sqlStore *gorm.DB,
	userServer *plebs.Server,
	adminServer *admin.Server,
	emitter *eventEmitter.EventEmitter,
	conf *config.Config,
	smsMgr *sms.Manager,
	authServer *auth.Server,
	bugServer *bug.Server,
) (*gin.Engine, error) {
	jwtIdentityKey = conf.JWTIdentityKey
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		gin.Recovery(),
		Logger,
		cors.New(corsConfig()),
	)

	sqlDB = sqlStore
	existingTokens = &sync.Map{}
	var err error

	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		SigningAlgorithm: "HS512",
		Realm:            "dreamteam",
		Key:              []byte(config.VerySecretKey),
		Timeout:          config.JWTDuration,
		MaxRefresh:       config.JWTDuration,
		IdentityKey:      jwtIdentityKey,
		IdentityHandler:  identityHandler,
		PayloadFunc:      payload,
		Authorizator:     authorizator,
		Unauthorized:     unauthorized,
		Authenticator:    authenticator,
		LoginResponse:    loginOnSuccess,
		TokenLookup:      "header: Authorization, query: token",
	})

	if err != nil {
		return nil, err
	}

	smsGroup := r.Group("/sms")
	{
		smsGroup.GET("/status", smsStatus)
	}

	bugGroup := r.Group("/bug")
	{
		bugGroup.POST("", func(c *gin.Context) {
			bugServer.ServeHTTP(c.Writer, c.Request)
		})
	}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("", authServer.HTTPHandler)
	}

	userGroup := r.Group("")
	{
		userGroup.GET("/ws", CheckIsAuthMiddleware(), userServer.WSHandler)
		if conf.DebugMode {
			userGroup.POST("/rpc", userServer.HTTPHandler)
		}
	}

	if conf.DebugMode {
		bugGroup.GET("/smd", func(c *gin.Context) {
			c.JSON(http.StatusOK, bugServer.SMD())
		})

		authGroup.GET("/smd", func(c *gin.Context) {
			c.JSON(http.StatusOK, authServer.SMD())
		})

		userGroup.GET("/smd", func(c *gin.Context) {
			c.JSON(http.StatusOK, userServer.SMD())
		})

		r.GET("/user/:uid/auth/code", func(c *gin.Context) {
			id, _ := strconv.Atoi(c.Param("uid"))
			code, _ := smsMgr.Get(uint(id))
			c.String(00, "%d", code)
		})

		//TODO когда сделаем админ сервер вынести из дебага
		adminGroup := r.Group("/admin")
		{
			adminGroup.POST("/rpc", adminServer.HTTPHandler)
			//TODO adminGroup.GET("/ws", CheckIsAuthMiddleware(), adminServer.WSHandler)
			adminGroup.GET("/smd", func(c *gin.Context) {
				c.JSON(http.StatusOK, adminServer.SMD())
			})
		}

		r.Any("/fnskey/:fnsKey", FNSKeyDevSetter())
		r.GET("/rating/start/:oid", func(c *gin.Context) {
			oid, _ := strconv.Atoi(c.Param("oid"))
			now := time.Now()
			event := &events.RatingStartedEvent{
				Start:          now,
				End:            now.Add(time.Duration(time.Second.Nanoseconds() * int64(conf.RatingEventDebugDuration))),
				OrganizationID: uint(oid),
			}

			emitter.Emit(event)
		})

		r.GET("/box", func(c *gin.Context) {
			zenrpc.SMDBoxHandler(c.Writer, c.Request)
		})
		r.GET("/hi", func(c *gin.Context) {c.String(200, "hi")})
	}

	return r, nil
}

func corsConfig() cors.Config {
	cc := cors.DefaultConfig()
	cc.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"}
	cc.AllowHeaders = []string{"Origin", "Content-Type", "Content-Disposition", "X-Auth-Token", "Authorization"}
	cc.AllowCredentials = true
	cc.AllowWildcard = true
	cc.AllowFiles = true
	cc.AllowWebSockets = true
	cc.AllowAllOrigins = true
	return cc
}
