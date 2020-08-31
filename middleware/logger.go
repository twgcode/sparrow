/**
* @Author: wei-g
* @Date: 2020/3/20 4:42 下午
* @Description: 设置 gin 框架日志中间件
 */

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)

	}
}

// GinLoggerMap  记录gin的请求日志
// loggerMap中的2个key: access 代表web访问日志, business 代表业务日志
func GinLoggerMap(loggerMap map[string]*zap.Logger) gin.HandlerFunc {
	access, oka := loggerMap["access"]
	business, okb := loggerMap["business"]
	if oka && okb && access != nil && business != nil {
		return func(c *gin.Context) {
			path := c.Request.URL.Path
			query := c.Request.URL.RawQuery
			business.Info(path,
				zap.String("method", c.Request.Method), zap.String("path", path), zap.String("query", query),
				zap.String("ip", c.ClientIP()), zap.String("user-agent", c.Request.UserAgent()),
			)
			start := time.Now()
			c.Next()
			cost := time.Since(start)
			access.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)

		}
	} else if oka && access != nil { // 只有访问日志使用
		return func(c *gin.Context) {
			path := c.Request.URL.Path
			query := c.Request.URL.RawQuery
			start := time.Now()
			c.Next()
			cost := time.Since(start)
			access.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)

		}
	} else if okb && business != nil { // 只有业务日志
		return func(c *gin.Context) {
			path := c.Request.URL.Path
			query := c.Request.URL.RawQuery
			business.Info(path,
				zap.String("method", c.Request.Method), zap.String("path", path), zap.String("query", query),
				zap.String("ip", c.ClientIP()), zap.String("user-agent", c.Request.UserAgent()),
			)
			c.Next()
		}
	}
	// 不记录日志
	return func(c *gin.Context) {
		c.Next()
	}
}
