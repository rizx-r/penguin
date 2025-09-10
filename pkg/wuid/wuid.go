package wuid

import (
	"database/sql"
	"fmt"
	"github.com/edwingeng/wuid/mysql/wuid"
	"sort"
	"strconv"
)

// 全局的 WUID 实例，用来生成唯一 ID
var w *wuid.WUID

// Init /*
// dns: 数据库连接字符串
func Init(dsn string) {
	// newDB: 用来创建 MySQL 连接的函数
	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, false, err
		}
		return db, true, nil
	}

	w = wuid.NewWUID("default", nil) // 创建一个新的 WUID 生成器。
	// 从 MySQL 的 wuid 表里加载 Worker ID（类似雪花算法中的机器号），保证分布式环境下不会重复。
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

// GenUid 生成唯一 ID
func GenUid(dsn string) string {
	if w == nil {
		Init(dsn)
	}

	return fmt.Sprintf("%#016x", w.Next()) // 把整数转成 固定长度的 16 位十六进制字符串（带 0x 前缀）。例如：0x00000017f6a3b12d
}

func CombineId(aid, bid string) string {
	ids := []string{aid, bid}
	sort.Slice(ids, func(i, j int) bool {
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		return a < b
	})

	return fmt.Sprintf("%s_%s", ids[0], ids[1])
	/*
		用户 A 的 ID：0x0000000000000001
		用户 B 的 ID：0x0000000000000002
		CombineId 结果：0x0000000000000001_0x0000000000000002
	*/
}
