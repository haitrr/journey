create database journey;
use journey;
create table users (Id int not null AUTO_INCREMENT primary key, UserName varchar(255), PasswordHash varchar(255) )