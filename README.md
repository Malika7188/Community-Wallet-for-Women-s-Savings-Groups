# Community Wallet for Women's Savings Groups

A blockchain-powered savings and lending platform for Chamas (informal savings groups), built on Stellar and Soroban smart contracts. The platform enables group creation, wallet management, secure contributions, payouts, and notifications.

---

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Setup Instructions](#setup-instructions)
  - [Backend](#backend)
  - [Frontend](#frontend)
- [API Endpoints](#api-endpoints)
- [Smart Contract Operations](#smart-contract-operations)
- [Scripts](#scripts)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [License](#license)

---

## Features

- Create and manage savings groups (Chamas)
- Individual and group wallets (Stellar-based)
- Deposit, withdraw, and transfer funds (XLM)
- Group contributions and payouts
- Notifications for group activities and approvals
- Transaction history and balance queries
- Secure authentication and authorization

---

## Architecture

- **Backend:** Go (Fiber), Stellar SDK, Soroban CLI, GORM (database)
- **Frontend:** React, TypeScript, Tailwind CSS
- **Smart Contracts:** Rust (Soroban)

---

## Tech Stack

- Go, Fiber, GORM, Stellar SDK, Soroban CLI
- React, TypeScript, Tailwind CSS
- Rust (for Soroban contracts)
- PostgreSQL (recommended for production)

---

## Setup Instructions

### Backend

1. **Install Go and dependencies:**
   ```bash
   sudo apt update
   sudo apt install golang-go
   cd chama-wallet-backend
   go mod tidy
   ```

2. **Install Stellar CLI and Soroban:**
   ```bash
   sudo apt install -y libudev-dev pkg-config
   cargo install stellar-cli --locked --version 23.0.0
   stellar --version
   ```

3. **Run the backend server:**
   ```bash
   go run .
   # Server runs on http://localhost:3000
   ```

### Frontend

1. **Install Node.js and dependencies:**
   ```bash
   cd chama-wallet-frontend
   npm install
   ```

2. **Run the frontend app:**
   ```bash
   npm run dev
   # App runs on http://localhost:5173
   ```

---

## API Endpoints

### Wallet

- `POST /create-wallet` — Create a new wallet
- `GET /balance/:address` — Get wallet balance
- `GET /generate-keypair` — Generate Stellar keypair
- `POST /fund/:address` — Fund a wallet (testnet)
- `GET /transactions/:address` — Get transaction history
- `POST /transfer` — Transfer XLM between wallets

### Groups

- `POST /group/create` — Create a new group
- `GET /groups` — List all groups

### Smart Contract Operations

- `POST /api/contribute` — Contribute to a group contract
- `POST /api/balance` — Get contract balance
- `POST /api/withdraw` — Withdraw from contract
- `POST /api/history` — Get contract transaction history

---

## Smart Contract Operations (Soroban CLI)

- **Build contract:**
  ```bash
  stellar contract build
  ```
- **Deploy contract:**
  ```bash
  stellar contract deploy --source-account <account>
  ```
- **Invoke contract (initialize, contribute, get_balance):**
  ```bash
  soroban contract invoke --id <contract_id> --source-account <account> --network testnet -- initialize
  soroban contract invoke --id <contract_id> --source-account <account> --network testnet -- contribute --user <user> --amount <amount>
  soroban contract invoke --id <contract_id> --source-account <account> --network testnet -- get_balance --user <user>
  ```

---

## Scripts

- `create-group.sh` — Create a new group
- `test-operations.sh` — Run contract operations
- `manage-groups.sh` — Manage groups (list, show, contribute, balance)

---

## Testing

- **Fund test accounts:**
  ```bash
  curl "https://friendbot.stellar.org/?addr=YOUR_PUBLIC_KEY"
  soroban keys fund --rpc-url https://soroban-testnet.stellar.org:443 --network-passphrase "Test SDF Network ; September 2015" <account>
  ```
- **Check keys:**
  ```bash
  soroban keys show <wallet>
  soroban keys address <wallet>
  ```

---

## Troubleshooting

- Ensure all dependencies are installed (Go, Node.js, Cargo, Stellar CLI)
- Use testnet for development and testing
- If contract deployment fails, check Soroban CLI and network status
- For API errors, check backend logs and request payloads

---

## License

This project is licensed under the MIT License.

---

## Authors

- [**Malika**](https://github.com/Malika7188)
- [**Andrew**](https://github.com/andyosyndoh)

