package jsonController

import(

	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type Character struct {
    Name string `json:"name"`
    Height string `json:"height"`
    Mass string `json:"mass"`
    HairColor string `json:"hair_color"`
    SkinColor string `json:"skin_color"`
    EyeColor string `json:"eye_color"`
    BirthYear string `json:"birth_year"`
    Gender string `json:"gender"`
    HomeWorld string `json:"homeWorld"`
    Films []string `json:"films"`
    Species []string `json:"species"`
    Vehicles []string `json:"vehicles"`
    Starships []string `json:"starships"`
}

type CharacterAPIResponse struct {
	Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Characters []Character `json:"results"`
}

type Film struct {
    Title string `json:"title"`
}

type Planet struct {
    Name string `json:"name"`
}

type Specie struct {
    Name string `json:"name"`
}

type Vehicle struct {
    Name string `json:"name"`
}

type Starship struct {
    Name string `json:"name"`
}

//MAIN FUNCTION THAT IS CALLED THE FIRST TIME
func GetJson(database *sql.DB) {
	var characterAPIResponseStruct CharacterAPIResponse
	jsonIntoDB(database, "https://swapi.co/api/people/",&characterAPIResponseStruct);
}

//FUNCTION THAT ITERATES OVER AN URL AND RETRIEVE THE JSON INFO AND STORES IN THE DATABASE
func jsonIntoDB(database *sql.DB, url string, structInterface interface{}) interface{}{
    fmt.Printf(".")
    res, err := http.Get(url)
    if err != nil {
        panic(err.Error())
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err.Error())
    }
    //STORES THE INTERFACE INTO A VARIABLE SO WE CAN USE IT AS ANY STRUCT TYPE LATER ON
    structDummy := structInterface
    error := json.Unmarshal(body, &structDummy)
    if(error != nil){
        fmt.Println("An error has ocurred:", error)
    }
    //WE GET THE STRUCT TYPE, IF ITS A CharacterAPIResponse TYPE WE PROCESS IT BECAUSE IT'S A CHARACTER AND WE NEED TO STORE THE INFO INTO THE DATABASE
    //EVERYTIME WE CALL THE FUNCION jsonIntoDB() IT RETURNS THE DESIRED STRUCT SO WE CAN EXTRACT ANY INFORMATION FROM IT
    switch result := structDummy.(type) {
        case *CharacterAPIResponse:	
           for i:=0;i<len(result.Characters);i++{
                //GET THE PLANET NAME
                var planetStruct Planet
	    		planet:=jsonIntoDB(database,result.Characters[i].HomeWorld,&planetStruct).(*Planet);
	    		result.Characters[i].HomeWorld=planet.Name
	    		//GET THE SPECIE NAME
                var specieName string
                if len(result.Characters[i].Species)>0{
                    var specieStruct Specie
                    specie:=jsonIntoDB(database,result.Characters[i].Species[0],&specieStruct).(*Specie);
                    specieName=specie.Name
                }else{
                    specieName=""
                }
                //INSERT BASIC INFO INTO THE CHARACTER'S TABLE
                statement, error := database.Prepare("INSERT INTO character (name, height, mass, hairColor, skinColor, eyeColor, birthYear, gender, homeWorld, species) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
                if error != nil {
                    fmt.Println("Failed to insert a row into the source database \"character\" table:", error)
                }
                statement.Exec(result.Characters[i].Name,result.Characters[i].Height,result.Characters[i].Mass,result.Characters[i].HairColor,result.Characters[i].SkinColor,result.Characters[i].EyeColor,result.Characters[i].BirthYear,result.Characters[i].Gender,result.Characters[i].HomeWorld,specieName)
                
                //ADDING VEHICLES
		    	for j:=0;j<len(result.Characters[i].Films);j++{
		    		var filmStruct Film
		    		film:=jsonIntoDB(database,result.Characters[i].Films[j],&filmStruct).(*Film);
		    		result.Characters[i].Films[j]=film.Title
		    		
                    //CHECK IF THE FILM EXISTS IN THE DATABASE, IF DOESN'T EXISTS WE DONT INSERT IT INTO THE DB
                    rows, e := database.Query("SELECT COUNT(title) FROM film WHERE title='"+film.Title+"'")
                    if e != nil {
                        fmt.Println("Failed to reach the \"film\" table:", e)
                    }
                    var repeat int
                    for rows.Next() {
                        rows.Scan(&repeat)
                    }
                    if(repeat==0){
                        statement, error = database.Prepare("INSERT INTO film (title) VALUES (?)")
                        if error != nil {
                            fmt.Println("Failed to insert a row into the source database \"film\" table:", error)
                        }
                        statement.Exec(result.Characters[i].Films[j])
                    }
                    //WE GET THE FILM ID
                    rows, e = database.Query("SELECT idFilm FROM film WHERE title='"+film.Title+"'")
                    if e != nil {
                        fmt.Println("Failed to reach the \"film\" table:", e)
                    }
                    var idFilm int
                    for rows.Next() {
                        rows.Scan(&idFilm)
                    }
                    //INSERT BOTH, THE CHARACTER ID AND FILM ID, IN THEIR RELATIONAL TABLE
                    statement, error = database.Prepare("INSERT INTO filmCharacter (character,film) VALUES (?, ?)")
                    if error != nil {
                        fmt.Println("Failed to insert a row into the source database \"filmCharacter\" table:", error)
                    }
                    statement.Exec(countReg(database,"name", "character"),idFilm)
                
		    	}

		    	//ADDING VEHICLES
		    	for j:=0;j<len(result.Characters[i].Vehicles);j++{
		    		var vehicleStruct Vehicle
		    		vehicle:=jsonIntoDB(database,result.Characters[i].Vehicles[j],&vehicleStruct).(*Vehicle);
		    		result.Characters[i].Vehicles[j]=vehicle.Name
		    		
                    //CHECK IF THE VEHICLE EXISTS IN THE DATABASE, IF DOESN'T EXISTS WE DONT INSERT IT INTO THE DB
                    rows, e := database.Query("SELECT COUNT(vehicleName) FROM vehicle WHERE vehicleName='"+vehicle.Name+"'")
                    if e != nil {
                        fmt.Println("Failed to reach the \"vehicle\" table:", e)
                    }
                    var repeat int
                    for rows.Next() {
                        rows.Scan(&repeat)
                    }
                    if(repeat==0){
                        statement, error = database.Prepare("INSERT INTO vehicle (vehicleName) VALUES (?)")
                        if error != nil {
                            fmt.Println("Failed to insert a row into the source database \"vehicles\" table:", error)
                        }
                        statement.Exec(result.Characters[i].Vehicles[j])
                    }
                    //WE GET THE VEHICLE ID
                    rows, e = database.Query("SELECT idVehicle FROM vehicle WHERE vehicleName='"+vehicle.Name+"'")
                    if e != nil {
                        fmt.Println("Failed to reach the \"vehicle\" table:", e)
                    }
                    var idVehicle int
                    for rows.Next() {
                        rows.Scan(&idVehicle)
                    }
                    //INSERT BOTH, THE CHARACTER ID AND STARSHIP ID, IN THEIR RELATIONAL TABLE
                    statement, error = database.Prepare("INSERT INTO driver (character,vehicle) VALUES (?, ?)")
                    if error != nil {
                        fmt.Println("Failed to insert a row into the source database \"driver\" table:", error)
                    }
                    statement.Exec(countReg(database,"name", "character"),idVehicle)
		    	}

		    	//ADDING STARSHIPS
                for j:=0;j<len(result.Characters[i].Starships);j++{
		    		var StarshipStruct Starship
		    		starship:=jsonIntoDB(database,result.Characters[i].Starships[j],&StarshipStruct).(*Starship);
		    		result.Characters[i].Starships[j]=starship.Name
                    
                    //CHECK IF THE STARSHIP EXISTS IN THE DATABASE, IF DOESN'T EXISTS WE DONT INSERT IT INTO THE DB
                    rows, e := database.Query("SELECT COUNT(starshipName) FROM starship WHERE starshipName='"+starship.Name+"'")
                    if e != nil {
                        fmt.Println("Failed to reach the \"starshipName\" table:", e)
                    }
                    var repeat int
                    for rows.Next() {
                        rows.Scan(&repeat)
                    }
                    if(repeat==0){
                        statement, error = database.Prepare("INSERT INTO starship (starshipName) VALUES (?)")
                        if error != nil {
                            fmt.Println("Failed to insert a row into the source database \"starship\" table:", error)
                        }
                        statement.Exec(result.Characters[i].Starships[j])
                    }
                    //WE GET THE STARSHIP ID
                    rows, e = database.Query("SELECT idStarship FROM starship WHERE starshipName='"+starship.Name+"'")
                    if e != nil {
                        fmt.Println("Failed to reach the \"starshipName\" table:", e)
                    }
                
                    var idStars int
                    for rows.Next() {
                        rows.Scan(&idStars)
                    }
                    //INSERT BOTH, THE CHARACTER ID AND STARSHIP ID, IN THEIR RELATIONAL TABLE
                    statement, error = database.Prepare("INSERT INTO pilot (character,starship) VALUES (?, ?)")
                    if error != nil {
                        fmt.Println("Failed to insert a row into the source database \"pilot\" table:", error)
                    }
                    statement.Exec(countReg(database,"name", "character"),idStars)
		    	}
		    }
        //GET THE NEXT CHARACTER PAGE AND TAKE ITS INFO
        if result.Next!="" {
            var characterAPIResponseStruct CharacterAPIResponse
            jsonIntoDB(database, result.Next,&characterAPIResponseStruct);
        }
	}
    return structDummy
}

//FUNCTION THAT COUNTS THE REGISTERS IN A TABLE
func countReg(database *sql.DB, column string, table string) (count int){

    rows, e := database.Query("SELECT COUNT("+column+") FROM "+table)
    if e != nil {
        fmt.Println("Failed to reach the \""+table+"\" table:", e)
    }
    for rows.Next() {
        rows.Scan(&count)
    }
    return count
}