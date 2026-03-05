#!/bin/bash
# Gas Town preToolUse guard — filters tool invocations and enforces PR workflow policy.
# Fail-closed: if gt tap is unavailable, deny all bash commands rather than allowing
# unguarded execution.
INPUT=$(cat)
TOOL_NAME=$(echo "$INPUT" | jq -r '.toolName')
[ "$TOOL_NAME" = "bash" ] || exit 0

COMMAND=$(echo "$INPUT" | jq -r '.toolArgs' | jq -r '.command // empty')
[ -n "$COMMAND" ] || exit 0

# Fail-closed: deny all bash commands if gt is not available.
# Without gt tap, guard policies cannot be evaluated — block everything.
if ! command -v gt >/dev/null 2>&1; then
  jq -nc '{"permissionDecision":"deny","permissionDecisionReason":"gt not found in PATH — guard cannot evaluate policies (fail-closed)"}'
  exit 0
fi

if echo "$COMMAND" | grep -qE '(^|[;&|]\s*|&&\s*|\|\|\s*)(\s*)(gh pr create|git checkout -b|git switch -c)'; then
  RESULT=$(gt tap guard pr-workflow 2>&1)
  EXIT_CODE=$?
  if [ $EXIT_CODE -ne 0 ]; then
    jq -nc --arg reason "$RESULT" \
      '{"permissionDecision":"deny","permissionDecisionReason":$reason}'
  fi
fi
