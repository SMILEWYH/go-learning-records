-- 行子查询
-- 子查询返回的结果是一行（可以是多列），这种子查询称为行子查询。
-- 常用的操作符：=  <>  IN  NOT IN

-- 1. 查询与 "张无忌" 的薪资及直属领导相同的员工信息；
-- a. 查询 "张无忌" 的薪资及直属领导
select salary, managerid from emp where name = '张无忌';

-- b. 查询与 "张无忌" 的薪资及直属领导相同的员工信息；
select * from emp where (salary,managerid) = (select salary, managerid from emp where name = '张无忌');
