package utils

import (
	"net/http/httptest"
	"testing"
)

func TestPaginationDefaults(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	paginate, err := Paginate(req, 100)
	if err != nil {
		t.Error(err)
	}

	if paginate.Limit != 25 {
		t.Errorf("limit is not the default 25")
	}
	if paginate.Offset != 0 {
		t.Errorf("offset is not the default 0")
	}
	if paginate.Total != 4 {
		t.Errorf("total pages is not the default 4 with 100 records")
	}
}

func TestPaginationPageChange(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo?page=5", nil)

	paginate, err := Paginate(req, 500)
	if err != nil {
		t.Error(err)
	}

	if paginate.Limit != 25 {
		t.Errorf("limit is not 25")
	}
	if paginate.Offset != 125 {
		t.Errorf("offset is not 125")
	}
	if paginate.Total != 20 {
		t.Errorf("total pages is not the expected value")
	}
}

func TestPaginationPageSize(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo?page=2&per_page=50", nil)

	paginate, err := Paginate(req, 500)
	if err != nil {
		t.Error(err)
	}

	if paginate.Limit != 50 {
		t.Errorf("limit is not 50 (actual: %d)", paginate.Limit)
	}
	if paginate.Offset != 100 {
		t.Errorf("offset is not 100 (actual: %d)", paginate.Offset)
	}
	if paginate.Total != 10 {
		t.Errorf("total pages is not the expected value")
	}
}
