create table if not exists users(
        id bigint primary key auto_increment not null,
        username varchar(50) not null,
    password_hash varchar(100) not null,
    icon varchar(100),
    unique(id, username)
    )
    engine = InnoDB
    default charset = utf8;

create table if not exists chats(
      id bigint primary key auto_increment not null,
      name varchar(50) not null,
    types varchar(20) not null,
    unique(id)
    )
    engine = InnoDB
    default charset = utf8;

create table if not exists messages(
      id bigint  auto_increment not null,
      chat_id bigint not null ,
      author bigint not null ,
    text text(8191) not null,
      sent_at timestamp default current_timestamp,
    unique(id),
    primary key (id),
    foreign key (chat_id) references chats(id),
    foreign key (author) references users(id)
    )
    engine = InnoDB
    default charset = utf8;

create table if not exists users_relationship(
      id bigint primary key auto_increment not null,
      sender_id bigint not null,
    recipient_id bigint  not null,
      relationship varchar(50) not null,
    unique(id),
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY ( recipient_id) REFERENCES users(id)
    )
    engine = InnoDB
    default charset = utf8;

create table if not exists chat_users(
    id bigint primary key auto_increment  not null,
    chat_id bigint not null,
    user_id bigint not null,
    unique(id),
    foreign key (chat_id) references chats(id),
    foreign key (user_id) references users(id)
    )
    engine = InnoDB
    default charset = utf8;
