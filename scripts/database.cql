-- KEYSPACE
create KEYSPACE IF NOT EXISTS authx WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

-- TABLES
create table authx.credentials (username text, password blob, role_id text, organization_id text, PRIMARY KEY (username));
create table authx.roles (organization_id text, role_id text, name text, internal boolean, primitives list<text>, PRIMARY KEY (organization_id, role_id));
create table authx.tokens (username text, token_id text, refresh_token blob, expiration_date bigint, PRIMARY KEY (username, token_id));

create table IF NOT EXISTS authx.deviceTokens (device_id text, token_id text, refresh_token text, expiration_date bigint, organization_id text, device_group_id text, PRIMARY KEY (device_id, token_id));
create table IF NOT EXISTS authx.deviceCredentials (organization_id text, device_group_id text, device_id text, device_api_key text, enabled boolean, PRIMARY KEY ((organization_id, device_group_id), device_id));
create table IF NOT EXISTS authx.deviceGroupCredentials (organization_id text, device_group_id text,  device_group_api_key text, enabled boolean, default_device_connectivity boolean, secret text, PRIMARY KEY (organization_id, device_group_id));
create INDEX IF NOT EXISTS device_group_api ON authx.devicegroupcredentials ( device_group_api_key);
create INDEX IF NOT EXISTS device_api ON authx.devicecredentials ( device_api_key);
create INDEX IF NOT EXISTS device_group_secret ON authx.devicegroupcredentials ( secret);
create INDEX IF NOT EXISTS device_refresh_token ON authx.devicetokens ( refresh_token);
