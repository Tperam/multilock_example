/*
 * @Author: Tperam
 * @Date: 2022-04-28 23:40:58
 * @LastEditTime: 2022-04-29 00:13:36
 * @LastEditors: Tperam
 * @Description:
 * @FilePath: \multilock_example\main.go
 */
package main

import (
	"fmt"
)

func main() {

}

func BusinessFunc(insertData []string) error {
	// 假装插入数据库
	for i := range insertData {
		fmt.Println("insert data", insertData[i])
	}
	return nil
}
