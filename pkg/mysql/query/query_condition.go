// Package query is a library for mysql query, support for complex conditional paging queries.
// Deprecated: has been moved to package pkg/ggorm
package query

import (
	"fmt"
	"strings"
)

const (
	// Eq equal
	Eq = "eq"
	// Neq not equal
	Neq = "neq"
	// Gt greater than
	Gt = "gt"
	// Gte greater than or equal
	Gte = "gte"
	// Lt less than
	Lt = "lt"
	// Lte less than or equal
	Lte = "lte"
	// Like fuzzy lookup
	Like = "like"
	// In include
	In = "in"
	// NotIN not include
	NotIN = "notin"
	// IsNull is null
	IsNull = "isnull"
	// IsNotNull is not null
	IsNotNull = "isnotnull"

	// AND logic and
	AND string = "and"
	// OR logic or
	OR string = "or"
)

var expMap = map[string]string{
	Eq:        " = ",
	Neq:       " <> ",
	Gt:        " > ",
	Gte:       " >= ",
	Lt:        " < ",
	Lte:       " <= ",
	Like:      " LIKE ",
	In:        " IN ",
	NotIN:     " NOT IN ",
	IsNull:    " IS NULL ",
	IsNotNull: " IS NOT NULL ",

	"=":           " = ",
	"!=":          " <> ",
	">":           " > ",
	">=":          " >= ",
	"<":           " < ",
	"<=":          " <= ",
	"not in":      " NOT IN ",
	"is null":     " IS NULL ",
	"is not null": " IS NOT NULL ",
}

var logicMap = map[string]string{
	AND: " AND ",
	OR:  " OR ",

	"&":   " AND ",
	"&&":  " AND ",
	"|":   " OR ",
	"||":  " OR ",
	"AND": " AND ",
	"OR":  " OR ",

	"and:(": " AND ",
	"and:)": " AND ",
	"or:(":  " OR ",
	"or:)":  " OR ",
}

// Params query parameters
// Deprecated: moved to package pkg/gorm/query Params
type Params struct {
	Page int    `json:"page" form:"page" binding:"gte=0"`
	Size int    `json:"size" form:"size" binding:"gte=1"`
	Sort string `json:"sort,omitempty" form:"sort" binding:""`

	Columns []Column `json:"columns,omitempty" form:"columns"` // not required
}

// Column query info
// Deprecated: moved to package pkg/gorm/query Column
type Column struct {
	Name  string      `json:"name" form:"name"`   // column name
	Exp   string      `json:"exp" form:"exp"`     // expressions, default value is "=", support =, !=, >, >=, <, <=, like, in, notin, isnull, isnotnull
	Value interface{} `json:"value" form:"value"` // column value
	Logic string      `json:"logic" form:"logic"` // logical type, defaults to and when the value is null, with &(and), ||(or)
}

func (c *Column) checkValid() error {
	if c.Name == "" {
		return fmt.Errorf("field 'name' cannot be empty")
	}
	if c.Value == nil {
		v := expMap[strings.ToLower(c.Exp)]
		if v == " IS NULL " || v == " IS NOT NULL " {
			return nil
		}
		return fmt.Errorf("field 'value' cannot be nil")
	}
	return nil
}

// converting ExpType to sql expressions and LogicType to sql using characters
func (c *Column) convert() (string, error) {
	symbol := "?"
	if c.Exp == "" {
		c.Exp = Eq
	}
	if v, ok := expMap[strings.ToLower(c.Exp)]; ok { //nolint
		c.Exp = v
		switch c.Exp {
		case " LIKE ":
			val, ok1 := c.Value.(string)
			if !ok1 {
				return symbol, fmt.Errorf("invalid value type '%s'", c.Value)
			}
			l := len(val)
			if l > 2 {
				val2 := val[1 : l-1]
				val2 = strings.ReplaceAll(val2, "%", "\\%")
				val2 = strings.ReplaceAll(val2, "_", "\\_")
				val = string(val[0]) + val2 + string(val[l-1])
			}
			if strings.HasPrefix(val, "%") ||
				strings.HasPrefix(val, "_") ||
				strings.HasSuffix(val, "%") ||
				strings.HasSuffix(val, "_") {
				c.Value = val
			} else {
				c.Value = "%" + val + "%"
			}
		case " IN ", " NOT IN ":
			val, ok1 := c.Value.(string)
			if !ok1 {
				return symbol, fmt.Errorf("invalid value type '%s'", c.Value)
			}
			iVal := []interface{}{}
			ss := strings.Split(val, ",")
			for _, s := range ss {
				iVal = append(iVal, s)
			}
			c.Value = iVal
			symbol = "(?)"
		case " IS NULL ", " IS NOT NULL ":
			c.Value = nil
			symbol = ""
		}
	} else {
		return symbol, fmt.Errorf("unsported exp type '%s'", c.Exp)
	}

	if c.Logic == "" {
		c.Logic = AND
	} else {
		logic := strings.ToLower(c.Logic)
		if _, ok := logicMap[logic]; ok { //nolint
			c.Logic = logic
		} else {
			return symbol, fmt.Errorf("unsported logic type '%s'", c.Logic)
		}
	}

	return symbol, nil
}

