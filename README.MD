# shoppinglist-api

## 📖Описание
ShoppingList API — это сервис для совместного управления списками покупок. Пользователи могут объединяться в группы и совместно редактировать общий список продуктов в реальном времени.

Данный сервис был разработан в рамках университетского задания и имеет дополнительный компонент: [android-app](https://github.com/Clonsaldafon/ShoppingListApp)

## ⚙️Установка
1️⃣ Изменить название файла
```bash
.env.example -> .env
```
2️⃣ Поменять конфигурацию в зависимосит от своих потребностей
```bash
ENV=dev #can be prod

HTTP_HOST=app
HTTP_PORT=9090

POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DB=yourdb
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres

JWT_SECRET=yoursecret
JWT_EXPIRE=1h
```
3️⃣ Запустить сервис
```bash
docker-compose up --build
```

## 📚 Документация
Если ENV=dev, то документация и спецификация будут доступны [тут](http://localhost:9090/swagger/index.html)
