package catalog

// This package is a bit of a mess. It's a collection of things that are
// required for the database to function, specifically with respect to its
// internal schema. The catalog package is responsible for managing the schema,
// primarily through the Manager object. The Manager object is a composite
// object that contains the schema, the system objects, and the store. The
// schema is a collection of tables and sequences, and the system objects are
// the tables and sequences that are required for the database to function.

// The package layout:
// └── /catalog
//     ├── manager.go     - Responsible for schema management and persistence.
//     ├── sequence.go    - Functions for managing sequences in the schema.
//     ├── table.go       - Functions for managing tables in the schema.
//     ├── bootstrap.go   - Responsible for all logic around bootstrapping a fresh db.
//     ├── /desc          - Descriptor types and utility functions.
//     └── /schema        - Schema object and utility functions.

// Manager hierarchy:
// Manager
// ├── SystemObjects
// │   ├── MetaTable
// │   └── Sequences...
// ├── Schema
// │   ├── Tables
// │   └── Sequences (duplicates SystemObjects)
// └── Store
//     └── kv.Store (persists the schema)
