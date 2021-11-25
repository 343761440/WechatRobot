package common

// MySQLOptions
type MySQLOptions struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

// SQLITEOptions
type SQLITEOptions struct {
	Path string
}

type dbOptions struct {
	dialect       string
	mysqlOptions  *MySQLOptions
	sqliteOptions *SQLITEOptions
}

var defaultDBOptions = dbOptions{
	dialect: "mysql",
	mysqlOptions: &MySQLOptions{
		User:     "sakura",
		Password: "sakura",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "sakura",
	},
}

// A ServerOption sets options such as credentials, codec and keepalive parameters, etc.
type DBOption interface {
	apply(*dbOptions)
}

// funcServerOption wraps a function that modifies serverOptions into an
// implementation of the ServerOption interface.
type funcDBOption struct {
	f func(*dbOptions)
}

func (fdo *funcDBOption) apply(do *dbOptions) {
	fdo.f(do)
}

func newFuncDBOption(f func(*dbOptions)) *funcDBOption {
	return &funcDBOption{
		f: f,
	}
}

// MysqlOptions
func MysqlOptions(options *MySQLOptions) DBOption {
	return newFuncDBOption(func(o *dbOptions) {
		o.dialect = "mysql"
		o.mysqlOptions = options
	})
}

// SqliteOptions
func SqliteOptions(options *SQLITEOptions) DBOption {
	return newFuncDBOption(func(o *dbOptions) {
		o.dialect = "sqlite3"
		o.sqliteOptions = options
	})
}
