package response

import (
	"fmt"
	"net/http"

	"git.exclouds.org/ww/greet_ws/global"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
}

const (
	HTTP_STATUS_OK        int = 200
	STATUS_SUCCESS        int = 1
	STATUS_FAIL           int = -1
	ERROR                 int = -1
	INTERNAL_SERVER_ERROR int = 500
	AUTH_EXPIRED          int = 4001
)

func Result(code, status int, data interface{}, msg string, c *gin.Context) {

	if code != HTTP_STATUS_OK {
		errMsg := fmt.Sprintf("%v %v;参数：%v\n 详细信息：%v", c.Request.Method, c.FullPath(), c.Params, data)
		global.LOGGER.Error(errMsg)
	}

	// 原始
	c.JSON(http.StatusOK, Response{
		code,
		status,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(HTTP_STATUS_OK, STATUS_SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(HTTP_STATUS_OK, STATUS_SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(HTTP_STATUS_OK, STATUS_SUCCESS, data, "调用成功", c)
}

func OkDetailed(data interface{}, message string, c *gin.Context) {
	Result(HTTP_STATUS_OK, STATUS_SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	errMsg := fmt.Sprintf("%v %v;参数：%v", c.Request.Method, c.FullPath(), c.Params)
	global.LOGGER.Warn(errMsg)
	Result(HTTP_STATUS_OK, STATUS_FAIL, map[string]interface{}{}, "调用失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	errMsg := fmt.Sprintf("%v %v;参数：%v\n 详细信息：%v", c.Request.Method, c.FullPath(), c.Params, message)
	global.LOGGER.Warn(errMsg)
	Result(HTTP_STATUS_OK, STATUS_FAIL, map[string]interface{}{}, message, c)
}

func FailWithData(data interface{}, c *gin.Context) {
	errMsg := fmt.Sprintf("%v %v;参数：%v; 详细信息：%v", c.Request.Method, c.FullPath(), c.Params, data)
	global.LOGGER.Warn(errMsg)
	Result(HTTP_STATUS_OK, STATUS_FAIL, data, "调用失败", c)
}

func FailWithDetailed(httpStatusCode int, data interface{}, message string, c *gin.Context) {
	errMsg := fmt.Sprintf("%v %v;参数：%v;httpStatusCode:%v \n 详细信息：%v", c.Request.Method, c.FullPath(), c.Params, httpStatusCode, message)
	global.LOGGER.Warn(errMsg)
	Result(httpStatusCode, STATUS_FAIL, data, message, c)
}

func FailWithCustomCode(httpStatusCode int, data interface{}, c *gin.Context) {
	errMsg := fmt.Sprintf("%v %v;参数：%v;httpStatusCode:%v \n 详细信息：%v", c.Request.Method, c.FullPath(), c.Params, httpStatusCode, data)
	global.LOGGER.Error(errMsg)
	Result(httpStatusCode, STATUS_FAIL, map[string]interface{}{}, errMsg, c)
}

func InternalServerError(message interface{}, c *gin.Context) {
	msg := fmt.Sprintf("%v", message)
	errMsg := fmt.Sprintf("%v %v；参数：%v \n 详细信息：%v", c.Request.Method, c.FullPath(), c.Params, message)
	global.LOGGER.Error(errMsg)
	Result(INTERNAL_SERVER_ERROR, STATUS_FAIL, map[string]interface{}{}, msg, c)
}
