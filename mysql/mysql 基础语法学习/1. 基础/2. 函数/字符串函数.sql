/*
MySQL中内置了很多字符串函数，常用的几个如下：

    CONCAT(S1, S2, ... Sn)        字符串拼接，将S1、S2、... Sn拼接成一个字符串
    LOWER(str)                   将字符串str全部转为小写
    UPPER(str)                   将字符串str全部转为大写
    LPAD(str, n, pad)            左填充，用字符串pad对str的左边进行填充，达到n个字符的长度
    RPAD(str, n, pad)            右填充，用字符串pad对str的右边进行填充，达到n个字符的长度
    TRIM(str)                    去掉字符串头部和尾部的空格
    SUBSTRING(str, start, len)   返回字符串str从start位置起的len个字符

    简单的调用语法：
        SELECT 函数(参数);
*/
