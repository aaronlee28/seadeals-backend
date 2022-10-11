package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"seadeals-backend/dto"
	"strconv"
	"strings"
)

const (
	REG = "REG"
)

func findRegularPrice(res []dto.DeliveryCalculateRes) *dto.DeliveryCalculateReturn {
	for _, cost := range res[0].Costs {
		if cost.Service == REG {
			splitEta := strings.Split(cost.Cost[0].Etd, "-")
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

	client := &http.Client{}
	URL := config.Config.ShippingURL
	requestStr := `{` +
		`"origin_city":` + r.OriginCity +
		`, "destination_city":` + r.DestinationCity +
		`, "weight":` + r.Weight +
		`, "courier":"` + r.Courier + `"` +
		`}`
	var jsonStr = []byte(requestStr)
	fmt.Println(requestStr)

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
		fmt.Println(req.Header.Get("Api-Key"))
		type shippingError struct {
			StatusCode int    `json:"status_code"`
			Code       string `json:"code"`
			Message    string `json:"message"`
		}
		var j shippingError
		err = json.NewDecoder(resp.Body).Decode(&j)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return nil, apperror.BadRequestError(j.Message)
	}

	var dtoRes []dto.DeliveryCalculateRes
	err = json.NewDecoder(resp.Body).Decode(&dtoRes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	returnRes := findRegularPrice(dtoRes)
	if returnRes == nil {
		return nil, apperror.InternalServerError("No service available for that order")
	}

	return returnRes, nil
}
