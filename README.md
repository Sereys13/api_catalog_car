# Описание
Это REST api реализующее каталог автомобилей с использованием базы данных Postgres.  
# Задачи
Данное api решает следующие задачи:
1) Получение данных с фильтрацией по следующим полям каталога: марка автомобиля, модель автомобиля, владелец автомобиля, год выпуска автомобиля
2) Удаление записи из каталога по идентификатору
3) Изменение одного или нескольких полей каталога по идентификатору
4) Добавления новых автомобилей с помощью стороннего api
# В данном api реализовано:
1) Создание структуры базы данных(postgres) при запуске api путем миграции.  
Миграция реализована с помощью пакета [golang-migrate](https://github.com/golang-migrate/migrate)  
Драйвер для работы с postgres [pgx](https://github.com/jackc/pgx)
2) Логирование   
Реализовано с помощью пакета [logrus](github.com/sirupsen/logrus)
3) Маршрутизация запросов  
Реализовано с помощью  [gorilla/mux](https://github.com/gorilla/mux)
4) Чтение конфигурации api из config.env  
Работа с env файлами реализована с помощью пакета [joho/godotenv](https://github.com/joho/godotenv)
5) Добавление записей новый автомобилей путем обращения к стороннему api, описанного следующим сваггером

`openapi: 3.0.3`  
`info:`  
`  title: Car info`  
`  version: 0.0.1`   
`paths:`  
`  /info:`  
`    get:`  
`      parameters:`  
`        - name: reqNum`  
`          in: query`  
`          required: true`  
`          schema:`  
`            type: string`  
`      responses:`  
`        '200':`  
`          description: Ok`  
`          content:`  
`            application/json:`  
`              schema:`  
`                $ref: '#/components/schemas/Car'`  
`        '500':`  
`          description: Internal server error`  
`components:`  
`  schemas:`  
`    Car:`  
`      required:`  
`        - reqNum`  
`        - mark`  
`        - model`  
`        - owner`  
`      type: object`  
`      properties:`  
`        reqNum:`  
`          type: string`  
`          example: X123XX150`  
`        mark:`  
`          type: string`  
`          example: Lada`  
`        model:`  
`          type: string`  
`          example: Vesta`  
`        year:`  
`          type: integer`  
`          example: 2002`  
`         owner: `  
`          $ref: '#/components/schemas/People'`  
`    People:`  
`      required:`  
`        - name`  
`        - surname`  
`      type: object`  
`      properties:`  
`        name:`  
`          type: string`  
`        surname:`  
`          type: string`  
`        patronymic:`  
`          type: string`  
url данного api прописывается в файле config.env в параметре <code>URL_API_CAR_INFO</code>

6) Выдача существующих данных с пагинацией
7) Фильтрация данных по марке автомобиля, модели автомобиля, владельцу и году выпуска
8) Удаление записи по идентификатору
9) Обновление записи по идентификатору
# Rest методы
## GET
* `url/catalog` 
Выдача 10 строк данных каталога и индекса объекта последней строке (для пагинации)  в формате JSON
[Пример выдачи](https://github.com/Sereys13/api_catalog_car/blob/main/примерGetCatalog.json)
* `url/catalog?p=id  `
p параметр для выдачи следующих строк данных каталога
* `url/catalog?count=int `
count параметр для установки количества строк выдачи каталога(по умолчанию количество равно десяти)
* `url/catalog?filtr= `
filtr параметр, указывающий наличие параметров для фильтрации выдачи каталога
* `url/catalog?filtr=true$brand=[id...]`
brand параметр фильтрации выдачи товаров по определенным маркам автомобилей
* `url/catalog?filtr=true$model=[id...]`
model параметр фильтрации выдачи товаров по определенным моделям автомобилей
* `url/catalog?filtr=true$holder=[id...]`
holder параметр фильтрации выдачи товаров по определенным владельцам автомобилей
* `url/catalog?filtr=true$year=ravno&year=[id...]`
year параметр фильтрации выдачи товаров по определенным годам выпуска
* `url/catalog?filtr=true$year=ot&year=2010`
year параметр фильтрации выдачи товаров от(включительно) определенного года выпуска до текущего года(включительного)
* `url/catalog?filtr=true$year=ot&year=2010&year=2020`
year параметр фильтрации выдачи товаров от(включительно) определенного года выпуска до определенного года(включительно)
* `url/catalog?filtr=true$year=do&year=2010`
year параметр фильтрации выдачи товаров до определенного года(включительного)

* `url/filters` 
Выдача доступных фильтров каталога в формате JSON
[Пример выдачи](https://github.com/Sereys13/api_catalog_car/blob/main/примерGetFilters.json)

## POST
* `url/catalog`  
В теле запроса ожидается json в формате  
![image](https://github.com/Sereys13/api_catalog_car/assets/134072150/175c31db-008f-4e25-a46f-1564eb556e89)  
Добавление записи в каталог путем обращения к сторонему api
На сервере производится валидация данных перед обращением к стороннему api и валидация полученных данных от стороннего api
## DELETE
* `url/catalog/[id:[0-9]+]`  
Удаление записи в каталоге по идентификатору
## PUT
* `url/catalog/[id:[0-9]+]`  
Обновление записи в каталоге по идентификатору
В теле запроса ожидается json в следующих форматах:
- Для обновления номера автомобиля нужно указать следующее поле  
  ![image](https://github.com/user-attachments/assets/2fc6570c-4c20-4148-83f7-7079a7298b6d)
- Для обновления года выпуска автомобиля нужно указать следующее поле  
  ![image](https://github.com/user-attachments/assets/b4fdf17e-24b5-458b-a1f1-209d14368bae)
- Для обновления марки автомобиля нужно указать следующие поля  
  ![image](https://github.com/user-attachments/assets/93f9f4e0-a023-498b-8ecf-d8aedb7d0665)
- Для обновления модели автомобиля нужно указать следующие поля  
  ![image](https://github.com/user-attachments/assets/8d285b5e-1d8f-479b-8413-2e99481fcb99)
- Для обновления данных о владельца нужно указать следующие поля  
  ![image](https://github.com/user-attachments/assets/cfabdcdb-b410-4150-9c45-8a9fae5dcdfc)  
  или  
  ![image](https://github.com/user-attachments/assets/c8e1ed6a-d6ec-4558-9456-1e9b3623ed0e)
- Для обновления всех полей JSON выглядит следующим образом:
  ![image](https://github.com/user-attachments/assets/b5e8d19c-7575-4ba5-80cd-6e6dddf13d8a)










