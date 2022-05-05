package sdk

import (
	"errors"
	"fmt"
	"github.com/hopfenspace/matebot-web/conf"
	"github.com/hopfenspace/matebot-web/sdk/core"
	"github.com/hopfenspace/matebot-web/sdk/requests"
	"github.com/hopfenspace/matebot-web/sdk/responses"
	"github.com/monaco-io/request"
	"strings"
	"time"
)

type SDK interface {
	Login(username string, password string) error
	UpdateApplicationID() (uint, error)
	UserCreate() (*core.User, error)
	UserSetFlagExternal(aliasName string, external bool) error
	UserSetFlagPermission(aliasName string, permission bool) error
	UserSetVoucher(debtor string, voucher string) error
	UserDisable(aliasName string) error
	AliasCreate(userID uint, username string) (*core.Alias, error)
	AliasDelete(aliasID uint) error
}

type sdk struct {
	accessToken   string
	baseUrl       string
	applicationID uint
	config        *conf.Config
}

func New(baseUrl string, config *conf.Config) SDK {
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}

	if !strings.HasSuffix(baseUrl, "/v1/") {
		baseUrl += "v1/"
	}

	return &sdk{
		baseUrl: baseUrl,
		config:  config,
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

func (s *sdk) UpdateApplicationID() (uint, error) {
	var result []core.Application

	data := make(map[string]string)
	data["name"] = s.config.MateBot.User

	c := request.Client{
		URL:           s.baseUrl + "applications",
		Method:        "GET",
		MultipartForm: request.MultipartForm{Fields: data},
		Timeout:       time.Second,
		Header:        map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
	}

	if resp := c.Send().ScanJSON(&result); !resp.OK() {
		return 0, resp.Error()
	} else {
		if resp.Code() != 200 {
			return 0, errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			if len(result) != 1 {
				return 0, errors.New(fmt.Sprintf("Cannot unpack %d results", len(result)))
			} else {
				s.applicationID = result[0].ID
				return s.applicationID, nil
			}
		}
	}
}

func (s *sdk) UserCreate() (*core.User, error) {
	result := core.User{}

	c := request.Client{
		URL:    s.baseUrl + "users",
		Method: "POST",
		JSON: requests.CreateUser{
			External:   true,
			Permission: false,
		},
		Timeout: time.Second,
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
	}

	if resp := c.Send().ScanJSON(&result); !resp.OK() {
		return nil, resp.Error()
	} else {
		if resp.Code() != 201 {
			return nil, errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			return &result, nil
		}
	}
}

func (s *sdk) UserSetFlagExternal(aliasName string, external bool) error {
	c := request.Client{
		URL:    s.baseUrl + "users/setFlag",
		Method: "POST",
		JSON: requests.UserFlagExternal{
			User:     aliasName,
			External: external,
		},
		Timeout: time.Second,
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
	}

	if resp := c.Send(); !resp.OK() {
		return resp.Error()
	} else {
		if resp.Code() != 200 {
			return errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			return nil
		}
	}
}

func (s *sdk) UserSetFlagPermission(aliasName string, permission bool) error {
	c := request.Client{
		URL:    s.baseUrl + "users/setFlag",
		Method: "POST",
		JSON: requests.UserFlagPermission{
			User:       aliasName,
			Permission: permission,
		},
		Timeout: time.Second,
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
	}

	if resp := c.Send(); !resp.OK() {
		return resp.Error()
	} else {
		if resp.Code() != 200 {
			return errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			return nil
		}
	}
}

func (s *sdk) UserSetVoucher(debtor string, voucher string) error {
	c := request.Client{
		URL:    s.baseUrl + "users/setVoucher",
		Method: "POST",
		JSON: requests.UserVoucher{
			Debtor:  debtor,
			Voucher: voucher,
		},
		Timeout: time.Second,
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
	}

	if resp := c.Send(); !resp.OK() {
		return resp.Error()
	} else {
		if resp.Code() != 200 {
			return errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			return nil
		}
	}
}

func (s *sdk) UserDisable(aliasName string) error {
	c := request.Client{
		URL:    s.baseUrl + "users/setVoucher",
		Method: "POST",
		JSON: requests.UserDisable{
			User: aliasName,
		},
		Timeout: time.Second,
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
	}

	if resp := c.Send(); !resp.OK() {
		return resp.Error()
	} else {
		if resp.Code() != 200 {
			return errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			return nil
		}
	}
}

func (s *sdk) AliasCreate(userID uint, username string) (*core.Alias, error) {
	var result core.Alias

	c := request.Client{
		URL:    s.baseUrl + "aliases",
		Method: "POST",
		JSON: requests.CreateAlias{
			Username:      username,
			ApplicationID: s.applicationID,
			UserID:        userID,
		},
		Timeout: time.Second,
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
	}

	if resp := c.Send().ScanJSON(&result); !resp.OK() {
		return nil, resp.Error()
	} else {
		if resp.Code() != 201 {
			return nil, errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			return &result, nil
		}
	}
}

func (s *sdk) AliasDelete(aliasID uint) error {
	c := request.Client{
		URL:    s.baseUrl + "aliases",
		Method: "DELETE",
		JSON: requests.DeleteAlias{
			AliasID: aliasID,
		},
		Timeout: time.Second,
		Header:  map[string]string{"Authorization": fmt.Sprintf("Bearer %s", s.accessToken)},
	}

	if resp := c.Send(); !resp.OK() {
		return resp.Error()
	} else {
		if resp.Code() != 204 {
			return errors.New(fmt.Sprintf("%d: %s", resp.Code(), resp.String()))
		} else {
			return nil
		}
	}
}
