/*
 * wraith - the wraith game engine and server
 * Copyright (c) 2022 Michael D. Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

DROP TABLE IF EXISTS user;
CREATE TABLE user
(
    id     int         not null,
    effdt  datetime    not null,
    enddt  datetime    not null,
    email  varchar(64) not null,
    handle varchar(16) not null
);
ALTER TABLE user
    ADD PRIMARY KEY (id, effdt);

DROP TABLE IF EXISTS user_secret;
CREATE TABLE user_secret
(
    id            int         not null,
    hashed_secret varchar(64) not null
);
ALTER TABLE user_secret
    ADD PRIMARY KEY (id);

# CREATE TRIGGER ins_users BEFORE INSERT ON users
#     FOR EACH ROW
# BEGIN
#     SET NEW.handle_lower = lower(NEW.handle);
#     SET NEW.email = lower(NEW.email);
# END;

DROP TABLE IF EXISTS game;
CREATE TABLE game
(
    id         int         not null,
    effdt      datetime    not null,
    enddt      datetime    not null,
    short_name varchar(8)  not null,
    name       varchar(32) not null
);
ALTER TABLE game
    ADD PRIMARY KEY (id, effdt);

DROP TABLE IF EXISTS game_turn;
CREATE TABLE game_turn
(
    game_id int      not null,
    turn_no int      not null,
    effdt   datetime not null,
    enddt   datetime not null
);
ALTER TABLE game_turn
    ADD PRIMARY KEY (game_id, turn_no);

DROP TABLE IF EXISTS systems;
CREATE TABLE systems
(
    game_id varchar(8),
    id      int PRIMARY KEY AUTO_INCREMENT,
    x       int,
    y       int,
    z       int
);

DROP TABLE IF EXISTS stars;
CREATE TABLE stars
(
    game_id   varchar(8),
    system_id int,
    id        int PRIMARY KEY AUTO_INCREMENT,
    kind      varchar(255)
);

DROP TABLE IF EXISTS orbits;
CREATE TABLE orbits
(
    id        int PRIMARY KEY AUTO_INCREMENT,
    game_id   int,
    system_id int,
    star_id   int,
    orbit_no  int
);

DROP TABLE IF EXISTS planets;
CREATE TABLE planets
(
    id        int PRIMARY KEY AUTO_INCREMENT,
    game_id   int,
    system_id int,
    star_id   int,
    orbit_id  int,
    orbit_no  int
);

DROP TABLE IF EXISTS nations;
CREATE TABLE nations
(
    game_id varchar(8),
    id      int PRIMARY KEY AUTO_INCREMENT,
    name    varchar(255)
);

DROP TABLE IF EXISTS players;
CREATE TABLE players
(
    id        int PRIMARY KEY AUTO_INCREMENT,
    user_id   int,
    game_id   int,
    nation_id int
);

DROP TABLE IF EXISTS colonies;
CREATE TABLE colonies
(
    id            int PRIMARY KEY AUTO_INCREMENT,
    game_id       int,
    system_id     int,
    star_id       int,
    orbit_id      int,
    controlled_by int,
    location      int
);

DROP TABLE IF EXISTS inventory;
CREATE TABLE inventory
(
    colony_id       int,
    unit            varchar(255),
    tech_level      int,
    qty_operational int,
    qty_stowed      int
);

DROP TABLE IF EXISTS units;
CREATE TABLE units
(
    code  varchar(8)  not null,
    effdt varchar(10) not null,
    enddt varchar(10) not null,
    name  varchar(20) not null,
    descr varchar(64) not null
);

ALTER TABLE units
    ADD PRIMARY KEY (code, effdt);