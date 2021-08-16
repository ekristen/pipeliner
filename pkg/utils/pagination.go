package utils

import (
	"net/http"
	"strconv"
)

// Paginator --
type Paginator struct {
	Limit  int
	Offset int
	Total  int
}

// Paginate --
func Paginate(r *http.Request, count int64) (*Paginator, error) {
	paginator := &Paginator{
		Limit:  20,
		Offset: 0,
		Total:  0,
	}

	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("per_page")

	if perPage != "" {
		p, err := strconv.Atoi(perPage)
		if err != nil {
			return nil, err
		}
		if p == 0 {
			p = 1
		}

		paginator.Limit = p
	} else {
		paginator.Limit = 25
	}

	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}

		if p == 0 {
			p = 1
		}

		paginator.Offset = p * paginator.Limit
	} else {
		paginator.Offset = 0
	}

	pages := int(count) / paginator.Limit
	paginator.Total = pages

	return paginator, nil
}
