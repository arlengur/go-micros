package ufo

import (
	"context"

	"github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepository.Create(ctx, info)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
