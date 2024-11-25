package catalog

import "github.com/angles-n-daemons/popsql/pkg/kv/data"

type Manager struct {
	Schema Schema
	Store  data.Store
}
