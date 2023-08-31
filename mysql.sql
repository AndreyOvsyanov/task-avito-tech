CREATE TABLE user (
    id INT PRIMARY KEY AUTO_INCREMENT,
    fio VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE segment (
    id INT PRIMARY KEY AUTO_INCREMENT,
    slug VARCHAR(255) UNIQUE
);

CREATE TABLE user_segments (
    user_id INT,
    segment_id INT,
    PRIMARY KEY (user_id, segment_id)
);

CREATE TABLE history_operation (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    segment_id INT,
    type_of_operation VARCHAR(50),
    date_of_operation DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO user(fio) VALUES
("Овсянов Андрей Борисович"),
("Дмитриев Дмитрий Дмитриевич"),
("Андреев Андрей Андреевич");