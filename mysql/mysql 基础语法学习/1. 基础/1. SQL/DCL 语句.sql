/*
第一部分：管理用户
    1. 查询用户
        USE mysql;
        SELECT * FROM user;

    2. 创建用户
        CREATE USER '用户名'@'主机名' IDENTIFIED BY '密码';

    3. 修改用户密码
        ALTER USER '用户名'@'主机名' IDENTIFIED WITH mysql_native_password BY '新密码';

    4. 删除用户
        DROP USER '用户名'@'主机名';

    主机名决定哪些主机可以连接到数据库服务器,
    如果是本地连接, 主机名可以使用localhost, 
    如果是远程连接, 主机名可以使用IP地址或者域名,
    也可以使用通配符%, 表示所有主机都可以连接到数据库服务器。
*/

/*
第二部分：权限控制
    权限列表：
        ALL, ALL PRIVILEGES    所有权限
        SELECT                 查询数据
        INSERT                 插入数据
        UPDATE                 修改数据
        DELETE                 删除数据
        ALTER                  修改表
        DROP                   删除数据库/表/视图
        CREATE                 创建数据库/表

    1. 查询权限
        SHOW GRANTS FOR '用户名'@'主机名';

    2. 授予权限
        GRANT 权限列表 ON 数据库名.表名 TO '用户名'@'主机名';

    想给全部数据库和表的权限则使用 *.* 
    想给某个数据库的全部表的权限则使用 数据库名.* 
    想给某个数据库的某个表的权限则使用 数据库名.表名

    3. 撤销权限
        REVOKE 权限列表 ON 数据库名.表名 FROM '用户名'@'主机名';
*/
