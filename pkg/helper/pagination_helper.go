package helper

import (
	"net/http"
	"strconv"

	"github.com/savioruz/bake/internal/domain/model"
)

func ParsePagination(r *http.Request) *model.ProductPagination {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	pagination := &model.ProductPagination{
		Page:  1,
		Limit: 10,
		Sort:  "created_at",
		Order: "desc",
	}

	if page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
			pagination.Page = pageNum
		}
	}
	if limit != "" {
		if limitNum, err := strconv.Atoi(limit); err == nil && limitNum > 0 && limitNum <= 100 {
			pagination.Limit = limitNum
		}
	}
	if sort != "" {
		pagination.Sort = sort
	}
	if order != "" {
		pagination.Order = order
	}

	return pagination
}
