CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(50) NOT NULL,
    created TIMESTAMP NOT NULL
);

CREATE TABLE posts
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    message VARCHAR(150) NOT NULL,
    created TIMESTAMP NOT NULL,
    owner_id INT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE comments 
(
    id SERIAL PRIMARY KEY,
    message VARCHAR(100) NOT NULL,
    created TIMESTAMP NOT NULL,
    post_id INT,
    owner_id INT,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE tags 
(
    tags VARCHAR(30) NOT NULL,
    post_id INT,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);