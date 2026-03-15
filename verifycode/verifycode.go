package verifycode

import (
	"Monitoring-Pressure/dao/redis"
	"Monitoring-Pressure/util"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	VerifyCodeExpiration = 2 * time.Minute
	SendInterval         = 60 * time.Second
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateCode 生成验证码
func GenerateCode(telephone string) (string, error) {
	if telephone == "" {
		return "", fmt.Errorf("telephone empty")
	}

	// 发送频率限制
	if err := checkSendInterval(telephone); err != nil {
		return "", err
	}

	// 生成6位验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 存储验证码
	if err := storeCode(telephone, code); err != nil {
		return "", err
	}

	// 记录发送时间
	if err := storeSendInterval(telephone); err != nil {
		return "", err
	}

	// 发送短信
	content := fmt.Sprintf("您的验证码是：%s。请不要把验证码泄露给其他人。", code)
	if err := sendSMS(telephone, content); err != nil {
		return "", err
	}

	return code, nil
}

// VerifyCode 验证验证码
func VerifyCode(telephone, inputCode string) bool {
	storedCode, err := getStoredCode(telephone)
	if err != nil {
		return false
	}

	if storedCode != inputCode {
		return false
	}

	// 验证成功后删除，防止重复使用
	_ = deleteCode(telephone)
	return true
}

func storeCode(telephone, code string) error {
	if redis.RedisClient == nil {
		return fmt.Errorf("redis客户端未初始化")
	}
	return redis.RedisClient.Set(getCodeKey(telephone), code, VerifyCodeExpiration).Err()
}

func getStoredCode(telephone string) (string, error) {
	if redis.RedisClient == nil {
		return "", fmt.Errorf("redis客户端未初始化")
	}
	return redis.RedisClient.Get(getCodeKey(telephone)).Result()
}

func deleteCode(telephone string) error {
	if redis.RedisClient == nil {
		return fmt.Errorf("redis客户端未初始化")
	}
	return redis.RedisClient.Del(getCodeKey(telephone)).Err()
}

func checkSendInterval(telephone string) error {
	if redis.RedisClient == nil {
		return fmt.Errorf("redis客户端未初始化")
	}

	key := getSendLimitKey(telephone)
	_, err := redis.RedisClient.Get(key).Result()
	if err == nil {
		return fmt.Errorf("发送过于频繁，请稍后再试")
	}
	return nil
}

func storeSendInterval(telephone string) error {
	if redis.RedisClient == nil {
		return fmt.Errorf("redis客户端未初始化")
	}
	return redis.RedisClient.Set(getSendLimitKey(telephone), "1", SendInterval).Err()
}

func getCodeKey(telephone string) string {
	return "verify_code:" + telephone
}

func getSendLimitKey(telephone string) string {
	return "verify_send_limit:" + telephone
}

// 发送验证码短信
func sendSMS(mobile, content string) error {
	account := os.Getenv("SMS_ACCOUNT")
	password := os.Getenv("SMS_PASSWORD")

	if account == "" || password == "" {
		return fmt.Errorf("短信配置缺失")
	}

	v := url.Values{}
	nowStr := strconv.FormatInt(time.Now().Unix(), 10)

	v.Set("account", account)
	v.Set("password", util.GetMD5String(account+password+mobile+content+nowStr))
	v.Set("mobile", mobile)
	v.Set("content", content)
	v.Set("time", nowStr)

	body := strings.NewReader(v.Encode())
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", "http://106.ihuyi.com/webservice/sms.php?method=Submit&format=json", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	return err
}
