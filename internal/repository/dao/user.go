package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateUser  = errors.New("用户冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error

	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const duplicateErr uint16 = 1062 // 数据库 email 冲突
		if me.Number == duplicateErr {
			return ErrDuplicateUser
		}
	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDao) UpdateByUserId(ctx context.Context, entity User) error {
	// 使用 dao.db.WithContext(ctx) 的目的是为了实现上下文控制。这个机制允许你在处理数据库请求时，如果上下文 ctx 被取消（例如由于超时或其它原因），则可以取消正在进行的数据库操作。
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
		// 就算是有的字段没有也可以更新吗？因为 User 结构体里有 Password 和 Email 之类的字段，但是这里没有传入，只需要传入本次需要修改的？
		Updates(map[string]any{
			"utime":    time.Now().UnixMilli(),
			"nickname": entity.Nickname,
			"birthday": entity.Birthday,
			"about_me": entity.AboutMe,
		}).Error
}

func (dao *UserDao) FindById(ctx context.Context, uid int64) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&res).Error
	return res, err
}

func (dao *UserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&res).Error
	return res, err
}

type User struct {
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
	Id       int64  `gorm:"PrimaryKey,autoIncrement"`
	Nickname string `gorm:"type=varchar(128)"`
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`
}
