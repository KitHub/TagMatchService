package service

import (
	"context"
	"sync"

	"github.com/KitHub/TagMatchService/logic"
	"github.com/KitHub/protocols/TagMatchService"
)

var (
	tagMatchServiceAPIServiceInstance *TagMatchServiceAPIService
	tagMatchServiceAPIServiceOnce     sync.Once
)

type TagMatchServiceAPIService struct {
	TagMatchService.UnimplementedTagMatchServiceAPIServer
	demoLogic *logic.DemoLogic
}

// Match implements [TagMatchService.TagMatchServiceAPIServer].
func (t *TagMatchServiceAPIService) Match(context.Context, *TagMatchService.TagMatchServiceRequest) (*TagMatchService.TagMatchServiceResponse, error) {
	panic("unimplemented")
}

func NewTagMatchServiceAPIService(tagMatchLogic *logic.DemoLogic) *TagMatchServiceAPIService {
	tagMatchServiceAPIServiceOnce.Do(func() {
		tagMatchServiceAPIServiceInstance = &TagMatchServiceAPIService{
			demoLogic: tagMatchLogic,
		}
	})
	return tagMatchServiceAPIServiceInstance
}
