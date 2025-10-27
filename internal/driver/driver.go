package driver

// Database driver constants
type Driver string

const (
	Postgres Driver = "postgres"
	MySQL    Driver = "mysql"
	SQLite   Driver = "sqlite3"
	MSSQL    Driver = "sqlserver"
)

// String returns the string representation of the driver
func (d Driver) String() string {
	return string(d)
}

// IsValid checks if the driver is a valid/supported driver
func (d Driver) IsValid() bool {
	switch d {
	case Postgres, MySQL, SQLite, MSSQL:
		return true
	default:
		return false
	}
}
