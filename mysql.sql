CREATE TABLE user (
    id int primary key auto_increment,
    fio varchar(255),
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp
);

CREATE TABLE segment (
    id int primary key auto_increment,
    slug varchar(255) unique
);

CREATE TABLE user_segments (
    user_id int,
    segment_id int,
    primary key (user_id, segment_id),
    foreign key (user_id) references user(id),
    foreign key (segment_id) references segment(id)
);

CREATE TABLE history_operation (
    id int primary key auto_increment,
    user_id int,
    segment_id int,
    type_of_operation varchar(50),
    date_of_operation datetime default current_timestamp,
    foreign key (user_id) references user(id),
    foreign key (segment_id) references segment(id)
);

INSERT INTO user(fio) VALUES
('Овсянов Андрей Борисович'),
('Дмитриев Дмитрий Дмитриевич'),
('Андреев Андрей Андреевич');