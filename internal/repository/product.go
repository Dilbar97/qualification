package repository

import (
	"context"
	"encoding/json"

	"qualification/internal/jaeger"

	dbModel "qualification/internal/repository/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

func SelectProducts(ctx context.Context, db *pgxpool.Pool) ([]dbModel.Product, error) {
	ctx, span := jaeger.NewSpan(ctx, "DB GetProducts")
	defer span.End()

	jaeger.AddSqlEvent(span, "SELECT * FROM products", nil)

	rows, err := db.Query(ctx, "SELECT * FROM products")
	if err != nil {
		jaeger.RecordError(span, err)

		return nil, err
	}
	defer rows.Close()

	var products []dbModel.Product
	for rows.Next() {
		var product dbModel.Product
		if err := rows.Scan(&product.ID, &product.Name); err != nil {
			jaeger.RecordError(span, err)

			return nil, err
		}

		products = append(products, product)
	}

	jsonProducts, _ := json.Marshal(products)
	jaeger.AddAttributes(span, "sql query result", string(jsonProducts))

	return products, nil
}

func InsertProduct(ctx context.Context, db *pgxpool.Pool, id int, name string) error {
	ctx, span := jaeger.NewSpan(ctx, "DB CreateProduct")
	defer span.End()

	jaeger.AddSqlEvent(span, "INSERT INTO products (id, name) VALUES ($1, $2)", id, name)

	if _, err := db.Exec(ctx, "INSERT INTO products (id, name) VALUES ($1, $2)", id, name); err != nil {
		jaeger.RecordError(span, err)
		return err
	}

	return nil
}

func SelectProduct(ctx context.Context, db *pgxpool.Pool, id int) (dbModel.Product, error) {
	ctx, span := jaeger.NewSpan(ctx, "DB SelectProduct")
	defer span.End()

	jaeger.AddSqlEvent(span, "SELECT * FROM products WHERE id = $1", id)

	var dbProduct dbModel.Product
	if err := db.QueryRow(ctx, "SELECT * FROM products WHERE id = $1", id).Scan(&dbProduct.ID, &dbProduct.Name); err != nil {
		return dbProduct, err
	}

	jsonProduct, _ := json.Marshal(dbProduct)
	jaeger.AddAttributes(span, "sql query result", string(jsonProduct))

	return dbProduct, nil
}
