package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	dbModel "qualification/internal/repository/model"
)

func SelectProducts(ctx context.Context, db *pgxpool.Pool) ([]dbModel.Product, error) {
	rows, err := db.Query(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []dbModel.Product
	for rows.Next() {
		var product dbModel.Product
		if err := rows.Scan(&product.ID, &product.Name); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func InsertProduct(ctx context.Context, db *pgxpool.Pool, id int, name string) error {
	if _, err := db.Exec(ctx, "INSERT INTO products (id, name) VALUES ($1, $2)", id, name); err != nil {
		return err
	}

	return nil
}

func SelectProduct(ctx context.Context, db *pgxpool.Pool, id int) (dbModel.Product, error) {
	var dbProduct dbModel.Product

	if err := db.QueryRow(ctx, "SELECT * FROM products WHERE id = $1", id).Scan(&dbProduct.ID, &dbProduct.Name); err != nil {
		return dbProduct, err
	}

	return dbProduct, nil
}
