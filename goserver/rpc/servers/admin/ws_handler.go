package admin

import (
	"dt/logwrap"
	"dt/models"
	"dt/rpc/servers/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

func (rpc *Server) WSHandler(c *gin.Context) {
	connection, err := common.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logwrap.NetworkError("ws upgrade", "", "",
			"upgrade connection failed with err=%v", err,
		)
		return
	}
	defer connection.Close()

	currentUser := c.Keys["currentUser"].(*models.User)

	for {
		mt, message, err := connection.ReadMessage()

		// normal closure
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			break
		}
		// abnormal closure
		if err != nil {
			logwrap.Debug("read message failed with err=%v", err)
			break
		}

		data, err := rpc.Do(common.NewRequestContext(c.Request.Context(), currentUser, connection, c.Request), message)
		if err != nil {
			logwrap.Debug("marshal json response failed with err=%v", err)
			connection.WriteControl(websocket.CloseInternalServerErr, nil, time.Time{})
			break
		}

		if err = connection.WriteMessage(mt, data); err != nil {
			logwrap.Debug("write response failed with err=%v", err)
			connection.WriteControl(websocket.CloseInternalServerErr, nil, time.Time{})
			break
		}
	}
}
