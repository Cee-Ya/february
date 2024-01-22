package memory

import (
	"february/common"
	"february/entity"
)

func InitMemoryCache() {
	common.MemoryCache = entity.NewMemoryCache()
}
