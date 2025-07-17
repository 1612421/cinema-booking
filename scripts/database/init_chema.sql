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
    updated_at TIMESTAMP                    NOT NULL,
    CONSTRAINT users_screen_row_num_uniq UNIQUE (screen_id, seat_row, number)
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

create
    definer = root@`%` procedure seed_seats(IN pScreenID varchar(255), IN pRow int, IN pSeatNum int)
BEGIN
    DECLARE row_idx INT DEFAULT 1;
    DECLARE seat_num INT;
    DECLARE row_letter VARCHAR(5);
    DECLARE letters CHAR(26) DEFAULT 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

    WHILE row_idx <= pRow
        DO
            SET seat_num = 1;

            -- Sinh tên hàng: A, B, ..., Z, AA, AB,...
            SET row_letter = IF(row_idx <= 26, SUBSTRING(letters, row_idx, 1), CONCAT(
                    SUBSTRING(letters, FLOOR((row_idx - 1) / 26), 1),
                    SUBSTRING(letters, MOD((row_idx - 1), 26) + 1, 1)
                                                                               ));

            WHILE seat_num <= pSeatNum
                DO
                    INSERT IGNORE INTO seats (id, screen_id, seat_row, number, class, created_at, updated_at)
                    VALUES (UUID(),
                            pScreenID,
                            row_letter,
                            seat_num,
                            'standard',
                            NOW(),
                            NOW());

                    SET seat_num = seat_num + 1;
                END WHILE;

            SET row_idx = row_idx + 1;
        END WHILE;
END;

