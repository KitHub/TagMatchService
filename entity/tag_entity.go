package entity

type TagEntity struct {
	Tag   string
	Value string
}

type BizEntity struct {
	Id   int64
	Tags []*TagEntity
}
