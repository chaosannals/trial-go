package xlsrd4

type xlsProtectionKind string

const (
	PROTECTION_INHERIT     xlsProtectionKind = "inherit"
	PROTECTION_PROTECTED   xlsProtectionKind = "protected"
	PROTECTION_UNPROTECTED xlsProtectionKind = "unprotected"
)

type xlsProtection struct {
	Locked xlsProtectionKind
	Hidden xlsProtectionKind
}
