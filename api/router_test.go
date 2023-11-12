package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"phone-directory/api"
	"phone-directory/internal/entities"
	"phone-directory/internal/service"
	serviceMock "phone-directory/internal/service/mock"
)

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))

}

type APITestSuite struct {
	suite.Suite

	ctrl *gomock.Controller

	us service.UserService
	ps service.PhoneService
	as service.AddressService
}

func (ts *APITestSuite) SetupTest() {
	ts.ctrl = gomock.NewController(ts.T())

	ts.us = serviceMock.NewMockUserService(ts.ctrl)
	ts.ps = serviceMock.NewMockPhoneService(ts.ctrl)
	ts.as = serviceMock.NewMockAddressService(ts.ctrl)
}

func (ts *APITestSuite) TearDownTest() {
	ts.ctrl.Finish()
}

func (ts *APITestSuite) TestGetUser() {

	h := api.NewHandler(ts.us, ts.ps, ts.as)
	router, err := api.Router(h)
	ts.NoError(err)

	ts.us.(*serviceMock.MockUserService).EXPECT().Get(gomock.Any(), uint(1)).Return(&entities.User{
		ID:   1,
		Name: "test",
	}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/users/1", nil)
	router.ServeHTTP(w, req)

	ts.Equal(200, w.Code)
	ts.Equal("{\"ID\":1,\"Name\":\"test\",\"Phones\":null,\"Addresses\":null}", w.Body.String())

}
