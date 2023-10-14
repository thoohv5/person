package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// FormCsv 导入
func FormCsv(resourceFile string, csvTemplate map[string]string, result interface{}) (err error) {

	defer func() {
		if rev := recover(); nil != rev {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("FormCsv: panic recovered: %s\n%s", rev, buf)
		}
	}()

	// 获取csv名称与pg标签的对应关系
	template := make(map[string]string)
	for item, field := range csvTemplate {
		items := strings.Split(item, ".")
		template[items[1]] = field
	}

	// 类型检查
	rt := reflect.TypeOf(result)
	if rt.Kind() != reflect.Ptr && rt.Elem().Kind() != reflect.Slice {
		return errors.New("result参数必须为切片的指针")
	}

	// 读取数据
	decodeString, err := base64.StdEncoding.DecodeString(resourceFile)
	if err != nil {
		return fmt.Errorf("%s", "导入数据格式错误")
	}

	// 数据编码
	var reader *csv.Reader
	if utf8.Valid(decodeString) {
		reader = csv.NewReader(transform.NewReader(bytes.NewReader(decodeString), unicode.UTF8BOM.NewDecoder()))
	} else {
		reader = csv.NewReader(transform.NewReader(bytes.NewReader(decodeString), simplifiedchinese.GBK.NewDecoder()))
	}

	// 获取结构Name与pg标签的对应关系
	mNamePg := make(map[string]string)
	rrt := rt.Elem().Elem().Elem()
	for i := 0; i < rrt.NumField(); i++ {
		f := rrt.Field(i)
		mNamePg[strings.Split(f.Tag.Get("pg"), ",")[0]] = f.Name
	}

	var (
		header []string
	)
	re := make([]reflect.Value, 0)
	for {
		// csv安行读取
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%s", "导入数据格式错误")
		}
		// 取标题
		if len(header) == 0 {
			header = line
			continue
		}

		// 内容数据绑定
		rv := reflect.New(rrt)
		for index, item := range line {
			// 内容数据与标题不对等
			if len(header) <= index {
				continue
			}
			// 获取pg标签
			sn, ok := template[header[index]]
			if !ok {
				continue
			}
			// 获取结构Name标签
			fn, ok := mNamePg[sn]
			if !ok {
				continue
			}

			// 按照类型渲染结构体
			v := rv.Elem().FieldByName(fn)
			item = strings.TrimPrefix(item, "\t")
			switch v.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				parseInt, _ := strconv.ParseInt(item, 10, 64)
				v.SetInt(parseInt)
			case reflect.String:
				v.SetString(item)
			case reflect.Slice:
				switch v.Type() {
				case reflect.TypeOf(net.IP{}):
					v.Set(reflect.ValueOf(net.ParseIP(item)))
				default:
					return fmt.Errorf("行号：%d \n字段：%s，类型：%v，未适配", index+1, item, v.Kind())
				}
			default:
				return fmt.Errorf("行号：%d 数据：%s，类型：%s，未适配", index+1, item, v.Kind())
			}

		}
		// 数据追加
		re = append(re, rv)
	}
	reflect.ValueOf(&result).Elem().Elem().Elem().Set(reflect.Append(reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(rrt)), 0, len(re)), re...))

	return nil
}

// ToCsvFile 导出
func ToCsvFile(
	ctx context.Context,
	name string,
	source func(ctx context.Context, fields []string) (interface{}, error),
	csvTemplate map[string]string) (str string, err error) {

	defer func() {
		if rev := recover(); nil != rev {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("TOCsvFile: panic recovered: %s\n%s", rev, buf)
		}
	}()

	titles := make([]string, len(csvTemplate))
	fields := make([]string, len(csvTemplate))
	for item, field := range csvTemplate {
		items := strings.Split(item, ".")
		num, _ := strconv.Atoi(items[0])
		titles[num] = items[1]
		fields[num] = field
	}

	bytesBuffer := &bytes.Buffer{}
	// 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码
	bytesBuffer.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(bytesBuffer)
	if err := writer.Write(titles); err != nil {
		return "", err
	}

	list, err := source(ctx, fields)
	if err != nil {
		return "", err
	}

	rt := reflect.TypeOf(list)
	if rt.Kind() != reflect.Slice {
		return "", errors.New("list参数必须为切片的指针")
	}

	mNamePg := make(map[string]string)
	rrt := rt.Elem().Elem()
	for i := 0; i < rrt.NumField(); i++ {
		mNamePg[strings.Split(rrt.Field(i).Tag.Get("pg"), ",")[0]] = rrt.Field(i).Name
	}

	rv := reflect.ValueOf(list)

	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i)
		rows := make([]string, 0, len(fields))
		for _, field := range fields {
			fn, ok := mNamePg[field]
			if !ok {
				continue
			}
			row := ""
			v := item.Elem().FieldByName(fn)
			switch v.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				row = strconv.FormatInt(v.Int(), 10)
			case reflect.String:
				row = v.String()
			case reflect.Slice:
				switch v.Type() {
				case reflect.TypeOf(net.IP{}):
					result := v.MethodByName("String").Call([]reflect.Value{})
					row = result[0].String()
				default:
					return "", fmt.Errorf("导出数据，字段：%s，类型：%v，未适配", field, v.Kind())
				}
			default:
				return "", fmt.Errorf("导出数据，字段：%s，类型：%v，未适配", field, v.Kind())
			}
			rows = append(rows, fmt.Sprintf("\t%s", row))
		}
		if err = writer.Write(rows); err != nil {
			return "", err
		}
	}

	writer.Flush()

	file := fmt.Sprintf("/tmp/%s", name)
	createFile, err := os.Create(file)
	if err != nil {
		return "", err
	}

	if _, err = createFile.Write(bytesBuffer.Bytes()); err != nil {
		return "", err
	}

	return file, nil
}
