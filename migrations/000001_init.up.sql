CREATE TABLE IF NOT EXISTS users (
    id INT GENERATED ALWAYS AS IDENTITY UNIQUE,
    telegram_user_id INT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS mangas (
    id INT GENERATED ALWAYS AS IDENTITY UNIQUE,
    name varchar(255) NOT NULL UNIQUE,
    link varchar(160) NOT NULL,
    last_chapter varchar(100),
    last_chapter_url varchar(160)
);

CREATE TABLE IF NOT EXISTS users_mangas (
    user_id int,
    manga_id int,
    CONSTRAINT fk_manga_id FOREIGN KEY (manga_id)
        REFERENCES mangas(id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id)
        REFERENCES users(id)
);