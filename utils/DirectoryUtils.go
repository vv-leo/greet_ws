package utils

import (
	"git.exclouds.org/ww/greet_ws/global"
	"go.uber.org/zap"
	"os"
)

type DirectoryUtils struct {
}

var directoryUtils DirectoryUtils

// @title    PathExists
// @description   文件目录是否存在
// @auth                     （2020/04/05  20:22）
// @param     path            string
// @return    err             error

func (_ *DirectoryUtils) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// @title    createDir
// @description   批量创建文件夹
// @auth                     （2020/04/05  20:22）
// @param     dirs            string
// @return    err             error

func (_ *DirectoryUtils) CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := directoryUtils.PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.LOGGER.Debug("create directoryUtils" + v)
			err = os.MkdirAll(v, os.ModePerm)
			if err != nil {
				global.LOGGER.Error("create directoryUtils"+v, zap.Any(" error:", err))
			}
		}
	}
	return err
}
