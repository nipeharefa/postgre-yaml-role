package types

// Role represents a PostgreSQL role
type Role struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       RoleSpec `yaml:"spec"`
}

// User represents a PostgreSQL user
type User struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       UserSpec `yaml:"spec"`
}

// Grant represents a PostgreSQL grant/permission
type Grant struct {
	APIVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       GrantSpec `yaml:"spec"`
}

// Metadata contains common metadata for all resources
type Metadata struct {
	Name string `yaml:"name"`
}

// RoleSpec defines the specification of a role
type RoleSpec struct {
	Login           bool   `yaml:"login"`
	Inherit         bool   `yaml:"inherit"`
	CreateRole      bool   `yaml:"createrole"`
	CreateDB        bool   `yaml:"createdb"`
	Replication     bool   `yaml:"replication"`
	BypassRLS       bool   `yaml:"bypassrls"`
	ConnectionLimit int    `yaml:"connectionLimit"`
	ValidUntil      string `yaml:"validUntil,omitempty"`
	Comment         string `yaml:"comment,omitempty"`
}

// UserSpec defines the specification of a user
type UserSpec struct {
	Login           bool     `yaml:"login"`
	Inherit         bool     `yaml:"inherit"`
	CreateRole      bool     `yaml:"createrole"`
	CreateDB        bool     `yaml:"createdb"`
	Replication     bool     `yaml:"replication"`
	BypassRLS       bool     `yaml:"bypassrls"`
	ConnectionLimit int      `yaml:"connectionLimit"`
	Password        string   `yaml:"password,omitempty"`
	ValidUntil      string   `yaml:"validUntil,omitempty"`
	Comment         string   `yaml:"comment,omitempty"`
	InRoles         []string `yaml:"inRoles,omitempty"`
}

// GrantSpec defines the specification of a grant
type GrantSpec struct {
	Role       string   `yaml:"role"`
	Database   string   `yaml:"database"`
	Schema     string   `yaml:"schema,omitempty"`
	Table      string   `yaml:"table,omitempty"`
	Privileges []string `yaml:"privileges"`
}