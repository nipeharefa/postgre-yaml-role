# PostgreSQL YAML Role Manager (pgctl)

A Kubernetes-inspired CLI tool for managing PostgreSQL roles and users using YAML manifests.

## Overview

`pgctl` is a command-line tool that allows you to manage PostgreSQL roles and users using YAML manifests, similar to how `kubectl` works with Kubernetes resources. This approach provides a declarative way to manage your PostgreSQL permissions and access control.

## Features

- Apply YAML manifests to create/update PostgreSQL roles and users
- Delete roles and users based on YAML manifests
- List existing roles and users
- Validate YAML manifests before applying
- Support for role grants and permissions

## Installation

### From Source

```bash
go install github.com/nipeharefa/postgre-yaml-role@latest
```

### Binary Installation

Download the appropriate binary for your platform from the releases page.

## Quick Start

1. Create a configuration file at `~/.pgctl/config.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: postgres
  sslmode: disable
```

2. Create a role manifest (`role.yaml`):

```yaml
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

3. Apply the manifest:

```bash
pgctl apply -f role.yaml
```

## Commands

### Apply
Apply a YAML manifest to create or update PostgreSQL roles and users.

```bash
pgctl apply -f <file.yaml>
pgctl apply -f <directory/>
```

### Delete
Delete PostgreSQL roles and users defined in a YAML manifest.

```bash
pgctl delete -f <file.yaml>
pgctl delete -f <directory/>
```

### List
List existing roles and users in the database.

```bash
pgctl get roles
pgctl get users
```

### Validate
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

### Environment Variables
- `PGCTL_CONFIG`: Path to configuration file
- `PGCTL_DATABASE_HOST`: Database host
- `PGCTL_DATABASE_PORT`: Database port
- `PGCTL_DATABASE_USER`: Database user
- `PGCTL_DATABASE_PASSWORD`: Database password
- `PGCTL_DATABASE_NAME`: Database name
- `PGCTL_DATABASE_SSLMODE`: SSL mode

## Examples

See the [examples](examples/) directory for sample YAML manifests.

## Building from Source

```bash
git clone https://github.com/nipeharefa/postgre-yaml-role.git
cd postgre-yaml-role
go build -o pgctl
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.