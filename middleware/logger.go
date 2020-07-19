/**
* @Author: wei-g
* @Date: 2020/3/20 4:42 下午
* @Description: 设置 gin 框架日志中间件
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 接收gin框架默认的日志
func GinLogger(loggerMap []map[string]*zap.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
	}
}
