package models

import "github.com/gofrs/uuid"

type VideoSplitterMessage struct {
	VideoID  uuid.UUID             `json:"video_id"`
	Metadata VideoSplitterMetadata `json:"metadata"`
}

type VideoSplitterMetadata struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}
