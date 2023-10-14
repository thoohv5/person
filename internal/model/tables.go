package model

import (
	"embed"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/thoohv5/person/internal/util"
)

// TableNameFromFilesystem 从文件系统中获取表名
func TableNameFromFilesystem(fs http.FileSystem, dir string) (tables []string, err error) {
	// 读取文件夹下所有文件
	f, err := fs.Open(dir)
	// 如果文件夹不存在，直接返回空
	if os.IsNotExist(err) {
		return tables, nil
	}
	if err != nil {
		return tables, err
	}
	defer func() {
		if fErr := f.Close(); fErr != nil {
			if err != nil {
				err = fmt.Errorf("%w, %v", err, fErr)
			} else {
				err = fErr
			}
			return
		}
	}()

	// 文件是否能存在
	if _, err = f.Stat(); os.IsNotExist(err) {
		return tables, nil
	}
	files, err := f.Readdir(-1)
	if err != nil {
		return tables, err
	}
	for _, s := range files {
		filePath := filepath.Join(dir, s.Name())
		content, err := fs.Open(filePath)
		if err != nil {
			return tables, err
		}

		fileContent, err := ioutil.ReadAll(content)
		if err != nil {
			return tables, err
		}

		// 正则表达式用于匹配CREATE TABLE语句并提取表名
		regCreateTable := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)`)
		// 正则表达式用于匹配CREATE TABLE IF NOT EXISTS语句并提取表名
		regCreateTableIfNotExists := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+IF\s+NOT\s+EXISTS\s+(\w+)`)

		// 提取CREATE TABLE语句中的表名
		matchesCreateTable := regCreateTable.FindAllStringSubmatch(string(fileContent), -1)
		for _, match := range matchesCreateTable {
			tableName := match[1]
			if strings.Contains(strings.ToLower(tableName), "if") {
				continue
			}
			tables = append(tables, tableName)
		}
		// 提取CREATE TABLE IF NOT EXISTS语句中的表名
		matchesCreateTableIfNotExists := regCreateTableIfNotExists.FindAllStringSubmatch(string(fileContent), -1)
		for _, match := range matchesCreateTableIfNotExists {
			tableName := match[1]
			tables = append(tables, tableName)
		}
	}

	return tables, nil
}

// GetTables 获取表名
func GetTables(models []interface{}, sQLMigrations embed.FS) (tableNames []string) {
	for _, m := range models {
		Type := reflect.TypeOf(m).Elem()
		name := GetTableName(reflect.New(Type).Elem().Interface())
		tableNames = append(tableNames, name)
	}
	if tables, err := TableNameFromFilesystem(http.FS(sQLMigrations), "/"); err == nil {
		tableNames = append(tableNames, tables...)
	} else {
		fmt.Printf("get table name from filesystem err: %v\n", err)
	}
	tableNames = util.UniqStrArr(tableNames)
	return
}
