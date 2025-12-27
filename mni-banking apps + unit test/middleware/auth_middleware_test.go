// buatkan mocking buat auth middleware
package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"enigmacamp.com/mini-banking/middleware"
	modelutils "enigmacamp.com/mini-banking/utils/model_utils"
	"enigmacamp.com/mini-banking/utils/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	jwtServiceMock *service.JwtServiceMock
	authMiddleware middleware.AuthMiddleware
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	s.jwtServiceMock = new(service.JwtServiceMock)
	s.authMiddleware = middleware.NewAuthMiddleware(s.jwtServiceMock)
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

func (s *AuthMiddlewareTestSuite) TestRequireToken_Success_WithRole() {
	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	mockClaims := &modelutils.JwtPayloadClaims{
		UserId: 1,
		Role:   "admin",
	}
	s.jwtServiceMock.On("VerifyToken", "valid-token").Return(mockClaims, nil)

	// Eksekusi middleware
	handler := s.authMiddleware.RequireToken("admin")
	handler(c)

	// Assertions
	assert.False(s.T(), c.IsAborted(), "Context should not be aborted")
	claims, exists := c.Get("claims")
	assert.True(s.T(), exists, "Claims should exist in context")
	assert.Equal(s.T(), mockClaims, claims)
	s.jwtServiceMock.AssertExpectations(s.T())
}

func (s *AuthMiddlewareTestSuite) TestRequireToken_Success_NoRoleRequired() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer valid-token")

	mockClaims := &modelutils.JwtPayloadClaims{Role: "user"}
	s.jwtServiceMock.On("VerifyToken", "valid-token").Return(mockClaims, nil)

	handler := s.authMiddleware.RequireToken() // Tidak ada role yang dibutuhkan
	handler(c)

	assert.False(s.T(), c.IsAborted())
	s.jwtServiceMock.AssertExpectations(s.T())
}

func (s *AuthMiddlewareTestSuite) TestRequireToken_Fail_NoAuthHeader() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil) // Tanpa header Authorization

	handler := s.authMiddleware.RequireToken("admin")
	handler(c)

	assert.True(s.T(), c.IsAborted(), "Context should be aborted")
	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)

}

func (s *AuthMiddlewareTestSuite) TestRequireToken_Fail_InvalidTokenFormat() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Invalid-token") // Format salah

	handler := s.authMiddleware.RequireToken("admin")
	handler(c)

	assert.True(s.T(), c.IsAborted())
	assert.Equal(s.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(s.T(), w.Body.String(), "invalid token format")
}

func (s *AuthMiddlewareTestSuite) TestRequireToken_Fail_ForbiddenRole() {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer user-token")

	mockClaims := &modelutils.JwtPayloadClaims{Role: "user"}
	s.jwtServiceMock.On("VerifyToken", "user-token").Return(mockClaims, nil)

	handler := s.authMiddleware.RequireToken("admin") // Membutuhkan 'admin', tapi tokennya 'user'
	handler(c)

	assert.True(s.T(), c.IsAborted())
	assert.Equal(s.T(), http.StatusForbidden, w.Code)
	assert.Contains(s.T(), w.Body.String(), "you don't have the right role")
	s.jwtServiceMock.AssertExpectations(s.T())
}
