package code

import (
	"fmt"
	"strconv"
)

// ProduceType 产品类型
type ProduceType string

const (
	// CommonProduceType 通用产品
	CommonProduceType = "10"
	// DHCPProduceType DHCP产品
	DHCPProduceType = "13"
)

// ServiceType 服务类型
type ServiceType string

const (
	// CommonServiceType 通用 01
	CommonServiceType = "01"
	// IPAMServiceType IPAM 03
	IPAMServiceType = "03"
)

type ModuleType string

const (
	// CommonModuleType 通用 01
	CommonModuleType = "01"
	// NetworkModuleType 网络 02
	NetworkModuleType = "02"
	// DeviceModuleType 设备 03
	DeviceModuleType = "03"
	// IPModuleType 地址 04
	IPModuleType = "04"
	// ProbeModuleType 探针 05
	ProbeModuleType = "05"
	// TerminalModuleType 终端 06
	TerminalModuleType = "06"
	// AttributeModuleType 自定义属性 07
	AttributeModuleType = "07"
	// ImportExportModuleType 导入导出 08
	ImportExportModuleType = "08"
	// TaskModuleType 任务 09
	TaskModuleType = "09"
	// BusinessType 组织 09
	BusinessType = "10"
)

// ErrType 错误类型
type ErrType string

const (
	// CommonErrType 通用 00
	CommonErrType = "00"
	// ParamErrType 参数 01
	ParamErrType = "01"
	// NetworkErrType 网络 02
	NetworkErrType = "02"
	// DatabaseErrType 数据库 03
	DatabaseErrType = "03"
	// FileErrType 文件 04
	FileErrType = "04"
	// BusinessErrType 业务 05
	BusinessErrType = "05"
	// OtherErrType 其他 06
	OtherErrType = "06"
)

type CCode struct {
	PT ProduceType
	ST ServiceType
	MT ModuleType
	ET ErrType
}

// Register 注册Code
func (cc *CCode) Register(status int, c int, msg string) *Code {
	if c < 0 || c >= 1000 {
		panic("exceed limit")
	}

	s := fmt.Sprintf("%s%s%s%s%03d", cc.PT, cc.ST, cc.MT, cc.ET, c)
	code, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("conversion err: %v", err))
	}

	co := &Code{
		status: status,
		code:   code,
		msg:    msg,
	}
	allCode[co.Code()] = co
	return co
}
