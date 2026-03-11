#!/bin/bash

# ============================================================================
# Digna Deployment Wrapper
# ============================================================================
# Wrapper script for easy deployment from project root
# ============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_SCRIPT="${SCRIPT_DIR}/scripts/deploy/vps_deploy.sh"

if [[ ! -f "${DEPLOY_SCRIPT}" ]]; then
    echo "Error: Deployment script not found at ${DEPLOY_SCRIPT}"
    echo "Please run from project root directory"
    exit 1
fi

# Pass all arguments to the actual deployment script
"${DEPLOY_SCRIPT}" "$@"