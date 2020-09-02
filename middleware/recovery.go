/**
* @Author: wei-g
* @Date:   2020/3/22 9:09 下午
* @Description: 设置 gin 框架recover中间件
 */

package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/twgcode/sparrow/util/log"
	"github.com/twgcode/sparrow/util/log/access"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	if stack {
		return func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					// Check for a broken connection, as it is not really a
					// condition that warrants a panic stack trace.
					var brokenPipe bool
					if ne, ok := err.(*net.OpError); ok {
						if se, ok := ne.Err.(*os.SyscallError); ok {
							if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
								brokenPipe = true
							}
						}
					}

					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					if brokenPipe {
						logger.Error(c.Request.URL.Path,
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
						// If the connection is dead, we can't write a status to it.
						c.Error(err.(error)) // nolint: errcheck
						c.Abort()
						return
					}
					// 根据 stack 的值的不同 记录log的内容不同
					logger.Error("[Recovery from panic]",
						zap.String("error", err.(string)),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())))
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}()
			c.Next()
		}
	} else {
		return func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					// Check for a broken connection, as it is not really a
					// condition that warrants a panic stack trace.
					var brokenPipe bool
					if ne, ok := err.(*net.OpError); ok {
						if se, ok := ne.Err.(*os.SyscallError); ok {
							if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
								brokenPipe = true
							}
						}
					}

					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					if brokenPipe {
						logger.Error(c.Request.URL.Path,
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
						// If the connection is dead, we can't write a status to it.
						c.Error(err.(error)) // nolint: errcheck
						c.Abort()
						return
					}
					// 根据 stack 的值的不同 记录log的内容不同
					logger.Error("[Recovery from panic]",
						zap.String("error", err.(string)),
						zap.String("request", string(httpRequest)),
					)
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}()
			c.Next()
		}
	}
}

// DefaultGinRecovery recover掉项目可能出现的panic
func DefaultGinRecovery(stack bool) gin.HandlerFunc {
	// 为了保证代码运行时的性能才把代码写的那么臃肿
	accessUp := false
	if access.LoggerMgr != nil && access.LoggerMgr.Logger != nil {
		accessUp = true
	}

	if accessUp { // 走 路由日志记录
		if stack {
			return func(c *gin.Context) {
				defer func() {
					if err := recover(); err != nil {
						// Check for a broken connection, as it is not really a
						// condition that warrants a panic stack trace.
						var brokenPipe bool
						if ne, ok := err.(*net.OpError); ok {
							if se, ok := ne.Err.(*os.SyscallError); ok {
								if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
									brokenPipe = true
								}
							}
						}

						httpRequest, _ := httputil.DumpRequest(c.Request, false)
						if brokenPipe {
							access.Error(c.Request.URL.Path,
								zap.Any("error", err),
								zap.String("request", string(httpRequest)),
							)
							// If the connection is dead, we can't write a status to it.
							c.Error(err.(error)) // nolint: errcheck
							c.Abort()
							return
						}
						// 根据 stack 的值的不同 记录log的内容不同
						access.Error("[Recovery from panic]",
							zap.String("error", err.(string)),
							zap.String("request", string(httpRequest)),
							zap.String("stack", string(debug.Stack())))
						c.AbortWithStatus(http.StatusInternalServerError)
					}
				}()
				c.Next()
			}
		} else {
			return func(c *gin.Context) {
				defer func() {
					if err := recover(); err != nil {
						// Check for a broken connection, as it is not really a
						// condition that warrants a panic stack trace.
						var brokenPipe bool
						if ne, ok := err.(*net.OpError); ok {
							if se, ok := ne.Err.(*os.SyscallError); ok {
								if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
									brokenPipe = true
								}
							}
						}

						httpRequest, _ := httputil.DumpRequest(c.Request, false)
						if brokenPipe {
							access.Error(c.Request.URL.Path,
								zap.Any("error", err),
								zap.String("request", string(httpRequest)),
							)
							// If the connection is dead, we can't write a status to it.
							c.Error(err.(error)) // nolint: errcheck
							c.Abort()
							return
						}
						// 根据 stack 的值的不同 记录log的内容不同
						access.Error("[Recovery from panic]",
							zap.String("error", err.(string)),
							zap.String("request", string(httpRequest)),
						)
						c.AbortWithStatus(http.StatusInternalServerError)
					}
				}()
				c.Next()
			}
		}
	} else { // 走 业务日志 路由日志记录
		if stack {
			return func(c *gin.Context) {
				defer func() {
					if err := recover(); err != nil {
						// Check for a broken connection, as it is not really a
						// condition that warrants a panic stack trace.
						var brokenPipe bool
						if ne, ok := err.(*net.OpError); ok {
							if se, ok := ne.Err.(*os.SyscallError); ok {
								if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
									brokenPipe = true
								}
							}
						}

						httpRequest, _ := httputil.DumpRequest(c.Request, false)
						if brokenPipe {
							log.Error(c.Request.URL.Path,
								zap.Any("error", err),
								zap.String("request", string(httpRequest)),
							)
							// If the connection is dead, we can't write a status to it.
							c.Error(err.(error)) // nolint: errcheck
							c.Abort()
							return
						}
						// 根据 stack 的值的不同 记录log的内容不同
						log.Error("[Recovery from panic]",
							zap.String("error", err.(string)),
							zap.String("request", string(httpRequest)),
							zap.String("stack", string(debug.Stack())))
						c.AbortWithStatus(http.StatusInternalServerError)
					}
				}()
				c.Next()
			}
		} else {
			return func(c *gin.Context) {
				defer func() {
					if err := recover(); err != nil {
						// Check for a broken connection, as it is not really a
						// condition that warrants a panic stack trace.
						var brokenPipe bool
						if ne, ok := err.(*net.OpError); ok {
							if se, ok := ne.Err.(*os.SyscallError); ok {
								if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
									brokenPipe = true
								}
							}
						}

						httpRequest, _ := httputil.DumpRequest(c.Request, false)
						if brokenPipe {
							log.Error(c.Request.URL.Path,
								zap.Any("error", err),
								zap.String("request", string(httpRequest)),
							)
							// If the connection is dead, we can't write a status to it.
							c.Error(err.(error)) // nolint: errcheck
							c.Abort()
							return
						}
						// 根据 stack 的值的不同 记录log的内容不同
						log.Error("[Recovery from panic]",
							zap.String("error", err.(string)),
							zap.String("request", string(httpRequest)),
						)
						c.AbortWithStatus(http.StatusInternalServerError)
					}
				}()
				c.Next()
			}
		}
	}
}
