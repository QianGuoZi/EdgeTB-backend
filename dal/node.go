package dal

import "log"

// AddNode 创建节点配置
func AddNode(node Node) (int64, error) {
	result := DB.Model(&Node{}).Create(&node)
	if result.Error != nil {
		log.Printf("[AddNode] 数据库创建节点配置失败")
		return 0, result.Error
	}
	return node.Id, nil
}
