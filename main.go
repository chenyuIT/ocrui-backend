package main 

import (
	"context"
	"fmt"
        "time"	
	"upload/biz/handler"
	
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/client"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/cors"

)

func main() {
	// WithMaxRequestBodySize can set the size of the body
	h := server.Default(server.WithHostPorts("0.0.0.0:5001"), server.WithMaxRequestBodySize(20<<20))

	
	
        h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allowed domains, need to bring schema
		AllowMethods:     []string{"PUT", "PATCH","GET","POST"},    // Allowed request methods
		AllowHeaders:     []string{"Origin"},          // Allowed request headers
		ExposeHeaders:    []string{"Content-Length"},  // Request headers allowed in the upload_file
		AllowCredentials: true,                        // Whether cookies are attached
		AllowOriginFunc: func(origin string) bool { // Custom domain detection with lower priority than AllowOrigins
	        return origin == "*"
		},
		MaxAge: 12 * time.Hour, // Maximum length of upload_file-side cache preflash requests (seconds)
	}))

	h.GET("/cors", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "Cross Domain OK!")
	})


        h.GET("/ping",handler.Ping)

	h.POST("/upload", func(ctx context.Context, c *app.RequestContext) {
		// single file
		file, _ := c.FormFile("file")
		fmt.Println(file.Filename)

		// Upload the file to specific dst
		c.SaveUploadedFile(file, fmt.Sprintf("./upload/%s", file.Filename))

		c.String(consts.StatusOK, fmt.Sprintf("'%s' uploaded,next scan ocr to txt file!", file.Filename))
		
		cc, err := client.NewClient()
        	if err != nil {
			return  
	        }
		status, body, _ := cc.Get(context.Background(), nil, "http://127.0.0.1:5002/get_ocr?image="+file.Filename)
		
		fmt.Println("http://127.0.0.1:5002/get_ocr?image="+file.Filename)
	        fmt.Printf("status=%v body=%v\n", status, string(body))
		

        	})

	h.POST("/uploads", func(ctx context.Context, c *app.RequestContext) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["file"]

		for _, file := range files {
			fmt.Println(file.Filename)

			// Upload the file to specific dst.
			c.SaveUploadedFile(file, fmt.Sprintf("./upload/%s", file.Filename))
		}
		c.String(consts.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})



	// File() will return the contents of the file directly
	h.GET("/download", func(ctx context.Context, c *app.RequestContext) {
		file := c.Query("filename")
		c.File("./public/"+file)
	})

	// FileAttachment() sets the "content-disposition" header and returns the file as an "attachment".
//	h.GET("/downloads", func(ctx context.Context, c *app.RequestContext) {
		// If you use Chinese, need to encode
//		file := c.Query("filename")
//		c.FileAttachment("./public/"+file)
//	})

	h.Spin()
}
