package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"
)

var (
	// from http://w3c.github.io/html-reference/datatypes.html#form.data.emailaddress
	emailReg = regexp.MustCompile("^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")
)

func IsEmail(s string) bool {
	return emailReg.MatchString(s)
}

// 验证码
func IsNumCode(s string, length int) bool {
	r := regexp.MustCompile(fmt.Sprintf(`^\d{%d}$`, length))
	return r.MatchString(s)
}

func IsString(s string, min, max int) bool {
	return utf8.RuneCountInString(s) >= min && utf8.RuneCountInString(s) <= max
}

// 车牌号
func IsPlateNo(s string) bool {
	r := regexp.MustCompile(`^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领A-Z]{1}[A-Z]{1}[A-Z0-9]{4}[A-Z0-9挂学警港澳]{1}$`)
	return r.MatchString(s)
}

func IsTimeYM(s string) bool {
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return false
	}

	return !t.IsZero()
}

func IsTimeYMD(s string) bool {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return false
	}

	return !t.IsZero()
}

func IsDatetime(s string) bool {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return false
	}

	return !t.IsZero()
}

// 公民身份号码是特征组合码，由十七位数字本体码和一位校验码组成
// 排列顺序从左至右依次为：六位数字地址码，八位数字出生日期码，三位数字顺序码和一位校验码
// 顺序码的奇数分配给男性，偶数分配给女性
func IsCardId(id string) (int8, string, bool) {
	tmp := []byte(id)
	if len(tmp) != 18 {
		return 0, "", false
	}

	weight := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}      //十七位数字本体码权重
	validate := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'} //mod11,对应校验码字符值

	mode := tmp[17]
	tmp = tmp[:17]

	sum := 0
	var gender int8 = 1
	for i := 0; i < len(tmp); i++ {
		n, err := strconv.Atoi(string(tmp[i]))
		if err != nil {
			return 0, "", false
		}
		if i == 16 && n%2 == 0 {
			gender = 2
		}

		sum += n * weight[i]
	}

	birthday := fmt.Sprintf("%s-%s-%s",
		string(tmp[6:10]), string(tmp[10:12]), string(tmp[12:14]))

	return gender, birthday, validate[sum%11] == mode
}

// 密码, 大小写字母、数字、符号至少包含2种
// golang unsupported Perl syntax: `(?=` `(?!`
// r := regexp.MustCompile(`^(?![A-Za-z]+$)(?!\d+$)(?![\W_]+$)\S{6,16}$`)
// r := regexp.MustCompile(`^(?=.*[a-zA-Z0-9].*)(?=.*[a-zA-Z\W].*)(?=.*[0-9\W].*).{6,16}$`)
func IsPasswordV2(s string) bool {
	if !IsString(s, 6, 16) {
		return false
	}

	rL := regexp.MustCompile(`.*[0-9].*`)
	rN := regexp.MustCompile(`.*[a-zA-Z].*`)
	rS := regexp.MustCompile(`.*\W.*`)

	var flag, tmp byte

	if rL.MatchString(s) {
		flag |= 1
	}
	if rN.MatchString(s) {
		flag |= 2
	}
	if rS.MatchString(s) {
		flag |= 4
	}

	var i uint
	for i = 0; i < 3; i++ {
		tmp += (flag >> i) & 1 // or use bits.OnesCount
	}

	return tmp > 1
}
