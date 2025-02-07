package services

import (
	"strings"
	"web-api/internal/pkg/models/request"

	"gorm.io/gorm"
)

func GetExtrasFromIDs(db *gorm.DB, extraIDs string) ([]request.Extrasrequest, error) {
	if extraIDs == "" {
		return []request.Extrasrequest{}, nil
	}

	var extras []request.Extrasrequest
	idList := strings.Split(extraIDs, ",") // Tách danh sách ID từ chuỗi "1,2,3"
	err := db.Where("id IN (?)", idList).Find(&extras).Error
	return extras, err
}
