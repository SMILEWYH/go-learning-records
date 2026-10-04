/*
子查询

    概念：SQL 语句中嵌套 SELECT 语句，称为嵌套查询，又称子查询。

    语法：
        SELECT * FROM t1 WHERE column1 = (SELECT column1 FROM t2);

    子查询外部的语句可以是 INSERT / UPDATE / DELETE / SELECT 的任何一个。

    根据子查询结果不同，分为：
        1. 标量子查询（子查询结果为单个值, 单个值等于一行一列）
        2. 列子查询（子查询结果为一列多行）
        3. 行子查询（子查询结果为一行多列）
        4. 表子查询（子查询结果为多行多列）

    根据子查询位置，分为：WHERE 之后、FROM 之后、SELECT 之后。
*/
