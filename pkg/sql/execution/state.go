package execution

import (
	"github.com/angles-n-daemons/popsql/pkg/kv"
	"github.com/angles-n-daemons/popsql/pkg/sql/plan"
)

// NewState creates a new cursor struct and walks the plan
// tree, opening cursors and creating offsets for inlined values.
func NewState(st kv.Store, p plan.Plan) (*State, error) {
	c := &State{
		store:       st,
		cursors:     make(map[string]kv.Cursor),
		valueOffset: make(map[string]int),
	}
	_, err := plan.VisitPlan(p, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// State is a utility struct which walks a plan, and initializes
// the cursors which will be used by the scan nodes.
type State struct {
	store        kv.Store
	cursors      map[string]kv.Cursor
	valueOffset  map[string]int
	tableCreated bool
}

// While CreateTable cannot have a scan, there's no need to return anything.
func (c *State) VisitCreateTable(*plan.CreateTable) (any, error) { return nil, nil }

// While Insert cannot have a scan, there's no need to return anything.
func (c *State) VisitInsert(*plan.Insert) (any, error) { return nil, nil }

// For scan, we open a cursor on the specified table (or index),
// assign it to the scan node, and save the reference in the
// internal map.
func (c *State) VisitScan(sc *plan.Scan) (any, error) {
	span := sc.Table.Span()
	cursor, err := c.store.Scan(span.Start.Encode(), span.End.Encode())
	if err != nil {
		return nil, err
	}

	c.cursors[sc.ID] = cursor
	return nil, nil
}

// For values, we initialize an offset to keep track of during
// execution so that we know which row is being used.
func (c *State) VisitValues(vals *plan.Values) (any, error) {
	c.valueOffset[vals.ID] = 0
	return nil, nil
}
