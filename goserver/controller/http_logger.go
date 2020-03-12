package controller

import (
	"bytes"
	"dt/logwrap"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"strings"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger(c *gin.Context) {
	var b strings.Builder
	buf, err := ioutil.ReadAll(c.Request.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	if err != nil {
		logwrap.Error("[error logging http]: %s", err.Error())
		c.Next()
		return
	}
	
	_, err = io.Copy(&b, rdr1)
	if err != nil {
		logwrap.Error("[error logging http]: %s", err.Error())
		c.Next()
		return
	}

	rbw := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = rbw
	c.Next()
	logwrap.Info(`{"url":"%s", "method":"%s", "code":%d}`, c.Request.URL, c.Request.Method, c.Writer.Status())
	input := b.String()
	if len(input) == 0 {
		logwrap.Debug(
			`{"url":"%s", "method":"%s", "code":%d, "output":%s}`,
			c.Request.URL,
			c.Request.Method,
			c.Writer.Status(),
			rbw.body.String(),
		)
	} else {
		logwrap.Debug(
			`{"url":"%s", "method":"%s", "code":%d, "input":%s, "output":%s}`,
			c.Request.URL,
			c.Request.Method,
			c.Writer.Status(),
			b.String(),
			rbw.body.String(),
		)
	}
}
