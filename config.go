package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Username          string `json:"username"`
	PasswordHash      string `json:"password_hash"`
	EncryptionKey     string `json:"encryption_key"`
	Port              string `json:"port"`
	IsDefaultPassword bool   `json:"is_default_password"`
	CaptchaMode       string `json:"captcha_mode"`
	SessionHours      int    `json:"session_hours"`
}

var cfg Config

const configFile = "config.json"

func loadConfig() {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		createDefaultConfig()
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Gagal baca config: %v", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Gagal parse config: %v", err)
	}
	if cfg.CaptchaMode == "" {
		cfg.CaptchaMode = "math"
	}
	if cfg.SessionHours <= 0 {
		cfg.SessionHours = 1
	}
	if err := initCrypto(cfg.EncryptionKey); err != nil {
		log.Fatalf("Crypto error: %v", err)
	}
}

func createDefaultConfig() {
	keyBytes := make([]byte, 32)
	rand.Read(keyBytes)

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	cfg = Config{
		Username:          "admin",
		PasswordHash:      string(hash),
		EncryptionKey:     hex.EncodeToString(keyBytes),
		Port:              "8080",
		IsDefaultPassword: true,
		CaptchaMode:       "math",
		SessionHours:      1,
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(configFile, data, 0600)

	log.Println("Config baru dibuat. Username: admin | Password: admin123")
	log.Println("PENTING: Ganti password lewat menu Pengaturan setelah login!")
}

func saveConfig() error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0600)
}
