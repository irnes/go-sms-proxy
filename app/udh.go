package app

// NewUDH create newd
func NewUDH(a, b, c, d, e byte) UDH {
	return UDH{a, b, c, d, e}
}

// UDH represents User Data Header
type UDH [5]byte
