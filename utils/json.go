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

	file.Truncate(0)
	file.Seek(0,0)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(*data)
	Check(err)
}

func IsCallerCreator(fileName string, guildID string, callerID string) bool {
	var jsonData map[string]structs.GuildData
	ReadJSONFile(fileName, &jsonData)

	data, ok := jsonData[guildID]

	return ok && (data.Creator == callerID)
}

func IsEventStarted(fileName string, guildID string) bool {
	var jsonData map[string]structs.GuildData
	ReadJSONFile(fileName, &jsonData)

	data, ok := jsonData[guildID]

	return ok && !data.Ended
}

func IsEventEnded(fileName string, guildID string) bool {
	var jsonData map[string]structs.GuildData
	ReadJSONFile(fileName, &jsonData)

	data, ok := jsonData[guildID]

	return ok && data.Ended
}

func IsUserRegistered(fileName string, guildID string, userID string) bool {
	var jsonData map[string]structs.GuildData
	ReadJSONFile(fileName, &jsonData)

	data, ok := jsonData[guildID]

	if !ok {
		return false
	}

	_, userOk := data.Responses[userID]

	return userOk
}