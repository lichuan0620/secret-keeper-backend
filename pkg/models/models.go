package models

import "time"

const (
	Version = "2021-12-23"
)

type Box struct {
	Id             string          `json:"Id,omitempty" bson:"_id,omitempty"`
	CreatedAt      *time.Time      `json:"CreatedAt,omitempty" bson:"CreatedAt,omitempty"`
	Body           string          `json:"Body,omitempty" bson:"Body,omitempty"`
	EmojiFeedbacks map[string]uint `json:"EmojiFeedbacks,omitempty" bson:"EmojiFeedbacks,omitempty"`
	LastViewed     *time.Time      `json:"LastViewed,omitempty" bson:"LastViewed,omitempty"`
}

type QueueItem struct {
	Id    string `json:"Id"`
	Score int64  `json:"Score"`
}

type CreateBoxRequest struct {
	Body string `json:"Body,omitempty"`
}

type CreateBoxResponse Box

type AddBoxEmoji struct {
	Id             string          `json:"Id,omitempty"`
	EmojiFeedbacks map[string]uint `json:"EmojiFeedbacks,omitempty"`
}

type AddBoxEmojiRequest AddBoxEmoji

type AddBoxEmojiResponse AddBoxEmoji

type ViewBoxResponse Box

type SyncRequest QueueItem

type SyncResponse struct{}

type DequeueResponse struct {
	Id string `json:"Id"`
}
