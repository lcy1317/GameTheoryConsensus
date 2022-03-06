package StartTest

import (
	bolt "bbolt"
	"colorout"
	"log"
	"os"
)

const testBucket = "test"
const testKey = "Luochengyu"
const testValue = "Mengyiyun"

func IfBoltDBExist(dbFile string) error {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return err
	}
	return nil
}
func BoltDBCreate(dbFile string) error {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close() // 及时关闭数据库
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(testBucket))
		if err != nil {
			log.Fatalf(colorout.Red("创建测试Bucket出错:")+"%s", err.Error())
			return err
		}
		if err = bucket.Put([]byte(testKey), []byte(testValue)); err != nil {
			log.Fatalf(colorout.Red("测试Bucket存放数据错误:")+"%s", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf(colorout.Red("更新数据库错误")+"%s", err.Error())
	}
	return nil
}

func BoltDBReadTest(dbFile string) error {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close() // 及时关闭数据库

	// 测试读取刚才的数据
	err = db.View(func(tx *bolt.Tx) error {
		//找到柜子
		bucket := tx.Bucket([]byte(testBucket))
		//找东西
		val := bucket.Get([]byte(testKey))
		log.Printf(colorout.Green("获取存在Key的值:")+"%s", val)
		val = bucket.Get([]byte("hello"))
		log.Printf(colorout.Yellow("获取不存在Key的值:")+"%s", val)
		return nil
	})
	if err != nil {
		log.Fatalf(colorout.Red("数据库读取错误:")+"%s", err.Error())
	}
	return nil
}
