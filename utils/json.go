package utils

import (
	"encoding/json"
	"os"
	"secretsanta/structs"
)

func Check(e error){
	if e != nil {
		panic(e)
	}
}	

func ReadJSONFile[T any](fileName string, data *map[string]T) {
	file, err := os.Open(fileName)
	Check(err)
	defer file.Close()

	decoder := json.NewDecoder(file)

	// decoder.Token()

	for decoder.More(){
		decoder.Decode(data)
	}
}

func WriteJSONFile[T any](fileName string, data *map[string]T) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	Check(err)
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(*data)
	Check(err)
}

func IsEventStarted(fileName string, guildID string) bool {
	var jsonData map[string]structs.GuildData
	ReadJSONFile(fileName, &jsonData)

	_, ok := jsonData[guildID]

	return ok
}