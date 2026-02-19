#!/usr/bin/env bash
# =============================================================================
# TMS — Keepalive ping (runs every 5 min via cron)
# Keeps the Go backend warm so the first real request never hits a cold start.
# =============================================================================
curl -sf -o /dev/null --max-time 10 http://127.0.0.1:8080/health || \
  echo "$(date): backend health check failed"
