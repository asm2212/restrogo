# restrogo

**restrogo** is a Restaurant Management System built with Go, Gin, and MongoDB.  
It is designed to manage users, food items, menus, tables, orders, order items, and invoices for a restaurant.

## Features

- User authentication & management
- Food item CRUD
- Menu management
- Table booking and management
- Order and order item handling
- Invoice generation
- Built-in middleware authentication

## Tech Stack

- [Go](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [MongoDB](https://www.mongodb.com/)
- [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)

## Getting Started

### Prerequisites

- Go 1.18+
- MongoDB instance (local or remote)

### Installation

1. **Clone the repo:**
    ```sh
    git clone https://github.com/asm2212/restrogo.git
    cd restrogo
    ```

2. **Install dependencies:**
    ```sh
    go mod tidy
    ```

3. **Set environment variables:**
    - `PORT` (optional, default is 8000)
    - MongoDB connection (see your `database` package for expected connection string)

4. **Run the server:**
    ```sh
    go run main.go
    ```

### API Endpoints

- **User:** `/users`  
- **Food:** `/foods`
- **Menu:** `/menus`
- **Table:** `/tables`
- **Order:** `/orders`
- **OrderItem:** `/orderitems`
- **Invoice:** `/invoices`

See the `routes/` directory for detailed route handlers.

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss.

## License

[MIT](LICENSE)
