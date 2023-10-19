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
		fmt.Println(fmt.Errorf("[esclient][GetResponse]err during query marshal=%w", err))
		return nil
	}

	queryJs, err := json.Marshal(queryStr)
	if err != nil {
		fmt.Println(fmt.Errorf("[esclient][GetResponse]err during query marshal=%w", err))
		return nil
	}

	fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))

	es.client.Search().Index("products").AllowPartialSearchResults(true).SearchSource(searchSource)
	searchService := es.client.Search().Index("products").SearchSource(searchSource)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return nil
	}

	for _, hit := range searchResult.Hits.Hits {
		var product model.Product
		if err = json.Unmarshal(hit.Source, &product); err != nil {
			fmt.Println("[Getting Products][Unmarshal] Err=", err)
			continue
		}

		products = append(products, product)
	}

	if err != nil {
		fmt.Println("Fetching product fail: ", err)
		return nil
	}

	return products
}
