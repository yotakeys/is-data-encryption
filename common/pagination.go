package common

import (
	"gin-gorm-clean-template/entity"
	"math"

	"gorm.io/gorm"
)

func Pagination(pagination *entity.Pagination, totalData int64) func(db *gorm.DB) *gorm.DB {
	pagination.TotalData = totalData
	maxPage := int(math.Ceil(float64(totalData) / float64(pagination.PerPage)))
	pagination.MaxPage = maxPage
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.PerPage)
	}
}