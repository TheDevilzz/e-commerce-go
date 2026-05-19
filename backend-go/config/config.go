package config // ประกาศ package config

import ( // เริ่มส่วน import package ที่ไฟล์นี้ต้องใช้
	"log"     // นำเข้า package log
	"os"      // นำเข้า package os
	"strconv" // นำเข้า package strconv
	"strings" // นำเข้า package strings

	"github.com/joho/godotenv" // นำเข้า package github.com/joho/godotenv
) // ปิด block หรือกลุ่มรายการก่อนหน้า

type Config struct { // ประกาศ struct Config
	Port             string   // ทำงานคำสั่ง Port string
	DBHost           string   // ทำงานคำสั่ง DBHost string
	DBPort           string   // ทำงานคำสั่ง DBPort string
	DBUser           string   // ทำงานคำสั่ง DBUser string
	DBPassword       string   // ทำงานคำสั่ง DBPassword string
	DBName           string   // ทำงานคำสั่ง DBName string
	JWTSecret        string   // ทำงานคำสั่ง JWTSecret string
	AllowedOrigins   []string // ทำงานคำสั่ง AllowedOrigins []string
	RefreshTokenDays int      // ทำงานคำสั่ง RefreshTokenDays int
	EnforceHTTPS     bool     // ทำงานคำสั่ง EnforceHTTPS bool
} // ปิด block การทำงานปัจจุบัน

func LoadConfig() Config { // ประกาศฟังก์ชัน LoadConfig
	if err := godotenv.Load(); err != nil { // ตรวจเงื่อนไข err := godotenv.Load(); err != nil
		log.Println("No .env file loaded; using environment variables") // ทำงานคำสั่ง log.Println("No .env file loaded; using environment variables")
	} // ปิด block การทำงานปัจจุบัน

	port := getEnv("PORT", "3001")                            // กำหนดค่า port จาก getEnv("PORT", "3001")
	jwtSecret := strings.TrimSpace(os.Getenv("TOKEN_SECRET")) // กำหนดค่า jwt Secret จาก strings.TrimSpace(os.Getenv("TOKEN_SECRET"))
	if jwtSecret == "" {                                      // ตรวจเงื่อนไข jwtSecret == ""
		log.Fatal("TOKEN_SECRET environment variable is required") // ทำงานคำสั่ง log.Fatal("TOKEN_SECRET environment variable is required")
	} // ปิด block การทำงานปัจจุบัน

	refreshTokenDays := getEnvInt("REFRESH_TOKEN_DAYS", 30) // กำหนดค่า refresh Token Days จาก getEnvInt("REFRESH_TOKEN_DAYS", 30)
	if refreshTokenDays < 7 {                               // ตรวจเงื่อนไข refreshTokenDays < 7
		refreshTokenDays = 7 // กำหนดค่า refresh Token Days จาก 7
	} // ปิด block การทำงานปัจจุบัน
	if refreshTokenDays > 30 { // ตรวจเงื่อนไข refreshTokenDays > 30
		refreshTokenDays = 30 // กำหนดค่า refresh Token Days จาก 30
	} // ปิด block การทำงานปัจจุบัน

	return Config{ // คืนค่า Config
		Port:             port,                                                                               // ระบุค่ารายการหนึ่งในชุดข้อมูล
		DBHost:           os.Getenv("DB_HOST"),                                                               // ระบุค่ารายการหนึ่งในชุดข้อมูล
		DBPort:           os.Getenv("DB_PORT"),                                                               // ระบุค่ารายการหนึ่งในชุดข้อมูล
		DBUser:           os.Getenv("DB_USER"),                                                               // ระบุค่ารายการหนึ่งในชุดข้อมูล
		DBPassword:       os.Getenv("DB_PASSWORD"),                                                           // ระบุค่ารายการหนึ่งในชุดข้อมูล
		DBName:           os.Getenv("DB_NAME"),                                                               // ระบุค่ารายการหนึ่งในชุดข้อมูล
		JWTSecret:        jwtSecret,                                                                          // ระบุค่ารายการหนึ่งในชุดข้อมูล
		AllowedOrigins:   splitCSV(getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://127.0.0.1:5173")), // ระบุค่ารายการหนึ่งในชุดข้อมูล
		RefreshTokenDays: refreshTokenDays,                                                                   // ระบุค่ารายการหนึ่งในชุดข้อมูล
		EnforceHTTPS:     getEnvBool("ENFORCE_HTTPS", true),                                                  // ระบุค่ารายการหนึ่งในชุดข้อมูล
	} // ปิด block การทำงานปัจจุบัน
} // ปิด block การทำงานปัจจุบัน

func getEnv(key string, fallback string) string { // ประกาศฟังก์ชัน getEnv
	value := strings.TrimSpace(os.Getenv(key)) // กำหนดค่า value จาก strings.TrimSpace(os.Getenv(key))
	if value == "" {                           // ตรวจเงื่อนไข value == ""
		return fallback // คืนค่า fallback
	} // ปิด block การทำงานปัจจุบัน
	return value // คืนค่า value
} // ปิด block การทำงานปัจจุบัน

func getEnvInt(key string, fallback int) int { // ประกาศฟังก์ชัน getEnvInt
	value := strings.TrimSpace(os.Getenv(key)) // กำหนดค่า value จาก strings.TrimSpace(os.Getenv(key))
	if value == "" {                           // ตรวจเงื่อนไข value == ""
		return fallback // คืนค่า fallback
	} // ปิด block การทำงานปัจจุบัน
	parsed, err := strconv.Atoi(value) // กำหนดค่า parsed, err จาก strconv.Atoi(value)
	if err != nil {                    // ตรวจเงื่อนไข err != nil
		return fallback // คืนค่า fallback
	} // ปิด block การทำงานปัจจุบัน
	return parsed // คืนค่า parsed
} // ปิด block การทำงานปัจจุบัน

func getEnvBool(key string, fallback bool) bool { // ประกาศฟังก์ชัน getEnvBool
	value := strings.TrimSpace(os.Getenv(key)) // กำหนดค่า value จาก strings.TrimSpace(os.Getenv(key))
	if value == "" {                           // ตรวจเงื่อนไข value == ""
		return fallback // คืนค่า fallback
	} // ปิด block การทำงานปัจจุบัน
	parsed, err := strconv.ParseBool(value) // กำหนดค่า parsed, err จาก strconv.ParseBool(value)
	if err != nil {                         // ตรวจเงื่อนไข err != nil
		return fallback // คืนค่า fallback
	} // ปิด block การทำงานปัจจุบัน
	return parsed // คืนค่า parsed
} // ปิด block การทำงานปัจจุบัน

func splitCSV(value string) []string { // ประกาศฟังก์ชัน splitCSV
	parts := strings.Split(value, ",")     // กำหนดค่า parts จาก strings.Split(value, ",")
	items := make([]string, 0, len(parts)) // กำหนดค่า items จาก make([]string, 0, len(parts))
	for _, part := range parts {           // วนลูปตาม _, part := range parts
		item := strings.TrimSpace(part) // กำหนดค่า item จาก strings.TrimSpace(part)
		if item != "" {                 // ตรวจเงื่อนไข item != ""
			items = append(items, item) // กำหนดค่า items จาก append(items, item)
		} // ปิด block การทำงานปัจจุบัน
	} // ปิด block การทำงานปัจจุบัน
	return items // คืนค่า items
} // ปิด block การทำงานปัจจุบัน
