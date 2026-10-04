-- 表子查询
-- 子查询返回的结果是多行多列，这种子查询称为表子查询。
-- 常用的操作符：IN
-- 表子查询一般常用于 FROM 后面，作为临时表参与查询，也可以配合 IN 使用。

-- 1. 查询与 "鹿杖客"，"宋远桥" 的职位和薪资相同的员工信息
-- a. 查询 "鹿杖客"，"宋远桥" 的职位和薪资
select job, salary from emp where name = '鹿杖客' or name = '宋远桥';

-- b. 查询与 "鹿杖客"，"宋远桥" 的职位和薪资相同的员工信息
select * from emp where (job,salary) in (select job, salary from emp where name = '鹿杖客' or name = '宋远桥');

-- 2. 查询入职日期是 "2006-01-01" 之后的员工信息，及其部门信息
-- a. 入职日期是 "2006-01-01" 之后的员工信息
select * from emp where entrydate > '2006-01-01';

-- b. 查询这部分员工，对应的部门信息；
select e.*, d.* from (select * from emp where entrydate > '2006-01-01') e left join dept d on e.dept_id = d.id;
