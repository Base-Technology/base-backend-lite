package chat

import (
	//"fmt"

	"github.com/Base-Technology/base-backend-lite/database"
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
	//fmt.Println("tokenCount: ", tokenCount)
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
