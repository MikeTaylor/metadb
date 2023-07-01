package server

import (
	"fmt"
	"strconv"

	"github.com/metadb-project/metadb/cmd/metadb/catalog"
	"github.com/metadb-project/metadb/cmd/metadb/command"
	"github.com/metadb-project/metadb/cmd/metadb/dbx"
	"github.com/metadb-project/metadb/cmd/metadb/sqlx"
	"github.com/metadb-project/metadb/cmd/metadb/sysdb"
)

func addTable(cmd *command.Command, cat *catalog.Catalog, source string) error {
	table := dbx.Table{S: cmd.SchemaName, T: cmd.TableName}
	parentTable := dbx.Table{S: cmd.ParentTable.Schema, T: cmd.ParentTable.Table}
	return cat.CreateNewTable(&table, cmd.Transformed, &parentTable, source)
	/*
		// if tracked, then assume the table exists
		if cat.TableExists(table) {
			return nil
		}
		// create tables
		if err := createSchemaIfNotExists(sqlx.NewTable(cmd.SchemaName, cmd.TableName), db, users); err != nil {
			return err
		}
		if err := createMainTableIfNotExists(sqlx.NewTable(cmd.SchemaName, cmd.TableName), db, users); err != nil {
			return err
		}
		// track new table
		parentTable := dbx.Table{S: cmd.ParentTable.Schema, T: cmd.ParentTable.Table}
		if err := cat.AddTableEntry(table, cmd.Transformed, parentTable); err != nil {
			return err
		}
		return nil
	*/
}

func addPartition(cat *catalog.Catalog, cmd *command.Command) error {
	yearStr := cmd.SourceTimestamp[0:4]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return fmt.Errorf("adding partition for table %q: invalid year format: %q",
			cmd.SchemaName+"."+cmd.TableName, yearStr)
	}
	if err = cat.AddPartYearIfNotExists(cmd.SchemaName, cmd.TableName, year); err != nil {
		return fmt.Errorf("adding partition for table %q year %q: %v", cmd.SchemaName+"."+cmd.TableName,
			yearStr, err)
	}
	return nil
}

//func createSchemaIfNotExists(table *sqlx.Table, db sqlx.DB, users *cache.Users) error {
//	_, err := db.Exec(nil, "CREATE SCHEMA IF NOT EXISTS "+db.IdentiferSQL(table.Schema)+"")
//	if err != nil {
//		return err
//	}
//	for _, u := range users.WithPerm(table) {
//		_, err := db.Exec(nil, "GRANT USAGE ON SCHEMA "+db.IdentiferSQL(table.Schema)+" TO "+u+"")
//		if err != nil {
//			log.Warning("%s", err)
//		}
//	}
//	return nil
//}

/*
func createCurrentTableIfNotExists(table *sqlx.T, db sqlx.DB, users *cache.Users) error {
	_, err := db.Exec(nil, ""+
		"CREATE TABLE IF NOT EXISTS "+db.TableSQL(table)+" ("+
		"    __id bigint "+db.AutoIncrementSQL()+" PRIMARY KEY,"+
		"    __cf boolean NOT NULL DEFAULT TRUE,"+
		"    __start timestamp with time zone NOT NULL,"+
		"    __origin varchar(63) NOT NULL DEFAULT ''"+
		")")
	if err != nil {
		return err
	}
	// Add indexes on new columns.
	_, err = db.Exec(nil, "CREATE INDEX ON "+db.TableSQL(table)+" (__start)")
	if err != nil {
		return err
	}
	_, err = db.Exec(nil, "CREATE INDEX ON "+db.TableSQL(table)+" (__origin)")
	if err != nil {
		return err
	}
	// Grant permissions on new table.
	for _, u := range users.WithPerm(table) {
		_, err := db.Exec(nil, "GRANT SELECT ON "+db.TableSQL(table)+" TO "+u+"")
		if err != nil {
			return err
		}
	}
	return nil
}
*/

