package domain

type SearchKeyword string

func ToSearchKeyword(q *string) *SearchKeyword {
	if q == nil || *q == "" {
		return nil
	}
	k := SearchKeyword(*q)
	return &k
}
