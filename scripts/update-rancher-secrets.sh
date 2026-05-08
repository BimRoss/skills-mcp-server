#!/usr/bin/env bash
set -euo pipefail

# Push skills-mcp-server runtime Secret to the admin Kubernetes cluster.
#
# The MCP tool create_google_doc runs in THIS service; Google OAuth env vars must live in
# Secret skills-mcp-server-runtime (namespace skills-mcp-server), not only on agents-mcp-server.
#
# Sources repo-root .env.prod if present, else .env, unless ENV_FILE is set.
# Do not point at .env.dev for production cluster sync.
#
# Keys written (see internal/config/config.go + internal/googledocs/config.go):
#   GEMINI_API_KEY — required for read_web / router-adjacent tooling (from ENV_FILE or existing cluster secret)
#   JOANNE_GOOGLE_CLIENT_ID, JOANNE_GOOGLE_CLIENT_SECRET, JOANNE_GOOGLE_REFRESH_TOKEN — preferred for create_google_doc
#   Optional fallbacks (only if set in ENV_FILE or already in cluster): GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET /
#     GOOGLE_REFRESH_TOKEN, GOOGLE_OAUTH_CLIENT_ID / GOOGLE_OAUTH_CLIENT_SECRET / GOOGLE_OAUTH_REFRESH_TOKEN
#
# By default copies dockerhub-pull into the namespace (same pattern as makeacompany-ai update-rancher-secrets).
# Set SYNC_PULL_SECRET=false to skip.
#
# Usage:
#   ./scripts/update-rancher-secrets.sh
#   ENV_FILE=/path/to/agents-mcp-server/.env.prod ./scripts/update-rancher-secrets.sh
#
# Kube: if KUBECONFIG is unset, uses ~/.kube/config/admin.yaml or grant-admin.yaml when present.
# Honor KUBECONFIG_HOST_PATH / KUBECONFIG_CONTEXT from ENV_FILE like makeacompany-ai.

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

if [[ -z "${KUBECONFIG:-}" ]]; then
  for _k in "${HOME}/.kube/config/admin.yaml" "${HOME}/.kube/config/grant-admin.yaml"; do
    if [[ -f "${_k}" ]]; then
      export KUBECONFIG="${_k}"
      break
    fi
  done
fi

if [[ -z "${ENV_FILE:-}" ]]; then
  if [[ -f "${ROOT}/.env.prod" ]]; then
    ENV_FILE="${ROOT}/.env.prod"
  else
    ENV_FILE="${ROOT}/.env"
  fi
fi

NAMESPACE="${NAMESPACE:-skills-mcp-server}"
SECRET_NAME="${SECRET_NAME:-skills-mcp-server-runtime}"
SYNC_PULL_SECRET="${SYNC_PULL_SECRET:-true}"
PULL_SECRET_NAME="${PULL_SECRET_NAME:-dockerhub-pull}"
PULL_SECRET_SOURCE_NAMESPACE="${PULL_SECRET_SOURCE_NAMESPACE:-bimross-web}"
PULL_SECRET_FALLBACK_NAMESPACE="${PULL_SECRET_FALLBACK_NAMESPACE:-employee-factory}"

kubectl_app() {
  local args=()
  if [[ -n "${KUBE_CONTEXT:-}" ]]; then
    args+=(--context "${KUBE_CONTEXT}")
  fi
  kubectl "${args[@]}" "$@"
}

