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

drop table if exists cors_factory_group_stages;
drop table if exists cors_factory_group_units;
drop table if exists cors_factory_group;

drop table if exists cors_farm_group_stages;
drop table if exists cors_farm_group_units;
drop table if exists cors_farm_group;

drop table if exists cors_mining_group_stages;
drop table if exists cors_mining_group_units;
drop table if exists cors_mining_group;

drop table if exists cors_loc;
drop table if exists cors_pay;
drop table if exists cors_population;
drop table if exists cors_rations;
drop table if exists cors_inventory;
drop table if exists cors_hull;
drop table if exists cors_dtl;

drop table if exists cors;

drop table if exists resources;

drop table if exists planet_dtl;
drop table if exists planets;

drop table if exists stars;
drop table if exists systems;

drop table if exists nation_player;

drop table if exists nation_research;
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
    id    int         not null auto_increment,
    code  varchar(5)  not null,
    name  varchar(25) not null,
    descr varchar(64) not null,
    primary key (id)
);


create table users
(
    id            int         not null auto_increment,
    handle        varchar(32) not null comment 'handle forced to lower-case',
    hashed_secret varchar(64) not null,
    primary key (id)
);

create table user_profile
(
    user_id int         not null,
    effdt   datetime    not null,
    enddt   datetime    not null,
    handle  varchar(32) not null comment 'display handle',
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


create table turns
(
    game_id  int        not null,
    no       int        not null,
    year     int        not null,
    quarter  int        not null,
    turn     varchar(6) not null comment 'formatted as yyyy/q',
    start_dt datetime   not null,
    end_dt   datetime   not null,
    primary key (game_id, turn),
    foreign key (game_id) references games (id)
        on delete cascade
);


create table players
(
    id      int not null auto_increment,
    game_id int not null,
    primary key (id),
    foreign key (game_id) references games (id)
        on delete cascade
);

create table player_dtl
(
    player_id     int         not null,
    efftn         varchar(6)  not null,
    endtn         varchar(6)  not null,
    handle        varchar(32) not null comment 'name in the game',
    controlled_by int comment 'user controlling the player',
    subject_of    int comment 'set if player is regent or viceroy',
    primary key (player_id, efftn),
    foreign key (player_id) references players (id)
        on delete cascade,
    foreign key (controlled_by) references users (id)
        on delete set null,
    foreign key (subject_of) references players (id)
        on delete set null
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
    name          varchar(64) not null,
    govt_name     varchar(64) not null,
    govt_kind     varchar(64) not null,
    controlled_by int comment 'player controlling the nation',
    primary key (nation_id, efftn),
    foreign key (nation_id) references nations (id)
        on delete cascade,
    foreign key (controlled_by) references players (id)
        on delete set null
);

create table nation_player
(
    nation_id int not null,
    player_id int not null,
    primary key (nation_id, player_id),
    unique key (player_id),
    foreign key (nation_id) references nations (id)
        on delete cascade,
    foreign key (player_id) references players (id)
        on delete cascade
);

create table nation_research
(
    nation_id            int        not null,
    efftn                varchar(6) not null,
    endtn                varchar(6) not null,
    tech_level           int        not null,
    research_points_pool int        not null,
    primary key (nation_id),
    unique key (nation_id, efftn),
    foreign key (nation_id) references nations (id)
        on delete cascade
);

create table nation_skills
(
    nation_id     int        not null,
    efftn         varchar(6) not null,
    endtn         varchar(6) not null,
    biology       int        not null,
    bureaucracy   int        not null,
    gravitics     int        not null,
    life_support  int        not null,
    manufacturing int        not null,
    military      int        not null,
    mining        int        not null,
    shields       int        not null,
    primary key (nation_id),
    unique key (nation_id, efftn),
    foreign key (nation_id) references nations (id)
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
    sequence  varchar(1) not null comment 'suffix appended to star location',
    kind      varchar(4) not null,
    primary key (id),
    unique key (system_id, sequence),
    foreign key (system_id) references systems (id)
        on delete cascade
);

create table planets
(
    id          int         not null auto_increment,
    star_id     int         not null,
    orbit_no    int         not null comment 'range 1..10',
    kind        varchar(13) not null comment 'kind of planet',
    home_planet varchar(1)  not null,
    primary key (id),
    foreign key (star_id) references stars (id)
        on delete cascade
);

create table planet_dtl
(
    planet_id       int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    controlled_by   int comment 'nation controlling planet',
    habitability_no int        not null,
    primary key (planet_id, efftn),
    foreign key (planet_id) references planets (id)
        on delete cascade,
    foreign key (controlled_by) references nations (id)
        on delete set null
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

create table cors
(
    id      int         not null auto_increment,
    game_id int         not null,
    msn     int         not null comment 'unique hull number',
    kind    varchar(13) not null,
    primary key (id),
    unique key (game_id, msn)
) comment 'contains colonies and ships';

create table cors_dtl
(
    cors_id       int         not null,
    efftn         varchar(6)  not null,
    endtn         varchar(6)  not null,
    name          varchar(32) not null comment 'name of colony or ship',
    tech_level    int         not null comment 'tech level of colony or ship',
    controlled_by int comment 'player controlling the colony or ship',
    primary key (cors_id, efftn),
    foreign key (cors_id) references cors (id)
        on delete cascade,
    foreign key (controlled_by) references players (id)
        on delete set null
);

create table cors_loc
(
    cors_id   int        not null,
    efftn     varchar(6) not null,
    endtn     varchar(6) not null,
    planet_id int        not null comment 'location of colony or ship',
    primary key (cors_id, efftn),
    foreign key (cors_id) references cors (id)
        on delete cascade,
    foreign key (planet_id) references planets (id)
        on delete cascade
);

create table cors_hull
(
    cors_id         int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    unit_id         int        not null,
    tech_level      int        not null,
    qty_operational int,
    primary key (cors_id, efftn, unit_id, tech_level),
    foreign key (cors_id) references cors (id)
        on delete cascade,
    foreign key (unit_id) references units (id)
        on delete cascade
) comment 'infrastructure of the colony or ship';

create table cors_inventory
(
    cors_id         int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    unit_id         int        not null,
    tech_level      int        not null,
    qty_operational int        not null,
    qty_stowed      int        not null,
    primary key (cors_id, efftn, unit_id, tech_level),
    foreign key (cors_id) references cors (id)
        on delete cascade,
    foreign key (unit_id) references units (id)
        on delete cascade
) comment 'cargo of the colony or ship';

create table cors_population
(
    cors_id                int        not null,
    efftn                  varchar(6) not null,
    endtn                  varchar(6) not null,
    qty_professional       int        not null,
    qty_soldier            int        not null,
    qty_unskilled          int        not null,
    qty_unemployed         int        not null,
    qty_construction_crews int        not null,
    qty_spy_teams          int        not null,
    rebel_pct              float      not null,
    primary key (cors_id, efftn),
    foreign key (cors_id) references cors (id)
        on delete cascade
);

create table cors_rations
(
    cors_id          int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    professional_pct float      not null,
    soldier_pct      float      not null,
    unskilled_pct    float      not null,
    unemployed_pct   float      not null,
    primary key (cors_id, efftn),
    foreign key (cors_id) references cors (id)
        on delete cascade
);

create table cors_pay
(
    cors_id          int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    professional_pct float      not null,
    soldier_pct      float      not null,
    unskilled_pct    float      not null,
    unemployed_pct   float      not null,
    primary key (cors_id, efftn),
    foreign key (cors_id) references cors (id)
        on delete cascade
);

create table cors_factory_group
(
    id         int        not null auto_increment,
    cors_id    int        not null,
    group_no   int        not null,
    efftn      varchar(6) not null,
    endtn      varchar(6) not null,
    unit_id    int        not null comment 'unit being manufactured',
    tech_level int        not null comment 'tech level of unit being manufactured',
    primary key (id),
    unique key (cors_id, group_no, efftn),
    foreign key (cors_id) references cors (id)
        on delete cascade,
    foreign key (unit_id) references units (id)
        on delete cascade
);

create table cors_factory_group_units
(
    factory_group_id int        not null,
    efftn            varchar(6) not null,
    endtn            varchar(6) not null,
    unit_id          int        not null,
    tech_level       int        not null,
    qty_operational  int        not null,
    primary key (factory_group_id, efftn, unit_id, tech_level),
    foreign key (factory_group_id) references cors_factory_group (id)
        on delete cascade,
    foreign key (unit_id) references units (id)
        on delete cascade
);

create table cors_factory_group_stages
(
    factory_group_id int        not null,
    turn             varchar(6) not null,
    qty_stage_1      int        not null,
    qty_stage_2      int        not null,
    qty_stage_3      int        not null,
    qty_stage_4      int        not null,
    primary key (factory_group_id, turn),
    foreign key (factory_group_id) references cors_factory_group (id)
        on delete cascade
);

create table cors_mining_group
(
    id          int        not null auto_increment,
    cors_id     int        not null,
    group_no    int        not null,
    efftn       varchar(6) not null,
    endtn       varchar(6) not null,
    resource_id int        not null,
    primary key (id),
    unique key (cors_id, group_no, efftn),
    foreign key (cors_id) references cors (id)
        on delete cascade,
    foreign key (resource_id) references resources (id)
        on delete cascade
);

create table cors_mining_group_units
(
    mining_group_id int        not null,
    efftn           varchar(6) not null,
    endtn           varchar(6) not null,
    unit_id         int        not null,
    tech_level      int        not null,
    qty_operational int,
    primary key (mining_group_id, efftn),
    foreign key (mining_group_id) references cors_mining_group (id)
        on delete cascade,
    foreign key (unit_id) references units (id)
        on delete cascade
);

create table cors_mining_group_stages
(
    mining_group_id int        not null,
    turn            varchar(6) not null,
    qty_stage_1     int        not null,
    qty_stage_2     int        not null,
    qty_stage_3     int        not null,
    qty_stage_4     int        not null,
    primary key (mining_group_id, turn),
    foreign key (mining_group_id) references cors_mining_group (id)
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
    foreign key (controlled_by) references cors (id)
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

insert into units (code, name, descr)
values ('ANM', 'anti-missile', 'anti-missile');
insert into units (code, name, descr)
values ('ASC', 'assault-craft', 'assault-craft');
insert into units (code, name, descr)
values ('ASW', 'assault-weapon', 'assault-weapon');
insert into units (code, name, descr)
values ('AUT', 'automation', 'automation');
insert into units (code, name, descr)
values ('CNGD', 'consumer-goods', 'consumer-goods');
insert into units (code, name, descr)
values ('ESH', 'energy-shield', 'energy-shield');
insert into units (code, name, descr)
values ('EWP', 'energy-weapon', 'energy-weapon');
insert into units (code, name, descr)
values ('FCT', 'factory', 'factory');
insert into units (code, name, descr)
values ('FOOD', 'food', 'food');
insert into units (code, name, descr)
values ('FRM', 'farm', 'farm');
insert into units (code, name, descr)
values ('FUEL', 'fuel', 'fuel');
insert into units (code, name, descr)
values ('GOLD', 'gold', 'gold');
insert into units (code, name, descr)
values ('HDR', 'hyper-drive', 'hyper-drive');
insert into units (code, name, descr)
values ('LSP', 'life-support', 'life-support');
insert into units (code, name, descr)
values ('LTSU', 'light-structural', 'light-structural');
insert into units (code, name, descr)
values ('MIN', 'mine', 'mine');
insert into units (code, name, descr)
values ('MLR', 'military-robots', 'military-robots');
insert into units (code, name, descr)
values ('MLSP', 'military-supplies', 'military-supplies');
insert into units (code, name, descr)
values ('MSS', 'missile', 'missile');
insert into units (code, name, descr)
values ('MSL', 'missile-launcher', 'missile-launcher');
insert into units (code, name, descr)
values ('MTLS', 'metallics', 'metallics');
insert into units (code, name, descr)
values ('NMTS', 'non-metallics', 'non-metallics');
insert into units (code, name, descr)
values ('SDR', 'space-drive', 'space-drive');
insert into units (code, name, descr)
values ('SNR', 'sensor', 'sensor');
insert into units (code, name, descr)
values ('SLSU', 'super-light-structural', 'super-light-structural');
insert into units (code, name, descr)
values ('STUN', 'structural', 'structural');
insert into units (code, name, descr)
values ('TPT', 'transport', 'transport');

insert into users (handle, hashed_secret)
values ('nobody', '*nobody*');
insert into users (handle, hashed_secret)
values ('sysop', '*sysop*');
insert into users (handle, hashed_secret)
values ('batch', '*batch*');

insert into user_profile (user_id, effdt, enddt, handle, email)
select id, str_to_date('2022/06/22', '%Y/%m/%d'), str_to_date('2099/12/31', '%Y/%m/%d'), handle, handle
from users;
