package database

type Query struct {
	Filters  []Filter
	Sorts    []Sort
	Preloads []Preload
	Limit    int64
	Offset   int64
	Distinct bool
	GroupBy  []string
	Having   []string
}

type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

type Sort struct {
	Field string
	Desc  bool
}

type Preload struct {
	Field      string
	Conditions []interface{}
}

func NewQuery() *Query {
	return &Query{
		Filters:  make([]Filter, 0),
		Sorts:    make([]Sort, 0),
		Preloads: make([]Preload, 0),
	}
}

func (q *Query) Where(field, operator string, value interface{}) *Query {
	q.Filters = append(q.Filters, Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return q
}

func (q *Query) Equal(field string, value interface{}) *Query {
	return q.Where(field, "=", value)
}

func (q *Query) NotEqual(field string, value interface{}) *Query {
	return q.Where(field, "!=", value)
}

func (q *Query) GreaterThan(field string, value interface{}) *Query {
	return q.Where(field, ">", value)
}

func (q *Query) GreaterThanOrEqual(field string, value interface{}) *Query {
	return q.Where(field, ">=", value)
}

func (q *Query) LessThan(field string, value interface{}) *Query {
	return q.Where(field, "<", value)
}

func (q *Query) LessThanOrEqual(field string, value interface{}) *Query {
	return q.Where(field, "<=", value)
}

func (q *Query) Like(field string, value interface{}) *Query {
	return q.Where(field, "LIKE", value)
}

func (q *Query) In(field string, values interface{}) *Query {
	return q.Where(field, "IN", values)
}

func (q *Query) NotIn(field string, values interface{}) *Query {
	return q.Where(field, "NOT IN", values)
}

func (q *Query) IsNull(field string) *Query {
	return q.Where(field, "IS NULL", nil)
}

func (q *Query) IsNotNull(field string) *Query {
	return q.Where(field, "IS NOT NULL", nil)
}

func (q *Query) OrderBy(field string) *Query {
	q.Sorts = append(q.Sorts, Sort{
		Field: field,
		Desc:  false,
	})
	return q
}

func (q *Query) OrderByDesc(field string) *Query {
	q.Sorts = append(q.Sorts, Sort{
		Field: field,
		Desc:  true,
	})
	return q
}

func (q *Query) WithPreload(field string, conditions ...interface{}) *Query {
	q.Preloads = append(q.Preloads, Preload{
		Field:      field,
		Conditions: conditions,
	})
	return q
}

func (q *Query) WithLimit(limit int64) *Query {
	q.Limit = limit
	return q
}

func (q *Query) WithOffset(offset int64) *Query {
	q.Offset = offset
	return q
}

func (q *Query) WithPagination(page, size int64) *Query {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	q.Limit = size
	q.Offset = (page - 1) * size
	return q
}

func (q *Query) WithDistinct() *Query {
	q.Distinct = true
	return q
}

func (q *Query) WithGroupBy(fields ...string) *Query {
	q.GroupBy = append(q.GroupBy, fields...)
	return q
}

func (q *Query) WithHaving(conditions ...string) *Query {
	q.Having = append(q.Having, conditions...)
	return q
}

func (q *Query) Reset() *Query {
	q.Filters = make([]Filter, 0)
	q.Sorts = make([]Sort, 0)
	q.Preloads = make([]Preload, 0)
	q.Limit = 0
	q.Offset = 0
	q.Distinct = false
	q.GroupBy = make([]string, 0)
	q.Having = make([]string, 0)
	return q
}

func (q *Query) Clone() *Query {
	clone := &Query{
		Filters:  make([]Filter, len(q.Filters)),
		Sorts:    make([]Sort, len(q.Sorts)),
		Preloads: make([]Preload, len(q.Preloads)),
		Limit:    q.Limit,
		Offset:   q.Offset,
		Distinct: q.Distinct,
		GroupBy:  make([]string, len(q.GroupBy)),
		Having:   make([]string, len(q.Having)),
	}

	copy(clone.Filters, q.Filters)
	copy(clone.Sorts, q.Sorts)
	copy(clone.Preloads, q.Preloads)
	copy(clone.GroupBy, q.GroupBy)
	copy(clone.Having, q.Having)

	return clone
}
