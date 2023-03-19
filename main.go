package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func CookieTool() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get cookie
		if cookie, err := c.Cookie("label"); err == nil {
			if cookie == "ok" {
				c.Next()
				return
			}
		}

		// Cookie verification failed
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden with no cookie"})
		c.Abort()
	}
}

func main() {
	// Redis client連線設定
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redis的伺服器地址
        Password: "",              // 如果設置了密碼，則需要提供密碼
        DB:       0,               // 使用的數據庫編號
    })

	route := gin.Default()

	route.GET("/login", func(c *gin.Context) {
		// Set cookie {"label": "ok" }, maxAge 30 seconds.
		c.SetCookie("label", "ok", 30, "/", "localhost", false, true)
		c.String(200, "Login success!")
	})

	route.GET("/home", CookieTool(), func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "Your home page"})
	})

	route.GET("/redis", func(c *gin.Context) {
        // 對Redis數據庫進行操作
        err := rdb.Incr(c.Request.Context(), "counter").Err()
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        // 讀取Redis數據庫的值
        val, err := rdb.Get(c.Request.Context(), "counter").Result()
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        c.JSON(200, gin.H{"counter": val})
    })

	route.Run(":8080")
}
