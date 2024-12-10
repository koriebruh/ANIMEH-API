
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
