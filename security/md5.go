package security

import (
	"crypto/md5"
	"errors"
	"fmt"
)

//md5加密接口
func SuMd5(salt ...string) (data string, err error) {
	if len(salt) > 0 {
		mima := salt[0]
		if len(salt) > 1 && salt[1] != "" {
			bytes := []byte(mima + salt[1])
			sum := md5.Sum(bytes)
			return fmt.Sprintf("%x", sum), nil
		} else {
			bytes := []byte(mima)
			sum := md5.Sum(bytes)
			return fmt.Sprintf("%x", sum), nil
		}
	} else {
		return data, errors.New("数据为空")
	}
}
