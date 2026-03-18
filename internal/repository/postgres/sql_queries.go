package postgres

import (
	"context"

	"github.com/klyakssa/go-diplom.git/internal/domain/auth"
	balancePkg "github.com/klyakssa/go-diplom.git/internal/domain/balance"
	"github.com/klyakssa/go-diplom.git/internal/domain/orders"
	"github.com/shopspring/decimal"
)

func (p *PostgresStorage) CreateUser(ctx context.Context, login, password string) (string, error) {
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`
	var userid string
	err := p.DB.QueryRowContext(ctx, query, login, password).Scan(&userid)
	return userid, err
}

func (p *PostgresStorage) GetUserByLogin(ctx context.Context, login string) (*auth.User, error) {
	query := `SELECT id, login, password FROM users WHERE login = $1`
	row := p.DB.QueryRowxContext(ctx, query, login)
	user := &auth.User{}
	err := row.StructScan(user)
	return user, err
}

func (p *PostgresStorage) GetPendingOrders(ctx context.Context) ([]orders.Order, error) {
	rows, err := p.DB.QueryxContext(ctx, `
		SELECT number, user_id, status, accrual, is_accrualed
		FROM orders
		WHERE status IN ('NEW', 'PROCESSING')
		   OR (status = 'PROCESSED' AND is_accrualed = FALSE)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordersArr []orders.Order

	for rows.Next() {
		var o orders.Order
		err := rows.StructScan(&o)
		if err != nil {
			return nil, err
		}
		ordersArr = append(ordersArr, o)
	}

	return ordersArr, nil
}

func (p *PostgresStorage) ApplyAccrual(ctx context.Context, order *orders.Order) error {
	tx, err := p.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var applied bool
	err = tx.QueryRowContext(ctx,
		`SELECT is_accrualed FROM orders WHERE number = $1 FOR UPDATE`,
		order.Number,
	).Scan(&applied)
	if err != nil {
		return err
	}

	if applied {
		return nil
	}

	_, err = tx.NamedExecContext(ctx,
		`UPDATE balances SET current = current + :balance WHERE id = :id`,
		map[string]interface{}{
			"balance": order.Accrual,
			"id":      order.UserID,
		},
	)
	if err != nil {
		return err
	}

	_, err = tx.NamedExecContext(ctx,
		`UPDATE orders SET is_accrualed = TRUE WHERE number = :number`,
		map[string]interface{}{
			"number": order.Number,
		},
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) UpdateOrder(ctx context.Context, number, status string, accrual decimal.Decimal) error {
	_, err := p.DB.NamedExecContext(ctx,
		`UPDATE orders SET status = :status, accrual = :accrual WHERE number = :number`,
		map[string]interface{}{
			"status":  status,
			"accrual": accrual,
			"number":  number,
		},
	)
	return err
}

func (p *PostgresStorage) CreateOrder(ctx context.Context, order *orders.Order) error {
	_, err := p.DB.NamedExecContext(ctx,
		`INSERT INTO orders (number, user_id) VALUES (:number, :user_id)`,
		map[string]interface{}{
			"number":  order.Number,
			"user_id": order.UserID,
		},
	)
	if err != nil {
		return err
	}

	_, err = p.DB.NamedExecContext(ctx,
		`INSERT INTO balances (user_id) VALUES (:user_id)`,
		map[string]interface{}{
			"user_id": order.UserID,
		},
	)
	return err
}

func (p *PostgresStorage) GetOrderByNumber(ctx context.Context, number string) (*orders.Order, error) {
	var o orders.Order
	err := p.DB.GetContext(ctx, &o, `SELECT * FROM orders WHERE number = $1`, number)
	return &o, err
}

func (p *PostgresStorage) GetOrdersByUserID(ctx context.Context, userID string) ([]orders.Order, error) {
	var o []orders.Order
	err := p.DB.SelectContext(ctx, &o, `SELECT * FROM orders WHERE user_id = $1 ORDER BY uploaded_at`, userID)
	return o, err
}

func (p *PostgresStorage) WithdrawBalance(ctx context.Context, userID string, orderNumber string, amount decimal.Decimal) error {
	tx, err := p.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance decimal.Decimal
	err = tx.QueryRowContext(ctx,
		`SELECT current FROM balances WHERE user_id = $1 FOR UPDATE`,
		userID,
	).Scan(&balance)
	if err != nil {
		return err
	}

	if balance.LessThan(amount) {
		return balancePkg.ErrInsufficientFunds
	}

	_, err = tx.NamedExecContext(ctx,
		`UPDATE balances SET current = current - :amount WHERE user_id = :user_id`,
		map[string]interface{}{
			"amount":  amount,
			"user_id": userID,
		},
	)
	if err != nil {
		return err
	}

	_, err = tx.NamedExecContext(ctx,
		`INSERT INTO withdraw_history (number, sum, user_id) VALUES (:number, :sum, :user_id)`,
		map[string]interface{}{
			"number":  orderNumber,
			"sum":     amount,
			"user_id": userID,
		},
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetBalanceWithdrawn(ctx context.Context, userID string) (int, decimal.Decimal, error) {
	tx, err := p.DB.BeginTxx(ctx, nil)
	if err != nil {
		return 0, decimal.Zero, err
	}
	defer tx.Rollback()

	var withdrawn decimal.Decimal
	err = tx.GetContext(ctx, &withdrawn, `SELECT COALESCE(SUM(sum), 0)::NUMERIC FROM withdraw_history WHERE user_id = $1`, userID)
	if err != nil {
		return 0, decimal.Zero, err
	}

	var balance int
	err = tx.GetContext(ctx, &balance, `SELECT current FROM balances WHERE user_id = $1`, userID)
	if err != nil {
		return 0, decimal.Zero, err
	}

	return balance, withdrawn, tx.Commit()
}

func (p *PostgresStorage) GetWithdrawls(ctx context.Context, userID string) ([]balancePkg.WithdrawHistory, error) {
	var wh []balancePkg.WithdrawHistory
	err := p.DB.SelectContext(ctx, &wh, `SELECT * FROM withdraw_history WHERE user_id = $1 ORDER BY processed_at`, userID)
	return wh, err
}
