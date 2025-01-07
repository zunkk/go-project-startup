// Code generated by SQLBoiler 4.18.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// /////////////////////////////// BEGIN EXTENSIONS /////////////////////////////////
// Expose table columns
var (
	UserAllColumns            = userAllColumns
	UserColumnsWithoutDefault = userColumnsWithoutDefault
	UserColumnsWithDefault    = userColumnsWithDefault
	UserPrimaryKeyColumns     = userPrimaryKeyColumns
	UserGeneratedColumns      = userGeneratedColumns
)

// GetID get ID from model object
func (o *User) GetID() int64 {
	return o.ID
}

// GetIDs extract IDs from model objects
func (s UserSlice) GetIDs() []int64 {
	result := make([]int64, len(s))
	for i := range s {
		result[i] = s[i].ID
	}
	return result
}

// GetIntfIDs extract IDs from model objects as interface slice
func (s UserSlice) GetIntfIDs() []interface{} {
	result := make([]interface{}, len(s))
	for i := range s {
		result[i] = s[i].ID
	}
	return result
}

// ToIDMap convert a slice of model objects to a map with ID as key
func (s UserSlice) ToIDMap() map[int64]*User {
	result := make(map[int64]*User, len(s))
	for _, o := range s {
		result[o.ID] = o
	}
	return result
}

// ToUniqueItems construct a slice of unique items from the given slice
func (s UserSlice) ToUniqueItems() UserSlice {
	result := make(UserSlice, 0, len(s))
	mapChk := make(map[int64]struct{}, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		o := s[i]
		if _, ok := mapChk[o.ID]; !ok {
			mapChk[o.ID] = struct{}{}
			result = append(result, o)
		}
	}
	return result
}

// FindItemByID find item by ID in the slice
func (s UserSlice) FindItemByID(id int64) *User {
	for _, o := range s {
		if o.ID == id {
			return o
		}
	}
	return nil
}

// FindMissingItemIDs find all item IDs that are not in the list
// NOTE: the input ID slice should contain unique values
func (s UserSlice) FindMissingItemIDs(expectedIDs []int64) []int64 {
	if len(s) == 0 {
		return expectedIDs
	}
	result := []int64{}
	mapChk := s.ToIDMap()
	for _, id := range expectedIDs {
		if _, ok := mapChk[id]; !ok {
			result = append(result, id)
		}
	}
	return result
}

// InsertAll inserts all rows with the specified column values, using an executor.
// IMPORTANT: this will calculate the widest columns from all items in the slice, be careful if you want to use default column values
func (o UserSlice) InsertAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	// Calculate the widest columns from all rows need to insert
	wlCols := make(map[string]struct{}, 10)
	for _, row := range o {
		wl, _ := columns.InsertColumnSet(
			userAllColumns,
			userColumnsWithDefault,
			userColumnsWithoutDefault,
			queries.NonZeroDefaultSet(userColumnsWithDefault, row),
		)
		for _, col := range wl {
			wlCols[col] = struct{}{}
		}
	}
	wl := make([]string, 0, len(wlCols))
	for _, col := range userAllColumns {
		if _, ok := wlCols[col]; ok {
			wl = append(wl, col)
		}
	}

	var sql string
	vals := []interface{}{}
	for i, row := range o {

		if i == 0 {
			sql = "INSERT INTO \"user\" " + "(\"" + strings.Join(wl, "\",\"") + "\")" + " VALUES "
		}
		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}
		valMapping, err := queries.BindMapping(userType, userMapping, wl)
		if err != nil {
			return 0, err
		}

		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, vals)
	}

	result, err := exec.ExecContext(ctx, sql, vals...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to insert all from user slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by insertall for user")
	}

	return rowsAff, nil
}

// InsertIgnoreAll inserts all rows with ignoring the existing ones having the same primary key values.
// NOTE: This function calls UpsertAll() with updateOnConflict=false and conflictColumns=<primary key columns>
// IMPORTANT: this will calculate the widest columns from all items in the slice, be careful if you want to use default column values
// IMPORTANT: if the table has `id` column of auto-increment type, this may not work as expected
func (o UserSlice) InsertIgnoreAll(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	return o.UpsertAll(ctx, exec, false, userPrimaryKeyColumns, boil.None(), columns)
}

