# Chama Wallet Backend

A blockchain-powered savings and lending platform for Chamas (informal savings groups) built on the Stellar network using Go and Fiber.

## 🌟 Features

### Core Functionality
- **User Authentication**: JWT-based authentication with registration and login
- **Wallet Management**: Create wallets, check balances, transfer funds
- **Group Management**: Create and manage savings groups (Chamas)
- **Stellar Integration**: Full integration with Stellar blockchain
- **Member Management**: Add members to groups and track contributions
- **Transaction History**: View transaction history for wallets

### Technical Features
- **RESTful API**: Clean REST endpoints for all operations
- **JWT Authentication**: Secure token-based authentication
- **PostgreSQL Database**: Reliable data persistence with GORM
- **Stellar Network**: Testnet integration for development
- **CORS Support**: Cross-origin resource sharing enabled
- **Error Handling**: Comprehensive error handling and validation

## 🚀 Getting Started

### Prerequisites
- Go 1.23+
- PostgreSQL 12+
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd chama-wallet-backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Database Setup**
   ```bash
   # Create PostgreSQL database
   createdb chama_wallet
   
   # Create user (optional)
   psql -c "CREATE USER chama_user WITH PASSWORD 'malika';"
   psql -c "GRANT ALL PRIVILEGES ON DATABASE chama_wallet TO chama_user;"
   ```

4. **Environment Configuration**
   Update the database connection string in `database/db.go`:
   ```go
   dsn := "host=localhost user=chama_user password=malika dbname=chama_wallet port=5432 sslmode=disable"
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

6. **Server will start on `http://localhost:3000`**

## 📁 Project Structure

```
chama-wallet-backend/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── database/
│   └── db.go              # Database connection and migration
├── models/
│   ├── group.go           # Group, Member, Contribution models
│   └── user.go            # User and auth models
├── handlers/
│   ├── wallet.go          # Wallet operation handlers
│   ├── group.go           # Group operation handlers
│   ├── group_handlers.go  # Additional group handlers
│   └── auth.go            # Authentication handlers
├── routes/
│   ├── routes.go          # Wallet routes
│   ├── group.go           # Group routes
│   └── auth.go            # Authentication routes
├── services/
│   ├── stellar.go         # Stellar blockchain operations
│   ├── stellar_service.go # Additional Stellar services
│   ├── balance.go         # Balance checking services
│   ├── fund.go            # Account funding services
│   ├── group_service.go   # Group management services
│   └── auth_service.go    # Authentication services
├── middleware/
│   └── auth.go            # JWT authentication middleware
└── utils/
    └── wallet.go          # Wallet utility functions
```

## 🔐 Authentication API

### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "wallet": "STELLAR_ADDRESS",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "jwt_token_here"
}
```

### Login User
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

### Get Profile
```http
GET /auth/profile
Authorization: Bearer <jwt_token>
```

### Update Profile
```http
PUT /auth/profile
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "John Updated"
}
```

### Logout
```http
POST /auth/logout
Authorization: Bearer <jwt_token>
```

## 💰 Wallet API

### Create Wallet
```http
POST /create-wallet
```

**Response:**
```json
{
  "address": "STELLAR_PUBLIC_KEY",
  "seed": "STELLAR_SECRET_KEY"
}
```

### Get Balance
```http
GET /balance/{address}
```

**Response:**
```json
{
  "balances": ["native: 1000.0000000"]
}
```

### Transfer Funds
```http
POST /transfer
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "from_seed": "SECRET_KEY",
  "to_address": "DESTINATION_ADDRESS",
  "amount": "100"
}
```

### Generate Keypair
```http
GET /generate-keypair
```

### Fund Account (Testnet)
```http
POST /fund/{address}
```

### Get Transaction History
```http
GET /transactions/{address}
```

## 👥 Group API

### Create Group
```http
POST /group/create
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "Alpha Chama",
  "description": "Our savings group for community development"
}
```

**Response:**
```json
{
  "ID": "group_uuid",
  "Name": "Alpha Chama",
  "Description": "Our savings group for community development",
  "Wallet": "GROUP_STELLAR_ADDRESS",
  "Members": [],
  "Contributions": []
}
```

### Get All Groups
```http
GET /groups
```

### Join Group
```http
POST /group/{id}/join
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "wallet": "USER_WALLET_ADDRESS"
}
```

### Contribute to Group
```http
POST /group/{id}/contribute
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "from": "USER_WALLET_ADDRESS",
  "secret": "USER_SECRET_KEY",
  "amount": "50"
}
```

### Get Group Balance
```http
GET /group/{id}/balance
```

**Response:**
```json
{
  "group_id": "group_uuid",
  "wallet": "GROUP_WALLET_ADDRESS",
  "balance": "500.0000000"
}
```

## 🔒 Authentication & Security

### JWT Authentication
- **Token Expiry**: 24 hours
- **Algorithm**: HS256
- **Claims**: User ID and email

### Protected Routes
Routes requiring authentication:
- `POST /transfer`
- `POST /group/create`
- `POST /group/:id/contribute`
- `POST /group/:id/join`
- `GET /auth/profile`
- `PUT /auth/profile`

### Optional Authentication
Routes with optional authentication (enhanced features when authenticated):
- `GET /groups`
- `GET /group/:id/balance`
- `GET /balance/:address`
- `POST /fund/:address`
- `GET /transactions/:address`

### Password Security
- **Hashing**: bcrypt with cost 14
- **Minimum Length**: 6 characters
- **Storage**: Passwords never stored in plain text

## 🌐 Stellar Integration

### Network Configuration
- **Environment**: Testnet (for development)
- **Horizon Server**: `https://horizon-testnet.stellar.org`
- **Network Passphrase**: Test SDF Network ; September 2015

