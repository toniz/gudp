create database accountdb;
use accountdb;
CREATE TABLE `t_user` (
          `user_id` varchar(20) DEFAULT NULL,
          `user_name` varchar(255) DEFAULT NULL,
          `type` varchar(10) DEFAULT NULL
);

insert into t_user(user_id, user_name, type) values("1", "xx","db1");


create database image;
USE image;
CREATE TABLE `t_user` (
          `user_id` varchar(20) DEFAULT NULL,
          `user_name` varchar(255) DEFAULT NULL,
          `type` varchar(10) DEFAULT NULL
);

CREATE TABLE `t_images` (
          `id` int(11) DEFAULT NULL,
          `name` varchar(50) DEFAULT NULL,
          `image` varchar(255) DEFAULT NULL
);

insert into t_user(user_id, user_name, type) values("1", "yy","db2");



