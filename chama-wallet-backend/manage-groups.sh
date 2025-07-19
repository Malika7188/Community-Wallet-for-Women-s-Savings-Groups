#!/bin/bash

# Group Management Script for Chama Wallet
# Provides easy commands to manage your groups

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# API Configuration
BASE_URL="http://localhost:3000/api"

# Functions
show_help() {
    echo -e "${BLUE}🏛️  Chama Group Manager${NC}"
    echo "======================="
    echo ""
    echo -e "${YELLOW}Usage:${NC}"
    echo "  $0 list                    - List all groups"
    echo "  $0 show <group_id>         - Show specific group details"
    echo "  $0 create                  - Create a new group (interactive)"
    echo "  $0 contribute <group_id>   - Contribute to a group (interactive)"
    echo "  $0 balance <contract_id>   - Check balance for a contract"
    echo "  $0 status                  - Check API status"
    echo ""
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 list"
    echo "  $0 show a0939bd4-6931-45ff-970d-c04e856bfad3"
    echo "  $0 create"
    echo ""
}

check_api_status() {
    if curl -s --max-time 3 "$BASE_URL/../" > /dev/null; then
        echo -e "${GREEN}✅ API is running${NC}"
        return 0
    else
        echo -e "${RED}❌ API is not responding${NC}"
        echo -e "${YELLOW}💡 Start your backend with: go run main.go${NC}"
        return 1
    fi
}

list_groups() {
    echo -e "${BLUE}📋 Listing all groups...${NC}"
    echo ""
    
    local response=$(curl -s "$BASE_URL/groups")
    
    if command -v jq &> /dev/null; then
        # Pretty format with jq
        echo "$response" | jq -r '.[] | "🏛️  \(.Name)\n   ID: \(.ID)\n   Wallet: \(.Wallet)\n   Contract: \(.ContractID)\n   Description: \(.Description)\n"'
    else
        # Basic formatting without jq
        echo "$response"
    fi
}

show_group() {
    local group_id="$1"
    
    if [[ -z "$group_id" ]]; then
        echo -e "${RED}❌ Please provide a group ID${NC}"
        echo "Usage: $0 show <group_id>"
        return 1
    fi
    
    echo -e "${BLUE}🔍 Fetching group details for: $group_id${NC}"
    echo ""
    
    local response=$(curl -s "$BASE_URL/groups/$group_id")
    
    if command -v jq &> /dev/null; then
        # Check if response contains error
        if echo "$response" | jq -e '.error' > /dev/null; then
            echo -e "${RED}❌ Error: $(echo "$response" | jq -r '.error')${NC}"
            return 1
        fi
        
        # Pretty format group details
        echo -e "${GREEN}📋 Group Details:${NC}"
        echo "$response" | jq -r '
            "🏛️  Name: " + .Name +
            "\n🆔 ID: " + .ID +
            "\n📝 Description: " + .Description +
            "\n💳 Wallet: " + .Wallet +
            "\n🔗 Contract ID: " + .ContractID'
    else
        echo "$response"
    fi
}

interactive_create() {
    echo -e "${BLUE}🏛️  Create New Group${NC}"
    echo "===================="
    echo ""
    
    # Get group details
    read -p "Enter group name: " name
    read -p "Enter description: " description
    read -p "Enter wallet address (G...): " wallet
    
    # Validate inputs
    if [[ -z "$name" || -z "$wallet" ]]; then
        echo -e "${RED}❌ Name and wallet address are required${NC}"
        return 1
    fi
    
    # Create JSON payload
    local payload=$(cat <<EOF
{
    "name": "$name",
    "description": "$description",
    "wallet": "$wallet"
}
EOF
)
    
    echo -e "${YELLOW}Creating group...${NC}"
    local response=$(curl -s -X POST "$BASE_URL/groups" \
        -H "Content-Type: application/json" \
        -d "$payload")
    
    if command -v jq &> /dev/null; then
        if echo "$response" | jq -e '.error' > /dev/null; then
            echo -e "${RED}❌ Error: $(echo "$response" | jq -r '.error')${NC}"
            return 1
        fi
        
        echo -e "${GREEN}✅ Group created successfully!${NC}"
        echo ""
        echo "$response" | jq -r '.group | 
            "🏛️  Name: " + .Name +
            "\n🆔 ID: " + .ID +
            "\n💳 Wallet: " + .Wallet +
            "\n🔗 Contract ID: " + .ContractID'
    else
        echo "$response"
    fi
}

