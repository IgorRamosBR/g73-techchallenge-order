package dto

type AuthorizerResponse struct {
	IsAuthorized bool           `json:"isAuthorized"`
	Message      string         `json:"message"`
	User         AuthorizedUser `json:"user"`
}

type AuthorizedUser struct {
	CPF   string `json:"cpf"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
