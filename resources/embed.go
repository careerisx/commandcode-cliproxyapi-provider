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

// QuotaPage 提供直接输入 API key 的内置账户与额度页面。
//
//go:embed quota_page.html
var QuotaPage string
