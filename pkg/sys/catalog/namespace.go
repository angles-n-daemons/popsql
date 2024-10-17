package catalog

func NewNamespace(name string) *Namespace {
	return &Namespace{
		Name:          name,
		Tables:        map[string]*Table{},
		TableIdByName: map[string]string{},
	}
}

type Namespace struct {
	Name          string
	Tables        map[string]*Table
	TableIdByName map[string]string
}

func (n *Namespace) AddTable(t *Table) error {
	// error if table already exists
	return nil
}

func (n *Namespace) GetTable(id string) (*Table, error) {
	return nil, nil
}

func (n *Namespace) GetTableByName(id string) (*Table, error) {
	return nil, nil
}
