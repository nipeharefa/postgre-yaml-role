package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/nipeharefa/postgre-yaml-role/pkg/types"
)

// DB represents a PostgreSQL database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new PostgreSQL database connection
func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// CreateRole creates a new PostgreSQL role
func (db *DB) CreateRole(role types.Role) error {
	// Build the CREATE ROLE statement
	var stmt strings.Builder
	stmt.WriteString("CREATE ROLE ")
	stmt.WriteString(pq.QuoteIdentifier(role.Metadata.Name))

	options := []string{}

	if role.Spec.Login {
		options = append(options, "LOGIN")
	} else {
		options = append(options, "NOLOGIN")
	}

	if role.Spec.Inherit {
		options = append(options, "INHERIT")
	} else {
		options = append(options, "NOINHERIT")
	}

	if role.Spec.CreateRole {
		options = append(options, "CREATEROLE")
	} else {
		options = append(options, "NOCREATEROLE")
	}

	if role.Spec.CreateDB {
		options = append(options, "CREATEDB")
	} else {
		options = append(options, "NOCREATEDB")
	}

	if role.Spec.Replication {
		options = append(options, "REPLICATION")
	} else {
		options = append(options, "NOREPLICATION")
	}

	if role.Spec.BypassRLS {
		options = append(options, "BYPASSRLS")
	} else {
		options = append(options, "NOBYPASSRLS")
	}

	if role.Spec.ConnectionLimit >= 0 {
		options = append(options, fmt.Sprintf("CONNECTION LIMIT %d", role.Spec.ConnectionLimit))
	}

	if role.Spec.ValidUntil != "" {
		options = append(options, fmt.Sprintf("VALID UNTIL '%s'", role.Spec.ValidUntil))
	}

	if len(options) > 0 {
		stmt.WriteString(" ")
		stmt.WriteString(strings.Join(options, " "))
	}

	if role.Spec.Comment != "" {
		stmt.WriteString(fmt.Sprintf("; COMMENT ON ROLE %s IS '%s'", 
			pq.QuoteIdentifier(role.Metadata.Name), role.Spec.Comment))
	}

	_, err := db.Exec(stmt.String())
	return err
}

// UpdateRole updates an existing PostgreSQL role
func (db *DB) UpdateRole(role types.Role) error {
	// For simplicity, we'll drop and recreate the role
	// In a production implementation, you might want to alter the role instead
	if err := db.DeleteRole(role.Metadata.Name); err != nil {
		return err
	}
	return db.CreateRole(role)
}

// DeleteRole deletes a PostgreSQL role
func (db *DB) DeleteRole(name string) error {
	_, err := db.Exec("DROP ROLE IF EXISTS " + pq.QuoteIdentifier(name))
	return err
}

// CreateUser creates a new PostgreSQL user
func (db *DB) CreateUser(user types.User) error {
	// Build the CREATE USER statement
	var stmt strings.Builder
	stmt.WriteString("CREATE USER ")
	stmt.WriteString(pq.QuoteIdentifier(user.Metadata.Name))

	options := []string{}

	if user.Spec.Login {
		options = append(options, "LOGIN")
	}

	if user.Spec.Inherit {
		options = append(options, "INHERIT")
	} else {
		options = append(options, "NOINHERIT")
	}

	if user.Spec.CreateRole {
		options = append(options, "CREATEROLE")
	} else {
		options = append(options, "NOCREATEROLE")
	}

	if user.Spec.CreateDB {
		options = append(options, "CREATEDB")
	} else {
		options = append(options, "NOCREATEDB")
	}

	if user.Spec.Replication {
		options = append(options, "REPLICATION")
	} else {
		options = append(options, "NOREPLICATION")
	}

	if user.Spec.BypassRLS {
		options = append(options, "BYPASSRLS")
	} else {
		options = append(options, "NOBYPASSRLS")
	}

	if user.Spec.ConnectionLimit >= 0 {
		options = append(options, fmt.Sprintf("CONNECTION LIMIT %d", user.Spec.ConnectionLimit))
	}

	if user.Spec.Password != "" {
		options = append(options, fmt.Sprintf("PASSWORD '%s'", user.Spec.Password))
	}

	if user.Spec.ValidUntil != "" {
		options = append(options, fmt.Sprintf("VALID UNTIL '%s'", user.Spec.ValidUntil))
	}

	if len(options) > 0 {
		stmt.WriteString(" ")
		stmt.WriteString(strings.Join(options, " "))
	}

	// Execute the CREATE USER statement
	_, err := db.Exec(stmt.String())
	if err != nil {
		return err
	}

	// Add comment if specified
	if user.Spec.Comment != "" {
		_, err = db.Exec(fmt.Sprintf("COMMENT ON ROLE %s IS '%s'", 
			pq.QuoteIdentifier(user.Metadata.Name), user.Spec.Comment))
		if err != nil {
			return err
		}
	}

	// Add user to roles if specified
	for _, roleName := range user.Spec.InRoles {
		_, err = db.Exec(fmt.Sprintf("GRANT %s TO %s", 
			pq.QuoteIdentifier(roleName), pq.QuoteIdentifier(user.Metadata.Name)))
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateUser updates an existing PostgreSQL user
func (db *DB) UpdateUser(user types.User) error {
	// For simplicity, we'll drop and recreate the user
	// In a production implementation, you might want to alter the user instead
	if err := db.DeleteUser(user.Metadata.Name); err != nil {
		return err
	}
	return db.CreateUser(user)
}

// DeleteUser deletes a PostgreSQL user
func (db *DB) DeleteUser(name string) error {
	_, err := db.Exec("DROP USER IF EXISTS " + pq.QuoteIdentifier(name))
	return err
}

// CreateGrant grants privileges to a role
func (db *DB) CreateGrant(grant types.Grant) error {
	privileges := strings.Join(grant.Spec.Privileges, ", ")
	
	var stmt string
	if grant.Spec.Table != "" {
		stmt = fmt.Sprintf("GRANT %s ON TABLE %s.%s.%s TO %s", 
			privileges, 
			pq.QuoteIdentifier(grant.Spec.Database), 
			pq.QuoteIdentifier(grant.Spec.Schema), 
			pq.QuoteIdentifier(grant.Spec.Table), 
			pq.QuoteIdentifier(grant.Spec.Role))
	} else if grant.Spec.Schema != "" {
		stmt = fmt.Sprintf("GRANT %s ON SCHEMA %s.%s TO %s", 
			privileges, 
			pq.QuoteIdentifier(grant.Spec.Database), 
			pq.QuoteIdentifier(grant.Spec.Schema), 
			pq.QuoteIdentifier(grant.Spec.Role))
	} else {
		stmt = fmt.Sprintf("GRANT %s ON DATABASE %s TO %s", 
			privileges, 
			pq.QuoteIdentifier(grant.Spec.Database), 
			pq.QuoteIdentifier(grant.Spec.Role))
	}

	_, err := db.Exec(stmt)
	return err
}

// ListRoles lists all roles in the database
func (db *DB) ListRoles() ([]string, error) {
	rows, err := db.Query("SELECT rolname FROM pg_roles WHERE rolname !~ '^pg_'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// ListUsers lists all users in the database
func (db *DB) ListUsers() ([]string, error) {
	rows, err := db.Query("SELECT usename FROM pg_user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}