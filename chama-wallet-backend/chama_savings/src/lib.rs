#![no_std]

use soroban_sdk::{contract, contractimpl, symbol_short, Env, Symbol, Vec, Map, Address, log};

#[contract]  
pub struct ChamaSavings;

#[contractimpl]
impl ChamaSavings {
    /// Initialize the contract with proper setup
    pub fn initialize(env: Env) {
        log!(&env, "Initializing ChamaSavings contract");
        
        let init_key = symbol_short!("init");
        env.storage().persistent().set(&init_key, &true);
        
        // Initialize empty contributions and balances
        let contrib_key = symbol_short!("contrib");
        let balance_key = symbol_short!("balance");
        
        let empty_contributions: Vec<(Address, i128)> = Vec::new(&env);
        let empty_balances: Map<Address, i128> = Map::new(&env);
        
        env.storage().persistent().set(&contrib_key, &empty_contributions);
        env.storage().persistent().set(&balance_key, &empty_balances);
        
        log!(&env, "Contract initialized successfully");
    }

    /// Contribute function with comprehensive validation and logging
    pub fn contribute(env: Env, user: Address, amount: i128) {
        log!(&env, "Contribute called - User: {}, Amount: {}", user, amount);
        
        // Validate amount is positive
        if amount <= 0 {
            log!(&env, "Error: Amount must be positive, got: {}", amount);
            panic!("Amount must be positive");
        }

        // Require user authorization - this is critical for security
        user.require_auth();
        log!(&env, "User authorization successful for: {}", user);

        let contrib_key = symbol_short!("contrib");
        let balance_key = symbol_short!("balance");

        // Get existing contributions with error handling
        let mut contributions = env
            .storage()
            .persistent()
            .get::<Symbol, Vec<(Address, i128)>>(&contrib_key)
            .unwrap_or_else(|| {
                log!(&env, "No existing contributions found, creating new vector");
                Vec::new(&env)
            });

        // Get existing balances with error handling
        let mut balances = env
            .storage()
            .persistent()
            .get::<Symbol, Map<Address, i128>>(&balance_key)
            .unwrap_or_else(|| {
                log!(&env, "No existing balances found, creating new map");
                Map::new(&env)
            });

        // Add new contribution to history
        contributions.push_back((user.clone(), amount));
        log!(&env, "Added contribution to history. Total contributions: {}", contributions.len());

        // Update user balance
        let current_balance = balances.get(user.clone()).unwrap_or(0);
        let new_balance = current_balance + amount;
        balances.set(user.clone(), new_balance);
        
        log!(&env, "Updated balance for user {} from {} to {}", user, current_balance, new_balance);

        // Save updated data to storage
        env.storage().persistent().set(&contrib_key, &contributions);
        env.storage().persistent().set(&balance_key, &balances);
        
        log!(&env, "Contribution completed successfully");
    }

    /// Get total contributions for a user with logging
    pub fn get_balance(env: Env, user: Address) -> i128 {
        log!(&env, "Getting balance for user: {}", user);
        
        let balance_key = symbol_short!("balance");
        let balances = env
            .storage()
            .persistent()
            .get::<Symbol, Map<Address, i128>>(&balance_key)
            .unwrap_or_else(|| {
                log!(&env, "No balances found, returning empty map");
                Map::new(&env)
            });

        let balance = balances.get(user.clone()).unwrap_or(0);
        log!(&env, "Balance for user {}: {}", user, balance);
        
        balance
    }

    /// Get all contributions for transparency
    pub fn get_all_contributions(env: Env) -> Vec<(Address, i128)> {
        log!(&env, "Getting all contributions");
        
        let contrib_key = symbol_short!("contrib");
        let contributions = env
            .storage()
            .persistent()
            .get::<Symbol, Vec<(Address, i128)>>(&contrib_key)
            .unwrap_or_else(|| {
                log!(&env, "No contributions found, returning empty vector");
                Vec::new(&env)
            });
            
        log!(&env, "Total contributions found: {}", contributions.len());
        contributions
    }

    /// Get total pool amount across all users
    pub fn get_total_pool(env: Env) -> i128 {
        log!(&env, "Calculating total pool");
        
        let balance_key = symbol_short!("balance");
        let balances = env
            .storage()
            .persistent()
            .get::<Symbol, Map<Address, i128>>(&balance_key)
            .unwrap_or_else(|| {
                log!(&env, "No balances found for total pool calculation");
                Map::new(&env)
            });

        let mut total = 0i128;
        for (user, amount) in balances.iter() {
            total += amount;
            log!(&env, "Adding {} from user {} to total pool", amount, user);
        }
        
        log!(&env, "Total pool calculated: {}", total);
        total
    }

