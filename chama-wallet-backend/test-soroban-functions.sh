#!/bin/bash

# Test script for Soroban contract functions
# Usage: ./test-soroban-functions.sh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="http://localhost:3000/api"
CONTRACT_ID="CADHKUC557DJ2F2XGEO4BGHFIYQ6O5QDVNG637ANRAGPBSWXMXXPMOI4"
TEST_USER_ADDRESS="GBVFW3I5RSN5GPMWUC73J7FUPPXX3QXLMDEGZXYMP56C36E6FQQLX6WU"
TEST_SECRET_KEY="SAHUN4JW2MYG4IGCBPGFGPBKWS7CV43V54SVWL4AZSKLDB7FQBMDWFVL"

echo -e "${BLUE}🧪 Testing Soroban Contract Functions${NC}"
echo "====================================="
echo ""

# Function to test API endpoint
test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local description="$4"

    echo -e "${BLUE}Testing: $description${NC}"
    echo -e "${YELLOW}$method $endpoint${NC}"
    
    if [[ -n "$data" ]]; then
        echo -e "${YELLOW}Data: $data${NC}"
    fi
    
    echo "----------------------------------------"

    local response
    if [[ -n "$data" ]]; then
        response=$(curl -s -X "$method" "$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data" \
            -w "\nHTTP_STATUS:%{http_code}")
    else
        response=$(curl -s -X "$method" "$endpoint" \
            -w "\nHTTP_STATUS:%{http_code}")
    fi

    local http_status=$(echo "$response" | grep "HTTP_STATUS" | cut -d: -f2)
    local json_response=$(echo "$response" | sed '/HTTP_STATUS/d')

    if [[ "$http_status" =~ ^2[0-9]{2}$ ]]; then
        echo -e "${GREEN}✅ Success (HTTP $http_status)${NC}"
        if command -v jq &> /dev/null; then
            echo "$json_response" | jq .
        else
            echo "$json_response"
        fi
    else
        echo -e "${RED}❌ Failed (HTTP $http_status)${NC}"
        echo "$json_response"
    fi
    
    echo ""
    return $http_status
}

# Check if API is running
echo -e "${BLUE}🔍 Checking API availability...${NC}"
if ! curl -s --max-time 3 "$API_URL/../" > /dev/null; then
    echo -e "${RED}❌ API is not responding. Please start your backend server.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ API is running${NC}"
echo ""

# Test 1: Contribute to contract
echo -e "${YELLOW}Test 1: Direct Contract Contribution${NC}"
contribute_data="{
    \"contract_id\": \"$CONTRACT_ID\",
    \"user_address\": \"$TEST_USER_ADDRESS\",
    \"amount\": \"100\",
    \"secret_key\": \"$TEST_SECRET_KEY\"
}"
test_endpoint "POST" "$API_URL/contribute" "$contribute_data" "Direct Soroban Contribution"

# Test 2: Check balance
echo -e "${YELLOW}Test 2: Check Contract Balance${NC}"
balance_data="{
    \"contract_id\": \"$CONTRACT_ID\",
    \"user_address\": \"$TEST_USER_ADDRESS\"
}"
test_endpoint "POST" "$API_URL/balance" "$balance_data" "Check Contract Balance"

# Test 3: Get contribution history
echo -e "${YELLOW}Test 3: Get Contribution History${NC}"
history_data="{
    \"contract_id\": \"$CONTRACT_ID\",
    \"user_address\": \"$TEST_USER_ADDRESS\"
}"
test_endpoint "POST" "$API_URL/history" "$history_data" "Get Contribution History"

# Test 4: Test withdrawal
echo -e "${YELLOW}Test 4: Test Withdrawal${NC}"
withdraw_data="{
    \"contract_id\": \"$CONTRACT_ID\",
    \"user_address\": \"$TEST_USER_ADDRESS\",
    \"amount\": \"50\",
    \"secret_key\": \"$TEST_SECRET_KEY\"
}"
test_endpoint "POST" "$API_URL/withdraw" "$withdraw_data" "Test Withdrawal"

# Test 5: Check balance after withdrawal
echo -e "${YELLOW}Test 5: Check Balance After Withdrawal${NC}"
test_endpoint "POST" "$API_URL/balance" "$balance_data" "Check Balance After Withdrawal"

echo -e "${BLUE}🎯 Test Summary${NC}"
echo "==============="
echo -e "${GREEN}✅ Soroban function tests completed${NC}"
echo ""
echo -e "${YELLOW}💡 Next steps:${NC}"
echo "1. Test group contributions through the group API"
echo "2. Test wallet-to-wallet transfers"
echo "3. Test payout functionality"
echo ""