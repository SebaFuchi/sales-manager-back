package teamRepository

import (
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/user"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"

	"gorm.io/gorm"
)

func GetAll(tenantID uint) ([]user.User, response.Status) {
	var users []user.User
	db := databaseHelper.Db

	result := db.Where("tenant_id = ?", tenantID).Find(&users)

	if err := result.Error; err != nil {
		return users, response.StatusInternalServerError
	}

	return users, response.StatusOk
}

func GetByID(tenantID, userID uint) (user.User, response.Status) {
	var userItem user.User
	db := databaseHelper.Db

	result := db.Where("tenant_id = ? AND id = ?", tenantID, userID).
		First(&userItem)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return userItem, response.StatusNotFound
		}
		return userItem, response.StatusInternalServerError
	}

	return userItem, response.StatusOk
}

func GetByFirebaseUID(firebaseUID string) (user.User, response.Status) {
	var userItem user.User
	db := databaseHelper.Db

	result := db.Where("firebase_uid = ?", firebaseUID).
		First(&userItem)

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return userItem, response.StatusNotFound
		}
		return userItem, response.StatusInternalServerError
	}

	return userItem, response.StatusOk
}

func Create(newUser *user.User) (*user.User, response.Status) {
	db := databaseHelper.Db

	result := db.Create(newUser)
	if err := result.Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return newUser, response.StatusConflict
		}
		return newUser, response.StatusInternalServerError
	}

	return newUser, response.StatusCreated
}

func Update(tenantID, userID uint, updates map[string]interface{}) response.Status {
	db := databaseHelper.Db

	result := db.Model(&user.User{}).
		Where("tenant_id = ? AND id = ?", tenantID, userID).
		Updates(updates)

	if err := result.Error; err != nil {
		return response.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		return response.StatusNotFound
	}

	return response.StatusOk
}
