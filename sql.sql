CREATE DATABASE api_anime;

USE api_anime;

SHOW TABLES ;

SELECT * FROM users;
SELECT * FROM favorites;

select * from favorites where user_id =2;

SELECT anime_id FROM favorites WHERE user_id = 3;


# DROP TABLE users;
# DROP TABLE favorites;
