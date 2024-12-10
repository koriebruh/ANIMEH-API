
Untuk menjalankan web server, MySQL, dan Elasticsearch, ikuti langkah-langkah berikut:

1. **Menjalankan Docker Compose**  
   Jalankan perintah berikut untuk mengaktifkan web server, MySQL, dan Elasticsearch menggunakan Docker Compose:

   ```bash
   docker compose up -d
   ```


2. **Tunggu Beberapa Saat**  
   Tunggu sekitar 10 detik untuk memastikan Elasticsearch dan MySQL selesai melakukan proses booting.


3. **Masuk ke dalam Container**  
   Setelah sistem siap, masuk ke dalam container untuk menjalankan unit test yang secara otomatis akan menyisipkan data ke Elasticsearch:

   ```bash
   docker exec -it apinime-api_anime-1 bash
   ```

4. **Jalankan Unit Test**  
   Setelah berada di dalam container, jalankan perintah berikut untuk menjalankan unit test:

   ```bash
   go test -v -run TestInsertDB .
   ```

---

# 📑 API Documentation

### Root Endpoint 🌍
```
GET http://localhost:3000/
Accept: application/json
```

### 🔍 **Auto Complete**
```
GET http://localhost:3000/autocomplete?q=One%20Piece
Accept: application/json
Content-Type: application/json
```

```
GET http://localhost:3000/autocomplete?q=sousou
Accept: application/json
Content-Type: application/json
```

### 🔎 **Search Anime**
```
GET http://localhost:3000/search/anime?name=One+Piece&from=0&size=20&genre=Comedy
Accept: application/json
Content-Type: application/json
```

```
GET http://localhost:3000/search/anime?genre=Comedy
Accept: application/json
Content-Type: application/json
```

```
GET http://localhost:3000/search/anime?name=Sousou%20no%20Frieren
Accept: application/json
Content-Type: application/json
```

```
GET http://localhost:3000/search/anime?min_score=8
Accept: application/json
Content-Type: application/json
```

### ⭐ **Top Anime**
```
GET http://localhost:3000/anime/top?top_year=2022
Accept: application/json
Content-Type: application/json
```

### 📚 **Find Anime by ID**
```
GET http://localhost:3000/anime/100
Accept: application/json
Content-Type: application/json
```

### 🤖 **Recommend Anime by ID**
```
GET http://localhost:3000/anime/100/recommend?page=1
Accept: application/json
Content-Type: application/json
```

### 📝 **Create New User**
```
POST http://localhost:3000/users
Accept: application/json
Content-Type: application/json

{
  "username": "fren",
  "email": "frenm@gmail.com",
  "password": "fren123"
}
```
### 🔐 **Login User**
```
POST http://localhost:3000/users/login
Accept: application/json
Content-Type: application/json

{
  "email": "frenm@gmail.com",
  "password": "fren123"
}
```

### 🔑 **Change Password**
```
POST http://localhost:3000/users/change
Accept: application/json
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "email": "frenm@gmail.com",
  "new_password": "fren47"
}
```

### ✅ **Confirm Password Change**
```
POST http://localhost:3000/users/change-confirm
Accept: application/json
Content-Type: application/json
Authorization: Bearer <your_jwt_token>

{
  "token": "<confirmation_token>"
}
```


### 💖 **Add to Favorites**
```
POST http://localhost:3000/users/fav/771
Accept: application/json
Content-Type: application/json
Authorization: Bearer <your_jwt_token>
```

### ❌ **Remove from Favorites**
```
DELETE http://localhost:3000/users/fav/771
Accept: application/json
Content-Type: application/json
Authorization: Bearer <your_jwt_token>
```

### 📜 **List All Favorites**
```
GET http://localhost:3000/users/fav
Accept: application/json
Content-Type: application/json
Authorization: Bearer <your_jwt_token>
```

---

