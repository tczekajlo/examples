package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/kataras/iris"
)

func main() {
	// you can debug path with get working directory
	// s, _ := os.Getwd()
	// println(s)
	//
	iris.Templates("./_examples/file_upload_simple/*")
	// Maximum request body size.
	//
	// The server rejects requests with bodies exceeding this limit.
	//
	// By default request body size is unlimited.
	iris.SetMaxRequestBodySize(32 << 20) //32MB max upload filesize

	// Serve the form.html to the user
	iris.Get("/upload", func(ctx *iris.Context) {
		//these are optionaly you can just call RenderFile("form.html",{})
		//create the token
		now := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(now, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		//render the form with the token for any use you like
		ctx.RenderFile("form.html", token)
	})

	// Handle the post request from the form.html to the server
	iris.Post("/upload", func(ctx *iris.Context) {

		// Get the file from the request
		info, err := ctx.FormFile("uploadfile")
		file, err := info.Open()
		defer file.Close()
		fname := info.Filename

		// Create a file with the same name
		// assuming that you have a folder named 'uploads'
		out, err := os.OpenFile("./uploads/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer out.Close()

		io.Copy(out, file)

	})

	fmt.Println("Iris is listening on :8080")
	iris.Listen("8080")

}