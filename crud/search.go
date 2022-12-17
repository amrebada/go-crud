package crud

type SortOrder string

var SortDirection = []SortOrder{"ASC", "DESC"}

var DefaultSearchableFields = []string{}

var DefaultSortableFields = []string{"id", "created_at", "updated_at"}

var SearchOperators = []string{"exact", "contains", "startswith", "endswith"}

type PaginationDetails struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Pages int64 `json:"pages"`
}
