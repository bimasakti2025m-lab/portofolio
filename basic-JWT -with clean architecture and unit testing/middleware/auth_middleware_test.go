package middleware

import (
	"basic-JWT/mock/service_mock"
	"basic-JWT/model"
	modelutils "basic-JWT/utils/model_utils"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareSuite struct {
	suite.Suite
	authMiddleware *AuthMiddleware
	jwtService     *service_mock.JWTServiceMock // Use the mock service

}

func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareSuite))
}

func (a *AuthMiddlewareSuite) SetupTest() {
	a.jwtService = new(service_mock.JWTServiceMock)
	a.authMiddleware = NewAuthMiddleware(a.jwtService)
}

func (a *AuthMiddlewareSuite) TestRequireToken_Success() {
	// Create a dummy Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer valid_token")

	// Mock the JWT service to return a valid token
	a.jwtService.On("VerifyToken", "valid_token").Return(&modelutils.JwtPayloadClaims{Role: "user"}, nil).Once()

	// Call the middleware
	handler := a.authMiddleware.RequireToken("user")
	handler(c)

	// Assert that the request was not aborted and no error was returned
	a.False(c.IsAborted())
	a.Equal(http.StatusOK, w.Code) // Assuming Next() was called, so no response written by middleware
}

func (a *AuthMiddlewareSuite) TestRequireToken_NoToken() {
	// Create a dummy Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	// Call the middleware
	handler := a.authMiddleware.RequireToken("user")
	handler(c)

	// Assert that the request was aborted and an unauthorized response was returned
	a.True(c.IsAborted())
	a.Equal(http.StatusUnauthorized, w.Code)
	a.Contains(w.Body.String(), "unauthorized")
}

func (a *AuthMiddlewareSuite) TestRequireToken_InvalidToken() {
	// Create a dummy Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid_token")

	// Mock the JWT service to return an error for an invalid token
	a.jwtService.On("VerifyToken", "invalid_token").Return(&modelutils.JwtPayloadClaims{}, errors.New("invalid token")).Once()

	// Call the middleware
	handler := a.authMiddleware.RequireToken("user")
	handler(c)

	// Assert that the request was aborted and an unauthorized response was returned
	a.True(c.IsAborted())
	a.Equal(http.StatusUnauthorized, w.Code)
	a.Contains(w.Body.String(), "unauthorized")
}

func (a *AuthMiddlewareSuite) TestRequireToken_ForbiddenRole() {
	// Create a dummy Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer valid_token")

	// Mock the JWT service to return a token with a 'user' role
	a.jwtService.On("VerifyToken", "valid_token").Return(&modelutils.JwtPayloadClaims{Role: "user"}, nil).Once()

	// Call the middleware, requiring an 'admin' role
	handler := a.authMiddleware.RequireToken("admin")
	handler(c)

	// Assert that the request was aborted and a forbidden response was returned
	a.True(c.IsAborted())
	a.Equal(http.StatusForbidden, w.Code)
	a.Contains(w.Body.String(), "forbidden")
}

func (a *AuthMiddlewareSuite) TestRequireToken_MultipleRoles_Success() {
	// Create a dummy Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer valid_token")

	// Mock the JWT service to return a token with a 'user' role
	a.jwtService.On("VerifyToken", "valid_token").Return(&modelutils.JwtPayloadClaims{Role: "user"}, nil).Once()

	// Call the middleware, requiring either 'admin' or 'user' role
	handler := a.authMiddleware.RequireToken("admin", "user")
	handler(c)

	// Assert that the request was not aborted
	a.False(c.IsAborted())
	a.Equal(http.StatusOK, w.Code)
}

func (a *AuthMiddlewareSuite) TestCreateToken_Success() {
	user := model.User{ID: 1, Username: "username", Password: "password", Role: "user"}
	a.jwtService.On("CreateToken", user).Return("test_token").Once()

	token := a.jwtService.CreateToken(user)
	a.Equal("test_token", token)
}

func (a *AuthMiddlewareSuite) TestVerifyToken_Success() {
	tokenString := "valid_token"
	claims := &modelutils.JwtPayloadClaims{UserId: 1, Role: "admin"}
	a.jwtService.On("VerifyToken", tokenString).Return(claims, nil).Once()

	resultClaims, err := a.jwtService.VerifyToken(tokenString)
	// Assert
	a.NoError(err)
	a.Equal(claims, resultClaims)

}




	
