package mgr

// GProject 项目
type GProject struct {
	ID          string
	ProjectName string
	Owner       string
}

// GUser 用户
type GUser struct {
	UserID    string
	UserName  string
	LoginName string
	Pwd       string
	UserKey   string
	OpenID    string
	Email     string
}

// GDataSourceURL 数据源配置
type GDataSourceURL struct {
	ID        string
	DbType    string
	DbURL     string
	UserName  string
	Pwd       string
	ProjectID string
	DbAlias   string
}

// Gids 数据源
type Gids struct {
	ID        string
	Meta      string
	ProjectID string
}

// GService 服务
type GService struct {
	ID          string
	BodyType    string
	ServiceType string
	NameSpace   string
	Enabled     int
	MsgLog      int
	Security    int
	Meta        string
	ProjectID   string
	Context     string
}
