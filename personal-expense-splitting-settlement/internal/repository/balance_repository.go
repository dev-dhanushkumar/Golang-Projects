package repository

import (
	"personal-expense-splitting-settlement/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BalanceRepository defines methods for balance calculations
type BalanceRepository interface {
	CalculateUserBalance(userID, otherUserID uuid.UUID) (float64, error)
	CalculateAllUserBalances(userID uuid.UUID) ([]dto.BalanceItem, error)
	CalculateGroupBalances(groupID uuid.UUID) ([]dto.GroupBalanceItem, error)
}

type balanceRepository struct {
	db *gorm.DB
}

// NewBalanceRepository creates a new balance repository instance
func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepository{db: db}
}

// BalanceResult is a helper struct for SQL query results
type BalanceResult struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Balance   float64
}

// GroupBalanceResult is a helper struct for group balance SQL query results
type GroupBalanceResult struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	TotalPaid float64
	TotalOwed float64
}

// CalculateUserBalance calculates the balance between two users
// Returns positive if otherUser owes userID, negative if userID owes otherUser
func (r *balanceRepository) CalculateUserBalance(userID, otherUserID uuid.UUID) (float64, error) {
	// Calculate from expenses: amount user paid for other - amount other paid for user
	var balance float64
	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN ep1.user_id = ? AND ep2.user_id = ? 
				THEN ep2.owed_amount 
				ELSE 0 
			END - 
			CASE 
				WHEN ep1.user_id = ? AND ep2.user_id = ? 
				THEN ep1.owed_amount 
				ELSE 0 
			END
		), 0) as balance
		FROM expense_participants ep1
		JOIN expense_participants ep2 ON ep1.expense_id = ep2.expense_id
		WHERE ((ep1.user_id = ? AND ep2.user_id = ?) OR (ep1.user_id = ? AND ep2.user_id = ?))
		AND ep1.user_id != ep2.user_id
	`
	err := r.db.Raw(query, userID, otherUserID, otherUserID, userID, userID, otherUserID, otherUserID, userID).
		Scan(&balance).Error
	if err != nil {
		return 0, err
	}

	// Subtract settlements: payments made by userID to otherUser - payments made by otherUser to userID
	var settlementBalance float64
	settlementQuery := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN payer_id = ? AND payee_id = ? 
				THEN -amount 
				WHEN payer_id = ? AND payee_id = ? 
				THEN amount 
				ELSE 0 
			END
		), 0) as settlement_balance
		FROM settlements
		WHERE ((payer_id = ? AND payee_id = ?) OR (payer_id = ? AND payee_id = ?))
		AND is_confirmed = true
		AND deleted_at IS NULL
	`
	err = r.db.Raw(settlementQuery, userID, otherUserID, otherUserID, userID, userID, otherUserID, otherUserID, userID).
		Scan(&settlementBalance).Error
	if err != nil {
		return balance, err
	}

	return balance + settlementBalance, nil
}

// CalculateAllUserBalances calculates balances with all users for a given user
func (r *balanceRepository) CalculateAllUserBalances(userID uuid.UUID) ([]dto.BalanceItem, error) {
	// Complex query to calculate balances with all users
	var results []BalanceResult
	query := `
		WITH user_expenses AS (
			SELECT 
				CASE 
					WHEN ep1.user_id = ? THEN ep2.user_id 
					ELSE ep1.user_id 
				END as other_user_id,
				SUM(
					CASE 
						WHEN ep1.user_id = ? THEN ep1.paid_amount - ep1.owed_amount
						ELSE ep2.paid_amount - ep2.owed_amount
					END
				) as expense_balance
			FROM expense_participants ep1
			JOIN expense_participants ep2 ON ep1.expense_id = ep2.expense_id
			WHERE (ep1.user_id = ? OR ep2.user_id = ?)
			AND ep1.user_id != ep2.user_id
			GROUP BY other_user_id
		),
		user_settlements AS (
			SELECT 
				CASE 
					WHEN payer_id = ? THEN payee_id 
					ELSE payer_id 
				END as other_user_id,
				SUM(
					CASE 
						WHEN payer_id = ? THEN -amount 
						ELSE amount 
					END
				) as settlement_balance
			FROM settlements
			WHERE (payer_id = ? OR payee_id = ?)
			AND is_confirmed = true
			AND deleted_at IS NULL
			GROUP BY other_user_id
		)
		SELECT 
			u.id as user_id,
			u.first_name,
			u.last_name,
			u.email,
			COALESCE(ue.expense_balance, 0) + COALESCE(us.settlement_balance, 0) as balance
		FROM users u
		LEFT JOIN user_expenses ue ON u.id = ue.other_user_id
		LEFT JOIN user_settlements us ON u.id = us.other_user_id
		WHERE (ue.expense_balance IS NOT NULL OR us.settlement_balance IS NOT NULL)
		AND (COALESCE(ue.expense_balance, 0) + COALESCE(us.settlement_balance, 0)) != 0
		ORDER BY ABS(COALESCE(ue.expense_balance, 0) + COALESCE(us.settlement_balance, 0)) DESC
	`
	err := r.db.Raw(query, userID, userID, userID, userID, userID, userID, userID, userID).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	balances := make([]dto.BalanceItem, len(results))
	for i, result := range results {
		balances[i] = dto.BalanceItem{
			UserID:    result.UserID,
			UserName:  result.FirstName + " " + result.LastName,
			UserEmail: result.Email,
			Amount:    result.Balance,
		}
	}

	return balances, nil
}

// CalculateGroupBalances calculates balances for all members in a group
func (r *balanceRepository) CalculateGroupBalances(groupID uuid.UUID) ([]dto.GroupBalanceItem, error) {
	var results []GroupBalanceResult
	query := `
		SELECT 
			u.id as user_id,
			u.first_name,
			u.last_name,
			COALESCE(SUM(ep.paid_amount), 0) as total_paid,
			COALESCE(SUM(ep.owed_amount), 0) as total_owed
		FROM users u
		JOIN group_members gm ON u.id = gm.user_id
		LEFT JOIN expense_participants ep ON u.id = ep.user_id
		LEFT JOIN expenses e ON ep.expense_id = e.id AND e.group_id = ?
		WHERE gm.group_id = ?
		AND gm.left_at IS NULL
		AND (e.deleted_at IS NULL OR e.deleted_at IS NULL)
		GROUP BY u.id, u.first_name, u.last_name
		ORDER BY (COALESCE(SUM(ep.paid_amount), 0) - COALESCE(SUM(ep.owed_amount), 0)) DESC
	`
	err := r.db.Raw(query, groupID, groupID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	balances := make([]dto.GroupBalanceItem, len(results))
	for i, result := range results {
		balances[i] = dto.GroupBalanceItem{
			UserID:     result.UserID,
			UserName:   result.FirstName + " " + result.LastName,
			TotalPaid:  result.TotalPaid,
			TotalOwed:  result.TotalOwed,
			NetBalance: result.TotalPaid - result.TotalOwed,
		}
	}

	return balances, nil
}
