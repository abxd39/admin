package backstage

import (
	"admin/app/models"
	"admin/errors"
	"admin/utils"
	"fmt"
	"strings"
)

// 后台配置
type Config struct {
	models.BaseModel `xorm:"-"`
	Name             string `xorm:"name pk" json:"name"`
	Value            string `xorm:"value" json:"value"`
}

// 表名
func (*Config) TableName() string {
	return "config"
}

const (
	CONFIG_SITE = "site" // 基础配置
	CONFIG_SMS  = "sms"  // 短信配置
	CONFIG_KEFU = "kefu" // 客服配置
)

// 获取配置
func (c *Config) Get(name string) (*Config, error) {
	engine := utils.Engine_backstage
	config := new(Config)
	has, err := engine.Id(name).Get(config)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	if !has {
		return nil, errors.NewNormal("配置不存在")
	}

	return config, nil
}

// 设置配置
func (c *Config) Set(config *Config) error {
	// 转义json字符串里包含的引号
	config.Value = strings.Replace(config.Value, `"`, `\"`, -1)
	config.Value = strings.Replace(config.Value, `'`, `\'`, -1)

	// 开始写入，已存在就更新
	engine := utils.Engine_backstage
	_, err := engine.Exec(fmt.Sprintf("INSERT INTO %s (name, value) VALUES ('%s', '%s') ON DUPLICATE KEY UPDATE value='%[3]s'", c.TableName(), config.Name, config.Value))
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}
