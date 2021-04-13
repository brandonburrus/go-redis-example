package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var port = getEnv("PORT", "8080")
var redisHost = getEnv("REDIS_HOST", "127.0.0.1")
var redisPort = getEnv("REDIS_PORT", "6379")

var rdb = getRedisClient()

func getEnv(env, defaultValue string) string {
	envValue, isEnvValuePresent := os.LookupEnv(env)
	if isEnvValuePresent {
		return envValue
	}
	return defaultValue
}

func getRedisClient() *redis.Client {
	redisAddr := redisHost + ":" + redisPort
	log.Printf("event=redis_connect redis_addr=%s", redisAddr)
	r := redis.NewClient(&redis.Options{
		Addr:	  redisAddr,
		Password: "",
		DB: 	  0,
	})
	if r == nil {
		log.Printf("event=redis_fail")
	}
	return r
}

func helloWorld(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello, world!")
}

type Count struct {
	Count string
}

func getCounter(ctx *gin.Context) {
	val, err := rdb.Get("count").Result()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "")
		return
	} 
	if val == "" {
		rdb.Set("count", "0", -1)
		ctx.String(http.StatusOK, "Count: <empty>")
		return
	} 
	ctx.HTML(http.StatusOK, "count.tmpl", Count{
		Count: val,
	})
}

func increaseCounter(ctx *gin.Context) {
	_, err := rdb.Incr("count").Result()
	if err == nil {
		getCounter(ctx)
		return
	}
	ctx.JSON(http.StatusInternalServerError, struct{
		status int
	}{
		status: http.StatusInternalServerError,
	})
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", helloWorld)
	router.GET("/count", getCounter)
	router.POST("/count", increaseCounter)

	router.Run(":" + port)
}