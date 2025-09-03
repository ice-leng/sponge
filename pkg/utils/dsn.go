package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// AdaptiveMysqlDsn adaptation of various mysql format dsn address
func AdaptiveMysqlDsn(dsn string) string {
	// remove optional scheme prefix
	dsn = strings.ReplaceAll(dsn, "mysql://", "")

	// ensure a valid network/address section for go-sql-driver/mysql
	// Expected forms:
	//   user:pass@tcp(127.0.0.1:3306)/db
	//   user:pass@unix(/path/mysql.sock)/db
	// If it's like '@(127.0.0.1:3306)' → add 'tcp'
	// If it's like '@127.0.0.1:3306' → wrap to '@tcp(127.0.0.1:3306)'
	at := strings.Index(dsn, "@")
	if at != -1 {
		afterAt := dsn[at+1:]
		slashIdx := strings.Index(afterAt, "/")
		if slashIdx != -1 {
			addrPart := afterAt[:slashIdx]
			// If empty addrPart, nothing to fix
			if addrPart != "" {
				if strings.HasPrefix(addrPart, "(") {
					// missing protocol
					dsn = strings.Replace(dsn, "@(", "@tcp(", 1)
				} else if !(strings.HasPrefix(addrPart, "tcp(") || strings.HasPrefix(addrPart, "unix(")) {
					// no parentheses and no protocol → wrap with tcp()
					dsn = strings.Replace(dsn, "@"+addrPart, "@tcp("+addrPart+")", 1)
				}
			}
		}
	}

	// ensure the connection prefers utf8mb4 to avoid collation mismatch
	// issues with MySQL 8 (e.g. mixing utf8mb3_general_ci and utf8mb4_0900_ai_ci).
	qIdx := strings.Index(dsn, "?")
	if qIdx == -1 {
		// no query string → add charset parameter
		return dsn + "?charset=utf8mb4"
	}

	prefix := dsn[:qIdx]
	queryStr := dsn[qIdx+1:]
	parts := strings.Split(queryStr, "&")

	hasCharset := false
	hasCollation := false
	for i, p := range parts {
		if strings.HasPrefix(p, "charset=") {
			hasCharset = true
			val := strings.TrimPrefix(p, "charset=")
			// split by comma and de-duplicate while ensuring utf8mb4 comes first if present/added
			charsets := []string{}
			for _, cs := range strings.Split(val, ",") {
				cs = strings.TrimSpace(cs)
				if cs == "" {
					continue
				}
				// skip duplicates
				dup := false
				for _, existing := range charsets {
					if strings.EqualFold(existing, cs) {
						dup = true
						break
					}
				}
				if !dup {
					charsets = append(charsets, cs)
				}
			}

			// ensure utf8mb4 is present and at the first position
			containsUtf8mb4 := false
			for _, cs := range charsets {
				if strings.EqualFold(cs, "utf8mb4") {
					containsUtf8mb4 = true
					break
				}
			}
			if !containsUtf8mb4 {
				charsets = append([]string{"utf8mb4"}, charsets...)
			} else if len(charsets) > 0 && !strings.EqualFold(charsets[0], "utf8mb4") {
				// move utf8mb4 to front
				newOrder := []string{"utf8mb4"}
				for _, cs := range charsets {
					if !strings.EqualFold(cs, "utf8mb4") {
						newOrder = append(newOrder, cs)
					}
				}
				charsets = newOrder
			}

			parts[i] = "charset=" + strings.Join(charsets, ",")
			break
		}
		if strings.HasPrefix(p, "collation=") {
			hasCollation = true
		}
	}

	if !hasCharset {
		parts = append(parts, "charset=utf8mb4")
	}
	if !hasCollation {
		// default to a broadly compatible utf8mb4 collation
		parts = append(parts, "collation=utf8mb4_general_ci")
	}

	return prefix + "?" + strings.Join(parts, "&")
}

// AdaptivePostgresqlDsn convert postgres dsn to kv string
func AdaptivePostgresqlDsn(dsn string) string {
	if strings.Count(dsn, " ") > 3 {
		return dsn
	}

	if !strings.Contains(dsn, "postgres://") {
		dsn = "postgres://" + dsn
	}

	dsn = DeleteBrackets(dsn)

	u, err := url.Parse(dsn)
	if err != nil {
		panic(err)
	}

	password, _ := u.User.Password()

	if u.RawQuery == "" {
		u.RawQuery = "sslmode=disable"
	} else if u.Query().Get("sslmode") == "" {
		u.RawQuery = "sslmode=disable&" + u.RawQuery
	}
	ss := strings.Split(u.RawQuery, "&")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		u.Hostname(), u.Port(), u.User.Username(), password, u.Path[1:], strings.Join(ss, " "))
}

// AdaptiveSqlite adaptive sqlite
func AdaptiveSqlite(dbFile string) string {
	// todo convert to absolute path
	return dbFile
}

// AdaptiveMongodbDsn adaptive mongodb dsn
func AdaptiveMongodbDsn(dsn string) string {
	if !strings.Contains(dsn, "mongodb://") &&
		!strings.Contains(dsn, "mongodb+srv://") {
		dsn = "mongodb://" + dsn // default scheme
	}

	return DeleteBrackets(dsn)
}

// DeleteBrackets delete brackets in dsn
func DeleteBrackets(str string) string {
	start := strings.Index(str, "@(")
	end := strings.LastIndex(str, ")/")

	if start == -1 || end == -1 {
		return str
	}

	addr := str[start+2 : end]
	return strings.Replace(str, "@("+addr+")/", "@"+addr+"/", 1)
}
