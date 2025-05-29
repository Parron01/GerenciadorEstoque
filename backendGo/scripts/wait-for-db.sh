#!/bin/bash
set -e

host="$POSTGRES_HOST"
port="$POSTGRES_PORT"
user="$POSTGRES_USER"
password="$POSTGRES_PASSWORD"
dbname="$POSTGRES_DB"

echo "Waiting for PostgreSQL ($host:$port)..."

# Maximum number of attempts
max_tries=30
# Delay between tries (in seconds)
delay=2

# Counter for tries
tries=0

# Keep trying until max_tries is reached
while [ $tries -lt $max_tries ]; do
  tries=$((tries+1))
  
  # Try to connect to PostgreSQL
  export PGPASSWORD="$password"
  if psql -h "$host" -p "$port" -U "$user" -d "$dbname" -c "SELECT 1" >/dev/null 2>&1; then
    echo "PostgreSQL is ready! Starting application..."
    break
  fi
  
  # If not successful and max_tries is not reached yet, wait and try again
  if [ $tries -lt $max_tries ]; then
    echo "Attempt $tries failed. Trying again in $delay seconds..."
    sleep $delay
  else
    echo "Could not connect to PostgreSQL after $max_tries attempts. Giving up."
    exit 1
  fi
done

# Execute the command passed to docker run
exec "$@"