### Wallet Operations
- **Account Creation**: Automatic keypair generation
- **Funding**: Friendbot integration for testnet
- **Transactions**: Native XLM transfers
- **Balance Checking**: Real-time balance queries

### Group Wallets
- Each group gets its own Stellar wallet
- Transparent on-chain transactions
- Multi-signature support (future enhancement)

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id VARCHAR PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    name VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    wallet VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

### Groups Table
```sql
CREATE TABLE groups (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    wallet VARCHAR NOT NULL
);
```

### Members Table
```sql
CREATE TABLE members (
    id VARCHAR PRIMARY KEY,
    group_id VARCHAR REFERENCES groups(id),
    wallet VARCHAR NOT NULL
);
```

### Contributions Table
```sql
CREATE TABLE contributions (
    id VARCHAR PRIMARY KEY,
    group_id VARCHAR REFERENCES groups(id),
    member_id VARCHAR REFERENCES members(id),
    amount DECIMAL NOT NULL
);
```

## 🛠️ Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o chama-wallet main.go
```

### Environment Variables
For production, consider using environment variables:
```bash
export DB_HOST=localhost
export DB_USER=chama_user
export DB_PASSWORD=secure_password
export DB_NAME=chama_wallet
export JWT_SECRET=your-super-secure-secret
export STELLAR_NETWORK=testnet
```

### CORS Configuration
CORS is enabled for all origins in development. For production:
```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://yourdomain.com",
    AllowMethods: "GET,POST,PUT,DELETE",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
}))
```

## 📊 API Response Formats

### Success Response
```json
{
  "data": { ... },
  "message": "Operation successful"
}
```

### Error Response
```json
{
  "error": "Error message description"
}
```

### Authentication Response
```json
{
  "user": {
    "id": "user_id",
    "name": "User Name",
    "email": "user@example.com",
    "wallet": "STELLAR_ADDRESS",
    "created_at": "timestamp"
  },
  "token": "jwt_token"
}
```

## 🚀 Deployment

### Docker Deployment
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Environment Setup
1. Set up PostgreSQL database
2. Configure environment variables
3. Run database migrations
4. Start the application

## 🔧 Configuration

### Database Configuration
Update `database/db.go` with your database credentials:
```go
dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_PORT"),
)
```

### JWT Configuration
Update JWT secret in `services/auth_service.go`:
```go
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
```

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify connection credentials
   - Ensure database exists

2. **Stellar Network Issues**
   - Check internet connection
   - Verify Stellar Horizon server status
   - Ensure using correct network (testnet/mainnet)

3. **Authentication Issues**
   - Verify JWT secret configuration
   - Check token expiry
   - Ensure proper Authorization header format

### Logging
The application logs important events. Check console output for:
- Database connection status
- Route setup confirmation
- Error messages and stack traces

## 📚 Additional Resources

- [Stellar Documentation](https://developers.stellar.org/)
- [Fiber Framework](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [JWT Go Library](https://github.com/golang-jwt/jwt)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the troubleshooting section
- Review the API documentation

---

Built with ❤️ for the Chama community using Go, Fiber, PostgreSQL, and Stellar blockchain technology.