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
VALUES ('admin', 'admin', 1, true), ('manager', 'manager', 2, true), ('user', 'user', 3, true);


create table if not exists producto (
    producto_id int(10) primary key auto_increment,
    producto_nombre varchar(20),
    producto_codigo int(10) unique,
    producto_margen decimal(10,2),
    producto_precio int(10),
    producto_activado boolean
);

CREATE table if not exists distribuidor (
  distribuidor_id int(10) primary key auto_increment,
  distribuidor_nombre varchar(20)
);

create table if not exists entrada (
    entrada_id int(10) primary key auto_increment,
    entrada_fecha datetime,
    entrada_usuario int(10),
    entrada_distribuidor int(10),
    foreign key (entrada_usuario) references usuario(usuario_id),
    foreign key (entrada_distribuidor) references distribuidor(distribuidor_id)
);

CREATE table if NOT EXISTS producto_entrada (
    pro_ent_id int(10) primary key auto_increment,
    pro_ent_ent_fk int(10),
    pro_ent_pro_fk int(10),
    pro_ent_cantidad int(10),
    pro_ent_precio int(10),
    foreign key (pro_ent_ent_fk) references entrada(entrada_id),
    foreign key (pro_ent_pro_fk) references producto(producto_id)
);

CREATE table if NOT EXISTS salida_tipo (
    salida_tipo_id int(1) primary key auto_increment,
    salida_tipo_nombre varchar(20)
);

insert into salida_tipo (salida_tipo_id, salida_tipo_nombre) values (1, 'venta'), (2, 'merma');

create TABLE IF NOT EXISTS salida (
    salida_id int(10) primary key auto_increment,
    salida_fecha datetime,
    salida_tipo int(1),
    salida_usuario int(10),
    foreign key (salida_tipo) references salida_tipo(salida_tipo_id),
    foreign key (salida_usuario) references usuario(usuario_id)
);

create table IF NOT EXISTS producto_salida (
    pro_sal_id int(10) primary key auto_increment,
    pro_sal_salida int(10),
    pro_sal_producto int(10),
    pro_sal_cantidad int(10),
    pro_sal_precio int(10),
    foreign key (pro_sal_salida) references salida(salida_id),
    foreign key (pro_sal_producto) references producto(producto_id)
);
