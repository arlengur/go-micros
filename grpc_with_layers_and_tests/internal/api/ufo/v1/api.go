package v1

import (
	"github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/service"
	ufoV1 "github.com/arlengur/go-micros/grpc_with_layers_and_tests/pkg/proto/ufo/v1"
)

type api struct {
	ufoV1.UnimplementedUFOServiceServer

	ufoService service.UFOService
}

func NewAPI(ufoService service.UFOService) *api {
	return &api{
		ufoService: ufoService,
	}
}
