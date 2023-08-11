package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func Download(urls string) (int64, string) {
	header, err := http.Head(urls)
	if err != nil {
		log.Fatal(err)
	}
	file_size := header.ContentLength / 1024 / 1024
	fmt.Println("文件大小：", file_size, "MB")
	if file_size >= 100 {
		fmt.Println("文件大于100MB,不允许下载")
	} else {
		fmt.Println("文件小于100MB")
		file_name := path.Base(urls)
		fmt.Println(file_name)
		resp, _ := http.Get(urls)
		defer resp.Body.Close()
		out, err := os.Create("./download_files/" + file_name)
		if err != nil {
			panic(err)
		}
		defer out.Close()
		// 然后将响应流和文件流对接起来
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			panic(err)
		}
	}
	return file_size, path.Base(urls)
}

func listFiles(directory string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path[15:])
			fmt.Println(files)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.POST("/download", func(c *gin.Context) {
		urls := c.PostForm("urls")
		_, _ = Download(urls)
		//c.JSON(200, gin.H{
		//	"message":   "ok",
		//	"file_name": file_name,
		//})
		c.Redirect(http.StatusFound, "/downloadFile")
	})
	r.GET("/downloadFile", func(c *gin.Context) {
		var fileList []string
		fileList, _ = listFiles("./download_files")
		c.HTML(http.StatusOK, "downloadfile.html", gin.H{
			"Files": fileList,
		})
	})
	r.GET("/File", func(c *gin.Context) {
		fileName := c.Query("file_name")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		c.Header("Content-Type", "application/text/plain")
		c.File("./download_files/" + fileName)
		c.JSON(200, gin.H{
			"message":   "ok",
			"file_name": fileName,
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "download.html", nil)
	})
	err := r.Run(":" + os.Args[1])
	if err != nil {
		return
	}
}
