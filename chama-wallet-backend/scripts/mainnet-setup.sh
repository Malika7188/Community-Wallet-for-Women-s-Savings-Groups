#!/bin/bash

# Mainnet Setup Script for Chama Wallet
# This script helps configure the application for mainnet deployment

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🌐 Chama Wallet Mainnet Setup${NC}"
echo "============================="
echo ""

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}📝 Creating .env file from template...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✅ .env file created${NC}"
else
    echo -e "${YELLOW}⚠️  .env file already exists${NC}"
fi

echo ""
echo -e "${YELLOW}🔧 Mainnet Configuration Checklist:${NC}"
echo ""

# 1. Network Configuration
echo -e "${BLUE}1. Network Configuration${NC}"
echo "   Set STELLAR_NETWORK=mainnet in your .env file"
echo ""

# 2. Contract Deployment
echo -e "${BLUE}2. Smart Contract Deployment${NC}"
echo "   Deploy your contract to mainnet manually:"
echo "   ${YELLOW}cd chama_savings${NC}"
echo "   ${YELLOW}stellar contract build${NC}"
echo "   ${YELLOW}stellar contract deploy --source-account YOUR_MAINNET_ACCOUNT --network mainnet --wasm target/wasm32-unknown-unknown/release/chama_savings.wasm${NC}"
echo ""

# 3. Account Setup
echo -e "${BLUE}3. Account Setup${NC}"
echo "   Create a mainnet account with real XLM:"
echo "   - Generate keypair: ${YELLOW}stellar keys generate --global mainnet-account${NC}"
echo "   - Fund with real XLM from an exchange or another wallet"
echo "   - Set SOROBAN_PUBLIC_KEY and SOROBAN_SECRET_KEY in .env"
echo ""

# 4. Security Considerations
echo -e "${BLUE}4. Security Considerations${NC}"
echo "   ${RED}⚠️  IMPORTANT SECURITY NOTES:${NC}"
echo "   - Never commit .env files with real secret keys"
echo "   - Use environment variables in production"
echo "   - Consider using hardware wallets for large amounts"
echo "   - Enable transaction memos for compliance"
echo "   - Set appropriate transfer limits"
echo ""

# 5. Testing
echo -e "${BLUE}5. Testing on Mainnet${NC}"
echo "   Start with small amounts to test:"
echo "   - Create a test group with minimal contribution amounts"
echo "   - Test all operations with small XLM amounts first"
echo "   - Verify all transactions on Stellar Explorer"
echo ""

# 6. Environment Variables
echo -e "${BLUE}6. Required Environment Variables for Mainnet:${NC}"
echo "   ${YELLOW}STELLAR_NETWORK=mainnet${NC}"
echo "   ${YELLOW}SOROBAN_CONTRACT_ID=YOUR_DEPLOYED_CONTRACT_ID${NC}"
echo "   ${YELLOW}SOROBAN_PUBLIC_KEY=YOUR_MAINNET_PUBLIC_KEY${NC}"
echo "   ${YELLOW}SOROBAN_SECRET_KEY=YOUR_MAINNET_SECRET_KEY${NC}"
echo "   ${YELLOW}USDC_ASSET_CODE=USDC${NC}"
echo "   ${YELLOW}USDC_ASSET_ISSUER=GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN${NC}"
echo ""

# 7. Verification
echo -e "${BLUE}7. Verification Commands:${NC}"
echo "   Check network status: ${YELLOW}curl http://localhost:3000/network${NC}"
echo "   Test balance check: ${YELLOW}curl http://localhost:3000/balance/YOUR_MAINNET_ADDRESS${NC}"
echo ""

echo -e "${GREEN}✅ Setup guide complete!${NC}"
echo ""
echo -e "${RED}⚠️  Remember: Mainnet uses real money. Test thoroughly on testnet first!${NC}"