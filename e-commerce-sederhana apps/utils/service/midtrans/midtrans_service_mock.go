package midtrans

import (
	"E-commerce-Sederhana/model"

	"github.com/midtrans/midtrans-go/snap"
	"github.com/stretchr/testify/mock"
)

type MidtransServiceMock struct {
	mock.Mock
}

func (m *MidtransServiceMock) CreateTransaction(order model.Order) (*snap.Response, error) {
	args := m.Called(order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*snap.Response), args.Error(1)
}

func (m *MidtransServiceMock) HandleNotification(notificationPayload map[string]interface{}) error {
	args := m.Called(notificationPayload)
	return args.Error(0)
}

func (m *MidtransServiceMock) verifySignatureKey(orderID, statusCode, grossAmount, signatureKey string) bool {
	args := m.Called(orderID, statusCode, grossAmount, signatureKey)
	return args.Bool(0)
}