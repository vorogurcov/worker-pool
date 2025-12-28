package internal

import (
	"context"
	"traceroute-optimised/internal/domain"
)

type DataService struct {
	DataRepo domain.IDataRepository
}

func (ds *DataService) GetData(ctx context.Context, dataKey string, keyName string) (*domain.DataModel, error) {
	return ds.DataRepo.GetData(ctx, dataKey, keyName)
}

func (ds *DataService) GetBulkData(ctx context.Context, dataKeys []string, keyName string) (*[]domain.DataModel, error) {
	return ds.DataRepo.GetBulkData(ctx, dataKeys, keyName)
}

func (ds *DataService) UpdateData(ctx context.Context, newData domain.DataModel, keyName string) (*domain.DataModel, error) {
	return ds.DataRepo.UpdateData(ctx, newData, keyName)
}
