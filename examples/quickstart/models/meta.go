// Auto generated by sqlw-mysql (https://github.com/huangjunwen/sqlw-mysql) default template.
// DON NOT EDIT.

package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gopkg.in/volatiletech/null.v6"
	"gopkg.in/volatiletech/null.v6/convert"
)

// TableMeta contains meta information of a database table.
type TableMeta struct {
	// Basic information.
	tableName   string
	columnNames []string // column pos -> column name

	// Optinal information.
	primaryColumnsPos []int  // len(primaryColumnsPos) == 0 if the table has no primary key
	autoIncColumnPos  int    // autoIncColumnPos == -1 if the table has no auto increment column
	hasDefault        []bool // len(hasDefault) == len(columnNames), true if the column has not NULL server side default

	// Pre calculate information.
	columnsPos map[string]int // column name -> column pos

	// Post calculate information.
	isPrimary            []bool // column pos -> true if the column is part of primary key
	primaryCond          string // "`id1`=? AND `id2`=? ..."
	selectQuery          string // "SELECT `col`, .... FROM `xxx`"
	deleteByPrimaryQuery string // "DELETE FROM `xxx` WHERE `id`=? AND ..."
}

// TableMetaOption is used in creating TableMeta.
type TableMetaOption func(*TableMeta)

// OptColumnsWithDefault sets the columns that have not NULL server side default, including:
//  - AUTO_INCREMENT
//  - NOW()
//  - Other not NULL constant defaults.
func OptColumnsWithDefault(columnNames ...string) TableMetaOption {
	return func(meta *TableMeta) {
		for _, columnName := range columnNames {
			pos, ok := meta.columnsPos[columnName]
			if !ok {
				panic(fmt.Errorf("Table `%s` has no column named `%s`", meta.tableName, columnName))
			}
			meta.hasDefault[pos] = true
		}
	}
}

// OptPrimaryColumns sets the primary key columns.
func OptPrimaryColumns(columnNames ...string) TableMetaOption {
	return func(meta *TableMeta) {
		for _, columnName := range columnNames {
			pos, ok := meta.columnsPos[columnName]
			if !ok {
				panic(fmt.Errorf("Table `%s` has no column named `%s`", meta.tableName, columnName))
			}
			meta.primaryColumnsPos = append(meta.primaryColumnsPos, pos)
		}
	}
}

// OptAutoIncColumn sets the auto increment column.
func OptAutoIncColumn(columnName string) TableMetaOption {
	return func(meta *TableMeta) {
		pos, ok := meta.columnsPos[columnName]
		if !ok {
			panic(fmt.Errorf("Table `%s` has no column named `%s`", meta.tableName, columnName))
		}
		meta.autoIncColumnPos = pos
	}
}

// NewTableMeta creates a new TableMeta.
func NewTableMeta(tableName string, columnNames []string, opts ...TableMetaOption) *TableMeta {
	meta := &TableMeta{
		tableName:         tableName,
		columnNames:       columnNames,
		primaryColumnsPos: nil,
		autoIncColumnPos:  -1,
		hasDefault:        make([]bool, len(columnNames)), // All false
		columnsPos:        make(map[string]int),
		isPrimary:         make([]bool, len(columnNames)), // All false
	}

	// --- Pre calculate ---
	for i, columnName := range columnNames {
		meta.columnsPos[columnName] = i
	}

	// --- Apply options ---
	for _, opt := range opts {
		opt(meta)
	}

	// --- Post calculate ---
	if len(meta.primaryColumnsPos) == 0 {
		return meta
	}

	// isPrimary
	for _, pos := range meta.primaryColumnsPos {
		meta.isPrimary[pos] = true
	}

	// primaryCond
	conds := []string{}
	for _, pos := range meta.primaryColumnsPos {
		conds = append(conds, fmt.Sprintf("`%s`=?", meta.columnNames[pos]))
	}
	meta.primaryCond = strings.Join(conds, " AND ")

	// selectQuery
	cols := []string{}
	for _, columnName := range meta.columnNames {
		cols = append(cols, fmt.Sprintf("`%s`", columnName))
	}
	meta.selectQuery = fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(cols, ", "), meta.tableName)

	// deleteByPrimaryQuery
	meta.deleteByPrimaryQuery = fmt.Sprintf("DELETE FROM `%s` WHERE %s", meta.tableName, meta.primaryCond)

	return meta

}

// TableName returns the table name.
func (meta *TableMeta) TableName() string {
	return meta.tableName
}

// insertTR inserts `tr` into table. `tr` must be valid. Rules:
//   - If column has non zero value, then this column will be inserted.
//   - Othewise if the column doesn't have not NULL default, this column will be inserted.
//
// If `tr` has auto increment column then insertTR will also update it.
//
// `modifier` can be one of:
//   - "ignore": "INSERT IGNORE INTO ..."
//   - "replace": "REPLACE INTO ..."
//   - "": Normal "INSERT INTO ..."
func insertTR(ctx context.Context, e Execer, tr TableRow, modifier string) error {

	meta := tr.TableMeta()
	if !tr.Valid() {
		return fmt.Errorf("Insert: row is invalid")
	}

	cols := []string{}
	phs := []string{}
	args := []interface{}{}

	// Choose columns to insert.
	for i, col := range meta.columnNames {
		val := tr.ColumnValue(i)
		if isZero(val) && meta.hasDefault[i] {
			// Skip column value that is zero value and has not null default.
			continue
		}
		cols = append(cols, col)
		phs = append(phs, "?")
		args = append(args, val)
	}

	// Modifier.
	verb := ""
	switch modifier {
	case "ignore":
		verb = "INSERT IGNORE INTO"
	case "replace":
		verb = "REPLACE INTO"
	default:
		verb = "INSERT INTO"
	}

	// Construct query.
	query := ""
	if len(cols) == 0 {
		query = fmt.Sprintf("%s `%s` () VALUES ()", verb, meta.tableName)
	} else {
		query = fmt.Sprintf("%s `%s` (`%s`) VALUES (%s)", verb, meta.tableName, strings.Join(cols, "`, `"), strings.Join(phs, ", "))
	}

	// Query.
	result, err := e.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// If has auto increment column and the column's value is zero.
	if meta.autoIncColumnPos >= 0 && isZero(tr.ColumnValue(meta.autoIncColumnPos)) {
		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		if err := convert.ConvertAssign(tr.ColumnPointer(meta.autoIncColumnPos), lastInsertId); err != nil {
			return err
		}
	}

	return nil

}

