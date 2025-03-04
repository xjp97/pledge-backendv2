package services

import (
	"fmt"
	"pledge-backendv2/api/common/statecode"
	"pledge-backendv2/api/models"
	"pledge-backendv2/api/models/request"
	"pledge-backendv2/log"
)

type SearchService struct{}

func NewSearch() *SearchService {
	return &SearchService{}
}

// 查询 poolbases 表，分页查询
func (s *SearchService) Search(req *request.Search) (int, int64, []models.Pool) {

	whereCondition := fmt.Sprintf(`chain_id = '%v'`, req.ChainID)
	if req.LendTokenSymbol != "" {
		whereCondition += fmt.Sprintf(`lend_token_symbol = '%v'`, req.LendTokenSymbol)
	}
	if req.State != "" {
		whereCondition += fmt.Sprintf(`state = '%v'`, req.State)
	}
	err, tatol, pools := models.NewPool().Search(req, whereCondition)
	if err != nil {
		log.Logger.Error(err.Error())
		return statecode.CommonErrServerErr, 0, nil
	}
	return 0, tatol, pools

}
