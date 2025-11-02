# Database Migration Strategy

This document outlines the database migration strategy for the Ultimate Frisbee API project.

## Overview

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management. It provides a robust, production-ready migration system with support for multiple database engines.

## Migration File Structure

```
infra/database/migrations/
├── 000001_create_uuid_extension.up.sql
├── 000001_create_uuid_extension.down.sql
├── 000002_teams.up.sql
└── 000002_teams.down.sql
```

### Naming Convention

- **Sequential numbering**: `000001_`, `000002_`, etc.
- **Descriptive names**: `create_uuid_extension`, `teams`, `add_user_roles`
- **File pairs**: Each migration should have both `.up.sql` and `.down.sql` files

## Available Commands

### Using Make Commands (Recommended)

```bash
# Install migration tool (required first time)
brew install golang-migrate  # macOS
# or
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run all pending migrations
make db/migration/up

# Rollback all migrations
make db/migration/down

# Create a new migration
make db/migration/create MIGRATION_NAME=add_user_table

# Seed the database with test data
make db/seed

# Start/stop database service
make services/database/up
make services/database/down
```

### Using Application Binary

```bash
# Run migrations programmatically through the application
go run *.go -e migration
```

### Direct Migration Commands

```bash
# Run migrations
migrate -path ./infra/database/migrations -database "postgres://user:pass@localhost:42013/db?sslmode=disable" up

# Rollback one migration
migrate -path ./infra/database/migrations -database "postgres://user:pass@localhost:42013/db?sslmode=disable" down 1

# Check migration status
migrate -path ./infra/database/migrations -database "postgres://user:pass@localhost:42013/db?sslmode=disable" version
```

## Best Practices

### 1. Migration Safety

- **Always create both up and down migrations**
- **Test migrations on a copy of production data**
- **Use transactions when possible**
- **Make migrations idempotent** using `IF NOT EXISTS` clauses

```sql
-- Good: Idempotent migration
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL
);

-- Avoid: Non-idempotent migration
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL
);
```

### 2. Schema Changes

- **Additive changes**: Prefer adding columns/tables over modifying existing ones
- **Backward compatibility**: Ensure new migrations don't break existing code
- **Data migration**: Separate schema changes from data changes when possible

### 3. Production Considerations

- **Backup first**: Always backup production database before running migrations
- **Downtime planning**: Some migrations may require maintenance windows
- **Rollback plan**: Always have a tested rollback strategy
- **Monitoring**: Monitor migration performance on large tables

### 4. Column Changes

```sql
-- Safe: Adding a new nullable column
ALTER TABLE teams ADD COLUMN website VARCHAR(255);

-- Potentially unsafe: Changing column type
-- Better to: Add new column, migrate data, drop old column in separate migrations
ALTER TABLE teams ADD COLUMN new_description TEXT;
-- (migrate data in application code or separate migration)
-- ALTER TABLE teams DROP COLUMN description;
-- ALTER TABLE teams RENAME COLUMN new_description TO description;
```

### 5. Index Management

```sql
-- Create indexes concurrently to avoid locking (PostgreSQL)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_teams_origin_country 
ON teams(origin_country);
```

## Migration Workflow

### Development

1. **Create migration**: `make db/migration/create MIGRATION_NAME=your_change`
2. **Edit up/down files** with your schema changes
3. **Test locally**: `make db/migration/up`
4. **Verify changes** work with your application
5. **Test rollback**: `make db/migration/down` then `make db/migration/up`

### Production Deployment

1. **Review migrations** in staging environment
2. **Backup production database**
3. **Deploy application** with migration code
4. **Run migrations**: `make db/migration/up` or programmatically
5. **Verify application** works correctly
6. **Monitor** for any issues

## Environment Configuration

### Local Development
```yaml
# config/local.yaml
database:
  connectionString: postgres://ultimate_frisbee_manager_user:some_password@localhost:42013/ultimate_frisbee_manager?sslmode=disable
```

### Docker Compose
```yaml
# docker-compose.yaml
postgres:
  image: postgres:12.3
  ports:
    - "42013:5432"
  environment:
    POSTGRES_USER: ultimate_frisbee_manager_user
    POSTGRES_PASSWORD: some_password
    POSTGRES_DB: ultimate_frisbee_manager
```

## Troubleshooting

### Common Issues

1. **Permission denied**: Ensure database user has CREATE/ALTER privileges
2. **Migration already applied**: Check current version with `migrate version`
3. **Dirty state**: Force version if migration failed midway (use with caution)

```bash
# Check current migration status
migrate -path ./infra/database/migrations -database $DATABASE_CONNECTION_STRING version

# Force version (dangerous - only if you're sure)
migrate -path ./infra/database/migrations -database $DATABASE_CONNECTION_STRING force VERSION
```

### Recovery

If migrations fail in production:

1. **Don't panic** - assess the situation
2. **Check logs** for specific error messages
3. **Verify database state** manually
4. **Consider rollback** if safe to do so
5. **Fix forward** if rollback is not possible

## Seeding Data

Seed files are located in `infra/database/seeds/` and should be run after migrations:

```bash
make db/seed
```

Seed files are numbered and should be idempotent using `INSERT ... ON CONFLICT` or similar patterns.

## Integration with Application

The migration system is integrated into the main application binary:

```go
// Run migrations programmatically
go run *.go -e migration
```

This ensures migrations use the same configuration as the application and can be easily automated in deployment pipelines.
