# Go Coffeeder Admin API

A backend API for the Coffeeder application built using Go (Gin Framework) and PostgreSQL. It supports product, user, order, cart, and transaction management with CRUD, pagination, search, caching, and authentication features.

## Tech Stack

- **Framework**: Gin Web Framework
- **Database**: PostgreSQL dengan pgx/v5 (pgxpool)
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: Argon2
- **Caching**: Redis
- **Validation**: go-playground/validator v10
- **Email Service**: Gomail
- **File Upload**: Cloudinary
- **Environment Variables**: godotenv
- **Database Migration**: golang-migrate
- **API Documentation**: HATEOAS (Hypermedia as the Engine of Application State)

## Features

-  User authentication & authorization (JWT)
-  Password reset dengan OTP email
-  Product management (CRUD with images)
-  Shopping cart system
-  Transaction & order management
-  Category management
-  Redis caching untuk performa optimal
-  Pagination & search functionality
-  Role-based access control (Admin/User)
-  Image upload ke Cloudinary
-  Favorite products
-  Product variants & sizes
-  Promo management

## Preview ERD

```mermaid
erDiagram
    users {
        bigint id PK
        varchar(100) fullname
        varchar(100) email
        varchar(100) password
        varchar(100) role
        text reset_token
        timestamp reset_expires
        varchar(25) reset_otp
        timestamp created_at
        timestamp updated_at
    }

    profile {
        bigint id PK
        bigint user_id FK
        varchar(250) image
        varchar(50) phone
        varchar(250) address
        timestamp created_at
        timestamp updated_at
    }

    forgot_password {
        bigint id PK
        bigint user_id FK
        varchar(100) token
        timestamp expires_at
        timestamp created_at
    }

    variants {
        bigint id PK
        varchar(50) name
        int additional_price
    }

    sizes{
        bigint id PK
        varchar(100) name
        numeric additional_price
    }

    categories {
        bigint id PK
        varchar(100) name
        timestamp created_at
        timestamp updated_at
    }

    products {
        bigint id PK
        varchar(100) title
        varchar(250) description
        bigint category_id FK
        int stock
        numeric base_price
        boolean is_favorite
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    product_variants {
        bigint id PK
        bigint product_id FK
        bigint variant_id FK
        timestamp created_at
        timestamp updated_at
    }

   product_sizes{
        bigint product_id FK
        int size_id FK
    }

    products_categories {
        bigint product_id FK
        bigint category_id FK
    }

    product_images {
        bigint id PK
        bigint product_id FK
        text image
        timestamp updated_at
        timestamp deleted_at
    }

    promos {
        bigint id PK
        varchar(100) title
        varchar(100) description
        float discount
        timestamp start
        timestamp end
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    product_promos {
        bigint id PK
        bigint promo_id FK
        bigint product_id FK
    }

    payment_methods {
        bigint id PK
        varchar(100) name
        varchar(250) image
        timestamp created_at
        timestamp updated_at
    }

    shippings{
        bigint id PK
        varchar(50) name
    }

    status {
        bigint id PK
        varchar(20) name
    }

    recommended_products {
        bigint product_id FK
        bigint recommended_id FK
        timestamp created_at
    }

   carts {
        bigint id PK
        bigint user_id FK
        bigint product_id FK
        int size_id FK
        bigint variant_id FK
        int quantity
        timestamp created_at
        timestamp updated_at
    }

    transactions {
        bigint id PK
        bigint user_id FK
        varchar(100) fullname
        varchar(100) email
        varchar(20) phone
        varchar(250) address
        bigint payment_method_id FK
        bigint shipping_id FK
        varchar(50) invoice_number
        numeric total
        varchar(50) status
        timestamp created_at
        timestamp updated_at
    }

    transaction_items {
        bigint id PK
        bigint transaction_id FK
        bigint product_id FK
        bigint variant_id FK
        int size_id FK
        int quantity
        numeric subtotal
    }

    %% RELATIONSHIPS
    users ||--|| profile : "has"
    users ||--|| forgot_password : "has"
    users ||--o{carts : "owns"
    users ||--o{ transactions : "makes"

    products ||--o{ product_variants : "has variants"
    variants ||--o{ product_variants : "used in"

    products ||--o{product_sizes: "has sizes"
    sizes||--o{product_sizes: "available for"

    products ||--o{ products_categories : "belongs to"
    categories ||--o{ products_categories : "categorizes"

    products ||--o{ product_images : "has images"

    products ||--o{ product_promos : "applies promo"
    promos ||--o{ product_promos : "applied to"

    products ||--o{ recommended_products : "recommended"
    recommended_products ||--|| products : "linked"

   carts ||--|| products : "contains product"
   carts ||--|| sizes: "selected size"
   carts ||--|| variants : "selected variant"

    transactions ||--|| users : "belongs to"
    transactions ||--|| payment_methods : "paid via"
    transactions ||--|| shippings: "shipped via"
    transaction_items ||--|| transactions : "part of"
    transaction_items ||--|| products : "includes"
    transaction_items ||--|| sizes: "with size"
    transaction_items ||--|| variants : "with variant"
    transactions ||--o{ status : "status reference"
```

