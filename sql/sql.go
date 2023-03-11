package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	shared "micron/shared"
	"strconv"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func SqlOpen() *sql.DB {
	db, err := sql.Open("sqlserver", shared.Config.SQLURL)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// func GetSqlContent(db *sql.DB) ([]string, []string, []float64, []string, []string, []float64, []int8, error) {
func GetSqlContent(db *sql.DB) ([]shared.AddGame, error) {
	var (
		Games []shared.AddGame
		// Name          []string
		// Description   []string
		// Price         []float64
		// Platform      []string
		// Publisher     []string
		// Storage       []float64
		// UnitOfStorage []int8
		ctx context.Context
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select Name, Description, Price, Platform, Publisher, Storage, UnitOfStorage from Micron.Games")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var _gameName string
		var _gameDescription string
		var _gamePrice float32
		var _gamePlatform int8
		var _gamePublisher string
		var _gameStorage float32
		var _gameUnitOfStorage int8
		var _gameReleaseDate time.Time

		err := rows.Scan(&_gameName, &_gameDescription, &_gamePrice, &_gamePlatform, &_gamePublisher, &_gameStorage, &_gameUnitOfStorage)

		if err != nil {
			// return Name, Description, Price, Platform, Publisher, Storage, UnitOfStorage, err
			return Games, err
		} else {
			game := shared.AddGame{
				Name:          _gameName,
				Description:   _gameDescription,
				Price:         _gamePrice,
				Platform:      _gamePlatform,
				Publisher:     _gamePublisher,
				Storage:       _gameStorage,
				UnitOfStorage: _gameUnitOfStorage,
				ReleaseDate:   _gameReleaseDate,
				IsDeleted:     false,
			}
			Games = append(Games, game)
			// Name = append(Name, _gameName)
			// Description = append(Description, _gameDescription)
			// Price = append(Price, _gamePrice)
			// Platform = append(Platform, _gamePlatform)
			// Publisher = append(Publisher, _gamePublisher)
			// Storage = append(Storage, _gameStorage)
			// UnitOfStorage = append(UnitOfStorage, _gameUnitOfStorage)
		}
	}
	// return Name, Description, Price, Platform, Publisher, Storage, UnitOfStorage, nil
	return Games, nil
}

func IntersSqlContent(db *sql.DB, game *shared.AddGame) (int64, error) {
	stmt, err := db.Prepare(`INSERT INTO Micron.Games(Name, Description, Price, Platform, Publisher, Storage, UnitOfStorage)
							VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7); 
							select ID = convert(bigint, SCOPE_IDENTITY())`)

	if err != nil {
		handleError(err, "Could nor insert SqlDB")
		return 0, err
	}

	var ctx context.Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	defer stmt.Close()

	rows := stmt.QueryRowContext(ctx, game.Name, game.Description, game.Price, game.Platform, game.Publisher, game.Storage, game.UnitOfStorage)
	if rows.Err() != nil {
		return 0, err
	}

	var _id int64
	rows.Scan(&_id)
	return _id, nil
}

func GetGames(db *sql.DB) {
	// SQL List All Product
	// names, descriptions, prices, platforms, publishers, storages, unitOfStorages, err := sql2.GetSqlContent(db)
	Games, err := GetSqlContent(db)
	if err != nil {
		fmt.Println("(sqlmicron) Error getting content: " + err.Error())
	}
	fmt.Println(strings.Repeat("-", 100))

	// Now read the contents
	for _, value := range Games {
		platform := value.Platform
		platformValue := shared.GetPlatformName(platform)
		// platformValue := ""
		// for v := range platform {
		// 	platformValue += shared.GetPlatformName(v)
		// }

		// unitOfStorage := shared.GetUnitOfStorageName(unitOfStorages[i])
		unitOfStorage := shared.GetUnitOfStorageName(value.UnitOfStorage)

		fmt.Println("Game: " + value.Name + "\nDescription: " + value.Description +
			"\nPrice: " + strconv.FormatFloat(float64(value.Price), 'f', 2, 64) + "\nPlatforms: " + platformValue +
			"\nPublisher: " + value.Publisher + "\nStorage: " + strconv.FormatFloat(float64(value.Storage), 'f', 2, 64) + " " + unitOfStorage)
		fmt.Println(strings.Repeat("-", 100))
	}
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
