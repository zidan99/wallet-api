package tx_repository

import (
	"database/sql"

	"gorm.io/gorm"
)

type Contract interface {
	// StartTransaction memulai transaksi database baru.
	// Jika logic di dalam StartTransaction menghasilkan error, maka seluruh perubahan akan dibatalkan (rollback).
	// Sedangkan jika tidak ada error, maka seluruh perubahan akan diterapkan (commit).
	StartTransaction(fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error

	// Begin memulai transaksi database baru.
	Begin(opts ...*sql.TxOptions) *gorm.DB

	// Commit menerapkan semua perubahan di dalam transaksi database.
	Commit() *gorm.DB

	// Rollback membatalkan semua perubahan di dalam transaksi database.
	Rollback() *gorm.DB
}

type repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Contract {
	return &repository{DB: db}
}

func (r *repository) StartTransaction(fn func(*gorm.DB) error, opts ...*sql.TxOptions) error {
	return r.DB.Transaction(fn, opts...)
}

func (r *repository) Begin(opts ...*sql.TxOptions) *gorm.DB {
	return r.DB.Begin(opts...)
}

func (r *repository) Commit() *gorm.DB {
	return r.DB.Commit()
}

func (r *repository) Rollback() *gorm.DB {
	return r.DB.Rollback()
}
