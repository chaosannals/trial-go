package xlsrd2

type XlsBookProperty struct {
	TypeId   int
	StartPos int32
	Name     string
	Size     int32
}

type XlsBookPropertySets struct {
	WorkBookId            *int
	RootEntryId           *int
	SummaryInfoId         *int
	DocumentSummaryInfoId *int
	All                   []XlsBookProperty
}

type XlsBook struct {
	// bigBlockDepotBlockCount     int
	// rootStartBlockId            int
	// smallBlockDepotStartBlockId int
	// extensionBlockId            int
	// extensionBlockCount         int

	bigBlockChain   []byte
	smallBlockChain []byte
	// data            []byte
	PropertySets *XlsBookPropertySets
}
