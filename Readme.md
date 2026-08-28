# 🍽️ Restaurant Backend API

A robust RESTful API for a restaurant management system built with **Go**, **GORM**, **PostgreSQL**, and **JWT authentication**. This backend handles restaurants, menus, tables, reservations, orders, and reviews with fine‑grained role‑based access control.

## ✨ Features

- **JWT Authentication** (Access + Refresh tokens)
- **Role‑Based Access Control** (Admin, Owner, Staff, Customer)
- **Restaurant Management** (CRUD, ownership verification)
- **Menu Management** (CRUD, availability, pagination, filters)
- **Table Management** (CRUD, availability)
- **Reservation System** (create, confirm, cancel, conflict detection)
- **Order Management** (create, cancel, status updates)
- **Review System** (rate restaurants, one review per user per restaurant, only after completed order)
- **Secure Password Hashing** (bcrypt)
- **Structured Logging** & **Panic Recovery**
- **Graceful Shutdown**

---

## 🛠️ Tech Stack

| Category           | Technology                          |
|--------------------|-------------------------------------|
| Language           | Go 1.21+                            |
| Web Framework      | Net/http (standard library)         |
| ORM                | GORM v2                             |
| Database           | PostgreSQL 15+                      |
| Authentication     | JWT (golang-jwt/jwt/v5)             |
| Password Hashing   | bcrypt (golang.org/x/crypto/bcrypt) |
| Routing            | Standard library (Go 1.22+ patterns) |
| Environment        | godotenv                            |

---

## 📦 Setup Instructions

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Make (optional)

### 1. Clone the repository
```bash
git clone https://github.com/yourusername/restaurant-backend.git
cd restaurant-backend
```

### 2. Create a `.env` file
```env
DB_URL=postgres://user:password@localhost:5432/restaurant_db?sslmode=disable
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h   # 7 days
JWT_SIGNING_KEY=your-256-bit-secret-key
```

### 3. Install dependencies
```bash
go mod tidy
```

### 4. Run migrations (auto‑migrate is enabled)
The database schema will be created automatically on startup.

### 5. Start the server
```bash
go run cmd/main.go
```

The server will start on `:8080`.

---

## 🔐 Environment Variables

| Variable             | Description                                 | Default      |
|----------------------|---------------------------------------------|--------------|
| `DB_URL`             | PostgreSQL connection string                | *required*   |
| `JWT_ACCESS_EXPIRY`  | Access token lifetime                       | `15m`        |
| `JWT_REFRESH_EXPIRY` | Refresh token lifetime                      | `168h`       |
| `JWT_SIGNING_KEY`    | Secret key for signing JWT tokens           | *required*   |

---

## 📖 API Documentation

All endpoints are prefixed with `/api` (except auth). Authentication is required for protected endpoints via `Authorization: Bearer <token>` header.

### 🔑 Authentication

| Method | Endpoint                 | Description                           |
|--------|--------------------------|---------------------------------------|
| POST   | `/auth/register`         | Register a new user (role optional)   |
| POST   | `/auth/login`            | Login & get access/refresh tokens     |
| POST   | `/auth/refresh`          | Get new access token using refresh    |
| POST   | `/auth/logout`           | Revoke refresh token & logout         |
| GET    | `/auth/me`               | Get current user profile (protected)  |

