package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Name string
	DB   *sql.DB
}

func OpenPostgres(dsn *DSN) (*PostgresDB, error) {
	s := "host=" + dsn.Host + " port=" + dsn.Port + " user=" + dsn.User + " password=" + dsn.Password + " dbname=" + dsn.DBName + " sslmode=" + dsn.SSLMode
	db, err := sql.Open("postgres", s)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{DB: db}, nil
}

func (d *PostgresDB) Close() {
	_ = d.DB.Close()
}

func (d *PostgresDB) Ping() error {
	return d.DB.Ping()
}

func (d *PostgresDB) VacuumAnalyzeTable(table *Table) error {
	_, err := d.DB.ExecContext(context.TODO(), "VACUUM ANALYZE "+d.TableSQL(table))
	if err != nil {
		return err
	}
	return nil
}

func (d *PostgresDB) EncodeString(s string) string {
	var b strings.Builder
	b.WriteString("E'")
	for _, c := range s {
		switch c {
		case '\\':
			b.WriteString("\\\\")
		case '\'':
			b.WriteString("''")
		case '\b':
			b.WriteString("\\b")
		case '\f':
			b.WriteString("\\f")
		case '\n':
			b.WriteString("\\n")
		case '\r':
			b.WriteString("\\r")
		case '\t':
			b.WriteString("\\t")
		default:
			b.WriteRune(c)
		}
	}
	b.WriteRune('\'')
	return b.String()
}

func (d *PostgresDB) ExecMultiple(tx *sql.Tx, sql []string) error {
	var m strings.Builder
	for _, q := range sql {
		m.WriteString(q)
		m.WriteRune(';')
	}
	_, err := d.Exec(tx, m.String())
	if err != nil {
		return err
	}
	return nil
	// For Snowflake:
	//ctx, err := multiStatementContext(db.Type, 2)
	//if err != nil {
	//	return err
	//}
	//if _, err := tx.ExecContext(ctx, b.String()); err != nil {
	//	return err
	//}
}

func (d *PostgresDB) Exec(tx *sql.Tx, sql string) (sql.Result, error) {
	if tx == nil {
		result, err := d.DB.ExecContext(context.TODO(), sql)
		if err != nil {
			return nil, fmt.Errorf(d.Name + ": SQL: " + sql)
		}
		return result, nil
	} else {
		result, err := tx.ExecContext(context.TODO(), sql)
		if err != nil {
			return nil, fmt.Errorf(d.Name + ": SQL: " + sql)
		}
		return result, nil
	}
}

func (d *PostgresDB) Query(tx *sql.Tx, query string) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	if tx == nil {
		rows, err = d.DB.QueryContext(context.TODO(), query)
	} else {
		rows, err = tx.QueryContext(context.TODO(), query)
	}
	if err != nil {
		return nil, fmt.Errorf(d.Name + ": SQL: " + query)
	}
	return rows, nil
}

func (d *PostgresDB) QueryRow(tx *sql.Tx, query string) *sql.Row {
	if tx == nil {
		return d.DB.QueryRowContext(context.TODO(), query)
	}
	return tx.QueryRowContext(context.TODO(), query)
}

func (d *PostgresDB) HistoryTableSQL(table *Table) string {
	return d.TableSQL(d.HistoryTable(table))
}

func (d *PostgresDB) HistoryTable(table *Table) *Table {
	return &Table{
		Schema: table.Schema,
		Table:  table.Table + "__",
	}
}

func (d *PostgresDB) TableSQL(table *Table) string {
	return d.IdentiferSQL(table.Schema) + "." + d.IdentiferSQL(table.Table)
}

func (d *PostgresDB) IdentiferSQL(id string) string {
	return "\"" + id + "\""
}

func (d *PostgresDB) AutoIncrementSQL() string {
	return "GENERATED BY DEFAULT AS IDENTITY"
}

func (d *PostgresDB) BeginTx() (*sql.Tx, error) {
	tx, err := d.DB.BeginTx(context.TODO(), &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

//type Postgres struct {
//	//Database *sql.DB
//}

//func OpenPostgres(dsn *DSN) (*DB, error) {
//	s := "host=" + dsn.Host + " port=" + dsn.Port + " user=" + dsn.User + " password=" + dsn.Password + " dbname=" + dsn.DBName + " sslmode=" + dsn.SSLMode
//	db, err := sql.Open("postgres", s)
//	if err != nil {
//		return nil, err
//	}
//	return &DB{DB: db, Type: &Postgres{}}, nil
//}
//
//func (d *Postgres) DB() *sql.DB {
//	return d.Database
//}

//func (d *Postgres) String() string {
//	return "postgresql"
//}
//
//func (d *Postgres) EncodeString(s string) string {
//	return encodeStringPostgres(s, true)
//}
//
//func (d *Postgres) Id(identifier string) string {
//	return "\"" + identifier + "\""
//}
//
//func (d *Postgres) Identity() string {
//	return "GENERATED BY DEFAULT AS IDENTITY"
//}
//
//func (d *Postgres) SupportsIndexes() bool {
//	return true
//}

//func (d *Postgres) CreateIndex(name string, table *Table, columns []string) string {
//	var clist strings.Builder
//	for i, c := range columns {
//		if i != 0 {
//			clist.WriteString(",")
//		}
//		clist.WriteString(d.Id(c))
//	}
//	return "CREATE INDEX " + name + " ON " + table.Id(d) + "(" + clist.String() + ")"
//}

//func (d *Postgres) JSONType() string {
//	return "JSON"
//}
