package logger

// Config 日志配置
type Config struct {
	Out   string `json:"out,omitempty"`
	Level string `json:"level,omitempty"`
	File  *File  `json:"file,omitempty"`
	Type  string `json:"type,omitempty"`
	// 默认：0-default, 1-simple
	Model int32 `json:"model,omitempty"`
}

func (x *Config) GetOut() string {
	if x != nil {
		return x.Out
	}
	return ""
}

func (x *Config) GetLevel() string {
	if x != nil {
		return x.Level
	}
	return ""
}

func (x *Config) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

func (x *Config) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Config) GetModel() int32 {
	if x != nil {
		return x.Model
	}
	return 0
}

type File struct {
	Path       string `json:"path,omitempty"`
	FileName   string `json:"file_name,omitempty"`
	MaxSize    int32  `json:"max_size,omitempty"`
	MaxBackups int32  `json:"max_backups,omitempty"`
	MaxAge     int32  `json:"max_age,omitempty"`
	Compress   bool   `json:"compress,omitempty"`
}

func (x *File) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *File) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *File) GetMaxSize() int32 {
	if x != nil {
		return x.MaxSize
	}
	return 0
}

func (x *File) GetMaxBackups() int32 {
	if x != nil {
		return x.MaxBackups
	}
	return 0
}

func (x *File) GetMaxAge() int32 {
	if x != nil {
		return x.MaxAge
	}
	return 0
}

func (x *File) GetCompress() bool {
	if x != nil {
		return x.Compress
	}
	return false
}