**Register Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123",
  "role": "customer"   // optional, defaults to customer
}
```

**Login Response:**
```json
{
  "access_token": "eyJ...",
  "user": {
    "id": "...",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "customer"
  }
}
```

---

### 🏢 Restaurants

#### Public Endpoints
| Method | Endpoint                         | Description                     |
|--------|----------------------------------|---------------------------------|
| GET    | `/api/restaurants`               | List all restaurants (paginated) |
| GET    | `/api/restaurants/{id}`          | Get restaurant by ID            |

#### Owner/Admin Only
| Method | Endpoint                         | Description                     |
|--------|----------------------------------|---------------------------------|
| POST   | `/api/restaurants`               | Create a new restaurant         |
| GET    | `/api/my/restaurant`             | Get the owner's own restaurant  |
| PUT    | `/api/restaurants/{id}`          | Update restaurant details       |
| DELETE | `/api/restaurants/{id}`          | Delete a restaurant             |

**Create / Update Request:**
```json
{
  "name": "Garden Way",
  "description": "Authentic Indian cuisine",
  "address": "MG Road, Thane",
  "city": "Thane",
  "state": "Maharashtra",
  "pincode": "400601",
  "phone": "9876543211",
  "email": "contact@gardenway.com",
  "image": "https://...",
  "opening_time": "10:00",
  "closing_time": "22:00",
  "is_open": true
}
```

---

### 📋 Menu

#### Public
| Method | Endpoint                               | Description                          |
|--------|----------------------------------------|--------------------------------------|
| GET    | `/api/restaurants/{restaurantId}/menu` | Get menu items (paginated, filter)   |
| GET    | `/api/menu/{id}`                       | Get a single menu item               |

**Query params:** `page`, `limit`, `search`, `category`, `available`

#### Owner/Admin Only
| Method | Endpoint                               | Description                          |
|--------|----------------------------------------|--------------------------------------|
| POST   | `/api/restaurants/{restaurantId}/menu` | Create a menu item                   |
| PUT    | `/api/menu/{id}`                       | Update menu item                     |
| PATCH  | `/api/menu/{id}/availability`          | Toggle availability                  |
| DELETE | `/api/menu/{id}`                       | Delete menu item                     |

**Create/Update Request:**
```json
{
  "name": "Butter Chicken",
  "category": "Main Course",
  "description": "Creamy tomato curry",
  "price": 450.50,
  "image": "https://...",
  "is_available": true
}
```

---

### 🪑 Tables

#### Public
| Method | Endpoint                               | Description                          |
|--------|----------------------------------------|--------------------------------------|
| GET    | `/api/restaurants/{restaurantId}/tables`| List all tables of a restaurant      |
| GET    | `/api/tables/{id}`                     | Get a single table                   |

#### Owner/Admin Only
| Method | Endpoint                               | Description                          |
|--------|----------------------------------------|--------------------------------------|
| POST   | `/api/restaurants/{restaurantId}/tables`| Create a table                       |
| PUT    | `/api/tables/{id}`                     | Update table (number, capacity)      |
| PATCH  | `/api/tables/{id}/availability`        | Toggle table availability            |
| DELETE | `/api/tables/{id}`                     | Delete a table                       |

**Create/Update Request:**
```json
{
  "table_number": 5,
  "capacity": 4
}
```

---

### 📅 Reservations

#### Customer/Admin
| Method | Endpoint                         | Description                          |
|--------|----------------------------------|--------------------------------------|
| POST   | `/api/reservations`              | Create a reservation                 |
| GET    | `/api/my/reservations`           | Get current user's reservations      |
| GET    | `/api/reservations/{id}`         | Get a reservation by ID              |
| PATCH  | `/api/reservations/{id}/cancel`  | Cancel a pending reservation         |

#### Owner/Admin Only
| Method | Endpoint                                 | Description                          |
|--------|------------------------------------------|--------------------------------------|
| GET    | `/api/restaurants/{restaurantId}/reservations`| List all reservations for restaurant |
| PATCH  | `/api/reservations/{id}/confirm`         | Confirm a pending reservation        |

**Create Request:**
```json
{
  "restaurant_id": "uuid",
  "table_id": "uuid",
  "reservation_time": "2026-12-25T19:00:00Z",
  "guest_count": 2,
  "notes": "Window seat please"
}
```

---

### 🧾 Orders

#### Customer/Admin
| Method | Endpoint                         | Description                          |
|--------|----------------------------------|--------------------------------------|
| POST   | `/api/orders`                    | Create an order                      |
| GET    | `/api/my/orders`                 | Get user's orders                    |
| GET    | `/api/orders/{id}`               | Get an order by ID                   |
| PATCH  | `/api/orders/{id}/cancel`        | Cancel a pending order               |

#### Owner/Admin Only
| Method | Endpoint                                 | Description                          |
|--------|------------------------------------------|--------------------------------------|
| GET    | `/api/restaurants/{restaurantId}/orders` | List orders for a restaurant         |
| PATCH  | `/api/orders/{id}/status`                | Update order status                  |

**Create Request:**
```json
{
  "restaurant_id": "uuid",
  "notes": "No onions",
  "items": [
    {
      "menu_item_id": "uuid",
      "quantity": 2
    }
  ]
}
```

**Status values:** `pending`, `confirmed`, `preparing`, `ready`, `completed`, `cancelled`

---

### ⭐ Reviews

#### Public
| Method | Endpoint                                 | Description                          |
|--------|------------------------------------------|--------------------------------------|
| GET    | `/api/restaurants/{restaurantId}/reviews`| Get all reviews for a restaurant     |
| GET    | `/api/reviews/{id}`                      | Get a single review                  |

#### Customer/Admin (must have completed an order)
| Method | Endpoint                                 | Description                          |
|--------|------------------------------------------|--------------------------------------|
| POST   | `/api/restaurants/{restaurantId}/reviews`| Create a review                      |
| GET    | `/api/my/reviews`                        | Get user's reviews                   |
| PUT    | `/api/reviews/{id}`                      | Update own review                    |
| DELETE | `/api/reviews/{id}`                      | Delete own review                    |

**Create/Update Request:**
```json
{
  "rating": 5,
  "comment": "Excellent food!"
}
```

---

## 🔄 Authentication Flow

1. **Register** – user provides `name`, `email`, `password`, optional `role`.
2. **Login** – server validates credentials, returns `access_token` (short-lived) and sets an `http‑only` cookie with `refresh_token`.
3. **Access Protected Endpoints** – client sends `Authorization: Bearer <access_token>`.
4. **Refresh** – when access token expires, client calls `/auth/refresh`; the refresh token cookie is automatically sent. A new access token is returned.
5. **Logout** – revokes the refresh token and clears the cookie.

**Role mapping:**
- `admin` – full access (can manage any restaurant)
- `owner` – owns one restaurant; can manage its tables, menu, reservations, orders
- `staff` – (currently not used; can be extended)
- `customer` – can create reservations, orders, and reviews after completing an order

---

## 🗄️ Database Schema (ORM Models)

```text
User
  - id (UUID, PK)
  - name, email, password_hash
  - role (customer, owner, staff, admin)
  - restaurant_id (nullable, FK to Restaurant)
  - timestamps

