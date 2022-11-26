package response

import "mybox/model/system"

type RepositoryResponse struct {
	Repository system.SysUserRepository `json:"repository"`
}
