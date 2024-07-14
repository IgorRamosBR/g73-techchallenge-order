package authorizer

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-order/internal/core/usecases/dto"
	mock_http "github.com/IgorRamosBR/g73-techchallenge-order/internal/infra/drivers/http/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthorizer_AuthorizerUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_http.NewMockHttpClient(ctrl)

	type args struct {
		cpf string
	}
	type want struct {
		response dto.AuthorizerResponse
		err      error
	}
	type httpCall struct {
		times    int
		response *http.Response
		err      error
	}
	tests := []struct {
		name string
		args
		want
		httpCall
	}{
		{
			name: "should fail to authorize user when client returns error",
			args: args{
				cpf: "123456789",
			},
			want: want{
				response: dto.AuthorizerResponse{},
				err:      errors.New("internal server error"),
			},
			httpCall: httpCall{
				times:    1,
				response: &http.Response{},
				err:      errors.New("internal server error"),
			},
		},
		{
			name: "should fail to authorize user when response status code is non-2xx",
			args: args{
				cpf: "123456789",
			},
			want: want{
				response: dto.AuthorizerResponse{},
				err:      ErrUnauthorized,
			},
			httpCall: httpCall{
				times: 1,
				response: &http.Response{
					StatusCode: 500,
				},
				err: nil,
			},
		},
		{
			name: "should fail to authorize user when json decoder returns error",
			args: args{
				cpf: "123456789",
			},
			want: want{
				response: dto.AuthorizerResponse{},
				err:      errors.New("invalid character '<' looking for beginning of value"),
			},
			httpCall: httpCall{
				times: 1,
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("<invalid json>")),
				},
				err: nil,
			},
		},
		{
			name: "should authorize user successfully",
			args: args{
				cpf: "123456789",
			},
			want: want{
				response: dto.AuthorizerResponse{
					IsAuthorized: true,
					Message:      "user is authorized",
					User:         dto.AuthorizedUser{},
				},
				err: nil,
			},
			httpCall: httpCall{
				times: 1,
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"isAuthorized": true, "message": "user is authorized"}`)),
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		httpClient.
			EXPECT().
			DoPost(gomock.Eq("/authorize"), gomock.Any()).
			Return(tt.httpCall.response, tt.httpCall.err)

		authorizer := NewAuthorizer(httpClient, "/authorize")
		response, err := authorizer.AuthorizeUser(tt.args.cpf)

		if err != nil {
			assert.EqualError(t, err, tt.want.err.Error())
		} else {
			assert.Equal(t, tt.want.response, response)
		}
	}
}
