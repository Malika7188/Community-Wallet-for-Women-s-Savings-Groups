# Community-Wallet-for-Women-s-Savings-Groups
A blockchain-powered savings and lending platform for Chamas or informal savings groups using Stellar.


- To view balance 
```bash
  curl http://127.0.0.1:3000/wallet/GDRGMQ56QAK3DARIWK4BQ6UUC62SKWWUQ5XR7AU5UFV3OIDY5RP2S3F7/balance
```


to generate new keys
```bash
curl http://localhost:3000/wallet/new
```

From: GASDU7PUYFG2ONPMZDD3POTHC4NV7M7ZCIQLHSRSWYDXGMQK445IFVQE
To: GC6HUPHYW33367DVMFX72AFPREAA53QD3SMTCVUQ5BKMHHWGKVE4U3JQ

- create a new group
```bash
 curl -X POST http://localhost:3000/group/create \

  -H "Content-Type: application/json" \
  -d '{
    "name": "Alpha Chama"
  }'
  ```

  get groups
  curl http://localhost:3000/groups

  - stellar instalation
  ```bash
  sudo apt update && sudo apt install -y libudev-dev pkg-config

  cargo install stellar-cli --locked --version 23.0.0

  stellar contract build

  stellar --version

    stellar contract deploy \
  --source-account malika \
 ``` 


 ``bash
 # Initialize the contract (if needed)
soroban contract invoke \
  --id CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4 \
  --source-account malika \
  --network testnet \
  -- initialize

# Make a contribution
soroban contract invoke \
  --id CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4 \
  --source-account malika \
  --network testnet \
  -- contribute --user malika --amount 1000

# Check balance
soroban contract invoke \
  --id CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4 \
  --source-account malika \
  --network testnet \
  -- get_balance --user malika

  ``

  cargo clean

soroban keys show wallet1
secret key: SAHUN4JW2MYG4IGCBPGFGPBKWS7CV43V54SVWL4AZSKLDB7FQBMDWFVL
soroban keys address wallet1 
Account id: GCIVECFQHJI3BOLOTFUPUFQ3OD7YC2SA5VYPTNY5DQOZYGYDXWZQXXI5

to fund existing account
```bash
soroban keys fund --rpc-url https://soroban-testnet.stellar.org:443 --network-passphrase "Test SDF Network ; September 2015" malika
```

curl -X POST http://localhost:3000/api/contribute   -H "Content-Type: application/jsonCADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4", "amount":"100"