package teamHandler

import (
	"sales-manager-back/internal/data/infrastructure/teamRepository"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/user"
)

func GetAll(tenantID uint) (interface{}, response.Status) {
	return teamRepository.GetAll(tenantID)
}

func GetByID(tenantID, userID uint) (interface{}, response.Status) {
	return teamRepository.GetByID(tenantID, userID)
}

func GetByFirebaseUID(firebaseUID string) (interface{}, response.Status) {
	return teamRepository.GetByFirebaseUID(firebaseUID)
}

func Create(newUser *user.User) (interface{}, response.Status) {
	return teamRepository.Create(newUser)
}

func Update(tenantID, userID uint, updates map[string]interface{}) response.Status {
	return teamRepository.Update(tenantID, userID, updates)
}
