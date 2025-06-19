package utils

import (
	"fmt"
	"git.exclouds.org/ww/greet_ws/global"
	"github.com/avct/uasurfer"
	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"strings"
)

type ClientUtils struct {
}

const ipDbFile = "utils/ip2region.db"

// 获取浏览器+版本
func (*ClientUtils) UserDevice(c *gin.Context) string {
	userAgent := c.Request.Header.Get("User-Agent")
	ua := uasurfer.Parse(userAgent)
	browser := ua.Browser
	return strings.ReplaceAll(browser.Name.String(), "Browser", "") + fmt.Sprintf(" %d.%d.%d", browser.Version.Major, browser.Version.Minor, browser.Version.Patch)
}

// 获取IP
func (*ClientUtils) ClientIp(c *gin.Context) string {
	return c.ClientIP()
}

// 根据IP获取地点信息
func (*ClientUtils) ClientInfo(ip string) (info ip2region.IpInfo, err error) {
	region, err := ip2region.New(ipDbFile)
	defer region.Close()
	if err != nil {
		return info, err
	}

	return region.MemorySearch(ip)
}

// 根据IP获取最近的地点
func (c *ClientUtils) ClientAddress(ip string) string {
	info, err := c.ClientInfo(ip)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("%v", err))
		return ""
	}

	if info.City != "0" {
		return info.City
	}

	if info.Province != "0" {
		return info.Province
	}

	if info.Region != "0" {
		return info.Region
	}
	return info.Country
}
