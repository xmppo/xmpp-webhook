#!/usr/bin/env bash

curl -H "Content-Type: application/json" -d '[{"labels":{"alertname":"XMPPTest"}}]' localhost:9093/api/v1/alerts
