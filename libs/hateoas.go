package libs

import (
	"fmt"
	"net/url"
)

type HateoasLinks struct {
	Next *string `json:"next"`
	Back *string `json:"back"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

func BuildHateoasGlobal(basePath string, page, limit, totalItems int, query url.Values) (Pagination, HateoasLinks) {
	totalPages := (totalItems + limit - 1) / limit

	qNext := cloneQuery(query)
	qBack := cloneQuery(query)

	var nextLink *string
	var backLink *string

	if page < totalPages {
		qNext.Set("page", fmt.Sprintf("%d", page+1))
		qNext.Set("limit", fmt.Sprintf("%d", limit))
		nextStr := basePath + "?" + qNext.Encode()
		nextLink = &nextStr
	} else {
		nextLink = nil
	}

	if page > 1 {
		qBack.Set("page", fmt.Sprintf("%d", page-1))
		qBack.Set("limit", fmt.Sprintf("%d", limit))
		backStr := basePath + "?" + qBack.Encode()
		backLink = &backStr
	} else {
		backLink = nil
	}

	pagination := Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	links := HateoasLinks{
		Next: nextLink,
		Back: backLink,
	}

	return pagination, links
}

func cloneQuery(v url.Values) url.Values {
	clone := url.Values{}
	for key, vals := range v {
		for _, val := range vals {
			clone.Add(key, val)
		}
	}
	return clone
}
