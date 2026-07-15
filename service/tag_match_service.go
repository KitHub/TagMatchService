package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/KitHub/TagMatchService/logic"
	"github.com/KitHub/protocols/TagMatchService"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	tagMatchServiceAPIServiceInstance *TagMatchServiceAPIService
	tagMatchServiceAPIServiceOnce     sync.Once
)

type TagMatchServiceAPIService struct {
	TagMatchService.UnimplementedTagMatchServiceAPIServer
	tagMatchLogic *logic.TagMatchLogic
}

// AddEntities implements [TagMatchService.TagMatchServiceAPIServer].
func (t *TagMatchServiceAPIService) AddEntities(ctx context.Context, req *TagMatchService.AddEntitiesRequest) (rsp *TagMatchService.AddEntitiesResponse, err error) {
	slog.InfoContext(ctx, "add entities", slog.Any("req", req))
	err = req.Validate()
	if err != nil {
		slog.DebugContext(ctx, "invalid request", slog.Any("req", req), slog.Any("error", err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters")
	}

	rsp = &TagMatchService.AddEntitiesResponse{
		ErrCode: 0,
		ErrMsg:  "ok",
	}
	slog.InfoContext(ctx, "add entities done", slog.Any("rsp", rsp))
	return rsp, nil
}

// MatchEntites implements [TagMatchService.TagMatchServiceAPIServer].
func (t *TagMatchServiceAPIService) MatchEntites(ctx context.Context, req *TagMatchService.MatchEntitesRequest) (rsp *TagMatchService.MatchEntitesResponse, err error) {
	slog.InfoContext(ctx, "match entities", slog.Any("req", req))
	slog.InfoContext(ctx, "match entities done", slog.Any("rsp", rsp))
	return nil, nil
}

func NewTagMatchServiceAPIService(tagMatchLogic *logic.TagMatchLogic) *TagMatchServiceAPIService {
	tagMatchServiceAPIServiceOnce.Do(func() {
		tagMatchServiceAPIServiceInstance = &TagMatchServiceAPIService{
			tagMatchLogic: tagMatchLogic,
		}
	})
	return tagMatchServiceAPIServiceInstance
}
