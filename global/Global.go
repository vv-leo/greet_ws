package global

import (
	"errors"
	"fmt"
	"git.exclouds.org/ww/greet_ws/config"
	"github.com/go-redis/redis"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

const (
	WsTimeOut    = 30
	TableListKey = "TableList"

	ORIGINAL_CONTACT_TABLE = "contact"
	ORIGINAL_MESSAGE_TABLE = "message"

	HOT_DATA_EXPIRE = time.Minute
)

var (
	SERVER_CONFIG config.ServerConfig
	CONTACT_DB    *gorm.DB
	MESSAGE_DB    *gorm.DB
	LOGGER        *zap.Logger
	REDIS         *redis.Client
	VIPER         *viper.Viper
)

func ShardedTableName(originalTableName string, seatId int64) string {

	return originalTableName + "_" + strconv.FormatInt(seatId%SERVER_CONFIG.SystemConfig.ShardedTableCount, 10)
}

func GenerateUniqueId() string {
	charset := "0123456789"
	length := 10

	timestamp := time.Now().UnixMilli() // 毫秒级时间戳
	randomID, err := gonanoid.Generate(charset, length)
	if err != nil {
		LOGGER.Error(err.Error())
		return "-1"
	}

	id := fmt.Sprintf("%d%s", timestamp, randomID)
	return id
}

func GetKeepTime() int64 {
	return time.Now().Unix() - SERVER_CONFIG.SystemConfig.MaxKeepTime
}

func IsDuplicateRecord(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "Duplicate entry")
}
