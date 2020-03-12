package requestContext

import (
	"context"
	"dt/models"
	"github.com/gorilla/websocket"
)

const (
	ConnectionContextKey  = "currentConnection"
	UserContextKey        = "currentUser"
	UserQueryKey          = "me"
	HTTPRequestContextKey = "Request"
)

func WebsocketFromContext(ctx context.Context) (*websocket.Conn, bool) {
	con, ok := ctx.Value(ConnectionContextKey).(*websocket.Conn)
	return con, ok
}

func CurrentUser(ctx context.Context) (u *models.User) {
	u, _ = ctx.Value(UserContextKey).(*models.User)
	return
}
