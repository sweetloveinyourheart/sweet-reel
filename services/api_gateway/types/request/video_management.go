package request

import (
	"errors"

	"github.com/gofrs/uuid"
)

type PresignedUrlRequestBody struct {
	ChannelID   string `json:"channel_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FileName    string `json:"file_name"`
}

func (r PresignedUrlRequestBody) Validate() error {
	if r.ChannelID == "" {
		return errors.New("channel_id should not be empty")
	}

	if _, err := uuid.FromString(r.ChannelID); err != nil {
		return errors.New("channel_id should be a valid uuid string")
	}

	if r.Title == "" {
		return errors.New("title should not be empty")
	}

	if r.Description == "" {
		return errors.New("description should not be empty")
	}

	if r.FileName == "" {
		return errors.New("file_name should not be empty")
	}

	return nil
}
