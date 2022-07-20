-- name: GetUsers :many
select * from ttms.public."user";

-- name: GetUserByName :one
select * from ttms.public."user"
where username = @username::varchar;

-- name: GetUserById :one
select * from ttms.public."user"
where id = @id::bigint;

-- name: UpdateUserAvatar :exec
update ttms.public."user"
set avatar = @avatar
where id = @id;

-- name: CheckUserRepeat :one
select count(*) from ttmsz.public."user"
where username = @username
or email = @email;

-- name: DeleteUserById :exec
delete from
ttms.public."user"
where id = @id;

-- name: CreateUser :one
insert into ttms.public."user"
(username,
 password,
 avatar,
 email,
 birthday,
 signature,
 privilege)
values (@username,
        @password,
        @avatar,
        @email,
        @birthday,
        @signature,
        @privilege)
returning *;

-- name: DeleteUser :exec
delete from ttms.public."user"
where id = @id::bigint;


-- name: UpdateUser :exec
update ttms.public."user"
set username = @username,
    email = @email,
    birthday = @birthday,
    gender = @gender,
    signature = @signature,
    hobby = @hobby,
    lifestate = @lifestate
where id = @id;

-- name: UpdatePassword :exec
update ttms.public."user"
set password = @password
where email = @email;

-- name: ListUserInfo :many
select id,username,email,privilege
from ttms.public."user"
order by id desc
limit $1 offset $2;

-- name: ListNum :one
select count(*) from ttms.public."user";

-- name: SearchUserByName :many
select id,username,email,privilege
from ttms.public."user"
where username like @username
order by id desc
limit $1 offset $2;

-- name: ListNameNum :one
select count(*) from
ttms.public."user"
where username like @username;

