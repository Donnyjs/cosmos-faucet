package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var chain string
var recaptchaSecretKey string
var amountFaucet string
var amountSteak string
var key string
var pass string
var node string
var publicUrl string

type claim_struct struct {
	Address  string
	Response string
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(key, "=", value)
		return value
	} else {
		log.Fatal("Error loading environment variable: ", key)
		return ""
	}
}

func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		defer func() {
			if err := recover(); err != nil {

			}
		}()
		c.Next()
	})
	r.GET("/claim", getCoinsHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "success")
	})
	r.Run(":8000")
}

func executeCmd(command string, writes ...string) {
	cmd, wc, _ := goExecute(command)

	for _, write := range writes {
		wc.Write([]byte(write + "\n"))
	}
	cmd.Wait()
}

func goExecute(command string) (cmd *exec.Cmd, pipeIn io.WriteCloser, pipeOut io.ReadCloser) {
	cmd = getCmd(command)
	pipeIn, _ = cmd.StdinPipe()
	pipeOut, _ = cmd.StdoutPipe()
	go cmd.Start()
	time.Sleep(time.Second)
	return cmd, pipeIn, pipeOut
}

func getCmd(command string) *exec.Cmd {
	// split command into command and args
	split := strings.Split(command, " ")

	var cmd *exec.Cmd
	if len(split) == 1 {
		cmd = exec.Command(split[0])
	} else {
		cmd = exec.Command(split[0], split[1:]...)
	}

	return cmd
}

func getCoinsHandler(c *gin.Context) {
	fmt.Println("enter....")
	address := c.Query("address")
	if address == "" {
		c.JSON(500, "param err")
		return
	}
	fmt.Println("address : ", address)
	sendFaucet := fmt.Sprintf(
		"echo 'y' | spiked tx bank send spike1fstqaepz9sl0jmm95f2dstmucakzyzxxqj5kla %v 100000000000000000uspike --keyring-backend test --gas-prices 1000000000uspike --chain-id spiketestnet_9090-1",
		address)
	executeCmd(sendFaucet, pass)
	c.JSON(200, "success")
}
