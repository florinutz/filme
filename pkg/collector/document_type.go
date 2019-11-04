package collector

type DocType int8

const (
	TypeImdbMovie DocType = iota
	TypeImdbSeries
	TypeImdbSeriesUnfinished

	Type1337xListNoPagination
	Type1337xListWithPagination
	Type1337xDetail
)

type DocumentHasType interface {
	GetDocumentType() DocType
}
