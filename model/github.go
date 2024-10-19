package model

import "github.com/LasramR/sclng-backend-test-lasramR/util"

// Represents a response from https://api.github.com/search/repositories
type GithubApiProjectsResponse struct {
	TotalCount        int                             `json:"total_count"`
	IncompleteResults bool                            `json:"incomplete_results"`
	Items             []GithubApiProjectsResponseItem `json:"items"`
}

// Represents an item from a RepositoriesResponse
type GithubApiProjectsResponseItem struct {
	Id           int                                                         `json:"id"`
	Name         string                                                      `json:"name"`
	FullName     string                                                      `json:"full_name"`
	Owner        GithubApiProjectsReponseItemOwner                           `json:"owner"`
	Url          string                                                      `json:"url"`
	LanguagesUrl string                                                      `json:"languages_url"`
	License      util.NullableJsonField[GithubApiProjectsReponseItemLicense] `json:"license"`
	UpdatedAt    string                                                      `json:"updated_at"` // TODO replace with time or custom marshaller
	Size         int                                                         `json:"size"`
}

// Represent an item's owner from a RepositoriesResponseItem
type GithubApiProjectsReponseItemOwner struct {
	Login string `json:"login"`
}

type GithubApiProjectsReponseItemLicense struct {
	Key string `json:"key"`
}

// Represent a response from a language_url of a RepositoriesResponseItem
type GithubApiLanguageURLResponse map[string]int
