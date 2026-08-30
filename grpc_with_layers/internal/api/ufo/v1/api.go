package v1

import (
	"github.com/arlengur/go-micros/layers/internal/service"
	ufoV1 "github.com/arlengur/go-micros/layers/pkg/proto/ufo/v1"
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
