package usecase

import (
	"context"
	"encoding/json"

	"qualification/internal/jaeger"
	"qualification/internal/repository"
	blModel "qualification/internal/usecase/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetProducts(ctx context.Context, db *pgxpool.Pool) ([]blModel.Product, error) {
	ctx, span := jaeger.NewSpan(ctx, "UC GetProducts")
	defer span.End()

	dbProducts, err := repository.SelectProducts(ctx, db)
	if err != nil {
		jaeger.RecordError(span, err)

		return nil, err
	}

	var ucProducts []blModel.Product
	for _, dbProduct := range dbProducts {
		ucProducts = append(ucProducts, blModel.Product{ID: dbProduct.ID, Name: dbProduct.Name})
	}

	jsonProducts, err := json.Marshal(ucProducts)
	if err != nil {
		jaeger.RecordError(span, err)
	}

	jaeger.AddAttributes(span, "uc result", string(jsonProducts))

	return ucProducts, nil
}

func CreateProduct(ctx context.Context, db *pgxpool.Pool, blProduct blModel.Product) error {
	ctx, span := jaeger.NewSpan(ctx, "UC CreateProduct")
	defer span.End()

	return repository.InsertProduct(ctx, db, blProduct.ID, blProduct.Name)
}

func GetProduct(ctx context.Context, db *pgxpool.Pool, id int) (*blModel.Product, error) {
	ctx, span := jaeger.NewSpan(ctx, "UC GetProduct")
	defer span.End()

	dbProduct, err := repository.SelectProduct(ctx, db, id)
	if err != nil {
		jaeger.RecordError(span, err)
		return nil, err
	}

	return &blModel.Product{ID: dbProduct.ID, Name: dbProduct.Name}, nil
}
