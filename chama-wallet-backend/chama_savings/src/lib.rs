#![no_std]

use soroban_sdk::{contract, contractimpl, symbol_short, Env, Symbol, Vec, Map, Address};

#[contract]  
pub struct ChamaSavings;

#[contractimpl]
impl ChamaSavings {
    // Initialize the contract
    pub fn initialize(env: Env) {
        // You might want to add initialization logic here
        let key = symbol_short!("init");
        env.storage().persistent().set(&key, &true);
    }

    // Contribute function with proper validation
    pub fn contribute(env: Env, user: Address, amount: i128) {
        // Validate amount is positive
        if amount <= 0 {
            panic!("Amount must be positive");
        }

        // Require user authorization
        user.require_auth();

        let key = symbol_short!("contrib");

        // Get existing contributions
        let mut contributions = env
            .storage()
            .persistent()
            .get::<Symbol, Vec<(Address, i128)>>(&key)
            .unwrap_or(Vec::new(&env));

        // Add new contribution
        contributions.push_back((user.clone(), amount));
        env.storage().persistent().set(&key, &contributions);

        // Optionally, maintain a separate balance mapping for efficiency
        let balance_key = symbol_short!("balance");
        let mut balances = env
            .storage()
            .persistent()
            .get::<Symbol, Map<Address, i128>>(&balance_key)
            .unwrap_or(Map::new(&env));

        let current_balance = balances.get(user.clone()).unwrap_or(0);
        balances.set(user, current_balance + amount);
        env.storage().persistent().set(&balance_key, &balances);
    }

    // Get total contributions for a user
    pub fn get_balance(env: Env, user: Address) -> i128 {
        let balance_key = symbol_short!("balance");
        let balances = env
            .storage()
            .persistent()
            .get::<Symbol, Map<Address, i128>>(&balance_key)
            .unwrap_or(Map::new(&env));

        balances.get(user).unwrap_or(0)
    }

    // Get all contributions (for transparency)
    pub fn get_all_contributions(env: Env) -> Vec<(Address, i128)> {
        let key = symbol_short!("contrib");
        env.storage()
            .persistent()
            .get::<Symbol, Vec<(Address, i128)>>(&key)
            .unwrap_or(Vec::new(&env))
    }

    // Get total pool amount
    pub fn get_total_pool(env: Env) -> i128 {
        let balance_key = symbol_short!("balance");
        let balances = env
            .storage()
            .persistent()
            .get::<Symbol, Map<Address, i128>>(&balance_key)
            .unwrap_or(Map::new(&env));

        let mut total = 0i128;
        for (_, amount) in balances.iter() {
            total += amount;
        }
        total
    }
}