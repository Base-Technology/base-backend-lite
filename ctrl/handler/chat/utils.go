package chat

import (
	"fmt"
	"time"

	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func updateBalance(limit *database.ChatGPTLimit, response *string) {
	// TODO: use tokenizer from openai to count token
	tokenCount := len(*response)
	seelog.Infof("tokenCount: %d", tokenCount)
	limit.DailyLeftCallCount = max(0, limit.DailyLeftCallCount-1)
	limit.DailyLeftTokenCount = max(0, limit.DailyLeftTokenCount-tokenCount)
	limit.TotalTokenLeftCount = max(0, limit.TotalTokenLeftCount-tokenCount)
}

// verify if there is enough token to use
func enoughBalance(limit *database.ChatGPTLimit, query *string) bool {
	// TODO: use tokenizer from openai to count token
	if limit.DailyLeftCallCount <= 0 || limit.DailyLeftTokenCount <= 0 || limit.TotalTokenLeftCount <= 0 {
		return false
	}
	return true
}

// increase balance for referer
func IncreaseBalanceForReferer(referrer_id uint) {
	limit := &database.ChatGPTLimit{}
	if err := database.GetInstance().Where(database.ChatGPTLimit{UserID: referrer_id}).First(&limit).Error; err != nil {
		msg := fmt.Sprintf("IncreaseBalanceForReferer: UserId=%v not found [%v]", referrer_id, err)
		seelog.Errorf(msg)
		return
	}
	// increase maximum by 25% & refresh daily limit
	limit.MaxDailyCallCount += 5
	limit.MaxDailyTokenCount += 3000
	limit.TotalTokenLeftCount += 15000
	limit.DailyLeftCallCount = limit.MaxDailyCallCount
	limit.DailyLeftTokenCount = limit.MaxDailyTokenCount
	limit.LastResetTime = time.Now()
	database.GetInstance().Save(limit)
}
