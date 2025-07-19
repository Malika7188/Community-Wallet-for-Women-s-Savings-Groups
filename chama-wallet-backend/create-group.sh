#!/bin/bash

# Script to create a new Chama group via the API
# Usage: ./create-group.sh [group_name] [description] [wallet_address]

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default API endpoint
API_URL="http://localhost:3000/api/groups"

# Function to display usage
show_usage() {
    echo -e "${BLUE}Usage: $0 [group_name] [description] [wallet_address]${NC}"
    echo -e "${BLUE}   or: $0 --interactive${NC}"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 \"Sar Queens\" \"Women-led savings group\" \"GA...WALLET\""
    echo "  $0 --interactive"
    echo ""
    exit 1
}

# Function to validate Stellar address
validate_stellar_address() {
    local address=$1
    if [[ ! $address =~ ^G[A-Z2-7]{55}$ ]]; then
        echo -e "${RED}❌ Invalid Stellar address format${NC}"
        echo -e "${YELLOW}Stellar addresses should start with 'G' and be 56 characters long${NC}"
        return 1
    fi
    return 0
}

# Function to create group via API
create_group() {
    local name="$1"
    local description="$2"
    local wallet="$3"

    echo -e "${BLUE}Creating group: ${name}${NC}"
    echo -e "${BLUE}Description: ${description}${NC}"
    echo -e "${BLUE}Wallet: ${wallet}${NC}"
    echo ""

    # Create JSON payload
    local json_payload=$(cat <<EOF
{
    "name": "$name",
    "description": "$description",
    "wallet": "$wallet"
}
EOF
)

    # Make API request
    echo -e "${YELLOW}Sending request to API...${NC}"
    local response=$(curl -s -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d "$json_payload" \
        -w "\nHTTP_STATUS:%{http_code}")

    # Parse response and status
    local http_status=$(echo "$response" | grep "HTTP_STATUS" | cut -d: -f2)
    local json_response=$(echo "$response" | sed '/HTTP_STATUS/d')

    # Check if request was successful
    if [ "$http_status" = "200" ] || [ "$http_status" = "201" ]; then
        echo -e "${GREEN}✅ Group created successfully!${NC}"
        echo ""
        
        # Parse and display response using jq if available
        if command -v jq &> /dev/null; then
            echo -e "${GREEN}📋 Group Details:${NC}"
            echo "$json_response" | jq -r '
                "🆔 ID: " + .group.ID + 
                "\n📛 Name: " + .group.Name + 
                "\n📝 Description: " + .group.Description + 
                "\n💳 Wallet: " + .group.Wallet + 
                "\n🔗 Contract ID: " + .group.ContractID'
        else
            echo -e "${GREEN}📋 Raw Response:${NC}"
            echo "$json_response"
        fi
    else
        echo -e "${RED}❌ Failed to create group (HTTP $http_status)${NC}"
        echo -e "${RED}Response: $json_response${NC}"
        exit 1
    fi
}

# Interactive mode
interactive_mode() {
    echo -e "${BLUE}🏛️  Chama Group Creation (Interactive Mode)${NC}"
    echo "=============================================="
    echo ""

    # Get group name
    read -p "Enter group name: " group_name
    if [[ -z "$group_name" ]]; then
        echo -e "${RED}❌ Group name cannot be empty${NC}"
        exit 1
    fi

    # Get description
    read -p "Enter group description: " group_description
    if [[ -z "$group_description" ]]; then
        group_description="Chama savings group powered by Stellar"
    fi

    # Get wallet address
    echo -e "${YELLOW}💡 You can generate a test wallet at: https://lab.stellar.org/account-creator${NC}"
    read -p "Enter Stellar wallet address (starts with G): " wallet_address
    
    # Validate wallet address
    if ! validate_stellar_address "$wallet_address"; then
        exit 1
    fi

    # Confirm details
    echo ""
    echo -e "${YELLOW}Please confirm the details:${NC}"
    echo "Name: $group_name"
    echo "Description: $group_description"
    echo "Wallet: $wallet_address"
    echo ""
    read -p "Create this group? (y/N): " confirm

    if [[ $confirm =~ ^[Yy]$ ]]; then
        create_group "$group_name" "$group_description" "$wallet_address"
    else
        echo -e "${YELLOW}Operation cancelled${NC}"
        exit 0
    fi
}

# Generate sample wallet address for testing
generate_sample_wallet() {
    echo -e "${BLUE}🎲 Generating sample wallet address...${NC}"
    if command -v soroban &> /dev/null; then
        soroban keys generate --global temp-wallet-$(date +%s) 2>/dev/null
        local sample_address=$(soroban keys address temp-wallet-$(date +%s) 2>/dev/null || echo "GA7LYFBRHPF3WOJTUCQIWC3RDRBORDCXVWAYWAXDG4BT2XIPKDEXNJXL")
        echo -e "${GREEN}Sample address: $sample_address${NC}"
        echo -e "${YELLOW}💡 Use this for testing, or generate your own at: https://lab.stellar.org/account-creator${NC}"
    else
        echo -e "${GREEN}Sample address: GA7LYFBRHPF3WOJTUCQIWC3RDRBORDCXVWAYWAXDG4BT2XIPKDEXNJXL${NC}"
        echo -e "${YELLOW}💡 This is a sample address. Generate your own at: https://lab.stellar.org/account-creator${NC}"
    fi
}

# Check if API is running
check_api() {
    echo -e "${BLUE}🔍 Checking if API is running...${NC}"
    if curl -s --max-time 3 "http://localhost:3000" > /dev/null; then
        echo -e "${GREEN}✅ API is running${NC}"
    else
        echo -e "${RED}❌ API is not responding. Please make sure your backend is running on port 3000${NC}"
        echo -e "${YELLOW}Start your backend with: go run main.go${NC}"
        exit 1
    fi
}

# Main script logic
main() {
    echo -e "${BLUE}🏛️  Chama Group Creator${NC}"
    echo "======================="
    echo ""

    # Check if API is running
    check_api
    echo ""

    # Check command line arguments
    case "$1" in
        "--help"|"-h")
            show_usage
            ;;
        "--interactive"|"-i")
            interactive_mode
            ;;
        "--sample-wallet")
            generate_sample_wallet
            ;;
        "")
            interactive_mode
            ;;
        *)
            if [[ $# -eq 3 ]]; then
                # Validate wallet address
                if ! validate_stellar_address "$3"; then
                    exit 1
                fi
                create_group "$1" "$2" "$3"
            else
                echo -e "${RED}❌ Invalid number of arguments${NC}"
                show_usage
            fi
            ;;
    esac
}

# Run main function
main "$@"