package resources

import (
	_ "embed"
	"encoding/base64"
)

// logoPNG 保存随插件编译的官方静态图标，运行时不再访问外部站点。
//
//go:embed logo.png
var logoPNG []byte

// LogoDataURI 供 CPA 直接展示内置图标，并兼容反向代理前缀及未加载资源路由的场景。
var LogoDataURI = "data:image/png;base64," + base64.StdEncoding.EncodeToString(logoPNG)

// OAuthPage 是完成宿主登录会话所需的浏览器授权引导页。
//
//go:embed oauth_page.html
var OAuthPage string
