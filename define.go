package five9

type ChangePasswordInfo struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type Org struct {
	OrgID   string `json:"orgId"`
	OrgName string `json:"orgName"`
}

type Auth struct {
	Credentials *PasswordCredentials `json:"passwordCredentials"`
	Org         *Org                 `json:"org,omitempty"`
	AppKey      string               `json:"appKey"`
	Policy      string               `json:"policy"`
}

type PasswordCredentials struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	TenantName string `json:"tenantName,omitempty"`
}

type Token struct {
	TokenID string `json:"tokenId"`
	UserID  string `json:"userId"`
	Context struct {
		FarmID string `json:"farmId"`
	} `json:"context"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	DataCenters []DataCenterInfo `json:"dataCenters"`
}

type DataCenterInfo struct {
	Name      string        `json:"name"`
	Active    bool          `json:"active"`
	APIURLs   []ClusterInfo `json:"apiUrls"`
	LoginURLs []ClusterInfo `json:"loginUrls"`
	UIURLs    []ClusterInfo `json:"uiUrls"`
}

type ClusterInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	RouteKey string `json:"routeKey"`
	Version  string `json:"version"`
}

type StationInfo struct {
	StationType string `json:"stationType"`
}

type MaintenanceNoticeInfo struct {
	ID         string `json:"id"`
	Accepted   bool   `json:"accepted"`
	Annotation string `json:"annotation,omitempty"`
	Text       string `json:"text,omitempty"`
}
