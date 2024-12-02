drop database if exists inventario;

CREATE DATABASE IF NOT EXISTS inventario;
use inventario;

create table if not exists privilegio (
    privilegio_id int(1) primary key auto_increment,
    privilegio_nombre varchar(20)
);

insert into privilegio (privilegio_nombre) values ('admin'), ('manager'), ('user');

create table if not exists usuario (
    usuario_id int(10) primary key auto_increment,
    usuario_nombre varchar(20) unique,
    usuario_psswd varchar(20),
    usuario_privilegio int(1),
    usuario_activo boolean,
    foreign key (usuario_privilegio) references privilegio(privilegio_id)
);

INSERT into usuario (usuario_nombre, usuario_psswd, usuario_privilegio, usuario_activo)
VALUES 
('admin', 'admin', 1, true), 
('manager', 'manager', 2, true), 
('user', 'user', 3, true),
('tester', 'tester', 1, true);


create table if not exists producto (
    producto_id int(10) primary key auto_increment,
    producto_nombre varchar(20),
    producto_codigo int(10) unique,
    producto_margen decimal(10,2),
    producto_precio int(10),
    producto_activado boolean
);

insert into producto
(producto_nombre, producto_codigo, producto_margen, producto_precio, producto_activado)
values
("producto1", 111111, 10, 2000, 1),
("producto2", 222222, 11, 2500, 1),
("producto3", 333333, 12, 3000, 1),
("producto4", 444444, 13, 3500, 1),
("producto5", 555555, 14, 4000, 1),
("producto6", 666666, 15, 4500, 1)
;

CREATE table if not exists distribuidor (
  distribuidor_id int(10) primary key auto_increment,
  distribuidor_nombre varchar(20)
);

insert into distribuidor
(distribuidor_nombre)
values
("distribuidor1"),
("distribuidor2"),
("distribuidor3")
;

create table if not exists entrada (
    entrada_id int(10) primary key auto_increment,
    entrada_fecha datetime,
    entrada_usuario int(10),
    entrada_distribuidor int(10),
    foreign key (entrada_usuario) references usuario(usuario_id),
    foreign key (entrada_distribuidor) references distribuidor(distribuidor_id)
);

insert into entrada
(entrada_fecha, entrada_usuario, entrada_distribuidor)
values
('2024-11-28 12:00:00', 1, 1),
('2024-11-28 13:00:00', 1, 1),
('2024-11-29 12:00:00', 2, 2),
('2024-11-29 13:00:00', 2, 2),
('2024-11-30 12:00:00', 3, 3),
('2024-11-30 13:00:00', 3, 3),
('2024-12-01 12:00:00', 3, 3),
('2024-12-01 13:00:00', 3, 3)
;


CREATE table if NOT EXISTS producto_entrada (
    pro_ent_id int(10) primary key auto_increment,
    pro_ent_ent_fk int(10),
    pro_ent_pro_fk int(10),
    pro_ent_cantidad int(10),
    pro_ent_precio int(10),
    foreign key (pro_ent_ent_fk) references entrada(entrada_id),
    foreign key (pro_ent_pro_fk) references producto(producto_id)
);

insert into producto_entrada
(pro_ent_ent_fk, pro_ent_pro_fk, pro_ent_cantidad, pro_ent_precio)
values
(1, 1, 2, 2000),
(1, 1, 2, 2000),
(2, 2, 2, 2500),
(2, 2, 2, 2500),
(3, 3, 2, 3000),
(3, 3, 2, 3000),
(4, 4, 2, 3500),
(4, 4, 2, 3500),
(5, 5, 2, 4000),
(5, 5, 2, 4000),
(6, 6, 2, 4500),
(6, 6, 2, 4500),
(7, 1, 2, 2000),
(7, 1, 2, 2000),
(8, 2, 2, 2500),
(8, 2, 2, 2500)
;

CREATE table if NOT EXISTS salida_tipo (
    salida_tipo_id int(1) primary key auto_increment,
    salida_tipo_nombre varchar(20)
);

insert into salida_tipo
(salida_tipo_id, salida_tipo_nombre) 
values 
(1, 'venta'), 
(2, 'merma');

create TABLE IF NOT EXISTS salida (
    salida_id int(10) primary key auto_increment,
    salida_fecha datetime,
    salida_tipo int(1),
    salida_usuario int(10),
    foreign key (salida_tipo) references salida_tipo(salida_tipo_id),
    foreign key (salida_usuario) references usuario(usuario_id)
);
insert into salida
(salida_fecha, salida_tipo, salida_usuario)
values
('2024-11-28 15:00:00', 1, 1),
('2024-11-28 16:00:00', 1, 1),
('2024-11-29 15:00:00', 1, 2),
('2024-11-29 16:00:00', 1, 2),
('2024-11-30 15:00:00', 1, 3),
('2024-11-30 16:00:00', 2, 3),
('2024-12-01 15:00:00', 2, 3),
('2024-12-01 16:00:00', 2, 3)
;


create table IF NOT EXISTS producto_salida (
    pro_sal_id int(10) primary key auto_increment,
    pro_sal_salida int(10),
    pro_sal_producto int(10),
    pro_sal_cantidad int(10),
    pro_sal_precio int(10),
    foreign key (pro_sal_salida) references salida(salida_id),
    foreign key (pro_sal_producto) references producto(producto_id)
);
insert into producto_salida
(pro_sal_salida, pro_sal_producto, pro_sal_cantidad, pro_sal_precio)
values
(1, 1, 2, 2000),
(1, 2, 2, 2500),
(2, 3, 2, 3000),
(2, 4, 2, 3500),
(3, 5, 2, 4000),
(3, 6, 2, 4500),
(4, 1, 2, 2000),
(4, 2, 2, 2500),
(5, 3, 2, 3000),
(5, 4, 2, 3500),
(6, 5, 2, 4000),
(6, 6, 2, 4500),
(7, 1, 2, 2000),
(7, 2, 2, 2500),
(8, 3, 2, 3000),
(8, 4, 2, 3500)
;