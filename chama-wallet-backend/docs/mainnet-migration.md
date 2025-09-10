# Mainnet Migration Guide

This guide explains how to migrate your Chama Wallet from Stellar Testnet to Mainnet.

## 🚨 Important Security Notice

**MAINNET USES REAL MONEY!** 
- All transactions involve real XLM and USDC
- Test thoroughly on testnet before switching to mainnet
- Keep secret keys secure and never share them
- Start with small amounts for testing

## 🔧 Configuration Changes

### 1. Environment Variables

Update your `.env` file with mainnet configuration:

```bash
# Network Configuration
STELLAR_NETWORK=mainnet

# Mainnet Endpoints
STELLAR_HORIZON_URL=https://horizon.stellar.org
STELLAR_SOROBAN_RPC_URL=https://soroban-rpc.mainnet.stellar.org:443
STELLAR_NETWORK_PASSPHRASE=Public Global Stellar Network ; September 2015

# Your Mainnet Account (replace with your actual keys)
SOROBAN_PUBLIC_KEY=YOUR_MAINNET_PUBLIC_KEY
SOROBAN_SECRET_KEY=YOUR_MAINNET_SECRET_KEY

# Your Deployed Contract ID (replace with actual contract ID)
SOROBAN_CONTRACT_ID=YOUR_MAINNET_CONTRACT_ID

# USDC Configuration (optional)
USDC_ASSET_CODE=USDC
USDC_ASSET_ISSUER=GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN

# Security Settings
REQUIRE_MEMO_FOR_TRANSFERS=true
MIN_TRANSFER_AMOUNT=0.0000001
MAX_TRANSFER_AMOUNT=10000
```

### 2. Switch Back to Testnet

To switch back to testnet, simply change:

```bash
STELLAR_NETWORK=testnet
```

All other configuration will automatically adjust.

## 🚀 Deployment Steps

### 1. Prepare Mainnet Account

```bash
# Generate a new mainnet keypair
stellar keys generate --global mainnet-account

# Get the public key
stellar keys address mainnet-account

# Fund the account with real XLM (minimum 1 XLM recommended)
# You can buy XLM from exchanges like Coinbase, Kraken, etc.
```

### 2. Deploy Smart Contract

```bash
# Run the automated deployment script
chmod +x scripts/deploy-mainnet-contract.sh
./scripts/deploy-mainnet-contract.sh
```

Or deploy manually:

```bash
cd chama_savings
stellar contract build

# Deploy to mainnet
stellar contract deploy \
  --source-account mainnet-account \
  --network mainnet \
  --wasm target/wasm32-unknown-unknown/release/chama_savings.wasm

# Initialize the contract
stellar contract invoke \
  --id YOUR_CONTRACT_ID \
  --source-account mainnet-account \
  --network mainnet \
  -- initialize
```

### 3. Update Configuration

Update your `.env` file with the deployed contract ID and restart your server.

### 4. Test the Migration

```bash
# Check network status
curl http://localhost:3000/network

# Test balance check (should show mainnet)
curl http://localhost:3000/balance/YOUR_MAINNET_ADDRESS
```

## 📊 API Changes

### Network Information Endpoint

New endpoint to check current network configuration:

```bash
GET /network
```

**Response:**
```json
{
  "network": "mainnet",
  "horizon_url": "https://horizon.stellar.org",
  "soroban_rpc_url": "https://soroban-rpc.mainnet.stellar.org:443",
  "network_passphrase": "Public Global Stellar Network ; September 2015",
  "contract_id": "CXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "is_mainnet": true,
  "supported_assets": {
    "XLM": {
      "code": "XLM",
      "issuer": "native",
      "type": "native"
    },
    "USDC": {
      "code": "USDC",
      "issuer": "GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN",
      "type": "credit_alphanum4"
    }
  }
}
```

### Enhanced API Responses

All API responses now include network information:

```json
{
  "message": "Transfer completed successfully",
  "transaction_hash": "abc123...",
  "network": "mainnet",
  "explorer_url": "https://stellar.expert/explorer/public/tx/abc123..."
}
```

