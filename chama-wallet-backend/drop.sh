#!/bin/bash
echo "🗑️  Resetting database..."
psql -U chama_user -d chama_wallet -c "
DROP TABLE IF EXISTS contributions CASCADE;
DROP TABLE IF EXISTS members CASCADE;
DROP TABLE IF EXISTS group_invitations CASCADE;
DROP TABLE IF EXISTS admin_nominations CASCADE;
DROP TABLE IF EXISTS payout_approvals CASCADE;
DROP TABLE IF EXISTS payout_requests CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS groups CASCADE;
DROP TABLE IF EXISTS users CASCADE;
"

echo "🔄 Setting up database permissions..."
sudo -u postgres psql -d chama_wallet -c "
GRANT ALL ON SCHEMA public TO chama_user;
GRANT CREATE ON SCHEMA public TO chama_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO chama_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO chama_user;
"

echo "✅ Database reset complete. Restart your server to recreate tables."