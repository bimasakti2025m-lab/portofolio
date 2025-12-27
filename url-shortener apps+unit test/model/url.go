package model

import "time"

type Url struct {
	Id        uint32    `json:"id"`
	LongUrl   string    `json:"long_url" binding:"required"`
	ShortCode string    `json:"short_code"`
	UserId    uint32    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
