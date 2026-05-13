package authHandler

import (
	"log"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/tenant"
	"sales-manager-back/pkg/domain/user"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"
	"sales-manager-back/pkg/useCases/Helpers/firebaseHelper"
	"time"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	CompanyName string `json:"companyName"`
}

func Register(uid string, req RegisterRequest) (interface{}, response.Status) {
	db := databaseHelper.Db

	// ── Check if a user with this FirebaseUID already exists ──
	var existingUser user.User
	if err := db.Where("firebase_uid = ?", uid).First(&existingUser).Error; err == nil {
		// User already exists — just re-set claims and return the tenant
		var existingTenant tenant.Tenant
		db.First(&existingTenant, existingUser.TenantID)

		setClaims(uid, existingTenant.ID, existingUser.ID, existingUser.Role)
		log.Printf("Register: user already exists (UID=%s, email=%s), re-set claims", uid, req.Email)
		return existingTenant, response.StatusOk
	}

	// ── Check if a user with this email already exists (created by admin or previous partial attempt) ──
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		// User record exists but has no FirebaseUID — link it
		db.Model(&user.User{}).Where("id = ?", existingUser.ID).Update("firebase_uid", uid)

		var existingTenant tenant.Tenant
		db.First(&existingTenant, existingUser.TenantID)

		setClaims(uid, existingTenant.ID, existingUser.ID, existingUser.Role)
		log.Printf("Register: linked existing user (email=%s) to Firebase UID=%s", req.Email, uid)
		return existingTenant, response.StatusOk
	}

	// ── No existing user — create Tenant + User from scratch ──
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create the Tenant (Agency)
	newTenant := tenant.Tenant{
		Name:             req.CompanyName,
		Owner:            req.DisplayName,
		Email:            req.Email,
		Plan:             tenant.PlanStarter,
		Status:           tenant.EstadoTrial,
		RegistrationDate: time.Now().Format("2006-01-02"),
		LastActivity:     time.Now().Format("2006-01-02"),
	}

	if err := tx.Create(&newTenant).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating tenant: %v", err)
		return nil, response.StatusInternalServerError
	}

	// 2. Create the User (Agency Owner)
	newUser := user.User{
		TenantID:           newTenant.ID,
		Name:               req.DisplayName,
		Email:              req.Email,
		TeamRole:           user.EquipoRolTitular,
		Role:               user.RoleAgency,
		Status:             user.EstadoActivo,
		SplitPercentageSub: 0,
		FirebaseUID:        uid,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating user: %v", err)
		return nil, response.StatusInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction: %v", err)
		return nil, response.StatusInternalServerError
	}

	// 3. Update Firebase Custom Claims
	setClaims(uid, newTenant.ID, newUser.ID, user.RoleAgency)

	return newTenant, response.StatusCreated
}

// setClaims sets Firebase custom claims for the given user
func setClaims(uid string, tenantID, userID uint, role user.UserRole) {
	claims := map[string]interface{}{
		"tenantId": float64(tenantID),
		"userId":   float64(userID),
		"role":     string(role),
	}

	err := firebaseHelper.SetCustomUserClaims(uid, claims)
	if err != nil {
		log.Printf("CRITICAL: Failed to set Firebase claims for UID %s: %v", uid, err)
	}
}
