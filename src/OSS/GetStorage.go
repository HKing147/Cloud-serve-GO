package OSS

import (
	"fmt"
)

func GetStorage() (int64, error) {
	client := Client()
	stat, err := client.GetBucketStat(BucketName)
	if err != nil {
		return 0, err
	}
	// 获取Bucket的总存储量，单位为字节。
	fmt.Println("Bucket Stat Storage:", stat.Storage)
	return stat.Storage, nil
}
