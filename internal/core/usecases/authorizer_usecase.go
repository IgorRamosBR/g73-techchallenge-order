package usecases

import (
	"github.com/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/g73-techchallenge-order/internal/infra/drivers/authorizer"
	log "github.com/sirupsen/logrus"
)

type AuthorizerUsecase interface {
	AuthorizeUser(cpf string) (dto.AuthorizedUser, error)
}

type authorizerUsecase struct {
	authorizer authorizer.Authorizer
}

func NewAuthorizerUsecase(authorizer authorizer.Authorizer) AuthorizerUsecase {
	return authorizerUsecase{
		authorizer: authorizer,
	}
}

func (u authorizerUsecase) AuthorizeUser(cpf string) (dto.AuthorizedUser, error) {
	authorizerResponse, err := u.authorizer.AuthorizeUser(cpf)
	if err != nil {
		log.Errorf("failed to authorize user, error: %v", err)
		return dto.AuthorizedUser{}, err
	}

	return authorizerResponse.User, nil
}
