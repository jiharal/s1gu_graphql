CREATE table user (
    id serial primary key,
    username string,
    password string,
    email string,
    phone_number string
);

CREATE TABLE question (
    id serial PRIMARY KEY,
    iduser int not null REFERENCES user(id),
    question string,
    answer string,
    action string,
    status bool default 'false',
    created_at TIMESTAMP default now(),
    updated_at TIMESTAMP default now(),
    INDEX(iduser)
);

