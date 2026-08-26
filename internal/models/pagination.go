package models

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedRestaurants struct {
	Data       []Restaurant `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

type PaginatedMenus struct {
	Data       []Menu `json:"data"`
	Pagination Pagination `json:"pagination"`
}