package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 通常情况下, 我们会使用Model结尾表示这个结构体类型是表的结构体
type UserModel struct {
	ID   int64  `gorm:"primaryKey"`      // 主键
	Name string `gorm:"not null;unique"` // 不能为 nil, 且唯一
	// 更多标签参考这个 https://gorm.io/zh_CN/docs/models.html#字段标签
	CreatedAt time.Time // 通过 GORM 创建该行数据时, 如果它是零值, GORM 自动填入当前时间; 如果已经手动赋值, 则保留该值
}

// 自定义 UserModel 结构体生成表时的表名
func (*UserModel) TableName() string {
	return "user_model"
}

// 如果在插入的钩子函数中，返回错误，则数据不会被插入
func (u *UserModel) BeforeCreate(tx *gorm.DB) error {
	fmt.Println("创建的钩子函数")
	return nil
}

// 如果在更新的钩子函数中, 返回错误, 则数据不会被插入
func (u *UserModel) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("更新的钩子函数")
	return nil
}

func main() {
	db, err := gorm.Open(mysql.Open("root:b1Ol49pbgjEPLvQOO-gvgDwG2t7X7DAb@tcp(127.0.0.1:3306)/vvyh_sql"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// 通用方法:
	// 执行数据库操作的方法之前链式调用 Debug(), 可输出该操作的 SQL 语句、耗时、影响行数和日志
	// 在查询、单表更新或单表删除方法之前链式调用 Limit(n), 可限制本次操作的记录数量, 比如查询最多返回 n 条, 更新最多匹配 n 条, 删除最多删除 n 条
	// 查询时接收结果的单个结构体、更新时 Model 传入的结构体、删除时 Delete 传入的结构体, 若已有非零主键, GORM 会将其作为 AND 条件加入 WHERE

	// 注意事项:
	// 1. 单个结构体传递指针, 结构体切片传递值, 这样可以让 Gorm 做回填操作, 也就是将数据库表中的数据回填到结构体中, 方便后续使用

	// 1. Gorm 的方法都是链式调用的, 所以每次调用方法后都要接收返回值, 否则会报错
	// 2. Gorm 的方法都是返回一个 *gorm.DB 类型的实例, 所以可以继续链式调用
	// 3. Gorm 的方法都是返回一个 *gorm.DB 类型的实例, 所以可以继续链式调用, 但是如果你不接收返回值, 那么就会丢失之前的链式调用的结果, 所以要接收返回值

	// 第一步: 通过结构体初始化表结构
	initTable(db)

	// 第二步: 增删改查
	// 1.1 单次插入
	createData_sigle(db)
	// 1.2 批量插入
	createData_batch(db)
	// 1.3 插入数据前的钩子函数
	// 看结构体的 BeforeCreate 方法

	// 2.1 查询全部数据
	findAllData(db)
	// 2.2 条件查询
	findDataByRules(db)
	// Find 未查到符合条件的行时不会返回 Error, 但可以通过返回的 RowsAffected == 0 判断是否未查到记录
	// 使用切片接收结果时, 也可以通过 len(userList) == 0 判断

	// 2.3 查询一条数据, 同样可以写查询条件; 以把 Take 当成Limit(1).Find(), 但是两者对“没查到记录”的处理不同, 前者会返回 Error, 后者不会
	// db.Take()
	// 2.4 根据主键排序查询第一个和最后一个
	// db.First()
	// db.Last()

	// 这里还未完善 start
	// 3.1 Save 更新方法: 传入的结构体没有主键值时创建记录; 有主键值时更新主键值对应行数据的全部字段: 传递新值则覆盖原值, 没传递新值的字段会使用 go 数据类型的零值来进行覆盖
	// 创建了一条新数据
	SaveNewData(db)
	// 更新主键为 9 的数据
	SaveData(db)

	// 3.2 Update 更新方法: 不像 Save 那样通过主键进行分支操作, 而是必须要有查询条件找到符合查询条件的行数据, 然后更新其中单个字段
	db.Model(&UserModel{}).
		Where("id = ?", 1).
		Update("name", "张三") // 走了更新钩子

	// 3.3 UpdateColumn 更新方法: 和 Update 差不多, Update 会执行更新钩子并自动维护更新时间, UpdateColumn, 会跳过这两项
	db.Model(&UserModel{}).
		Where("id = ?", 1).
		UpdateColumn("name", "张三")

	// 3.4 Updates 更新方法: 和 Update 差不多, 它的特殊点在于不仅可以传递结构体, 还可以传递 map
	user = UserModel{ID: 1}

	db.Model(&user).Updates(UserModel{
		Name: "张三丰",
	})

	user = UserModel{ID: 2}
	// 不会更新零值
	db.Model(&user).Updates(UserModel{
		Name: "",
	})

	user = UserModel{ID: 2}
	// 会更新零值
	db.Model(&user).Updates(map[string]any{
		"name": "",
	})

	// 再看看这些方法的零值边界问题

	// gorm.Expr 方法通常用于获取原字段的数据
	db.Model(&UserModel{}).
		Where("id = ?", 1).
		UpdateColumn("age", gorm.Expr("age + ？", 1)) // 这样可以防止 sql 注入
	// 这里还未完善 end
}

func initTable(db *gorm.DB) {
	// AutoMigrate 方法是根据你定义的 Go 结构体, 在数据库中自动创建表, 或者更新已有表的结构: (不怎么使用, 使用原生 sql 建表、修改表更好)
	// 首次调用并传递 UserModel 结构体它会创建一个 user_model 表
	// 接下来如果结构体新增了一个字段, 然后再调用一次, 则可以给已存在的表新增对应的字段
	// 接下来如果删除了结构体的一个字段, 然后再调用一次, 这不会删除已存在的表中对应的字段和数据
	// 接下来如果修改了结构体的字段, 然后再调用一次, 那么这里的规则有点复杂, 略
	err := db.AutoMigrate(&UserModel{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func createData_sigle(db *gorm.DB) {
	user := UserModel{Name: "张三"}

	// var result *gorm.DB
	result := db.Create(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	// 通过回填, 可以得倒这一次操作的对应行数据
	fmt.Println(user.ID, user.Name, user.CreatedAt)
}

func createData_batch(db *gorm.DB) {
	var users = []UserModel{
		{
			Name: "王五",
		},
		{
			Name: "李四",
		},
	}

	result := db.Create(users) // 切片按值复制复制的是切片首个下索引数据的内存地址、len、cap
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
}

func findAllData(db *gorm.DB) {
	var userList []UserModel

	// 因为 CreatedAt 的类型是 time.Time, 而 MySQL 的时间类型返回的是日期时间文本
	// 所以需要在 DSN 中设置 parseTime=true, 使驱动将数据库的日期时间值解析为 time.Time, 否则读取该字段时可能报类型转换错误
	result := db.Find(&userList)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	fmt.Println(userList)
}

func findDataByRules(db *gorm.DB) {
	var userList []UserModel

	result1 := db.Find(&userList, "name = ?", "张三")
	result2 := db.Find(&userList, "name = ? AND id > ?", "张三", 10) // 多条件
	if result1.Error != nil {
		fmt.Println(result1.Error)
		return
	}
	if result2.Error != nil {
		fmt.Println(result2.Error)
		return
	}
}

func SaveNewData(db *gorm.DB) {
	user := UserModel{
		Name: "枫枫1",
	}

	result := db.Save(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
}

func SaveData(db *gorm.DB) {
	user := UserModel{
		ID:        9,
		Name:      "枫枫2",
		CreatedAt: time.Now(), // 这里和 AutoMigrate 创建的自动添加不同, Save 修改不会自动修改, 所以需要你手动传递一个新的时间点过去
	}

	result := db.Save(&user)
	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
}
