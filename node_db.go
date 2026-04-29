package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed n8n-nodes.db
var nodeDBBytes []byte

type NodeType struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	GroupName   string `json:"group_name"`
	Version     int    `json:"version"`
}

type NodeSchema struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	GroupName   string `json:"group_name"`
	Version     int    `json:"version"`
	Inputs      string `json:"inputs"`
	Outputs     string `json:"outputs"`
	Properties  string `json:"properties"`
}

// InitNodeDB writes the embedded SQLite bytes to a temp file and opens it read-only.
// Returns a cleanup func — caller must defer it.
func InitNodeDB() (*sql.DB, func(), error) {
	tmp, err := os.CreateTemp("", "n8n-nodes-*.db")
	if err != nil {
		return nil, nil, fmt.Errorf("node db: create temp: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(nodeDBBytes); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return nil, nil, fmt.Errorf("node db: write temp: %w", err)
	}
	tmp.Close()

	dsn := fmt.Sprintf("file:%s?mode=ro&_journal_mode=off", tmpPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		os.Remove(tmpPath)
		return nil, nil, fmt.Errorf("node db: open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		os.Remove(tmpPath)
		return nil, nil, fmt.Errorf("node db: ping: %w", err)
	}
	db.SetMaxOpenConns(1)

	cleanup := func() {
		db.Close()
		os.Remove(tmpPath)
	}
	return db, cleanup, nil
}

// SearchNodes returns up to limit node types whose name or display_name contains any of the comma-separated keywords.
// If group is non-empty, results are filtered to that group_name value (t=trigger, i=action, o=output).
func SearchNodes(db *sql.DB, keywords, group string, limit int) ([]NodeType, error) {
	if limit <= 0 {
		limit = 20
	}

	parts := strings.Split(keywords, ",")
	var conditions []string
	var args []interface{}
	for _, kw := range parts {
		kw = strings.TrimSpace(kw)
		if kw == "" {
			continue
		}
		pattern := "%" + kw + "%"
		conditions = append(conditions, "(name LIKE ? OR display_name LIKE ?)")
		args = append(args, pattern, pattern)
	}
	if len(conditions) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(
		`SELECT name, display_name, group_name, version FROM node_types WHERE %s`,
		strings.Join(conditions, " OR "),
	)
	if group != "" {
		query += " AND group_name = ?"
		args = append(args, group)
	}
	query += " ORDER BY display_name LIMIT ?"
	args = append(args, limit)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []NodeType
	for rows.Next() {
		var n NodeType
		var versionRaw interface{}
		if err := rows.Scan(&n.Name, &n.DisplayName, &n.GroupName, &versionRaw); err != nil {
			return nil, err
		}
		n.Version = toInt(versionRaw)
		results = append(results, n)
	}
	return results, rows.Err()
}

// toInt converts a raw SQLite value to int. Some rows have "[object Object]" stored
// as the version column (JS serialization artifact) — those fall back to 1.
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
		return 1
	case []byte:
		if n, err := strconv.Atoi(string(val)); err == nil {
			return n
		}
		return 1
	}
	return 1
}

// GetNodeSchema returns the full record for a single node type by exact name.
func GetNodeSchema(db *sql.DB, name string) (*NodeSchema, error) {
	row := db.QueryRow(
		`SELECT name, display_name, group_name, version, inputs, outputs, properties
		 FROM node_types WHERE name = ?`,
		name,
	)
	var s NodeSchema
	var versionRaw interface{}
	if err := row.Scan(&s.Name, &s.DisplayName, &s.GroupName, &versionRaw,
		&s.Inputs, &s.Outputs, &s.Properties); err != nil {
		return nil, err
	}
	s.Version = toInt(versionRaw)
	return &s, nil
}
