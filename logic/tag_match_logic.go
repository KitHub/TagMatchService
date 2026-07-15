package logic

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/KitHub/TagMatchService/component"
	"github.com/KitHub/TagMatchService/entity"
	"github.com/RoaringBitmap/roaring/v2/roaring64"
)

var tagMatchLogic *TagMatchLogic
var onceTagMatchLogic sync.Once

type tagEntityIndexes struct {
	data *component.SyncMap[string, *roaring64.Bitmap] // key=tag#tagValue
}

type TagMatchLogic struct {
	indexes *component.SyncMap[int64, *tagEntityIndexes] // key=projectId
}

func (t *TagMatchLogic) Hello(ctx context.Context) error {
	panic("unimplemented")
}

func NewTagMatchLogic() *TagMatchLogic {
	onceTagMatchLogic.Do(func() {
		tagMatchLogic = &TagMatchLogic{}
	})
	return tagMatchLogic
}

func (t *TagMatchLogic) AddEntities(ctx context.Context, projectId int64, entities []*entity.BizEntity) error {
	slog.InfoContext(ctx, "add entities", slog.Any("entities", entities))

	if len(entities) == 0 {
		slog.WarnContext(ctx, "the count of entities is 0")
		return nil
	}

	var index *tagEntityIndexes
	index, ok := t.indexes.Load(projectId)
	if !ok {
		index = &tagEntityIndexes{
			data: &component.SyncMap[string, *roaring64.Bitmap]{},
		}
		t.indexes.Store(projectId, index)
	}

	for _, tmpEntity := range entities {
		for _, tmpTag := range tmpEntity.Tags {
			key := fmt.Sprintf("%s#%s", tmpTag.Tag, tmpTag.Value)
			tmpBitmap, ok := index.data.Load(key)
			if !ok {
				tmpBitmap = roaring64.New()
			}
			index.data.Store(key, tmpBitmap)

			tmpBitmap.Add(uint64(tmpEntity.Id))
		}
	}

	slog.InfoContext(ctx, "add entities done", slog.Any("entities", entities))
	return nil
}

func (t *TagMatchLogic) MatchEntities(ctx context.Context, projectId int64, tags []*entity.TagEntity) (entityIds []int64, err error) {
	return nil, nil
}
