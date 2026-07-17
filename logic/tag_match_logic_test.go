package logic_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/KitHub/TagMatchService/entity"
	"github.com/KitHub/TagMatchService/logic"
)

type tagAndValues struct {
	tag    string
	values []string
}

func BenchmarkTagMatchLogic_MatchEntities(b *testing.B) {
	var projectId int64 = 0

	var globalTagsCount int32 = 10000
	var globalTagValuesCount int32 = 10

	var entitiesCount int32 = 100
	var entityTagsCount int32 = 10

	var requestedTagsCount int32 = 10

	ctx := b.Context()
	now := time.Now()
	seed1, seed2 := now.UnixMilli(), now.UnixMilli()
	rnd := rand.New(rand.NewPCG(uint64(seed1), uint64(seed2)))
	globalTagAndValuesSlice := makeGlobalTagAndValues(ctx, globalTagsCount, globalTagValuesCount)
	entities := makeEntities(ctx, rnd, entitiesCount, entityTagsCount, globalTagAndValuesSlice)
	tagMatchLogic := logic.NewTagMatchLogic(ctx)
	err := tagMatchLogic.AddEntities(ctx, projectId, entities)
	if err != nil {
		b.Fatal(err)
	}

	matchTagEntities := makeMatchTagEntities(ctx, rnd, requestedTagsCount, globalTagAndValuesSlice)

	b.ResetTimer()

	tagMatchLogic.MatchEntities(ctx, projectId, matchTagEntities)

}

func makeTagValues(ctx context.Context, valuesCount int32) []string {
	var result []string
	for i := range valuesCount {
		result = append(result, fmt.Sprintf("%d", i))
	}
	return result
}

func makeGlobalTagAndValues(ctx context.Context, tagsCount int32, tagValuesCount int32) []*tagAndValues {
	var result []*tagAndValues
	for i := range tagsCount {
		result = append(result, &tagAndValues{
			tag:    fmt.Sprintf("Tag$%d", i),
			values: makeTagValues(ctx, tagValuesCount),
		})
	}
	return result
}

func makeEntities(ctx context.Context, rnd *rand.Rand, entitiesCount int32, entityTagsCount int32, globalTagAndValuesSlice []*tagAndValues) []*entity.BizEntity {
	var result []*entity.BizEntity
	for range entitiesCount {
		var tmpTags []*entity.TagEntity
		for range entityTagsCount {
			tmpTagAndValues := globalTagAndValuesSlice[rnd.Int32N(int32(len(globalTagAndValuesSlice)))]
			tmpTags = append(tmpTags, &entity.TagEntity{
				Tag:   tmpTagAndValues.tag,
				Value: tmpTagAndValues.values[rnd.Int32N(int32(len(tmpTagAndValues.values)))],
			})
		}
		tmpEntity := &entity.BizEntity{
			Id:   rnd.Int64(),
			Tags: tmpTags,
		}
		result = append(result, tmpEntity)
	}
	return result
}

func makeMatchTagEntities(ctx context.Context, rnd *rand.Rand, tagEntitiesCount int32, tagAndValuesSlice []*tagAndValues) []*entity.TagEntity {
	var result []*entity.TagEntity
	for range tagEntitiesCount {
		tmpTagAndValues := tagAndValuesSlice[rnd.Int32N(int32(len(tagAndValuesSlice)))]
		tmpTagEntity := &entity.TagEntity{
			Tag:   tmpTagAndValues.tag,
			Value: tmpTagAndValues.values[rnd.Int32N(int32(len(tmpTagAndValues.values)))],
		}
		result = append(result, tmpTagEntity)
	}
	return result
}
