package usecases

import (
	"errors"
	"evermos-mini/models"
	"evermos-mini/repositories"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TrxItemInput struct {
	ProductID   uint `json:"product_id"`
	ProdukID    uint `json:"produk_id"`
	LogProdukID uint `json:"log_produk_id"`
	Kuantitas   int  `json:"kuantitas"`
}

type TrxCreateInput struct {
	AlamatPengiriman uint           `json:"alamat_pengiriman"`
	AlamatKirim      uint           `json:"alamat_kirim"`
	MetodeBayar      string         `json:"metode_bayar"`
	MethodBayar      string         `json:"method_bayar"`
	DetailTrx        []TrxItemInput `json:"detail_trx"`
	Items            []TrxItemInput `json:"items"`
}

type TrxUsecase interface {
	CreateTrx(userID uint, input TrxCreateInput) (*models.Transaksi, error)
	GetAllTrx(userID uint) ([]models.Transaksi, error)
	GetTrxByID(id, userID uint) (*models.Transaksi, error)
	UpdateStatusTrx(id uint, status string) (*models.Transaksi, error)
}

type trxUsecase struct {
	trxRepo    repositories.TrxRepository
	alamatRepo repositories.AlamatRepository
}

func NewTrxUsecase(trxRepo repositories.TrxRepository, alamatRepo repositories.AlamatRepository) TrxUsecase {
	return &trxUsecase{
		trxRepo:    trxRepo,
		alamatRepo: alamatRepo,
	}
}

func (u *trxUsecase) CreateTrx(userID uint, input TrxCreateInput) (*models.Transaksi, error) {
	// Normalisasi parameter alamat kirim
	alamatID := input.AlamatPengiriman
	if alamatID == 0 {
		alamatID = input.AlamatKirim
	}

	if alamatID == 0 {
		return nil, errors.New("alamat pengiriman wajib dipilih")
	}

	// Validasi: Alamat diperlukan untuk alamat kirim produk dan harus milik user yang bersangkutan
	_, err := u.alamatRepo.FindByIDAndUserID(alamatID, userID)
	if err != nil {
		return nil, errors.New("alamat pengiriman tidak valid atau bukan milik Anda")
	}

	// Normalisasi metode bayar
	metodeBayar := input.MetodeBayar
	if metodeBayar == "" {
		metodeBayar = input.MethodBayar
	}
	if metodeBayar == "" {
		metodeBayar = "Transfer Bank"
	}

	// Normalisasi item transaksi
	items := input.DetailTrx
	if len(items) == 0 {
		items = input.Items
	}
	if len(items) == 0 {
		return nil, errors.New("daftar produk yang dibeli tidak boleh kosong")
	}

	db := u.trxRepo.GetDB()
	var createdTrx models.Transaksi

	// Eksekusi proses transaksi dalam database transaction untuk menjamin ACID
	err = db.Transaction(func(tx *gorm.DB) error {
		var totalHarga float64
		var details []models.DetailTransaksi

		for _, it := range items {
			if it.Kuantitas <= 0 {
				return errors.New("kuantitas produk harus lebih dari 0")
			}

			// Cari referensi produk ID
			targetProductID := it.ProductID
			if targetProductID == 0 {
				targetProductID = it.ProdukID
			}

			var product models.Produk

			// Jika yang diberikan adalah log_produk_id
			if targetProductID == 0 && it.LogProdukID > 0 {
				var log models.LogProduk
				if err := tx.First(&log, it.LogProdukID).Error; err != nil {
					return fmt.Errorf("data log produk ID %d tidak ditemukan", it.LogProdukID)
				}
				targetProductID = log.ProdukID
			}

			if targetProductID == 0 {
				return errors.New("identitas produk tidak valid")
			}

			// Lock row produk untuk konsistensi stok
			if err := tx.First(&product, targetProductID).Error; err != nil {
				return fmt.Errorf("produk ID %d tidak ditemukan", targetProductID)
			}

			// Validasi ketersediaan stok
			if product.Stok < it.Kuantitas {
				return fmt.Errorf("stok produk '%s' tidak mencukupi (tersisa: %d, diminta: %d)", product.NamaProduk, product.Stok, it.Kuantitas)
			}

			// Pengurangan stok otomatis
			product.Stok -= it.Kuantitas
			if err := tx.Save(&product).Error; err != nil {
				return errors.New("gagal mengupdate stok produk")
			}

			// Ketentuan: Tabel log product diisi ketika melakukan transaksi untuk menyimpan snapshot produk
			snapshotLog := models.LogProduk{
				ProdukID:      product.ID,
				NamaProduk:    product.NamaProduk,
				Slug:          product.Slug,
				HargaReseller: product.HargaReseller,
				HargaKonsumen: product.HargaKonsumen,
				Deskripsi:     product.Deskripsi,
				TokoID:        product.TokoID,
				CategoryID:    product.CategoryID,
			}
			if err := tx.Create(&snapshotLog).Error; err != nil {
				return errors.New("gagal mencatat log snapshot produk transaksi")
			}

			subtotal := product.HargaKonsumen * float64(it.Kuantitas)
			totalHarga += subtotal

			details = append(details, models.DetailTransaksi{
				LogProdukID: snapshotLog.ID,
				Kuantitas:   it.Kuantitas,
				HargaTotal:  subtotal,
			})
		}

		invoiceCode := fmt.Sprintf("INV-%d-%d", userID, time.Now().Unix())

		trx := models.Transaksi{
			UserID:           userID,
			AlamatPengiriman: alamatID,
			MetodeBayar:      metodeBayar,
			HargaTotal:       totalHarga,
			Invoice:          invoiceCode,
			Status:           "Menunggu Pembayaran",
		}

		if err := tx.Create(&trx).Error; err != nil {
			return errors.New("gagal membuat data transaksi")
		}

		for i := range details {
			details[i].TransaksiID = trx.ID
		}

		if err := tx.Create(&details).Error; err != nil {
			return errors.New("gagal menyimpan rincian detail transaksi")
		}

		createdTrx = trx
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Ambil transaksi lengkap beserta relasi snapshot
	fullTrx, err := u.trxRepo.FindByIDAndUserID(createdTrx.ID, userID)
	if err == nil {
		return fullTrx, nil
	}
	return &createdTrx, nil
}

func (u *trxUsecase) GetAllTrx(userID uint) ([]models.Transaksi, error) {
	// Proteksi: User hanya dapat melihat daftar transaksinya sendiri
	return u.trxRepo.FindAllByUserID(userID)
}

func (u *trxUsecase) GetTrxByID(id, userID uint) (*models.Transaksi, error) {
	// Proteksi: User tidak dapat melihat data transaksi dari user lain
	trx, err := u.trxRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan atau Anda tidak memiliki hak akses")
	}
	return trx, nil
}

func (u *trxUsecase) UpdateStatusTrx(id uint, status string) (*models.Transaksi, error) {
	return u.trxRepo.UpdateStatus(id, status)
}
