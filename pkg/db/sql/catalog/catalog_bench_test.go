package catalog_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/angles-n-daemons/popsql/pkg/db/kv/store"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/catalog/desc"
	"github.com/angles-n-daemons/popsql/pkg/db/sql/parser/scanner"
	"github.com/angles-n-daemons/popsql/pkg/test/assert"
)

// ai wrote most of this file.

// Common names and types for random generation
var (
	columnNames = []string{"id", "name", "email", "age", "created_at", "updated_at",
		"status", "price", "quantity", "description", "address", "phone", "user_id"}

	tableNamePrefixes = []string{"user", "product", "order", "customer", "account",
		"transaction", "invoice", "item", "category", "payment"}

	tableNameSuffixes = []string{"data", "info", "records", "details", "list", "items", "logs"}

	// Token types that can be converted to data types
	dataTypeTokens = []scanner.TokenType{
		scanner.DATATYPE_NUMBER, scanner.DATATYPE_STRING,
	}
)

// generateRandomTable creates a table with random columns and optionally sequences
// ai made this
func generateRandomTable() (*desc.Table, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random table name
	prefix := tableNamePrefixes[r.Intn(len(tableNamePrefixes))]
	suffix := tableNameSuffixes[r.Intn(len(tableNameSuffixes))]
	tableName := fmt.Sprintf("%s_%s", prefix, suffix)

	// Generate random columns (3-8 columns)
	numColumns := r.Intn(6) + 3
	columns := make([]*desc.Column, 0, numColumns)
	usedNames := make(map[string]bool)

	for i := 0; i < numColumns; i++ {
		// Generate unique column name
		var colName string
		for {
			if i < len(columnNames) {
				colName = columnNames[r.Intn(len(columnNames))]
			} else {
				colName = fmt.Sprintf("column_%d", i)
			}
			if !usedNames[colName] {
				usedNames[colName] = true
				break
			}
		}

		// Random data type
		tokenType := dataTypeTokens[r.Intn(len(dataTypeTokens))]
		dataType, err := desc.GetDataType(tokenType)
		if err != nil {
			return nil, err
		}
		col := desc.NewColumn(colName, dataType)

		// 20% chance to make it a sequence column
		if r.Float32() < 0.2 {
			col = desc.NewSequenceColumn(colName, fmt.Sprintf("%s_%s_seq", tableName, colName))
		}

		columns = append(columns, col)
	}

	// Generate primary key (either pick a column or generate a new internal key)
	var primaryKey []string
	if r.Float32() < 0.7 && len(columns) > 0 {
		// Choose random column as primary key
		keyIdx := r.Intn(len(columns))
		primaryKey = []string{columns[keyIdx].Name}
	}

	// Create the table
	table, err := desc.NewTable(tableName, columns, primaryKey)
	if err != nil {
		// If we couldn't create with random primary key, try with internal key
		table, err = desc.NewTable(tableName, columns, nil)
		if err != nil {
			return nil, err
		}
	}

	// Random table ID
	table.TID = (uint64(r.Intn(100000000) + 1))

	return table, nil
}

// Common prefixes and suffixes for sequence names
var (
	sequencePrefixes = []string{"seq", "id", "counter", "increment", "serial"}
	sequenceSuffixes = []string{"gen", "sequence", "counter", "val", "num"}
)

// generateRandomSequence creates a sequence with random name and ID
func generateRandomSequence() *desc.Sequence {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random sequence name
	prefix := sequencePrefixes[r.Intn(len(sequencePrefixes))]
	suffix := sequenceSuffixes[r.Intn(len(sequenceSuffixes))]
	entityName := tableNamePrefixes[r.Intn(len(tableNamePrefixes))]

	// Format as either "prefix_entity_suffix" or "entity_suffix"
	var seqName string
	if r.Float32() < 0.5 {
		seqName = fmt.Sprintf("%s_%s_%s", prefix, entityName, suffix)
	} else {
		seqName = fmt.Sprintf("%s_%s", entityName, suffix)
	}

	// Create the sequence
	seq := desc.NewSequence(seqName)

	// Set random ID between 1 and 1000
	seq.SID = (uint64(r.Intn(10000000000) + 1))

	// Occasionally start with a non-zero value (20% chance)
	if r.Float32() < 0.2 {
		seq.V = uint64(r.Intn(100))
	}

	return seq
}

func BenchmarkCatalogWrites(b *testing.B) {
	tables := []*desc.Table{}
	sequences := []*desc.Sequence{}
	store := store.NewMemStore()
	ct, err := catalog.NewManager(store)
	assert.NoError(b, err)
	for i := 0; i < 1000000; i++ {
		t, err := generateRandomTable()
		assert.NoError(b, err)
		tables = append(tables, t)

		s := generateRandomSequence()
		sequences = append(sequences, s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := rand.Intn(len(tables))
		catalog.Create(ct, tables[j])
		catalog.Create(ct, sequences[j])

	}
}
