package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //匿名导入包 —— 只导入包但不使用包内类型和数值。匿名包和其他包一样会让导入包编译到可执行文件中，同时导入包也会触发 init()函数调用
	"strings"
	"time"
)

//我们先将数据库配置信息定义成为常量
const (
	userName = "root"
	password = "your_password"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "test"
)

//初始化数据库连接，返回数据库连接的指针引用
func InitDB() *sql.DB {
	//Golang数据连接："用户名:密码@tcp(IP:端口号)/数据库名?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	// a hander for a database. 它并没有真正的去打开一个链接。因此也并不会报错
	db, _ := sql.Open("mysql", path)
	//defer db.Close()  在适当的时候要对数据库进行关闭
	//the maximum amount of time a connection may be reused. 简单理解就是多长时间不用会回收此链接
	//官方推荐时间 小于 5 minutes
	db.SetConnMaxLifetime(time.Minute * 3)
	//设置数据库的最大链接个数， 强力推荐根据server情况进行设置。但是没有默认的推荐值
	db.SetMaxOpenConns(5)
	//设置上数据库最大闲置连接数。建议和 MaxOpenConns相同，这样可以避免数据库频繁打开和关闭
	db.SetMaxIdleConns(5)
	//Ping verifies a connection to the database is still alive, establishing a connection if necessary.
	//用来验证链接是否还存活，必要的话会建立一个链接。
	if err := db.Ping(); err != nil {
		fmt.Println("连接失败～～1")
		panic(err)
	}
	//将数据库连接的指针引用返回
	return db
}

//插入操作
func main() {
	//使用工具获取数据库连接
	db := InitDB()
	//开启事务 [这里不使用事务也是可以的]
	tx, err := db.Begin()
	if err != nil {
		//事务开启失败，直接panic
		panic(err)
	}
	//准备SQL语句 ? = placeholder
	sql := "insert into user (`name`, `age`) values (?, ?)"
	//对SQL语句进行预处理
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close() // Close the statement when we leave main() / the program terminates

	//Execute stmt
	result, err := stmt.Exec("Rock", 30)
	if err != nil {
		//SQL执行失败，直接panic
		panic(err)
	}
	//返回插入记录的id
	fmt.Println(result.LastInsertId())
	//提交事务
	tx.Commit()
}
/*
// 查询操作语句
func main() {
   db := InitDB()
   sql := "SELECT name FROM user WHERE id = ?"

   stmt, err := db.Prepare(sql)
   if err != nil {
      panic(err)
   }
   defer stmt.Close()
   var ret string
   //执行查询语句, QueryRow 用来赋值查询语句。Scan取出结果
   err = stmt.QueryRow(1).Scan(&ret)
   if err != nil {
      panic(err)
   }
   fmt.Println(ret)
}
//更新操作
func main() {
   db := InitDB()
   sql := "UPDATE user SET name=?  WHERE id = ?"

   stmt, err := db.Prepare(sql)
   if err != nil {
      panic(err)
   }
   defer stmt.Close()
   var ret string
   //执行查询语句, QueryRow 用来赋值查询语句。Scan取出结果
   _, error := stmt.Exec("张三", 1)
   if error != nil {
      panic(error)
   }
   fmt.Println(ret)
}
//删除操作
func main() {
   db := InitDB()
   sql := "Delete from user  WHERE id = ?;"

   stmt, err := db.Prepare(sql)
   if err != nil {
      panic(err)
   }
   defer stmt.Close()
   var ret string
   //执行查询语句, QueryRow 用来赋值查询语句。Scan取出结果
   _, error := stmt.Exec(1)
   if error != nil {
      panic(error)
   }
   fmt.Println(ret)
}
*/