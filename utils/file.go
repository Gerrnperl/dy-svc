package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveFile 保存文件
//
// saves a file to the specified folder with the specified filename.
func SaveFile(data *multipart.FileHeader, folder string, filename string) error {
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	src, err := data.Open()

	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(folder + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

// RemoveFile 删除文件
func RemoveFile(path string) error {
	return os.Remove(path)
}

// GetExt 获取文件扩展名
func GetExt(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}
	return ext[1:]
}