RefreshToken
  - id (PK)
  - user_id (FK)
  - token_hash (unique)
  - expires_at, revoked, created_at

Restaurant
  - id (UUID, PK)
  - owner_id (FK to User)
  - name, description, address, city, state, pincode
  - phone, email, image
  - opening_time, closing_time, is_open
  - timestamps

Menu
  - id (UUID, PK)
  - restaurant_id (FK)
  - name, category, description, price
  - image, is_available
  - timestamps

RestaurantTable
  - id (UUID, PK)
  - restaurant_id (FK)
  - table_number, capacity, is_available
  - timestamps

Reservation
  - id (UUID, PK)
  - restaurant_id, table_id, user_id (FKs)
  - reservation_time, guest_count
  - status (pending, confirmed, cancelled, completed)
  - notes, timestamps

Order
  - id (UUID, PK)
  - restaurant_id, user_id (FKs)
  - total_amount, status, notes
  - timestamps
  - has many OrderItems

OrderItem
  - id (UUID, PK)
  - order_id (FK)
  - menu_item_id (FK)
  - quantity, price (snapshot)

Review
  - id (UUID, PK)
  - restaurant_id, user_id (FKs)
  - rating (1-5), comment
  - timestamps
```

---

## 🚀 Testing the API

You can test endpoints using **cURL** or tools like **Postman**.

Example workflow:

1. **Register a customer** (role default)
2. **Login** → get `access_token`
3. **Create a restaurant** (as admin/owner) → get `restaurant_id`
4. **Create tables** for that restaurant
5. **Create menu items**
6. **Place an order** (as customer) → get `order_id`
7. **Update order status** to `completed` (as owner)
8. **Create a reservation**
9. **Review the restaurant** (only after completed order)

---

## 📝 License

This project is open‑source and available under the [MIT License](LICENSE).