/*
第一部分：数据库操作
    1. 查询所有数据库：
        SHOW DATABASES;
    2. 查询当前位置处于哪个数据库：
        SELECT DATABASE();
    3. 创建数据库：
        CREATE DATABASE [IF NOT EXISTS] 数据库名 [DEFAULT CHARSET 字符集][COLLATE 排序规则];
    4. 删除数据库：
        DROP DATABASE [IF EXISTS] 数据库名;
    5. 使用数据库：
        USE database_name;
*/

/*
第二部分：表操作-查询
    1. 查询当前数据库下所有表：
        SHOW TABLES;
    2. 查询表结构
        DESC 表名;
    3. 查询表的创建语句：
        SHOW CREATE TABLE 表名;
*/

/*
第三部分：表操作-创建
    1. 创建表：
        CREATE TABLE [IF NOT EXISTS] 表名(
            字段名 数据类型 [COMMENT '表注释'],
            ...  最后一列不要加逗号
        ) [COMMENT '表注释'] [ENGINE=存储引擎] [DEFAULT CHARSET=字符集] [COLLATE=排序规则];
*/

/*
第四部分：表操作-数据类型
    MySQL中的数据类型有很多，主要分为三类:数值类型、字符串类型、日期时间类型。

    1. 数值类型：
        1.1 整数类型：
            TINYINT：1字节，范围-128~127
            SMALLINT：2字节，范围-32768~32767
            MEDIUMINT：3字节，范围-8388608~8388607
            INT：4字节，范围-2147483648~2147483647
            BIGINT：8字节，范围-9223372036854775808~9223372036854775807
        1.2 浮点数类型：
            FLOAT：4字节，单精度浮点数
            DOUBLE：8字节，双精度浮点数
        1.3 定点数类型：
            DECIMAL：用于存储精确的小数值，使用的时候需要指定精度和小数位数

        默认情况为有正负符号的范围, 如果想要无符号的范围，可以在数据类型后加上UNSIGNED关键字。

    2. 字符串类型：
        2.1 定长字符串类型：
            CHAR(n)：固定长度的字符串，长度为n，最大长度为255
        2.2 变长字符串类型：
            VARCHAR(n)：可变长度的字符串，长度为n，最大长度为65535
        
        CHAR和VARCHAR的区别：
            CHAR是定长的，如果存储的字符串长度小于n，则会在字符串后面补空格，存储时会占用固定的空间。
            VARCHAR是变长的，如果存储的字符串长度小于n，则只占用实际存储的空间加上1个字节的长度信息。
        并且前者性能好一些，后者性能差一些，可以根据不同的场景需要来使用：
            1. CHAR适合存储长度固定的字符串，如身份证号、手机号等。
            2. VARCHAR适合存储长度不固定的字符串，如用户名、地址等。

        2.3 文本类型：
            TINYBLOB：0-255 bytes，不超过255个字符的二进制数据
            TINYTEXT：0-255 bytes，短文本字符串
            BLOB：0-65535 bytes，二进制形式的长文本数据
            TEXT：0-65535 bytes，长文本数据
            MEDIUMBLOB：0-16777215 bytes，二进制形式的中等长度文本数据
            MEDIUMTEXT：0-16777215 bytes，中等长度文本数据
            LONGBLOB：0-4294967295 bytes，二进制形式的极大文本数据
            LONGTEXT：0-4294967295 bytes，极大文本数据

    3. 日期类型：
        DATE：3字节，日期值，范围1000-01-01~9999-12-31，格式YYYY-MM-DD
        TIME：3字节，时间值或持续时间，范围-838:59:59~838:59:59，格式HH:MM:SS
        YEAR：1字节，年份值，范围1901~2155（另支持0000），格式YYYY
        DATETIME：5字节，日期和时间值，范围1000-01-01 00:00:00~9999-12-31 23:59:59，格式YYYY-MM-DD HH:MM:SS
        TIMESTAMP：4字节，日期和时间值（时间戳），范围1970-01-01 00:00:01~2038-01-19 03:14:07（UTC），格式YYYY-MM-DD HH:MM:SS
*/

/*
第五部分：表操作-修改、删除表和字段
    1. 添加字段：
        ALTER TABLE 表名 ADD COLUMN 字段名 数据类型 [COMMENT '字段注释'] [AFTER 字段名];
    2. 修改字段类型：
        ALTER TABLE 表名 MODIFY COLUMN 字段名 新数据类型;
    3. 修改字段名：
        ALTER TABLE 表名 CHANGE COLUMN 旧字段名 新字段名 数据类型;
    4. 修改字段注释：
        ALTER TABLE 表名 MODIFY COLUMN 字段名 数据类型 COMMENT '新注释'; 
    5. 删除字段：
        ALTER TABLE 表名 DROP COLUMN 字段名; 
    6. 修改表名：
        ALTER TABLE 表名 RENAME TO 新表名;
    7. 删除表：
        DROP TABLE [IF EXISTS] 表名;
    8. 清空表数据：
        TRUNCATE TABLE 表名;
*/