//func createMainTableIfNotExists(table *sqlx.Table, db sqlx.DB, users *cache.Users) error {
//	q := "CREATE TABLE IF NOT EXISTS " + db.HistoryTableSQL(table) + " (" +
//		"__id bigint GENERATED BY DEFAULT AS IDENTITY, " +
//		"__cf boolean NOT NULL DEFAULT TRUE, " +
//		"__start timestamp with time zone NOT NULL, " +
//		"__end timestamp with time zone NOT NULL, " +
//		"__current boolean NOT NULL, " +
//		"__origin varchar(63) NOT NULL DEFAULT ''" +
//		") PARTITION BY LIST (__current)"
//	if _, err := db.Exec(nil, q); err != nil {
//		return err
//	}
//	q = "CREATE TABLE IF NOT EXISTS " + db.TableSQL(table) + " PARTITION OF " + db.HistoryTableSQL(table) + " FOR VALUES IN (TRUE)"
//	if _, err := db.Exec(nil, q); err != nil {
//		return err
//	}
//	nctable := "\"" + table.Schema + "\".\"zzz___" + table.Table + "___\""
//	q = "CREATE TABLE IF NOT EXISTS " + nctable + " PARTITION OF " + db.HistoryTableSQL(table) + " FOR VALUES IN (FALSE) " +
//		"PARTITION BY RANGE (__start)"
//	if _, err := db.Exec(nil, q); err != nil {
//		return err
//	}
//	// Grant permissions on new tables.
//	for _, u := range users.WithPerm(table) {
//		if _, err := db.Exec(nil, "GRANT SELECT ON "+db.HistoryTableSQL(table)+" TO "+u+""); err != nil {
//			return err
//		}
//		if _, err := db.Exec(nil, "GRANT SELECT ON "+db.TableSQL(table)+" TO "+u+""); err != nil {
//			return err
//		}
//	}
//	return nil
//}

/*
func alterColumnVarcharSize(cat *catalog.Catalog, table *sqlx.Table, column string, datatype command.DataType, typesize int64, db sqlx.DB) error {
	var err error
	// Remove index if type size too large.
	if typesize > util.MaximumTypeSizeIndex {
		log.Trace("disabling index: value too large")
		_, err = db.Exec(nil, "DROP INDEX IF EXISTS "+db.IdentiferSQL(table.Schema)+"."+db.IdentiferSQL(indexName(db.HistoryTable(table).Table, column)))
		if err != nil {
			return fmt.Errorf("changing varchar size on column %q in table %q: drop index: %v",
				column, table, err)
		}
	}
	// Alter table.
	dtypesql, dataType, charMaxLen := command.DataTypeToSQL(datatype, typesize)
	_, err = db.Exec(nil, "ALTER TABLE "+db.HistoryTableSQL(table)+" ALTER COLUMN \""+column+"\" TYPE "+dtypesql)
	if err != nil {
		return fmt.Errorf("changing varchar size on column %q in table %q: alter column: %v",
			column, table, err)
	}
	// Update schema.
	cat.UpdateColumn(&sqlx.Column{Schema: table.Schema, Table: table.Table, Column: column}, dataType, charMaxLen)
	return nil
}
*/

/*
func alterColumnIntegerSize(table *sqlx.T, column string, typesize int64, db sqlx.DB, schema *cache.S) error {
	dtypesql, _, _ := command.DataTypeToSQL(command.IntegerType, typesize)
	_, err := db.Exec(nil, "ALTER TABLE "+db.TableSQL(table)+" ALTER COLUMN \""+column+"\" TYPE "+dtypesql)
	if err != nil {
		return err
	}
	_, err = db.Exec(nil, "ALTER TABLE "+db.HistoryTableSQL(table)+" ALTER COLUMN \""+column+"\" TYPE "+dtypesql)
	if err != nil {
		return err
	}
	// Update schema.
	schema.Update(&sqlx.Column{S: table.S, T: table.T, Column: column}, dtypesql, 0)
	return nil
}
*/

/*
func alterColumnToVarchar(table *sqlx.T, column string, typesize int64, db sqlx.DB, schema *cache.S) error {
	dtypesql, dataType, charMaxLen := command.DataTypeToSQL(command.TextType, typesize)
	_, err := db.Exec(nil, "ALTER TABLE "+db.TableSQL(table)+" ALTER COLUMN \""+column+"\" TYPE "+dtypesql+
		" USING \""+column+"\"::varchar")
	if err != nil {
		return err
	}
	_, err = db.Exec(nil, "ALTER TABLE "+db.HistoryTableSQL(table)+" ALTER COLUMN \""+column+"\" TYPE "+dtypesql+
		" USING \""+column+"\"::varchar")
	if err != nil {
		return err
	}
	// Update schema.
	schema.Update(&sqlx.Column{S: table.S, T: table.T, Column: column}, dataType, charMaxLen)
	return nil
}
*/

