package logic

import (
	"context"
	"errors"
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

func NewTagMatchLogic(ctx context.Context) *TagMatchLogic {
	onceTagMatchLogic.Do(func() {
		tagMatchLogic = &TagMatchLogic{
			indexes: &component.SyncMap[int64, *tagEntityIndexes]{},
		}
	})
	return tagMatchLogic
}

func (t *TagMatchLogic) AddEntities(ctx context.Context, projectId int64, entities []*entity.BizEntity) error {
	slog.InfoContext(ctx, "add entities", slog.Int64("projectId", projectId), slog.Any("entities", entities))

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
			key := makeIndexKey(tmpTag.Tag, tmpTag.Value)
			tmpBitmap, ok := index.data.Load(key)
			if !ok {
				tmpBitmap = roaring64.New()
			}
			index.data.Store(key, tmpBitmap)

			tmpBitmap.Add(uint64(tmpEntity.Id))
		}
	}

	slog.InfoContext(ctx, "add entities done", slog.Int64("projectId", projectId), slog.Any("entities", entities))
	return nil
}

func (t *TagMatchLogic) MatchEntities(ctx context.Context, projectId int64, tags []*entity.TagEntity) (entityIds []int64, err error) {
	slog.InfoContext(ctx, "match entities", slog.Int64("projectId", projectId), slog.Any("tags", tags))

	index, ok := t.indexes.Load(projectId)
	if !ok {
		slog.ErrorContext(ctx, "project index not found", slog.Int64("projectId", projectId))
		return nil, errors.New("project not found")
	}

	var bitmaps []*roaring64.Bitmap
	for _, tmpTag := range tags {
		tmpKey := makeIndexKey(tmpTag.Tag, tmpTag.Value)
		tmpBitmap, ok := index.data.Load(tmpKey)
		if !ok {
			slog.WarnContext(ctx, "not bitmap data for tag", slog.String("tag", tmpTag.Tag), slog.String("value", tmpTag.Value))
			return nil, nil
		}
		bitmaps = append(bitmaps, tmpBitmap)
	}

	if len(bitmaps) == 0 {
		slog.WarnContext(ctx, "no result data for tag")
		return nil, nil
	}

	resultBitmap := bitmaps[0]
	for i := 1; i < len(bitmaps); i++ {
		i++
	}

	entityIdsArrayInUint64 := resultBitmap.ToArray()
	for _, tmpEntityIdInUint64 := range entityIdsArrayInUint64 {
		entityIds = append(entityIds, int64(tmpEntityIdInUint64))
	}

	slog.InfoContext(ctx, "match entities done", slog.Int64("projectId", projectId), slog.Any("tags", tags), slog.Any("entityIds", entityIds))
	return entityIds, nil
}

func makeIndexKey(tag string, value string) string {
	return fmt.Sprintf("%s#%s", tag, value)
}
