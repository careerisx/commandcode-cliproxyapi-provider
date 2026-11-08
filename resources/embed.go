package resources

import _ "embed"

// OAuthPage 是完成宿主登录会话所需的浏览器授权引导页。
//
//go:embed oauth_page.html
var OAuthPage string
