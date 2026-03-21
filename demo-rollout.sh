#!/bin/bash
# Demo: Rollout operator - gradual migration over 30 seconds

NUM_USERS=20
DURATION=30
FLAG_FILE="$(dirname "$0")/rollout-demo-flags.json"
OFREP_PORT=8016

# Kill any existing flagd
pkill -9 -f "flagd start" 2>/dev/null || true
sleep 1

# Calculate timestamps - start 5 seconds in the future to allow setup
START_TIME=$(($(date +%s) + 5))
END_TIME=$((START_TIME + DURATION))

# Create flag config - NOTE THE CONFIG DOES NOT CHANGE DURING THE DEMO - the rule itself is time-aware
cat > "$FLAG_FILE" << EOF
{
  "flags": {
    "demo": {
      "state": "ENABLED",
      "variants": { "old": "OLD", "new": "NEW" },
      "defaultVariant": "old",
      "targeting": { "rollout": [$START_TIME, $END_TIME, "new"] }
    }
  }
}
EOF

# Generate random user keys
declare -a USERS
for i in $(seq 1 $NUM_USERS); do
  USERS[$i]="user_$(head -c4 /dev/urandom | xxd -p)"
done

echo "Starting flagd (rollout: $START_TIME -> $END_TIME)..."
./bin/flagd start -f "file:$FLAG_FILE" &>/dev/null &
FLAGD_PID=$!
trap "kill $FLAGD_PID 2>/dev/null" EXIT

# Wait until rollout starts
now=$(date +%s)
wait_time=$((START_TIME - now))
if [[ $wait_time -gt 0 ]]; then
  echo "Waiting ${wait_time}s for rollout to begin..."
  sleep $wait_time
fi

# Evaluate flag for a user
eval_flag() {
  curl -s "http://localhost:$OFREP_PORT/ofrep/v1/evaluate/flags/demo" \
    -H "Content-Type: application/json" \
    -d "{\"context\":{\"targetingKey\":\"$1\"}}" | grep -o '"value":"[^"]*"' | cut -d'"' -f4
}

# Header
printf "\n%-6s" "TIME"
for i in $(seq 1 $NUM_USERS); do printf "  U%-2d" "$i"; done
printf "  %%NEW\n"
printf "%s\n" "$(printf '=%.0s' {1..60})"

# Poll every 2 seconds
for t in $(seq 0 2 $DURATION); do
  [[ $t -gt 0 ]] && sleep 2
  
  printf "%-6s" "${t}s"
  new_count=0
  
  for i in $(seq 1 $NUM_USERS); do
    result=$(eval_flag "${USERS[$i]}")
    if [[ "$result" == "NEW" ]]; then
      printf "  \033[32m%-3s\033[0m" "NEW"
      new_count=$((new_count + 1))
    else
      printf "  \033[33m%-3s\033[0m" "OLD"
    fi
  done
  
  pct=$((new_count * 100 / NUM_USERS))
  printf "  %3d%%\n" "$pct"
done

echo -e "\nRollout complete!"
