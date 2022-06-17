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

DROP TABLE IF EXISTS `inventory`;
DROP TABLE IF EXISTS `colonies`;
DROP TABLE IF EXISTS `planets`;
DROP TABLE IF EXISTS `orbits`;
DROP TABLE IF EXISTS `stars`;
DROP TABLE IF EXISTS `systems`;
DROP TABLE IF EXISTS `players`;
DROP TABLE IF EXISTS `nations`;
DROP TABLE IF EXISTS `games`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `units`;

CREATE TABLE `users`
(
    `id`            int PRIMARY KEY AUTO_INCREMENT,
    `email`         varchar(255),
    `handle`        varchar(255),
    `hashed_secret` varchar(255)
);

ALTER TABLE `users`
    ADD UNIQUE (`email`);

ALTER TABLE `users`
    ADD UNIQUE (`handle`);

CREATE TABLE `games`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `name`       varchar(255),
    `short_name` varchar(8),
    `turn_no`    int
);

CREATE TABLE `nations`
(
    `id`      int PRIMARY KEY AUTO_INCREMENT,
    `game_id` int,
    `name`    varchar(255)
);

CREATE TABLE `players`
(
    `id`        int PRIMARY KEY AUTO_INCREMENT,
    `user_id`   int,
    `game_id`   int,
    `nation_id` int
);

CREATE TABLE `systems`
(
    `id`      int PRIMARY KEY AUTO_INCREMENT,
    `game_id` int,
    `x`       int,
    `y`       int,
    `z`       int
);

CREATE TABLE `stars`
(
    `id`        int PRIMARY KEY AUTO_INCREMENT,
    `game_id`   int,
    `system_id` int,
    `kind`      varchar(255)
);

CREATE TABLE `orbits`
(
    `id`        int PRIMARY KEY AUTO_INCREMENT,
    `game_id`   int,
    `system_id` int,
    `star_id`   int,
    `orbit_no`  int
);

CREATE TABLE `planets`
(
    `id`        int PRIMARY KEY AUTO_INCREMENT,
    `game_id`   int,
    `system_id` int,
    `star_id`   int,
    `orbit_id`  int,
    `orbit_no`  int
);

CREATE TABLE `colonies`
(
    `id`            int PRIMARY KEY AUTO_INCREMENT,
    `game_id`       int,
    `system_id`     int,
    `star_id`       int,
    `orbit_id`      int,
    `controlled_by` int,
    `location`      int
);

CREATE TABLE `inventory`
(
    `colony_id`       int,
    `unit`            varchar(255),
    `tech_level`      int,
    `qty_operational` int,
    `qty_stowed`      int
);

CREATE TABLE `units`
(
    `code` varchar(255) PRIMARY KEY,
    `name` varchar(255)
);

ALTER TABLE `nations`
    ADD FOREIGN KEY (`game_id`) REFERENCES `games` (`id`);

ALTER TABLE `players`
    ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `players`
    ADD FOREIGN KEY (`game_id`) REFERENCES `games` (`id`);

ALTER TABLE `players`
    ADD FOREIGN KEY (`nation_id`) REFERENCES `nations` (`id`);

ALTER TABLE `systems`
    ADD FOREIGN KEY (`game_id`) REFERENCES `games` (`id`);

ALTER TABLE `stars`
    ADD FOREIGN KEY (`game_id`) REFERENCES `games` (`id`);

ALTER TABLE `stars`
    ADD FOREIGN KEY (`system_id`) REFERENCES `systems` (`id`);

ALTER TABLE `orbits`
    ADD FOREIGN KEY (`game_id`) REFERENCES `games` (`id`);

ALTER TABLE `orbits`
    ADD FOREIGN KEY (`system_id`) REFERENCES `systems` (`id`);

ALTER TABLE `orbits`
    ADD FOREIGN KEY (`star_id`) REFERENCES `stars` (`id`);

ALTER TABLE `planets`
    ADD FOREIGN KEY (`game_id`) REFERENCES `games` (`id`);

ALTER TABLE `planets`
    ADD FOREIGN KEY (`system_id`) REFERENCES `systems` (`id`);

ALTER TABLE `planets`
    ADD FOREIGN KEY (`star_id`) REFERENCES `stars` (`id`);

ALTER TABLE `planets`
    ADD FOREIGN KEY (`orbit_id`) REFERENCES `orbits` (`id`);

ALTER TABLE `colonies`
    ADD FOREIGN KEY (`game_id`) REFERENCES `games` (`id`);

ALTER TABLE `colonies`
    ADD FOREIGN KEY (`system_id`) REFERENCES `systems` (`id`);

ALTER TABLE `colonies`
    ADD FOREIGN KEY (`star_id`) REFERENCES `stars` (`id`);

ALTER TABLE `colonies`
    ADD FOREIGN KEY (`orbit_id`) REFERENCES `orbits` (`id`);

ALTER TABLE `colonies`
    ADD FOREIGN KEY (`controlled_by`) REFERENCES `nations` (`id`);

ALTER TABLE `colonies`
    ADD FOREIGN KEY (`location`) REFERENCES `planets` (`id`);

ALTER TABLE `inventory`
    ADD FOREIGN KEY (`colony_id`) REFERENCES `colonies` (`id`);

ALTER TABLE `inventory`
    ADD FOREIGN KEY (`unit`) REFERENCES `units` (`code`);
