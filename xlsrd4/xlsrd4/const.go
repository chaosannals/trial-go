package xlsrd4

const (
	// 这些相对于整个文件偏移量
	BIG_BLOCK_DEPOT_BLOCKS_POS     = 0x4C
	NUM_BIG_BLOCK_DEPOT_BLOCKS_POS = 0x2C // 块总量位置
	ROOT_START_BLOCK_POS           = 0x30 // 根开始块ID 位置
	SMALL_BLOCK_DEPOT_BLOCK_POS    = 0x3C //
	EXTENSION_BLOCK_POS            = 0x44
	NUM_EXTENSION_BLOCK_POS        = 0x48

	// 这些是相对与 data 各个块内 的偏移量
	SIZE_OF_NAME_POS = 0x40
	TYPE_POS         = 0x42
	START_BLOCK_POS  = 0x74
	SIZE_POS         = 0x78

	BIG_BLOCK_SIZE              = 0x200
	SMALL_BLOCK_SIZE            = 0x40
	PROPERTY_STORAGE_BLOCK_SIZE = 0x80

	SMALL_BLOCK_THRESHOLD = 0x1000

	REKEY_BLOCK = 0x400

	XLS_TYPE_FORMULA         = 0x0006
	XLS_TYPE_EOF             = 0x000A
	XLS_TYPE_EXTERNSHEET     = 0x0017
	XLS_TYPE_DEFINEDNAME     = 0x0018
	XLS_TYPE_DATEMODE        = 0x0022
	XLS_TYPE_EXTERNNAME      = 0x0023
	XLS_TYPE_FILEPASS        = 0x002F
	XLS_TYPE_FONT            = 0x0031
	XLS_TYPE_CODEPAGE        = 0x0042
	XLS_TYPE_SHEET           = 0x0085
	XLS_TYPE_PALETTE         = 0x0092
	XLS_TYPE_XF              = 0x00E0
	XLS_TYPE_MSODRAWINGGROUP = 0x00EB
	XLS_TYPE_SST             = 0x00FC
	XLS_TYPE_LABELSST        = 0x00FD
	XLS_TYPE_EXTERNALBOOK    = 0x01AE
	XLS_TYPE_NUMBER          = 0x0203
	XLS_TYPE_LABEL           = 0x0204
	XLS_TYPE_BOOLERR         = 0x0205
	XLS_TYPE_RK              = 0x027E
	XLS_TYPE_STYLE           = 0x0293
	XLS_TYPE_FORMAT          = 0x041E
	XLS_TYPE_BOF             = 0x0809
	XLS_TYPE_XFEXT           = 0x087D

	XLS_WORKBOOKGLOBALS = 0x0005
	XLS_WORKSHEET       = 0x0010

	XLS_BIFF8 = 0x0600
	XLS_BIFF7 = 0x0500

	MS_BIFF_CRYPTO_NONE = 0
	MS_BIFF_CRYPTO_XOR  = 1
	MS_BIFF_CRYPTO_RC4  = 2

	WORKSHEET_SHEETSTATE_VISIBLE    = "visible"
	WORKSHEET_SHEETSTATE_HIDDEN     = "hidden"
	WORKSHEET_SHEETSTATE_VERYHIDDEN = "veryHidden"
)