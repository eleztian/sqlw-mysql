package infos

import (
	"fmt"

	"github.com/huangjunwen/sqlw-mysql/datasrc"
	"github.com/huangjunwen/sqlw-mysql/utils"
)

// DBInfo contains information of a database.
type DBInfo struct {
	tables     []*TableInfo
	tableNames map[string]int
}

// TableInfo contains information of a table.
type TableInfo struct {
	db            *DBInfo
	tableName     string
	columns       []*ColumnInfo
	columnNames   map[string]int
	indices       []*IndexInfo
	indexNames    map[string]int
	fks           []*FKInfo
	fkNames       map[string]int
	primary       *IndexInfo  // nil if not exists
	autoIncColumn *ColumnInfo // nil if not exists
	camel         utils.CamelName
}

// ColumnInfo contains information of a table column.
type ColumnInfo struct {
	table  *TableInfo
	column datasrc.Column
	camel  utils.CamelName
}

// IndexInfo contains information of an index.
type IndexInfo struct {
	table     *TableInfo
	indexName string
	columns   []*ColumnInfo
	isPrimary bool
	isUnique  bool
	camel     utils.CamelName
}

// FKInfo contains information of a foreign key constraint.
type FKInfo struct {
	fkName         string
	table          *TableInfo
	columns        []*ColumnInfo
	refTableName   string
	refColumnNames []string
	camel          utils.CamelName
}

// NewDBInfo extracts information from current database.
func NewDBInfo(loader *datasrc.Loader) (*DBInfo, error) {

	db := &DBInfo{
		tableNames: make(map[string]int),
	}

	tableNames, err := loader.LoadTableNames()
	if err != nil {
		return nil, err
	}

	for _, tableName := range tableNames {

		table := &TableInfo{
			db:          db,
			tableName:   tableName,
			columnNames: make(map[string]int),
			indexNames:  make(map[string]int),
			fkNames:     make(map[string]int),
			camel:       utils.NewCamelName(tableName),
		}

		// Columns info
		cols, err := loader.LoadColumns(tableName)
		if err != nil {
			return nil, err
		}

		for _, col := range cols {
			column := &ColumnInfo{
				table:  table,
				column: col,
				camel:  utils.NewCamelName(col.Name),
			}
			table.columns = append(table.columns, column)
			table.columnNames[col.Name] = len(table.columns) - 1
		}

		// Auto increment column
		autoIncColumnName, err := loader.LoadAutoIncColumn(tableName)
		if err != nil {
			return nil, err
		}
		if autoIncColumnName != "" {
			table.autoIncColumn = table.columns[table.columnNames[autoIncColumnName]]
		}

		// Index info
		indexNames, err := loader.LoadIndexNames(tableName)
		if err != nil {
			return nil, err
		}

		for _, indexName := range indexNames {
			columnNames, isPrimary, isUnique, err := loader.LoadIndex(tableName, indexName)
			if err != nil {
				return nil, err
			}

			index := &IndexInfo{
				table:     table,
				indexName: indexName,
				isPrimary: isPrimary,
				isUnique:  isUnique,
				camel:     utils.NewCamelName(indexName),
			}

			for _, columnName := range columnNames {
				index.columns = append(index.columns, table.columns[table.columnNames[columnName]])
			}

			table.indices = append(table.indices, index)
			table.indexNames[indexName] = len(table.indices) - 1

			// This is primary index
			if isPrimary {
				table.primary = index
			}
		}

		// FK info
		fkNames, err := loader.LoadFKNames(tableName)
		if err != nil {
			return nil, err
		}

		for _, fkName := range fkNames {
			columnNames, refTableName, refColumnNames, err := loader.LoadFK(tableName, fkName)
			if err != nil {
				return nil, err
			}

			fk := &FKInfo{
				fkName:         fkName,
				table:          table,
				refTableName:   refTableName,
				refColumnNames: refColumnNames,
				camel:          utils.NewCamelName(fkName),
			}

			for _, columnName := range columnNames {
				fk.columns = append(fk.columns, table.columns[table.columnNames[columnName]])
			}

			table.fks = append(table.fks, fk)
			table.fkNames[fkName] = len(table.fks) - 1

		}

		db.tables = append(db.tables, table)
		db.tableNames[tableName] = len(db.tables) - 1

	}

	return db, nil

}

// Valid returns true if info != nil.
func (info *DBInfo) Valid() bool {
	return info != nil
}

// NumTable returns the number of table in the database. It returns 0 if info is nil.
func (info *DBInfo) NumTable() int {
	if info == nil {
		return 0
	}
	return len(info.tables)
}

