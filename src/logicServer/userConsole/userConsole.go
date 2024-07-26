package userConsole

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"log"
	"myGameDemo/logicServer/msg"
	"sync"
	"time"
)

type UserConsole struct {
	rdb *redis.Client
}

var userConsole *UserConsole
var once sync.Once

func GetUserConsole() *UserConsole {
	once.Do(func() {
		userConsole = &UserConsole{rdb: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // 没有密码，默认值
			DB:       0,  // 默认DB 0
		})}
	})
	return userConsole
}

const (
	TOKEN_EXP_TIME = 10 * time.Second
)

// hashPassword 函数用于生成密码的哈希摘要
func hashPassword(password string) (string, error) {
	// 使用bcrypt算法进行密码加密
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// checkPassword 函数用于验证密码是否匹配
func checkPassword(password, hashedPassword string) bool {
	// 将哈希摘要转换为字节数组
	hashedBytes := []byte(hashedPassword)

	// 使用bcrypt算法验证密码
	err := bcrypt.CompareHashAndPassword(hashedBytes, []byte(password))
	return err == nil
}

func (C *UserConsole) Register(auth UserInfo) (*msg.Res, error) {
	exists, err := C.rdb.SIsMember(context.Background(), "Users", auth.Username).Result()
	if err != nil {
		log.Fatalf("Error checking if username %s exists in set %s: %v", auth.Username, "Users", err)
		return nil, err
	}
	var result msg.Res = msg.Res{}
	if exists {
		result.Code = msg.REGISTERED
		result.Msg = "用户名已注册"
		return &result, nil
	}

	hashPwd, _ := hashPassword(auth.Pwd)
	fields := map[string]interface{}{
		"username": auth.Username, // 用户名字段
		"pwd":      hashPwd,       // 密码字段
	}
	txn := C.rdb.TxPipeline() //创建事务
	txn.SAdd(context.Background(), "Users", auth.Username)
	txn.HSet(context.Background(), "User:"+auth.Username, fields)
	_, _ = txn.Exec(context.Background())

	result.Code = msg.SUCCESS
	result.Msg = "注册成功"
	return &result, nil
}

func (C *UserConsole) Login(auth UserInfo) (*msg.Res, error) {
	var result msg.Res = msg.Res{}
	exists, err := C.rdb.SIsMember(context.Background(), "Users", auth.Username).Result()
	if err != nil {
		return &result, err
	}

	if !exists {
		result.Code = msg.NOUSER
		result.Msg = "用户名不存在"
		return &result, nil
	}

	hashedPassword, _ := C.rdb.HGet(context.Background(), "User:"+auth.Username, "pwd").Result()
	if !checkPassword(auth.Pwd, hashedPassword) {

		result.Code = msg.PWDERR
		result.Msg = "用户名或密码错误"
		return &result, nil
	}

	val, errOnline := C.rdb.Get(context.Background(), "Online:"+auth.Username).Result()
	UUID := uuid.New()
	Session := UUID.String()
	if errors.Is(errOnline, redis.Nil) {
		txn := C.rdb.TxPipeline()
		txn.Set(context.Background(), "Session:"+Session, auth.Username, TOKEN_EXP_TIME)
		txn.Set(context.Background(), "Online:"+auth.Username, Session, TOKEN_EXP_TIME)
		_, _ = txn.Exec(context.Background())
	} else {
		txn := C.rdb.TxPipeline()
		txn.Del(context.Background(), "Session:"+val)
		txn.Del(context.Background(), "Online:"+auth.Username)
		txn.Set(context.Background(), "Session:"+Session, auth.Username, TOKEN_EXP_TIME)
		txn.Set(context.Background(), "Online:"+auth.Username, Session, TOKEN_EXP_TIME)
		_, _ = txn.Exec(context.Background())
	}

	result.Code = msg.SUCCESS
	result.Msg = Session
	return &result, nil
}

func (C *UserConsole) GetOnlineUser() (*msg.Res, error) {
	ctx := context.Background()

	// 使用SCAN命令遍历所有符合条件的key
	var keys []string
	var cursor uint64 = 0
	pattern := "Online:*" // 匹配的前缀模式

	for {
		var err error
		var result []string
		result, cursor, err = C.rdb.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range result {
			keys = append(keys, r[7:])
		}
		if cursor == 0 {
			break
		}
	}
	return &msg.Res{Code: msg.SUCCESS, Msg: "在线玩家列表", Data: keys}, nil
}

func (C *UserConsole) GetUsersList() (*msg.Res, error) {
	members, err := C.rdb.SMembers(context.Background(), "Users").Result()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	userInfo := make(map[string]map[string]interface{})
	for _, m := range members {
		userInfo[m] = make(map[string]interface{})
		userInfo[m]["username"] = m
		userInfo[m]["online"] = false
		exi, err := C.rdb.Exists(context.Background(), "Online:"+m).Result()
		if err != nil {
			return nil, err
		}
		if exi == 1 {
			userInfo[m]["online"] = true
		}
	}
	return &msg.Res{Code: msg.SUCCESS, Msg: "玩家列表", Data: userInfo}, nil

}

func (C *UserConsole) GetUsername(sessionID string) (string, error) {
	result, err := C.rdb.Get(context.Background(), "Session:"+sessionID).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (C *UserConsole) Heart(session SessionInfo) (*msg.Res, error) {
	var result msg.Res = msg.Res{}
	username, err := C.rdb.Get(context.Background(), "Session:"+session.SessionId).Result()
	if errors.Is(err, redis.Nil) {
		result.Code = msg.OUTTIMESESSION
		result.Msg = "会话已过期"
		return &result, nil
	}
	txn := C.rdb.TxPipeline()
	txn.Set(context.Background(), "Session:"+session.SessionId, username, TOKEN_EXP_TIME)
	txn.Set(context.Background(), "Online:"+username, session.SessionId, TOKEN_EXP_TIME)
	_, _ = txn.Exec(context.Background())
	result.Code = msg.SUCCESS
	result.Msg = "会话已更新"
	return &result, nil
}

func (C *UserConsole) SessionCheck(session SessionInfo) (*msg.Res, error) {
	username, err := C.rdb.Get(context.Background(), "Session:"+session.SessionId).Result()
	if errors.Is(err, redis.Nil) {
		return &msg.Res{Code: msg.OUTTIMESESSION, Msg: "会话已过期"}, nil
	}
	return &msg.Res{Code: msg.SUCCESS, Msg: username}, nil
}
