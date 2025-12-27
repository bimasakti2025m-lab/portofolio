package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CartUsecaseTestSuite struct {
	suite.Suite
	cartRepo *repository.CartRepositoryMock
	usecase  CartUseCase
}

func (suite *CartUsecaseTestSuite) SetupTest() {
	suite.cartRepo = new(repository.CartRepositoryMock)
	suite.usecase = NewCartUseCase(suite.cartRepo)
}

func TestCartUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CartUsecaseTestSuite))
}

func (suite *CartUsecaseTestSuite) TestCreateCart_Success() {
	cart := model.Cart{UserID: 1}
	createdCart := cart
	createdCart.ID = 1

	suite.cartRepo.On("CreateCart", &cart).Return(&createdCart, nil)

	result, err := suite.usecase.CreateCart(cart)
	suite.NoError(err)
	suite.Equal(&createdCart, result)
}

func (suite *CartUsecaseTestSuite) TestGetAllCarts_Success() {
	carts := []model.Cart{{ID: 1, UserID: 1}}
	suite.cartRepo.On("GetAllCarts").Return(carts, nil)

	result, err := suite.usecase.GetAllCarts()
	suite.NoError(err)
	suite.Equal(carts, result)
}

func (suite *CartUsecaseTestSuite) TestGetCartByID_Success() {
	cart := model.Cart{ID: 1, UserID: 1}
	suite.cartRepo.On("GetCartByID", 1).Return(cart, nil)

	result, err := suite.usecase.GetCartByID(1)
	suite.NoError(err)
	suite.Equal(cart, result)
}

func (suite *CartUsecaseTestSuite) TestUpdateCart_Success() {
	cart := model.Cart{ID: 1, UserID: 2}
	updatedCart := cart

	suite.cartRepo.On("UpdateCart", &cart).Return(&updatedCart, nil)

	result, err := suite.usecase.UpdateCart(cart)
	suite.NoError(err)
	suite.Equal(&updatedCart, result)
}

func (suite *CartUsecaseTestSuite) TestDeleteCart_Success() {
	suite.cartRepo.On("DeleteCart", 1).Return(nil)

	err := suite.usecase.DeleteCart(1)
	suite.NoError(err)
}

func (suite *CartUsecaseTestSuite) TestDeleteCart_Error() {
	suite.cartRepo.On("DeleteCart", 1).Return(errors.New("failed"))

	err := suite.usecase.DeleteCart(1)
	suite.Error(err)
}
