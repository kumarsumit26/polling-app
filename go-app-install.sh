#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# === Variables ===
GO_VERSION="1.21.5"
APP_DIR="/opt/polling-app/backend"  # Change to your backend project path
MAIN_FILE="main.go"                 # Entry point of your Go app
BACKEND_IP=$(hostname -I | awk '{print $1}') # Automatically get the server's IP address

echo "### Updating system packages"
sudo yum update -y

echo "### Installing dependencies for Go and Git"
sudo yum install -y wget git tar bash-completion

echo "### Removing any existing Go installation"
sudo rm -rf /usr/local/go

echo "### Downloading Go $GO_VERSION"
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -P /tmp

echo "### Installing Go"
sudo tar -C /usr/local -xzf /tmp/go${GO_VERSION}.linux-amd64.tar.gz

echo "### Setting up Go environment variables"
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh
export PATH=$PATH:/usr/local/go/bin

echo "### Verifying Go installation"
go version

git clone https://github.com/kumarsumit26/polling-app.git
cd polling-app/backend

echo "### Initializing Go modules and downloading dependencies"
go mod tidy
go mod vendor

export DB_HOST=@@{Postgres.address}@@
export DB_USER=sumit_kumar2

echo "### Running the backend server"
nohup go run $MAIN_FILE > backend.log 2>&1 &


echo "Installing Node.js and npm..."
curl -sL https://rpm.nodesource.com/setup_16.x | sudo bash -
sudo yum install -y nodejs

# Verify Node.js and npm installation
echo "Verifying Node.js and npm installation..."
node -v
npm -v

cd ../ui

echo "REACT_APP_BACKEND_URL=http://$BACKEND_IP:8080" > .env

# Install dependencies
echo "Installing dependencies..."
npm install

# Start the React development server
echo "Starting the React development server..."
setsid env BROWSER=none HOST=0.0.0.0 PORT=3000 npx react-scripts start > frontend.log 2>&1 < /dev/null &

sleep 20

echo "### Setup complete"
echo "Backend is running on http://$BACKEND_IP:8080"
echo "Frontend is running on http://$BACKEND_IP:3000"
