package repository

import (
	"context"
	"log"

	pg "github.com/jackc/pgx/v5"
	"github.com/yakupovdev/FoodStore/internal/model"
)

type OrdersRepo struct {
	Conn *pg.Conn
}

func NewOrdersRepo(conn *pg.Conn) *OrdersRepo {
	return &OrdersRepo{
		Conn: conn,
	}
}

func (or *OrdersRepo) GetProfileByID(clientID int64) (model.Client, error) {
	ctx := context.Background()
	stmt := `SELECT client_id, name, rating FROM clients WHERE client_id=$1`
	row := or.Conn.QueryRow(ctx, stmt, clientID)

	var client model.Client
	if err := row.Scan(&client.ID, &client.Name, &client.Rating); err != nil {
		log.Println(err)
		return model.Client{}, ErrGetProfile
	}
	stmt = `SELECT email,balance  FROM users WHERE userid=$1`
	row1 := or.Conn.QueryRow(ctx, stmt, clientID)
	if err := row1.Scan(&client.Email, &client.Balance); err != nil {
		log.Println(err)
		return model.Client{}, ErrGetProfile
	}

	return client, nil
}

func (or *OrdersRepo) GetOrdersByClientID(clientID int64) ([]model.ClientOrders, error) {
	ctx := context.Background()
	stmt := `SELECT order_id, client_id, status, created_at FROM orders WHERE client_id=$1`

	rows, err := or.Conn.Query(ctx, stmt, clientID)
	if err != nil {
		log.Println(err)
		return nil, ErrGetOrders
	}
	defer rows.Close()

	var orders []model.ClientOrders

	for rows.Next() {
		var order model.ClientOrders
		if err := rows.Scan(&order.OrderID, &order.ClientID, &order.Status, &order.CreatedAt); err != nil {
			log.Println(err)
			return nil, ErrGetOrders
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, ErrGetOrders
	}

	return orders, nil
}

func (or *OrdersRepo) GetOrderItemsByOrderID(orderID int64) ([]model.ClientOrdersItems, error) {
	ctx := context.Background()
	stmt := `SELECT order_item_id, order_id, seller_id, product_id, quantity, price_at_purchase FROM orders_items WHERE order_id=$1`

	rows, err := or.Conn.Query(ctx, stmt, orderID)
	if err != nil {
		log.Println(err)
		return nil, ErrGetOrderItems
	}
	defer rows.Close()

	var orderItems []model.ClientOrdersItems

	for rows.Next() {
		var item model.ClientOrdersItems
		if err := rows.Scan(&item.OrderItemsId, &item.OrderID, &item.SellerID, &item.ProductID, &item.Quantity, &item.PriceAtPurchase); err != nil {
			log.Println(err)
			return nil, ErrGetOrderItems
		}
		orderItems = append(orderItems, item)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, ErrGetOrderItems
	}

	return orderItems, nil
}

func (or *OrdersRepo) GetCategories() ([]model.Category, error) {
	ctx := context.Background()
	stmt := `SELECT categories_id, name FROM categories WHERE parent_id IS NULL`
	rows, err := or.Conn.Query(ctx, stmt)
	if err != nil {
		log.Println(err)
		return nil, ErrGetCategories
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Println(err)
			return nil, ErrGetCategories
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, ErrGetCategories
	}

	return categories, nil
}

func (or *OrdersRepo) GetSubCategoriesByCategoryID(categoryID int64) ([]model.SubCategory, error) {
	ctx := context.Background()
	stmt := `SELECT categories_id, name FROM categories WHERE parent_id=$1`
	rows, err := or.Conn.Query(ctx, stmt, categoryID)
	if err != nil {
		log.Println(err)
		return nil, ErrGetCategories
	}
	defer rows.Close()

	var subcategories []model.SubCategory
	for rows.Next() {
		var subcategory model.SubCategory
		if err := rows.Scan(&subcategory.ID, &subcategory.Name); err != nil {
			log.Println(err)
			return nil, ErrGetCategories
		}
		subcategories = append(subcategories, subcategory)
	}

	return subcategories, nil
}

func (or *OrdersRepo) GetProductsBySubCategoryID(subCategoryID int64) ([]model.Product, error) {
	ctx := context.Background()
	stmt := `SELECT product_id, name, description, img FROM products WHERE categories_id=$1`
	rows, err := or.Conn.Query(ctx, stmt, subCategoryID)
	if err != nil {
		log.Println(err)
		return nil, ErrGetProducts
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image); err != nil {
			log.Println(err)
			return nil, ErrGetProducts
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, ErrGetProducts
	}

	return products, nil
}
