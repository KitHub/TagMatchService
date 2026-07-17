package entity

import (
	"encoding/json"
	"fmt"
)

type TagEntity struct {
	Tag   string
	Value string
}

func (t *TagEntity) String() string {
	return fmt.Sprintf("%s:%s", t.Tag, t.Value)
}

type BizEntity struct {
	Id   int64
	Tags []*TagEntity
}

func (b *BizEntity) String() string {
	bytes, _ := json.Marshal(b)
	return string(bytes)
}
