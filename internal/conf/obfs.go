package conf

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/rclone/rclone/fs/config/obscure"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

const obfuscatedPrefix = "___Obfuscated___"

func Obfuscate(str string) string {
	temp := str
	if !strings.HasPrefix(temp, obfuscatedPrefix) {
		temp = obfuscatedPrefix + obscure.MustObscure(temp)
	}
	return temp
}
func Revel(str string) string {
	temp := str
	if strings.HasPrefix(temp, obfuscatedPrefix) {
		temp, _ = strings.CutPrefix(temp, obfuscatedPrefix)
		temp = obscure.MustReveal(temp)
	}
	return temp
}

/*func (c CInfo) Value() (driver.Value, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}*/

// ObfusText is a transparent text encryption type for gorm. the string is decrypted when reading from DB,
// so it's plaintext when used anywhere else
//type ObfusText string

type ObfusText struct {
	text string
}

func (o *ObfusText) GormDataType() string {
	return "ObfusText"
}
func NewObfusText(str string) *ObfusText {
	return &ObfusText{str}
}
func (o ObfusText) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  obfuscatedPrefix + "?",
		Vars: []interface{}{obscure.MustObscure(o.text)},
	}
}

// Scan a value into struct from a database driver
func (o *ObfusText) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New(fmt.Sprint("Failed to parse string:", value))
	}
	if strings.HasPrefix(str, obfuscatedPrefix) {
		str, _ = strings.CutPrefix(str, obfuscatedPrefix)
		str = obscure.MustReveal(str)
	}
	*o = *NewObfusText(str)
	return nil
}

func (o ObfusText) Value() (driver.Value, error) {
	str := Obfuscate(o.text)
	//if !strings.HasPrefix(str, obfuscatedPrefix) {
	//	temp, err := obscure.Obscure(str)
	//	if err != nil {
	//		return nil, err
	//	}
	//	str = obfuscatedPrefix + temp
	//}
	return str, nil
}
func (o ObfusText) String() string {
	//str := string(*o)
	//if strings.HasPrefix(str, obfuscatedPrefix) {
	//	str, _ = strings.CutPrefix(str, obfuscatedPrefix)
	//	str = obscure.MustReveal(str)
	//}
	return o.text
}
