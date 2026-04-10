package util

import (
	"log"
	"os"
)

func InitLoggers() {
	file, err := os.OpenFile("storage/logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Gagal membuka file log:", err)
	}

	log.SetOutput(file)
	log.Println("--- Server IMOCA Started ---")
}