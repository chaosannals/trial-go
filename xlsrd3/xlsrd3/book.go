package xlsrd3

type XlsBook struct {
}

func LoadFormFile(xlsPath string) (book *XlsBook, err error) {
	reader, err := readXlsFile(xlsPath)
	if err != nil {
		return
	}

	reader.traceInfo()

	book = &XlsBook{}
	return
}
