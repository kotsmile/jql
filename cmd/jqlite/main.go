package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/kotsmile/jql/internal/engine"
	"github.com/kotsmile/jql/internal/lexer"
	"github.com/kotsmile/jql/internal/lexer/token"
	"github.com/kotsmile/jql/util"

	_ "github.com/mattn/go-sqlite3"
)

func processCmd(cmd string, db *sql.DB, logger util.Logger) error {
	if strings.HasPrefix(cmd, "load") {
		lexer := lexer.New(logger)
		lexer.Lex(cmd)
		tokens, err := lexer.Collect()
		if err != nil {
			return fmt.Errorf("failed to tokenize expression: %s", err)
		}

		filenameToken, ok := util.At(tokens, 1)
		if !ok {
			return fmt.Errorf("filename is not specified")
		}

		filename := filenameToken.Value()

		asToken, ok := util.At(tokens, 2)
		if !ok {
			return fmt.Errorf("as 'as' not specified")
		}
		if !(asToken.Is(token.Word) && asToken.Value() == "as") {
			return fmt.Errorf("as 'as' is not specified")
		}

		tableNameToken, ok := util.At(tokens, 3)
		if !ok {
			return fmt.Errorf("table name is not specified")
		}
		tableName := tableNameToken.Value()
		e := engine.New(os.Stdout)

		if err := e.LoadTable(filename, tableName); err != nil {
			return fmt.Errorf("failed to load table: %s", err)
		}

		table, _ := e.GetTable(tableName)
		typedColumns := table.ToSqliteTypes()

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("CREATE TABLE \"%s\" (\n", tableName))
		for i, c := range typedColumns {
			if i == len(typedColumns)-1 {
				sb.WriteString(fmt.Sprintf("\t\"%s\" %s\n", c.Name, c.SqliteType))
			} else {
				sb.WriteString(fmt.Sprintf("\t\"%s\" %s,\n", c.Name, c.SqliteType))
			}
		}
		sb.WriteString(");")
		logger.Debugf("creating table:\n%s", sb.String())
		if _, err = db.Exec(sb.String()); err != nil {
			return fmt.Errorf("failed to create table:\n%s", err)
		}

		for _, row := range table.Rows() {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("INSERT INTO \"%s\" (\n", tableName))
			for i, c := range typedColumns {
				if i == len(typedColumns)-1 {
					sb.WriteString(fmt.Sprintf("\t\"%s\"\n", c.Name))
				} else {
					sb.WriteString(fmt.Sprintf("\t\"%s\",\n", c.Name))
				}
			}

			sb.WriteString(") VALUES (\n")

			for i, c := range typedColumns {
				value, ok := row[c.Name]

				if ok {
					if c.SqliteType == "TEXT" {
						valueJSONBytes, err := json.Marshal(value)
						if err != nil {
							return fmt.Errorf("failed to marshal value: %s", err)
						}

						valueJSON := string(valueJSONBytes)
						valueString, ok := value.(string)
						if ok {
							valueJSON = valueString
						}

						if i == len(typedColumns)-1 {
							sb.WriteString(fmt.Sprintf("\t'%s'\n", valueJSON))
						} else {
							sb.WriteString(fmt.Sprintf("\t'%s',\n", valueJSON))
						}
					} else {
						if i == len(typedColumns)-1 {
							sb.WriteString(fmt.Sprintf("\t%v\n", value))
						} else {
							sb.WriteString(fmt.Sprintf("\t%v,\n", value))
						}
					}
				} else {
					if i == len(typedColumns)-1 {
						sb.WriteString("\tNULL\n")
					} else {
						sb.WriteString("\tNULL,\n")
					}
				}
			}
			sb.WriteString(");")

			logger.Debugf("inserting into table:\n%s", sb.String())
			if _, err = db.Exec(sb.String()); err != nil {
				return fmt.Errorf("failed to insert into table:\n%s", err)
			}
		}

		return nil
	} else if strings.HasPrefix(cmd, "save") {
		return fmt.Errorf("save is not implemented yet")
	}

	rows, err := db.Query(cmd)
	if err != nil {
		return fmt.Errorf("failed to execute query: %s", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %s", err)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	defer writer.Flush()

	for _, col := range columns {
		fmt.Fprintf(writer, "%s\t", col)
	}
	fmt.Fprintln(writer)

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %s", err)
		}

		for _, val := range values {
			var value string
			if b, ok := val.([]byte); ok {
				value = string(b)
			} else {
				value = fmt.Sprintf("%v", val)
			}
			fmt.Fprintf(writer, "%s\t", value)
		}
		fmt.Fprintln(writer)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error during rows iteration: %s", err)
	}

	return nil
}

func main() {
	debug := false

	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	logger := util.NewLogger(debug)

	file, err := os.Create(".jqlite.db")
	if err != nil {
		logger.Fatalf("failed to create file: %s", err)
	}
	defer os.Remove(file.Name())

	db, err := sql.Open("sqlite3", file.Name())
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// if debug {
	// 	execute("load \"./examples/simple.json\" as simple;", db, logger)
	// }

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if err := processCmd(cmd, db, logger); err != nil {
			logger.Errorf("failed to execute command: %s", err)
		}
	}
}
