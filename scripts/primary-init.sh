#!/bin/bash
set -e

echo "Configuring PostgreSQL Primary for replication..."

# Create replication user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = '$POSTGRES_REPLICATION_USER') THEN
            CREATE USER $POSTGRES_REPLICATION_USER WITH REPLICATION ENCRYPTED PASSWORD '$POSTGRES_REPLICATION_PASSWORD';
        END IF;
    END
    \$\$;
EOSQL

# Configure pg_hba.conf for replication
# Allow replication connections from any host in Docker network
cat >> "$PGDATA/pg_hba.conf" <<EOF

# Replication connections
host    replication     $POSTGRES_REPLICATION_USER      0.0.0.0/0               md5
host    replication     $POSTGRES_REPLICATION_USER      ::/0                    md5
EOF

echo "Primary PostgreSQL configured for replication successfully"