// Table returns the i-th table in the database. It returns nil if info is nil or i is out of range.
func (info *DBInfo) Table(i int) *TableInfo {
	if info == nil {
		return nil
	}
	if i < 0 || i >= len(info.tables) {
		return nil
	}
	return info.tables[i]
}

// Tables returns all tables in the database. It returns nil if info is nil.
func (info *DBInfo) Tables() []*TableInfo {
	if info == nil {
		return nil
	}
	return info.tables
}

// TableByName returns the named table in the database. It returns nil if info is nil or table not found.
func (info *DBInfo) TableByName(tableName string) *TableInfo {
	if info == nil {
		return nil
	}
	i, found := info.tableNames[tableName]
	if !found {
		return nil
	}
	return info.tables[i]
}

// Valid returns true if info != nil.
func (info *TableInfo) Valid() bool {
	return info != nil
}

// UName is the upper camel case of the table name.
func (info *TableInfo) UName() string {
	if info == nil {
		return ""
	}
	return info.camel.UName
}

// LName is the lower camel case of the table name.
func (info *TableInfo) LName() string {
	if info == nil {
		return ""
	}
	return info.camel.LName
}

// TableName returns the table name or "" if info is nil.
func (info *TableInfo) TableName() string {
	if info == nil {
		return ""
	}
	return info.tableName
}

// NumColumn returns the number of columns in the table or 0 if info is nil.
func (info *TableInfo) NumColumn() int {
	if info == nil {
		return 0
	}
	return len(info.columns)
}

// Column returns the i-th column of the table. It returns nil if info is nil or i is out of range.
func (info *TableInfo) Column(i int) *ColumnInfo {
	if info == nil {
		return nil
	}
	if i < 0 || i >= len(info.columns) {
		return nil
	}
	return info.columns[i]
}

// Columns returns all columns in the table or nil if info is nil.
func (info *TableInfo) Columns() []*ColumnInfo {
	if info == nil {
		return nil
	}
	return info.columns
}

// ColumnByName returns the named column. It returns nil if info is nil or not found.
func (info *TableInfo) ColumnByName(columnName string) *ColumnInfo {
	if info == nil {
		return nil
	}
	i, found := info.columnNames[columnName]
	if !found {
		return nil
	}
	return info.columns[i]
}

// NumIndex returns the number of indices in the table. It returns 0 if info is nil.
func (info *TableInfo) NumIndex() int {
	if info == nil {
		return 0
	}
	return len(info.indices)
}

// Index returns the i-th index in the table. It returns nil if info is nil.
func (info *TableInfo) Index(i int) *IndexInfo {
	if info == nil {
		return nil
	}
	return info.indices[i]
}

// Indices returns all indices in the table. It returns nil if info is nil.
func (info *TableInfo) Indices() []*IndexInfo {
	if info == nil {
		return nil
	}
	return info.indices
}

// IndexByName return the named index in the table. It returns nil if info is nil or not found.
func (info *TableInfo) IndexByName(indexName string) *IndexInfo {
	if info == nil {
		return nil
	}
	i, found := info.indexNames[indexName]
	if !found {
		return nil
	}
	return info.indices[i]
}

// NumFK returns the number of foreign key in the table. It returns 0 if info is nil.
func (info *TableInfo) NumFK() int {
	if info == nil {
		return 0
	}
	return len(info.fks)
}

// FK returns the i-th foreign key in the table. It returns nil if info is nil.
func (info *TableInfo) FK(i int) *FKInfo {
	if info == nil {
		return nil
	}
	return info.fks[i]
}

// FKs returns all foreign keys in the table. It returns nil if info is nil.
func (info *TableInfo) FKs() []*FKInfo {
	if info == nil {
		return nil
	}
	return info.fks
}

// FKByName returns the named foreign key. It returns nil if info is nil or not found.
func (info *TableInfo) FKByName(fkName string) *FKInfo {
	if info == nil {
		return nil
	}
	i, found := info.fkNames[fkName]
	if !found {
		return nil
	}
	return info.fks[i]
}

// Primary returns the primary key of the table. It returns nil if info is nil or primary key not exists.
func (info *TableInfo) Primary() *IndexInfo {
	if info == nil {
		return nil
	}
	return info.primary
}

// AutoIncColumn returns the single 'auto increment' column of the table. It returns nil if info is nil or auto increment column not exists.
func (info *TableInfo) AutoIncColumn() *ColumnInfo {
	if info == nil {
		return nil
	}
	return info.autoIncColumn
}

// Valid returns true if info != nil.
func (info *ColumnInfo) Valid() bool {
	return info != nil
}

// UName is the upper camel case of the column name.
func (info *ColumnInfo) UName() string {
	if info == nil {
		return ""
	}
	return info.camel.UName
}

