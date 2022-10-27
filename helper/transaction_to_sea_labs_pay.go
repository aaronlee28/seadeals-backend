package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"seadeals-backend/apperror"
	"seadeals-backend/config"
	"strconv"
	"strings"
	"time"
)

func TransactionToSeaLabsPay(accountNumber string, amount string, sign string, callback string, trxType string) (string, uint, error) {

	if os.Getenv("ENV") == "testing" {
		return "", 0, nil
	}

	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	redirURL := config.Config.RedirectPaymentBase + "/transactions/post-slp-" + trxType

	data := url.Values{}
	data.Set("card_number", accountNumber)
	data.Set("amount", amount)
	data.Set("merchant_code", config.Config.SeaLabsPayMerchantCode)
	data.Set("redirect_url", redirURL)
	data.Set("callback_url", config.Config.NgrokURL+callback)
	data.Set("signature", sign)
	encodeData := data.Encode()

	fmt.Println(sign)
	req, err := http.NewRequest(http.MethodPost, config.Config.SeaLabsPayTransactionURL, strings.NewReader(encodeData))
	if err != nil {
		return "", 0, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer func(Body io.ReadCloser) {
		err2 := Body.Close()
		if err2 != nil {
			fmt.Println("Error Closing Client")
		}
	}(response.Body)

	if response.StatusCode == http.StatusSeeOther {
		redirectUrl, err2 := response.Location()
		if err2 != nil {
			return "", 0, err2
		}

		TxnID, err3 := strconv.ParseUint(redirectUrl.Query().Get("txn_id"), 10, 64)
		if err3 != nil {
			return "", 0, err3
		}
		return redirectUrl.String(), uint(TxnID), nil
	} else {
		type seaLabsPayError struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Data    struct {
			} `json:"data"`
		}
		var j seaLabsPayError
		err = json.NewDecoder(response.Body).Decode(&j)
		if err != nil {
			panic(err)
		}
		return "", 0, apperror.BadRequestError(j.Message)
	}
}
