package btree

type MemPage struct {
	isInit byte // true if previously initialized. must be first.

	intKey     byte // true if table b-trees. false for index b-trees
	intKeyLeaf byte // true if the leaf of an intKey table.

	pgno uint32 // page number for this page

	/*
		Only the first 8 bytes (above) are zeroed by pager.c when a new page
		is allocated. All the fields that follow must be initialized before use.
	*/

	leaf         byte // true if a leaf page
	hdrOffset    byte // 100 for page 1. 0 otherwise
	childPtrSize byte // 0 if leaf == 1. 4 if leaf == 0

	// -- unknown
	max1bytePayload byte //min(maxLocal, 127)
	nOverflow       byte // number of overflow cell bodies in aCell[]

	// -- unknown
	maxLocal uint16 // copy of BTShared.maxLocal or BTShared.maxLeaf
	minLocal uint16 // copy of BTShared.minLocal or BTShared.minLeaf

	cellOffset uint16 // index in aData of the first cell pointer
	nFree      int32  // number of free bytes on the page. -1 for unknown

	nCells uint16 // number of cells on this page, local and overflow
	// -- unknown
	maskPage uint16   // mask for the page offset
	aiOvfl   []uint16 // Insert the i-th overflow cell before the aiOvfl-th non-overflow cell
	apOvfl   []byte   // Pointers to the body of overflow cells

	// UNCOMMENT
	// pBt *BTShared // pointer to the BTreeShared that this page is a part of

	aData    []byte // pointer to the disk image of the page data
	aDataEnd []byte // one byte past the end of the entire page - not just the usable space, the entire page. used to prevent corruption-induced buffer overflow

	aCellIdx  []byte // cell index area
	aDataOfst []byte // Same as aData for leaves. aData+4 for interior

	// UNCOMMENT
	// pDbPage *DBPage // pager page handle
}

func (m *MemPage) xCellSize() uint16 {
	return 0
}

//func (m *MemPage) xParseCell()
