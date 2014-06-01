package models

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io"
	"time"

	"github.com/revel/revel"
)

type User struct {
	Id              int64  `xorm:"pk"`
	Name            string `xorm:"varchar(16)"`
	Email           string `xorm:"varchar(256)"`
	Password        string `xorm:"-"`
	ConfirmPassword string `xorm:"-"`
	HashedPassword  string `xorm:"varchar(64)"` //SHA256 输出32个字节 64长度的字符串
	Salt            string `xorm:"varchar(64)"`
	Code            string `xorm:"varchar(32)"`
	CodeCreatedTime time.Time
	JoinedAt        time.Time
	Problems        int64 //Number of solved problems
}

func (user *User) Validate(v *revel.Validation) {
	v.Required(user.Name)
	v.Required(user.Email)
	valid := v.Email(user.Email)
	if user.Password != user.ConfirmPassword {
		v.Errors = append(v.Errors, &revel.ValidationError{Message: "两次密码不一致", Key: "user.Password"})
	}

	if user.HasName() {
		err := &revel.ValidationError{
			Message: "该用户名已经被注册",
			Key:     "user.Name",
		}
		v.Errors = append(v.Errors, err)
	}

	if valid.Ok {
		if user.HasEmail() {
			err := &revel.ValidationError{
				Message: "该邮箱已经被用于注册",
				Key:     "user.Email",
			}
			v.Errors = append(v.Errors, err)
		}
	}
}

func (user *User) HasName() bool {
	u := new(User)
	has, _ := engine.Where("name = ?", user.Name).Get(u)
	if has {
		return true
	}
	return false
}

func (user *User) HasEmail() bool {
	u := new(User)
	has, _ := engine.Where("email = ?", user.Email).Get(u)
	if has {
		return true
	}
	return false
}

//验证登陆
func (user *User) LoginOk() bool {
	engine.Get(user)
	if user.Salt == "" {
		return false
	}
	hashedPwd := HashPassword(user.Password, user.Salt)
	fmt.Printf("%s\n", hashedPwd)
	if slowEquals(hashedPwd, user.HashedPassword) {
		return true
	}
	return false
}

// SHA256 加盐 加密
// params  raw password
// return hashed pwd and salt
func GenHashPasswordAndSalt(password string) (string, string) {
	r := make([]byte, sha256.Size)
	_, err := rand.Read(r)
	if err != nil {
		panic(err)
	}
	h := sha256.New()
	io.WriteString(h, password+fmt.Sprintf("%x", r))
	return fmt.Sprintf("%x", h.Sum(nil)), fmt.Sprintf("%x", r)
}

// SHA256 代盐hash
// params raw password, salt
// return hased password
func HashPassword(rawpwd, salt string) string {
	h := sha256.New()
	io.WriteString(h, rawpwd+salt)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//慢相等,使得比较时间恒定
func slowEquals(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func (u *User) Save() bool {
	if u.Password != "" {
		u.HashedPassword, u.Salt = GenHashPasswordAndSalt(u.Password)
	}
	u.JoinedAt = time.Now()
	_, err := engine.Insert(u)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetUserId(name string) (bool, int64) {
	user := new(User)
	has, _ := engine.Where("name = ?", name).Get(user)
	return has, user.Id
}
