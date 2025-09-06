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

	dsn = ensureNetworkAddress(dsn)
	return ensureCharsetAndCollation(dsn)
}

// helper: ensure network/address section is valid for go-sql-driver/mysql
func ensureNetworkAddress(dsn string) string {
    at := strings.Index(dsn, "@")
    if at == -1 {
        return dsn
    }

    afterAt := dsn[at+1:]
    slashIdx := strings.Index(afterAt, "/")
    if slashIdx == -1 {
        return dsn
    }

    addrPart := afterAt[:slashIdx]
    if addrPart == "" {
        return dsn
    }

    if strings.HasPrefix(addrPart, "(") {
        // missing protocol, add tcp
        return strings.Replace(dsn, "@(", "@tcp(", 1)
    }

    if strings.HasPrefix(addrPart, "tcp(") || strings.HasPrefix(addrPart, "unix(") {
        return dsn
    }

    // no parentheses and no protocol → wrap with tcp()
    return strings.Replace(dsn, "@"+addrPart, "@tcp("+addrPart+")", 1)
}

// helper: ensure charset utf8mb4 and a reasonable collation are present
func ensureCharsetAndCollation(dsn string) string {
    qIdx := strings.Index(dsn, "?")
    if qIdx == -1 {
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
            parts[i] = "charset=" + normalizeCharsets(strings.TrimPrefix(p, "charset="))
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
        parts = append(parts, "collation=utf8mb4_general_ci")
    }

    return prefix + "?" + strings.Join(parts, "&")
}

// normalizeCharsets deduplicates a comma-separated charset list and ensures utf8mb4 is first
func normalizeCharsets(val string) string {
    pieces := strings.Split(val, ",")
    seen := map[string]bool{}
    ordered := []string{}
    for _, cs := range pieces {
        cs = strings.TrimSpace(cs)
        if cs == "" {
            continue
        }
        lower := strings.ToLower(cs)
        if seen[lower] {
            continue
        }
        seen[lower] = true
        ordered = append(ordered, cs)
    }

    // ensure utf8mb4 is present and at the front (case-insensitive)
    found := -1
    for i, cs := range ordered {
        if strings.EqualFold(cs, "utf8mb4") {
            found = i
            break
        }
    }
    if found == -1 {
        ordered = append([]string{"utf8mb4"}, ordered...)
    } else if found != 0 {
        // move to front
        front := []string{"utf8mb4"}
        for i, cs := range ordered {
            if i == found {
                continue
            }
            if strings.EqualFold(cs, "utf8mb4") {
                continue
            }
            front = append(front, cs)
        }
        ordered = front
    }

    return strings.Join(ordered, ",")
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
