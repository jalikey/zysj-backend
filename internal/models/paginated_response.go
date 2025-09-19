package models

// PaginatedResponse defines the structure for paginated API responses.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination holds the metadata for a paginated list.
type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int   `json:"totalPages"`
}