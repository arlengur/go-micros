package ufo

import (
	"context"

	"github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/model"
)

func (s *service) Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	err := s.ufoRepository.Update(ctx, uuid, updateInfo)
	if err != nil {
		return err
	}

	return nil
}