// UpsertAll inserts or updates all rows
// Currently it doesn't support "NoContext" and "NoRowsAffected"
// IMPORTANT: this will calculate the widest columns from all items in the slice, be careful if you want to use default column values
// IMPORTANT: if the table has `id` column of auto-increment type, this may not work as expected
func (o UserSlice) UpsertAll(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	// Calculate the widest columns from all rows need to upsert
	insertCols := make(map[string]struct{}, 10)
	for _, row := range o {
		insert, _ := insertColumns.InsertColumnSet(
			userAllColumns,
			userColumnsWithDefault,
			userColumnsWithoutDefault,
			queries.NonZeroDefaultSet(userColumnsWithDefault, row),
		)
		for _, col := range insert {
			insertCols[col] = struct{}{}
		}
	}
	insert := make([]string, 0, len(insertCols))
	for _, col := range userAllColumns {
		if _, ok := insertCols[col]; ok {
			insert = append(insert, col)
		}
	}

	update := updateColumns.UpdateColumnSet(
		userAllColumns,
		userPrimaryKeyColumns,
	)

	if updateOnConflict && len(update) == 0 {
		return 0, errors.New("model: unable to upsert user, could not build update column list")
	}

	conflict := conflictColumns
	if len(conflict) == 0 {
		conflict = make([]string, len(userPrimaryKeyColumns))
		copy(conflict, userPrimaryKeyColumns)
	}

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	columns := "DEFAULT VALUES"
	if len(insert) != 0 {
		columns = fmt.Sprintf("(%s) VALUES %s",
			strings.Join(insert, ", "),
			strmangle.Placeholders(dialect.UseIndexPlaceholders, len(insert)*len(o), 1, len(insert)),
		)
	}

	fmt.Fprintf(
		buf,
		"INSERT INTO %s %s ON CONFLICT ",
		"\"user\"",
		columns,
	)

	if !updateOnConflict || len(update) == 0 {
		buf.WriteString("DO NOTHING")
	} else {
		buf.WriteByte('(')
		buf.WriteString(strings.Join(conflict, ", "))
		buf.WriteString(") DO UPDATE SET ")

		for i, v := range update {
			if i != 0 {
				buf.WriteByte(',')
			}
			quoted := strmangle.IdentQuote(dialect.LQ, dialect.RQ, v)
			buf.WriteString(quoted)
			buf.WriteString(" = EXCLUDED.")
			buf.WriteString(quoted)
		}
	}

	query := buf.String()
	valueMapping, err := queries.BindMapping(userType, userMapping, insert)
	if err != nil {
		return 0, err
	}

	var vals []interface{}
	for _, row := range o {

		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valueMapping)...)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, vals)
	}

	result, err := exec.ExecContext(ctx, query, vals...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to upsert for user")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by upsert for user")
	}

	return rowsAff, nil
}

// DeleteAllByPage delete all User records from the slice.
// This function deletes data by pages to avoid exceeding Postgres limitation (max parameters: 65535)
func (s UserSlice) DeleteAllByPage(ctx context.Context, exec boil.ContextExecutor, limits ...int) (int64, error) {
	length := len(s)
	if length == 0 {
		return 0, nil
	}

	// max number of parameters = 65535
	chunkSize := DefaultPageSize
	if len(limits) > 0 && limits[0] > 0 && limits[0] <= MaxPageSize {
		chunkSize = limits[0]
	}
	if length <= chunkSize {
		return s.DeleteAll(ctx, exec)
	}

	rowsAffected := int64(0)
	start := 0
	for {
		end := start + chunkSize
		if end > length {
			end = length
		}
		rows, err := s[start:end].DeleteAll(ctx, exec)
		if err != nil {
			return rowsAffected, err
		}

		rowsAffected += rows
		start = end
		if start >= length {
			break
		}
	}
	return rowsAffected, nil
}

// UpdateAllByPage update all User records from the slice.
// This function updates data by pages to avoid exceeding Postgres limitation (max parameters: 65535)
func (s UserSlice) UpdateAllByPage(ctx context.Context, exec boil.ContextExecutor, cols M, limits ...int) (int64, error) {
	length := len(s)
	if length == 0 {
		return 0, nil
	}

	// max number of parameters = 65535
	// NOTE: len(cols) should not be too big
	chunkSize := DefaultPageSize
	if len(limits) > 0 && limits[0] > 0 && limits[0] <= MaxPageSize {
		chunkSize = limits[0]
	}
	if length <= chunkSize {
		return s.UpdateAll(ctx, exec, cols)
	}

	rowsAffected := int64(0)
	start := 0
	for {
		end := start + chunkSize
		if end > length {
			end = length
		}
		rows, err := s[start:end].UpdateAll(ctx, exec, cols)
		if err != nil {
			return rowsAffected, err
		}

		rowsAffected += rows
		start = end
		if start >= length {
			break
		}
	}
	return rowsAffected, nil
}

