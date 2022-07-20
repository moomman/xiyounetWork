
create type Gender As ENUM ('男','女','未知');
create type Privilege As ENUM ('BAN','管理员','用户');
create type LifeState As ENUM ('单身','热恋','已婚','为人父母','未知');

create table "user"
(
    id        bigserial primary key,
    username  varchar(255) not null unique,
    password  varchar(255) not null,
    avatar    varchar(255) not null,
    lifeState LifeState             default '未知',
    hobby     varchar(255)          default null,
    email     varchar(255) not null unique,
    birthday  timestamptz  not null,
    gender    Gender       not null default '未知',
    signature text         not null default '日常摆烂',
    privilege Privilege    not null default '用户'
);