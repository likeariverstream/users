package environment

import (
	"os"
)

type Env struct {
	App App
	Db  Db
}

type App struct {
	Port string
}

type Db struct {
	Dsn string
}

func LoadEnv() *Env {
	env := &Env{
		App: App{
			Port: os.Getenv("APP_PORT"),
		},
		Db: Db{
			Dsn: os.Getenv("DB_DATA_SOURCE_NAME"),
		},
	}

	return env
}
