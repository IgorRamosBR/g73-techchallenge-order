package authorizer

import (
	"encoding/json"
	"errors"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/http"
)

var ErrUnauthorized = errors.New("customer unauthorized")

type Authorizer interface {
	AuthorizeUser(cpf string) (dto.AuthorizerResponse, error)
}

type authorizer struct {
	client        http.HttpClient
	authorizerUrl string
}

func NewAuthorizer(client http.HttpClient, authorizerUrl string) Authorizer {
	return authorizer{
		client:        client,
		authorizerUrl: authorizerUrl,
	}
}

func (a authorizer) AuthorizeUser(cpf string) (dto.AuthorizerResponse, error) {
	reqBody := struct {
		CPF string `json:"cpf"`
	}{
		CPF: cpf,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return dto.AuthorizerResponse{}, err
	}

	response, err := a.client.DoPost(a.authorizerUrl, body)
	if err != nil {
		return dto.AuthorizerResponse{}, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return dto.AuthorizerResponse{}, ErrUnauthorized
	}

	var authorizeResponse dto.AuthorizerResponse
	err = json.NewDecoder(response.Body).Decode(&authorizeResponse)
	if err != nil {
		return dto.AuthorizerResponse{}, err
	}

	return authorizeResponse, nil
}
