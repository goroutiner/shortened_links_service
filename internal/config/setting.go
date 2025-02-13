package config

import (
	"os"
	"time"

	"golang.org/x/time/rate"
)

var (
	Port         = os.Getenv("PORT")
	Mode         = os.Getenv("MODE")
	PsqlUrl      = os.Getenv("DATABASE_URL")
	CharList     = "_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	LenShortLink = 10

	RateLimit   = rate.Limit(20) // ограничение RPS для пользователя
	BufferLimit = 40             // емкость "ведра" запросов, которые могут обрабатываться поверх RPS ограничения за раз

	CleanupInterval = 1 * time.Minute // интервал для чистки словаря с лимитерами неактивных пользователей
	InactivityLimit = 5 * time.Minute // время, через которое пользователь становится неактивным
)
