package infobook

type Info struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Book   string `json:"Book"`
}

func PrepareResponse() []Info {
	var infos []Info

	var book Info
	book.Id = 1
	book.Author = "Charlotte Bronte"
	book.Book = "Jane Eyr"
	infos = append(infos, book)

	book.Id = 2
	book.Author = "Mikhail Bulgakov"
	book.Book = "Master and Margarita"
	infos = append(infos, book)

	book.Id = 3
	book.Author = "Jack London"
	book.Book = "Martin Iden"
	infos = append(infos, book)

	book.Id = 4
	book.Author = "Ray Bradbury"
	book.Book = "Fahrenheit 451"
	infos = append(infos, book)
	return infos
}
