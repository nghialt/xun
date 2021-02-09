package sql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/yaoapp/xun/grammar"
	"github.com/yaoapp/xun/utils"
)

// Builder the SQL builder
type Builder struct{}

// SQLTableExists return the SQL for checking table exists.
func (builder Builder) SQLTableExists(db *sqlx.DB, name string, quoter grammar.Quoter) string {
	return fmt.Sprintf("SHOW TABLES like %s", quoter.VAL(name, db))
}

// SQLRenameTable return the SQL for the renaming table.
func (builder Builder) SQLRenameTable(db *sqlx.DB, old string, new string, quoter grammar.Quoter) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME %s", quoter.ID(old, db), quoter.ID(new, db))
}

// SQLCreateColumn return the add column sql for table create
func (builder Builder) SQLCreateColumn(db *sqlx.DB, Column *grammar.Column, types map[string]string, quoter grammar.Quoter) string {
	// `id` bigint(20) unsigned NOT NULL,
	typ, has := types[Column.Type]
	if !has {
		typ = "VARCHAR"
	}
	if Column.Precision > 0 && Column.Scale > 0 {
		typ = fmt.Sprintf("%s(%d,%d)", typ, Column.Precision, Column.Scale)
	} else if Column.DatetimePrecision > 0 {
		typ = fmt.Sprintf("%s(%d)", typ, Column.DatetimePrecision)
	} else if Column.Length > 0 {
		typ = fmt.Sprintf("%s(%d)", typ, Column.Length)
	}

	nullable := utils.GetIF(Column.Nullable, "NOT NULL", " NULL").(string)
	defaultValue := utils.GetIF(Column.Default != nil, fmt.Sprintf("DEFAULT %v", Column.Default), "").(string)
	comment := utils.GetIF(Column.Comment != "", fmt.Sprintf("COMMENT %s", quoter.VAL(Column.Comment, db)), "").(string)
	collation := utils.GetIF(Column.Collation != "", fmt.Sprintf("COLLATE %s", Column.Collation), "").(string)
	sql := fmt.Sprintf(
		"%s %s %s %s %s %s %s",
		quoter.ID(Column.Name, db), typ, nullable, defaultValue, Column.Extra, comment, collation)

	sql = strings.Trim(sql, " ")
	return sql
}

// SQLCreateIndex  return the add index sql for table create
func (builder Builder) SQLCreateIndex(db *sqlx.DB, index *grammar.Index, indexTypes map[string]string, quoter grammar.Quoter) string {
	typ, has := indexTypes[index.Type]
	if !has {
		typ = "KEY"
	}

	// UNIQUE KEY `unionid` (`unionid`) COMMENT 'xxxx'
	columns := []string{}
	for _, Column := range index.Columns {
		columns = append(columns, quoter.ID(Column.Name, db))
	}

	comment := ""
	if index.Comment != "" {
		comment = fmt.Sprintf("COMMENT %s", quoter.VAL(index.Comment, db))
	}

	sql := fmt.Sprintf(
		"%s %s (%s) %s",
		typ, quoter.ID(index.Name, db), strings.Join(columns, "`,`"), comment)

	return sql
}
