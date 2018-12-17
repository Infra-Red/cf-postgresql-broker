package pgp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
)

// PGP is a postgresql manipulation entity implementation
type PGP struct {
	conn   *sql.DB
	host   string
	port   string
	prefix string
}

// Credentials contains all information that is needed to connect to created databases
type Credentials struct {
	DBName   string `json:"dbname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	URL      string `json:"url"`
}

// defaultPort is PostgreSQL default port
const defaultPort = "5432"

// New creates a new PGP entity
func New(source string) (*PGP, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "postgresql" {
		return nil, errors.New("malformed url")
	}

	chunks := strings.SplitN(u.Host, ":", 2)
	if len(chunks) < 2 {
		chunks = append(chunks, defaultPort)
	}

	conn, err := sql.Open("postgres", source)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &PGP{
		conn:   conn,
		host:   chunks[0],
		port:   chunks[1],
		prefix: "sb_",
	}, nil
}

// CreateDB creates the named database
func (b *PGP) CreateDB(ctx context.Context, d string) (string, error) {
	dbname := b.dbname(d)
	_, err := b.conn.ExecContext(ctx, "CREATE DATABASE "+de(dbname))
	return dbname, err
}

// DropDB deletes the named database
func (b *PGP) DropDB(ctx context.Context, d string) error {
	dbname := b.dbname(d)
	if _, err := b.conn.ExecContext(ctx, "UPDATE pg_database SET datallowconn = $1 WHERE datname = $2", false, dbname); err != nil {
		return err
	}

	if _, err := b.conn.ExecContext(ctx, "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1", dbname); err != nil {
		return err
	}

	// TODO: drop users
	_, err := b.conn.Exec("DROP DATABASE " + de(dbname))
	return err
}

// CreateUser creates a user for the named database
func (b *PGP) CreateUser(ctx context.Context, d, u string) (*Credentials, error) {
	dbname := b.dbname(d)
	if !b.databaseExists(ctx, dbname) {
		return nil, fmt.Errorf("database %q doesn't exist", dbname)
	}

	username := b.username(u)
	password, err := password.Generate(20, 10, 0, false, false)
	if err != nil {
		return nil, err
	}

	if !b.userExists(ctx, username) {
		if _, err := b.conn.ExecContext(ctx, "CREATE USER "+de(username)+" WITH PASSWORD "+se(password)); err != nil {
			return nil, err
		}
	}

	if _, err := b.conn.ExecContext(ctx, "GRANT ALL PRIVILEGES ON DATABASE "+de(dbname)+" TO "+de(username)); err != nil {
		return nil, err
	}

	return &Credentials{
		DBName:   dbname,
		Username: username,
		Password: password,
		Host:     b.host,
		Port:     b.port,
		URL:      fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, b.host, b.port, dbname),
	}, nil
}

// DropUser removes the named user
func (b *PGP) DropUser(ctx context.Context, d, u string) error {
	dbname := b.dbname(d)
	username := b.username(u)
	if _, err := b.conn.Exec("REVOKE ALL PRIVILEGES ON DATABASE " + de(dbname) + " FROM " + de(username)); err != nil {
		return err
	}
	_, err := b.conn.Exec("DROP USER " + de(username))
	return err
}

func (b *PGP) DisableUser(ctx context.Context, d, u string) error {
	dbname := b.dbname(d)
	username := b.username(u)
	if _, err := b.conn.Exec("REVOKE ALL PRIVILEGES ON DATABASE " + de(dbname) + " FROM " + de(username)); err != nil {
		return err
	}
	_, err := b.conn.Exec("ALTER ROLE " + de(username) + " WITH NOLOGIN")
	return err
}

// dbname prefixes the named database name
func (b *PGP) dbname(d string) string {
	return b.prefix + d
}

// username prefixes the named username
func (b *PGP) username(u string) string {
	return b.prefix + u
}

// password generates a random password
/*func (b *PGP) password(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}*/

// userExists checks whether the named database exists
func (b *PGP) databaseExists(ctx context.Context, dbname string) bool {
	return b.exists(ctx, "pg_database", "datname", dbname)
}

// userExists checks whether the named user exists
func (b *PGP) userExists(ctx context.Context, username string) bool {
	return b.exists(ctx, "pg_user", "usename", username)
}

// exists checks whether the named column is exists in the provided table name
// and it equals to the specified value
func (b *PGP) exists(ctx context.Context, table, column, value string) bool {
	var num string
	b.conn.QueryRowContext(ctx, "SELECT 1 FROM "+de(table)+" WHERE "+de(column)+" = $1 LIMIT 1", value).Scan(&num)
	return num != ""
}

// de double-quotes the named string safely escaping it
func de(s string) string {
	return fmt.Sprintf("%q", s)
}

// se single-quotes the named string safely escaping it
func se(s string) string {
	return "'" + strings.Replace(s, "'", "\\'", -1) + "'"
}
