package database

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/pudongping/goblog/pkg/logger"
)

// DB 数据库对象
var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
	initDB()
	createTables()
}

func initDB() {

	var err error

	// 设置数据库连接信息
	config := mysql.Config{
		User:                 "root",
		Passwd:               "123456",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 准备数据库连接池
	// DSN 全称为 Data Source Name，表示【数据源信息】
	// [用户名[:密码]@][协议(数据库服务器地址)]]/数据库名称?参数列表
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// root:123456@tcp(127.0.0.1:3306)/goblog?checkConnLiveness=false&maxAllowedPacket=0
	DB, err = sql.Open("mysql", config.FormatDSN())
	logger.LogError(err)

	// 设置最大连接数，0 表示无限制，默认为 0
	// 在高并发的情况下，将值设为大于 10，可以获得比设置为 1 接近六倍的性能提升。而设置为 10 跟设置为 0（也就是无限制），在高并发的情况下，性能差距不明显
	// 最大连接数不要大于数据库系统设置的最大连接数 show variables like 'max_connections';
	// 这个值是整个系统的，如有其他应用程序也在共享这个数据库，这个可以合理地控制小一点
	DB.SetMaxOpenConns(25)

	// 设置最大空闲连接数，0 表示不设置空闲连接数，默认为 2
	// 在高并发的情况下，将值设为大于 0，可以获得比设置为 0 超过 20 倍的性能提升
	// 这是因为设置为 0 的情况下，每一个 SQL 连接执行任务以后就销毁掉了，执行新任务时又需要重新建立连接。很明显，重新建立连接是很消耗资源的一个动作
	// 此值不能大于 SetMaxOpenConns 的值，大于的情况下 mysql 驱动会自动将其纠正
	DB.SetMaxIdleConns(25)

	// 设置每个连接的过期时间
	// 设置连接池里每一个连接的过期时间，过期会自动关闭。理论上来讲，在并发的情况下，此值越小，连接就会越快被关闭，也意味着更多的连接会被创建。
	// 设置的值不应该超过 MySQL 的 wait_timeout 设置项（默认情况下是 8 个小时）
	// 此值也不宜设置过短，关闭和创建都是极耗系统资源的操作。
	// 设置此值时，需要特别注意 SetMaxIdleConns 空闲连接数的设置。假如设置了 100 个空闲连接，过期时间设置了 1 分钟，在没有任何应用的 SQL 操作情况下，数据库连接每 1.6 秒就销毁和新建一遍。
	// 这里的推荐，比较保守的做法是设置五分钟
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = DB.Ping()
	logger.LogError(err)

}

// 创建数据表
func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}
