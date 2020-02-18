package dbFactory

import (
    "StarWarsApp/lib/jsonController"
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

func CreateDataBase(database *sql.DB) {
    /*database, err := sql.Open("sqlite3", "../../Database/StarWarsDB.db")
    if err != nil {
        fmt.Println("Failed to open the source database:", err)
    }*/

    /*var doesTestTableExist bool
    err = database.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'character' LIMIT 1)").Scan(&doesTestTableExist)
    if err != nil {
        fmt.Println("Failed to check if the \"character\" table exists in the destination database:", err)
    }
    if !doesTestTableExist {
        fmt.Println("The \"character\" table could not be found in the destination database.")
    }*/

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
    //statement, _ = database.Prepare("INSERT INTO character (idPersonaje, name, height) VALUES (?, ?, ?)")
    //statement.Exec(1,"Guille","1'72")
    //statement.Exec(2,"Harry", "2'00")
    //statement.Exec(3,"Chewi", "1'50")

    /*fmt.Println("MOSTRAR LOS PERSONAJES")
    rows1, _ := database.Query("SELECT * FROM character")
    var id int
    var name string
    var height string
    for rows1.Next() {
        rows1.Scan(&id, &name, &height)
        fmt.Println(strconv.Itoa(id) + ": " + name + " " + height)
    }*/

    statement, error = database.Prepare("CREATE TABLE IF NOT EXISTS starship (idStarship INTEGER PRIMARY KEY AUTOINCREMENT, starshipName TEXT UNIQUE)")
    if error != nil {
        fmt.Println("Failed to create the source database \"starship\" table:", error)
    }
    statement.Exec()
    /*statement2, _ = database.Prepare("INSERT INTO nave (idNave, name) VALUES (?,?)")
    statement2.Exec(1,"X-Wing")
    statement2.Exec(2,"Halcón Milenario")

    fmt.Println("MOSTRAR LAS NAVES")
    rows2, _ := database.Query("SELECT * FROM nave")
    var idNave int
    var nombreNave string
    for rows2.Next() {
        rows2.Scan(&idNave, &nombreNave)
        fmt.Println(strconv.Itoa(idNave) + ": " + nombreNave)
    }*/

    statement, error = database.Prepare("CREATE TABLE IF NOT EXISTS pilot (idPilot INTEGER PRIMARY KEY AUTOINCREMENT,"+
                                                                        "character TEXT,"+
                                                                        "starship TEXT,"+
                                                                        "FOREIGN KEY(character) REFERENCES character(idCharacter),"+
                                                                        "FOREIGN KEY(starship) REFERENCES starship(idStarship))")
    if error != nil {
        fmt.Println("Failed to create the source database \"pilot\" table:", error)
    }
    statement.Exec()
    /*statement, _ = database.Prepare("INSERT INTO piloto (personaje, starship) VALUES (?, ?)")
    statement.Exec("1", "1")
    statement.Exec("2", "2")
    statement.Exec("3", "2")

    fmt.Println("MOSTRAR QUIEN PILOTA QUE starship FORMATEADO")
    rows3, _ := database.Query("SELECT piloto.id, piloto.personaje, character.name, character.height, piloto.starship, starship.name FROM piloto,character,starship WHERE piloto.personaje=character.idPersonaje AND piloto.starship=starship.idNave")
    var idPiloto int
    var idPersonaje int
    var idNavePiloto int
    var nombrePersonaje string
    var heightPersonaje string
    var nombreNavePiloto string
    for rows3.Next() {
        rows3.Scan(&idPiloto, &idPersonaje, &nombrePersonaje, &heightPersonaje, &idNavePiloto, &nombreNavePiloto)
        fmt.Println(strconv.Itoa(idPiloto) + ": " + strconv.Itoa(idPersonaje) + " " + nombrePersonaje+ " " + heightPersonaje + ": " + strconv.Itoa(idNavePiloto) + " " + nombreNavePiloto)
    }*/

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
    fmt.Println("Retrieving json data.")
    jsonController.GetJson(database);
}