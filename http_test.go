package mta

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

var DoFunc func(req *http.Request) (*http.Response, error)

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	return DoFunc(req)
}

type mockReadCloser struct {
	mock.Mock
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}
