package logic

import (
	"context"
	"log/slog"
	"sync"

	"github.com/KitHub/TagMatchService/entity"
)

var tagMatchLogic *TagMatchLogic
var onceTagMatchLogic sync.Once

type TagMatchLogic struct {
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

func (t *TagMatchLogic) AddEntities(ctx context.Context, entities []*entity.Entity) error {
	slog.InfoContext(ctx, "add entities", slog.Any("entities", entities))

	if len(entities) == 0 {
		slog.WarnContext(ctx, "the count of entities is 0")
		return nil
	}

	slog.InfoContext(ctx, "add entities done", slog.Any("entities", entities))
	return nil
}
