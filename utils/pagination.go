package utils

import (
	"net/url"
	"strconv"

	"gopkg.in/gorp.v2"

	"github.com/juju/errors"
)

// Direction represents the sort direction
// swagger:strfmt Direction
// enum:ASC,DESC
type Direction string

// sort directions
const (
	Asc  Direction = "ASC"
	Desc Direction = "DESC"
)

const (
	DefaultPageSize = 20
)

// NewPage creates *Page from query parameters(e.g. pageStart=2&pageSize=50&sort=id&order=DESC).
func NewPage(params url.Values) *Page {
	if params == nil {
		return &Page{
			PageStart: 0,
			PageSize:  DefaultPageSize,
			Order:     "DESC",
			Sort:      "",
			Keyword:   "",
		}
	}

	page := &Page{}
	n, err := strconv.ParseInt(params.Get("pageStart"), 10, 64)
	if err == nil {
		page.PageStart = int(n)
	}
	n, err = strconv.ParseInt(params.Get("pageSize"), 10, 64)
	if err == nil {
		page.PageSize = int(n)
	}
	if page.PageSize == 0 {
		page.PageSize = DefaultPageSize
	}
	page.Sort = params.Get("sort")
	page.Order = params.Get("order")
	if len(page.Order) == 0 {
		page.Order = "DESC"
	}
	page.Keyword = params.Get("keyword")

	return page
}

// Page contains pagination information
type Page struct {
	// page start index, start from 0
	// example:1
	// default:1
	PageStart int `db:"pageStart" json:"pageStart"`

	// size of page
	// example:10
	// default:10
	PageSize int `db:"pageSize" json:"pageSize"`

	// sort order, ASC or DESC, default to DESC
	// example:DESC
	Order string `json:"order"`

	// field name to sort
	// example:id
	Sort string `json:"sort"`

	// keyword to query
	// example:arch
	Keyword string `db:"keyword" json:"keyword"`
}

func (p Page) ToParams() url.Values {
	params := make(url.Values)
	params.Add("pageStart", strconv.Itoa(p.PageStart))
	params.Add("pageSize", strconv.Itoa(p.PageSize))
	if len(p.Sort) > 0 {
		params.Add("sort", p.Sort)
	}
	if len(p.Order) > 0 {
		params.Add("order", p.Order)
	}
	return params
}

// String returns readable string represents the page
func (p Page) String() string {
	return p.ToParams().Encode()
}

// Validate implements Validatable interface
func (p Page) Validate() error {
	if p.PageStart < 0 {
		return errors.Errorf("pageStart must greater or equal than 0, got %v", p.PageStart)
	}
	if p.PageSize < 1 {
		return errors.Errorf("size must greater than 0, got %v", p.PageSize)
	}

	if len(p.Order) > 0 {
		if p.Order != "ASC" && p.Order != "DESC" {
			return errors.Errorf("invalid sort order: %v", p.Order)
		}
	}

	return nil
}

// PageResponse is the response to a page request
// swagger:model
type PageResponse struct {
	*Page
	// total item
	//   example:100
	Total int `db:"total" json:"total"`

	// payload
	Data interface{} `db:"data" json:"data"`
}

// NewPageResponse creates a PageResponse
func NewPageResponse(page *Page, total int, data interface{}) *PageResponse {
	return &PageResponse{
		Page:  page,
		Total: total,
		Data:  data,
	}
}

func LoadPage(tx gorp.SqlExecutor, q *Query, page *Page, holder interface{}, countQueries ...*Query) (int, error) {
	cq := q
	if len(countQueries) > 0 {
		if v := countQueries[0]; v != nil {
			cq = v
		}
	}
	n, err := cq.Count(tx)
	if err != nil {
		return -1, errors.Trace(err)
	}
	if page != nil {
		if err := page.Validate(); err != nil {
			return -1, errors.Trace(err)
		}

		if len(page.Sort) > 0 {
			if len(page.Order) == 0 {
				page.Order = "DESC"
			}
			q.ClearOrderBy().OrderByString(page.Sort, page.Order)
		}

		q.Offset(page.PageStart).Limit(page.PageSize)
	}
	if _, err = q.FetchAll(tx, holder); err != nil {
		return -1, errors.Trace(err)
	}
	return int(n), nil
}

func LoadPageResponse(tx gorp.SqlExecutor, q *Query, page *Page, holder interface{}, countQueries ...*Query) (*PageResponse, error) {
	count, err := LoadPage(tx, q, page, holder, countQueries...)
	if err != nil {
		return nil, err
	}

	return NewPageResponse(page, count, holder), nil
}

func LoadPage2(tx gorp.SqlExecutor, q *Query, page *Page, holder interface{}) error {
	if page != nil {
		if err := page.Validate(); err != nil {
			return errors.Trace(err)
		}

		if len(page.Sort) > 0 {
			if len(page.Order) == 0 {
				page.Order = "DESC"
			}
			q.ClearOrderBy().OrderByString(page.Sort, page.Order)
		}

		q.Offset(page.PageStart).Limit(page.PageSize)
	}
	if _, err := q.FetchAll(tx, holder); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func LoadPageResponse2(tx gorp.SqlExecutor, q *Query, page *Page, holder interface{}, countQuery *Query) (*PageResponse, error) {
	countQueryStr, vals, err := countQuery.ValQuery()
	if err != nil {
		return nil, q.QueryValError(err, countQueryStr, nil)
	}
	count, err := tx.SelectInt(countQueryStr, vals...)
	if err != nil {
		return nil, q.QueryValError(err, countQueryStr, nil)
	}

	if err := LoadPage2(tx, q, page, holder); err != nil {
		return nil, err
	}

	return NewPageResponse(page, int(count), holder), nil
}
