# Описание
(Проект находится на стадии доработки)
Это REST api реализующее каталог автомобилей.
# В данном api реализовано:
1) Создание структуры базы данных(postgres) при запуске api путем миграции.
Реализована с помощью [golang-migrate](https://github.com/golang-migrate/migrate)
Драйвер для работы с postgres [pgx](https://github.com/jackc/pgx)
3) Логирование уровнем info
Реализовано с помощью [logrus](github.com/sirupsen/logrus)
4) Маршрутизация запросов
Реализовано с помощью  [gorilla/mux](https://github.com/gorilla/mux)
5) Добавление записей новый автомобилей путем обращения к сторонему api
6) Выдача существующих данных с пагинацией
7) Фильтрация данных по марке автомобиля, модели автомобиля и году выпуска
8) Удаление записи
9) Обновление записи
# На данный момент нереализовано, но появится:
1) Фильтрация по остальным полям
2) Сортировка по полям
3) Swagger файл
4) Покрытие кода debug логами
# Rest методы
![image](https://github.com/Sereys13/api_catalog_car/assets/134072150/a80c2b68-9333-44ea-9d22-7c46e311543c)
