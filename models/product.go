package models

import (
	"fmt"
	"strconv"

	"github.com/bandros/framework"
)

func GetAllVariantProduct() ([]map[string]interface{}, error) {
	db := framework.Database{}
	defer db.Close()

	db.Select("v.id, m.nama brand, p.sku, v.tipe tipe, v.warna").From("variant v")
	db.Join("product p", "p.id=v.id_product", "")
	db.Join("merk m", "m.id=p.id_merk", "")
	return db.Result()
}

func BuatTransaksi(idvariant, qty, alamat, nama_pembeli string) (map[string]interface{}, error) {
	nqty, _ := strconv.Atoi(qty)
	db := framework.Database{}
	defer db.Close()
	db.Transaction()

	db.Select("*").From("variant").Where("id", idvariant)
	dataProduct, err := db.Row()
	if err != nil {
		db.Rollback()
		return nil, err
	}
	// buat transaksi
	db.From("transaksi")
	dataTransaksi := map[string]interface{}{
		"status":            0,
		"alamat_pengiriman": alamat,
		"nama_pembeli":      nama_pembeli,
	}
	transaksi, err := db.Insert(dataTransaksi)
	if err != nil {
		db.Rollback()
		return nil, err
	}

	harga, _ := strconv.Atoi(dataProduct["harga"].(string))
	detail_transaksi_data := map[string]interface{}{
		"id_transaksi": transaksi,
		"id_variant":   idvariant,
		"qty":          nqty,
		"harga_total":  nqty * harga,
	}
	fmt.Println(detail_transaksi_data)
	db.From("detail_transaksi")
	_, err = db.Insert(detail_transaksi_data)
	if err != nil {
		db.Rollback()
		return nil, err
	}

	// update variant
	stock, _ := strconv.Atoi(dataProduct["stock"].(string))
	updateVariantData := map[string]interface{}{
		"stock": stock - nqty,
	}
	db.From("variant")
	err = db.Update(updateVariantData)
	if err != nil {
		db.Rollback()
		return nil, err
	}

	err = db.Commit()
	if err != nil {
		db.Rollback()
		return nil, err
	}

	return map[string]interface{}{"id_transaksi": transaksi}, err
}

func UpdateTransaksi(id_transaksi, status string) (bool, error) {
	db := framework.Database{}
	defer db.Close()

	update := map[string]interface{}{
		"status": status,
	}
	db.From("transaksi").Where("id", id_transaksi)
	err := db.Update(update)
	if err != nil {
		return false, err
	}
	return true, nil
}
