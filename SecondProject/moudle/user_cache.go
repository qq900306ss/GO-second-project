package moudle

import (
	"context"
	"github/qq900306ss/SecondProject/utils"
	"time"
)

/*
*
設計再現用戶到緩存
*
*/
func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	utils.Red.Set(ctx, key, val, timeTTL)
}
