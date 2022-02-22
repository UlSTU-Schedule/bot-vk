CREATE TABLE IF NOT EXISTS vk_students
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(50)    NOT NULL,
    last_name  VARCHAR(50)    NOT NULL,
    user_id    INTEGER UNIQUE NOT NULL,
    group_name VARCHAR(20)    NOT NULL,
    faculty_id SMALLINT       NOT NULL
);

CREATE TABLE IF NOT EXISTS telegram_students
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(100)   NOT NULL,
    last_name  VARCHAR(100),
    user_id    INTEGER UNIQUE NOT NULL,
    group_name VARCHAR(20)    NOT NULL,
    faculty_id SMALLINT       NOT NULL
);