// InsertAllByPage insert all User records from the slice.
// This function inserts data by pages to avoid exceeding Postgres limitation (max parameters: 65535)
func (s UserSlice) InsertAllByPage(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns, limits ...int) (int64, error) {
	length := len(s)
	if length == 0 {
		return 0, nil
	}

	// max number of parameters = 65535
	chunkSize := MaxPageSize / reflect.ValueOf(&UserColumns).Elem().NumField()
	if len(limits) > 0 && limits[0] > 0 && limits[0] < chunkSize {
		chunkSize = limits[0]
	}
	if length <= chunkSize {
		return s.InsertAll(ctx, exec, columns)
	}

	rowsAffected := int64(0)
	start := 0
	for {
		end := start + chunkSize
		if end > length {
			end = length
		}
		rows, err := s[start:end].InsertAll(ctx, exec, columns)
		if err != nil {
			return rowsAffected, err
		}

		rowsAffected += rows
		start = end
		if start >= length {
			break
		}
	}
	return rowsAffected, nil
}

// InsertIgnoreAllByPage insert all User records from the slice.
// This function inserts data by pages to avoid exceeding Postgres limitation (max parameters: 65535)
func (s UserSlice) InsertIgnoreAllByPage(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns, limits ...int) (int64, error) {
	length := len(s)
	if length == 0 {
		return 0, nil
	}

	// max number of parameters = 65535
	chunkSize := MaxPageSize / reflect.ValueOf(&UserColumns).Elem().NumField()
	if len(limits) > 0 && limits[0] > 0 && limits[0] < chunkSize {
		chunkSize = limits[0]
	}
	if length <= chunkSize {
		return s.InsertIgnoreAll(ctx, exec, columns)
	}

	rowsAffected := int64(0)
	start := 0
	for {
		end := start + chunkSize
		if end > length {
			end = length
		}
		rows, err := s[start:end].InsertIgnoreAll(ctx, exec, columns)
		if err != nil {
			return rowsAffected, err
		}

		rowsAffected += rows
		start = end
		if start >= length {
			break
		}
	}
	return rowsAffected, nil
}

// UpsertAllByPage upsert all User records from the slice.
// This function upserts data by pages to avoid exceeding Postgres limitation (max parameters: 65535)
func (s UserSlice) UpsertAllByPage(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, limits ...int) (int64, error) {
	length := len(s)
	if length == 0 {
		return 0, nil
	}

	// max number of parameters = 65535
	chunkSize := MaxPageSize / reflect.ValueOf(&UserColumns).Elem().NumField()
	if len(limits) > 0 && limits[0] > 0 && limits[0] < chunkSize {
		chunkSize = limits[0]
	}
	if length <= chunkSize {
		return s.UpsertAll(ctx, exec, updateOnConflict, conflictColumns, updateColumns, insertColumns)
	}

	rowsAffected := int64(0)
	start := 0
	for {
		end := start + chunkSize
		if end > length {
			end = length
		}
		rows, err := s[start:end].UpsertAll(ctx, exec, updateOnConflict, conflictColumns, updateColumns, insertColumns)
		if err != nil {
			return rowsAffected, err
		}

		rowsAffected += rows
		start = end
		if start >= length {
			break
		}
	}
	return rowsAffected, nil
}

///////////////////////////////// END EXTENSIONS /////////////////////////////////

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *User) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("model: no user provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(userColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	userUpsertCacheMut.RLock()
	cache, cached := userUpsertCache[key]
	userUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			userAllColumns,
			userColumnsWithDefault,
			userColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			userAllColumns,
			userPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("model: unable to upsert user, could not build update column list")
		}

		ret := strmangle.SetComplement(userAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(userPrimaryKeyColumns) == 0 {
				return errors.New("model: unable to upsert user, could not build conflict column list")
			}

			conflict = make([]string, len(userPrimaryKeyColumns))
			copy(conflict, userPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"user\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(userType, userMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(userType, userMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "model: unable to upsert user")
	}

	if !cached {
		userUpsertCacheMut.Lock()
		userUpsertCache[key] = cache
		userUpsertCacheMut.Unlock()
	}

	return nil
}

