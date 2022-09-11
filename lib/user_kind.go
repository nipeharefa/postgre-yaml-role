package lib

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"

	"github.com/nipeharefa/postgre-yaml-role/action"
	"gopkg.in/yaml.v3"
)

type (
	UserData struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Role     string `yaml:"roleRef"`
	}
	UserKindYml struct {
		Kind string   `yaml:"kind"`
		Data UserData `yaml:"data"`
	}

	userKind struct {
		db *sql.DB
	}
)

func NewUserKind(db *sql.DB) *userKind {

	return &userKind{
		db: db,
	}
}

func (u *userKind) Parser(ctx context.Context, r io.Reader, a action.Action) error {
	user := UserKindYml{}
	var buf bytes.Buffer

	buf.ReadFrom(r)

	yaml.Unmarshal(buf.Bytes(), &user)
	var sqlBuf bytes.Buffer

	sqlBuf.WriteString(`CREATE USER `)

	fmt.Fprintf(&sqlBuf, "%s ", user.Data.Username)

	if paswd := user.Data.Password; paswd != "" {
		fmt.Fprintf(&sqlBuf, `WITH PASSWORD '%s'`, paswd)
	}

	fmt.Fprintf(&sqlBuf, ";")

	if role := user.Data.Role; role != "" {
		fmt.Fprintf(&sqlBuf, "GRANT %s TO %s;", role, user.Data.Username)
	}

	res, err := u.db.ExecContext(ctx, sqlBuf.String())
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(res.RowsAffected())

	return nil
}

func (u *userKind) ParseAndExecute(ctx context.Context, r io.Reader) {

	// u.Parser()
}
