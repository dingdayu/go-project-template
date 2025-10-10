package contextkey

// UserContext 是用于在 context 中传递用户信息的键类型
type UserContext string

var (
	// Email 用户邮箱
	Email UserContext = "email"
	// RealName 用户真实姓名
	RealName UserContext = "real_name"
	// UserName 用户名（登录名）
	UserName UserContext = "user_name"
	// Role 用户角色列表
	Role UserContext = "role"
	// IP 客户端 IP 地址
	IP UserContext = "ip"
)
