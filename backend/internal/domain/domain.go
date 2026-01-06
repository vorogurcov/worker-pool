package domain

import (
	"context"
	"encoding/json"
)

type DataModel struct {
	ProductID string          `json:"product_id"`
	Payload   json.RawMessage `json:"payload"`
}

type IDataRepository interface {
	GetData(ctx context.Context, dataKey string, keyName string) (*DataModel, error)
	GetBulkData(ctx context.Context, dataKeys []string, keyName string) (*[]DataModel, error)
	//UpdateData(ctx context.Context, newData DataModel, keyName string) (*DataModel, error)
}

type IDataService interface {
	GetData(ctx context.Context, dataKey string, keyName string) (*DataModel, error)
	GetBulkData(ctx context.Context, dataKeys []string, keyName string) (*[]DataModel, error)
	//UpdateData(ctx context.Context, newData DataModel, keyName string) (*DataModel, error)
}
