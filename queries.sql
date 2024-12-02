select * from transaccion_salida_producto;
select * from producto_salida;
select * from salida;

INSERT INTO salida (salida_fecha, salida_tipo, salida_usuario)
values (NOW(), 1, 1);
INSERT INTO producto_salida (pro_sal_salida, pro_sal_producto, pro_sal_cantidad, pro_sal_precio)
select (LAST_INSERT_ID(), p.producto_id, tsp.quantity, tsp.price)
from transaccion_salida_producto tsp
JOIN producto p on p.producto_id = (
    SELECT pr.producto_id from producto
    where pr.producto_codigo = tsp.code);