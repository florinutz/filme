package imdb

import "github.com/florinutz/filme/pkg/collector"

const (
	TypeImdbMovie collector.DocType = iota
	TypeImdbSeries
	TypeImdbSeriesUnfinished
)
