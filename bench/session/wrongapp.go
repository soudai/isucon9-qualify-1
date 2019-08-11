package session

import (
	"bytes"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/isucon/isucon9-qualify/bench/fails"
)

const (
	ItemMinPrice    = 100
	ItemMaxPrice    = 1000000
	ItemPriceErrMsg = "商品価格は100円以上、1,000,000円以下にしてください"
)

type resErr struct {
	Error string `json:"error"`
}

func secureRandomStr(b int) string {
	k := make([]byte, b)
	if _, err := crand.Read(k); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", k)
}

func (s *Session) LoginWithWrongPassword(accountName, password string) error {
	b, _ := json.Marshal(reqLogin{
		AccountName: accountName,
		Password:    password,
	})

	req, err := s.newPostRequest(ShareTargetURLs.AppURL, "/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fails.NewError(err, "POST /login: リクエストに失敗しました")
	}

	res, err := s.Do(req)
	if err != nil {
		return fails.NewError(err, "POST /login: リクエストに失敗しました")
	}
	defer res.Body.Close()

	msg, err := checkStatusCode(res, http.StatusUnauthorized)
	if err != nil {
		return fails.NewError(err, "POST /login: "+msg)
	}

	re := resErr{}
	err = json.NewDecoder(res.Body).Decode(&re)
	if err != nil {
		return fails.NewError(err, "POST /login: JSONデコードに失敗しました")
	}

	return nil
}

func (s *Session) SellWithWrongCSRFToken(name string, price int, description string, categoryID int) error {
	b, _ := json.Marshal(reqSell{
		CSRFToken:   secureRandomStr(20),
		Name:        name,
		Price:       price,
		Description: description,
		CategoryID:  categoryID,
	})
	req, err := s.newPostRequest(ShareTargetURLs.AppURL, "/sell", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fails.NewError(err, "POST /sell: リクエストに失敗しました")
	}

	res, err := s.Do(req)
	if err != nil {
		return fails.NewError(err, "POST /sell: リクエストに失敗しました")
	}
	defer res.Body.Close()

	msg, err := checkStatusCode(res, http.StatusUnprocessableEntity)
	if err != nil {
		return fails.NewError(err, "POST /sell: "+msg)
	}

	re := resErr{}
	err = json.NewDecoder(res.Body).Decode(&re)
	if err != nil {
		return fails.NewError(err, "POST /sell: JSONデコードに失敗しました")
	}

	return nil
}

func (s *Session) SellWithWrongPrice(name string, price int, description string, categoryID int) error {
	b, _ := json.Marshal(reqSell{
		CSRFToken:   s.csrfToken,
		Name:        name,
		Price:       price,
		Description: description,
		CategoryID:  categoryID,
	})
	req, err := s.newPostRequest(ShareTargetURLs.AppURL, "/sell", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fails.NewError(err, "POST /sell: リクエストに失敗しました")
	}

	res, err := s.Do(req)
	if err != nil {
		return fails.NewError(err, "POST /sell: リクエストに失敗しました")
	}
	defer res.Body.Close()

	msg, err := checkStatusCode(res, http.StatusBadRequest)
	if err != nil {
		return fails.NewError(err, "POST /sell: "+msg)
	}

	re := resErr{}
	err = json.NewDecoder(res.Body).Decode(&re)
	if err != nil {
		return fails.NewError(err, "POST /sell: JSONデコードに失敗しました")
	}

	if re.Error != ItemPriceErrMsg {
		return fails.NewError(err, "POST /sell: 商品価格は100円以上、1,000,000円以下しか出品できません")
	}

	return nil
}

func (s *Session) BuyWithWrongCSRFToken(itemID int64, token string) error {
	b, _ := json.Marshal(reqBuy{
		CSRFToken: secureRandomStr(20),
		ItemID:    itemID,
		Token:     token,
	})
	req, err := s.newPostRequest(ShareTargetURLs.AppURL, "/buy", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fails.NewError(err, "POST /buy: リクエストに失敗しました")
	}

	res, err := s.Do(req)
	if err != nil {
		return fails.NewError(err, "POST /buy: リクエストに失敗しました")
	}
	defer res.Body.Close()

	msg, err := checkStatusCode(res, http.StatusUnprocessableEntity)
	if err != nil {
		return fails.NewError(err, "POST /buy: "+msg)
	}

	re := resErr{}
	err = json.NewDecoder(res.Body).Decode(&re)
	if err != nil {
		return fails.NewError(err, "POST /buy: JSONデコードに失敗しました")
	}

	return nil
}

func (s *Session) BuyWithFailed(itemID int64, token string, expectedStatus int, expectedMsg string) error {
	b, _ := json.Marshal(reqBuy{
		CSRFToken: s.csrfToken,
		ItemID:    itemID,
		Token:     token,
	})
	req, err := s.newPostRequest(ShareTargetURLs.AppURL, "/buy", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fails.NewError(err, "POST /buy: リクエストに失敗しました")
	}

	res, err := s.Do(req)
	if err != nil {
		return fails.NewError(err, "POST /buy: リクエストに失敗しました")
	}
	defer res.Body.Close()

	msg, err := checkStatusCode(res, expectedStatus)
	if err != nil {
		return fails.NewError(err, "POST /buy: "+msg)
	}

	re := resErr{}
	err = json.NewDecoder(res.Body).Decode(&re)
	if err != nil {
		return fails.NewError(err, "POST /buy: JSONデコードに失敗しました")
	}

	if re.Error != expectedMsg {
		return fails.NewError(err, "POST /buy: "+expectedMsg+"というエラーではありません")
	}

	return nil
}

func (s *Session) ShipWithWrongCSRFToken(itemID int64) error {
	b, _ := json.Marshal(reqShip{
		CSRFToken: secureRandomStr(20),
		ItemID:    itemID,
	})
	req, err := s.newPostRequest(ShareTargetURLs.AppURL, "/ship", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fails.NewError(err, "POST /ship: リクエストに失敗しました")
	}

	res, err := s.Do(req)
	if err != nil {
		return fails.NewError(err, "POST /ship: リクエストに失敗しました")
	}
	defer res.Body.Close()

	msg, err := checkStatusCode(res, http.StatusUnprocessableEntity)
	if err != nil {
		return fails.NewError(err, "POST /ship: "+msg)
	}

	re := resErr{}
	err = json.NewDecoder(res.Body).Decode(&re)
	if err != nil {
		return fails.NewError(err, "POST /ship: JSONデコードに失敗しました")
	}

	return nil
}

func (s *Session) ShipWithWrongSeller(itemID int64) error {
	b, _ := json.Marshal(reqShip{
		CSRFToken: secureRandomStr(20),
		ItemID:    itemID,
	})
	req, err := s.newPostRequest(ShareTargetURLs.AppURL, "/ship", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fails.NewError(err, "POST /ship: リクエストに失敗しました")
	}

	res, err := s.Do(req)
	if err != nil {
		return fails.NewError(err, "POST /ship: リクエストに失敗しました")
	}
	defer res.Body.Close()

	msg, err := checkStatusCode(res, http.StatusForbidden)
	if err != nil {
		return fails.NewError(err, "POST /ship: "+msg)
	}

	re := resErr{}
	err = json.NewDecoder(res.Body).Decode(&re)
	if err != nil {
		return fails.NewError(err, "POST /ship: JSONデコードに失敗しました")
	}

	if re.Error != "権限がありません" {
		return fails.NewError(err, "POST /ship: 権限がないエラーが発生していません")
	}

	return nil
}