// updateTR updates `tr` to `newTr` in table by its primary key. `tr`/`newTr` must be valid.
// Rules:
//   - Only columns that have different values will be updated.
func updateTR(ctx context.Context, e Execer, tr, newTr TableRowWithPrimary) error {

	// Check the table rows' validity.
	meta := tr.TableMeta()
	if meta != newTr.TableMeta() {
		return fmt.Errorf("Update: tr/newTr must be rows of the same table")
	}
	if !tr.Valid() || !newTr.Valid() {
		return fmt.Errorf("Update: row(s) is(are) invalid")
	}

	asgmts := []string{}
	args := []interface{}{}

	// Choose columns to update.
	for i, col := range meta.columnNames {

		if meta.isPrimary[i] {
			// Skip primary column.
			continue
		}

		val := tr.ColumnValue(i)
		newVal := newTr.ColumnValue(i)
		if val == newVal {
			// Skip if no difference.
			continue
		}

		asgmts = append(asgmts, fmt.Sprintf("`%s`=?", col))
		args = append(args, newVal)

	}

	// Nothing to update.
	if len(args) == 0 {
		return nil
	}

	// Construct query.
	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", meta.tableName, strings.Join(asgmts, ", "), meta.primaryCond)

	// Append primary value.
	appendPrimaryValues(&args, tr)

	// Query.
	_, err := e.ExecContext(ctx, query, args...)
	return err

}

// deleteTR deletes `tr` in table by its primary key. `tr` must be valid.
func deleteTR(ctx context.Context, e Execer, tr TableRowWithPrimary) error {

	// Check the table row's validity.
	meta := tr.TableMeta()
	if !tr.Valid() {
		return fmt.Errorf("Delete: row is invalid")
	}

	// Prepare args.
	args := []interface{}{}
	appendPrimaryValues(&args, tr)

	// Query.
	_, err := e.ExecContext(ctx, meta.deleteByPrimaryQuery, args...)
	return err
}

// selectTR selects `tr` from table by its primary key. `tr` must be valid.
// If `lock` is true, then "SELECT ... FOR UPDATE" will be used.
// It returns no error if a row is successfully found and returned.
func selectTR(ctx context.Context, q Queryer, tr TableRowWithPrimary, lock bool) error {
	return selectTRCond(ctx, q, tr, lock, "")
}

// selectTRCond selects `tr` from table by custom condition. `tr` must be valid.
// If `lock` is true, then "SELECT ... FOR UPDATE" will be used.
// It returns no error if a row is successfully found and returned.
func selectTRCond(ctx context.Context, q Queryer, tr TableRowWithPrimary, lock bool, cond string, args ...interface{}) error {

	// Check the table row's validity.
	meta := tr.TableMeta()
	if !tr.Valid() {
		return fmt.Errorf("SelectOne: row is invalid")
	}

	// If cond is empty, use primary condition.
	if cond == "" {
		if len(args) != 0 {
			return fmt.Errorf("SelectTRCond: cond is empty but args is not")
		}
		cond = meta.primaryCond
		appendPrimaryValues(&args, tr)
	}

	// Construct query.
	query := fmt.Sprintf("%s WHERE %s", meta.selectQuery, cond)
	if lock {
		query += " FOR UPDATE"
	}

	// Query.
	row := q.QueryRowContext(ctx, query, args...)

	// Process result.
	result := []interface{}{}
	tr.nxPreScan(&result)
	if err := row.Scan(result...); err != nil {
		return err
	}

	return tr.nxPostScan()

}

func appendPrimaryValues(dest *[]interface{}, tr TableRowWithPrimary) {
	meta := tr.TableMeta()
	for _, pos := range meta.primaryColumnsPos {
		*dest = append(*dest, tr.ColumnValue(pos))
	}
}

func isZero(val interface{}) bool {

	switch v := val.(type) {
	case float32:
		return v == 0
	case float64:
		return v == 0
	case bool:
		return v == false
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case time.Time:
		return v.IsZero()
	case string:
		return v == ""
	case null.Float32:
		return !v.Valid
	case null.Float64:
		return !v.Valid
	case null.Bool:
		return !v.Valid
	case null.Int8:
		return !v.Valid
	case null.Int16:
		return !v.Valid
	case null.Int32:
		return !v.Valid
	case null.Int64:
		return !v.Valid
	case null.Uint8:
		return !v.Valid
	case null.Uint16:
		return !v.Valid
	case null.Uint32:
		return !v.Valid
	case null.Uint64:
		return !v.Valid
	case null.Time:
		return !v.Valid
	case null.String:
		return !v.Valid
	default:
		panic(fmt.Errorf("isZero: Not support test for type %T", val))
	}

}