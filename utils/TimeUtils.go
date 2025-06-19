package utils

import "time"

type TimeUtils struct {
}

// 时间字符串转换为unix时间戳、时间
// Str2TimeAndStamp("2020-11-26 05:00:00")
func (*TimeUtils) Str2TimeAndStamp(toBeConvertedTimeStr string) (timestamp int64, theTime time.Time, err error) {
	//获取本地location
	//toBeConvertedTimeStr := "2020-11-26 05:00:00"                            //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02 15:04:05"                                        //转化所需模板
	loc, err := time.LoadLocation("Local")                                     //重要：获取时区
	theTime, err = time.ParseInLocation(timeLayout, toBeConvertedTimeStr, loc) //使用模板在对应时区转化为time.time类型
	timestamp = theTime.Unix()                                                 //转化为时间戳 类型是int64
	//fmt.Println(theTime)                 //打印输出theTime 2015-01-01 15:15:00 +0800 CST
	//fmt.Println(timestamp, theTime, err) //打印输出时间戳 1420041600
	return
}

// 时间转换为字符串。 theTime时间；cutFrom：要截取的开始位置 cutTo：要截取的结束位置
// Time2Str(time.Now(), 0, 10)
func (*TimeUtils) Time2Str(theTime time.Time, cutFrom int, cutTo int) string {
	layout := "2006-01-02 15:04:05.000Z"
	str := theTime.Format(layout[cutFrom:cutTo])
	return str
}
