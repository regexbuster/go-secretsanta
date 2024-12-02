package main
	
import (
	"fmt"

	"secretsanta/utils"
	"secretsanta/structs"
)

// type GuildData struct {
// 	Creator 	string				`json:creator`
// 	EmbedID		string				`json:embedID`
// 	EmbedData	EmbededData			`json:embedData`
// 	Responses	map[string]string	`json:responses`
// 	Santas		map[string]string	`json:santas`
// 	Ended		bool				`json:ended`
// }

// type EmbededData struct {
// 	Name		string				`json:name`
// 	Price		string				`json:price`
// 	Register	string				`json:register`
// 	Deadline	string				`json:deadline`
// 	Notes		string				`json:notes`
// }

// func check(e error){
// 	if e != nil {
// 		panic(e)
// 	}
// }	

// func readJSONFile[T any](fileName string, data *map[string]T) {
// 	file, err := os.Open(fileName)
// 	check(err)
// 	defer file.Close()

// 	decoder := json.NewDecoder(file)

// 	// decoder.Token()

// 	for decoder.More(){
// 		decoder.Decode(data)
// 	}
// }

// func writeJSONFile[T any](fileName string, data *map[string]T) {
// 	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
// 	check(err)
// 	defer file.Close()

// 	encoder := json.NewEncoder(file)
// 	err = encoder.Encode(*data)
// 	check(err)
// }

func main() {
	// dat, err := os.ReadFile("./tmp/data.json")
	// check(err)
	// // fmt.Println(string(dat))

	// var jsonData map[string]GuildData

	// if err := json.Unmarshal(dat, &jsonData); err != nil{
	// 	panic(err)
	// }

	// fmt.Println(jsonData)

	// for key, val := range jsonData {
	// 	fmt.Printf("%s | %s\n", key, val.Responses)
	// }

	var jsonData map[string]structs.GuildData

	utils.ReadJSONFile("./data/secretsanta.json", &jsonData)

	fmt.Printf("%s \n %T\n", jsonData, jsonData)

	utils.WriteJSONFile("./tmp/data2.json", &jsonData)
}