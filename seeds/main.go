package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

func createSchema(db *sql.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL CHECK (name <> ''),
        email VARCHAR(100) UNIQUE NOT NULL CHECK (email <> ''),
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS products (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL CHECK (name <> ''),
        price NUMERIC(10, 2) NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS orders (
        id SERIAL PRIMARY KEY,
        user_id INT REFERENCES users(id) ON DELETE CASCADE,
        product_id INT REFERENCES products(id) ON DELETE CASCADE,
        quantity INT NOT NULL CHECK (quantity > 0),
        created_at TIMESTAMPTZ DEFAULT NOW()
    );
    `
	_, err := db.Exec(schema)
	return err
}

func addUser(db *sql.DB) {
	var name, email string

	fmt.Print("Enter a name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter an email: ")
	fmt.Scanln(&email)

	fmt.Println("Inserting user into database")
	_, err := db.Exec(
		"INSERT INTO users (name, email) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING",
		name,
		email,
	)
	if err != nil {
		fmt.Printf("Failed to add user: %v\n", err)
	} else {
		fmt.Println("User added successfully!")
	}
}

func addProduct(db *sql.DB) {
	var name string
	var price float64

	fmt.Print("Enter a product name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter a product price: ")
	fmt.Scanln(&price)

	_, err := db.Exec(
		"INSERT INTO products (name, price) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		name,
		price,
	)
	if err != nil {
		fmt.Printf("Failed to add product: %v\n", err)
	} else {
		fmt.Println("Product added successfully!")
	}
}

func addOrder(db *sql.DB) {
	var userID, productID, quantity int

	fmt.Print("Enter a user id: ")
	fmt.Scanln(&userID)

	fmt.Print("Enter a product id: ")
	fmt.Scanln(&productID)

	fmt.Print("Enter a quantity: ")
	fmt.Scanln(&quantity)

	_, err := db.Exec(
		"INSERT INTO orders (user_id, product_id, quantity) VALUES ($1, $2, $3)",
		userID,
		productID,
		quantity,
	)
	if err != nil {
		fmt.Printf("Failed to add order: %v\n", err)
	} else {
		fmt.Println("Order added successfully!")
	}
}

func updateUser(db *sql.DB) {
	var userID int
	var name, email string

	fmt.Print("Enter user ID to update: ")
	fmt.Scanln(&userID)

	fmt.Print("Enter new name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter new email: ")
	fmt.Scanln(&email)

	_, err := db.Exec(
		"UPDATE users SET name = $1, email = $2 WHERE id = $3",
		name,
		email,
		userID,
	)
	if err != nil {
		fmt.Printf("Failed to update user: %v\n", err)
	} else {
		fmt.Println("User updated successfully!")
	}
}

func updateProduct(db *sql.DB) {
	var productID int
	var name string
	var price float64

	fmt.Print("Enter product ID to update: ")
	fmt.Scanln(&productID)

	fmt.Print("Enter new product name: ")
	fmt.Scanln(&name)

	fmt.Print("Enter new product price: ")
	fmt.Scanln(&price)

	_, err := db.Exec(
		"UPDATE products SET name = $1, price = $2 WHERE id = $3",
		name,
		price,
		productID,
	)
	if err != nil {
		fmt.Printf("Failed to update product: %v\n", err)
	} else {
		fmt.Println("Product updated successfully!")
	}
}

func updateOrder(db *sql.DB) {
	var orderID, userID, productID, quantity int

	fmt.Print("Enter order ID to update: ")
	fmt.Scanln(&orderID)

	fmt.Print("Enter new user ID: ")
	fmt.Scanln(&userID)

	fmt.Print("Enter new product ID: ")
	fmt.Scanln(&productID)

	fmt.Print("Enter new quantity: ")
	fmt.Scanln(&quantity)

	_, err := db.Exec(
		"UPDATE orders SET user_id = $1, product_id = $2, quantity = $3 WHERE id = $4",
		userID,
		productID,
		quantity,
		orderID,
	)
	if err != nil {
		fmt.Printf("Failed to update order: %v\n", err)
	} else {
		fmt.Println("Order updated successfully!")
	}
}

func deleteUser(db *sql.DB) {
	var userID int

	fmt.Print("Enter user ID to delete: ")
	fmt.Scanln(&userID)

	_, err := db.Exec(
		"DELETE FROM users WHERE id = $1",
		userID,
	)
	if err != nil {
		fmt.Printf("Failed to delete user: %v\n", err)
	} else {
		fmt.Println("User deleted successfully!")
	}
}

func deleteProduct(db *sql.DB) {
	var productID int

	fmt.Print("Enter product ID to delete: ")
	fmt.Scanln(&productID)

	_, err := db.Exec(
		"DELETE FROM products WHERE id = $1",
		productID,
	)
	if err != nil {
		fmt.Printf("Failed to delete product: %v\n", err)
	} else {
		fmt.Println("Product deleted successfully!")
	}
}

func deleteOrder(db *sql.DB) {
	var orderID int

	fmt.Print("Enter order ID to delete: ")
	fmt.Scanln(&orderID)

	_, err := db.Exec(
		"DELETE FROM orders WHERE id = $1",
		orderID,
	)
	if err != nil {
		fmt.Printf("Failed to delete order: %v\n", err)
	} else {
		fmt.Println("Order deleted successfully!")
	}
}

func main() {
	db, err := sql.Open("postgres", env.GetEnv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()
	log.Printf("Connected to database successfully")

	// Create tables
	if err := createSchema(db); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	fmt.Println("Database schema created and seeded successfully!")

	var choice int

	for {
		fmt.Println("Enter a choice")
		fmt.Println("1. Add a user")
		fmt.Println("2. Update a user")
		fmt.Println("3. Delete a user")
		fmt.Println("4. Add a product")
		fmt.Println("5. Update a product")
		fmt.Println("6. Delete a product")
		fmt.Println("7. Add an order")
		fmt.Println("8.Update an order")
		fmt.Println("9. Delete an order")
		fmt.Println("10. Exit")

		fmt.Printf("Enter a choice: ")
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("Invalid input, please enter a number.")
			continue
		}

		switch choice {
		case 1:
			addUser(db)
		case 2:
			updateUser(db)
		case 3:
			deleteUser(db)
		case 4:
			addProduct(db)
		case 5:
			updateProduct(db)
		case 6:
			deleteProduct(db)
		case 7:
			addOrder(db)
		case 8:
			updateOrder(db)
		case 9:
			deleteOrder(db)
		case 10:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