// User is an object representing the database table.
type User struct {
	ID                           int64     `csv:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	CreateTime                   time.Time `csv:"create_time" boil:"create_time" json:"create_time" toml:"create_time" yaml:"create_time"`
	UpdateTime                   time.Time `csv:"update_time" boil:"update_time" json:"update_time" toml:"update_time" yaml:"update_time"`
	DeleteTime                   time.Time `csv:"delete_time" boil:"delete_time" json:"delete_time" toml:"delete_time" yaml:"delete_time"`
	DelState                     int64     `csv:"del_state" boil:"del_state" json:"del_state" toml:"del_state" yaml:"del_state"`
	Version                      int64     `csv:"version" boil:"version" json:"version" toml:"version" yaml:"version"`
	Username                     string    `csv:"username" boil:"username" json:"username" toml:"username" yaml:"username"`
	Info                         string    `csv:"info" boil:"info" json:"info" toml:"info" yaml:"info"`
	Role                         string    `csv:"role" boil:"role" json:"role" toml:"role" yaml:"role"`
	LastUsedWalletSendTXConfigID int64     `csv:"last_used_wallet_send_tx_config_id" boil:"last_used_wallet_send_tx_config_id" json:"last_used_wallet_send_tx_config_id" toml:"last_used_wallet_send_tx_config_id" yaml:"last_used_wallet_send_tx_config_id"`

	R *userR `csv:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
	L userL  `csv:"-" boil:"-" json:"-" toml:"-" yaml:"-"`
}

var UserColumns = struct {
	ID                           string
	CreateTime                   string
	UpdateTime                   string
	DeleteTime                   string
	DelState                     string
	Version                      string
	Username                     string
	Info                         string
	Role                         string
	LastUsedWalletSendTXConfigID string
}{
	ID:                           "id",
	CreateTime:                   "create_time",
	UpdateTime:                   "update_time",
	DeleteTime:                   "delete_time",
	DelState:                     "del_state",
	Version:                      "version",
	Username:                     "username",
	Info:                         "info",
	Role:                         "role",
	LastUsedWalletSendTXConfigID: "last_used_wallet_send_tx_config_id",
}

var UserTableColumns = struct {
	ID                           string
	CreateTime                   string
	UpdateTime                   string
	DeleteTime                   string
	DelState                     string
	Version                      string
	Username                     string
	Info                         string
	Role                         string
	LastUsedWalletSendTXConfigID string
}{
	ID:                           "user.id",
	CreateTime:                   "user.create_time",
	UpdateTime:                   "user.update_time",
	DeleteTime:                   "user.delete_time",
	DelState:                     "user.del_state",
	Version:                      "user.version",
	Username:                     "user.username",
	Info:                         "user.info",
	Role:                         "user.role",
	LastUsedWalletSendTXConfigID: "user.last_used_wallet_send_tx_config_id",
}

