#!/usr/bin/env bash
# =============================================================================
# TMS — Initial Deploy Script (Ubuntu 24.04, run as rhawez)
#
# Prerequisites already done on this VPS:
#   ✅ Docker + Compose plugin (rhawez in docker group)
#   ✅ Nginx with generous timeouts
#   ✅ Certbot with auto-renewal
#   ✅ UFW, Fail2ban, hardened SSH on port 2202
#
# Usage:
#   1. Fill in REPO_URL and EMAIL below
#   2. scp or git-clone this repo onto the VPS, then:
#      bash deploy/setup.sh
# =============================================================================
set -euo pipefail

DOMAIN="training.inova.krd"
REPO_DIR="$HOME/tms"
EMAIL="rawand.hawez@inov.krd"
REPO_URL="https://github.com/Inova-Deep/tms"

echo ""
echo "=== [1/4] Clone repository ==="
git clone "$REPO_URL" "$REPO_DIR"
cd "$REPO_DIR"

echo ""
echo "=== [2/4] Configure nginx site ==="
sudo cp deploy/nginx.conf /etc/nginx/sites-available/$DOMAIN
sudo ln -sf /etc/nginx/sites-available/$DOMAIN /etc/nginx/sites-enabled/$DOMAIN
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t && sudo systemctl reload nginx

echo ""
echo "=== [3/4] Obtain SSL certificate ==="
# Certbot will auto-edit nginx.conf to add the ssl_certificate lines
sudo certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos -m "$EMAIL"
sudo systemctl reload nginx

echo ""
echo "=== [4/4] Build and start containers ==="
cd "$REPO_DIR"
docker compose -f docker-compose.prod.yml up -d --build

echo ""
echo "=== Install keepalive cron (prevents cold-start timeouts) ==="
sudo cp deploy/keepalive.sh /usr/local/bin/tms-keepalive
sudo chmod +x /usr/local/bin/tms-keepalive
# Add cron job for current user (rhawez) - pings backend every 5 minutes
(crontab -l 2>/dev/null; echo "*/5 * * * * /usr/local/bin/tms-keepalive >> /var/log/tms-keepalive.log 2>&1") | crontab -

echo ""
echo "========================================"
echo " Setup complete!"
echo " Visit: https://$DOMAIN"
echo "========================================"
