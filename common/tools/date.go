package tools

import (
	"ai-report/common/consts"
	"time"
)

func GetNowStr() string {
	return time.Now().Format(consts.DateFormatYmdhms)
}
