/*
约束是作用于表中字段上的规则，用于限制存储在表中的数据。
保证数据库中数据的正确、有效性和完整性。

分类：
    非空约束：NOT NULL
        限制该字段的数据不能为NULL。

    唯一约束：UNIQUE
        保证该字段的所有数据都是唯一、不重复的。

    主键约束：PRIMARY KEY
        主键是一行数据的唯一标识，要求非空且唯一。
    经常配合 AUTO_INCREMENT 使用来生成自增编号，通常用于整数类型的主键。
    提示：插入失败、事务回滚或删除记录等都可能造成自增编号不连续。

    默认约束：DEFAULT
        保存数据时，如果未指定该字段的值，则采用默认值。

    检查约束（8.0.16版本之后）：CHECK ()
        保证字段值满足某一个条件。

    外键约束：FOREIGN KEY
        用来让两张表的数据之间建立连接，保证数据的一致性和完整性。
    具有外键的表叫子表，外键指向的表叫父表。
    相关更多细节看下面的外键约束例子。

注意：约束是作用于表中字段上的，可以在创建表/修改表的时候添加约束。
*/

-- 例子(不包含外键约束)：
CREATE TABLE `user` (
    id INT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID唯一标识',
    name VARCHAR(10) NOT NULL UNIQUE COMMENT '姓名',
    age INT CHECK age > 0 AND age <= 120 COMMENT '年龄',
    status CHAR(1) DEFAULT '1' COMMENT '状态',
    gender CHAR(1) COMMENT '性别'
) COMMENT '用户表';

-- 外键约束例子：
/*
1. 添加外键约束的两种方式：
（1）创建表时添加外键：
    CREATE TABLE 表名 (
        字段名 数据类型,
        ...
        [CONSTRAINT 外键名称] FOREIGN KEY (外键字段名) REFERENCES 主表 (主表列名)
    );
[CONSTRAINT 外键名称]表示可选，实际使用时不写方括号。

（2）修改表时添加外键：
    ALTER TABLE 表名 ADD CONSTRAINT 外键名称 FOREIGN KEY (外键字段名) REFERENCES 主表 (主表列名);

注意：外键设置的前提条件：父表必须要有主键或唯一键，但子表可以不需要。

2. 删除外键：
ALTER TABLE 表名 DROP FOREIGN KEY 外键名称;

3. 删除/更新行为：
（1）NO ACTION：
    当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有则不允许删除/更新。（与 RESTRICT 一致）

（2）RESTRICT：
    当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有则不允许删除/更新。（与 NO ACTION 一致）

（3）CASCADE：
    当在父表中删除/更新对应记录时，首先检查该记录是否有对应外键，如果有，则也删除/更新外键在子表中的记录。

（4）SET NULL：
    当在父表中删除对应记录时，首先检查该记录是否有对应外键，如果有则设置子表中该外键值为 NULL。（这就要求该外键允许取 NULL）

（5）SET DEFAULT：
    父表有变更时，子表将外键列设置成一个默认的值。（InnoDB 不支持）

设置级联更新和级联删除：
    ALTER TABLE 表名 ADD CONSTRAINT 外键名称 FOREIGN KEY (外键字段) REFERENCES 主表名 (主表字段名) ON UPDATE 删除/更新行为字段 ON DELETE 删除/更新行为字段;

*/
