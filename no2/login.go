package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)


type User struct {
	RealName string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ctx = context.Background()
var rdb *redis.Client

//utama
func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	app := fiber.New()

	app.Post("/login", loginHandler)

	app.Post("/register", registerHandler)

	log.Fatal(app.Listen(":3000"))
}


func sha1Hash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func registerHandler(c *fiber.Ctx) error {
	type Req struct {
		Username string `json:"username"`
		RealName string `json:"realname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "sorry invalid request"})
	}

	user := User{
		RealName: req.RealName,
		Email:    req.Email,
		Password: sha1Hash(req.Password),
	}

	data, _ := json.Marshal(user)
	key := fmt.Sprintf("login_%s", req.Username) //format key regis

	err := rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save user, try again please."})
	}

	return c.JSON(fiber.Map{"message": "user registered successfully"})
}

func loginHandler(c *fiber.Ctx) error {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request,check again the body."})
	}

	key := fmt.Sprintf("login_%s", req.Username)
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found,try register again."})
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse user data"})
	}

	if user.Password != sha1Hash(req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid password"})
	}

	return c.JSON(fiber.Map{
		"message":  "login success",
		"realname": user.RealName,
		"email":    user.Email,
	})
}
