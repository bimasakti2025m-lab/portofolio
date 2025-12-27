package midtrans

import (
	"E-commerce-Sederhana/config"
	"E-commerce-Sederhana/model"
	"E-commerce-Sederhana/repository"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService interface {
	CreateTransaction(order model.Order) (*snap.Response, error)
	HandleNotification(notificationPayload map[string]interface{}) error
}

type midtransService struct {
	snapClient      snap.Client
	orderRepository repository.OrderRepository // Tambahkan dependensi ke order repository
	serverKey       string                     // Tambahkan server key untuk verifikasi
}

func (m *midtransService) CreateTransaction(order model.Order) (*snap.Response, error) {

	// Buat request untuk Midtrans Snap
	req := &snap.Request{

		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(order.ID) + "-" + order.TransactionIDMidtrans, // Order ID unik
			GrossAmt: int64(order.Total),
		},
		// Anda bisa menambahkan detail pelanggan jika ada
		// CustomerDetails: &midtrans.CustomerDetails{
		// 	FName: "John",
		// 	LName: "Doe",
		// 	Email: "john.doe@example.com",
		// },
	}

	// Buat transaksi Snap
	snapResp, err := m.snapClient.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return snapResp, nil
}

func (m *midtransService) HandleNotification(notificationPayload map[string]interface{}) error {
	// 1. Parse payload dari notifikasi
	orderID, exists := notificationPayload["order_id"].(string)
	if !exists || orderID == "" {
		return nil // Notifikasi tidak valid atau order_id kosong
	}

	// Midtrans orderID bisa berupa "ID_ORDER-TRANSACTION_ID"
	// Kita perlu memisahkan ID_ORDER yang merupakan ID dari tabel order kita
	parts := strings.Split(orderID, "-")
	orderIDInt, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	transactionStatus, exists := notificationPayload["transaction_status"].(string)
	if !exists {
		return nil
	}

	// 2. Lakukan verifikasi signature key (WAJIB di production)
	err = m.verifySignature(notificationPayload)
	if err != nil {
		// Jika signature tidak valid, abaikan notifikasi
		return fmt.Errorf("signature key verification failed: %w", err)
	}

	// 3. Dapatkan order dari database berdasarkan orderID
	// Di sini kita asumsikan orderID dari Midtrans adalah ID dari tabel order Anda
	// Jika formatnya berbeda (misal: "order-123"), Anda perlu parsing dulu
	order, err := m.orderRepository.GetOrderById(orderIDInt)
	if err != nil || order.ID == 0 { // Periksa juga jika order tidak ditemukan (ID=0)
		return fmt.Errorf("order with ID %d not found: %w", orderIDInt, err)
	}

	// 4. Update status order berdasarkan status transaksi dari Midtrans
	if transactionStatus == "settlement" {
		// Transaksi berhasil dan dana sudah masuk
		order.StatusPesanan = "PAID" // Pastikan string ini sesuai dengan enum/status di sistem Anda
	} else if transactionStatus == "expire" {
		// Transaksi kadaluarsa
		order.StatusPesanan = "EXPIRED"
	} else if transactionStatus == "cancel" {
		// Transaksi dibatalkan oleh customer
		order.StatusPesanan = "CANCELLED"
	} else {
		// Untuk status lain (pending, deny, dll), kita tidak melakukan apa-apa
		return nil
	}

	// 5. Simpan perubahan status ke database
	_, err = m.orderRepository.UpdateOrder(order)
	return err
}

func (m *midtransService) verifySignature(payload map[string]interface{}) error {
	orderID, _ := payload["order_id"].(string)
	statusCode, _ := payload["status_code"].(string)
	grossAmount, _ := payload["gross_amount"].(string)
	signatureKey, _ := payload["signature_key"].(string)

	// Pastikan grossAmount memiliki format "10000.00" jika tidak ada
	if !strings.HasSuffix(grossAmount, ".00") {
		grossAmount = grossAmount + ".00"
	}

	hash := sha512.New()
	hash.Write([]byte(orderID + statusCode + grossAmount + m.serverKey))
	expectedSignature := hex.EncodeToString(hash.Sum(nil))

	if signatureKey != expectedSignature {
		// Tambahkan log untuk debugging agar lebih mudah melacak masalah
		log.Printf("Signature Mismatch. Expected: %s, Got: %s", expectedSignature, signatureKey)
		return fmt.Errorf("invalid signature key")
	}
	return nil
}

func NewMidtransService(cfg config.MidtransConfig, orderRepo repository.OrderRepository) MidtransService {
	var env midtrans.EnvironmentType
	if cfg.Env == "production" {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}

	var snapClient snap.Client
	snapClient.New(cfg.ServerKey, env)

	return &midtransService{
		snapClient:      snapClient,
		orderRepository: orderRepo, // Inisialisasi repository
		serverKey:       cfg.ServerKey,
	}
}
