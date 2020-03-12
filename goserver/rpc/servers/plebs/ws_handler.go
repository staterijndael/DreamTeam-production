package plebs

import (
	"dt/logwrap"
	"dt/models"
	"dt/rpc/servers/common"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strings"
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
	userAgent := strings.ToLower(c.Request.Header.Get("User-Agent"))
	if strings.Contains(userAgent, "androiddreamteam") || strings.Contains(userAgent, "iosdreamteam") {
		logwrap.Debug("pleb; user-agent: %s", c.Request.Header.Get("User-Agent"))
		rpc.ucm.Add(currentUser, connection)
	} else {
		logwrap.Debug("dash; user-agent: %s", c.Request.Header.Get("User-Agent"))
		rpc.dcm.Add(currentUser, connection)
	}

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
