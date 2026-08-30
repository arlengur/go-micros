package ufo

import (
	"context"

	"github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/model"
	repoConverter "github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Sighting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoSighting, ok := r.data[uuid]
	if !ok {
		return model.Sighting{}, model.ErrSightingNotFound
	}

	return repoConverter.SightingToModel(repoSighting), nil
}
