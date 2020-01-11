package mgr

type GProject struct {
	Id          string
	ProjectName string
	Owner       string
}

type GUser struct {
	UserId    string
	UserName  string
	LoginName string
	Pwd       string
	UserKey   string
	OpenId    string
	Email     string
}

type GDataSourceURL struct {
	Id        string
	DbType    string
	DbUrl     string
	UserName  string
	Pwd       string
	ProjectId string
	DbAlias   string
}

type Gids struct {
	Id        string
	Meta      string
	ProjectId string
}

type GService struct {
	Id          string
	BodyType    string
	ServiceType string
	NameSpace   string
	Enabled     int
	MsgLog      int
	Security    int
	Meta        string
	ProjectId   string
	Context     string
}
