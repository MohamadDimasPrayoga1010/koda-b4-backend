# Go E-Commerce Admin API

API ini adalah backend untuk aplikasi e-commerce, dibangun menggunakan Go (Gin) dan PostgreSQL. Mendukung manajemen produk, user, order, dan lainnya dengan fitur CRUD, pagination, dan search.

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

## API ENDPOINT

### Admin
Method	Endpoint	        Deskripsi
- GET	/admin/users	    List user dengan pagination dan search
- GET /admin/users/:id      Dapatkan user berdasrkan id
- POST	/admin/users	    Tambah user baru
- PATCH	/admin/users/:id	Edit user
- DELETE /admin/users/:id   Delete user

### Product Admin
Method	Endpoint	                            Deskripsi
- POST	/admin/products	                        Tambah produk
- GET	/admin/products	                        List produk (pagination & search)
- GETH	/admin/products/:id	                    produk by id
- DELETE	/admin/products/:id	                Hapus produk
- PATCH /admin/products/:id                     Edit product
- GET /products/:id/images                      Image id product
- GET /products/:id/images/:image_id            Product id images id
- PATCH /products/:id/images/:image_id          Edit Product id images id
- DELETE GET /products/:id/images/:image_id     Delete Product id images id


### Order Admin
Method	Endpoint	Deskripsi
- GET	/admin/orders	List orders
- PATCH	/admin/orders/:id/status	Update status order
- DELETE	/admin/orders/:id	Hapus order

### Preview waktu sebelum caching
Waktu 33 ms

![img](images/sebelumCaching.png)

### Preview waktu sesudah caching
Waktu 6 ms

![img](images/sesudah.png)

