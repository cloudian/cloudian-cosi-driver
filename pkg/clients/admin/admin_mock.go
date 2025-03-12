package admin

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/cloudian/cosi-driver/pkg/clients/admin/api"
)

// Mock client

type ClientMock struct {
	UserCreatorMock
	UserDeletorMock
}

// Grant bucket access

type UserCreatorMock struct{}

func (UserCreatorMock) PutUser(_ context.Context, body api.PutUserJSONRequestBody, _ ...api.RequestEditorFn) (*http.Response, error) {
	if body.UserId == "grant_bucket_access_create_user_fail" {
		return nil, errors.New("grant_bucket_access_create_user_fail")
	}

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(
			strings.NewReader(`{"canonicalUserId": "testCanonicalUserID"}`),
		),
	}

	return response, nil
}

func (UserCreatorMock) PutUserCredentialsWithResponse(_ context.Context, params *api.PutUserCredentialsParams, _ ...api.RequestEditorFn) (*api.PutUserCredentialsResponse, error) {
	if params.UserId == "grant_bucket_access_create_credentials_fail" {
		return nil, errors.New("grant_bucket_access_create_credentials_fail")
	}

	response := &api.PutUserCredentialsResponse{
		HTTPResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(
				strings.NewReader("mock response"),
			),
		},
		JSON200: &api.SecurityInfo{
			AccessKey: "mockAccessKey",
			SecretKey: "mockSecretKey",
		},
	}

	return response, nil
}

// Revoke bucket access

type UserDeletorMock struct{}

func (UserDeletorMock) GetUser(_ context.Context, params *api.GetUserParams, _ ...api.RequestEditorFn) (*http.Response, error) {
	if params.UserId == "revoke_bucket_access_get_user_fail" {
		return nil, errors.New("revoke_bucket_access_get_user_fail")
	}

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(
			strings.NewReader(`{"canonicalUserId": "testCanonicalUserID"}`),
		),
	}

	return response, nil
}

func (UserDeletorMock) DeleteUser(_ context.Context, params *api.DeleteUserParams, _ ...api.RequestEditorFn) (*http.Response, error) {
	if params.UserId == "revoke_bucket_access_delete_user_fail" {
		return nil, errors.New("revoke_bucket_access_delete_user_fail")
	}

	response := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(
			strings.NewReader("mock response"),
		),
	}

	return response, nil
}
