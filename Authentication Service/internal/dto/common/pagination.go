package common

type Pagination struct {
	PerPage int  `json:"perPage"`
	Page    int  `json:"page"`
	Total   int  `json:"total"`
	HasNext bool `json:"hasNext"`
}