// LName is the lower camel case of the column name.
func (info *ColumnInfo) LName() string {
	if info == nil {
		return ""
	}
	return info.camel.LName
}

// Table returns the tabe. It returns nil if info is nil.
func (info *ColumnInfo) Table() *TableInfo {
	if info == nil {
		return nil
	}
	return info.table
}

// ColumnName returns the table column name. It returns "" if info is nil.
func (info *ColumnInfo) ColumnName() string {
	if info == nil {
		return ""
	}
	return info.column.Name
}

// DataType returns the data type of the table column. It returns "" if info is nil.
func (info *ColumnInfo) DataType() string {
	if info == nil {
		return ""
	}
	return info.column.DataType
}

// Nullable returns the nullability of the table column. It returns true if info is nil.
func (info *ColumnInfo) Nullable() bool {
	if info == nil {
		return true
	}
	return info.column.Nullable
}

// Pos returns the position of the column in table. It returns -1 if info is nil.
func (info *ColumnInfo) Pos() int {
	if info == nil {
		return -1
	}
	return info.column.Pos
}

// HasDefaultValue returns true if the table column has default value (including 'AUTO_INCREMENT'/'NOW()'). It returns false if info is nil.
func (info *ColumnInfo) HasDefaultValue() bool {
	if info == nil {
		return false
	}
	return info.column.HasDefaultValue
}

// Col returns the underly datasrc.Column. It returns nil if info is nil.
func (info *ColumnInfo) Col() *datasrc.Column {
	if info == nil {
		return nil
	}
	return &info.column
}

// Valid returns true if info != nil.
func (info *IndexInfo) Valid() bool {
	return info != nil
}

// UName is upper camel case of the index name.
func (info *IndexInfo) UName() string {
	if info == nil {
		return ""
	}
	return info.camel.UName
}

// LName is upper camel case of the index name.
func (info *IndexInfo) LName() string {
	if info == nil {
		return ""
	}
	return info.camel.LName
}

// IndexName returns the name of the index. It returns "" if info is nil.
func (info *IndexInfo) IndexName() string {
	if info == nil {
		return ""
	}
	return info.indexName
}

// Table returns the table. It returns nil if info is nil.
func (info *IndexInfo) Table() *TableInfo {
	if info == nil {
		return nil
	}
	return info.table
}

// Columns returns the composed columns. It returns nil if info is nil.
func (info *IndexInfo) Columns() []*ColumnInfo {
	if info == nil {
		return nil
	}
	return info.columns
}

// IsPrimary returns true if this is a valid primary index.
func (info *IndexInfo) IsPrimary() bool {
	if info == nil {
		return false
	}
	return info.isPrimary
}

// IsUnique returns true if this is a valid unique index.
func (info *IndexInfo) IsUnique() bool {
	if info == nil {
		return false
	}
	return info.isUnique
}

// Valid returns true if info != nil.
func (info *FKInfo) Valid() bool {
	return info != nil
}

// UName is upper camel case of the fk name.
func (info *FKInfo) UName() string {
	if info == nil {
		return ""
	}
	return info.camel.UName
}

// LName is lower camel case of the fk name.
func (info *FKInfo) LName() string {
	if info == nil {
		return ""
	}
	return info.camel.LName
}

// FKName returns the name of foreign key. It returns "" if info is nil.
func (info *FKInfo) FKName() string {
	if info == nil {
		return ""
	}
	return info.fkName
}

// Table returns the table. It returns nil if info is nil.
func (info *FKInfo) Table() *TableInfo {
	if info == nil {
		return nil
	}
	return info.table
}

// Columns returns the composed columns. It returns nil if info is nil.
func (info *FKInfo) Columns() []*ColumnInfo {
	if info == nil {
		return nil
	}
	return info.columns
}

// RefTable returns the referenced table. It returns nil if info is nil or ref table not found in current database.
func (info *FKInfo) RefTable() *TableInfo {
	if info == nil {
		return nil
	}
	return info.table.db.TableByName(info.refTableName)
}

// RefColumns returns the referenced columns. It returns nil if info is nil or ref table not found in current database.
func (info *FKInfo) RefColumns() []*ColumnInfo {
	if info == nil {
		return nil
	}
	refTable := info.RefTable()
	if refTable == nil {
		return nil
	}
	refColumns := []*ColumnInfo{}
	for _, refColumnName := range info.refColumnNames {
		refColumn := refTable.ColumnByName(refColumnName)
		if refColumn == nil {
			panic(fmt.Errorf("Can't find column %+q in ref table %+q", refColumnName, info.refTableName))
		}
		refColumns = append(refColumns, refColumn)
	}
	return refColumns
}
