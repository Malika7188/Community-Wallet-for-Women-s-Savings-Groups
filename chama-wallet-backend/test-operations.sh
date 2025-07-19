#!/bin/bash

# Script to test all group operations
# Usage: ./test-operations.sh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API base URL
BASE_URL="http://localhost:3000/api"

# Test data
TEST_GROUP_NAME="Test Chama $(date +%s)"
TEST_DESCRIPTION="Test savings group created by script"
TEST_WALLET="GA7LYFBRHPF3WOJTUCQIWC3RDRBORDCXVWAYWAXDG4BT2XIPKDEXNJXL"
TEST_USER_ADDRESS="GBVFW3I5RSN5GPMWUC73J7FUPPXX3QXLMDEGZXYMP56C36E6FQQLX6WU"

# Store group ID for subsequent tests
GROUP_ID=""
CONTRACT_ID=""

# Function to make API request and display results
api_test() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local description="$4"

    echo -e "${BLUE}🧪 Testing: $description${NC}"
    echo -e "${YELLOW}$method $endpoint${NC}"
    
    if [[ -n "$data" ]]; then
        echo -e "${YELLOW}Data: $data${NC}"
    fi
    
    echo "----------------------------------------"

    # Make the request
    local response
    if [[ -n "$data" ]]; then
        response=$(curl -s -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data" \
            -w "\nHTTP_STATUS:%{http_code}")
    else
        response=$(curl -s -X "$method" "$BASE_URL$endpoint" \
            -w "\nHTTP_STATUS:%{http_code}")
    fi

    # Parse response
    local http_status=$(echo "$response" | grep "HTTP_STATUS" | cut -d: -f2)
    local json_response=$(echo "$response" | sed '/HTTP_STATUS/d')

    # Display results
    if [[ "$http_status" =~ ^2[0-9]{2}$ ]]; then
        echo -e "${GREEN}✅ Success (HTTP $http_status)${NC}"
        
        # Pretty print JSON if jq is available
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

# Function to extract group ID from response
extract_group_id() {
    local response="$1"
    if command -v jq &> /dev/null; then
        echo "$response" | jq -r '.group.ID // .ID // empty'
    else
        # Fallback parsing without jq
        echo "$response" | grep -o '"ID":"[^"]*"' | head -1 | cut -d'"' -f4
    fi
}

# Function to extract contract ID from response
extract_contract_id() {
    local response="$1"
    if command -v jq &> /dev/null; then
        echo "$response" | jq -r '.group.ContractID // .ContractID // .contract_id // empty'
    else
        # Fallback parsing without jq
        echo "$response" | grep -o '"ContractID":"[^"]*"' | head -1 | cut -d'"' -f4
    fi
}

# Test 1: Create a group
test_create_group() {
    local data="{
        \"name\": \"$TEST_GROUP_NAME\",
        \"description\": \"$TEST_DESCRIPTION\",
        \"wallet\": \"$TEST_WALLET\"
    }"
    
    local response_body
    response_body=$(curl -s -X POST "$BASE_URL/groups" \
        -H "Content-Type: application/json" \
        -d "$data")
    
    # Try to get response and parse it
    api_test "POST" "/groups" "$data" "Create Group"
    
    # Extract group ID for subsequent tests
    GROUP_ID=$(extract_group_id "$response_body")
    CONTRACT_ID=$(extract_contract_id "$response_body")
    
    if [[ -n "$GROUP_ID" ]]; then
        echo -e "${GREEN}📝 Group ID saved: $GROUP_ID${NC}"
    fi
    
    if [[ -n "$CONTRACT_ID" ]]; then
        echo -e "${GREEN}📝 Contract ID saved: $CONTRACT_ID${NC}"
    fi
    
    echo ""
}

# Test 2: Get all groups
test_get_all_groups() {
    api_test "GET" "/groups" "" "Get All Groups"
}

# Test 3: Get specific group
test_get_group() {
    if [[ -n "$GROUP_ID" ]]; then
        api_test "GET" "/groups/$GROUP_ID" "" "Get Specific Group"
    else
        echo -e "${YELLOW}⚠️  Skipping get specific group test - no group ID available${NC}"
        echo ""
    fi
}

# Test 4: Contribute to group
test_contribute() {
    if [[ -n "$GROUP_ID" ]]; then
        local data="{
            \"from\": \"$TEST_USER_ADDRESS\",
            \"secret\": \"test-secret\",
            \"amount\": \"100\"
        }"
        api_test "POST" "/groups/$GROUP_ID/contribute" "$data" "Contribute to Group"
    else
        echo -e "${YELLOW}⚠️  Skipping contribution test - no group ID available${NC}"
        echo ""
    fi
}

# Test 5: Direct Soroban contribute
test_soroban_contribute() {
    if [[ -n "$CONTRACT_ID" ]]; then
        local data="{
            \"user_address\": \"$TEST_USER_ADDRESS\",
            \"amount\": \"50\"
        }"
        api_test "POST" "/contribute" "$data" "Direct Soroban Contribute"
    else
        echo -e "${YELLOW}⚠️  Skipping Soroban contribute test - no contract ID available${NC}"
        echo ""
    fi
}

# Test 6: Check balance
test_check_balance() {
    if [[ -n "$CONTRACT_ID" ]]; then
        local data="{
            \"contract_id\": \"$CONTRACT_ID\",
            \"user_address\": \"$TEST_USER_ADDRESS\"
        }"
        api_test "POST" "/balance" "$data" "Check Balance"
    else
        echo -e "${YELLOW}⚠️  Skipping balance check - no contract ID available${NC}"
        echo ""
    fi
}

# Main test runner
main() {
    echo -e "${BLUE}🏛️  Chama Group Operations Tester${NC}"
    echo "===================================="
    echo ""
    
    # Check if API is running
    echo -e "${BLUE}🔍 Checking API availability...${NC}"
    if ! curl -s --max-time 3 "$BASE_URL/../" > /dev/null; then
        echo -e "${RED}❌ API is not responding. Please start your backend server.${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ API is running${NC}"
    echo ""
    
    # Run tests
    echo -e "${YELLOW}Running test suite...${NC}"
    echo ""
    
    test_create_group
    test_get_all_groups
    test_get_group
    test_contribute
    test_soroban_contribute
    test_check_balance
    
    # Summary
    echo -e "${BLUE}🎯 Test Summary${NC}"
    echo "==============="
    echo -e "${GREEN}✅ Test suite completed${NC}"
    
    if [[ -n "$GROUP_ID" ]]; then
        echo -e "${BLUE}📋 Created Group ID: $GROUP_ID${NC}"
    fi
    
    if [[ -n "$CONTRACT_ID" ]]; then
        echo -e "${BLUE}📋 Contract ID: $CONTRACT_ID${NC}"
    fi
    
    echo ""
    echo -e "${YELLOW}💡 You can now use these IDs to test your frontend integration${NC}"
}

# Handle command line arguments
case "$1" in
    "--help"|"-h")
        echo "Usage: $0 [options]"
        echo ""
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --quick, -q    Run only basic tests"
        echo ""
        echo "This script tests all group operations including:"
        echo "- Creating a group"
        echo "- Fetching groups"
        echo "- Contributing to groups"
        echo "- Soroban contract interactions"
        exit 0
        ;;
    "--quick"|"-q")
        echo -e "${YELLOW}Running quick tests only...${NC}"
        test_create_group
        test_get_all_groups
        ;;
    *)
        main
        ;;
esac