package model

import (
	"gorm.io/gorm"
)

// 问题的测试用例
type JudgeCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// 问题的配置信息
type JudgeConfig struct {
	MemoryLimit int   `json:"memoryLimit"`
	StackLimit  int   `json:"stackLimit"`
	TimeLimit   int64 `json:"timeLimit"`
}

const TableNameQuestion = "question"

// Question 题目
type Question struct {
	gorm.Model
	Identity    string `gorm:"column:identity;type:varchar(36);not null;index:idx_identity,priority:1;comment:唯一ID" json:"identity"` // 唯一ID
	JudgeCase   string `gorm:"column:judgeCase;type:text;comment:判题用例（json 数组）" json:"judgeCase"`                                    // 判题用例（json 数组）
	JudgeConfig string `gorm:"column:judgeConfig;type:text;comment:判题配置（json 对象）" json:"judgeConfig"`                                // 判题配置（json 对象）
}

// TableName Question's table name
func (*Question) TableName() string {
	return TableNameQuestion
}

/*

id  bigint auto_increment comment 'id' primary key
identity    varchar(36)   not null comment '唯一ID'
judgeCase   text          null comment '判题用例（json 数组）'
judgeConfig text          null comment '判题配置（json 对象）'
created_at  datetime      not null comment '创建时间'
updated_at  datetime      not null comment '更新时间'
deleted_at  datetime      null comment '删除时间'
*/
