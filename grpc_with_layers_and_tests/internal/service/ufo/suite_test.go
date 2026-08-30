package ufo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/arlengur/go-micros/grpc_with_layers_and_tests/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	ufoRepository *mocks.UFORepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.ufoRepository = mocks.NewUFORepository(s.T())

	s.service = NewService(
		s.ufoRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
