CREATE TABLE users
(
    ID       serial PRIMARY KEY,
    username varchar(255) not null,
    email    varchar(255) not null unique,
    platform varchar(255) not null
);

CREATE TABLE tokens
(
    ID serial PRIMARY KEY,
    userID serial REFERENCES users(ID) ON DELETE CASCADE not null,
    tokenValue varchar not null unique,
    revoked varchar(255) not null
);

CREATE TABLE users_spotify
(
    ID            serial PRIMARY KEY,
    userID serial REFERENCES users(ID) ON DELETE CASCADE not null,
    spotifyUri varchar not null unique,
    accessToken varchar not null unique,
    refreshToken varchar not null unique
);

CREATE TABLE user_youtubemusic
(
    ID            serial PRIMARY KEY,
    userID serial REFERENCES users(ID) ON DELETE CASCADE not null,
    accessToken varchar not null unique
);

CREATE TABLE groups
(
    ID serial PRIMARY KEY,
    name varchar(255) not null,
    description varchar(255) not null
);

CREATE TABLE roles
(
    ID serial PRIMARY KEY,
    name varchar(255) not null
);

CREATE TABLE user_group
(
    ID serial PRIMARY KEY,
    userID serial REFERENCES users(ID) ON DELETE CASCADE not null,
    groupID serial REFERENCES groups(ID) ON DELETE CASCADE not null,
    roleID serial REFERENCES roles(ID) ON DELETE CASCADE not null,
    playListID varchar(255) not null unique
);

CREATE TABLE playlists
(
    ID serial PRIMARY KEY,
    name varchar(255) not null,
    description varchar(255) not null,
    groupID serial REFERENCES groups(ID) ON DELETE CASCADE not null
);

CREATE TABLE tracks
(
    ID serial PRIMARY KEY,
    spotifyUri varchar(255) not null unique,
    youTubeMusicID varchar(255) not null unique,
    name varchar(255) not null
);

CREATE TABLE playlist_track
(
    ID serial PRIMARY KEY,
    trackID serial REFERENCES tracks(ID) ON DELETE CASCADE not null,
    playListID serial REFERENCES playlists(ID) ON DELETE CASCADE not null
);

CREATE TABLE youtube_music_tracks
(
    ID serial PRIMARY KEY,
    userID serial REFERENCES users(ID) ON DELETE CASCADE not null,
    trackID serial REFERENCES tracks(ID) ON DELETE CASCADE not null,
    playListID serial REFERENCES playlists(ID) ON DELETE CASCADE not null,
    playListYouTubeMusicID varchar unique
);



