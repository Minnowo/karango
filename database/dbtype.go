package database

type DBType int

const (
	POSTGRES DBType = iota
)

func DBTypeFromStr(s string) DBType {

	switch s {
	case "postgres":
		return POSTGRES
	}

	return POSTGRES
}
