package v1

import (
	"context"

	"github.com/arlengur/go-micros/layers/internal/converter"
	ufoV1 "github.com/arlengur/go-micros/layers/pkg/proto/ufo/v1"
)

func (a *api) Create(ctx context.Context, req *ufoV1.CreateRequest) (*ufoV1.CreateResponse, error) {
	uuid, err := a.ufoService.Create(ctx, converter.UFOInfoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &ufoV1.CreateResponse{
		Uuid: uuid,
	}, nil
}
