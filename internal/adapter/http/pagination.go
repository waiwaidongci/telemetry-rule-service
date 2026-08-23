package httpadapter

import (
	"net/http"
	"strconv"
)

type Page struct {
	Offset int
	Limit  int
}

func pageFrom(r *http.Request) Page {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	return Page{Offset: offset, Limit: limit}
}
func pageBounds(length int, p Page) (int, int) {
	start := p.Offset
	if start > length {
		start = length
	}
	end := start + p.Limit
	if end > length {
		end = length
	}
	return start, end
}
func setPagination(w http.ResponseWriter, total int, p Page) {
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	w.Header().Set("X-Page-Offset", strconv.Itoa(p.Offset))
	w.Header().Set("X-Page-Limit", strconv.Itoa(p.Limit))
}
