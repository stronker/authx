-- KEYSPACE
create KEYSPACE IF NOT EXISTS authx WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

-- TABLES
create table authx.credentials (username text, password blob, role_id text, organization_id text, PRIMARY KEY (username));
create table authx.roles (organization_id text, role_id text, name text, internal boolean, primitives list<text>, PRIMARY KEY (organization_id, role_id));
create table authx.tokens (username text, token_id text, refresh_token blob, expiration_date bigint, PRIMARY KEY (username, token_id));