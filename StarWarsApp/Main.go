package main

import (
    "database/sql"
    "StarWarsApp/StarWarsApp/lib/dbFactory"
    "fmt"
    "strconv"
    _ "github.com/mattn/go-sqlite3"
    "os"
    "bufio"
    "github.com/fatih/color"
)

func main() {
    
    database, err := sql.Open("sqlite3", "Database/StarWarsDB.db")
    defer database.Close()
    if err != nil {
        fmt.Println("Failed to open the source database:", err)
    }

    //CHECK IF THE DATABASE EXISTS, IF NOT, CREATE IT
    var doesTestTableExist bool
    err = database.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'character' LIMIT 1)").Scan(&doesTestTableExist)
    if err != nil {
        fmt.Println("Failed to check if the \"character\" table exists in the destination database:", err)
    }
    if !doesTestTableExist {
        dbFactory.CreateDataBase(database);

    }
    
    //CHECK IF ANY ARGUMENT HAS BEEN INPUT
    if len(os.Args)>1{
        showCharacter(database,os.Args[1])
		menu(database)

    }else{
        fmt.Println("You haven't input any character. Type: StarWarsApp nameOfCharacter")
    }
}

//SHOWS THE MENU AND MANAGE THE OPTIONS
func menu(database *sql.DB){
	var option string
        for option!="0"{
            fmt.Println("\nChoose one of the following options?")
            fmt.Println("\tEnter another character's name   (1)")
            fmt.Println("\tView the list of character       (2)")
            fmt.Println("\tView the list of vehicles        (3)")
            fmt.Println("\tView the list of starships       (4)")
            fmt.Println("\tView the list of films           (5)")
            fmt.Println("\tExit                             (0)")
            fmt.Scanf("%s\n", &option)
           
            switch option{
                case "1":   
                            var charName string
                            fmt.Println("Character name:")
                            scanner := bufio.NewScanner(os.Stdin)
                            scanner.Scan()
                            charName = scanner.Text()
                            showCharacter(database,charName)
                case "2":   color.Red("\n\tCHARACTERS:")   
                            selectAll(database, "name", "character")
                case "3":   color.Red("\n\tVEHICLES:")   
                            selectAll(database, "vehicleName", "vehicle")
                case "4":   color.Red("\n\tSTARSHIPS:")   
                            selectAll(database, "starshipName", "starship")
                case "5":   color.Red("\n\tFILMS:")   
                            selectAll(database, "title", "film")
                case "0":   color.Red("Good Bye! May the force be with you!")
                default:    fmt.Println("Wrong option, try again")
            }
        }
}

// SELECT ALL QUERY
func selectAll(database *sql.DB, regName string, table string){
    rows, e := database.Query("SELECT "+regName+" FROM "+table+" ORDER BY "+regName)
    if e != nil {
        fmt.Println("Failed to reach the \""+table+"\" table:", e)
    }
    var name string
    for rows.Next() {
        rows.Scan(&name)
        fmt.Println("\t-"+name)
    }
}

// SHOWS A CHARACTER FROM THE DATABASE
func showCharacter(database *sql.DB, character string){
    rows, e := database.Query("SELECT * FROM character Where name LIKE "+"'"+character+"%'")
    if e != nil {
        fmt.Println("Failed to reach the \"character\" table:", e)
    }
    defer rows.Close()
    var id int
    var name, height, mass, hairColor, skinColor, eyeColor, birthYear, gender, homeWorld, species string
    count:=0
    for rows.Next() {
        count++
        e=rows.Scan(&id, &name, &height,&mass, &hairColor, &skinColor, &eyeColor, &birthYear, &gender, &homeWorld, &species)
        if e != nil {
            fmt.Println(e)
        }
        
        color.Red("\n\t\t"+name)
        color.Yellow("\tHeight: \t\t"+ height)
        color.Yellow("\tMass: \t\t\t"+ mass)
        color.Yellow("\tHair color: \t\t"+ hairColor)
        color.Yellow("\tSkin color: \t\t"+ skinColor)
        color.Yellow("\tEye color: \t\t"+ eyeColor)
        color.Yellow("\tBirth year: \t\t"+ birthYear)
        color.Yellow("\tGender: \t\t"+ gender)
        color.Yellow("\tHome world: \t\t"+ homeWorld)
        color.Yellow("\tSpecies: \t\t"+ species)
        
        //GET STARSHIPS
        color.Green("\tStarship/s:")
        rows, e = database.Query("SELECT starshipName FROM starship s INNER JOIN pilot p ON p.starship=s.idStarship INNER JOIN character c ON c.idCharacter=p.character and c.idCharacter='"+strconv.Itoa(id)+"'")
        if e != nil {
            fmt.Println("Failed to reach the \"pilot\" table:", e)
        }
        var starshipName string
        for rows.Next() {
            rows.Scan(&starshipName)
             color.Green("\t\t\t\t"+starshipName)
        }

        //GET VEHICLES
        color.Cyan("\n\tVehicle/s:")
        rows, e = database.Query("SELECT vehicleName FROM vehicle v INNER JOIN driver d ON d.vehicle=v.idVehicle INNER JOIN character c ON c.idCharacter=d.character and c.idCharacter='"+strconv.Itoa(id)+"'")
        if e != nil {
            fmt.Println("Failed to reach the \"pilot\" table:", e)
        }
        var vehicleName string
        for rows.Next() {
            rows.Scan(&vehicleName)
             color.Cyan("\t\t\t\t"+vehicleName)
        }

        //GET FILMS
        color.Magenta("\n\tFilm/s:")
        rows, e = database.Query("SELECT title FROM film f INNER JOIN filmCharacter fc ON fc.film=f.idFilm INNER JOIN character c ON c.idCharacter=fc.character and c.idCharacter='"+strconv.Itoa(id)+"'")
        if e != nil {
            fmt.Println("Failed to reach the \"pilot\" table:", e)
        }
        var filmTitle string
        for rows.Next() {
            rows.Scan(&filmTitle)
            color.Magenta("\t\t\t\t"+filmTitle)
        }

    }
    if count==0{
        fmt.Println("The character you entered doesn't exist")
    }
}
