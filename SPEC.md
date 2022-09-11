# PostgreSQL YAML Role Manager - CLI Specification

## Overview
A Kubernetes-inspired CLI tool for managing PostgreSQL roles and users using YAML manifests. This tool will allow users to define PostgreSQL roles and users in YAML files and apply them to a PostgreSQL database, similar to how `kubectl` works with Kubernetes resources.

## Features
1. Apply YAML manifests to create/update PostgreSQL roles and users
2. Delete roles and users based on YAML manifests
3. List existing roles and users
4. Validate YAML manifests before applying
5. Support for role grants and permissions

## Commands

### 1. Apply
Apply a YAML manifest to create or update PostgreSQL roles and users.

```bash
pgctl apply -f <file.yaml>
pgctl apply -f <directory/>
```

### 2. Delete
Delete PostgreSQL roles and users defined in a YAML manifest.

```bash
pgctl delete -f <file.yaml>
pgctl delete -f <directory/>
```

### 3. List
List existing roles and users in the database.

```bash
pgctl get roles
pgctl get users
```

### 4. Describe
Show detailed information about a specific role or user.

```bash
pgctl describe role <role-name>
pgctl describe user <user-name>
```

### 5. Validate
Validate a YAML manifest without applying it.

```bash
pgctl validate -f <file.yaml>
```

## YAML Manifest Structure

### Role Definition
```yaml
apiVersion: pgctl/v1
kind: Role
metadata:
  name: myrole
spec:
  login: false
  inherit: true
  createrole: false
  createdb: false
  replication: false
  bypassrls: false
  connectionLimit: -1
  validUntil: "2024-12-31T23:59:59Z"
  comment: "A sample role"
```

### User Definition
```yaml
apiVersion: pgctl/v1
kind: User
metadata:
  name: myuser
spec:
  login: true
  inherit: true
  createrole: false
  createdb: false
  replication: false
  bypassrls: false
  connectionLimit: 10
  password: "mypassword"
  validUntil: "2024-12-31T23:59:59Z"
  comment: "A sample user"
  inRoles:
    - myrole
```

### Grant Definition
```yaml
apiVersion: pgctl/v1
kind: Grant
metadata:
  name: mygrant
spec:
  role: myrole
  database: mydatabase
  schema: myschema
  table: mytable
  privileges:
    - SELECT
    - INSERT
```

## Configuration

The tool will look for configuration in the following order:
1. `--config` flag
2. `PGCTL_CONFIG` environment variable
3. `$HOME/.pgctl/config.yaml`
4. Default configuration

### Configuration File Structure
```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: postgres
  sslmode: disable
```

### Environment Variables
- `PGCTL_CONFIG`: Path to configuration file
- `PGCTL_DATABASE_HOST`: Database host
- `PGCTL_DATABASE_PORT`: Database port
- `PGCTL_DATABASE_USER`: Database user
- `PGCTL_DATABASE_PASSWORD`: Database password
- `PGCTL_DATABASE_NAME`: Database name
- `PGCTL_DATABASE_SSLMODE`: SSL mode

## Installation

### Binary Installation
Download the appropriate binary for your platform from the releases page.

### From Source
```bash
go install github.com/yourusername/postgre-yaml-role@latest
```

## Examples

### Create a Role
```yaml
# role.yaml
apiVersion: pgctl/v1
kind: Role
metadata:
  name: readonly
spec:
  login: false
  inherit: true
  createrole: false
  createdb: false
  replication: false
  bypassrls: false
  connectionLimit: -1
  comment: "Read-only role"
```

Apply the role:
```bash
pgctl apply -f role.yaml
```

### Create a User
```yaml
# user.yaml
apiVersion: pgctl/v1
kind: User
metadata:
  name: appuser
spec:
  login: true
  inherit: true
  createrole: false
  createdb: false
  replication: false
  bypassrls: false
  connectionLimit: 5
  password: "securepassword"
  validUntil: "2025-12-31T23:59:59Z"
  comment: "Application user"
  inRoles:
    - readonly
```

Apply the user:
```bash
pgctl apply -f user.yaml
```

### Grant Permissions
```yaml
# grant.yaml
apiVersion: pgctl/v1
kind: Grant
metadata:
  name: readonly-grant
spec:
  role: readonly
  database: myapp
  schema: public
  table: "*"
  privileges:
    - SELECT
```

Apply the grant:
```bash
pgctl apply -f grant.yaml
```

## Implementation Plan

1. Create Go project structure
2. Implement YAML parsing and validation
3. Implement PostgreSQL connection and operations
4. Implement CLI commands
5. Add configuration management
6. Add validation and error handling
7. Write tests
8. Create documentation

## Dependencies

- Go 1.19+
- PostgreSQL driver (github.com/lib/pq)
- YAML parser (gopkg.in/yaml.v3)
- CLI library (github.com/spf13/cobra)
- Configuration library (github.com/spf13/viper)

## Project Structure
```
postgre-yaml-role/
├── cmd/
│   ├── root.go
│   ├── apply.go
│   ├── delete.go
│   ├── get.go
│   ├── describe.go
│   └── validate.go
├── pkg/
│   ├── config/
│   ├── postgres/
│   ├── yaml/
│   └── types/
├── examples/
├── SPEC.md
├── README.md
├── go.mod
└── go.sum
```