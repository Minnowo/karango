package database

type DBFlag int

const (
	FLAG_INVALID              DBFlag = -1
	FLAG_DEFAULT_DATA_CREATED DBFlag = iota
)

func (f DBFlag) String() string {

	switch f {
	case FLAG_DEFAULT_DATA_CREATED:
		return "FLAG_DEFAULT_DATA_CREATED"

	case FLAG_INVALID:
	default:
		return "FLAG_INVALID"
	}

	panic("unreachable")
}
