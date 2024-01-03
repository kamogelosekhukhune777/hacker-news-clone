package models

import (
	"errors"
	"math"
	"strings"
)

type Filter struct {
	Page     int
	PageSize int
	OrderBy  string
	Query    string
}

type Metadata struct {
	CurrentPage  int
	PageSize     int
	FirstPage    int
	NextPage     int
	PrevPage     int
	LastPage     int
	TotalRecords int
}

func (f *Filter) Validate() error {
	if f.Page <= 0 || f.Page > 10000 {
		return errors.New("invalid page range:1 to 10 thousand")
	}

	if f.PageSize <= 0 || f.PageSize > 100 {
		return errors.New("invalid page size: 1 to 100")
	}

	return nil
}

// sets odering by clause to query template string
func (f *Filter) addOrdering(query string) string {
	if f.OrderBy == "popular" {
		return strings.Replace(query, "#orderby#", "ORDER BY votes desc, p.created_at desc", 1)
	}

	return strings.Replace(query, "#orderby#", "ORDER BY p.created_at desc", 1)
}

// sets WHERE the clause to query template string
func (f *Filter) addWhere(query string) string {
	if len(f.Query) > 0 {
		return strings.Replace(query, "#where#", "WHERE LOWER(P.title) LIKE $1", 1)
	}

	return strings.Replace(query, "#where#", "", 1)
}

// sets the limitoffset of query template
func (f *Filter) addLimitOffset(query string) string {
	if len(f.Query) > 0 {
		return strings.Replace(query, "#limit#", "LIMIT $2 OFFSET $3", 1)
	}

	return strings.Replace(query, "#where#", "LIMIT $1 OFFSET $2", 1)
}

// adds the limitOffset, Where, and, order by clause
func (f *Filter) applyTemplate(query string) string {
	return f.addLimitOffset(f.addWhere(f.addOrdering(query)))
}

func (f *Filter) limit() int {
	return f.PageSize
}

func (f *Filter) offset() int {
	return (f.Page - 1) * f.PageSize
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	meta := Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}

	meta.NextPage = meta.CurrentPage + 1
	meta.PrevPage = meta.CurrentPage - 1

	if meta.CurrentPage <= meta.FirstPage {
		meta.PrevPage = 0
	}

	return meta
}
