package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type EnhanceUtils struct{}

// RandomValue 从切片中随机返回一个值
func RandomValue[T any](list []T) (T, error) {
	if len(list) == 0 {
		var zeroValue T // 返回类型的零值
		return zeroValue, fmt.Errorf("the slice is empty")
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(list))
	return list[randomIndex], nil
}

// ConvertToStringSlice 将 []interface{} 转换为 []string
func ConvertToStringSlice(input []interface{}) ([]string, error) {
	var result []string
	for _, v := range input {
		// 类型断言
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("non-string element found: %v", v)
		}
		result = append(result, str)
	}
	return result, nil
}

// ConvertInterfaceToStringSlice 将 interface{} 转换为 []string
func ConvertInterfaceToStringSlice(input interface{}) ([]string, error) {
	// 先断言为切片类型
	slice, ok := input.([]interface{})
	if !ok {
		return nil, errors.New("input is not a []interface{}")
	}

	// 转换 []interface{} 为 []string
	var result []string
	for _, v := range slice {
		// 检查每个元素是否是 string
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("non-string element found: %v", v)
		}
		result = append(result, str)
	}

	return result, nil
}
