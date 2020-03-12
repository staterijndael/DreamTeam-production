package auth

import "github.com/gin-gonic/gin"

func (rpc *Server) HTTPHandler(c *gin.Context) {
	rpc.ServeHTTP(c.Writer, c.Request)
}
