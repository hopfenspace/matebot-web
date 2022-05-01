package sdk

import (
	"errors"
	"fmt"
	"github.com/hopfenspace/matebot-web/sdk/responses"
	"github.com/monaco-io/request"
	"strings"
	"time"
)

type SDK interface {
	Login(username string, password string) error
	SetApplicationID(name string) error
	Consume(userID uint, consumableID uint, amount uint) error
	SearchUserByAliasName(aliasName string) (uint, error)
	CreateUser(name string, permission bool, external bool, voucherID uint) error
}

type sdk struct {
	accessToken   string
	baseUrl       string
	applicationID uint
}

func New(baseUrl string) SDK {
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	if !strings.HasSuffix(baseUrl, "/v1/") {
		baseUrl += "v1/"
	}

	return &sdk{
		baseUrl: baseUrl,
	}
}

func (s *sdk) Login(username string, password string) error {
	data := make(map[string]string)
	data["grant_type"] = "password"
	data["username"] = username
	data["password"] = password

	result := responses.Login{}

	c := request.Client{
		URL:           s.baseUrl + "login",
		Method:        "POST",
		MultipartForm: request.MultipartForm{Fields: data},
		Timeout:       time.Second,
	}

	if resp := c.Send().ScanJSON(&result); !resp.OK() {
		return resp.Error()
	} else {
		if resp.Code() != 200 {
			return errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			s.accessToken = result.AccessToken
		}
	}

	return nil
}

func (s *sdk) Consume(clientID uint, consumableID uint, amount uint) error {
	data := make(map[string]string)
	data["user_id"] = fmt.Sprint(clientID)
	data["consumable_id"] = fmt.Sprint(consumableID)
	data["amount"] = fmt.Sprint(amount)

	c := request.Client{
		URL:     s.baseUrl + "transactions",
		Method:  "POST",
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
		JSON:    data,
		Timeout: time.Second,
	}

	if resp := c.Send(); !resp.OK() {
		return resp.Error()
	} else {
		if resp.Code() != 201 {
			return errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			return nil
		}
	}
}

func (s *sdk) SetApplicationID(name string) error {
	var applications []responses.Application
	c := request.Client{
		URL:     s.baseUrl + "applications",
		Method:  "GET",
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
		Query:   map[string]string{"name": name},
		Timeout: time.Second,
	}

	if resp := c.Send().ScanJSON(&applications); !resp.OK() {
		return resp.Error()
	} else {
		if resp.Code() != 200 {
			return errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			if len(applications) != 1 {
				return errors.New(fmt.Sprintf("Unable to determine application ID: Found %d results", len(applications)))
			} else {
				s.applicationID = applications[0].Id
				return nil
			}
		}
	}
}

func (s *sdk) SearchUserByAliasName(aliasName string) (uint, error) {
	return 0, nil
}

func (s *sdk) CreateUser(name string, permission bool, external bool, voucherID uint) error {
	return nil
}
