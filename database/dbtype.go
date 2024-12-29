package database

type DBType int

const (
	POSTGRES DBType = iota
	MOCK     DBType = iota
)

func DBTypeFromStr(s string) DBType {

	switch s {
	case "postgres":
		return POSTGRES
	case "mock":
		return MOCK
	}

	return POSTGRES
}
