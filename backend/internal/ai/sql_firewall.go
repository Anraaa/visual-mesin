package ai

import (
	"strings"
)

type FirewallResult struct {
	Valid  bool
	Reason string
	SQL    string
}

type SQLFirewall struct {
	blacklistWords []string
}

func NewSQLFirewall() *SQLFirewall {
	return &SQLFirewall{
		blacklistWords: []string{
			"INSERT", "UPDATE", "DELETE", "DROP", "ALTER",
			"TRUNCATE", "CREATE", "REPLACE", "GRANT", "REVOKE",
			"EXECUTE", "EXEC", "CALL", "LOAD", "INTO OUTFILE",
			"INTO DUMPFILE", "SLEEP", "BENCHMARK", "INFORMATION_SCHEMA.TABLES",
		},
	}
}

func (f *SQLFirewall) Validate(sql string) *FirewallResult {
	sql = strings.TrimSpace(sql)
	upper := strings.ToUpper(sql)

	if !strings.HasPrefix(upper, "SELECT") {
		return &FirewallResult{
			Valid:  false,
			Reason: "Hanya query SELECT yang diizinkan",
			SQL:    sql,
		}
	}

	for _, word := range f.blacklistWords {
		// Check for word boundary
		upperWithSpaces := " " + upper + " "
		if strings.Contains(upperWithSpaces, " "+word+" ") {
			return &FirewallResult{
				Valid:  false,
				Reason: "Query mengandung keyword terlarang: " + word,
				SQL:    sql,
			}
		}
		if strings.HasPrefix(upper, word+" ") {
			return &FirewallResult{
				Valid:  false,
				Reason: "Query mengandung keyword terlarang: " + word,
				SQL:    sql,
			}
		}
		if strings.HasSuffix(upper, " "+word) {
			return &FirewallResult{
				Valid:  false,
				Reason: "Query mengandung keyword terlarang: " + word,
				SQL:    sql,
			}
		}
		if upper == word {
			return &FirewallResult{
				Valid:  false,
				Reason: "Query mengandung keyword terlarang: " + word,
				SQL:    sql,
			}
		}
	}

	// Check for SQL comments which could hide malicious SQL
	if strings.Contains(sql, "--") || strings.Contains(sql, "/*") {
		return &FirewallResult{
			Valid:  false,
			Reason: "Query mengandung komentar SQL yang tidak diizinkan",
			SQL:    sql,
		}
	}

	// Check for semicolons (multiple statements)
	if strings.Count(sql, ";") > 1 {
		return &FirewallResult{
			Valid:  false,
			Reason: "Hanya satu statement yang diizinkan",
			SQL:    sql,
		}
	}

	return &FirewallResult{
		Valid:  true,
		Reason: "OK",
		SQL:    sql,
	}
}
