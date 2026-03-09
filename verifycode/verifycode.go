package verifycode

import (
	"Monitoring-Pressure/dao/redis"
	"Monitoring-Pressure/util"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const VerifyCodeExpiration = 2 * time.Minute

// GenerateCode 生成验证码
func GenerateCode(telephone string) (string, error) {
	// 生成6位数的随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 将验证码存储在内存、数据库或其他存储介质中
	err := storeCode(telephone, code)
	if err != nil {
		return "", err
	}

	// 发送验证码短信
	content := fmt.Sprintf("您的验证码是：%s。请不要把验证码泄露给其他人。", code)
	err = sendSMS(telephone, content)
	if err != nil {
		return "", err
	}

	return code, nil
}

// VerifyCode 验证用户输入的验证码是否正确
func VerifyCode(telephone, inputCode string) bool {
	// 从存储中获取正确的验证码
	storedCode, err := getStoredCode(telephone)
	if err != nil {
		return false
	}
	// 比较用户输入的验证码和存储的验证码
	return storedCode == inputCode
}

// 将验证码存储在内存、数据库或其他存储介质中
func storeCode(telephone string, code string) error {
	if redis.RedisClient == nil {
		return fmt.Errorf("redis客户端未初始化")
	}
	return redis.RedisClient.Set(telephone, code, VerifyCodeExpiration).Err()
}

// 从存储中获取正确的验证码
func getStoredCode(telephone string) (string, error) {
	code, err := redis.RedisClient.Get(telephone).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}

// 发送验证码短信
func sendSMS(mobile, content string) error {
	v := url.Values{}
	_now := strconv.FormatInt(time.Now().Unix(), 10)
	_account := "C29894121"
	_password := "615f2126ef7af443c9c334e37f5b9885"
	v.Set("account", _account)
	v.Set("password", util.GetMD5String(_account+_password+mobile+content+_now))
	v.Set("mobile", mobile)
	v.Set("content", content)
	v.Set("time", _now)

	body := strings.NewReader(v.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://106.ihuyi.com/webservice/sms.php?method=Submit&format=json", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}
