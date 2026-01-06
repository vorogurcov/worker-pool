package internal

import (
	"fmt"
	"net/http"
	"strings"
	"traceroute-optimised/internal/domain"
)

func NewHandleGetData(dataService *DataService, r *http.Request) domain.JobExecFunc {
	route := "/v1/get"

	return func() domain.Result {

		if r.Method != http.MethodGet {
			return domain.Result{
				Status: 405,
				Body:   nil,
				Error:  fmt.Errorf("method not allowed"),
				Path:   route,
			}
		}

		dataKey := r.URL.Query().Get("key")
		dataOtherKey := r.URL.Query().Get("other_key")

		var keyName string
		var searchValue string

		if dataKey != "" {
			keyName = "product_id"
			searchValue = dataKey
		} else if dataOtherKey != "" {
			keyName = "category_id"
			searchValue = dataOtherKey
		} else {
			return domain.Result{
				Status: 400,
				Body:   nil,
				Error:  fmt.Errorf("missing key parameter"),
				Path:   route,
			}
		}

		data, err := dataService.GetData(r.Context(), searchValue, keyName)
		if err != nil {
			return domain.Result{
				Status: 500,
				Body:   nil,
				Error:  fmt.Errorf("error fetching data: %v", err),
				Path:   route,
			}
		}

		return domain.Result{
			Status: 200,
			Body:   data,
			Error:  nil,
			Path:   route,
		}
	}

}

func NewHandleGetBulkData(dataService *DataService, r *http.Request) domain.JobExecFunc {
	route := "/v1/bulk"

	return func() domain.Result {
		if r.Method != http.MethodGet {
			return domain.Result{
				Status: 405,
				Body:   nil,
				Error:  fmt.Errorf("method not allowed"),
				Path:   route,
			}
		}

		key := r.URL.Query().Get("product_id")
		otherKey := r.URL.Query().Get("category_id")

		var keyName string

		if key != "" {
			keyName = "product_id"
		} else if otherKey != "" {
			keyName = "category_id"
		} else {
			return domain.Result{
				Status: 400,
				Body:   nil,
				Error:  fmt.Errorf("missing key parameter"),
				Path:   route,
			}
		}
		dataKeys := strings.Split(r.URL.Query().Get("data_keys"), ",")

		data, err := dataService.GetBulkData(r.Context(), dataKeys, keyName)
		if err != nil {
			return domain.Result{
				Status: 500,
				Body:   nil,
				Error:  fmt.Errorf("error fetching data"),
				Path:   route,
			}
		}

		return domain.Result{
			Status: 200,
			Body:   data,
			Error:  nil,
			Path:   route,
		}
	}

}
