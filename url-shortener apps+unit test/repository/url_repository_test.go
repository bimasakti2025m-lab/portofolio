package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"enigmacamp.com/url-shortener/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UrlRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    UrlRepository
}

func (s *UrlRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(s.T(), err)
	s.mockDB = db
	s.mockSQL = mock
	s.repo = NewUrlRepository(s.mockDB)
}

func TestUrlRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UrlRepoTestSuite))
}

func (s *UrlRepoTestSuite) TestCreate_Success() {
	now := time.Now()
	payload := model.Url{
		LongUrl:   "https://www.google.com",
		ShortCode: "abcdefg",
		UserId:    1,
	}
	expectedID := uint32(1)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_urls (long_url, short_code, user_id) VALUES ($1, $2, $3) RETURNING id, created_at")).
		WithArgs(payload.LongUrl, payload.ShortCode, payload.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(expectedID, now))

	createdUrl, err := s.repo.Create(payload)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedID, createdUrl.Id)
	assert.Equal(s.T(), now, createdUrl.CreatedAt)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UrlRepoTestSuite) TestCreate_Fail() {
	payload := model.Url{
		LongUrl:   "https://www.google.com",
		ShortCode: "abcdefg",
		UserId:    1,
	}
	expectedError := errors.New("database error")

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("INSERT INTO mst_urls (long_url, short_code, user_id) VALUES ($1, $2, $3) RETURNING id, created_at")).
		WithArgs(payload.LongUrl, payload.ShortCode, payload.UserId).
		WillReturnError(expectedError)

	createdUrl, err := s.repo.Create(payload)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), uint32(0), createdUrl.Id)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UrlRepoTestSuite) TestGetByShortCode_Success() {
	expectedUrl := model.Url{
		Id:        1,
		LongUrl:   "https://www.enigmacamp.com",
		ShortCode: "xyz123",
		UserId:    1,
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "long_url", "short_code", "user_id", "created_at"}).
		AddRow(expectedUrl.Id, expectedUrl.LongUrl, expectedUrl.ShortCode, expectedUrl.UserId, expectedUrl.CreatedAt)

	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, long_url, short_code, user_id, created_at FROM mst_urls WHERE short_code = $1")).
		WithArgs(expectedUrl.ShortCode).
		WillReturnRows(rows)

	url, err := s.repo.GetByShortCode(expectedUrl.ShortCode)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedUrl, url)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UrlRepoTestSuite) TestGetByShortCode_NotFound() {
	shortCode := "notfound"
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT id, long_url, short_code, user_id, created_at FROM mst_urls WHERE short_code = $1")).
		WithArgs(shortCode).
		WillReturnError(sql.ErrNoRows)

	url, err := s.repo.GetByShortCode(shortCode)

	assert.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, sql.ErrNoRows))
	assert.Equal(s.T(), model.Url{}, url)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UrlRepoTestSuite) TestIsShortCodeExist_True() {
	shortCode := "exists"
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM mst_urls WHERE short_code = $1)")).
		WithArgs(shortCode).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := s.repo.IsShortCodeExist(shortCode)

	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}

func (s *UrlRepoTestSuite) TestIsShortCodeExist_False() {
	shortCode := "notexists"
	s.mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM mst_urls WHERE short_code = $1)")).
		WithArgs(shortCode).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := s.repo.IsShortCodeExist(shortCode)

	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)
	assert.NoError(s.T(), s.mockSQL.ExpectationsWereMet())
}
