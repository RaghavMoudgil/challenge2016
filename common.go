package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func InsertCSVToMySQL(db *sql.DB, csvPath, tableName string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("cannot open CSV: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("cannot read CSV: %w", err)
	}
	if len(rows) < 2 {
		return fmt.Errorf("no data rows found")
	}

	headers := rows[0]
	placeholders := make([]string, len(headers))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	insert := fmt.Sprintf("INSERT INTO %s (`City Code`, `Province Code`, `Country Code`, `City Name`, `Province name`, `Country Name`) VALUES (?,?,?,?,?,?)", tableName)

	for i, row := range rows[1:] {
		if len(row) != len(headers) {
			return fmt.Errorf("row %d: mismatched column count", i+2)
		}
		vals := make([]interface{}, len(row))
		for j, val := range row {
			vals[j] = val
		}
		if _, err := db.Exec(insert, vals...); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert row %d failed: %w", i+2, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

type Permissions struct {
	Include map[string]struct{}
	Exclude map[string]struct{}
	Parent  *Permissions
}

func (p *Permissions) CanDistribute(region string) bool {
	if !p.includes(region) {
		return false
	}
	if p.excludes(region) {
		return false
	}
	return true
}

func (p *Permissions) includes(region string) bool {
	for r := region; r != ""; r = parentRegion(r) {
		if _, ok := p.Include[r]; ok {
			if p.Parent != nil {
				if p.Parent.includes(r) {
					return true
				}
			} else {
				return true
			}
		}
	}
	return false
}

func (p *Permissions) excludes(region string) bool {
	for r := region; r != ""; r = parentRegion(r) {
		if _, ok := p.Exclude[r]; ok {
			return true
		}
	}
	return false
}

func parentRegion(region string) string {
	idx := strings.Index(region, "-")
	if idx == -1 {
		return ""
	}
	return region[idx+1:]
}

type Distributor struct {
	ID          int
	Name        string
	ParentID    sql.NullInt64
	Permissions *Permissions
}

func loadDistributors(db *sql.DB) (map[string]*Distributor, error) {
	rows, err := db.Query("SELECT id, name, parent_id FROM distributors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	distributors := make(map[int]*Distributor)
	for rows.Next() {
		var d Distributor
		if err := rows.Scan(&d.ID, &d.Name, &d.ParentID); err != nil {
			return nil, err
		}
		d.Permissions = &Permissions{
			Include: make(map[string]struct{}),
			Exclude: make(map[string]struct{}),
		}
		distributors[d.ID] = &d
	}

	permRows, err := db.Query("SELECT distributor_id, permission_type, region_code FROM distributor_permissions")
	if err != nil {
		return nil, err
	}
	defer permRows.Close()

	for permRows.Next() {
		var distributorID int
		var permType, regionCode string
		if err := permRows.Scan(&distributorID, &permType, &regionCode); err != nil {
			return nil, err
		}
		dist, ok := distributors[distributorID]
		if !ok {
			continue
		}
		switch permType {
		case "INCLUDE":
			dist.Permissions.Include[regionCode] = struct{}{}
		case "EXCLUDE":
			dist.Permissions.Exclude[regionCode] = struct{}{}
		}
	}

	for _, dist := range distributors {
		if dist.ParentID.Valid {
			parent, ok := distributors[int(dist.ParentID.Int64)]
			if ok {
				dist.Permissions.Parent = parent.Permissions
			}
		}
	}

	distributorMap := make(map[string]*Distributor)
	for _, dist := range distributors {
		distributorMap[dist.Name] = dist
	}

	return distributorMap, nil
}
