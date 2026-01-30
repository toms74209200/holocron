package domain

type Pagination struct {
	limit  int
	offset int
}

func ToPagination(limit, offset *int) Pagination {
	l := 20
	if limit != nil && *limit > 0 && *limit <= 100 {
		l = *limit
	}
	o := 0
	if offset != nil && *offset >= 0 {
		o = *offset
	}
	return Pagination{limit: l, offset: o}
}

func (p Pagination) Limit() int {
	return p.limit
}

func (p Pagination) Offset() int {
	return p.offset
}
