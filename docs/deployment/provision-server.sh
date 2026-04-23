#!/usr/bin/env bash
# =============================================================================
# UmrohOS — Server Provisioning Script
# Fase 1: System Setup + Docker + Firewall + Security
# Fase 2: Nginx Reverse Proxy Config
#
# Jalankan dari terminal Mac Lutfi:
#   chmod +x provision-server.sh
#   ./provision-server.sh
#
# Server: infra@216.176.238.161
# SSH Key: /Users/lutfiaf/.ssh/lutfi_id_ed25519
# =============================================================================

set -euo pipefail

SSH_KEY="/Users/lutfiaf/.ssh/lutfi_id_ed25519"
SSH_USER="infra"
SSH_HOST="216.176.238.161"
SSH_OPTS="-i ${SSH_KEY} -o StrictHostKeyChecking=no -o BatchMode=yes"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info()    { echo -e "${BLUE}[INFO]${NC} $*"; }
log_success() { echo -e "${GREEN}[OK]${NC} $*"; }
log_warn()    { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error()   { echo -e "${RED}[ERROR]${NC} $*"; }

run_remote() {
    local description="$1"
    shift
    log_info "→ ${description}"
    ssh ${SSH_OPTS} "${SSH_USER}@${SSH_HOST}" "$@"
}

# =============================================================================
echo ""
echo "=============================================="
echo "  UmrohOS — Server Provisioning"
echo "  Target: ${SSH_USER}@${SSH_HOST}"
echo "=============================================="
echo ""

# Verify SSH key exists
if [ ! -f "${SSH_KEY}" ]; then
    log_error "SSH key tidak ditemukan: ${SSH_KEY}"
    exit 1
fi

# Test SSH connectivity
log_info "Testing SSH connectivity..."
if ssh ${SSH_OPTS} "${SSH_USER}@${SSH_HOST}" "echo 'SSH OK'" 2>/dev/null; then
    log_success "SSH connection berhasil"
else
    log_error "SSH connection GAGAL ke ${SSH_HOST}"
    log_error "Pastikan SSH key valid dan server dapat diakses"
    exit 1
fi

# =============================================================================
echo ""
echo "--- FASE 1: System Provisioning ---"
echo ""

# 1. Update system packages
run_remote "Update & upgrade system packages (ini mungkin butuh waktu beberapa menit)..." \
    "sudo DEBIAN_FRONTEND=noninteractive apt-get update -qq && \
     sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y -qq"
log_success "System packages updated"

# 2. Install prerequisites
run_remote "Install prerequisites (ca-certificates, curl, gnupg)..." \
    "sudo DEBIAN_FRONTEND=noninteractive apt-get install -y -qq \
        ca-certificates curl gnupg lsb-release apt-transport-https"
log_success "Prerequisites installed"

# 3. Install Docker Engine
run_remote "Setup Docker APT repository..." \
    "sudo install -m 0755 -d /etc/apt/keyrings && \
     sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg \
        -o /etc/apt/keyrings/docker.asc && \
     sudo chmod a+r /etc/apt/keyrings/docker.asc && \
     echo \"deb [arch=\$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] \
        https://download.docker.com/linux/ubuntu \
        \$(. /etc/os-release && echo \"\$VERSION_CODENAME\") stable\" | \
        sudo tee /etc/apt/sources.list.d/docker.list > /dev/null && \
     sudo apt-get update -qq"
log_success "Docker APT repository configured"

run_remote "Install Docker Engine + Compose plugin..." \
    "sudo DEBIAN_FRONTEND=noninteractive apt-get install -y -qq \
        docker-ce docker-ce-cli containerd.io \
        docker-buildx-plugin docker-compose-plugin"
log_success "Docker Engine installed"

run_remote "Add infra user to docker group..." \
    "sudo usermod -aG docker infra"
log_success "User 'infra' added to docker group"

run_remote "Enable & start Docker service..." \
    "sudo systemctl enable docker && sudo systemctl start docker"
log_success "Docker service enabled"

# 4. Setup UFW Firewall
run_remote "Install UFW..." \
    "sudo DEBIAN_FRONTEND=noninteractive apt-get install -y -qq ufw"

run_remote "Configure UFW rules..." \
    "sudo ufw --force reset && \
     sudo ufw default deny incoming && \
     sudo ufw default allow outgoing && \
     sudo ufw allow 22/tcp comment 'SSH' && \
     sudo ufw allow 80/tcp comment 'HTTP Web' && \
     sudo ufw allow 4000/tcp comment 'API Gateway' && \
     sudo ufw --force enable"
log_success "UFW firewall configured (22, 80, 4000 open)"

# 5. Disable SSH password authentication (key-only)
run_remote "Harden SSH: disable password auth, enable key-only..." \
    "sudo sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config && \
     sudo sed -i 's/^#*PubkeyAuthentication.*/PubkeyAuthentication yes/' /etc/ssh/sshd_config && \
     sudo sed -i 's/^#*ChallengeResponseAuthentication.*/ChallengeResponseAuthentication no/' /etc/ssh/sshd_config && \
     sudo systemctl restart ssh"
log_success "SSH hardened: password auth disabled"

# 6. Install fail2ban
run_remote "Install fail2ban..." \
    "sudo DEBIAN_FRONTEND=noninteractive apt-get install -y -qq fail2ban"

run_remote "Configure fail2ban for SSH..." \
    "sudo tee /etc/fail2ban/jail.local > /dev/null << 'FAIL2BAN_EOF'
[DEFAULT]
bantime  = 3600
findtime = 600
maxretry = 5
backend  = systemd

[sshd]
enabled  = true
port     = ssh
logpath  = %(sshd_log)s
maxretry = 3
FAIL2BAN_EOF
sudo systemctl enable fail2ban && sudo systemctl restart fail2ban"
log_success "fail2ban installed and configured"

# 7. Install Nginx
run_remote "Install Nginx..." \
    "sudo DEBIAN_FRONTEND=noninteractive apt-get install -y -qq nginx && \
     sudo systemctl enable nginx && sudo systemctl start nginx"
log_success "Nginx installed"

# 8. Create project directory
run_remote "Create project directory /home/infra/umrohos..." \
    "mkdir -p /home/infra/umrohos && chmod 750 /home/infra/umrohos"
log_success "Project directory created"

# =============================================================================
echo ""
echo "--- FASE 2: Nginx Configuration ---"
echo ""

# Nginx config: umrohos (port 80 → 3001)
run_remote "Create Nginx site config: umrohos (port 80)..." \
    "sudo tee /etc/nginx/sites-available/umrohos > /dev/null << 'NGINX_WEB_EOF'
server {
    listen 80;
    server_name 216.176.238.161;

    # Security headers
    add_header X-Frame-Options SAMEORIGIN;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection \"1; mode=block\";

    # Logs
    access_log /var/log/nginx/umrohos-access.log;
    error_log  /var/log/nginx/umrohos-error.log;

    location / {
        proxy_pass         http://localhost:3001;
        proxy_http_version 1.1;
        proxy_set_header   Host \$host;
        proxy_set_header   X-Real-IP \$remote_addr;
        proxy_set_header   X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto \$scheme;
        proxy_read_timeout 60s;
        proxy_connect_timeout 10s;

        # Return 503 with friendly message if upstream is down
        proxy_next_upstream error timeout http_502 http_503 http_504;
    }
}
NGINX_WEB_EOF"
log_success "Nginx config umrohos created"

# Nginx config: umrohos-api (port 4000 → 4001)
run_remote "Create Nginx site config: umrohos-api (port 4000)..." \
    "sudo tee /etc/nginx/sites-available/umrohos-api > /dev/null << 'NGINX_API_EOF'
server {
    listen 4000;
    server_name 216.176.238.161;

    # Logs
    access_log /var/log/nginx/umrohos-api-access.log;
    error_log  /var/log/nginx/umrohos-api-error.log;

    location / {
        proxy_pass         http://localhost:4001;
        proxy_http_version 1.1;
        proxy_set_header   Host \$host;
        proxy_set_header   X-Real-IP \$remote_addr;
        proxy_set_header   X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto \$scheme;
        proxy_read_timeout 60s;
        proxy_connect_timeout 10s;

        # WebSocket support untuk API
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection \"upgrade\";
    }
}
NGINX_API_EOF"
log_success "Nginx config umrohos-api created"

# Enable sites & remove default
run_remote "Enable Nginx sites, disable default..." \
    "sudo ln -sf /etc/nginx/sites-available/umrohos /etc/nginx/sites-enabled/umrohos && \
     sudo ln -sf /etc/nginx/sites-available/umrohos-api /etc/nginx/sites-enabled/umrohos-api && \
     sudo rm -f /etc/nginx/sites-enabled/default"
log_success "Nginx sites enabled"

run_remote "Test Nginx config and reload..." \
    "sudo nginx -t && sudo systemctl reload nginx"
log_success "Nginx reloaded with new config"

# =============================================================================
echo ""
echo "--- VERIFIKASI ---"
echo ""

# Docker version
DOCKER_VER=$(ssh ${SSH_OPTS} "${SSH_USER}@${SSH_HOST}" "docker --version 2>/dev/null || echo 'NOT INSTALLED'")
log_info "Docker: ${DOCKER_VER}"

COMPOSE_VER=$(ssh ${SSH_OPTS} "${SSH_USER}@${SSH_HOST}" "docker compose version 2>/dev/null || echo 'NOT INSTALLED'")
log_info "Docker Compose: ${COMPOSE_VER}"

# UFW status
UFW_STATUS=$(ssh ${SSH_OPTS} "${SSH_USER}@${SSH_HOST}" "sudo ufw status | head -5")
log_info "UFW Status:\n${UFW_STATUS}"

# Nginx status
NGINX_STATUS=$(ssh ${SSH_OPTS} "${SSH_USER}@${SSH_HOST}" "sudo systemctl is-active nginx")
log_info "Nginx: ${NGINX_STATUS}"

# fail2ban status
F2B_STATUS=$(ssh ${SSH_OPTS} "${SSH_USER}@${SSH_HOST}" "sudo systemctl is-active fail2ban")
log_info "fail2ban: ${F2B_STATUS}"

# HTTP response test (upstream not up yet, so 502 is expected and OK)
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time 10 "http://${SSH_HOST}/" 2>/dev/null || echo "timeout")
log_info "HTTP GET http://${SSH_HOST}/ → ${HTTP_CODE} (502 = Nginx OK, upstream belum running)"

# =============================================================================
echo ""
echo "=============================================="
log_success "PROVISIONING SELESAI!"
echo "=============================================="
echo ""
echo "Status server:"
echo "  - Docker Engine    : INSTALLED"
echo "  - Docker Compose   : INSTALLED (plugin)"
echo "  - UFW Firewall     : ACTIVE (22, 80, 4000 open)"
echo "  - SSH hardening    : ACTIVE (key-only)"
echo "  - fail2ban         : ACTIVE"
echo "  - Nginx            : ACTIVE"
echo "  - Project dir      : /home/infra/umrohos"
echo ""
echo "Next steps:"
echo "  1. Fase 3 - Buat docker-compose.prod.yml"
echo "  2. Fase 4 - Setup .env.prod secrets"
echo "  3. Fase 5 - GitHub Actions CD pipeline"
echo ""
echo "Akses monitoring via SSH tunnel:"
echo "  ssh -L 3000:localhost:3000 -L 9090:localhost:9090 \\"
echo "      -i ${SSH_KEY} infra@${SSH_HOST}"
echo ""
