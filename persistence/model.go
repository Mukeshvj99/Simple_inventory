package persistence

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	//"github.com/jackc/pgx/v4"
)

type Product struct {
	ID    int64   `json:id`
	Name  string  `json:name`
	Price float64 `json:price`
}

var (
	getallquery string
	getquery    string
	deletequery string
	updatequery string
	createquery string
)

func InitializeQueries() {
	getquery = `select * from ` + table + ` where id=$1;`
	getallquery = `select * from ` + table + ` limit $1 offset $2;`
	deletequery = `delete from ` + table + ` products where id=$1;`
	updatequery = `update ` + table + ` set name=$1,price=$2 where id=$3`
	createquery = `insert into ` + table + ` values($1,$2,$3);`
}

func (p *Product) GetProductData(id int64) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()

	}()

	row, err := dbconn.db.Query(ctx, getquery, id)

	if err != nil {
		log.Println("error occurred in fetching the rows", err, id)
		return errors.New("Error in fetching rows")
	}

	defer row.Close()
	var err1 error

	if row != nil && row.Next() {
		err1 = row.Scan(&p.ID, &p.Name, &p.Price)
	} else {
		return fmt.Errorf("invalid Product id ")
	}

	return err1

}

func GetAllProductsData(count, start int) ([]Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()

	}()
	fmt.Println("tablname", table)

	row, err := dbconn.db.Query(ctx, getallquery, count, start)

	if err != nil {
		log.Println("error occurred in fetching the rows", err)
		return nil, errors.New("Error in fetching rows")
	}

	defer row.Close()
	var err1 error
	var products []Product
	if row != nil && row.Next() {
		var p Product
		err1 = row.Scan(&p.ID, &p.Name, &p.Price)
		products = append(products, p)
		for row.Next() {
			var p Product
			err1 = row.Scan(&p.ID, &p.Name, &p.Price)
			products = append(products, p)
		}
	} else {
		return nil, nil
	}

	return products, err1
}

func DeleteProductData(id int) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()

	}()

	cmd, err := dbconn.db.Exec(ctx, deletequery, id)

	if err != nil {
		log.Println("error occurred in deleting the  rows", err, id)
		return 0, errors.New("Error in deleting rows")
	}

	return int(cmd.RowsAffected()), err
}

func UpdataProductData(id int, data Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()

	}()

	cmd, err := dbconn.db.Exec(ctx, updatequery, data.Name, data.Price, id)

	if err != nil {
		log.Println("error occurred in Updating the  rows", err, id)
		return 0, errors.New("Error in Updating rows")
	}

	return int(cmd.RowsAffected()), err
}

func CreateProductData(products []Product) (int, []int) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()

	}()

	var failedproducts []int
	count := 0
	for _, p := range products {
		_, err := dbconn.db.Exec(ctx, createquery, p.ID, p.Name, p.Price)

		if err != nil {
			log.Println("error occurred in Creating the Product", err, "product id-", p.ID)
			failedproducts = append(failedproducts, int(p.ID))
		} else {
			count++
		}

	}

	return count, failedproducts
}
