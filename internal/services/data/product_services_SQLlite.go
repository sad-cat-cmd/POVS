package data

import (
	"database/sql"
	"errors"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sad-cat-cmd/WebApi/internal/config"
	"github.com/sad-cat-cmd/WebApi/internal/models"
)

const (
	createTableSQL = `
        CREATE TABLE IF NOT EXISTS products (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            definition TEXT,
            price REAL NOT NULL,
            image TEXT,
            created_at DATETIME NOT NULL,
            updated_at DATETIME NOT NULL
        );
    `

	createIndexSQL = `
        CREATE INDEX IF NOT EXISTS idx_products_id ON products(id);
    `

	selectAllSQL  = `SELECT id, name, definition, price, image, created_at, updated_at FROM products;`
	selectByIDSQL = `SELECT id, name, definition, price, image, created_at, updated_at FROM products WHERE id = ?;`
	insertSQL     = `INSERT INTO products (id, name, definition, price, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?);`
	updateSQL     = `UPDATE products SET name = ?, definition = ?, price = ?, image = ?, updated_at = ? WHERE id = ?;`
	deleteSQL     = `DELETE FROM products WHERE id = ?;`
)

type ProductServiceSQLlite struct {
	cfg   *config.Configuration
	db    *sql.DB
	mutex sync.RWMutex
}

func NewProductServiceSQLlite(c *config.Configuration) (*ProductServiceSQLlite, error) {
	db, errOpen := sql.Open("sqlite3", c.DataBaseFilePath)
	if errOpen != nil {
		return nil, errors.New("failed to open database:" + errOpen.Error())
	}
	errPing := db.Ping()
	if errPing != nil {
		return nil, errors.New("failed to ping database:" + errPing.Error())
	}
	service := &ProductServiceSQLlite{
		db: db,
	}
	errInitDB := service.initDBTable()
	if errInitDB != nil {
		return nil, errors.New("failed to init database:" + errInitDB.Error())
	}
	return service, nil
}

func (s *ProductServiceSQLlite) initDBTable() error {

	_, errCreateTable := s.db.Exec(createTableSQL)
	if errCreateTable != nil {
		return errors.New("Failed to create table:" + errCreateTable.Error())
	}
	_, errCreateIndex := s.db.Exec(createIndexSQL)
	if errCreateIndex != nil {
		return errors.New("Failed to create index:" + errCreateIndex.Error())
	}
	return nil
}
func (s *ProductServiceSQLlite) CloseDB() error {
	return s.db.Close()
}
func (s *ProductServiceSQLlite) getById(id string) (*models.Product, error) {
	row := s.db.QueryRow(selectByIDSQL, id)
	var product models.Product
	errScan := row.Scan(&product.ID,
		&product.Name,
		&product.Definition,
		&product.Price,
		&product.Image,
		&product.CreateAt,
		&product.UpdatedAt)
	if errScan == sql.ErrNoRows {
		return nil, nil
	}
	if errScan != nil {
		return nil, errors.New("failed to scan product:" + errScan.Error())
	}
	return &product, nil
}
func (s *ProductServiceSQLlite) Add(product *models.Product) (*models.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, errInsert := s.db.Exec(insertSQL,
		product.ID,
		product.Name,
		product.Definition,
		product.Price,
		product.Image,
		product.CreateAt,
		product.UpdatedAt)
	if errInsert != nil {
		return nil, errors.New("failed to insert product:" + errInsert.Error())
	}

	return product, nil
}
func (s *ProductServiceSQLlite) Remove(id string) (*models.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	product, err := s.getById(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("Object doesn't exist")
	}

	result, err := s.db.Exec(deleteSQL, id)
	if err != nil {
		return nil, errors.New("failed to delete product:" + err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, errors.New("failed to get rows affected:" + err.Error())
	}
	if rows == 0 {
		return nil, errors.New("Object doesn't exist")
	}

	return product, nil
}
func (s *ProductServiceSQLlite) Edit(new_product *models.Product) (*models.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	edited, err := s.db.Exec(updateSQL,
		new_product.Name,
		new_product.Definition,
		new_product.Price,
		new_product.Image,
		time.Now(),
		new_product.ID,
	)
	if err != nil {
		return nil, errors.New("failed to update product:" + err.Error())
	}
	rows, err := edited.RowsAffected()
	if rows == 0 {
		return nil, errors.New("object doesn't exist")
	}
	if err != nil {
		return nil, errors.New("failed to rowAffected")
	}
	return new_product, nil
}
func (s *ProductServiceSQLlite) Search(id string) (*models.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	searched, err := s.getById(id)
	if searched == nil && err == nil {
		return nil, errors.New("object doesn't exist")
	}
	if err != nil {
		return nil, errors.New("failed to rowAffected")
	}
	return searched, nil
}
func (s *ProductServiceSQLlite) GetAll() ([]*models.Product, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	rows, err := s.db.Query(selectAllSQL)
	if err != nil {
		return nil, errors.New("failed to query products:" + err.Error())
	}
	defer rows.Close()
	var products []*models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID,
			&p.Name,
			&p.Definition,
			&p.Price,
			&p.Image,
			&p.CreateAt,
			&p.UpdatedAt)
		if err != nil {
			return nil, errors.New("failed to scan product:" + err.Error())
		}
		products = append(products, &p)
	}

	return products, nil
}
