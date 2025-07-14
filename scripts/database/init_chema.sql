USE cinema_booking;

CREATE TABLE IF NOT EXISTS movies
(
    id         VARCHAR(255) CHARSET utf8mb3 NOT NULL PRIMARY KEY,
    title      VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    duration   INT                          NOT NULL,
    created_at TIMESTAMP                    NOT NULL,
    updated_at TIMESTAMP                    NOT NULL
);

CREATE TABLE IF NOT EXISTS seats
(
    id         VARCHAR(255) CHARSET utf8mb3 NOT NULL PRIMARY KEY,
    screen_id  VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    seat_row   VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    number     INT                          NOT NULL,
    class      VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    created_at TIMESTAMP                    NOT NULL,
    updated_at TIMESTAMP                    NOT NULL
);

CREATE TABLE IF NOT EXISTS users
(
    id           VARCHAR(255) CHARSET utf8mb3 NOT NULL PRIMARY KEY,
    username     VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    password     VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    status       VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    address      VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    phone_number VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    created_at   TIMESTAMP                    NOT NULL,
    updated_at   TIMESTAMP                    NOT NULL,
    CONSTRAINT users_username_uniq UNIQUE (username)
);


CREATE TABLE IF NOT EXISTS showtime
(
    id         VARCHAR(255) CHARSET utf8mb3 NOT NULL PRIMARY KEY,
    movie_id   VARCHAR(255) CHARSET utf8mb3 NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
    screen_id  VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    start_time TIMESTAMP                    NOT NULL,
    created_at TIMESTAMP                    NOT NULL,
    updated_at TIMESTAMP                    NOT NULL
);

CREATE TABLE IF NOT EXISTS bookings
(
    id          VARCHAR(255) CHARSET utf8mb3 NOT NULL PRIMARY KEY,
    showtime_id VARCHAR(255) CHARSET utf8mb3 NOT NULL REFERENCES showtime (id) ON DELETE CASCADE,
    user_id     VARCHAR(255) CHARSET utf8mb3 NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status      VARCHAR(255) CHARSET utf8mb3 NOT NULL,
    created_at  TIMESTAMP                    NOT NULL,
    updated_at  TIMESTAMP                    NOT NULL
);

CREATE TABLE IF NOT EXISTS booking_seats
(
    id          VARCHAR(255) CHARSET utf8mb3 NOT NULL PRIMARY KEY,
    booking_id  VARCHAR(255) CHARSET utf8mb3 NOT NULL REFERENCES bookings (id) ON DELETE CASCADE,
    showtime_id VARCHAR(255) CHARSET utf8mb3 NOT NULL REFERENCES showtime (id) ON DELETE CASCADE,
    seat_id     VARCHAR(255) CHARSET utf8mb3 NOT NULL REFERENCES seats (id) ON DELETE CASCADE,
    created_at  TIMESTAMP                    NOT NULL,
    updated_at  TIMESTAMP                    NOT NULL,
    CONSTRAINT booking_seats_showtime_seat_uniq UNIQUE (showtime_id, seat_id)
);