CREATE TABLE worships (id varchar(255), title varchar(255), image_url varchar(255), worship_order int);

CREATE TABLE albums (id varchar(255), title varchar(255), image_url varchar(255), album_order int);

CREATE TABLE songs (id varchar(255), title varchar(255), song_number int, file varchar(255), song_order int, album_id int, FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE);

CREATE TABLE sections (id varchar(255), section_name varchar(255), section_order int, worship_id varchar(255));

CREATE TABLE moments (id varchar(255), title varchar(255), song_order int, song_id int, section_id int,  FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE);
