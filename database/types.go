package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone        string `gorm:"size:11;unique"`
	Name         string `gorm:"size:20"`
	Password     string `gorm:"size:100"`
	Area         string `gorm:"size:20"`
	School       string `gorm:"size:20;index"`
	PrivateKey   string `gorm:"size:100"`
	Introduction string `gorm:"size:20"`
	Avatar       string
}

type Post struct {
	gorm.Model
	CreatorID     uint
	Title         string `gorm:"size:100"`
	Content       string
	PostPointedID uint

	Creator     *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostPointed *Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Image struct {
	gorm.Model
	PostID uint
	Source string

	Post *Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Comment struct {
	gorm.Model
	PostID           uint
	CreatorID        uint
	CommentPointedID uint
	Content          string

	Post           *Post    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Creator        *User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CommentPointed *Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Follow struct {
	gorm.Model
	UserID      uint `gorm:"uniqueIndex:distinct_follow"`
	FollowingID uint `gorm:"uniqueIndex:distinct_follow"`

	User      *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Following *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Like struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex:distinct_like"`
	PostID uint `gorm:"uniqueIndex:distinct_like"`

	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post *Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Collect struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex:distinct_collect"`
	PostID uint `gorm:"uniqueIndex:distinct_collect"`

	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post *Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ChatGPTLimit struct {
	UserID uint `gorm:"primaryKey"`

	DailyLeftCallCount  int `gorm:"default:50"`
	DailyLeftTokenCount int `gorm:"default:4000"`
	TotalTokenLeftCount int `gorm:"default:20000"`

	MaxDailyCallCount  int `gorm:"default:50"`
	MaxDailyTokenCount int `gorm:"default:4000"`

	LastResetTime time.Time
}
