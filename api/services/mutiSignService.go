package services

import (
	"encoding/json"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models"
	"pledge-backendv2/api/models/request"
	"pledge-backendv2/api/models/response"
)

type MutiSignService struct{}

func NewMutiSign() *MutiSignService {
	return &MutiSignService{}
}

func (c *MutiSignService) SetMultiSign(mutiSign *request.SetMultiSign) (int, error) {
	err := models.NewMultiSign().Set(mutiSign)
	if err != nil {
		return statecode.CommonErrServerErr, err
	}
	return statecode.CommonSuccess, nil
}

// GetMultiSign Get Multi-Sign
func (c *MutiSignService) GetMultiSign(mutiSign *response.MultiSign, chainId int) (int, error) {
	//db get
	multiSignModel := models.NewMultiSign()
	err := multiSignModel.Get(chainId)
	if err != nil {
		return statecode.CommonErrServerErr, err
	}
	var multiSignAccount []string
	_ = json.Unmarshal([]byte(multiSignModel.MultiSignAccount), &multiSignAccount)

	mutiSign.SpName = multiSignModel.SpName
	mutiSign.SpToken = multiSignModel.SpToken
	mutiSign.JpName = multiSignModel.JpName
	mutiSign.JpToken = multiSignModel.JpToken
	mutiSign.SpAddress = multiSignModel.SpAddress
	mutiSign.JpAddress = multiSignModel.JpAddress
	mutiSign.SpHash = multiSignModel.SpHash
	mutiSign.JpHash = multiSignModel.JpHash
	mutiSign.MultiSignAccount = multiSignAccount
	return statecode.CommonSuccess, nil
}
