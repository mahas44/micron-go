package shared

import "time"

type Configuration struct {
	AMQPURL string
	SQLURL  string
}

var Config = Configuration{
	AMQPURL: "amqp://guest:guest@localhost:5672/",
	SQLURL:  "sqlserver=localhost;user id=sa;password=Password12*;port=1433;database=Game;",
}

type Platform int8

const (
	PC Platform = iota
	PS
	XBOX
	MOBILE
)

type UnitOfStorage int8

const (
	KB UnitOfStorage = iota
	MB
	GB
)

type AddGame struct {
	Name          string
	Description   string
	Price         float32
	Platform      int8
	Publisher     string
	Storage       float32
	UnitOfStorage int8
	ReleaseDate   time.Time
	IsDeleted     bool
}

func GetUnitOfStorageName(value int8) string {
	switch value {
	case 0:
		return "KB"
	case 1:
		return "MB"
	case 2:
		return "GB"
	default:
		return ""
	}
}

func GetPlatformName(value int8) string {
	switch value {
	case 0:
		return "PC"
	case 1:
		return "PS"
	case 2:
		return "XBOX"
	case 3:
		return "MOBILE"
	default:
		return ""
	}
}
