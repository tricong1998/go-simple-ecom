package dto

type MetadataDto struct {
	Total   int64 `json:"total"`
	Page    int32 `json:"page"`
	PerPage int32 `json:"per_page"`
}