    /// Withdraw function for payouts (admin only in real implementation)
    pub fn withdraw(env: Env, user: Address, amount: i128) -> i128 {
        log!(&env, "Withdraw called - User: {}, Amount: {}", user, amount);
        
        // Validate amount is positive
        if amount <= 0 {
            log!(&env, "Error: Withdraw amount must be positive, got: {}", amount);
            panic!("Withdraw amount must be positive");
        }

        // Require user authorization
        user.require_auth();
        log!(&env, "User authorization successful for withdrawal: {}", user);

        let balance_key = symbol_short!("balance");
        let mut balances = env
            .storage()
            .persistent()
            .get::<Symbol, Map<Address, i128>>(&balance_key)
            .unwrap_or_else(|| Map::new(&env));

        let current_balance = balances.get(user.clone()).unwrap_or(0);
        
        // Check sufficient balance
        if current_balance < amount {
            log!(&env, "Error: Insufficient balance. Current: {}, Requested: {}", current_balance, amount);
            panic!("Insufficient balance for withdrawal");
        }

        let new_balance = current_balance - amount;
        balances.set(user.clone(), new_balance);
        
        // Save updated balances
        env.storage().persistent().set(&balance_key, &balances);
        
        log!(&env, "Withdrawal successful. New balance for {}: {}", user, new_balance);
        new_balance
    }

    /// Get contribution history for a specific user
    pub fn get_contribution_history(env: Env, user: Address) -> Vec<i128> {
        log!(&env, "Getting contribution history for user: {}", user);
        
        let contrib_key = symbol_short!("contrib");
        let contributions = env
            .storage()
            .persistent()
            .get::<Symbol, Vec<(Address, i128)>>(&contrib_key)
            .unwrap_or_else(|| Vec::new(&env));

        let mut user_contributions = Vec::new(&env);
        for (contributor, amount) in contributions.iter() {
            if contributor == user {
                user_contributions.push_back(amount);
            }
        }
        
        log!(&env, "Found {} contributions for user {}", user_contributions.len(), user);
        user_contributions
    }

    /// Check if contract is initialized
    pub fn is_initialized(env: Env) -> bool {
        let init_key = symbol_short!("init");
        env.storage().persistent().get(&init_key).unwrap_or(false)
    }

    /// Get contract statistics
    pub fn get_stats(env: Env) -> (i128, u32, u32) {
        let total_pool = Self::get_total_pool(env.clone());
        
        let contrib_key = symbol_short!("contrib");
        let contributions = env
            .storage()
            .persistent()
            .get::<Symbol, Vec<(Address, i128)>>(&contrib_key)
            .unwrap_or_else(|| Vec::new(&env));
            
        let balance_key = symbol_short!("balance");
        let balances = env
            .storage()
            .persistent()
            .get::<Symbol, Map<Address, i128>>(&balance_key)
            .unwrap_or_else(|| Map::new(&env));

        let total_contributions = contributions.len();
        let unique_contributors = balances.len();
        
        (total_pool, total_contributions, unique_contributors)
    }
}

#[cfg(test)]
mod test {
    use super::*;
    use soroban_sdk::{testutils::Address as _, Address, Env};

    #[test]
    fn test_initialize() {
        let env = Env::default();
        let contract_id = env.register(ChamaSavings, ());
        let client = ChamaSavingsClient::new(&env, &contract_id);

        client.initialize();
        assert!(client.is_initialized());
    }

    #[test]
    fn test_contribute_and_balance() {
        let env = Env::default();
        let contract_id = env.register(ChamaSavings, ());
        let client = ChamaSavingsClient::new(&env, &contract_id);

        let user = Address::generate(&env);
        
        client.initialize();
        
        // Mock authorization for testing
        env.mock_all_auths();
        
        client.contribute(&user, &1000);
        
        let balance = client.get_balance(&user);
        assert_eq!(balance, 1000);
    }

    #[test]
    fn test_multiple_contributions() {
        let env = Env::default();
        let contract_id = env.register(ChamaSavings, ());
        let client = ChamaSavingsClient::new(&env, &contract_id);

        let user1 = Address::generate(&env);
        let user2 = Address::generate(&env);
        
        client.initialize();
        env.mock_all_auths();
        
        client.contribute(&user1, &500);
        client.contribute(&user2, &300);
        client.contribute(&user1, &200);
        
        assert_eq!(client.get_balance(&user1), 700);
        assert_eq!(client.get_balance(&user2), 300);
        assert_eq!(client.get_total_pool(), 1000);
    }

    #[test]
    fn test_withdraw() {
        let env = Env::default();
        let contract_id = env.register(ChamaSavings, ());
        let client = ChamaSavingsClient::new(&env, &contract_id);

        let user = Address::generate(&env);
        
        client.initialize();
        env.mock_all_auths();
        
        client.contribute(&user, &1000);
        let new_balance = client.withdraw(&user, &300);
        
        assert_eq!(new_balance, 700);
        assert_eq!(client.get_balance(&user), 700);
    }
}