#!/bin/bash

# Create Soroban identity and fund it automatically

echo "Creating Soroban identity..."
soroban keys generate --global testing

echo "Getting public key..."
PUBLIC_KEY=$(soroban keys address test-account)
echo "Public Key: $PUBLIC_KEY"

echo "Funding account via Friendbot..."
FUND_RESPONSE=$(curl -s "https://friendbot.stellar.org/?addr=$PUBLIC_KEY")

if [[ $FUND_RESPONSE == *"successful"* ]]; then
    echo "✅ Account funded successfully!"
    echo "Account: $PUBLIC_KEY"
else
    echo "❌ Funding failed. Response:"
    echo "$FUND_RESPONSE"
fi

echo "Adding testnet network configuration..."
soroban network add \
  --global testnet \
  --rpc-url https://soroban-testnet.stellar.org:443 \
  --network-passphrase "Test SDF Network ; September 2015"

echo "Setup complete! You can now use 'test-account' as your source account."