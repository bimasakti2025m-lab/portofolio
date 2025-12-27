package usecase

import (
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
	"E-commerce-Sederhana/utils/service/midtrans"
	"time"
)

type OrderUsecase interface {
	GetAllOrders() ([]model.Order, error)
	GetOrderById(id int) (model.Order, error)
	CreateOrder(order model.Order) (model.Order, error)
	UpdateOrder(order model.Order) (model.Order, error)
	DeleteOrder(id int) error
}

type orderUsecase struct {
	repo            repository.OrderRepository
	midtransService midtrans.MidtransService
}

func (ou *orderUsecase) GetAllOrders() ([]model.Order, error) {
	return ou.repo.GetAllOrders()
}

func (ou *orderUsecase) GetOrderById(id int) (model.Order, error) {
	return ou.repo.GetOrderById(id)
}

func (ou *orderUsecase) CreateOrder(order model.Order) (model.Order, error) {
	// Tambahkan status awal dan ID unik untuk Midtrans
	order.StatusPesanan = "pending"
	order.TransactionIDMidtrans = "ORDER-" + time.Now().Format("20060102150405")

	// Simpan order ke database terlebih dahulu
	createdOrder, err := ou.repo.CreateOrder(order)
	if err != nil {
		return model.Order{}, err
	}

	// Buat transaksi di Midtrans
	snapResp, err := ou.midtransService.CreateTransaction(createdOrder)
	if err != nil {
		// Disini Anda bisa memutuskan, apakah mau menghapus order yang sudah dibuat
		// atau membiarkannya dengan status pending untuk dicoba lagi nanti.
		return model.Order{}, err
	}

	// Kembalikan URL pembayaran ke controller
	createdOrder.TransactionIDMidtrans = snapResp.RedirectURL
	return createdOrder, nil
}

func (ou *orderUsecase) UpdateOrder(order model.Order) (model.Order, error) {
	return ou.repo.UpdateOrder(order)
}

func (ou *orderUsecase) DeleteOrder(id int) error {
	return ou.repo.DeleteOrder(id)
}

func NewOrderUsecase(repo repository.OrderRepository, midtransService midtrans.MidtransService) OrderUsecase {
	return &orderUsecase{
		repo:            repo,
		midtransService: midtransService,
	}
}