sync_pull_secret() {
  local source_ns="${PULL_SECRET_SOURCE_NAMESPACE}"

  kubectl_app get namespace "${NAMESPACE}" >/dev/null 2>&1 || kubectl_app create namespace "${NAMESPACE}"

  if ! kubectl_app get secret "${PULL_SECRET_NAME}" -n "${source_ns}" >/dev/null 2>&1; then
    echo "Pull secret '${PULL_SECRET_NAME}' not found in '${source_ns}', trying '${PULL_SECRET_FALLBACK_NAMESPACE}'..."
    source_ns="${PULL_SECRET_FALLBACK_NAMESPACE}"
    kubectl_app get secret "${PULL_SECRET_NAME}" -n "${source_ns}" >/dev/null 2>&1 || {
      echo "Unable to find '${PULL_SECRET_NAME}' in '${PULL_SECRET_SOURCE_NAMESPACE}' or '${PULL_SECRET_FALLBACK_NAMESPACE}'." >&2
      exit 1
    }
  fi

  kubectl_app get secret "${PULL_SECRET_NAME}" -n "${source_ns}" -o json \
    | python3 -c 'import json,sys; src=json.load(sys.stdin); out={"apiVersion":"v1","kind":"Secret","metadata":{"name":src["metadata"]["name"],"namespace":"'"${NAMESPACE}"'"},"type":src.get("type"),"data":src.get("data",{})}; print(json.dumps(out))' \
    | kubectl_app apply -f -

  echo "Synced '${PULL_SECRET_NAME}' into namespace '${NAMESPACE}' from '${source_ns}'."
}

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "missing ${ENV_FILE}" >&2
  exit 1
fi

set -a
# shellcheck source=/dev/null
source "${ENV_FILE}"
set +a

if [[ -n "${KUBECONFIG_HOST_PATH:-}" && -f "${KUBECONFIG_HOST_PATH}" ]]; then
  export KUBECONFIG="${KUBECONFIG_HOST_PATH}"
fi
if [[ -n "${KUBECONFIG_CONTEXT:-}" && -z "${KUBE_CONTEXT:-}" ]]; then
  export KUBE_CONTEXT="${KUBECONFIG_CONTEXT}"
fi

read_existing_secret_key() {
  local key="$1"
  kubectl_app get secret "${SECRET_NAME}" -n "${NAMESPACE}" -o "jsonpath={.data.${key}}" 2>/dev/null \
    | python3 -c 'import sys,base64; raw=sys.stdin.read().strip(); print(base64.b64decode(raw).decode() if raw else "")' 2>/dev/null || true
}

if [[ "${SYNC_PULL_SECRET}" == "true" ]]; then
  sync_pull_secret
fi

GEMINI_EFFECTIVE="${GEMINI_API_KEY:-}"
if [[ -z "${GEMINI_EFFECTIVE}" ]]; then
  GEMINI_EFFECTIVE="$(read_existing_secret_key GEMINI_API_KEY)"
fi
if [[ -z "${GEMINI_EFFECTIVE}" ]]; then
  echo "need GEMINI_API_KEY in ${ENV_FILE} or already in cluster secret ${SECRET_NAME}" >&2
  exit 1
fi

merge_key() {
  local env_val="$1"
  local key="$2"
  local from_cluster=""
  if [[ -z "${env_val}" ]]; then
    from_cluster="$(read_existing_secret_key "${key}")"
    printf '%s' "${from_cluster}"
  else
    printf '%s' "${env_val}"
  fi
}

# Google: env wins; then existing cluster value for the same key name.
JOANNE_GOOGLE_CLIENT_ID_EFFECTIVE="$(merge_key "${JOANNE_GOOGLE_CLIENT_ID:-}" JOANNE_GOOGLE_CLIENT_ID)"
JOANNE_GOOGLE_CLIENT_SECRET_EFFECTIVE="$(merge_key "${JOANNE_GOOGLE_CLIENT_SECRET:-}" JOANNE_GOOGLE_CLIENT_SECRET)"
JOANNE_GOOGLE_REFRESH_TOKEN_EFFECTIVE="$(merge_key "${JOANNE_GOOGLE_REFRESH_TOKEN:-}" JOANNE_GOOGLE_REFRESH_TOKEN)"
GOOGLE_CLIENT_ID_EFFECTIVE="$(merge_key "${GOOGLE_CLIENT_ID:-}" GOOGLE_CLIENT_ID)"
GOOGLE_CLIENT_SECRET_EFFECTIVE="$(merge_key "${GOOGLE_CLIENT_SECRET:-}" GOOGLE_CLIENT_SECRET)"
GOOGLE_REFRESH_TOKEN_EFFECTIVE="$(merge_key "${GOOGLE_REFRESH_TOKEN:-}" GOOGLE_REFRESH_TOKEN)"
GOOGLE_OAUTH_CLIENT_ID_EFFECTIVE="$(merge_key "${GOOGLE_OAUTH_CLIENT_ID:-}" GOOGLE_OAUTH_CLIENT_ID)"
GOOGLE_OAUTH_CLIENT_SECRET_EFFECTIVE="$(merge_key "${GOOGLE_OAUTH_CLIENT_SECRET:-}" GOOGLE_OAUTH_CLIENT_SECRET)"
GOOGLE_OAUTH_REFRESH_TOKEN_EFFECTIVE="$(merge_key "${GOOGLE_OAUTH_REFRESH_TOKEN:-}" GOOGLE_OAUTH_REFRESH_TOKEN)"

