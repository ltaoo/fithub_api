package models

import "time"

type MediaResource struct {
	Id        int       `json:"id"`
	MediaType int       `json:"media_type"` // 大致类型，1image 2video 3audio 4markdown 这样
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Size      int       `json:"size"` // 字节
	Duration  int       `json:"duration"`
	CoverKey  string    `json:"cover_key"`
	Title     string    `json:"title"`
	Filename  string    `json:"filename"` // 原始文件名
	Filetype  string    `json:"filetype"` // image/png 这样格式的，就是 mime-type
	Hash      string    `json:"hash"`
	Key       string    `json:"key"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`

	CreatorId int `json:"creator_id"`
}

func (MediaResource) TableName() string {
	return "MEDIA_RESOURCE"
}