interactive_contribute() {
    local group_id="$1"
    
    if [[ -z "$group_id" ]]; then
        echo -e "${RED}❌ Please provide a group ID${NC}"
        echo "Usage: $0 contribute <group_id>"
        return 1
    fi
    
    echo -e "${BLUE}💰 Contribute to Group${NC}"
    echo "====================="
    echo ""
    
    # Get contribution details
    read -p "Enter your address (G...): " from_address
    read -p "Enter amount to contribute: " amount
    read -p "Enter your secret key (optional): " secret
    
    if [[ -z "$from_address" || -z "$amount" ]]; then
        echo -e "${RED}❌ Address and amount are required${NC}"
        return 1
    fi
    
    # Create JSON payload
    local payload=$(cat <<EOF
{
    "from": "$from_address",
    "secret": "$secret",
    "amount": "$amount"
}
EOF
)
    
    echo -e "${YELLOW}Processing contribution...${NC}"
    local response=$(curl -s -X POST "$BASE_URL/groups/$group_id/contribute" \
        -H "Content-Type: application/json" \
        -d "$payload")
    
    if command -v jq &> /dev/null; then
        if echo "$response" | jq -e '.error' > /dev/null; then
            echo -e "${RED}❌ Error: $(echo "$response" | jq -r '.error')${NC}"
            return 1
        fi
        
        echo -e "${GREEN}✅ Contribution successful!${NC}"
        echo "$response" | jq .
    else
        echo "$response"
    fi
}

check_balance() {
    local contract_id="$1"
    
    if [[ -z "$contract_id" ]]; then
        echo -e "${RED}❌ Please provide a contract ID${NC}"
        echo "Usage: $0 balance <contract_id>"
        return 1
    fi
    
    echo -e "${BLUE}💰 Checking Balance${NC}"
    echo "=================="
    echo ""
    
    read -p "Enter user address to check balance for: " user_address
    
    if [[ -z "$user_address" ]]; then
        echo -e "${RED}❌ User address is required${NC}"
        return 1
    fi
    
    local payload=$(cat <<EOF
{
    "contract_id": "$contract_id",
    "user_address": "$user_address"
}
EOF
)
    
    local response=$(curl -s -X POST "$BASE_URL/balance" \
        -H "Content-Type: application/json" \
        -d "$payload")
    
    if command -v jq &> /dev/null; then
        if echo "$response" | jq -e '.error' > /dev/null; then
            echo -e "${RED}❌ Error: $(echo "$response" | jq -r '.error')${NC}"
            return 1
        fi
        
        echo -e "${GREEN}💰 Balance: $(echo "$response" | jq -r '.balance')${NC}"
        echo -e "${BLUE}👤 User: $(echo "$response" | jq -r '.user')${NC}"
    else
        echo "$response"
    fi
}

# Main script logic
main() {
    # Check API first
    if ! check_api_status; then
        exit 1
    fi
    
    case "$1" in
        "list"|"ls")
            list_groups
            ;;
        "show"|"get")
            show_group "$2"
            ;;
        "create"|"new")
            interactive_create
            ;;
        "contribute"|"pay")
            interactive_contribute "$2"
            ;;
        "balance"|"bal")
            check_balance "$2"
            ;;
        "status"|"health")
            echo -e "${GREEN}✅ API is healthy${NC}"
            ;;
        "help"|"--help"|"-h"|"")
            show_help
            ;;
        *)
            echo -e "${RED}❌ Unknown command: $1${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"