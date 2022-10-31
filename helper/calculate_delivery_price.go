package helper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"strconv"
	"strings"
)

const (
	JTR = "jtr"
)

func findRegularPrice(res []dto.DeliveryCalculateRes) *dto.DeliveryCalculateReturn {
	for _, cost := range res[0].Costs {
		if cost.Service != JTR {
			splitSpace := strings.Split(cost.Cost[0].Etd, " ")
			splitEta := strings.Split(splitSpace[0], "-")
			etaInt, _ := strconv.Atoi(splitEta[0])
			result := &dto.DeliveryCalculateReturn{
				Total: cost.Cost[0].Value,
				Eta:   etaInt,
			}
			return result
		}
	}
	return nil
}

func CalculateDeliveryPrice(r *dto.DeliveryCalculateReq) (*dto.DeliveryCalculateReturn, error) {
	var err error
	var req *http.Request
	var resp *http.Response

	if os.Getenv("ENV") == "testing" {
		return &dto.DeliveryCalculateReturn{Total: 1, Eta: 1}, nil
	}

	client := &http.Client{}
	URL := config.Config.ShippingURL
	requestStr := `{` +
		`"origin_city":` + r.OriginCity +
		`, "destination_city":` + r.DestinationCity +
		`, "weight":` + r.Weight +
		`, "courier":"` + r.Courier + `"` +
		`}`
	var jsonStr = []byte(requestStr)

	req, err = http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Api-Key", config.Config.ShippingKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		type shippingError struct {
			StatusCode int    `json:"status_code"`
			Code       string `json:"code"`
			Message    string `json:"message"`
		}
		var j shippingError
		err = json.NewDecoder(resp.Body).Decode(&j)
		if err != nil {
			return nil, err
		}

		if config.Config.ShippingActionOnError == "ignore" {
			return &dto.DeliveryCalculateReturn{
				Total: 20000,
				Eta:   2,
			}, nil
		} else {
			return nil, apperror.BadRequestError(j.Message + ": Error dalam delivery API")
		}
	}

	var dtoRes []dto.DeliveryCalculateRes
	err = json.NewDecoder(resp.Body).Decode(&dtoRes)
	if err != nil {
		return nil, err
	}

	returnRes := findRegularPrice(dtoRes)
	if returnRes == nil {
		return nil, apperror.InternalServerError("No service available for that order")
	}

	return returnRes, nil
}