## 💰 Asset Support

### XLM (Native Asset)
- Available on both testnet and mainnet
- Default asset for all operations
- Minimum amount: 0.0000001 XLM

### USDC (Mainnet Only)
- Available only on mainnet
- Requires trustline setup
- Issued by Centre.io

### Transfer with Asset Type

```bash
POST /transfer
{
  "from_seed": "SECRET_KEY",
  "to_address": "DESTINATION_ADDRESS",
  "amount": "100",
  "asset_type": "USDC"  // Optional: "XLM" (default) or "USDC"
}
```

## 🔒 Security Features

### 1. Transaction Limits
- Minimum transfer amount validation
- Maximum transfer amount limits
- Configurable via environment variables

### 2. Memo Requirements
- Automatic memo addition for mainnet transactions
- Compliance with exchange requirements
- Configurable via `REQUIRE_MEMO_FOR_TRANSFERS`

### 3. Account Validation
- Enhanced error messages for mainnet
- Proper handling of non-existent accounts
- Real-time balance validation

## 🧪 Testing on Mainnet

### Example: Real Mainnet Transaction

**Request:**
```bash
curl -X POST http://localhost:3000/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "from_seed": "SXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    "to_address": "GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    "amount": "1.0000000",
    "asset_type": "XLM"
  }'
```

**Response:**
```json
{
  "message": "Transfer completed successfully",
  "transaction_hash": "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
  "from": "GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "to": "GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "amount": "1.0000000",
  "asset_type": "XLM",
  "network": "mainnet",
  "ledger": 45123456,
  "explorer_url": "https://stellar.expert/explorer/public/tx/a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456"
}
```

### Example: Group Contribution on Mainnet

**Request:**
```bash
curl -X POST http://localhost:3000/group/GROUP_ID/contribute \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "from": "GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    "secret": "SXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
    "amount": "50.0000000"
  }'
```

**Response:**
```json
{
  "message": "Contribution successful",
  "group_id": "group-uuid-here",
  "group_name": "Alpha Chama",
  "from": "GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "to": "GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "amount": "50.0000000",
  "tx_hash": "soroban-transaction-result",
  "network": "mainnet",
  "contribution": {
    "ID": "contribution-uuid",
    "GroupID": "group-uuid",
    "UserID": "user-uuid",
    "Amount": 50,
    "Status": "confirmed",
    "TxHash": "soroban-transaction-result",
    "CreatedAt": "2025-01-27T10:30:00Z"
  }
}
```

## 🔄 Switching Between Networks

### Quick Switch Commands

**Switch to Mainnet:**
```bash
echo "STELLAR_NETWORK=mainnet" > .env.network
cat .env.network .env.example > .env.new && mv .env.new .env
```

**Switch to Testnet:**
```bash
echo "STELLAR_NETWORK=testnet" > .env.network
cat .env.network .env.example > .env.new && mv .env.new .env
```

### Verification

After switching networks, verify the configuration:

```bash
curl http://localhost:3000/network | jq .
```

## 🚨 Troubleshooting

### Common Issues

1. **"Account not found on mainnet"**
   - Solution: Fund the account with real XLM first

2. **"Contract does not exist"**
   - Solution: Deploy the contract to mainnet first

3. **"Insufficient balance"**
   - Solution: Ensure account has enough XLM for transaction + fees

4. **"Transaction failed"**
   - Solution: Check transaction limits and account status

### Support

- Check transaction status on [Stellar Explorer](https://stellar.expert/explorer/public)
- Monitor account balances and transaction history
- Review server logs for detailed error information

## 📚 Additional Resources

- [Stellar Mainnet Documentation](https://developers.stellar.org/docs/networks/mainnet)
- [Soroban Mainnet Guide](https://soroban.stellar.org/docs/getting-started/deploy-to-mainnet)
- [Stellar Explorer](https://stellar.expert/explorer/public)
- [USDC on Stellar](https://www.centre.io/usdc-multichain/stellar)

---

**Remember:** Always test on testnet first, then start with small amounts on mainnet!