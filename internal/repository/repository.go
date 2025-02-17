package repository

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"avito-shop-service/internal/model"
)

const (
	defaultGiftBalance = 1000
)

const (
	sysUserRegistrationGift = "registration_gift"
	sysUserBuyItem          = "buy_item"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetUser(ctx context.Context, username string) (model.User, error) {
	query := `
	select username, password_hash
	from "user"
	where username = $1`

	rows, err := r.pool.Query(ctx, query, username)
	if err != nil {
		return model.User{}, errors.WithStack(err)
	}

	row, err := pgx.CollectExactlyOneRow[userRow](rows, pgx.RowToStructByNameLax[userRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, errors.WithStack(err)
	}

	return r.userModel(row), nil
}

func (r *Repository) CreateUser(ctx context.Context, username string, passwordHash string) (model.User, error) {
	query := `
	with u as (select username, password_hash
	           from "user"
	           where username = $1 for update),
	     uu as (insert into "user" (username, password_hash)
	            select $1, $2
	            where not exists (select from u)
	            returning username, password_hash),
	     b as (insert into balance (username, amount)
	           select $1, $3
	           where not exists(select from u)),
	     t as (insert into transaction ("from", "to", amount)
	           values ($4, $1, $3))
	select username, password_hash
	from u
	union
	select username, password_hash
	from uu`

	rows, err := r.pool.Query(ctx, query, username, passwordHash, defaultGiftBalance, sysUserRegistrationGift)
	if err != nil {
		return model.User{}, errors.WithStack(err)
	}

	row, err := pgx.CollectExactlyOneRow[userRow](rows, pgx.RowToStructByNameLax[userRow])
	if err != nil {
		return model.User{}, errors.WithStack(err)
	}

	return r.userModel(row), nil
}

func (r *Repository) UpdateBalance(ctx context.Context, username string, itemName string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = r.updateBalance(ctx, tx, username, itemName)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repository) SwapBalance(ctx context.Context, fromUser, toUser string, amount int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = r.swapBalance(ctx, tx, fromUser, toUser, amount)
	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repository) Balance(ctx context.Context, username string) (int64, error) {
	query := `select amount from balance where username = $1`

	var balance int64

	err := r.pool.QueryRow(ctx, query, username).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return balance, model.ErrUserNotFound
		}
		return balance, errors.WithStack(err)
	}

	return balance, nil
}

func (r *Repository) Inventory(ctx context.Context, username string) ([]model.Inventory, error) {
	query := `
	select item_name, quantity from inventory
	where username = $1`

	rows, err := r.pool.Query(ctx, query, username)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	inventoryRows, err := pgx.CollectRows[inventoryRow](rows, pgx.RowToStructByNameLax[inventoryRow])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	inventory := make([]model.Inventory, 0, len(inventoryRows))
	for _, row := range inventoryRows {
		inventory = append(inventory, r.inventoryModel(row))
	}

	return inventory, nil
}

func (r *Repository) Transaction(ctx context.Context, username string) ([]model.Transaction, error) {
	query := `
	select "from", "to", amount
	from transaction
	where $1 in ("from", "to")`

	rows, err := r.pool.Query(ctx, query, username)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	transactionRows, err := pgx.CollectRows[transactionRow](rows, pgx.RowToStructByNameLax[transactionRow])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	transactions := make([]model.Transaction, 0, len(transactionRows))
	for _, row := range transactionRows {
		transactions = append(transactions, r.transactionModel(row))
	}

	return transactions, nil
}

func (r *Repository) swapBalance(ctx context.Context, tx pgx.Tx, fromUser, toUser string, amount int) error {
	query := `
	with t as (insert into transaction ("from", "to", amount)
	           values ($1, $2, $3))
	update balance b
	set amount = u.amount
	from (select username,
	             case when username = $1 then amount - $3 else amount + $3 end amount
	      from balance
	      where username in ($1, $2) for update) u
	where b.username = u.username
	returning b.username, b.amount`

	rows, err := tx.Query(ctx, query, fromUser, toUser, amount)
	if err != nil {
		return errors.WithStack(err)
	}

	balanceRows, err := pgx.CollectRows[balanceRow](rows, pgx.RowToStructByNameLax[balanceRow])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "transaction_to_fkey" {
				return model.ErrUserNotFound
			}
		}
		return err
	}

	for _, row := range balanceRows {
		if row.Username == fromUser {
			if row.Amount < 0 {
				return model.ErrNegativeBalance
			}
		}
	}

	return nil

}

func (r *Repository) updateBalance(ctx context.Context, tx pgx.Tx, username string, itemName string) error {
	query := `
	with i as (insert into inventory (username, item_name)
	           values ($1, $2)
	           on conflict (username, item_name) do update
	           set quantity = inventory.quantity + 1),
	     p as (select price
	           from item
	           where name = $2),
	     t as (insert into transaction ("from", "to", amount)
	           values ($1, $3, coalesce((select price from p), 0)))
	update balance b
	set amount = b.amount - coalesce((select price from p), 0)
	from (select username, amount
	      from balance
	      where username = $1 for update) u
	where b.username = u.username
	returning b.amount`

	var balance int64

	err := tx.QueryRow(ctx, query, username, itemName, sysUserBuyItem).Scan(&balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "inventory_item_name_fkey" {
				return model.ErrItemNotFound
			}
		}
		return err
	}

	if balance < 0 {
		return model.ErrNegativeBalance
	}

	return nil
}

func (r *Repository) userModel(row userRow) model.User {
	return model.User{
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
	}
}

func (r *Repository) itemModel(row itemRow) model.Item {
	return model.Item{
		Name:  row.Name,
		Price: row.Price,
	}
}

func (r *Repository) balanceModel(row balanceRow) model.Balance {
	return model.Balance{
		Username: row.Username,
		Amount:   row.Amount,
	}
}

func (r *Repository) inventoryModel(row inventoryRow) model.Inventory {
	return model.Inventory{
		ItemName: row.ItemName,
		Quantity: row.Quantity,
	}
}

func (r *Repository) transactionModel(row transactionRow) model.Transaction {
	return model.Transaction{
		FromUser: row.FromUser,
		ToUser:   row.ToUser,
		Amount:   row.Amount,
	}
}

type userRow struct {
	Username          string `db:"username"`
	PasswordHash      string `db:"password_hash"`
	PasswordHashMatch bool   `db:"ok"`
}

type itemRow struct {
	Name  string `db:"name"`
	Price int64  `db:"price"`
}

type balanceRow struct {
	Username string `db:"username"`
	Amount   int64  `db:"amount"`
}

type inventoryRow struct {
	ItemName string `db:"item_name"`
	Quantity int64  `db:"quantity"`
}

type transactionRow struct {
	FromUser string `db:"from"`
	ToUser   string `db:"to"`
	Amount   int64  `db:"amount"`
}
