-- Rollback UUID extension creation
-- Note: This is usually not safe to run in production if other objects depend on this extension
drop extension if exists "uuid-ossp";