// Generated where

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod     { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod    { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) LIKE(x string) qm.QueryMod   { return qm.Where(w.field+" LIKE ?", x) }
func (w whereHelperstring) NLIKE(x string) qm.QueryMod  { return qm.Where(w.field+" NOT LIKE ?", x) }
func (w whereHelperstring) ILIKE(x string) qm.QueryMod  { return qm.Where(w.field+" ILIKE ?", x) }
func (w whereHelperstring) NILIKE(x string) qm.QueryMod { return qm.Where(w.field+" NOT ILIKE ?", x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

var UserWhere = struct {
	ID                           whereHelperint64
	CreateTime                   whereHelpertime_Time
	UpdateTime                   whereHelpertime_Time
	DeleteTime                   whereHelpertime_Time
	DelState                     whereHelperint64
	Version                      whereHelperint64
	Username                     whereHelperstring
	Info                         whereHelperstring
	Role                         whereHelperstring
	LastUsedWalletSendTXConfigID whereHelperint64
}{
	ID:                           whereHelperint64{field: "\"user\".\"id\""},
	CreateTime:                   whereHelpertime_Time{field: "\"user\".\"create_time\""},
	UpdateTime:                   whereHelpertime_Time{field: "\"user\".\"update_time\""},
	DeleteTime:                   whereHelpertime_Time{field: "\"user\".\"delete_time\""},
	DelState:                     whereHelperint64{field: "\"user\".\"del_state\""},
	Version:                      whereHelperint64{field: "\"user\".\"version\""},
	Username:                     whereHelperstring{field: "\"user\".\"username\""},
	Info:                         whereHelperstring{field: "\"user\".\"info\""},
	Role:                         whereHelperstring{field: "\"user\".\"role\""},
	LastUsedWalletSendTXConfigID: whereHelperint64{field: "\"user\".\"last_used_wallet_send_tx_config_id\""},
}

// UserRels is where relationship names are stored.
var UserRels = struct {
}{}

// userR is where relationships are stored.
type userR struct {
}

// NewStruct creates a new relationship struct
func (*userR) NewStruct() *userR {
	return &userR{}
}

// userL is where Load methods for each relationship are stored.
type userL struct{}

var (
	userAllColumns            = []string{"id", "create_time", "update_time", "delete_time", "del_state", "version", "username", "info", "role", "last_used_wallet_send_tx_config_id"}
	userColumnsWithoutDefault = []string{"id", "create_time", "update_time", "delete_time"}
	userColumnsWithDefault    = []string{"del_state", "version", "username", "info", "role", "last_used_wallet_send_tx_config_id"}
	userPrimaryKeyColumns     = []string{"id"}
	userGeneratedColumns      = []string{}
)

type (
	// UserSlice is an alias for a slice of pointers to User.
	// This should almost always be used instead of []User.
	UserSlice []*User

	userQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	userType                 = reflect.TypeOf(&User{})
	userMapping              = queries.MakeStructMapping(userType)
	userPrimaryKeyMapping, _ = queries.BindMapping(userType, userMapping, userPrimaryKeyColumns)
	userInsertCacheMut       sync.RWMutex
	userInsertCache          = make(map[string]insertCache)
	userUpdateCacheMut       sync.RWMutex
	userUpdateCache          = make(map[string]updateCache)
	userUpsertCacheMut       sync.RWMutex
	userUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single user record from the query.
func (q userQuery) One(ctx context.Context, exec boil.ContextExecutor) (*User, error) {
	o := &User{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: failed to execute a one query for user")
	}

	return o, nil
}

// All returns all User records from the query.
func (q userQuery) All(ctx context.Context, exec boil.ContextExecutor) (UserSlice, error) {
	var o []*User

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "model: failed to assign all query results to User slice")
	}

	return o, nil
}

// Count returns the count of all User records in the query.
func (q userQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to count user rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q userQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "model: failed to check if user exists")
	}

	return count > 0, nil
}

// Users retrieves all the records using an executor.
func Users(mods ...qm.QueryMod) userQuery {
	mods = append(mods, qm.From("\"user\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"user\".*"})
	}

	return userQuery{q}
}

// FindUser retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindUser(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*User, error) {
	userObj := &User{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"user\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, userObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "model: unable to select from user")
	}

	return userObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *User) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("model: no user provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(userColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	userInsertCacheMut.RLock()
	cache, cached := userInsertCache[key]
	userInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			userAllColumns,
			userColumnsWithDefault,
			userColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(userType, userMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(userType, userMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"user\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"user\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "model: unable to insert into user")
	}

	if !cached {
		userInsertCacheMut.Lock()
		userInsertCache[key] = cache
		userInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the User.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *User) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	userUpdateCacheMut.RLock()
	cache, cached := userUpdateCache[key]
	userUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			userAllColumns,
			userPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("model: unable to update user, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"user\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, userPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(userType, userMapping, append(wl, userPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update user row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by update for user")
	}

	if !cached {
		userUpdateCacheMut.Lock()
		userUpdateCache[key] = cache
		userUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q userQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all for user")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected for user")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o UserSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("model: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"user\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, userPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to update all in user slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to retrieve rows affected all in update all user")
	}
	return rowsAff, nil
}

// Delete deletes a single User record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *User) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("model: no User provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), userPrimaryKeyMapping)
	sql := "DELETE FROM \"user\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete from user")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by delete for user")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q userQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("model: no userQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from user")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for user")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o UserSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"user\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "model: unable to delete all from user slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "model: failed to get rows affected by deleteall for user")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *User) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindUser(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *UserSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := UserSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), userPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"user\".* FROM \"user\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, userPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "model: unable to reload all in UserSlice")
	}

	*o = slice

	return nil
}

// UserExists checks if the User row exists.
func UserExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"user\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "model: unable to check if user exists")
	}

	return exists, nil
}

// Exists checks if the User row exists.
func (o *User) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return UserExists(ctx, exec, o.ID)
}
