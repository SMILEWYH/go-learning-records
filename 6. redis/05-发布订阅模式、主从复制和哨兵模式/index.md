# 5. 发布订阅模式、主从复制和哨兵模式

## 5.1 发布订阅模式

Redis 发布订阅（pub/sub）是一种消息通信模式：发送者发（pub）送消息，订阅者(sub)接收消息。

Redis 客户端可以订阅任意数量的频道。

订阅 / 发布消息图：

![img](img/Jx8lVS.png)

下图展示了频道 channel1，已经订阅这个频道的三个客户端。

![img](img/Jx8jZ8.png)

当有新消息通过 publish 命令发送给频道 channel1 时，这个消息就会被发送给订阅它的三个客户端。

> 命令

这些命令被广泛应用于构建即时通讯应用、比如网络聊天室和实时广播、实时提醒等。

| 序号 | 命令及描述                                                   |
| :--- | :----------------------------------------------------------- |
| 1    | [PSUBSCRIBE pattern [pattern ...\]](https://www.runoob.com/redis/pub-sub-psubscribe.html) 订阅一个或多个符合给定模式的频道。 |
| 2    | [PUBSUB subcommand [argument [argument ...\]]](https://www.runoob.com/redis/pub-sub-pubsub.html) 查看订阅与发布系统状态。 |
| 3    | [PUBLISH channel message](https://www.runoob.com/redis/pub-sub-publish.html) 将信息发送到指定的频道。 |
| 4    | [PUNSUBSCRIBE [pattern [pattern ...\]]](https://www.runoob.com/redis/pub-sub-punsubscribe.html) 退订所有给定模式的频道。 |
| 5    | [SUBSCRIBE channel [channel ...\]](https://www.runoob.com/redis/pub-sub-subscribe.html) 订阅给定的一个或多个频道的信息。 |
| 6    | [UNSUBSCRIBE [channel [channel ...\]]](https://www.runoob.com/redis/pub-sub-unsubscribe.html) 指退订给定的频道。 |

Redis 的频道不需要提前创建，直接使用频道名就可以。

> 测试

客户端 A 订阅了名为 **redisChat** 的频道。

```bash
redis 127.0.0.1:6379> SUBSCRIBE redisChat

Reading messages... (press Ctrl-C to quit)
1) "subscribe"  # 频道订阅成功
2) "redisChat"  # 频道名称
3) (integer) 1  # 当前客户端已经订阅频道的数量
```

客户端 B 在同一个频道 **redisChat** 发布两次消息，订阅者就能接收到消息。

```bash
redis 127.0.0.1:6379> PUBLISH redisChat "Redis is a great caching technique"

(integer) 1

redis 127.0.0.1:6379> PUBLISH redisChat "Learn redis by mszlu.com"

(integer) 1

# 订阅者的客户端会显示如下消息
1) "message"
2) "redisChat"
3) "Redis is a great caching technique"
1) "message"
2) "redisChat"
3) "Learn redis by mszlu.com"
```

> 原理

Redis 是使用 C 实现的，通过分析 Redis 源码里的 public.c 文件，了解发布和订阅机制的底层实现，借此加深对 Redis 的理解。Redis 通过 public 、subscribe 和 psubscribe 等命令实现发布和订阅功能。

**使用场景：**

1、实时消息系统

2、实时聊天

3、订阅、关注系统都可以

稍微复杂的场景更多的使用消息中间件 MQ。



目前没想到有合适的业务场景使用 Redis 发布订阅模式。

注意哨兵模式的部分通信会使用发布订阅，各个哨兵会在被监控的 Redis 主从节点上，发布和订阅这个特殊频道sentinel:hello，这样哨兵就能相互发现、交换部分配置信息。

## 5.2 主从复制

### 5.2.1 概念

主从复制，是指将一台 Redis 服务器的数据，复制到其他的 Redis 服务器。前者称之为主节点（master），后者称之为从节点（slave）；数据的复制都是单向的，只能从主节点到从节点。Master 以写为主，Slave 以读为主。

默认情况下，每台 Redis 服务器都是主节点。且一个主节点可以有多个从节点或者没有从节点，但是反过来一个从节点只能有一个主节点。

### 5.2.2 主从复制的作用

1、数据冗余：主从复制实现了数据的热备份，是持久化的之外的一种数据冗余方式。

2、故障恢复（需要哨兵）：当主节点出现问题时，可以由从节点提供服务，实现快速的故障恢复。实际也是一种服务的冗余。

3、负载均衡：在主从复制的基础上，配合读写分离，可以由主节点提供写服务，由从节点提供读服务，分担服务器负载；尤其是在写少读多的场景下，通过多个节点分担读负载，可以大大提高 Redis 服务器的并发量。

4、高可用（集群）的基石：除了上述作用以外，主从复制还是哨兵模式和集群能够实施的基础，因此说主从复制是 Redis 高可用的基础。

一般来说，要将Redis 运用于工程项目中，只使用一台 Redis 是万万不能的（可能会宕机），原因如下：

1、从结构上，单个 Redis 服务器会发生单点故障，并且一台服务器需要处理所有的请求负载，压力很大；

2、从容量上，单个 Redis 服务器内存容量有限，就算一台 Redis 服务器内存容量为 256G， 也不能将所有的内存用作 Redis 存储内存，一般来说，**单台 Redis最大使用内存不应该超过 20G**。



电商网站上的商品，一般都是一次上传，无数次浏览的，说专业点就是“读多写少”。

对于这种场景，我们可以使用如下这种架构：

![img](img/YsMvfP.png)

主从复制，读写分离！80% 的情况下，都是在进行读操作。这种架构可以减少服务器压力，经常使用实际生产环境中，最少是“一主二从”的配置。真实环境中不可能使用单机 Redis。

### 5.2.3 一主二从的搭建

默认情况下，每台 Redis 服务器都是主节点，我们一般情况下，只用配置从库就好了。

我们先启动三个 Redis 服务端（6381、6382、6383），然后随便查看其中一个的信息：

```bash
[root@coder bin]# redis-cli -p 6381
127.0.0.1:6381> ping
PONG
127.0.0.1:6381> info replication			# 查看当前库的信息
# Replication
role:master									    # 角色
connected_slaves:0							# 当前没有从库
master_replid:2467dd9bd1c252ce80df280c925187b3417055ad
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:0
repl_backlog_size:1048576
repl_backlog_first_byte_offset:0
repl_backlog_histlen:0
127.0.0.1:6381> 
```

接下来选择 6382、6383 作为从库。

配置从库的方式有两种：

​	1、直接使用命令配置，这种方式当 Redis 重启后配置会失效。

​	2、使用配置文件。

这里使用命令演示一下：

```bash
127.0.0.1:6382> SLAVEOF 127.0.0.1 6381		# SLAVEOF host  port
OK
127.0.0.1:6382> info replication
# Replication
role:slave			# 角色已经是从机了
master_host:127.0.0.1	# 主节点地址
master_port:6381			# 主节点端口
master_link_status:up
master_last_io_seconds_ago:6
master_sync_in_progress:0
slave_repl_offset:0
slave_priority:100
slave_read_only:1
connected_slaves:0
master_replid:907bcdf00c69d361ede43f4f6181004e2148efb7
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:0
second_repl_offset:-1
repl_backlog_active:1
repl_backlog_size:1048576
repl_backlog_first_byte_offset:1
repl_backlog_histlen:0
```

配置好了之后，看主机：

```bash
127.0.0.1:6381> info replication
# Replication
role:master
connected_slaves:2		# 主节点下有两个从节点
slave0:ip=127.0.0.1,port=6382,state=online,offset=420,lag=1
slave1:ip=127.0.0.1,port=6383,state=online,offset=420,lag=1
master_replid:907bcdf00c69d361ede43f4f6181004e2148efb7
master_replid2:0000000000000000000000000000000000000000
master_repl_offset:420
second_repl_offset:-1
repl_backlog_active:1
repl_backlog_size:1048576
repl_backlog_first_byte_offset:1
repl_backlog_histlen:420
```

这里使用命令是暂时的。真实的主从配置应该是在 redis.conf 配置文件中配置，这样才是永久的。

```bash
################################# REPLICATION #################################

# Master-Replica replication. Use replicaof to make a Redis instance a copy of
# another Redis server. A few things to understand ASAP about Redis replication.
#
#   +------------------+      +---------------+
#   |      Master      | ---> |    Replica    |
#   | (receive writes) |      |  (exact copy) |
#   +------------------+      +---------------+
#
# 1) Redis replication is asynchronous, but you can configure a master to
#    stop accepting writes if it appears to be not connected with at least
#    a given number of replicas.
# 2) Redis replicas are able to perform a partial resynchronization with the
#    master if the replication link is lost for a relatively small amount of
#    time. You may want to configure the replication backlog size (see the next
#    sections of this file) with a sensible value depending on your needs.
# 3) Replication is automatic and does not need user intervention. After a
#    network partition replicas automatically try to reconnect to masters
#    and resynchronize with them.
#
# replicaof <masterip><masterport>			# 这里配置 master 

# If the master is password protected (using the "requirepass" configuration
# directive below) it is possible to tell the replica to authenticate before
# starting the replication synchronization process, otherwise the master will
# refuse the replica request.
#
# masterauth <master-password>
```

当前资料的主从配置的 redis-6382.conf、redis-6383.conf 也是配置了 replicaof 字段。

### 5.2.4 几个问题

1、主机可以写，从机不能写只能读。主机中的所有信息和数据都会保存在从机中。如果从机尝试进行写操作就会报错。

```bash
127.0.0.1:6382> get k1			# k1的值是在主机中写入的，从机中可以读取到。
"v1"
127.0.0.1:6382> set k2 v2			# 从机尝试写操作，报错了
(error) READONLY You can't write against a read only replica.
```

2、普通主从复制中，主机宕机，从机不会自动变成主节点（需要哨兵）。因此主机宕机了，从机仍保留已经同步到自己内存中的数据，默认可以继续提供读取（但数据可能不是最新的），但是还是没有写操作，并且它会不断尝试重新连接主机。

3、使用命令行配置的主从机，如果从机重启了，就会变回主机。如果再通过命令变回从机的话，立马就可以从主机中获取值。这是复制原理决定的。

### 5.2.6 复制原理

可以分三个阶段理解：

1、 首次同步（全量复制）：主机通过子进程生成 rdb 格式的数据快照，然后发送给从机，从机再加载到内存；并且快照生成、传输期间，主机可以处理写请求，所以还要把这段时间的新修改补发给从机。

2、 正常运行（增量复制）：主机发送数据变更的命令流到从机，更新从机内存数据。

3、 断线重连（全量复制、部分重同步）：主机和从机之间会对比复制历史和复制偏移量，根据情况选择进行全量复制还是部分重同步。



主从复制的两种模式：

### 第一种模式（常规）

![img](img/Ys7G5T.png)

### 第二种模式

![img](img/Ys7gRe.png)

第二种模式下 master 复制压力会降低，但复制链条上后续 slave 数据同步也越慢。

这种模式的话，将 6381 的主节点配置为 6380 。主节点 6379 只有一个从机。

如果现在 6379 节点宕机了， 6380 和 6381 节点都是从节点，只能进行读操作，都不会自动变为主节点。需要手动将其中一个变为主节点，使用如下命令：

```bash
SLAVEOF no one
```

## 5.3 哨兵模式

### 5.3.1 概述

主从切换技术的方式是：当主机服务器宕机之后，需要手动将一台服务器切换为主服务器，这需要人工干预，费时费力，还会造成一段时间内的服务不可用。这不是一种推荐的方式，更多的时候我们优先考虑的的是哨兵模式。

> 哨兵模式简略流程

哨兵通过不断向主库 ping 消息来判断主机是否故障。如果故障了，哨兵节点之间会进行投票，选出一个 leader 哨兵节点来选择把哪个从库转换为主库。

> 哨兵本质

哨兵本质上是以 Sentinel 模式运行的一个独立 Redis。

![img](img/YsvNuD.png)

这里的哨兵有两个作用

- 通过 ping 发送消息，让 Redis 服务器响应，以此来监控其运行状态，包括主服务器和从服务器
- 当哨兵检测到 master 宕机，会自动将 slave 切换为 master，然后通过发布订阅模式通知其他的从放服务器，修改配置文件，让他们切换主机。

然而一个哨兵进程对 Redis 服务器进行监控，可能会出现问题，为此，我们可以使用多个哨兵进行监控。各个哨兵之间还会进行监控，这样就形成了**多哨兵模式**。

![img](img/YypgWn.png)

> （多）哨兵模式详细流程

假设主服务器宕机了，其中一个哨兵先检测到这个结果，还不能够进行故障恢复，因为仅仅是多个哨兵中的一个认为主库不可用，这个现象称之为**主观下线**。当后面的哨兵也检测到主库不可用，并且数量达到 sentinel.conf 哨兵配置文件的数量时，那么才会进行故障恢复，这个过程称之为**客观下线**。

故障恢复流程：哨兵之间就会进行一次投票，选出来一个 leader 哨兵，它来根据从库条件，选择其中一个从库转换为主库。切换成功后，leader 哨兵还需要命令其他从节点连接新主节点，再通过发布订阅模式，让各个哨兵们更新自己的记录。

> 多哨兵之间为什么要选择出 leader 哨兵，并且是按照什么规则选的？



> leader 哨兵按照什么规则选哪个从库转换为主库



### 5.3.2 配置一个一主两从的哨兵模式

1、配置哨兵模式配置文件，新建文件 `/usr/local/bin/redis-config/sentinel.conf`。

```bash
# sentinel monitor 被监控的主节点起的名字（随便写） 主库地址 	判断客观下线需要的哨兵数量
sentinel monitor mymaster 127.0.0.1 6401 2
```

2、`redis-sentinel`  启动哨兵

```bash
[root@coder bin]# ls
6379.log  6381.log      dump6380.rdb  dump.rdb  redis-benchmark  redis-check-rdb  redis-sentinel
6380.log  dump6379.rdb  dump6381.rdb  redis-config   redis-check-aof  redis-cli        redis-server
[root@coder bin]# redis-sentinel redis-config/sentinel.conf 
2421:X 15 May 2020 20:24:06.847 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
2421:X 15 May 2020 20:24:06.847 # Redis version=5.0.5, bits=64, commit=00000000, modified=0, pid=2421, just started
2421:X 15 May 2020 20:24:06.847 # Configuration loaded
                _._                                                  
           _.-``__ ''-._                                             
      _.-``    `.  `_.  ''-._           Redis 5.0.5 (00000000/0) 64 bit
  .-`` .-```.  ```\/    _.,_ ''-._                                   
 (    '      ,       .-`  | `,    )     Running in sentinel mode
 |`-._`-...-` __...-.``-._|'` _.-'|     Port: 26379
 |    `-._   `._    /     _.-'    |     PID: 2421
  `-._    `-._  `-./  _.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |           http://redis.io        
  `-._    `-._`-.__.-'_.-'    _.-'                                   
 |`-._`-._    `-.__.-'    _.-'_.-'|                                  
 |    `-._`-._        _.-'_.-'    |                                  
  `-._    `-._`-.__.-'_.-'    _.-'                                   
      `-._    `-.__.-'    _.-'                                       
          `-._        _.-'                                           
              `-.__.-'                                               

2421:X 15 May 2020 20:24:06.848 # WARNING: The TCP backlog setting of 511 cannot be enforced because /proc/sys/net/core/somaxconn is set to the lower value of 128.
2421:X 15 May 2020 20:24:06.851 # Sentinel ID is 100430af0018d23bd1ae2fe57e71e0d45f64d9a5
2421:X 15 May 2020 20:24:06.851 # +monitor master myredis 127.0.0.1 6379 quorum 1
2421:X 15 May 2020 20:24:06.852 * +slave slave 127.0.0.1:6381 127.0.0.1 6381 @ myredis 127.0.0.1 6379
```

启动成功~！

如果现在 Master 节点宕机了，这个时候会从从库中选择一个作为主库。

如果原来的主库恢复运行了，只能归到新的主库下，作为从库， 这就是哨兵模式的规则。



**待优化**

> 哨兵模式的优点

1、哨兵集群，基于主从复制模式，所有的主从配置优点，它全有。

2、主从可以切换，故障可以转移，系统的可用性就会更好。

3、哨兵模式就是主从模式的升级，手动到自动，更加健壮。

> 哨兵模式的缺点

1、Redis 不方便在线扩容，集群达到一定的上限，在线扩容就会十分麻烦。

2、实现哨兵模式的配置其实也很麻烦，里面有甚多的配置项。