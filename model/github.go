package model

import "github.com/LasramR/sclng-backend-test-lasramR/util"

type Repository struct {
	FullName      string                         `json:"full_name"`
	Owner         string                         `json:"owner"`
	Repository    string                         `json:"repository"`
	RepositoryUrl string                         `json:"repository_url"`
	Languages     Language                       `json:"languages"`
	License       util.NullableJsonField[string] `json:"license"`
	Size          int                            `json:"size"`
	UpdatedAt     string                         `json:"updated_at"`
}

type Language map[string]LanguageStats

type LanguageStats struct {
	Bytes int `json:"bytes"`
}
