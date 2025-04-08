package productService

import (
	"context"
	"database/sql"
	"strconv"

	"grpcShop.com/backend/apis/product"
)

type ProductServer struct {
	product.UnimplementedProductServiceServer
	// for pourpose of learning about grpc and go. I don't implement repository pattern here
	Db *sql.DB
}

// Implement the methods of the ProductServiceServer interface here
func (productServer *ProductServer) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.ProductResponse, error) {
	// Implement the logic to get a product from the database

	var id int
	var name, description string
	var price float64

	productServer.Db.QueryRow("SELECT * FROM products WHERE id = $1", req.Id).Scan(&id, &name, &description, &price)

	// This is just a placeholder implementation
	return &product.ProductResponse{
		Id:          req.Id,
		Name:        name,
		Description: description,
		Price:       price,
	}, nil
}

func (productServer *ProductServer) ListProducts(ctx context.Context, _ *product.Empty) (*product.ListProductResponse, error) {
	// Implement the logic to list products from the database
	rows, err := productServer.Db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*product.ProductResponse
	for rows.Next() {
		var id string
		var name, description string
		var price float64

		err := rows.Scan(&id, &name, &description, &price)
		if err != nil {
			return nil, err
		}

		product := &product.ProductResponse{
			Id:          id,
			Name:        name,
			Description: description,
			Price:       price,
		}
		products = append(products, product)
	}

	return &product.ListProductResponse{
		Products: products,
	}, nil
}

func (productServer *ProductServer) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.ProductResponse, error) {

	var id int
	err := productServer.Db.QueryRow("INSERT INTO products (name,description,price) VALUES ($1,$2,$3) RETURNING id", req.Name, req.Description, req.Price).Scan(&id)
	if err != nil {
		return nil, err
	}

	idString := strconv.Itoa(id)

	return &product.ProductResponse{
		Id:          idString,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}, nil

}
