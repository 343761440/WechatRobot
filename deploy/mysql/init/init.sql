create database wechat;
use wechat;

ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'miyawaki';
GRANT all privileges on wechat.* to 'sakura'@'%' identified by 'kdf82dhsx';
flush privileges;