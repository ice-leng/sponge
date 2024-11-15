package parser

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/huandu/xstrings"
	"time"

	_ "github.com/go-sql-driver/mysql" //nolint
)

// GetMysqlTableInfo get table info from mysql
func GetMysqlTableInfo(dsn, tableName string) (string, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", fmt.Errorf("GetMysqlTableInfo error, %v", err)
	}
	defer db.Close() //nolint

	rows, err := db.Query("SHOW CREATE TABLE `" + tableName + "`")
	if err != nil {
		return "", fmt.Errorf("query show create table error, %v", err)
	}

	defer rows.Close() //nolint
	if !rows.Next() {
		return "", fmt.Errorf("not found found table '%s'", tableName)
	}

	var table string
	var info string
	err = rows.Scan(&table, &info)
	if err != nil {
		return "", err
	}

	return info, nil
}

// GetTableInfo get table info from mysql
// Deprecated: replaced by GetMysqlTableInfo
func GetTableInfo(dsn, tableName string) (string, error) {
	return GetMysqlTableInfo(dsn, tableName)
}

func getUnallocatedMenuId(db *sql.DB, date time.Time) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM `t_menu` WHERE parent_id = ? and name = ?", 0, "未分配").Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			insertQuery := "insert into `t_menu` (`parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			result, err := db.Exec(insertQuery, 0, "未分配", "CATALOG", "/unallocated", "Layout", "", 999, 1, "system", "/unallocated", date, date)
			if err != nil {
				return 0, fmt.Errorf("getUnallocatedMenuId 插入数据失败: %v", err)
			}
			id, err = result.LastInsertId()
			if err != nil {
				return 0, fmt.Errorf("getUnallocatedMenuId 获取插入数据的ID失败: %v", err)
			}
			return id, nil
		}
	}
	return id, err
}

func getMenuId(db *sql.DB, date time.Time, tableName, tableComment string) (int64, bool, error) {
	pid, err := getUnallocatedMenuId(db, date)
	if err != nil {
		return pid, true, fmt.Errorf("AddMenu search unallocated menu error, %v", err)
	}
	var id int64
	err = db.QueryRow("SELECT id FROM `t_menu` WHERE parent_id = ? and path = ?", pid, tableName).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			insertQuery := "insert into `t_menu` (`parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `created_at`, `updated_at`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			result, err := db.Exec(insertQuery, pid, tableComment+"管理", "MENU", tableName, tableName+"/index", "", 1, 1, "el-icon-User", "", date, date)
			if err != nil {
				return 0, true, fmt.Errorf("getMenuId 插入数据失败: %v", err)
			}
			id, err = result.LastInsertId()
			if err != nil {
				return 0, true, fmt.Errorf("getMenuId 获取插入数据的ID失败: %v", err)
			}
			return id, false, nil
		}
	}
	return id, true, err
}

func AddMenu(dsn, tableName, tableComment string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("AddMenu error, %v", err)
	}
	defer db.Close() //nolint

	if err = db.Ping(); err != nil {
		return fmt.Errorf("AddMenu db ping error, %v", err)
	}

	tableName = xstrings.FirstRuneToLower(tableName)
	date := time.Now()
	menuId, exist, err := getMenuId(db, date, tableName, tableComment)
	if err != nil {
		return fmt.Errorf("AddMenu error, %v", err)
	}
	if exist {
		return nil
	}

	menus := []struct {
		parent_id  int64
		name       string
		type2      string
		path       string
		component  string
		perm       string
		sort       string
		visible    string
		icon       string
		redirect   string
		created_at time.Time
		updated_at time.Time
	}{
		{menuId, tableComment + "新增", "BUTTON", "", "", tableName + ":add", "1", "1", "", "", date, date},
		{menuId, tableComment + "编辑", "BUTTON", "", "", tableName + ":edit", "2", "1", "", "", date, date},
		{menuId, tableComment + "删除", "BUTTON", "", "", tableName + ":delete", "3", "1", "", "", date, date},
	}
	query := "insert into `t_menu` (`parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `created_at`, `updated_at`) values "
	values := ""
	args := []interface{}{}
	for i, menu := range menus {
		if i > 0 {
			values += ", "
		}
		values += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		args = append(args, menu.parent_id, menu.name, menu.type2, menu.path, menu.component, menu.perm, menu.sort, menu.visible, menu.icon, menu.redirect, menu.created_at, menu.updated_at)
	}
	query += values
	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("AddMenu error, %v", err)
	}
	return nil
}
