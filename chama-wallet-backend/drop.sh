#!/bin/bash
echo "🗑️  Resetting database..."
sudo -u postgres psql -d chama_wallet "
DROP TABLE IF EXISTS contributions CASCADE;
DROP TABLE IF EXISTS members CASCADE;
DROP TABLE IF EXISTS group_invitations CASCADE;
DROP TABLE IF EXISTS admin_nominations CASCADE;
DROP TABLE IF EXISTS payout_approvals CASCADE;
DROP TABLE IF EXISTS payout_requests CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS groups CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS payout_schedules CASCADE;
" 

echo "🔄 Setting up database permissions..."
sudo -u postgres psql -d chama_wallet "
GRANT ALL ON SCHEMA public TO chama_user;
GRANT CREATE ON SCHEMA public TO chama_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO chama_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO chama_user;
"

echo "✅ Database reset complete. Restart your server to recreate tables."