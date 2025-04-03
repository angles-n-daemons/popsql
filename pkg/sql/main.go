package sql

import (
	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/catalog"
	"github.com/angles-n-daemons/popsql/pkg/sql/exec"
)

func NewExecutor(st kv.Store, cat *catalog.Manager) *exec.Executor {
	return &exec.Executor{
		Store:   st,
		Catalog: cat,
	}
}
