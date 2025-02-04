package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func Connect() *sql.DB {
	dsn := "postgresql://root:WM57FrgaTpjKOm7p9RO4DG241wmQeH7H@dpg-cug7gilds78s738d2pr0-a/db1_vcec" // Читаем переменную окружения
	if dsn == "" {
		log.Fatal("DATABASE_URL не найден!")
	}
	//test
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	log.Println("DB Connect Success")
	return db
}
