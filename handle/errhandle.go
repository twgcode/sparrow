/**
@Author: wei-g
@Date:   2020/7/21 4:51 下午
@Description: 一些对外输出错误的信息的处理函数
*/

package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/twgcode/sparrow/util/data"
	"net/http"
)

// NoRoute 设置404
func NoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, data.RequestErrJson("not found"))
}

// NoRoute 设置405
func NoMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, data.RequestErrJson("method not allowed"))
}
