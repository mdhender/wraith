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

drop table if exists resource_dtl;

drop table if exists colony_factory_group_orders;
drop table if exists colony_factory_group_inventory;
drop table if exists colony_factory_group_stages;
drop table if exists colony_factory_group_units;
drop table if exists colony_factory_group;

drop table if exists colony_mining_group_stages;
drop table if exists colony_mining_group_units;
drop table if exists colony_mining_group;

drop table if exists colony_pay;
drop table if exists colony_population;
drop table if exists colony_rations;
drop table if exists colony_inventory;
drop table if exists colony_hull;
drop table if exists colony_dtl;
drop table if exists colonies;

drop table if exists resources;
drop table if exists planets;

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

drop table if exists units;

create table units
(
    id    varchar(5)  not null,
    name  varchar(25) not null,
    descr varchar(64) not null,
    primary key (id)
);

insert into units (id, name, descr)
values ('ANM', 'anti-missile', 'anti-missile');
insert into units (id, name, descr)
values ('ASC', 'assault-craft', 'assault-craft');
insert into units (id, name, descr)
values ('ASW', 'assault-weapon', 'assault-weapon');
insert into units (id, name, descr)
values ('AUT', 'automation', 'automation');
insert into units (id, name, descr)
values ('CNGD', 'consumer-goods', 'consumer-goods');
insert into units (id, name, descr)
values ('ESH', 'energy-shield', 'energy-shield');
insert into units (id, name, descr)
values ('EWP', 'energy-weapon', 'energy-weapon');
insert into units (id, name, descr)
values ('FCT', 'factory', 'factory');
insert into units (id, name, descr)
values ('FOOD', 'food', 'food');
insert into units (id, name, descr)
values ('FRM', 'farm', 'farm');
insert into units (id, name, descr)
values ('FUEL', 'fuel', 'fuel');
insert into units (id, name, descr)
values ('GOLD', 'gold', 'gold');
insert into units (id, name, descr)
values ('HDR', 'hyper-drive', 'hyper-drive');
insert into units (id, name, descr)
values ('LSP', 'life-support', 'life-support');
insert into units (id, name, descr)
values ('LTSU', 'light-structural', 'light-structural');
insert into units (id, name, descr)
values ('MIN', 'mine', 'mine');
insert into units (id, name, descr)
values ('MLR', 'military-robots', 'military-robots');
insert into units (id, name, descr)
values ('MLSP', 'military-supplies', 'military-supplies');
insert into units (id, name, descr)
values ('MSS', 'missile', 'missile');
insert into units (id, name, descr)
values ('MSL', 'missile-launcher', 'missile-launcher');
insert into units (id, name, descr)
values ('MTLS', 'metallics', 'metallics');
insert into units (id, name, descr)
values ('NMTS', 'non-metallics', 'non-metallics');
insert into units (id, name, descr)
values ('SDR', 'space-drive', 'space-drive');
insert into units (id, name, descr)
values ('SNR', 'sensor', 'sensor');
insert into units (id, name, descr)
values ('SLSU', 'super-light-structural', 'super-light-structural');
insert into units (id, name, descr)
values ('STUN', 'structural', 'structural');
insert into units (id, name, descr)
values ('TPT', 'transport', 'transport');

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

create table planets
(
    id              int         not null auto_increment,
    star_id         int         not null,
    orbit_no        int         not null comment 'range 1..10',
    kind            varchar(13) not null comment 'kind of planet',
    habitability_no int         not null,
    home_planet     varchar(1)  not null,
    primary key (id),
    foreign key (star_id) references stars (id)
        on delete cascade
);

create table resources
(
    id          int         not null auto_increment,
    planet_id   int         not null,
    deposit_no  int         not null,
    kind        varchar(14) not null comment 'natural resource produced from deposit',
    qty_initial int         not null,
    yield_pct   float       not null comment 'range 0..1',
    primary key (id),
    unique key (planet_id, deposit_no),
    foreign key (planet_id) references planets (id)
        on delete cascade
);

