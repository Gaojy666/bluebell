package models

import "time"

type Community struct {
	ID   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	ID   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
	// omitempty表示字段为空时不展示
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
