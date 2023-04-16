package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lightsaid/short-net/util"
)

func (app *application) formUpload(w http.ResponseWriter, r *http.Request) (filename string, err error) {
	// 控制Body数据大小，包括文件和Form表单其他字段数据，假如想控制文件上传大小不能超过2M,
	// 需要多设置512kb或者1MB给表单其他数据
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20+512)

	// 上传的文件存储在maxMemory大小的内存里面，如果文件大小超过了maxMemory，
	// 那么剩下的部分将存储在系统的临时文件中。
	err = r.ParseMultipartForm(4 << 20)
	if err != nil {
		err = fmt.Errorf("r.ParseMultipartForm error: %w", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		// fmt.Println(">>>>> ", errors.Is(err, http.ErrMissingFile))
		err = fmt.Errorf("r.FormFile error: %w", err)
		return
	}
	defer file.Close()
	fileExt := filepath.Ext(header.Filename)
	filename = "./static/images/" + util.RandomString(8) + fileExt
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		err = fmt.Errorf("os.OpenFile error: %w", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		err = fmt.Errorf("io.Copy error: %w", err)
		return
	}

	filename = strings.TrimPrefix(filename, ".")

	return
}
