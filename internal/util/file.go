package util

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GetCurrentAbPath 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return GetCurrentAbPathByCaller(0)
	}
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// GetCurrentAbPathByExecutable 获取当前执行文件绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// GetCurrentAbPathByCaller 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller(skip int) string {
	var abPath string
	_, filename, _, ok := runtime.Caller(skip)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func AbPath(file string) string {

	if strings.HasPrefix(file, "/") {
		return file
	}

	return fmt.Sprintf("%s/%s", GetCurrentAbPathByCaller(2), file)
}

// PathExists 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetFileSize 获取文件大小
func GetFileSize(file string) int64 {
	f, err := os.Stat(file)
	if err != nil {
		return 0
	}
	return f.Size()
}

// Write 写文件
func Write(path string, data []byte) error {
	if ok, err := PathExists(path); !ok || err != nil {
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

// ToANSI UTF8转ANSI
func ToANSI(str string) (string, error) {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest([]byte(str))
	if err != nil {
		return "", err
	}
	if result.Charset == "UTF-8" {
		return str, nil
	}
	// 尝试使用不同的编码解码文件内容
	var decodedContent string
	decodedContent, err = encodeToUTF8(str, simplifiedchinese.GB18030.NewEncoder()) // 尝试GB2312编码
	if err != nil {
		decodedContent, err = encodeToUTF8(str, charmap.Windows1252.NewEncoder()) // 尝试ISO-8859-1 (ANSI) 编码
		if err != nil {
			return "", fmt.Errorf("invalid charset: %s. err: %s", result.Charset, err.Error())
		}
	}
	return decodedContent, nil
}

// ToUtf8 ANSI or GB2312 转UTF8
func ToUtf8(str string) (string, error) {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest([]byte(str))
	if err != nil {
		return "", err
	}
	if result.Charset == "UTF-8" {
		return str, nil
	}
	// 尝试使用不同的编码解码文件内容
	var decodedContent string
	decodedContent, err = decodeToUTF8(str, simplifiedchinese.GB18030.NewDecoder()) // 尝试GB2312编码解码
	if err != nil {
		decodedContent, err = decodeToUTF8(str, charmap.Windows1252.NewDecoder()) // 尝试ISO-8859-1 (ANSI) 编码解码
		if err != nil {
			fmt.Println("无法编码文件内容:", err)
			return "", fmt.Errorf("invalid charset: %s. err: %s", result.Charset, err.Error())
		}
	}
	return decodedContent, nil
}

// decodeToUTF8 使用指定的解码器解码文件内容
func decodeToUTF8(content string, decoder transform.Transformer) (string, error) {
	// 将字符串转换为字节数组
	bytes := []byte(content)

	// 使用指定的解码器解码字节数组
	decodedBytes, _, err := transform.Bytes(decoder, bytes)
	if err != nil {
		return "", err
	}
	// 将字节数组转换为字符串
	return string(decodedBytes), nil
}

// encodeToUTF8 使用指定的编码器编码文件内容
func encodeToUTF8(content string, encoder transform.Transformer) (string, error) {
	// 将字符串转换为字节数组
	bytes := []byte(content)

	// 使用指定的编码器编码字节数组
	encodedBytes, _, err := transform.Bytes(encoder, bytes)
	if err != nil {
		return "", err
	}
	// 将字节数组转换为字符串
	return string(encodedBytes), nil
}
