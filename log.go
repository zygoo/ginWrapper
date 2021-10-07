package ginWrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type GinLog struct {
	FilterParams []string
	EncryptParam []string
	SkipPaths    []string
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func APILogger(ginLog GinLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		bodyLogWriter := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = bodyLogWriter
		c.Next()
		end := time.Now()
		logger.WithFields(logger.Fields{
			"p":              bodyLogWriter.body.String(),
			"q":              getRequestInfo,
			"client_ip":      c.ClientIP(),
			"status":         c.Writer.Status(),
			"method":         c.Request.Method,
			"path":           c.FullPath(),
			"raw_query":      c.Request.URL.Query(),
			"user_agent":     c.Request.UserAgent(),
			"end":            end.Format("2006/01/02-15:04:05"),
			"duration":       end.Sub(start).Seconds(),
			"content_type":   c.ContentType(),
			"content_length": c.Request.ContentLength,
			"x-forward-for":  c.GetHeader("X-Forward-For"),
		}).Info("APILogger")
	}
}

func getRequestInfo(c *gin.Context) gin.H {
	params := map[string]interface{}{}
	if c.Request.Method != "GET" {
		switch c.ContentType() {
		case "application/json":
			data, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				defer c.Request.Body.Close()
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(data))
			m := map[string]interface{}{}
			err = json.Unmarshal(data, &m)
			if err != nil {
				params["raw"] = string(data)
			} else {
				for k, v := range m {
					params[k] = fmt.Sprint(v)
				}
			}
		case "application/x-www-form-urlencoded":
			_ = c.Request.ParseForm()
			if c.Request.Form != nil {
				for k, v := range c.Request.Form {
					params[k] = v
				}
			}
		}
	}
	return params
}
