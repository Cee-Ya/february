package tools

import (
	"february/common/consts"
	"time"
)

func GetNowStr() string {
	return time.Now().Format(consts.DateFormatYmdhms)
}
