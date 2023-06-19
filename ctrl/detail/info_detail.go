package detail

import "time"

type UserDetail struct {
	ID     uint   `json:"id"`
	Name   string `json:"username"`
	Avatar string `json:"avatar"`
	Sex    string `json:"sex"`
}
type UserDeatailMore struct {
	ID           uint   `json:"id"`
	Name         string `json:"username"`
	Area         string `json:"area"`
	School       string `json:"school"`
	Introduction string `json:"introduction"`
	Avatar       string `json:"avatar"`
	IMTPUserID   string `json:"imtp_user_id"`
	Sex          string `json:"sex"`
}
type RequestDetail struct {
	ID       uint      `json:"id"`
	Name     string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Message  string    `json:"message"`
	Status   string    `json:"status"`
	CreateAt time.Time `json:"create_at"`
	Sex      string    `json:"sex"`
}
type GroupDetail struct {
	ID          uint   `json:"group_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatorID   uint   `json:"creator_id"`
	MembersNum  int    `json:"member_num"`
	Avatar      string `json:"avatar"`
}
type PostDeatail struct {
	ID            uint      `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	CreateAt      time.Time `json:"create_at"`
	CreatorId     uint      `json:"creatorid"`
	CreatorName   string    `json:"Creator_name"`
	CreatorAvatar string    `json:"creator_avatar"`
	CommentCount  int64     `json:"comment_count"`
	LikeCount     int64     `json:"like_count"`
	CollectCount  int64     `json:"collect_count"`
}
type CommentDeatail struct {
	CommentID        uint      `json:"comment_id"`
	CreatorID        uint      `json:"creator_id"`
	CreatorName      string    `json:"creator_name"`
	CreatorAvatar    string    `json:"creator_avatar"`
	CreateAt         time.Time `json:"create_at"`
	Content          string    `json:"content"`
	Level            uint      `json:"level"`
	CommentPointedID uint      `json:"commentpointed_id"`
}
type ChatGPTLimitDetail struct {
	DailyLeftCallCount  int `json:"daily_left_call_count"`
	DailyLeftTokenCount int `json:"daily_left_token_count"`
	TotalTokenLeftCount int `json:"total_token_left_count"`

	MaxDailyCallCount  int `json:"max_daily_call_count"`
	MaxDailyTokenCount int `json:"max_daily_token_count"`
}
