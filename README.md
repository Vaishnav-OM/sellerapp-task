SellerApp Backend Task

## Prerequisites

- Go 1.16 or later
- MongoDB instance (local or cloud)
- Git (for cloning the repository)
- Docker

## Setup

1. **Clone the Repository**

   ```sh
   git clone https://github.com/yourusername/product-service.git
   cd product-service
   ```

2. **Install Dependencies**

   ```sh
   go mod download
   ```

3. **Set Up MongoDB**

   Ensure your MongoDB instance is running and accessible. Update the connection details in `db/mongo.go` if necessary.

4. **Initialize the Collection**

   Added a init-mongo file so that db and collection will be initialised when `docker-compose up --build` is called.
   The collection is initialized in `main.go`.

## Running the Service

1.  **Start the Service**

    run `main.go`

    ```sh
    go run main.go
    ```

    The service will be available at `http://localhost:8080`.

2.  **Docker Setup**

    ->Build and Start the Service

          ```sh
          docker-compose up --build
          ```

    This command will build the Docker images and start the containers. The service will be available at `http://localhost:8080`.

    ->Stop the Service

    ```sh
      docker-compose down
    ```

## Endpoints

- `POST /products`: Create a new product
- `PUT /products/:id`: Update an existing product
- `DELETE /products/:id`: Delete a product
- `GET /products`: List products with optional query parameters for filtering and pagination
- `GET /products/:id`: Get a product by ID

## Testing

1. **Run Unit Tests**

   Make sure your MongoDB instance is running and accessible, also make sure you are in the handlers folder, then run the tests:

   ```sh
   go test -v ./...
   ```

   This will execute the test cases defined in `handlers` package.
