package verifycode

import (
	"Monitoring-Pressure/dao/redis"
	"Monitoring-Pressure/util"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	VerifyCodeExpiration = 2 * time.Minute
	SendInterval         = 60 * time.Second
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ===================== 内存兜底存储 =====================
// 当 Redis 未初始化时使用，行为与 Redis 一致：带 TTL 的 KV
type memEntry struct {
	value    string
	expireAt time.Time
}

var (
	memMu    sync.Mutex
	memStore = map[string]memEntry{}
)

func memSet(key, val string, ttl time.Duration) {
	memMu.Lock()
	defer memMu.Unlock()
	memStore[key] = memEntry{value: val, expireAt: time.Now().Add(ttl)}
}

func memGet(key string) (string, bool) {
	memMu.Lock()
	defer memMu.Unlock()
	e, ok := memStore[key]
	if !ok {
		return "", false
	}
	if time.Now().After(e.expireAt) {
		delete(memStore, key)
		return "", false
	}
	return e.value, true
}

func memDel(key string) {
	memMu.Lock()
	defer memMu.Unlock()
	delete(memStore, key)
}

// ===================== 公共接口 =====================

// GenerateCode 生成验证码
func GenerateCode(telephone string) (string, error) {
	if telephone == "" {
		return "", fmt.Errorf("telephone empty")
	}

	if err := checkSendInterval(telephone); err != nil {
		return "", err
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	if err := storeCode(telephone, code); err != nil {
		return "", err
	}
	if err := storeSendInterval(telephone); err != nil {
		return "", err
	}

	// 始终将验证码打印到服务端控制台，便于开发期查看
	log.Printf("[verifycode] telephone=%s code=%s (有效期 %s)", telephone, code, VerifyCodeExpiration)

	// 短信发送是 best-effort，失败也不阻断接口（已经把 code 存好了）
	go func() {
		content := fmt.Sprintf("您的验证码是：%s。请不要把验证码泄露给其他人。", code)
		if err := sendSMS(telephone, content); err != nil {
			log.Printf("[verifycode] 短信发送失败 telephone=%s err=%v", telephone, err)
		}
	}()

	return code, nil
}

// VerifyCode 校验验证码
func VerifyCode(telephone, inputCode string) bool {
	stored, err := getStoredCode(telephone)
	if err != nil {
		return false
	}
	if stored != inputCode {
		return false
	}
	_ = deleteCode(telephone)
	return true
}

// ===================== 存储抽象（Redis 优先，内存兜底） =====================

func storeCode(telephone, code string) error {
	key := getCodeKey(telephone)
	if redis.RedisClient != nil {
		return redis.RedisClient.Set(key, code, VerifyCodeExpiration).Err()
	}
	memSet(key, code, VerifyCodeExpiration)
	return nil
}

func getStoredCode(telephone string) (string, error) {
	key := getCodeKey(telephone)
	if redis.RedisClient != nil {
		return redis.RedisClient.Get(key).Result()
	}
	v, ok := memGet(key)
	if !ok {
		return "", fmt.Errorf("code not found")
	}
	return v, nil
}

func deleteCode(telephone string) error {
	key := getCodeKey(telephone)
	if redis.RedisClient != nil {
		return redis.RedisClient.Del(key).Err()
	}
	memDel(key)
	return nil
}

func checkSendInterval(telephone string) error {
	key := getSendLimitKey(telephone)
	if redis.RedisClient != nil {
		if _, err := redis.RedisClient.Get(key).Result(); err == nil {
			return fmt.Errorf("发送过于频繁，请稍后再试")
		}
		return nil
	}
	if _, ok := memGet(key); ok {
		return fmt.Errorf("发送过于频繁，请稍后再试")
	}
	return nil
}

func storeSendInterval(telephone string) error {
	key := getSendLimitKey(telephone)
	if redis.RedisClient != nil {
		return redis.RedisClient.Set(key, "1", SendInterval).Err()
	}
	memSet(key, "1", SendInterval)
	return nil
}

func getCodeKey(telephone string) string      { return "verify_code:" + telephone }
func getSendLimitKey(telephone string) string { return "verify_send_limit:" + telephone }

// ===================== 短信发送（best-effort） =====================

func sendSMS(mobile, content string) error {
	account := os.Getenv("SMS_ACCOUNT")
	password := os.Getenv("SMS_PASSWORD")
	if account == "" || password == "" {
		return fmt.Errorf("短信配置缺失（SMS_ACCOUNT / SMS_PASSWORD）")
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
