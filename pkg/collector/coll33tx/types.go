package coll33tx

import "github.com/florinutz/filme/pkg/collector"

const (
	TypeListNoPagination collector.DocType = iota
	TypeListWithPagination
	TypeDetail
)
