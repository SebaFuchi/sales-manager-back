package saleHandler

import (
	"fmt"
	"log"
	"sales-manager-back/internal/data/infrastructure/saleRepository"
	"sales-manager-back/pkg/domain/client"
	"sales-manager-back/pkg/domain/commission"
	"sales-manager-back/pkg/domain/principal"
	"sales-manager-back/pkg/domain/response"
	"sales-manager-back/pkg/domain/sale"
	"sales-manager-back/pkg/domain/user"
	"sales-manager-back/pkg/useCases/Helpers/databaseHelper"
)

func GetAll(tenantID uint) (interface{}, response.Status) {
	return saleRepository.GetAll(tenantID)
}

func GetByVendedor(tenantID, vendedorID uint) (interface{}, response.Status) {
	return saleRepository.GetByVendedor(tenantID, vendedorID)
}

func GetByID(tenantID, saleID uint) (interface{}, response.Status) {
	return saleRepository.GetByID(tenantID, saleID)
}

// Create creates a new sale and auto-generates:
// 1. A Commission record based on the principal's base commission and agent's split
// 2. An AccountMovement (invoice) on the client's current account
func Create(newSale *sale.Sale) (interface{}, response.Status) {
	db := databaseHelper.Db

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create the Sale with its Items
	if err := tx.Create(newSale).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating sale: %v", err)
		return nil, response.StatusInternalServerError
	}

	// 2. Lookup the Principal to get the base commission percentage
	var principalItem principal.Principal
	if err := tx.Where("tenant_id = ? AND id = ?", newSale.TenantID, newSale.PrincipalID).
		First(&principalItem).Error; err != nil {
		log.Printf("Warning: could not find principal %d for commission calculation: %v", newSale.PrincipalID, err)
		// Don't fail the sale creation, just skip commission
	} else {
		// 3. Lookup the Agent to get split percentage
		var agent user.User
		var splitPercentage float64
		var subStatus commission.EstadoComisionSub = commission.EstadoSubNoAplica

		if err := tx.Where("tenant_id = ? AND id = ?", newSale.TenantID, newSale.AgentID).
			First(&agent).Error; err != nil {
			log.Printf("Warning: could not find agent %d: %v", newSale.AgentID, err)
		} else {
			splitPercentage = agent.SplitPercentageSub
			if agent.TeamRole == user.EquipoRolSubvendedor && splitPercentage > 0 {
				subStatus = commission.EstadoSubRetenida
			}
		}

		// Calculate commission amounts
		percentage := principalItem.BaseCommission
		if newSale.CommissionOverride != nil {
			percentage = *newSale.CommissionOverride
		}
		gross := newSale.NetAmount * (percentage / 100)
		expenses := 0.0 // Can be extended later
		net := gross - expenses

		// Distribution
		var ownerDist, subDist float64
		if splitPercentage > 0 {
			subDist = net * (splitPercentage / 100)
			ownerDist = net - subDist
		} else {
			ownerDist = net
			subDist = 0
		}

		newCommission := commission.Commission{
			TenantID:          newSale.TenantID,
			SaleID:            newSale.ID,
			PrincipalID:       newSale.PrincipalID,
			PrincipalName:     newSale.PrincipalName,
			ResponsibleID:     newSale.AgentID,
			ResponsibleName:   newSale.AgentName,
			CalculationBase:   newSale.NetAmount,
			Percentage:        percentage,
			Gross:             gross,
			Expenses:          expenses,
			Net:               net,
			OwnerDistribution: ownerDist,
			SubDistribution:   subDist,
			FactoryStatus:     commission.EstadoPendCobro,
			SubStatus:         subStatus,
		}

		if err := tx.Create(&newCommission).Error; err != nil {
			log.Printf("Warning: failed to auto-create commission for sale %d: %v", newSale.ID, err)
			// Don't fail the sale, commission can be created manually
		}
	}

	// 4. Auto-create AccountMovement (invoice) on the client's current account
	// Calculate the client's running balance
	var lastBalance float64
	var lastMovement client.AccountMovement
	if err := tx.Where("client_id = ?", newSale.ClientID).
		Order("id DESC").
		First(&lastMovement).Error; err == nil {
		lastBalance = lastMovement.Balance
	}

	newBalance := lastBalance + newSale.NetAmount

	movement := client.AccountMovement{
		ClientID:         newSale.ClientID,
		Date:             newSale.Date,
		VoucherNumber:    fmt.Sprintf("OP-%d", newSale.ID),
		Type:             client.TipoFactura,
		Detail:           fmt.Sprintf("Operación #%d — %s", newSale.ID, newSale.PrincipalName),
		Debit:            newSale.NetAmount,
		Credit:           0,
		Balance:          newBalance,
		DueDate:          nil, // Can be extended with payment terms later
		CollectionStatus: client.EstadoAlDia,
	}

	if err := tx.Create(&movement).Error; err != nil {
		log.Printf("Warning: failed to auto-create account movement for sale %d: %v", newSale.ID, err)
		// Don't fail the sale, movement can be created manually
	}

	// Commit
	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing sale transaction: %v", err)
		return nil, response.StatusInternalServerError
	}

	return newSale, response.StatusCreated
}

func Update(tenantID, saleID uint, updates map[string]interface{}) response.Status {
	return saleRepository.Update(tenantID, saleID, updates)
}

func Delete(tenantID, saleID uint) response.Status {
	db := databaseHelper.Db

	tx := db.Begin()

	// Delete associated commissions first
	if err := tx.Where("tenant_id = ? AND sale_id = ?", tenantID, saleID).
		Delete(&commission.Commission{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting commissions for sale %d: %v", saleID, err)
		return response.StatusInternalServerError
	}

	// Delete the sale (Items cascade via GORM)
	result := tx.Where("tenant_id = ? AND id = ?", tenantID, saleID).Delete(&sale.Sale{})
	if err := result.Error; err != nil {
		tx.Rollback()
		return response.StatusInternalServerError
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return response.StatusNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return response.StatusInternalServerError
	}

	return response.StatusOk
}