create table colonies
(
    id        int         not null auto_increment,
    game_id   int         not null,
    colony_no int         not null,
    kind      varchar(13) not null,
    planet_id int         not null,
    primary key (id),
    unique key (game_id, colony_no),
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
    unit_id         varchar(5) not null,
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
    unit_id         varchar(5) not null,
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
    rebel_pct              float      not null,
    primary key (colony_id, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_rations
(
    colony_id        int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    professional_pct float      not null,
    soldier_pct      float      not null,
    unskilled_pct    float      not null,
    unemployed_pct   float      not null,
    primary key (colony_id, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_pay
(
    colony_id        int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    professional_pct float      not null,
    soldier_pct      float      not null,
    unskilled_pct    float      not null,
    unemployed_pct   float      not null,
    primary key (colony_id, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade
);

create table colony_factory_group
(
    id         int        not null auto_increment,
    colony_id  int        not null,
    group_no   int        not null,
    efftn      varchar(6) not null,
    endtn      varchar(6) not null,
    unit_id    varchar(5) not null comment 'unit being manufactured',
    tech_level int        not null comment 'tech level of unit being manufactured',
    primary key (id),
    unique key (colony_id, group_no, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade,
    foreign key (unit_id) references units (id)
        on delete cascade
);

create table colony_factory_group_units
(
    factory_group_id int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    unit_id          varchar(5) not null,
    tech_level       int        not null,
    qty_operational  int        not null,
    primary key (factory_group_id, efftn, unit_id, tech_level),
    foreign key (factory_group_id) references colony_factory_group (id)
        on delete cascade,
    foreign key (unit_id) references units (id)
        on delete cascade
);

create table colony_factory_group_stages
(
    factory_group_id int        not null,
    turn             varchar(6) not null,
    qty_stage_1      int        not null,
    qty_stage_2      int        not null,
    qty_stage_3      int        not null,
    qty_stage_4      int        not null,
    primary key (factory_group_id, turn),
    foreign key (factory_group_id) references colony_factory_group (id)
        on delete cascade
);

create table colony_mining_group
(
    id          int        not null auto_increment,
    colony_id   int        not null,
    group_no    int        not null,
    efftn       varchar(6) not null,
    endtn       varchar(6) not null,
    resource_id int        not null,
    primary key (id),
    unique key (colony_id, group_no, efftn),
    foreign key (colony_id) references colonies (id)
        on delete cascade,
    foreign key (resource_id) references resources (id)
        on delete cascade
);

create table colony_mining_group_units
(
    mining_group_id int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    unit_id         varchar(5) not null,
    tech_level      int        not null,
    qty_operational int,
    primary key (mining_group_id, efftn),
    foreign key (mining_group_id) references colony_mining_group (id)
        on delete cascade,
    foreign key (unit_id) references units (id)
        on delete cascade
);

create table colony_mining_group_stages
(
    mining_group_id int        not null,
    turn            varchar(6) not null,
    qty_stage_1     int        not null,
    qty_stage_2     int        not null,
    qty_stage_3     int        not null,
    qty_stage_4     int        not null,
    primary key (mining_group_id, turn),
    foreign key (mining_group_id) references colony_mining_group (id)
        on delete cascade
);

create table resource_dtl
(
    resource_id   int        not null,
    efftn         varchar(6) not null,
    endtn         varchar(6) not null,
    remaining_qty int        not null,
    controlled_by int comment 'colony controlling the resource deposit',
    primary key (resource_id, efftn),
    foreign key (controlled_by) references colonies (id)
        on delete set null,
    foreign key (resource_id) references resources (id)
        on delete cascade
);


# CREATE TRIGGER ins_users BEFORE INSERT ON users
#     FOR EACH ROW
# BEGIN
#     SET NEW.handle_lower = lower(NEW.handle);
#     SET NEW.email = lower(NEW.email);
# END;

insert into users (hashed_secret)
values ('*nobody*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'nobody', 'nobody'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*sysop*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'sysop', 'sysop'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*batch*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'batch', 'batch'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user01*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user01', 'user01'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user02*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user02', 'user02'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user03*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user03', 'user03'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user04*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user04', 'user04'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user05*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user05', 'user05'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user06*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user06', 'user06'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user07*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user07', 'user07'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user08*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user08', 'user08'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user09*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user09', 'user09'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user10*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user10', 'user10'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user11*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user11', 'user11'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user12*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user12', 'user12'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user13*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user13', 'user13'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user14*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user14', 'user14'
from users
where id = (select max(id) from users);

insert into users (hashed_secret)
values ('*user15*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), 'user15', 'user15'
from users
where id = (select max(id) from users);

