package model

type XUSTUser struct {
	Uid   string `gorm:"type:char(44);primaryKey;not null"`         // UID
	Sno   string `gorm:"type:char(11);unique;uniqueIndex;not null"` // 学号
	Email string `gorm:"not null"`                                  // 邮箱
	Name  string `gorm:"not null"`                                  // 姓名
	Time  int    `gorm:"not null;default=1"`                        // 打卡次数
}
