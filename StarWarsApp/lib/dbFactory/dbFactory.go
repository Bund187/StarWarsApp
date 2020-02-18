package dbFactory

import (
    "StarWarsApp/StarWarsApp/lib/jsonController"
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

func CreateDataBase(database *sql.DB) {
    //CREATE CHARACTER TABLE
    statement, error := database.Prepare("CREATE TABLE IF NOT EXISTS character (idCharacter INTEGER PRIMARY KEY AUTOINCREMENT,"+
                                                                            "name TEXT NOT NULL,"+
                                                                            "height TEXT,"+
                                                                            "mass TEXT,"+
                                                                            "hairColor TEXT,"+
                                                                            "skinColor TEXT,"+
                                                                            "eyeColor TEXT,"+
                                                                            "birthYear TEXT,"+
                                                                            "gender TEXT,"+
                                                                            "homeWorld TEXT,"+
                                                                            "species TEXT)")
    if error != nil {
        fmt.Println("Failed to create the source database \"character\" table:", error)
    }
    statement.Exec()

    //CREATE STARSHIP TABLE
    statement, error = database.Prepare("CREATE TABLE IF NOT EXISTS starship (idStarship INTEGER PRIMARY KEY AUTOINCREMENT, starshipName TEXT UNIQUE)")
    if error != nil {
        fmt.Println("Failed to create the source database \"starship\" table:", error)
    }
    statement.Exec()

    //CREATE PILOT TABLE (RELATION BETWEEN CHARACTERS AND STARSHIPS)
    statement, error = database.Prepare("CREATE TABLE IF NOT EXISTS pilot (idPilot INTEGER PRIMARY KEY AUTOINCREMENT,"+
                                                                        "character TEXT,"+
                                                                        "starship TEXT,"+
                                                                        "FOREIGN KEY(character) REFERENCES character(idCharacter),"+
                                                                        "FOREIGN KEY(starship) REFERENCES starship(idStarship))")
    if error != nil {
        fmt.Println("Failed to create the source database \"pilot\" table:", error)
    }
    statement.Exec()
    
    //CREATE VEHICLE TABLE
    statement, error = database.Prepare("CREATE TABLE IF NOT EXISTS vehicle (idVehicle INTEGER PRIMARY KEY AUTOINCREMENT, vehicleName TEXT UNIQUE)")
    if error != nil {
        fmt.Println("Failed to create the source database \"vehicle\" table:", error)
    }
    statement.Exec()

    statement, error = database.Prepare("CREATE TABLE IF NOT EXISTS driver (idDriver INTEGER PRIMARY KEY AUTOINCREMENT,"+
                                                                            "character TEXT,"+
                                                                            "vehicle TEXT,"+
                                                                            "FOREIGN KEY(character) REFERENCES character(idCharacter),"+
                                                                            "FOREIGN KEY(vehicle) REFERENCES vehicle(idVehicle))")
    if error != nil {
        fmt.Println("Failed to create the source database \"driver\" table:", error)
    }
    statement.Exec()

    //CREATE FILM TABLE
    statement, error = database.Prepare("CREATE TABLE IF NOT EXISTS film (idFilm INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT UNIQUE)")
    if error != nil {
        fmt.Println("Failed to create the source database \"film\" table:", error)
    }
    statement.Exec()

    statement, error = database.Prepare("CREATE TABLE IF NOT EXISTS filmCharacter (idfilmCharacter INTEGER PRIMARY KEY AUTOINCREMENT,"+
                                                                                    "character TEXT,"+
                                                                                    "film TEXT,"+
                                                                                    "FOREIGN KEY(character) REFERENCES character(idCharacter),"+
                                                                                    "FOREIGN KEY(film) REFERENCES film(idFilm))")
    if error != nil {
        fmt.Println("Failed to create the source database \"filmCharacter\" table:", error)
    }
    statement.Exec()

    fmt.Println("DATABASE CREATED SUCCESFULLY")
    fmt.Printf("Retrieving json data. Wait a moment please.")
    jsonController.GetJson(database);
}