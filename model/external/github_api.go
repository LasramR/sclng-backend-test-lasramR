package external

import "github.com/LasramR/sclng-backend-test-lasramR/util"

// Represents a response from https://api.github.com/search/repositories
type RepositoriesResponse struct {
	TotalCount        int                        `json:"total_count"`
	IncompleteResults bool                       `json:"incomplete_results"`
	Repositories      []RepositoriesResponseItem `json:"items"`
}

func (r RepositoriesResponse) Items() []RepositoriesResponseItem {
	return r.Repositories
}

func (r RepositoriesResponse) Count() int {
	return r.TotalCount
}

// Represents an item from a RepositoriesResponse
type RepositoriesResponseItem struct {
	Id           int                                 `json:"id"`
	Name         string                              `json:"name"`
	FullName     string                              `json:"full_name"`
	Owner        ItemOwner                           `json:"owner"`
	Url          string                              `json:"url"`
	LanguagesUrl string                              `json:"languages_url"`
	License      util.NullableJsonField[ItemLicense] `json:"license"`
	UpdatedAt    string                              `json:"updated_at"` // TODO replace with time or custom marshaller
	Size         int                                 `json:"size"`
}

// Represents an item's owner from a RepositoriesResponseItem
type ItemOwner struct {
	Login string `json:"login"`
}

// Represents an item's license from a RepositoriesResponseItem
type ItemLicense struct {
	Key string `json:"key"`
}

// Represents a response from a language_url of a RepositoriesResponseItem
type Languages map[string]int
