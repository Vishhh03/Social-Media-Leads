#!/bin/bash
# ============================================
# LeadPilot — VPS Setup Script
# Run: curl -fsSL <your-repo>/setup.sh | bash
# Or:  chmod +x setup.sh && ./setup.sh
# ============================================

set -e

echo "🚀 LeadPilot VPS Setup"
echo "======================"

# 1. Update system
echo "📦 Updating system packages..."
sudo apt-get update -y && sudo apt-get upgrade -y

# 2. Install Docker if not present
if ! command -v docker &> /dev/null; then
  echo "🐳 Installing Docker..."
  curl -fsSL https://get.docker.com | sh
  sudo usermod -aG docker $USER
  echo "✅ Docker installed"
else
  echo "✅ Docker already installed"
fi

# 3. Install Docker Compose plugin if not present
if ! docker compose version &> /dev/null; then
  echo "🐳 Installing Docker Compose..."
  sudo apt-get install -y docker-compose-plugin
  echo "✅ Docker Compose installed"
else
  echo "✅ Docker Compose already installed"
fi

# 4. Create .env if not exists
if [ ! -f .env ]; then
  echo "📄 Creating .env from template..."
  cp .env.example .env

  # Generate a random JWT secret
  JWT=$(openssl rand -hex 32)
  sed -i "s/CHANGE_ME_generate_a_random_64_char_string/$JWT/" .env

  # Generate a random DB password
  DBPASS=$(openssl rand -hex 16)
  sed -i "s/CHANGE_ME_strong_password_here/$DBPASS/" .env

  echo "✅ .env created with random secrets"
  echo "⚠️  Edit .env to add your Meta API tokens"
else
  echo "✅ .env already exists"
fi

# 5. Create required directories
mkdir -p nginx certbot/conf certbot/www

# 6. Build and start
echo "🔨 Building and starting services..."
docker compose up -d --build

echo ""
echo "============================================"
echo "✅ LeadPilot is running!"
echo "============================================"
echo ""
echo "📌 Dashboard: http://$(hostname -I | awk '{print $1}')"
echo "📌 API:       http://$(hostname -I | awk '{print $1}')/api/v1/health"
echo ""
echo "Next steps:"
echo "  1. Open the dashboard and create an account"
echo "  2. Connect your Meta channels (WhatsApp/Instagram/Facebook)"
echo "  3. Set up automation rules"
echo ""
echo "To add SSL later:"
echo "  1. Point your domain to this server's IP"
echo "  2. Set DOMAIN in .env"
echo "  3. Run: sudo certbot certonly --webroot -w ./certbot/www -d yourdomain.com"
echo "  4. Uncomment the SSL block in nginx/nginx.conf"
echo "  5. Restart: docker compose restart nginx"