secret_args=(--namespace "${NAMESPACE}")
secret_args+=(--from-literal=GEMINI_API_KEY="${GEMINI_EFFECTIVE}")

append_if_nonempty() {
  local k="$1"
  local v="$2"
  if [[ -n "${v}" ]]; then
    secret_args+=(--from-literal="${k}=${v}")
  fi
}

append_if_nonempty JOANNE_GOOGLE_CLIENT_ID "${JOANNE_GOOGLE_CLIENT_ID_EFFECTIVE}"
append_if_nonempty JOANNE_GOOGLE_CLIENT_SECRET "${JOANNE_GOOGLE_CLIENT_SECRET_EFFECTIVE}"
append_if_nonempty JOANNE_GOOGLE_REFRESH_TOKEN "${JOANNE_GOOGLE_REFRESH_TOKEN_EFFECTIVE}"
append_if_nonempty GOOGLE_CLIENT_ID "${GOOGLE_CLIENT_ID_EFFECTIVE}"
append_if_nonempty GOOGLE_CLIENT_SECRET "${GOOGLE_CLIENT_SECRET_EFFECTIVE}"
append_if_nonempty GOOGLE_REFRESH_TOKEN "${GOOGLE_REFRESH_TOKEN_EFFECTIVE}"
append_if_nonempty GOOGLE_OAUTH_CLIENT_ID "${GOOGLE_OAUTH_CLIENT_ID_EFFECTIVE}"
append_if_nonempty GOOGLE_OAUTH_CLIENT_SECRET "${GOOGLE_OAUTH_CLIENT_SECRET_EFFECTIVE}"
append_if_nonempty GOOGLE_OAUTH_REFRESH_TOKEN "${GOOGLE_OAUTH_REFRESH_TOKEN_EFFECTIVE}"

# Match googledocs.LoadFromEnv precedence for diagnostics only
GOOGLE_DOC_ID="${JOANNE_GOOGLE_CLIENT_ID_EFFECTIVE:-${GOOGLE_CLIENT_ID_EFFECTIVE:-${GOOGLE_OAUTH_CLIENT_ID_EFFECTIVE:-}}}"
if [[ -z "${GOOGLE_DOC_ID}" ]]; then
  echo "WARNING: no Google OAuth client id in ${ENV_FILE} or cluster — create_google_doc will fail until JOANNE_GOOGLE_CLIENT_ID (or GOOGLE_* / GOOGLE_OAUTH_*) is set." >&2
fi

kubectl_app create secret generic "${SECRET_NAME}" \
  "${secret_args[@]}" \
  --dry-run=client -o yaml | kubectl_app apply -f -

echo "applied secret ${SECRET_NAME} in namespace ${NAMESPACE}"

ROLLOUT_AFTER_SECRET_SYNC="${ROLLOUT_AFTER_SECRET_SYNC:-true}"
if [[ "${ROLLOUT_AFTER_SECRET_SYNC}" == "true" ]]; then
  if kubectl_app get deployment skills-mcp-server -n "${NAMESPACE}" >/dev/null 2>&1; then
    kubectl_app rollout restart deployment/skills-mcp-server -n "${NAMESPACE}"
    echo "rollout restart: skills-mcp-server"
  fi
fi