## Installation

### Prerequisites

- Go 1.25
- PostgreSQL 
- Redis
- Cloudinary Account

### Setup

1. Clone repository:
```bash
git clone <repository-url>
cd coffeeder-backend
```

2. Install dependencies:
```bash
go mod download
```

3. Setup environment variables (`.env`):
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=ecommerce_db

# Redis
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key

# Cloudinary
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

4. Run database migrations:
```bash
migrate -path db/migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up
```

5. Run the application:
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Authentication
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/auth/register` | Register user baru | - |
| POST | `/auth/login` | Login user | - |
| POST | `/auth/forgot-password` | Request reset password | - |
| POST | `/auth/verify-otp` | Verify OTP code | - |
| PATCH | `/auth/reset-password` | Reset password dengan token | - |

### Admin - Users
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/admin/users` | List users (pagination & search) | Admin |
| GET | `/admin/users/:id` | Get user by ID | Admin |
| POST | `/admin/users` | Create new user | Admin |
| PATCH | `/admin/users/:id` | Update user | Admin |
| DELETE | `/admin/users/:id` | Delete user | Admin |

### Admin - Products
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/admin/products` | Create product | Admin |
| GET | `/admin/products` | List products (pagination & search) | Admin |
| GET | `/admin/products/:id` | Get product by ID | Admin |
| PATCH | `/admin/products/:id` | Update product | Admin |
| DELETE | `/admin/products/:id` | Delete product | Admin |
| GET | `/admin/products/:id/images` | Get product images | Admin |
| GET | `/admin/products/:id/images/:image_id` | Get specific image | Admin |
| PATCH | `/admin/products/:id/images/:image_id` | Update product image | Admin |
| DELETE | `/admin/products/:id/images/:image_id` | Delete product image | Admin |
| GET | `/admin/type-products` | Get product types | Admin |

### Admin - Categories
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/admin/categories` | List categories | Admin |
| GET | `/admin/categories/:id` | Get category by ID | Admin |
| POST | `/admin/categories` | Create category | Admin |
| PATCH | `/admin/categories/:id` | Update category | Admin |
| DELETE | `/admin/categories/:id` | Delete category | Admin |

### Admin - Transactions
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/admin/transactions` | List all transactions | Admin |
| GET | `/admin/transactions/:id` | Get transaction detail | Admin |
| PATCH | `/admin/transactions/:id/status` | Update transaction status | Admin |
| DELETE | `/admin/transactions/:id` | Delete transaction | Admin |

### Public - Products
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/products` | Filter & search products | - |
| GET | `/products/:id` | Get product detail | - |
| GET | `/favorite-products` | Get favorite products | - |

### User - Cart
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/cart` | Add item to cart | User |
| GET | `/cart` | Get user cart | User |
| DELETE | `/deletecart` | Remove item from cart | User |

### User - Transactions
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/transactions` | Create new transaction | User |
| GET | `/history` | Get transaction history | User |
| GET | `/history/:id` | Get transaction detail | User |
| GET | `/shippings` | Get shipping methods | User |
| GET | `/payment-methods` | Get payment methods | User |

## Performance

### Redis Caching Implementation

API ini menggunakan Redis untuk caching data yang sering diakses, meningkatkan performa secara signifikan:

**Sebelum Caching:**
- Response time: ~33ms

![Before Caching](images/sebelumCaching.png)

**Setelah Caching:**
- Response time: ~6ms (82% lebih cepat)

![After Caching](images/sesudah.png)


## Security Features

- Password hashing dengan Argon2
- JWT-based authentication
- Role-based access control (RBAC)
- Request validation dengan validator v10
- OTP verification untuk password reset
- Secure token management

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Contact

For questions or support, please open an issue in the repository.