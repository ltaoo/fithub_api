package models

type Pagination struct {
	Page       int    `json:"page" form:"page,default=1"`
	PageSize   int    `json:"page_size" form:"page_size,default=10"`
	NextMarker string `json:"next_marker" form:"next_marker"`
}
