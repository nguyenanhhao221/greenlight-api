package data

type Filters struct {
	Page     int
	PageSize int
	Sort     string
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return f.Page
}