// ConvertToPage converted to conform to gorm rules based on the page size sort parameter
// Deprecated: will be moved to package pkg/gorm/query ConvertToPage
func (p *Params) ConvertToPage() (order string, limit int, offset int) { //nolint
	page := NewPage(p.Page, p.Size, p.Sort)
	order = page.sort
	limit = page.size
	offset = page.page * page.size
	return //nolint
}

// ConvertToGormConditions conversion to gorm-compliant parameters based on the Columns parameter
// ignore the logical type of the last column, whether it is a one-column or multi-column query
// Deprecated: will be moved to package pkg/gorm/query ConvertToGormConditions
func (p *Params) ConvertToGormConditions() (string, []interface{}, error) {
	str := ""
	args := []interface{}{}
	l := len(p.Columns)
	if l == 0 {
		return "", nil, nil
	}

	isUseIN := true
	if l == 1 {
		isUseIN = false
	}
	field := p.Columns[0].Name

	for i, column := range p.Columns {
		if err := column.checkValid(); err != nil {
			return "", nil, err
		}

		symbol, err := column.convert()
		if err != nil {
			return "", nil, err
		}

		if i == l-1 { // ignore the logical type of the last column
			switch column.Logic {
			case "or:)", "and:)":
				str += column.Name + column.Exp + symbol + " ) "
			default:
				str += column.Name + column.Exp + symbol
			}
		} else {
			switch column.Logic {
			case "or:(", "and:(":
				str += " ( " + column.Name + column.Exp + symbol + logicMap[column.Logic]
			case "or:)", "and:)":
				str += column.Name + column.Exp + symbol + " ) " + logicMap[column.Logic]
			default:
				str += column.Name + column.Exp + symbol + logicMap[column.Logic]
			}
		}
		if column.Value != nil {
			args = append(args, column.Value)
		}
		// when multiple columns are the same, determine whether the use of IN
		if isUseIN {
			if field != column.Name {
				isUseIN = false
				continue
			}
			if column.Exp != expMap[Eq] {
				isUseIN = false
			}
		}
	}

	if isUseIN {
		str = field + " IN (?)"
		args = []interface{}{args}
	}

	return str, args, nil
}

func getExpsAndLogics(keyLen int, paramSrc string) ([]string, []string) { //nolint
	exps, logics := []string{}, []string{}
	param := strings.Replace(paramSrc, " ", "", -1)
	sps := strings.SplitN(param, "?", 2)
	if len(sps) == 2 {
		param = sps[1]
	}

	num := keyLen
	if num == 0 {
		return exps, logics
	}

	fields := []string{}
	kvs := strings.Split(param, "&")
	for _, kv := range kvs {
		if strings.Contains(kv, "page=") || strings.Contains(kv, "size=") || strings.Contains(kv, "sort=") {
			continue
		}
		fields = append(fields, kv)
	}

	// divide into num groups based on non-repeating keys, and determine in each group whether exp and logic exist
	group := map[string]string{}
	for _, field := range fields {
		split := strings.SplitN(field, "=", 2)
		if len(split) != 2 {
			continue
		}

		if _, ok := group[split[0]]; ok {
			// if exp does not exist, the default value of null is filled, and if logic does not exist, the default value of null is filled.
			exps = append(exps, group["exp"])
			logics = append(logics, group["logic"])

			group = map[string]string{}
			continue
		}
		group[split[0]] = split[1]
	}

	// handling the last group
	exps = append(exps, group["exp"])
	logics = append(logics, group["logic"])

	return exps, logics
}

// Conditions query conditions
type Conditions struct {
	Columns []Column `json:"columns" form:"columns" binding:"min=1"` // columns info
}

// ConvertToGorm conversion to gorm-compliant parameters based on the Columns parameter
// ignore the logical type of the last column, whether it is a one-column or multi-column query
// Deprecated: will be moved to package pkg/gorm/query ConvertToGorm
func (c *Conditions) ConvertToGorm() (string, []interface{}, error) {
	p := &Params{Columns: c.Columns}
	return p.ConvertToGormConditions()
}

// CheckValid check valid
func (c *Conditions) CheckValid() error {
	if len(c.Columns) == 0 {
		return fmt.Errorf("field 'columns' cannot be empty")
	}

	return nil
}
