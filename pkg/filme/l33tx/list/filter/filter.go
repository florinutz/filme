package filter

import "github.com/florinutz/filme/pkg/collector/coll33tx/list"

type Filter struct {
	MaxItems uint
	Seeders  IntVal
	Leechers IntVal
	Size     IntVal
}

type IntVal struct {
	Min uint
	Max uint
}

func (f *Filter) IsValid(item list.Item) bool {
	return true
}
