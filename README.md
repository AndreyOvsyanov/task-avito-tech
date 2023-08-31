task-avito-tech

Этот репозиторий является открытым кодом к решению тестового задания на стажировку в Avito

* В файле `mysql.sql` приведён sql код всех таблиц, которые находились у меня в локальной БД ( MySQL )
* В файлах `docker-compose.yaml` и `Dockerfile` код поднятия среды с помощью docker desktop, который запускает docker daemon
* `.env` - переменные окружения, которые используются для поднятия контейнера app и контейнера db в docker desktop

`server.go`


1. Подключение к базе данных на localhost по порту 3307 c именем my_database
   Последующие команды для работы docker - стенда
    ```
   db, error_db = sql.Open("mysql", "root:G12e70891@tcp(localhost:3307)/my_database")
    if error_db != nil {
      log.Fatal(error_db)
    }

   db.SetMaxOpenConns(1000)
   db.tMaxIdleConns(100)
   db.tConnMaxIdleTime(time.Minute * time.Duration(3))
   db.tConnMaxLifetime(time.Hour * time.Duration(1))

   defer func(result *sql.DB) { _ = result.Close() }(db)

   if  := db.Ping(); err != nil {
      log.Fatal(err)
   }
   ```
2. Route's для работы приложения
    ```
   router := mux.NewRouter()
   router.HandleFunc("/", defaultPage)
    router.HandleFunc("/users", getUserInfo)
    router.HandleFunc("/segments", getSegmentInfo)
    router.HandleFunc("/segment/create", createSegment)
    router.HandleFunc("/segment/delete", deleteSegment)
    router.HandleFunc("/segments/user", getUserSegments)
    router.HandleFunc("/useraddremove", actionUser)
    router.HandleFunc("/history/user", getUserHistory)

    log.Fatal(http.ListenAndServe(":8080", router))
   ```
3. `func getUserInfo(w http.ResponseWriter, r *http.Request)`
   
   Функция для вывода всех пользователей в JSON формате на веб-страницу с помощью json.MarshalIndent

   `Пример API: /users`


4. `func getSegmentInfo(w http.ResponseWriter, r *http.Request)`
   
   Аналогично п.3, только с segments
   
   `Пример API: /segments` 

ЗАДАНИЕ №1
5. `func createSegment(w http.ResponseWriter, r *http.Request)`

   Функция создаёт сегмент и принимает на вход имя сегмента: slug
   
   Если сегмент уже есть в базе, выводит сообщение

   `Пример API: /segment/create?slug=AVITO_PERSENTAGE_30`


ЗАДАНИЕ №2
6. `func deleteSegment(w http.ResponseWriter, r *http.Request)`

   Функция удаляет сегмент и принимает на вход имя сегмента: slug
   
   Если сегмента нет в базе, выдаёт сообщение

   `Пример API: /segment/delete?slug=AVITO_PERSENTAGE_30`

ЗАДАНИЕ №3

7. `func actionUser(w http.ResponseWriter, r *http.Request)`

   Функция, которая добавляет и удаляет сегменты у пользователя, принимает на вход идентификатор пользователя, сегменты, которые нужно добавить и которые нужно удалить
   
   Вспомогательная функция: `func userAddRemove(user_id string, added_segments []string, remove_segments []string)`
   
   Пример API
   ```
   /useraddremove?user_id=1&add=AVITO_TEST1 AVITO_TEST2 AVITO_TEST3&remove=AVITO_PERSENTAGE_30
   ```

   Здесь идентификатор пользователя = 1

   Добавленные сегменты это [AVITO_TEST1, AVITO_TEST2, AVITO_TEST3], которые достаются и преобразуются в дальнейшем в массив

   Аналогично с удалёнными сегментами [AVITO_PERSENTAGE_30]

ЗАДАНИЕ №4

8. `func getUserSegments(w http.ResponseWriter, r *http.Request)`

   Функция получает список сегментов, которые есть у пользователя по user_id

   Если пользователя в базе нет - выводит ошибку

   `Пример API: /segments/user?user_id=1`

Так как все задания на минимум сделал, постарался выполнить ещё 1 дополнительное задание
Взял историю операций ( статистику )

ДОП. ЗАДАНИЕ №1
```
func getUserHistory(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	user_id := params.Get("user_id")
	year := params.Get("year")
	month := params.Get("month")

	// Создаем слайс данных, которые нужно записать в CSV-файл
	query := `SELECT h.user_id, s.slug, h.type_of_operation, h.date_of_operation FROM segment s
   			  JOIN history_operation h ON s.id = h.segment_id 
              WHERE h.user_id = ? AND YEAR(h.date_of_operation) = ? AND MONTH(h.date_of_operation) = ?`
	
	history, err := db.Query(query, user_id, year, month)

	defer func(result *sql.Rows) { _ = result.Close() }(history)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "text/csv; charset=utf-8")
	w.Header().Set(
		"Content-Disposition", fmt.Sprintf("attachment; filename=\"historyUser(%s).csv\"", user_id))

	sqltocsv.Write(w, history)
}
```

Здесь я выбираю переменные айдишника пользователя, год и месяц и по ним возвращаю выборку, которая с помощью text/csv сохраняет её в csv формат с кодировкой utf-8
После записывает данные на страничку и "возвращает нам .csv файл"

`Пример API: /history/user?user_id=1&year=2023&month=8`