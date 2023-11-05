package config

import (
	"context"
	"fmt"
	"go-restfull-api-inventory/entity"
	"testing"
)

func TestExecSql(t *testing.T) {
	db := ConnectDB()
	defer db.Close()

	ctx := context.Background()

	query := "INSERT INTO Customer(id, name, phoneNumber) VALUES ('C100', 'Rizki', '034234');"

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}
}

func TestQuerySql(t *testing.T) {
	db := ConnectDB()
	defer db.Close()

	ctx := context.Background()

	query := "SELECT * FROM Customer;"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		customers := entity.Customer{}
		err := rows.Scan(&customers.Id, &customers.Name, &customers.Phone, &customers.Address)

		if err != nil {
			panic(err)
		}
		fmt.Println(customers)
	}
}

func TestSqlInjection(t *testing.T) {
	db := ConnectDB()
	defer db.Close()

	ctx := context.Background()

	id := "C1"

	query := "SELECT * FROM Customer WHERE id = $1;"
	var cust entity.Customer

	err := db.QueryRowContext(ctx, query, id).Scan(&cust.Id, &cust.Name, &cust.Phone, &cust.Address)
	if err != nil {
		panic(err)
	}
	fmt.Println(cust)

}
