package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/angles-n-daemons/popsql/pkg/db"
	"github.com/angles-n-daemons/popsql/pkg/db/execution"
	wire "github.com/jeroenrinzema/psql-wire"
	"github.com/lib/pq/oid"
)

func Server() {
	s := &server{
		db: db.GetEngine(),
	}

	s.ListenAndServe()
}

type server struct {
	db *db.Engine
}

func (s *server) ListenAndServe() {
	wire.ListenAndServe("127.0.0.1:5432", s.handler)
}

func toQueries(query string) []string {
	queries := []string{}
	for _, query := range strings.Split(query, ";") {
		if strings.TrimSpace(query) == "" {
			continue
		}
		queries = append(queries, query)
	}
	return queries
}

func toWireColumns(columns []string, row execution.Row) (wire.Columns, error) {
	wCols := wire.Columns{}
	for i, col := range columns {
		var wType oid.Oid
		switch row[i].(type) {
		case int:
			wType = oid.T_int8
		case float64:
			wType = oid.T_float8
		case string:
			wType = oid.T_text
		case bool:
			wType = oid.T_bool
		default:
			return nil, fmt.Errorf("unable to create wire type for column %s, %d type %T", col, i, row[i])
		}
		wCols = append(wCols, wire.Column{
			Table: 0,
			Name:  col,
			Oid:   wType,
			Width: 256,
		})

	}
	return wCols, nil
}

func (s *server) handler(ctx context.Context, query string) (wire.PreparedStatements, error) {
	queries := toQueries(query)
	var result *execution.Result
	var err error
	for _, q := range queries {
		result, err = s.db.Query(q, nil)
		if err != nil {
			return nil, err
		}
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("no rows returned")
	}

	write := func(ctx context.Context, writer wire.DataWriter, parameters []wire.Parameter) error {
		for _, row := range result.Rows {
			if err := writer.Row(row); err != nil {
				return writer.Complete(err.Error())
			}
		}
		return writer.Complete("OK")
	}

	cols, err := toWireColumns(result.Columns, result.Rows[0])
	if err != nil {
		return nil, err
	}

	return wire.Prepared(wire.NewStatement(write, wire.WithColumns(cols))), nil
}
