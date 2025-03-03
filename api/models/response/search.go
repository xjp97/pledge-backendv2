package response

import "pledge-backendv2/api/models"

type Search struct {
	Count int64         `json:"count"`
	Rows  []models.Pool `json:"rows"`
}
