package request

import "errors"

type PresignedUrlRequestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	FileName    string `json:"file_name"`
}

func (r PresignedUrlRequestBody) Validate() error {
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
