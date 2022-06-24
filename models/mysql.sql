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

drop table if exists deposit_dtl;

drop table if exists colony_factory_group_orders;
drop table if exists colony_factory_group_inventory;
drop table if exists colony_factory_group;

drop table if exists colony_mining_group_inventory;
drop table if exists colony_mining_group_deposit;
drop table if exists colony_mining_group;

drop table if exists colony_pay;
drop table if exists colony_population;
drop table if exists colony_rations;
drop table if exists colony_inventory;
drop table if exists colony_hull;
drop table if exists colony_dtl;
drop table if exists colonies;

drop table if exists deposits;
drop table if exists planets;

drop table if exists orbits;
drop table if exists stars;
drop table if exists systems;

drop table if exists nation_skills;
drop table if exists nation_dtl;
drop table if exists nations;

drop table if exists player_dtl;
drop table if exists players;

drop table if exists turns;
drop table if exists games;

drop table if exists user_profile;
drop table if exists users;

create table users
(
    id            int         not null auto_increment,
    hashed_secret varchar(64) not null,
    primary key (id)
);

create table user_profile
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
    id           int         not null auto_increment,
    short_name   varchar(8)  not null comment 'code showed on report',
    name         varchar(32) not null comment 'full name of game',
    current_turn varchar(6)  not null,
    descr        varchar(256) comment 'details about game',
    primary key (id),
    unique key (short_name)
);

create table players
(
    id            int not null auto_increment,
    game_id       int not null,
    controlled_by int comment 'user controlling the player',
    subject_of    int comment 'set if player is regent or viceroy',
    primary key (id),
    foreign key (game_id) references games (id)
        on delete cascade,
    foreign key (controlled_by) references users (id)
        on delete set null,
    foreign key (subject_of) references players (id)
        on delete set null
);

create table turns
(
    game_id  int        not null,
    turn     varchar(6) not null comment 'formatted as yyyy/q',
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

create table nation_dtl
(
    nation_id     int         not null,
    efftn         varchar(6)  not null,
    endtn         varchar(6)  not null,
    controlled_by int comment 'player controlling the nation',
    govt_name     varchar(64) not null,
    govt_kind     varchar(64) not null,
    primary key (nation_id, efftn),
    foreign key (nation_id) references nations (id)
        on delete cascade,
    foreign key (controlled_by) references players (id)
        on delete set null
);

create table nation_skills
(
    nation_id            int        not null,
    efftn                varchar(6) not null,
    endtn                varchar(6) not null,
    tech_level           int        not null,
    research_points_pool int        not null,
    biology              int        not null comment 'not used currently',
    bureaucracy          int        not null comment 'not used currently',
    gravitics            int        not null comment 'not used currently',
    life_support         int        not null comment 'not used currently',
    manufacturing        int        not null comment 'not used currently',
    military             int        not null comment 'not used currently',
    mining               int        not null comment 'not used currently',
    shields              int        not null comment 'not used currently',
    primary key (nation_id),
    unique key (nation_id, efftn),
    foreign key (nation_id) references nations (id)
        on delete cascade
);

create table player_dtl
(
    player_id int         not null,
    efftn     varchar(6)  not null,
    endtn     varchar(6)  not null,
    handle    varchar(32) not null comment 'name in the game',
    primary key (player_id, efftn),
    foreign key (player_id) references players (id)
        on delete cascade
);

create table systems
(
    id        int not null auto_increment,
    game_id   int not null,
    x         int,
    y         int,
    z         int,
    qty_stars int comment 'number of stars in system',
    primary key (id),
    unique key (game_id, x, y, z),
    foreign key (game_id) references games (id)
        on delete cascade
);

create table stars
(
    id        int        not null auto_increment,
    system_id int        not null,
    suffix    varchar(1) not null comment 'suffix appended to star location',
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
    orbit_no int         not null comment 'range 1..10',
    kind     varchar(13) not null comment 'kind of planet in this orbit',
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
    kind        varchar(14) not null comment 'natural resource produced from deposit',
    qty_initial int         not null,
    yield_pct   int         not null,
    primary key (id),
    unique key (planet_id, deposit_no),
    foreign key (planet_id) references planets (id)
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

create table colony_dtl
(
    colony_id     int         not null,
    efftn         varchar(6)  not null,
    endtn         varchar(6)  not null,
    name          varchar(32) not null comment 'name of colony',
    controlled_by int comment 'player controlling the colony',
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
) comment 'infrastructure of the colony';

create table colony_inventory
(
    colony_id       int        not null,
    unit_id         int        not null,
    tech_level      int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    qty_operational int        not null,
    qty_stowed      int        not null,
    primary key (colony_id, unit_id, tech_level, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
) comment 'cargo of the colony';

create table colony_population
(
    colony_id              int        not null,
    efftn                  varchar(6) not null,
    endtn                  varchar(6) not null,
    qty_professional       int        not null,
    qty_soldier            int        not null,
    qty_unskilled          int        not null,
    qty_unemployed         int        not null,
    qty_construction_crews int        not null,
    qty_spy_teams          int        not null,
    rebel_pct              int        not null,
    primary key (colony_id, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_rations
(
    colony_id        int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    qty_professional int        not null,
    qty_soldier      int        not null,
    qty_unskilled    int        not null,
    qty_unemployed   int        not null,
    primary key (colony_id, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_pay
(
    colony_id        int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    qty_professional int        not null,
    qty_soldier      int        not null,
    qty_unskilled    int        not null,
    qty_unemployed   int        not null,
    primary key (colony_id, efftn),
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
    qty_operational  int        not null,
    primary key (factory_group_id, unit_id, tech_level, efftn),
    foreign key (factory_group_id) references colony_factory_group (id)
        on delete cascade
);

create table colony_factory_group_orders
(
    factory_group_id int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    unit_id          int        not null comment 'unit being manufactured',
    tech_level       int        not null comment 'tech level of unit being manufactured',
    primary key (factory_group_id, efftn),
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

create table deposit_dtl
(
    deposit_id    int        not null,
    efftn         varchar(6) not null,
    endtn         varchar(6) not null,
    remaining_qty int        not null,
    controlled_by int comment 'colony controlling the deposit',
    primary key (deposit_id, efftn),
    foreign key (controlled_by) references colonies (id)
        on delete set null,
    foreign key (deposit_id) references deposits (id)
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

insert into users (hashed_secret)
values ('*login-not-permitted*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'sysop', 'sysop'
from users
where id = 1;

insert into users (hashed_secret)
values ('*login-not-permitted*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'batch', 'batch'
from users
where id = 2;

