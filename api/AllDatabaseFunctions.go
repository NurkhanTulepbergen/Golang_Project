package api

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Nurkhan05"
	dbname   = "ayazhanchert"
)

// -------------------------------------------------------------------------------------------------------
func createUser(name1, password1 string) { // добавление юзера к базе данных
	// Строка подключения к базе данных PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Установка соединения с базой данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id int // Объявляем переменную id для сохранения возвращаемого значения из запроса
	insertStmt := `INSERT INTO "user"(name, password) VALUES($1, $2) RETURNING id`
	// Выполняем запрос и сканируем результат в переменную id
	err = db.QueryRow(insertStmt, name1, password1).Scan(&id)
	if err != nil {
		panic(err)
	}
	// Используем формат вывода %d для вывода значения id
	fmt.Printf("Добавлен user с id: %d\n", id)
}

func deleteUser(id1 int) { // удаление юзера с базы данных
	// Строка подключения к базе данных PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Установка соединения с базой данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id int // Объявляем переменную id для сохранения возвращаемого значения из запроса

	deleteStmt := `DELETE FROM "user" WHERE id = $1`
	_, err = db.Exec(deleteStmt, id1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Пользователь с ID %d удален\n", id)
}

func viewUsers() {
	// Строка подключения к базе данных PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Установка соединения с базой данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT \"user\".id, name, password FROM \"user\"")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Проходимся по каждой строке результата запроса
	for rows.Next() {
		var id int
		var name, password string
		// Считываем значения столбцов из текущей строки
		if err := rows.Scan(&id, &name, &password); err != nil {
			panic(err)
		}
		// Выводим информацию о пользователе
		fmt.Printf("ID: %d, Name: %s, Password: %s\n", id, name, password)
	}
	// Проверяем наличие ошибок после завершения цикла
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

//-------------------------------------------------------------------------------------------------------

func viewProducts() {
	// Database connection string for PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Establishing a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description, price FROM products")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Iterating through each row of the query result
	for rows.Next() {
		var id, price int
		var title, description string
		// Reading values from the current row's columns
		if err := rows.Scan(&id, &title, &description, &price); err != nil {
			panic(err)
		}
		// Displaying information about the product
		fmt.Printf("ID: %d, Title: %s, Description: %s, Price: %d\n", id, title, description, price)
	}
	// Checking for errors after the loop ends
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

func createProduct(title, description string, price int) {
	// Строка подключения к базе данных PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Установка соединения с базой данных
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id int // Объявляем переменную id для сохранения возвращаемого значения из запроса
	insertStmt := `INSERT INTO products(title, description, price) VALUES($1, $2, $3) RETURNING id`
	// Выполняем запрос и сканируем результат в переменную id
	err = db.QueryRow(insertStmt, title, description, price).Scan(&id)
	if err != nil {
		panic(err)
	}
	// Используем формат вывода %d для вывода значения id
	fmt.Printf("Добавлен продукт с id: %d\n", id)
}

func deleteProduct(id1 int) {
	// Database connection string for PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Establishing a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	deleteStmt := `DELETE FROM products WHERE id = $1`
	// Executing the delete query for the product with the specified ID
	_, err = db.Exec(deleteStmt, id1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Product with ID %d deleted\n", id1)
}

//-------------------------------------------------------------------------------------------------------

func createCategory(title, description string) {
	// Database connection string for PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Establishing a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var id int
	insertStmt := `INSERT INTO category(title, description) VALUES($1, $2) RETURNING id`
	// Executing the query and scanning the result into the id variable
	err = db.QueryRow(insertStmt, title, description).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Product with ID %d added\n", id)
}

func deleteCategory(id1 int) {
	// Database connection string for PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Establishing a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	deleteStmt := `DELETE FROM category WHERE id = $1`
	// Executing the delete query for the product with the specified ID
	_, err = db.Exec(deleteStmt, id1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Product with ID %d deleted\n", id1)
}

func viewCategory() {
	// Database connection string for PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Establishing a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description FROM category")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Iterating through each row of the query result
	for rows.Next() {
		var id int
		var title, description string
		// Reading values from the current row's columns
		if err := rows.Scan(&id, &title, &description); err != nil {
			panic(err)
		}
		// Displaying information about the product
		fmt.Printf("ID: %d, Title: %s, Description: %s\n", id, title, description)
	}
	// Checking for errors after the loop ends
	if err := rows.Err(); err != nil {
		panic(err)
	}
}

//-------------------------------------------------------------------------------------------------------

func main() {
	//// Строка подключения к базе данных PostgreSQL
	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	host, port, user, password, dbname)
	//
	//// Установка соединения с базой данных
	//db, err := sql.Open("postgres", psqlInfo)
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	viewCategory()
	deleteCategory(2)
	viewCategory()
}