package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/KitHub/TagMatchService/entity"
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

	var bizEntities []*entity.BizEntity
	for _, tmpEntity := range req.Data {
		var tags []*entity.TagEntity
		for _, tmpTag := range tmpEntity.Tags {
			tags = append(tags, &entity.TagEntity{
				Tag:   tmpTag.Tag,
				Value: tmpTag.Value,
			})
		}
		bizEntities = append(bizEntities, &entity.BizEntity{
			Id:   tmpEntity.EntityId,
			Tags: tags,
		})
	}

	err = t.tagMatchLogic.AddEntities(ctx, req.GetProjectId(), bizEntities)
	if err != nil {
		slog.ErrorContext(ctx, "add entities failed", slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "server error")
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

	err = req.Validate()
	if err != nil {
		slog.DebugContext(ctx, "invalid request", slog.Any("req", req), slog.Any("error", err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters")
	}

	var tags []*entity.TagEntity
	for _, tmpTag := range req.GetTags() {
		tags = append(tags, &entity.TagEntity{
			Tag:   tmpTag.Tag,
			Value: tmpTag.Value,
		})
	}

	entityIds, err := t.tagMatchLogic.MatchEntities(ctx, req.GetProjectId(), tags)
	if err != nil {
		slog.ErrorContext(ctx, "match entities failed", slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "server error")
	}

	rsp = &TagMatchService.MatchEntitesResponse{
		ErrCode: 0,
		ErrMsg:  "ok",
		Data: &TagMatchService.MatchEntitesResponseData{
			EntityIds: entityIds,
		},
	}

	slog.InfoContext(ctx, "match entities done", slog.Any("rsp", rsp))
	return rsp, nil
}

func NewTagMatchServiceAPIService(tagMatchLogic *logic.TagMatchLogic) *TagMatchServiceAPIService {
	tagMatchServiceAPIServiceOnce.Do(func() {
		tagMatchServiceAPIServiceInstance = &TagMatchServiceAPIService{
			tagMatchLogic: tagMatchLogic,
		}
	})
	return tagMatchServiceAPIServiceInstance
}
