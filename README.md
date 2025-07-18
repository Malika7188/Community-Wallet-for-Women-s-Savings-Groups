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