// Change column type to a specified new type, optionally casting data to the new type
func alterColumnType(cat *catalog.Catalog, db sqlx.DB, schema string, table string, column string, datatype command.DataType, typesize int64, cast bool) error {
	schemaTable := sqlx.NewTable(schema, table)
	sqltype := command.DataTypeToSQL(datatype, typesize)
	var caststr string
	if cast {
		caststr = " USING \"" + column + "\"::" + sqltype
	}
	var q = "ALTER TABLE %s ALTER COLUMN \"" + column + "\" TYPE " + sqltype + caststr
	if _, err := db.Exec(nil, fmt.Sprintf(q, db.HistoryTableSQL(schemaTable))); err != nil {
		return fmt.Errorf("changing type of column %q in table %q to %q: alter column: %v",
			column, table, sqltype, err)
	}
	// Update schema.
	cat.UpdateColumn(&sqlx.Column{Schema: schema, Table: table, Column: column}, sqltype)
	return nil
}

/*
func renameColumnOldType(table *sqlx.T, column string, datatype command.DataType, typesize int64, db sqlx.DB, schema *cache.S) error {
	var err error
	// Find new name for old column.
	var newName string
	if newName, err = newNumberedColumnName(table, column, schema); err != nil {
		return err
	}
	// Current table: rename column.
	_, err = db.Exec(nil, "ALTER TABLE "+db.TableSQL(table)+" RENAME COLUMN \""+column+"\" TO \""+newName+"\"")
	if err != nil {
		return err
	}
	// Current table: rename index.
	index := db.IdentiferSQL(table.S) + "." + db.IdentiferSQL(indexName(table.T, column))
	_, err = db.Exec(nil, "ALTER INDEX IF EXISTS "+index+" RENAME TO "+db.IdentiferSQL(indexName(table.T, newName)))
	if err != nil {
		return err
	}
	// History table: rename column.
	_, err = db.Exec(nil, "ALTER TABLE "+db.HistoryTableSQL(table)+" RENAME COLUMN \""+column+"\" TO \""+newName+"\"")
	if err != nil {
		return err
	}
	// History table: rename index.
	index = db.IdentiferSQL(table.S) + "." + db.IdentiferSQL(indexName(db.HistoryTable(table).T, column))
	_, err = db.Exec(nil, "ALTER INDEX IF EXISTS "+index+" RENAME TO "+db.IdentiferSQL(indexName(db.HistoryTable(table).T, newName)))
	if err != nil {
		return err
	}
	// Update schema.
	schema.Delete(&sqlx.Column{S: table.S, T: table.T, Column: column})
	_, dataType, charMaxLen := command.DataTypeToSQL(datatype, typesize)
	schema.Update(&sqlx.Column{S: table.S, T: table.T, Column: newName}, dataType, charMaxLen)
	return nil
}

func newNumberedColumnName(table *sqlx.T, column string, schema *cache.S) (string, error) {
	var columns []string = schema.TableColumns(table)
	maxn := 0
	regex := regexp.MustCompile(`^` + column + `__([0-9]+)$`)
	for _, c := range columns {
		n := 0
		var match []string = regex.FindStringSubmatch(c)
		if match != nil {
			var err error
			if n, err = strconv.Atoi(match[1]); err != nil {
				return "", fmt.Errorf("internal error: column number: strconf.Atoi(): %s", err)
			}
		}
		if n > maxn {
			maxn = n
		}
	}
	newName := fmt.Sprintf("%s__%d", column, maxn+1)
	return newName, nil
}
*/

func selectTableSchema(cat *catalog.Catalog, table *sqlx.Table) (*sysdb.TableSchema, error) {
	m := cat.TableSchema(table)
	ts := new(sysdb.TableSchema)
	for k, v := range m {
		name := k
		dtype, dtypesize := command.MakeDataType(v)
		cs := sysdb.ColumnSchema{Name: name, DType: dtype, DTypeSize: dtypesize}
		ts.Column = append(ts.Column, cs)
	}
	return ts, nil
}

func indexName(table, column string) string {
	return table + "_" + column + "_idx"
}
