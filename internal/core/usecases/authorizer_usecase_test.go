package usecases

import (
	"errors"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	mock_authorizer "github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/authorizer/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthorizerUsecase_AuthorizeUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	authorizer := mock_authorizer.NewMockAuthorizer(ctrl)

	authorizerUseCase := NewAuthorizerUsecase(authorizer)

	cpf := "123456789"

	authorizer.EXPECT().AuthorizeUser(cpf).Times(1).Return(dto.AuthorizerResponse{}, errors.New("internal server error"))

	authorizedUser, err := authorizerUseCase.AuthorizeUser(cpf)

	assert.Empty(t, authorizedUser)
	assert.Error(t, err, "failed to authorize user, error: internal server error")

	expectedAuthorizerResponse := dto.AuthorizerResponse{
		IsAuthorized: true,
		Message:      "user is authorized",
		User: dto.AuthorizedUser{
			CPF:   cpf,
			Name:  "user1",
			Email: "user1@gmail.com",
		},
	}

	authorizer.EXPECT().AuthorizeUser(cpf).Times(1).Return(expectedAuthorizerResponse, nil)

	authorizedUser, err = authorizerUseCase.AuthorizeUser(cpf)

	assert.Equal(t, expectedAuthorizerResponse.User, authorizedUser)
	assert.NoError(t, err)
}
