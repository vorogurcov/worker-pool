package internal

import (
	"fmt"
	"net/http"
	"strings"
	"traceroute-optimised/internal/domain"
)

func NewHandleGetData(dataService *DataService, r *http.Request) domain.JobExecFunc {
	return func() domain.Result {

		if r.Method != http.MethodGet {
			return domain.Result{
				Status: 405,
				Body:   nil,
				Error:  fmt.Errorf("method not allowed"),
			}
		}

		dataKey := r.URL.Query().Get("key")
		dataOtherKey := r.URL.Query().Get("other_key")

		var keyName string
		var searchValue string

		if dataKey != "" {
			keyName = "key"
			searchValue = dataKey
		} else if dataOtherKey != "" {
			keyName = "other_key"
			searchValue = dataOtherKey
		} else {
			return domain.Result{
				Status: 400,
				Body:   nil,
				Error:  fmt.Errorf("missing key parameter"),
			}
		}

		data, err := dataService.GetData(r.Context(), searchValue, keyName)
		if err != nil {
			return domain.Result{
				Status: 500,
				Body:   nil,
				Error:  fmt.Errorf("error fetching data: %v", err),
			}
		}

		return domain.Result{
			Status: 200,
			Body:   data,
			Error:  nil,
		}
	}

}

func NewHandleGetBulkData(dataService *DataService, r *http.Request) domain.JobExecFunc {
	return func() domain.Result {
		if r.Method != http.MethodGet {
			return domain.Result{
				Status: 405,
				Body:   nil,
				Error:  fmt.Errorf("method not allowed"),
			}
		}

		key := r.URL.Query().Get("key")
		otherKey := r.URL.Query().Get("other_key")

		var keyName string

		if key != "" {
			keyName = "key"
		} else if otherKey != "" {
			keyName = "other_key"
		} else {
			return domain.Result{
				Status: 400,
				Body:   nil,
				Error:  fmt.Errorf("missing key parameter"),
			}
		}
		dataKeys := strings.Split(r.URL.Query().Get("data_keys"), ",")

		data, err := dataService.GetBulkData(r.Context(), dataKeys, keyName)
		if err != nil {
			return domain.Result{
				Status: 500,
				Body:   nil,
				Error:  fmt.Errorf("error fetching data"),
			}
		}

		return domain.Result{
			Status: 200,
			Body:   data,
			Error:  nil,
		}
	}

}
