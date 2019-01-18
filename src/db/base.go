package db

import (
	"database/sql"
	"com/logging"
	"com/gmysql"
)

//定义结果类型
type RowsInfo struct {
	rowsId  int64
	rowsNum int64
}

//查询获取多个字段值
func Querys(rows *sql.Rows) ([]map[string]interface{}, error) {
	//获取查询的列
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	//types, err := rows.ColumnTypes()
	//for _, v := range types {
	//	fmt.Println("--------------")
	//	fmt.Println(v.Name())
	//	//fmt.Println(v.DatabaseTypeName())
	//	fmt.Printf("%#v",v.ScanType().Name())
	//	//fmt.Println(v.ScanType().String())
	//	//fmt.Println(v.Nullable())
	//	fmt.Println("--------------")
	//}
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
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
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

//查询删除插入更新操作
//返回影响的行数
func QRUDExec(sqlStr string, args ...interface{}) (affectNum int64, affectId int64, err error) {
	//自动释放链接
	result, err := gmysql.Con.Exec(sqlStr, args...)
	if err != nil {
		return 0, 0, err
	} else {
		affectNum, err = result.RowsAffected()
		affectId, err = result.LastInsertId()
		if err != nil {
			return 0, 0, err
		}
		return affectNum, affectId, nil
	}
}

//事务的操作
func TxQRUDExec(tx *sql.Tx, sqlStr string, args ...interface{}) (num int64, id int64, err error) {
	//自动释放链接
	result, err := tx.Exec(sqlStr, args...)
	if err != nil {
		return 0, 0, err
	} else {
		num, err = result.RowsAffected()
		id, err = result.LastInsertId()
		if err != nil {
			return 0, 0, err
		}
		return num, id, nil
	}
}
