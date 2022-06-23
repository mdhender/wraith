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

drop table if exists units;

drop table if exists colony_factory_group_orders;
drop table if exists colony_factory_group_inventory;
drop table if exists colony_factory_group;

drop table if exists colony_mining_group_inventory;
drop table if exists colony_mining_group_deposit;
drop table if exists colony_mining_group;

drop table if exists colony_inventory;
drop table if exists colony_hull;
drop table if exists colony_govt;
drop table if exists colonies;

drop table if exists deposit_dtl;
drop table if exists deposits;
drop table if exists planets;

drop table if exists orbits;
drop table if exists stars;
drop table if exists systems;

drop table if exists players;

drop table if exists governments;
drop table if exists nations;

drop table if exists turns;
drop table if exists games;

drop table if exists profiles;
drop table if exists users;

create table users
(
    id            int         not null auto_increment,
    hashed_secret varchar(64) not null,
    primary key (id)
);

create table profiles
(
    user_id int         not null,
    effdt   datetime    not null,
    enddt   datetime    not null,
    handle  varchar(32) not null,
    email   varchar(64) not null,
    primary key (user_id, effdt),
    foreign key (user_id) references users (id)
        on delete cascade
);

create table games
(
    id         int         not null auto_increment,
    short_name varchar(8)  not null,
    name       varchar(32) not null,
    turn       varchar(6)  not null,
    descr      varchar(256),
    primary key (id),
    unique key (short_name)
);

create table turns
(
    game_id  int        not null,
    turn     varchar(6) not null,
    start_dt datetime   not null,
    end_dt   datetime   not null,
    primary key (game_id, turn),
    foreign key (game_id) references games (id)
        on delete cascade
);

create table nations
(
    id         int          not null auto_increment,
    game_id    int          not null,
    nation_no  int          not null,
    speciality varchar(16)  not null,
    descr      varchar(256) not null,
    primary key (id),
    foreign key (game_id) references games (id)
        on delete cascade,
    unique key (game_id, nation_no)
);

create table players
(
    id        int         not null auto_increment,
    user_id   int         not null,
    nation_id int         not null,
    handle    varchar(32) not null,
    primary key (id),
    unique key (nation_id, handle),
    foreign key (user_id) references users (id)
        on delete cascade,
    foreign key (nation_id) references nations (id)
        on delete cascade
);

create table governments
(
    nation_id int         not null,
    efftn     varchar(6)  not null,
    endtn     varchar(6)  not null,
    govt_name varchar(64) not null,
    govt_kind varchar(64) not null,
    primary key (nation_id, efftn),
    foreign key (nation_id) references nations (id)
        on delete cascade
);

create table systems
(
    id      int not null auto_increment,
    game_id int not null,
    x       int,
    y       int,
    z       int,
    primary key (id),
    unique key (game_id, x, y, z),
    foreign key (game_id) references games (id)
        on delete cascade
);

create table stars
(
    id        int        not null auto_increment,
    system_id int        not null,
    suffix    varchar(1) not null,
    kind      varchar(4) not null,
    primary key (id),
    unique key (system_id, suffix),
    foreign key (system_id) references systems (id)
        on delete cascade
);

create table orbits
(
    id       int         not null auto_increment,
    star_id  int         not null,
    orbit_no int         not null,
    kind     varchar(13) not null,
    primary key (id),
    unique key (star_id, orbit_no),
    foreign key (star_id) references stars (id)
        on delete cascade
);

create table planets
(
    id              int not null auto_increment,
    orbit_id        int not null,
    habitability_no int not null,
    primary key (id),
    foreign key (orbit_id) references orbits (id)
        on delete cascade
);

create table deposits
(
    id          int         not null auto_increment,
    planet_id   int         not null,
    deposit_no  int         not null,
    kind        varchar(14) not null,
    initial_qty int         not null,
    yield_pct   int         not null,
    primary key (id),
    unique key (planet_id, deposit_no),
    foreign key (planet_id) references planets (id)
        on delete cascade
);

create table deposit_dtl
(
    deposit_id    int        not null,
    efftn         varchar(6) not null,
    endtn         varchar(6) not null,
    remaining_qty int        not null,
    primary key (deposit_id, efftn),
    foreign key (deposit_id) references deposits (id)
        on delete cascade
);

create table colonies
(
    id        int         not null auto_increment,
    colony_no int         not null,
    planet_id int         not null,
    kind      varchar(13) not null,
    primary key (id),
    foreign key (planet_id) references planets (id)
        on delete cascade
);

create table colony_govt
(
    colony_id     int        not null,
    efftn         varchar(6) not null,
    endtn         varchar(6) not null,
    controlled_by int,
    primary key (colony_id, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade,
    foreign key (controlled_by) references players (id)
        on delete set null
);


create table colony_hull
(
    colony_id       int        not null,
    unit_id         int        not null,
    tech_level      int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    qty_operational int,
    primary key (colony_id, unit_id, tech_level, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_inventory
(
    colony_id       int        not null,
    unit_id         int        not null,
    tech_level      int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    qty_operational int,
    qty_stowed      int,
    primary key (colony_id, unit_id, tech_level, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_factory_group
(
    id        int        not null auto_increment,
    colony_id int        not null,
    group_no  int        not null,
    efftn     varchar(6) not null,
    endtn     varchar(6) not null,
    primary key (id),
    unique key (colony_id, group_no, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_factory_group_inventory
(
    factory_group_id int        not null,
    unit_id          int        not null,
    tech_level       int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    qty_operational  int,
    primary key (factory_group_id, unit_id, tech_level, efftn),
    foreign key (factory_group_id) references colony_factory_group (id)
        on delete cascade
);

create table colony_factory_group_orders
(
    factory_group_id int        not null,
    unit_id          int        not null,
    tech_level       int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    primary key (factory_group_id, unit_id, tech_level, efftn),
    foreign key (factory_group_id) references colony_factory_group (id)
        on delete cascade
);

create table colony_mining_group
(
    id        int        not null auto_increment,
    colony_id int        not null,
    group_no  int        not null,
    efftn     varchar(6) not null,
    endtn     varchar(6) not null,
    primary key (id),
    unique key (colony_id, group_no, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_mining_group_deposit
(
    mining_group_id int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    deposit_id      int        not null,
    primary key (mining_group_id, efftn),
    foreign key (deposit_id) references deposits (id)
        on delete cascade
);


create table colony_mining_group_inventory
(
    mining_group_id int        not null,
    unit_id         int        not null,
    tech_level      int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    qty_operational int,
    primary key (mining_group_id, unit_id, tech_level, efftn),
    foreign key (mining_group_id) references colony_mining_group (id)
        on delete cascade
);

create table units
(
    id    int         not null auto_increment,
    code  varchar(5)  not null,
    effdt datetime    not null,
    enddt datetime    not null,
    name  varchar(20) not null,
    descr varchar(64) not null,
    primary key (id),
    unique key (code, effdt)
);


# CREATE TRIGGER ins_users BEFORE INSERT ON users
#     FOR EACH ROW
# BEGIN
#     SET NEW.handle_lower = lower(NEW.handle);
#     SET NEW.email = lower(NEW.email);
# END;

