package elasticsearch

import "context"

type Indexable interface {
	GetID() string
}

func IndexDocument(indexName string, document Indexable) error {
	_, err := GetClient().Index(indexName).Id(document.GetID()).Request(document).Do(context.TODO())
	return err
}
