#!/bin/bash

# Set all environment secrets for Fly.io deployment
echo "Setting environment secrets for chama-wallet-backend on Fly.io..."

# Database Configuration
flyctl secrets set DATABASE_URL="postgresql://neondb_owner:npg_ZVKPWQE85tae@ep-wispy-bush-aez0yr2l-pooler.c-2.us-east-2.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

# Stellar/Soroban Configuration
flyctl secrets set SOROBAN_NETWORK="testnet"
flyctl secrets set SOROBAN_RPC_URL="https://soroban-testnet.stellar.org"
flyctl secrets set SOROBAN_NETWORK_PASSPHRASE="Test SDF Network ; September 2015"

# Stellar Keys
flyctl secrets set SOROBAN_SECRET_KEY="SAHUN4JW2MYG4IGCBPGFGPBKWS7CV43V54SVWL4AZSKLDB7FQBMDWFVL"
flyctl secrets set SOROBAN_PUBLIC_KEY="GCIVECFQHJI3BOLOTFUPUFQ3OD7YC2SA5VYPTNY5DQOZYGYDXWZQXXI5"

# Server Configuration
flyctl secrets set PORT="8080"

# CORS Configuration (update with your frontend URL after deployment)
flyctl secrets set ALLOWED_ORIGINS="http://localhost:5173,https://your-frontend-domain.com"

# JWT Secret (you should change this to a secure random string)
flyctl secrets set JWT_SECRET="5cUiffEpgVSyk2rs/KDXy3eFFUYSCEzlSiWEQ8J9Te8="

echo "✅ All secrets have been set!"
echo "Remember to update ALLOWED_ORIGINS with your actual frontend URL after deployment"
