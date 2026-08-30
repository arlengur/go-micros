package ufo

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/model"
	repoConverter "github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/repository/converter"
	repoModel "github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/repository/model"
)

func (r *repository) Create(_ context.Context, info model.SightingInfo) (string, error) {
	newUUID := uuid.NewString()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[newUUID] = repoModel.Sighting{
		Uuid:      newUUID,
		Info:      repoConverter.SightingInfoToRepoModel(info),
		CreatedAt: time.Now(),
	}

	return newUUID, nil
}
