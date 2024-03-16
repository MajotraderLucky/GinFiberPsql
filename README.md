#### Проект, в котором применяется jwt-токенизация находится здесь:
https://github.com/MajotraderLucky/GinFiberPsql

В пакете [**postgrdb**](https://github.com/MajotraderLucky/GinFiberPsql/tree/main/postgrdb) есть описание структуры БД postgreSQL в файле **postgrdb/init.sql**:
`-- Creating the fio_data table`

`CREATE TABLE fio_data (`

`id SERIAL PRIMARY KEY,`

`name TEXT,`

`surname TEXT,`

`patronymic TEXT,`

`age INTEGER,`

`gender TEXT,`

`nationality TEXT,`

`error_reason TEXT,`

`created_at TIMESTAMP DEFAULT current_timestamp,`

`updated_at TIMESTAMP DEFAULT current_timestamp`

`);`

`-- Creating an index on the name field`

`CREATE INDEX idx_fio_data_name ON fio_data (name);`

`-- Creating a composite index on the name and surname fields`

`CREATE INDEX idx_fio_data_name_surname ON fio_data (name, surname);`

-------------------------------------------------------------

*Чтобы добавить данные в базу данных используем такой handler фреймворк Gin:*

i`func addDataHandler(db *sql.DB) gin.HandlerFunc {`

`return func(c *gin.Context) {`

`// Structure for parsing the request body`

`var data struct {`

`Name string json:"name"`

`Surname string json:"surname"`

`Patronymic string json:"patronymic"`

`Age int json:"age"`

`Gender string json:"gender"`

`Nationality string json:"nationality"`

`}`

`// Parsing the JSON body of the request`

`if err := c.ShouldBindJSON(&data); err != nil {`

`c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})`

`return`

`}`

`// SQL query to add data`

`query := INSERT INTO fio_data (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)`

`_, err := db.Exec(query, data.Name, data.Surname, data.Patronymic, data.Age, data.Gender, data.Nationality)`

`if err != nil {`

`log.Printf("Error inserting data: %v", err)`

`c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data"})`

`return`

`}`
`// Response about successful data addition`
`c.JSON(http.StatusOK, gin.H{"message": "Data added successfully"})`
`}`
`}`

----------------------------------------------------------------------
### Реализация в функции main() без использования Jwt
`// Registering the handler for adding data`
`router.POST("/add_data", addDataHandler(db))`

### Для авторизации с помощью jwt используем пакет: 
**github.com/golang-jwt/jwt/v5*

#### Создаём функцию, которая будет проверять запрос следующим образом:
==`var jwtKey = []byte("8GoPUxkoCEeKaEG381hL6p9RAfwgCaiDJhrwy+/k8Og=") // Replace with your key==

`// JWTAuthMiddleware checks for the presence and validity of a JWT token`

`func JWTAuthMiddleware() gin.HandlerFunc {`

`return func(c *gin.Context) {`

`authHeader := c.GetHeader("Authorization")`

`if authHeader == "" {`

`c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})`

`c.Abort()`

`return`

`}`

`tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))`

`token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {`

`// Ensure the token signature algorithm is expected`

`if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {`

`return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])`

`}`

`return jwtKey, nil`

`})`

`if err != nil {`

`c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})`

`c.Abort()`

`return`

`}`

`if !token.Valid {`

`c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})`

`c.Abort()`

`return`

`}`

`c.Next()`

`}`

`}`

### Теперь можем использовать авторизацию при запросе к handler:
`router.POST("/add_data", JWTAuthMiddleware(), addDataHandler(db))`

### Откуда берем ключ?
==var jwtKey = []byte("8GoPUxkoCEeKaEG381hL6p9RAfwgCaiDJhrwy+/k8Og=") // Replace with your key====

Можно использовать  bash-script:

`#!/bin/bash`
`openssl rand -base64 32`

### Как выглядит запрос к handler Gin без авторизации?

`curl -X POST http://localhost:8085/add_data      -H "Content-Type: application/json"      -d '{"name":"Иван", "surname":"Иванов", "patronymic":"Иванович", "age":30, "gender":"мужской", "nationality":"Россия"}'
`
#### Ключ jwtKey нужен для проверки и для создания ключа запроса.

### Создаём ключ запроса:
Для лучшего понимания я создаю ключ запроса вручную, но, разумеется, процесс можно автоматизировать. Берём пакет [tokengenerator](https://github.com/MajotraderLucky/GinFiberPsql/tree/main/tokengenerator)
`package main
`
``import (`
`"fmt"`
`"time"`
`"github.com/golang-jwt/jwt/v5"`
`)`
==//Здесь вставляем ключ, который генерировали с помощью bash-скрипта==
`var jwtKey = []byte("8GoPUxkoCEeKaEG381hL6p9RAfwgCaiDJhrwy+/k8Og=")`

`type MyCustomClaims struct {`

`Username string json:"username"`

`jwt.RegisteredClaims`

`}`

`func main() {`

`// Создаем утверждения`

`claims := MyCustomClaims{`

`Username: "exampleUser",`

`RegisteredClaims: jwt.RegisteredClaims{`

==`ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен истекает через 24 часа`==

`},`

`}`
`// Создаем токен`

`token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)`

`// Подписываем токен нашим ключом`

`tokenString, err := token.SignedString(jwtKey)`

`if err != nil {`

`fmt.Println("Ошибка при создании токена:", err)`

`return`

`}`
`fmt.Println("Сгенерированный JWT:", tokenString)`
`}`

-------------------------------------------------------------------------

==`go run tokengenerator.go`== 
`Сгенерированный JWT: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGVVc2VyIiwiZXhwIjoxNzEwNjcyNTkwfQ.zshKwa4fyW-MLMEetGKQymjYL9xkv3oouvDop0G1MXk`

**Этот сгенерированный ключ используем для запроса**

### Запрос с авторизацией выглядит так:

`curl -X POST http://localhost:8085/add_data \`
     `-H "Content-Type: application/json" \`
     `-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGVVc2VyIiwiZXhwIjoxNzEwNjY2ODY3fQ.uVIJTrrpak1HLD8LDopX_O2nLT3UC1c0-F49MeUUOTA" \`
     `-d '{"name":"Sergey", "surname":"Ivanov", "patronymic":"Viktorovich", "age":40, "gender":"мужской", "nationality":"Россия"}'`