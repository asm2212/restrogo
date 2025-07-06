# restrogo

a Restaurant Management System built with Go, Gin, and MongoDB.  
It is designed to manage users, food items, menus, tables, orders, order items, and invoices for a restaurant.

---

## Features

- **User authentication & management**
- **Food item CRUD**
- **Menu management**
- **Table booking and management**
- **Order and order item handling**
- **Invoice generation**
- **Built-in middleware authentication**

---

## Tech Stack

- [Go](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [MongoDB](https://www.mongodb.com/)
- [mongo-go-driver](https://github.com/mongodb/mongo-go-driver)

---

## API Endpoints

Below are the primary REST API endpoints as defined in the `routes` package.

### User
| Method | Endpoint               | Description          |
|--------|------------------------|----------------------|
| GET    | `/users`               | Get all users        |
| GET    | `/users/:user_id`      | Get single user      |
| POST   | `/users/signup`        | User signup          |
| POST   | `/users/login`         | User login           |

### Food
| Method | Endpoint                | Description          |
|--------|-------------------------|----------------------|
| GET    | `/foods`                | Get all foods        |
| GET    | `/foods/:food_id`       | Get single food      |
| POST   | `/foods`                | Create food          |
| PATCH  | `/foods/:food_id`       | Update food          |

### Menu
| Method | Endpoint                | Description          |
|--------|-------------------------|----------------------|
| GET    | `/menus`                | Get all menus        |
| GET    | `/menus/:menu_id`       | Get single menu      |
| POST   | `/menus`                | Create menu          |
| PATCH  | `/menus/:menu_id`       | Update menu          |

### Table
| Method | Endpoint                   | Description            |
|--------|----------------------------|------------------------|
| GET    | `/tables`                  | Get all tables         |
| GET    | `/tables/:table_id`        | Get single table       |
| POST   | `/tables`                  | Create table           |
| PATCH  | `/tables/:table_id`        | Update table           |

### Order
| Method | Endpoint                   | Description            |
|--------|----------------------------|------------------------|
| GET    | `/orders`                  | Get all orders         |
| GET    | `/orders/:order_id`        | Get single order       |
| POST   | `/orders`                  | Create order           |
| PATCH  | `/orders/:order_id`        | Update order           |

### OrderItem
| Method | Endpoint                            | Description                   |
|--------|-------------------------------------|-------------------------------|
| GET    | `/orderItems`                       | Get all order items           |
| GET    | `/orderItems/:orderItem_id`         | Get single order item         |
| GET    | `/orderItems-order/:order_id`       | Get order items by order      |
| POST   | `/orderItems`                       | Create order item             |
| PATCH  | `/orderItems/:orderItem_id`         | Update order item             |

### Invoice
| Method | Endpoint                    | Description            |
|--------|-----------------------------|------------------------|
| GET    | `/invoices`                 | Get all invoices       |
| GET    | `/invoices/:invoice_id`     | Get single invoice     |
| POST   | `/invoices`                 | Create invoice         |
| PATCH  | `/invoices/:invoice_id`     | Update invoice         |


---


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

---

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss.

---

## License

[MIT](LICENSE)
