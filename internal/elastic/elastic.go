package elastic

import (
	"context"
	"encoding/json"
	"fmt"

	"qualification/internal/usecase/model"

	"github.com/olivere/elastic/v7"
)

type es struct {
	client *elastic.Client
}

type ES interface {
	CreateProductIndex(ctx context.Context, product model.Product) error
	SearchDoc(ctx context.Context, productName string) []model.Product
}

func NewEsClient(esClient *elastic.Client) ES {
	return &es{client: esClient}
}

func (es *es) CreateProductIndex(ctx context.Context, product model.Product) error {
	dataJSON, err := json.Marshal(product)

	if _, err = es.client.Index().Index("products").BodyJson(string(dataJSON)).Do(ctx); err != nil {
		return err
	}

	return nil
}

func (es *es) SearchDoc(ctx context.Context, productName string) []model.Product {
	var products []model.Product

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchPhrasePrefixQuery("name", productName))

	queryStr, err := searchSource.Source()
	if err != nil {
		fmt.Println(fmt.Errorf("err during query marshal: %w", err))
		return nil
	}

	queryJs, err := json.Marshal(queryStr)
	if err != nil {
		fmt.Println(fmt.Errorf("err during query marshal: %w", err))
		return nil
	}

	fmt.Println("Final ESQuery=\n", string(queryJs))

	es.client.Search().Index("products").AllowPartialSearchResults(true).SearchSource(searchSource)
	searchService := es.client.Search().Index("products").SearchSource(searchSource)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("search err: %w", err))
		return nil
	}

	for _, hit := range searchResult.Hits.Hits {
		var product model.Product
		if err = json.Unmarshal(hit.Source, &product); err != nil {
			fmt.Println(fmt.Errorf("unmarshal err: %w", err))
			continue
		}

		products = append(products, product)
	}

	if err != nil {
		fmt.Println(fmt.Errorf("fetching err: %w", err))
		return nil
	}

	return products
}
