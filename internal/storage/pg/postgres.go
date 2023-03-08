package pg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/models"

	_ "github.com/lib/pq"
)

type PgStorage struct {
	connect *sql.DB
}

func New(uri string) (*PgStorage, error) {
	connect, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return &PgStorage{connect: connect}, nil
}

func (PS *PgStorage) Ping() error {
	if err := PS.connect.Ping(); err != nil {
		return err
	}
	err := createTables(PS.connect)
	if err != nil {
		return err
	}
	return nil
}

func (PS *PgStorage) Close() error {
	if err := PS.connect.Close(); err != nil {
		return err
	}
	return nil
}

func createTables(connect *sql.DB) error {
	_, err := connect.Exec(`
	create table if not exists public.users(
		login_user text primary key,
		password_user text,
		create_user timestamp default now()
	);
	
	create table if not exists public.orders(
		 number_order text primary key,
		 login_user text,
		 status_order varchar(50),
		 accrual_order double precision,
		 uploaded_order timestamp default now(),
		 created_order timestamp default now(),
		 foreign key (login_user) references public.users (login_user)
	);
	
	create table if not exists public.withdraws(
		 login_user text,
		 number_order text,
		 sum double precision,
		 uploaded_order timestamp default now()
	);
	`)
	if err != nil {
		return err
	}
	return nil
}

func (PS *PgStorage) Adduser(ctx context.Context, user models.User) error {
	result, err := PS.connect.ExecContext(ctx,
		`insert into public.users (login_user, password_user) values ($1, $2) on conflict do nothing`,
		user.Login, user.Password)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return consta.ErrorNoUNIQUE
	}
	return nil
}

func (PS *PgStorage) Authentication(ctx context.Context, user models.User) (bool, error) {
	var done int
	err := PS.connect.QueryRowContext(ctx, `select count(1) from public.users where login_user=$1 and password_user=$2`,
		user.Login, user.Password).Scan(&done)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if done == 0 {
		return false, nil
	}
	return true, nil
}

func (PS *PgStorage) GetOrder(ctx context.Context, numberOrder string) (models.Order, error) {
	var order models.Order
	err := PS.connect.QueryRowContext(ctx, `select number_order, login_user, status_order, accrual_order, uploaded_order
	from public.orders where number_order=$1 order by created_order desc`,
		numberOrder).Scan(&order.NumberOrder, &order.UserLogin, &order.Status,
		&order.Accrual, &order.Uploaded)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (PS *PgStorage) GetOrders(ctx context.Context, userLogin string) ([]models.Order, error) {
	var orders []models.Order
	rows, err := PS.connect.QueryContext(ctx, `select number_order, login_user, status_order, accrual_order, uploaded_order
	from public.orders where login_user=$1`, userLogin)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.NumberOrder, &order.UserLogin, &order.Status, &order.Accrual, &order.Uploaded)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (PS *PgStorage) AddOrder(ctx context.Context, numberOrder string, order models.Order) error {
	result, err := PS.connect.ExecContext(ctx, `insert into public.orders 
    (number_order, login_user, status_order, uploaded_order, accrual_order)  values ($1, $2, $3, $4, $5) on conflict do nothing`,
		numberOrder, order.UserLogin, order.Status, order.Uploaded, order.Accrual)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return consta.ErrorNoUNIQUE
	}
	return nil
}

func (PS *PgStorage) GetOrdersProcess(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	sliceStatus := []interface{}{consta.OrderStatusPROCESSING, consta.OrderStatusNEW, consta.OrderStatusREGISTERED, consta.OrderStatusInvalid}
	rows, err := PS.connect.QueryContext(ctx, `select number_order, login_user, status_order, accrual_order, uploaded_order
	from public.orders where status_order in ($1, $2, $3,$4)`, sliceStatus...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.NumberOrder, &order.UserLogin, &order.Status, &order.Accrual, &order.Uploaded)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (PS *PgStorage) UpdateOrder(ctx context.Context, loyality models.Loyalty) error {
	_, err := PS.connect.ExecContext(ctx, `update public.orders set accrual_order=$1, status_order=$2,
                         uploaded_order=now() where number_order =$3`,
		loyality.Accrual, loyality.Status, loyality.NumberOrder)
	if err != nil {
		return err
	}
	return nil
}

func (PS *PgStorage) GetUserBalance(ctx context.Context, userLogin string) (float64, float64, error) {
	var ordersSUM float64
	var withdrawsSUM float64
	err := PS.connect.QueryRowContext(ctx, `select (case when sum_order is null then 0.0 else sum_order end) as sum_order, (case when sum_withdraws is null then 0.0 else sum_withdraws end) as sum_withdraws from
	 (select sum(accrual_order) as  sum_order from public.orders where login_user = $1) as orders,
	 (select sum(sum) as  sum_withdraws from public.withdraws where login_user = $1) as withdraws`, userLogin).
		Scan(&ordersSUM, &withdrawsSUM)
	return ordersSUM, withdrawsSUM, err
}

func (PS *PgStorage) AddWithdraw(ctx context.Context, withdraw models.Withdraw) error {
	result, err := PS.connect.ExecContext(ctx, `

	insert into public.withdraws (login_user, number_order, sum, uploaded_order)
	select $1, $2, $3, $4
	where (
          select sum_order >= sum_withdraws + $3 from (
          select (case when sum_order is null then 0 else sum_order end ) as sum_order,
          (case when sum_withdraws is null then 0 else sum_withdraws end ) as sum_withdraws from
          (select sum(accrual_order) as  sum_order from public.orders where login_user = $1) as orders,
          (select sum(sum) as  sum_withdraws from public.withdraws where login_user = $1) as withdraws) as s
          );
	`, withdraw.UserLogin, withdraw.NumberOrder, withdraw.Sum, withdraw.ProcessedAT)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return consta.ErrorStatusShortfallAccount
	}
	return nil
}

func (PS *PgStorage) GetWithdraws(ctx context.Context, userLogin string) ([]models.Withdraw, error) {
	var withdraws []models.Withdraw
	rows, err := PS.connect.QueryContext(ctx, `select login_user, number_order, sum, uploaded_order from public.withdraws
	where login_user = $1
	order by uploaded_order`, userLogin)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	for rows.Next() {
		var withdraw models.Withdraw
		err = rows.Scan(&withdraw.UserLogin, &withdraw.NumberOrder, &withdraw.Sum, &withdraw.ProcessedAT)
		if err != nil {
			return nil, err
		}
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, rows.Err()
}
