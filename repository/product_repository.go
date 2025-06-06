package repository

import (
	"database/sql"
	"fmt"
	"product-go-api/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts(page, limit int, name string) ([]model.Product, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := "SELECT id, product_name, price FROM product"
	var args []interface{}
	argIdx := 1

	if name != "" {
		query += fmt.Sprintf(" WHERE product_name ILIKE $%d", argIdx)
		args = append(args, "%"+name+"%")
		argIdx++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := pr.connection.Query(query, args...)

	if err != nil {
		return []model.Product{}, err
	}

	var productList []model.Product
	var productObj model.Product

	for rows.Next() {
		if err := rows.Scan(&productObj.ID, &productObj.Name, &productObj.Price); err != nil {
			return []model.Product{}, err
		}
		productList = append(productList, productObj)
	}

	rows.Close()
	return productList, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {

	var id int
	query, err := pr.connection.Prepare(
		"INSERT INTO product" + "(product_name, price)" + "VALUES ($1, $2) RETURNING id;",
	)
	if err != nil {
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Price).Scan(&id)
	if err != nil {
		return 0, err
	}

	query.Close()
	return id, nil
}

func (pr *ProductRepository) GetProductById(id_product int) (*model.Product, error) {
	query, err := pr.connection.Prepare("SELECT * FROM product WHERE id = $1;")

	if err != nil {
		return nil, err
	}

	var product model.Product

	err = query.QueryRow(id_product).Scan(&product.ID, &product.Name, &product.Price)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	query.Close()
	return &product, nil
}

func (pr *ProductRepository) DeleteProduct(id_product int) error {
	query, err := pr.connection.Prepare("DELETE FROM product WHERE id = $1;")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(id_product)
	return err
}

func (pr *ProductRepository) UpdateProduct(product model.Product) (*model.Product, error) {
	query, err := pr.connection.Prepare(
		"UPDATE product SET product_name = $2, price = $3 WHERE id = $1 RETURNING id, product_name, price;",
	)
	if err != nil {
		return nil, err
	}

	var updatedProduct model.Product
	err = query.QueryRow(product.ID, product.Name, product.Price).Scan(
		&updatedProduct.ID,
		&updatedProduct.Name,
		&updatedProduct.Price,
	)

	if err != nil {
		return nil, err
	}

	query.Close()
	return &updatedProduct, nil
}
