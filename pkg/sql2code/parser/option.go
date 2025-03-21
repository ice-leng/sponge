package parser

// NullStyle null type
type NullStyle int

// nolint
const (
	NullDisable NullStyle = iota
	NullInSql
	NullInPointer
)

// Option function
type Option func(*options)

type options struct {
	DBDriver       string
	FieldTypes     map[string]string // name:type
	Charset        string
	Collation      string
	JSONTag        bool
	JSONNamedType  int
	TablePrefix    string
	ColumnPrefix   string
	NoNullType     bool
	NullStyle      NullStyle
	Package        string
	GormType       bool
	ForceTableName bool
	IsEmbed        bool // is gorm.Model embedded
	IsWebProto     bool // true: proto file include router path and swagger info, false: normal proto file without router and swagger
	IsExtendedAPI  bool // true: extended api (9 api), false: basic api (5 api)

	IsCustomTemplate bool // true: custom extend template, false: sponge template
}

var defaultOptions = options{
	DBDriver:   "mysql",
	FieldTypes: map[string]string{},
	NullStyle:  NullInSql,
	Package:    "model",
}

// WithDBDriver set db driver
func WithDBDriver(driver string) Option {
	return func(o *options) {
		if driver != "" {
			o.DBDriver = driver
		}
	}
}

// WithFieldTypes set field types
func WithFieldTypes(fieldTypes map[string]string) Option {
	return func(o *options) {
		if fieldTypes != nil {
			o.FieldTypes = fieldTypes
		}
	}
}

// WithCharset set charset
func WithCharset(charset string) Option {
	return func(o *options) {
		o.Charset = charset
	}
}

// WithCollation set collation
func WithCollation(collation string) Option {
	return func(o *options) {
		o.Collation = collation
	}
}

// WithTablePrefix set table prefix
func WithTablePrefix(p string) Option {
	return func(o *options) {
		o.TablePrefix = p
	}
}

// WithColumnPrefix set column prefix
func WithColumnPrefix(p string) Option {
	return func(o *options) {
		o.ColumnPrefix = p
	}
}

// WithJSONTag set json tag, 0 for underscore, other values for hump
func WithJSONTag(namedType int) Option {
	return func(o *options) {
		o.JSONTag = true
		o.JSONNamedType = namedType
	}
}

// WithNoNullType set NoNullType
func WithNoNullType() Option {
	return func(o *options) {
		o.NoNullType = true
	}
}

// WithNullStyle set NullType
func WithNullStyle(s NullStyle) Option {
	return func(o *options) {
		o.NullStyle = s
	}
}

// WithPackage set package name
func WithPackage(pkg string) Option {
	return func(o *options) {
		o.Package = pkg
	}
}

// WithGormType will write type in gorm tag
func WithGormType() Option {
	return func(o *options) {
		o.GormType = true
	}
}

// WithForceTableName set forceFloats
func WithForceTableName() Option {
	return func(o *options) {
		o.ForceTableName = true
	}
}

// WithEmbed is embed gorm.Model
func WithEmbed() Option {
	return func(o *options) {
		o.IsEmbed = true
	}
}

// WithWebProto set proto file type
func WithWebProto() Option {
	return func(o *options) {
		o.IsWebProto = true
	}
}

// WithExtendedAPI set extended api
func WithExtendedAPI() Option {
	return func(o *options) {
		o.IsExtendedAPI = true
	}
}

// WithCustomTemplate set custom template
func WithCustomTemplate() Option {
	return func(o *options) {
		o.IsCustomTemplate = true
	}
}

func parseOption(options []Option) options {
	o := defaultOptions
	for _, f := range options {
		f(&o)
	}
	if o.NoNullType {
		o.NullStyle = NullDisable
	}
	return o
}
