package code

import "net/http"

var (

	// 公共业务错误，前缀为：100101

	// 通用 00
	cc = &CCode{
		PT: CommonProduceType,
		ST: IPAMServiceType,
		MT: CommonModuleType,
		ET: CommonErrType,
	}
	// ErrCommon 通用错误
	ErrCommon = cc.Register(http.StatusInternalServerError, 1, "通用错误")
	// ErrOpCancel 操作取消
	ErrOpCancel = cc.Register(http.StatusBadRequest, 2, "操作取消")
	// ErrOpNotAllowed 操作不被允许
	ErrOpNotAllowed = cc.Register(http.StatusBadRequest, 3, "操作不被允许")

	// 参数 01
	cp = &CCode{
		PT: CommonProduceType,
		ST: IPAMServiceType,
		MT: CommonModuleType,
		ET: ParamErrType,
	}
	// ErrParamInvalid 参数无效错误
	ErrParamInvalid = cp.Register(http.StatusBadRequest, 1, "参数无效错误")
	// ErrTypeCast 类型转换错误
	ErrTypeCast = cp.Register(http.StatusInternalServerError, 2, "类型转换错误")

	// 网络 02
	cn = &CCode{
		PT: CommonProduceType,
		ST: IPAMServiceType,
		MT: CommonModuleType,
		ET: NetworkErrType,
	}
	// ErrNetworkException 网络异常
	ErrNetworkException = cn.Register(http.StatusInternalServerError, 1, "网络异常")
	// ErrRequestFail 请求失败
	ErrRequestFail = cn.Register(http.StatusBadRequest, 2, "请求失败")

	// 数据库 03
	cd = &CCode{
		PT: CommonProduceType,
		ST: IPAMServiceType,
		MT: CommonModuleType,
		ET: DatabaseErrType,
	}
	// ErrOpDB 操作数据库错误
	ErrOpDB = cd.Register(http.StatusInternalServerError, 1, "操作数据库错误")
	// ErrDataExist 数据已存在
	ErrDataExist = cd.Register(http.StatusBadRequest, 2, "数据已存在")
	// ErrDataNotExist 数据不存在
	ErrDataNotExist = cd.Register(http.StatusBadRequest, 3, "数据不存在")

	// 文件 04
	cf = &CCode{
		PT: CommonProduceType,
		ST: IPAMServiceType,
		MT: CommonModuleType,
		ET: FileErrType,
	}
	// ErrOpFile 操作文件错误
	ErrOpFile = cf.Register(http.StatusInternalServerError, 1, "操作文件错误")
	// ErrFileExt 不支持该文件类型
	ErrFileExt = cf.Register(http.StatusInternalServerError, 2, "不支持该文件类型")
	// ErrFileInvalid 非法文件，无法解析
	ErrFileInvalid = cf.Register(http.StatusInternalServerError, 3, "非法文件，无法解析")
	// ErrFileWriteError 文件写入异常
	ErrFileWriteError = cf.Register(http.StatusInternalServerError, 4, "文件写入异常")
	// ErrPermissionDenied 权限不足
	ErrPermissionDenied = cf.Register(http.StatusForbidden, 5, "权限不足")
	// ErrImportFileType 导入文件类型错误
	ErrImportFileType = cf.Register(http.StatusBadRequest, 6, "导入文件类型错误")
	// ErrImportFileSuffix 仅支持.cvs结尾的文件
	ErrImportFileSuffix = cf.Register(http.StatusBadRequest, 7, "仅支持.cvs结尾的文件")
	// ErrFileEncoding 不支持当前文件编码
	ErrFileEncoding = cf.Register(http.StatusBadRequest, 8, "不支持当前文件编码")

	// 业务 05

	// 其他 06
	cr = &CCode{
		PT: CommonProduceType,
		ST: IPAMServiceType,
		MT: CommonModuleType,
		ET: OtherErrType,
	}
	// ErrOpRedis 操作Redis错误
	ErrOpRedis = cr.Register(http.StatusInternalServerError, 1, "操作Redis错误")
)
