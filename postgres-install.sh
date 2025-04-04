#!/bin/bash

set -e

PG_USER="sumit_kumar2"
PG_PASSWORD="12345"
PG_DATABASE="polling_app_db"

echo "### Updating system packages"
sudo yum update -y

echo "### Installing required packages"
sudo yum install -y epel-release wget firewalld

echo "### Starting firewalld if not already running"
sudo systemctl enable --now firewalld

echo "### Installing PostgreSQL repository"
sudo yum install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm

echo "### Installing PostgreSQL server and client"
sudo yum install -y postgresql-server postgresql-contrib

echo "### Initializing PostgreSQL database"
sudo postgresql-setup initdb

echo "### Enabling and starting PostgreSQL service"
sudo systemctl enable postgresql
sudo systemctl start postgresql


# === PostgreSQL User, Database, and Tables ===
echo "### Creating PostgreSQL user and database"

# 1. Create the user if not exists
sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname = '$PG_USER'" | grep -q 1 || \
sudo -u postgres psql -c "CREATE USER $PG_USER WITH PASSWORD '$PG_PASSWORD';"

# 2. Create the database owned by that user
sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname = '$PG_DATABASE'" | grep -q 1 || \
sudo -u postgres psql -c "CREATE DATABASE $PG_DATABASE OWNER $PG_USER;"

# 3. Ensure role has full privileges
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $PG_DATABASE TO $PG_USER;"

# 4. Set default search_path
sudo -u postgres psql -d $PG_DATABASE -c "ALTER ROLE $PG_USER SET search_path TO public;"

# 5. Create tables as the new user (to ensure ownership)
echo "### Creating tables in $PG_DATABASE with $PG_USER as owner"

TABLE_SQL="
SET ROLE $PG_USER;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS polls (
    id SERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    yes_count INT DEFAULT 0,
    no_count INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    poll_id INT NOT NULL,
    username VARCHAR(50) NOT NULL,
    vote BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (poll_id) REFERENCES polls (id) ON DELETE CASCADE
);
"

echo "$TABLE_SQL" | sudo -u postgres psql -d "$PG_DATABASE" -v ON_ERROR_STOP=1

echo "### PostgreSQL installation and table setup complete!"


# Expose the port from outside
PG_CONF="/var/lib/pgsql/data/postgresql.conf"
sudo sed -i "s/^#listen_addresses = .*/listen_addresses = '*'/g" "$PG_CONF"

PG_HBA="/var/lib/pgsql/data/pg_hba.conf"
echo "host    all             all             0.0.0.0/0               md5" | sudo tee -a "$PG_HBA"

sudo sed -i 's/^\(local[[:space:]]\+all[[:space:]]\+all[[:space:]]\+\)peer/\1md5/' /var/lib/pgsql/data/pg_hba.conf

sudo systemctl restart postgresql

sudo firewall-cmd --permanent --add-port=5432/tcp

sudo firewall-cmd --reload
