CREATE TABLE usrs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(16) NOT NULL UNIQUE,
    pass VARCHAR(88) NOT NULL,

    CHECK(LENGTH(name) >= 8)
);

CREATE TABLE ads (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL, 
    image_url TEXT,
    price INT NOT NULL,
    user_id INT REFERENCES usrs (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()::TIMESTAMP,

    CHECK(LENGTH(title) >= 2),
    CHECK(LENGTH(content) >= 2),
    CHECK(price >= 1 and price <= 1000000)
);
