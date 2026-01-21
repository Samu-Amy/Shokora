package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		User:    &MockUserStore{},
		Product: &MockProductStore{},
	}
}

// ----- USERS -----

type MockUserStore struct {
}

func (m *MockUserStore) Create(context.Context, *sql.Tx, *User) error {
	return nil
}

func (m *MockUserStore) GetById(context.Context, int64) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) GetByEmail(context.Context, string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) Delete(context.Context, *sql.Tx, int64) error {
	return nil
}

func (m *MockUserStore) CreateUserAndSendVerification(context.Context, *User, string, time.Duration) error {
	return nil
}

func (m *MockUserStore) VerifyEmail(context.Context, string) error {
	return nil
}

func (m *MockUserStore) DeleteUserAndEmailVerificationToken(context.Context, int64) error {
	return nil
}

func (m *MockUserStore) createEmailVerification(context.Context, *sql.Tx, string, time.Duration, int64) error {
	return nil
}

func (m *MockUserStore) getUserFromEmailVerificationToken(context.Context, *sql.Tx, string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) setUserIsVerified(context.Context, *sql.Tx, int64) error {
	return nil
}

func (m *MockUserStore) deleteEmailVerificationToken(context.Context, *sql.Tx, int64) error {
	return nil
}

// ----- PRODUCTS -----

type MockProductStore struct {
}

func (m *MockProductStore) Create(context.Context, *Product) error {
	return nil
}
func (m *MockProductStore) GetById(context.Context, int64) (*Product, error) {
	return &Product{}, nil
}
func (m *MockProductStore) GetProducts(context.Context, QueryPaginationOptions, ProductsFilters) ([]Product, error) {
	return make([]Product, 1), nil
}
func (m *MockProductStore) GetMenuProducts(context.Context, QueryPaginationOptions, MenuFilters) ([]Product, error) {
	return make([]Product, 1), nil
}
func (m *MockProductStore) Update(context.Context, *Product) error {
	return nil
}
func (m *MockProductStore) Delete(context.Context, int64) error {
	return nil
}
