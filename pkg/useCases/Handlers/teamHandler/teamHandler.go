package teamHandler

import (
	"sales-manager-back/internal/data/infrastructure/teamRepository"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/user"
	"sales-manager-back/pkg/useCases/Helpers/firebaseHelper"
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

func Create(newUser *user.User) (*user.User, response.Status) {
	// First, create the user in our local database
	createdUser, status := teamRepository.Create(newUser)
	if status != response.StatusCreated && status != response.StatusOk {
		return createdUser, status
	}

	// Attempt to create or fetch in Firebase Auth
	if createdUser != nil {
		// Check if user already exists
		firebaseUser, err := firebaseHelper.GetUserByEmail(createdUser.Email)
		if err != nil || firebaseUser == nil {
			tempPassword := "Temporal123!" // In production, generate a secure random string or use Action Links
			firebaseUser, err = firebaseHelper.CreateUser(createdUser.Email, tempPassword, createdUser.Name)
		}

		if firebaseUser != nil {
			// Set custom claims
			claims := map[string]interface{}{
				"tenantId": float64(createdUser.TenantID),
				"userId":   float64(createdUser.ID),
				"role":     string(createdUser.Role),
			}
			firebaseHelper.SetCustomUserClaims(firebaseUser.UID, claims)

			// Save the FirebaseUID back to our local database so GetByFirebaseUID works on login
			teamRepository.Update(createdUser.TenantID, createdUser.ID, map[string]interface{}{
				"firebase_uid": firebaseUser.UID,
			})
			createdUser.FirebaseUID = firebaseUser.UID
		}
	}

	return createdUser, status
}

func Update(tenantID, userID uint, updates map[string]interface{}) response.Status {
	return teamRepository.Update(tenantID, userID, updates)
}
