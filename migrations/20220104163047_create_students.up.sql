CREATE TABLE IF NOT EXISTS vk_students
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(50)    NOT NULL,
    last_name  VARCHAR(50)    NOT NULL,
    user_id    INTEGER UNIQUE NOT NULL,
    group_name VARCHAR(20)    NOT NULL,
    faculty_id SMALLINT       NOT NULL
);