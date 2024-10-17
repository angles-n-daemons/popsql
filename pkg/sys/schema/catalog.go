package schema

type Catalog struct {
	TableLookup   map[string]*Table
	TableNameToID map[string]string
}

func (ts *Catalog) AddTable(name string) {

}

func (ts *Catalog) GetTable(name string) {

}
