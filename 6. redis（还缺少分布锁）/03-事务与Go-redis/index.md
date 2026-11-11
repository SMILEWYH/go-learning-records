# 3. Redis的事务和SpringBoot整合

## 3.1 事务

ACID：原子性，一致性，隔离性，持久性。



Redis的事务 不是原子性的，没法回滚，某一条命令失败其他命令也会执行。

redis这里不支持事务回滚，指的是不支持事务运行时x错误的事务回滚。



Redis 事务的本质：一组命令的集合！一个事务中的所有命令都会被序列化，在事务执行过程中，会按照顺序执行。

一次性、顺序性、排他性的执行一组命令。

Redis 事务没有隔离级别的概念。

所有的命令在事务中，并没有直接被执行，只有发起执行命令的时候才会执行（exec）。

Redis 单条命令是保证原子性的，但是事务不保证原子性。

Redis 事务的命令：

- 开启事务：multi
- 命令入队
- 执行事务：exec
- 撤销事务：discard

### 3.1.1 正常执行事务

~~~shell
127.0.0.1:6379> multi		# 开启事务
OK
127.0.0.1:6379> set k1 v1
QUEUED
127.0.0.1:6379> set k2 v2
QUEUED
127.0.0.1:6379> get k2
QUEUED
127.0.0.1:6379> set k3 v3
QUEUED
127.0.0.1:6379> exec			# 执行事务
1) OK
2) OK
3) "v2"
4) OK
127.0.0.1:6379> 
~~~

### 3.1.2 放弃事务

~~~shell
127.0.0.1:6379> multi
OK
127.0.0.1:6379> set m1 n1
QUEUED
127.0.0.1:6379> set m2 n2
QUEUED
127.0.0.1:6379> DISCARD		# 放弃事务
OK
127.0.0.1:6379> get m1		# 事务队列中命令都不会被执行
(nil)
~~~

> 编译型异常：命令有错，事务中所有的命令都不会被执行

~~~shell
127.0.0.1:6379> multi
OK
127.0.0.1:6379> set k1 v1
QUEUED
127.0.0.1:6379> set k2 v2
QUEUED
127.0.0.1:6379> setget k3 v3		# 错误的命令
(error) ERR unknown command `setget`, with args beginning with: `k3`, `v3`, 
127.0.0.1:6379> set k4 v4
QUEUED
127.0.0.1:6379> exec		# 执行事务报错
(error) EXECABORT Transaction discarded because of previous errors.
127.0.0.1:6379> get k4		# 所有的命令都不会被执行
(nil)
127.0.0.1:6379> 
~~~

> 运行时异常：如果事务中某条命令执行结果报错，其他命令是可以正常执行的，错误命令抛出异常

~~~shell
127.0.0.1:6379> set k1 "v1"
OK
127.0.0.1:6379> multi
OK
127.0.0.1:6379> incr k1		# 会执行失败
QUEUED
127.0.0.1:6379> set k2 v2
QUEUED
127.0.0.1:6379> set k3 v3
QUEUED
127.0.0.1:6379> get k3
QUEUED
127.0.0.1:6379> exec
1) (error) ERR value is not an integer or out of range		# 第一条命令执行失败，其余的正常执行
2) OK
3) OK
4) "v3"
127.0.0.1:6379> get k2
"v2"
127.0.0.1:6379> 
~~~

### 3.1.3 监视 Watch （面试常问）

悲观锁：很悲观，认为什么时候都会出问题，无论什么都会加锁。影响效率，实际情况一般会使用乐观锁。

乐观锁：很乐观，认为什么时候都不会出现问题，所以不上锁。更新数据的时候会判断一下，在此期间是否修改过监视的数据。

首先要了解redis事务中watch的作用，watch命令可以监控一个或多个key，一旦其中有一个key被修改（或删除），之后的事务就不会执行。监控一直持续到exec命令（事务中的命令是在exec之后才执行的，所以在multi命令后可以修改watch监控的键值）。假设我们通过watch命令在事务执行之前监控了多个Keys，倘若在watch之后有任何Key的值发生了变化，exec命令执行的事务都将被放弃，同时返回Null multi-bulk应答以通知调用者事务执行失败。

所以，需要注意的是watch监控键之后，再去操作这些键，否则watch可能会起不到效果。

> Redis 监视测试

正常测试：

```bash
127.0.0.1:6379> set money 100		
OK
127.0.0.1:6379> set out 0
OK
127.0.0.1:6379> watch money		# 监视 money 对象
OK
127.0.0.1:6379> multi		# 事务正常结束，执行期间，money 没有变动，这个时候就能执行成功了
OK
127.0.0.1:6379> DECRBY money 20
QUEUED
127.0.0.1:6379> INCRBY out 20
QUEUED
127.0.0.1:6379> exec
1) (integer) 80
2) (integer) 20
127.0.0.1:6379> 
```

测试多线程修改值，使用 watch 可以当做 Redis 的乐观锁操作。

```bash
127.0.0.1:6379> set money 100
OK
127.0.0.1:6379> set out 10
OK
127.0.0.1:6379> watch money	# 监视 money
OK
127.0.0.1:6379> multi
OK
127.0.0.1:6379> DECRBY money 10
QUEUED
127.0.0.1:6379> DECRBY out 10
QUEUED
127.0.0.1:6379> exec		# 执行之前，在另外一个线程 B 中修改 money 的值，下面就是执行失败。
(nil)
127.0.0.1:6379> 
```

B 线程：

```bash
[root@coder bin]# redis-cli -p 6379
127.0.0.1:6379> set money 30
OK
```

如果修改失败，获取最新的值就好。

```bash
127.0.0.1:6379> UNWATCH		# 事务执行失败，先解锁
OK
127.0.0.1:6379> WATCH money		# 获取最新的值，再次监视。相当于 MySQL 中的 select version
OK
127.0.0.1:6379> multi
OK
127.0.0.1:6379> DECRBY money 1
QUEUED
127.0.0.1:6379> INCRBY out 1
QUEUED
127.0.0.1:6379> exec		# 执行的时候会对比监视的值，如果发生变化会执行失败。
1) (integer) 29
2) (integer) 11
127.0.0.1:6379> 
```

