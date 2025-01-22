create table users (
                       ID serial primary key,
                       login varchar(255) not null,
                       pass varchar(255) not null,
                       email varchar(255) not null,
                       role varchar(255),
                       pfp varchar(255)
);

create table genre (
                       id serial primary key,
                       title text not null
);

create table albums (
                        ID serial primary key,
                        type_of varchar(100),
                        title text,
                        cover varchar (255),
                        author_id int references users(id) on delete cascade,
                        label varchar(255),
                        publication_date timestamp default current_timestamp
);

create table tracks (
                        ID serial primary key,
                        duration int,
                        tite text,
                        album_id int references albums(id) on delete cascade,
                        genre_id int references genre(id) on delete cascade,
                        tracklink varchar(255),
                        listen_count int
);

create table playlists (
                           ID serial primary key,
                           title text not null,
                           creation_date timestamp default current_timestamp,
                           author_id int references users(id) on delete cascade,
                           cover varchar(255),
                           description text,
                           isPublic bool
);

create table playlist_tracks (
                                 ID serial primary key,
                                 track_id int references tracks(id) on delete cascade,
                                 add_date timestamp default current_timestamp,
                                 playlist_id int references playlists(id) on delete cascade
);

create table playlist_folder (
                                 ID serial primary key,
                                 author_id int references users(id) on delete cascade,
                                 playlist_id int references playlists(id) on delete cascade,
                                 title text not null,
                                 color varchar(25)
);

create table users_artists (
                               ID serial primary key,
                               isLike bool,
                               isDisliked bool,
                               artist_id int references users(id) on delete cascade,
                               user_id int references users(id) on delete cascade
);
