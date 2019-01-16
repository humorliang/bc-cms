package db

import (
	"database/sql"
	"com/logging"
)

//信息结构体
type QueryInfo struct {
	info map[string]interface{}
}

//查询获取多个字段值
func Querys(rows *sql.Rows) ([]map[string]interface{}, error) {
	//获取查询的列
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	//长度
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	//结果值集合
	values := make([]interface{}, count)
	//结果值所对应的指针
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		//地址传递
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		//要写入的数据
		entry := make(map[string]interface{})
		for i, col := range columns {
			//定义空接口接受所有值
			var v interface{}
			//对值进行判断
			val := values[i]
			//类型判断
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

//事务查询
//返回不同sql结果集合
func TranscationQuerys(tx *sql.Tx, querySqls ...string) (map[int][]map[string]interface{}, error) {
	//结果存储容器
	txs := make(map[int][]map[string]interface{}, len(querySqls))
	//执行不同查询语句
	for i := 0; i < len(querySqls); i++ {
		//fmt.Println(querySqls[i])
		rows, err := tx.Query(querySqls[i])
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				logging.Error(err)
				return nil, err
			}
			return nil, err
		}
		res, err := Querys(rows) //[]map[string]interface{}
		//fmt.Println(res)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				logging.Error(err)
				return nil, err
			}
			return nil, err
		}
		txs[i] = res
	}
	//提交事务
	err := tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			logging.Error(err)
			return nil, err
		}
		return nil, err
	}
	return txs, nil
}
