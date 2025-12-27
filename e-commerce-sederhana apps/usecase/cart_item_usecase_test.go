package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CartItemUsecaseTestSuite struct {
	suite.Suite
	cartItemRepo *repository.CartItemRepositoryMock
	usecase      CartItemUseCase
}

func (suite *CartItemUsecaseTestSuite) SetupTest() {
	suite.cartItemRepo = new(repository.CartItemRepositoryMock)
	suite.usecase = NewCartItemUseCase(suite.cartItemRepo)
}

func TestCartItemUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CartItemUsecaseTestSuite))
}

func (suite *CartItemUsecaseTestSuite) TestCreateCartItem_Success() {
	cartItem := model.CartItem{ProductID: 1, Quantity: 2}
	createdCartItem := cartItem
	createdCartItem.ID = 1

	suite.cartItemRepo.On("CreateCartItem", &cartItem).Return(&createdCartItem, nil)

	result, err := suite.usecase.CreateCartItem(cartItem)
	suite.NoError(err)
	suite.Equal(&createdCartItem, result)
}

func (suite *CartItemUsecaseTestSuite) TestGetAllCartItems_Success() {
	cartItems := []model.CartItem{{ID: 1, ProductID: 1}}
	suite.cartItemRepo.On("GetAllCartItems").Return(cartItems, nil)

	result, err := suite.usecase.GetAllCartItems()
	suite.NoError(err)
	suite.Equal(cartItems, result)
}

func (suite *CartItemUsecaseTestSuite) TestGetCartItemByID_Success() {
	cartItem := model.CartItem{ID: 1, ProductID: 1}
	suite.cartItemRepo.On("GetCartItemByID", 1).Return(cartItem, nil)

	result, err := suite.usecase.GetCartItemByID(1)
	suite.NoError(err)
	suite.Equal(cartItem, result)
}

func (suite *CartItemUsecaseTestSuite) TestUpdateCartItem_Success() {
	cartItem := model.CartItem{ID: 1, Quantity: 5}
	updatedCartItem := cartItem

	suite.cartItemRepo.On("UpdateCartItem", &cartItem).Return(&updatedCartItem, nil)

	result, err := suite.usecase.UpdateCartItem(cartItem)
	suite.NoError(err)
	suite.Equal(&updatedCartItem, result)
}

func (suite *CartItemUsecaseTestSuite) TestDeleteCartItem_Success() {
	suite.cartItemRepo.On("DeleteCartItem", 1).Return(nil)

	err := suite.usecase.DeleteCartItem(1)
	suite.NoError(err)
}

func (suite *CartItemUsecaseTestSuite) TestDeleteCartItem_Error() {
	suite.cartItemRepo.On("DeleteCartItem", 1).Return(errors.New("failed"))

	err := suite.usecase.DeleteCartItem(1)
	suite.Error(err)
}
