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

// reset limit if last reset time is more than 24 hours
func resetBalance(limit *database.ChatGPTLimit) {
	if time.Since(limit.LastResetTime).Hours() > 24 {
		limit.DailyLeftCallCount = limit.MaxDailyCallCount
		limit.DailyLeftTokenCount = limit.MaxDailyTokenCount
		limit.LastResetTime = time.Now()
	}
}

func updateBalance(limit *database.ChatGPTLimit, usage *ChatGPTProxyUsage) {
	seelog.Infof("tokenUsage: %v", usage)
	limit.DailyLeftCallCount = max(0, limit.DailyLeftCallCount-1)
	limit.DailyLeftTokenCount = max(0, limit.DailyLeftTokenCount-usage.TotalTokens)
	limit.TotalTokenLeftCount = max(0, limit.TotalTokenLeftCount-usage.TotalTokens)
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
	tx := database.GetInstance().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err := tx.Commit().Error; err != nil {
			msg := fmt.Sprintf("IncreaseBalanceForReferer: UserId=%v commit error [%v]", referrer_id, err)
			seelog.Errorf(msg)
			tx.Rollback()
		}
	}()

	limit := &database.ChatGPTLimit{}
	if err := tx.Where(database.ChatGPTLimit{UserID: referrer_id}).First(limit).Error; err != nil {
		msg := fmt.Sprintf("IncreaseBalanceForReferer: UserId=%v not found [%v]", referrer_id, err)
		seelog.Errorf(msg)
		tx.Rollback()
		return
	}
	// increase maximum by 25% & refresh daily limit
	limit.MaxDailyCallCount += 5
	limit.MaxDailyTokenCount += 3000
	limit.TotalTokenLeftCount += 15000
	limit.DailyLeftCallCount = limit.MaxDailyCallCount
	limit.DailyLeftTokenCount = limit.MaxDailyTokenCount
	limit.LastResetTime = time.Now()
	if err := tx.Save(limit).Error; err != nil {
		msg := fmt.Sprintf("IncreaseBalanceForReferer: UserId=%v save error [%v]", referrer_id, err)
		seelog.Errorf(msg)
		tx.Rollback()
		return
	}
}
