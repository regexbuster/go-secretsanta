package structs

type GuildData struct {
	Creator 	string				`json:creator`
	EmbedID		string				`json:embedID`
	EmbedData	EmbededData			`json:embedData`
	Responses	map[string]UserData	`json:responses`
	Santas		map[string]string	`json:santas`
	Ended		bool				`json:ended`
}

type EmbededData struct {
	Name		string	`json:name`
	Price		string	`json:price`
	Register	string	`json:register`
	Deadline	string	`json:deadline`
	Notes		string	`json:notes`
}

type UserData struct {
	Name		string	`json:name`
	Wishlist	string	`json:wishlist`
}