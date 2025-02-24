package valueobject

type Pagination struct {
	Page     int
	PageSize int
	Total    int64
}

func NewPagination(page, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}
