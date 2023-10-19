package usecase

import (
	"context"

	"qualification/internal/repository"
	blModel "qualification/internal/usecase/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetProducts(ctx context.Context, db *pgxpool.Pool) ([]blModel.Product, error) {
	dbProducts, err := repository.SelectProducts(ctx, db)
	if err != nil {
		return nil, err
	}

	var ucProducts []blModel.Product
	for _, dbProduct := range dbProducts {
		ucProducts = append(ucProducts, blModel.Product{ID: dbProduct.ID, Name: dbProduct.Name})
	}

	return ucProducts, nil
}

func CreateProduct(ctx context.Context, db *pgxpool.Pool, blProduct blModel.Product) error {
	return repository.InsertProduct(ctx, db, blProduct.ID, blProduct.Name)
}

func GetProduct(ctx context.Context, db *pgxpool.Pool, id int) (*blModel.Product, error) {
	dbProduct, err := repository.SelectProduct(ctx, db, id)
	if err != nil {
		return nil, err
	}

	return &blModel.Product{ID: dbProduct.ID, Name: dbProduct.Name}, nil
}
