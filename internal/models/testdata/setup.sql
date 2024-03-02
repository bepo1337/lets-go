create table snippets (
    id integer not null primary key auto_increment,
    title varchar(100) not null,
    content text not null,
    created DATETIME not null,
    expires DATETIME not null
);

create index idx_snippets_created on snippets(created);
create table users (
    id integer not null primary key auto_increment,
    name varchar(255) not null,
    email varchar(255) not null unique,
    hashed_pw char(60) not null,
    created DATETIME not null
);

insert into users(name, email, hashed_pw, created) values(
                                                          'Alice Jones',
                                                          'alice@email',
                                                          '$2a$10$IZ9.M6iKpc48Nj.MW2752O0mz96178lqllvIZBrHKeAlWpoknfjw6',
                                                          '2022-01-01 10:00:00'
                                                         );

