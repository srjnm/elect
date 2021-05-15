package database

import (
	"elect/dto"
	"elect/models"
	"errors"
	"strconv"

	"github.com/biezhi/gorm-paginator/pagination"
	uuid "github.com/satori/go.uuid"
)

func (db *postgresDatabase) RegisterStudent(user models.User) error {
	res := db.connection.Create(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (db *postgresDatabase) RegisteredStudents(userId string, paginatorParams dto.PaginatorParams) ([]models.User, error) {
	var regStudents []models.User
	if paginatorParams.Page == "" {
		res := db.connection.Where("registered_by = ?", userId).Find(&regStudents)
		if res.Error != nil {
			return nil, res.Error
		}

		return regStudents, nil
	}

	page, _ := strconv.Atoi(paginatorParams.Page)
	limit, _ := strconv.Atoi(paginatorParams.Limit)

	_ = pagination.Paging(
		&pagination.Param{
			DB:      db.connection.Where("registered_by = ?", userId),
			Page:    page,
			Limit:   limit,
			OrderBy: []string{paginatorParams.OrderBy},
			ShowSQL: false,
		},
		&regStudents,
	)

	return regStudents, nil
}

func (db *postgresDatabase) DeleteRegisteredStudent(userId string, studentUserId string) error {
	var user models.User
	var count int

	if db.connection.Model(&user).Where("registered_by = ? AND user_id = ?", userId, studentUserId).Count(&count); count == 0 {
		return errors.New("Invalid Student!")
	}

	user.UserID = uuid.FromStringOrNil(studentUserId)
	res := db.connection.Where("registered_by = ? AND user_id = ?", userId, studentUserId).Delete(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
