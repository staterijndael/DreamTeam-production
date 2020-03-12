package common

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader = &websocket.Upgrader{
	ReadBufferSize:  102400,
	WriteBufferSize: 102400,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRequestContext(oldCtx context.Context, u *models.User, con *websocket.Conn, r *http.Request) context.Context {
	oldCtx = context.WithValue(oldCtx, requestContext.HTTPRequestContextKey, r)
	oldCtx = context.WithValue(oldCtx, requestContext.ConnectionContextKey, con)
	oldCtx = context.WithValue(oldCtx, requestContext.UserContextKey, u)
	return oldCtx
}
