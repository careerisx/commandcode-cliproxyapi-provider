# CommandCode Provider

为支持动态插件的 CLIProxyAPI v8 宿主接入 CommandCode，提供模型发现、对话请求转发、账户管理和额度查询。

- 插件 ID：`commandcode-cliproxyapi-provider`
- ProviderID：`commandcode`
- [源码仓库](https://github.com/careerisx/commandcode-cliproxyapi-provider) · [下载 Releases](https://github.com/careerisx/commandcode-cliproxyapi-provider/releases)

## 功能与使用条件

| 项目 | 说明 |
|---|---|
| 支持套餐 | 添加账户时联网核验身份与订阅，仅接受有效且未到期的个人 GOAT / Pro / Max；不支持 Go、未知或已失效套餐 |
| 模型发现 | 从上游动态获取模型目录，默认每 15 分钟刷新；默认使用 `commandcode/上游模型ID` |
| 请求协议 | 支持 Chat Completions、Messages、Responses，以及流式和非流式请求；实际可用模型、协议和能力以上游目录及账户权限为准 |
| 账户接入 | 手动输入 API key，或从宿主 OAuth 登录项进入浏览器授权引导 |
| 账户与额度 | 内置 CommandCode 页面；提供宿主原生额度查询，不支持重置额度 |
| 发布平台 | 仅提供 `linux/amd64` 的 `.so`，不提供 Windows、macOS 或 ARM 安装包 |

需要已启用管理中心、支持动态插件的宿主版本。普通无插件支持的宿主不能加载本插件。宿主还需能够访问 CommandCode 接口；商店安装需能访问商店索引与 GitHub Releases。

## 安装

### 推荐：通过自定义商店安装

将以下配置**合并到宿主现有 `plugins` 节**，保留其他插件配置；已有 `store-sources` 时追加 URL，不要重复定义 `plugins` 根节点：

```yaml
plugins:
  enabled: true
  dir: "./plugins"
  store-sources:
    - "https://raw.githubusercontent.com/careerisx/cpa-plusins/main/registry.json"
  configs:
    commandcode-cliproxyapi-provider:
      enabled: true
```

1. 保存宿主配置，按宿主要求重新加载或重启。
2. 在管理中心的插件商店找到 **CommandCode Provider**，安装适配当前平台的版本。
3. 确认插件已启用，再进入 **CommandCode** 资源菜单添加账户。

添加商店不等于安装插件。安装依赖仓库中实际存在的 Release 与匹配平台的资产；请以商店和 Releases 中可用的版本为准。

### 手动部署

1. 在 [Releases](https://github.com/careerisx/commandcode-cliproxyapi-provider/releases) 选择版本，下载 `commandcode-cliproxyapi-provider_<版本号>_linux_amd64.zip` 和该 Release 的 `checksums.txt`，核对 ZIP 的 SHA-256。
2. 将 ZIP 解压到新的空目录，再把其中的 `commandcode-cliproxyapi-provider.so` 放入宿主 `plugins.dir` 指定的目录。不要把 ZIP 当作插件文件，也不需要复制头文件。
3. 合并上面的配置，确保总开关 `plugins.enabled` 和 `plugins.configs.commandcode-cliproxyapi-provider.enabled` 均为 `true`。仅手动部署时可以不配置 `store-sources`。
4. 重新加载或重启宿主，确认插件已加载。

容器部署时，插件文件必须位于**容器内**配置指向的目录，并具有宿主进程可读权限。升级前备份配置与认证文件，停止宿主后替换 `.so`，再启动，避免覆盖正在加载的共享库。

## 添加与管理账户

先在同一浏览器、同一地址的管理中心登录，并勾选**记住密码**。插件页面使用已保存的管理凭据，不要把管理密码放进 URL。

### 已有 API key

1. 打开资源菜单中的 **CommandCode**（页面标题为“CommandCode 账户与额度”）。
2. 输入 CommandCode API key，点击**保存**。
3. 插件会联网核验账户邮箱和有效订阅，通过后交由宿主保存为认证文件。
4. 点击**刷新账户列表**确认账户出现。启用、禁用和删除账户均在宿主的**认证文件**页操作。

同一个 key 即使对应账户已禁用，也不能重复添加；应管理现有认证文件。若提示保存结果未确认，请先刷新认证文件列表，不要反复提交。

### 浏览器授权

此方式无需提前手动创建 key，但**不是自动 OAuth 回调**，最终仍使用 API key，也不提供自动续期。

1. 从宿主的 **CommandCode OAuth 登录项**开始，打开插件提供的浏览器授权页面。不要直接访问没有会话的授权页。
2. 点击**打开官方授权页**，登录个人账户并点击“授权”。
3. 官方自动传递失败后进入备用页，复制生成的 **API key**，不是地址栏网址。
4. 回到插件页面粘贴 key，点击**验证并提交**。
5. 等待页面确认 CPA 已完成保存并核对认证文件，随后再关闭页面。

登录会话有效期为 10 分钟。当前页面注明仅在已验证的 Chrome 默认设置下使用；不要解除浏览器对端口 `6000` 的限制。交付给宿主前可以取消，交付后不能保证取消会撤回保存。结果不确定时使用**重新确认状态**并查看认证文件，勿再次添加同一 key。

## 查看额度与调用模型

在 **CommandCode** 页面可查看账户邮箱、认证文件名、启用状态、凭据来源，以及上游返回的 5 小时、每周和月度 credits 用量、剩余额度与重置时间。

- **刷新账户列表**只更新账户清单；每张卡片的 **Refresh** 用于查询该账户的实时额度。
- 仅没有邮箱的账户会自动查询一次额度，并尝试将邮箱补充到认证文件；其他账户按需手动刷新。
- 页面会在当前浏览器会话中缓存查询结果，请结合更新时间判断；查询失败不代表额度为零。
- 支持相应功能的宿主也可使用原生额度入口；本插件不提供额度重置操作。

资源页路径为 `/v0/resource/plugins/commandcode-cliproxyapi-provider/quota`；有反向代理前缀时需带上该前缀，优先从资源菜单进入。

添加并启用账户后，在宿主模型列表选择实际出现的 `commandcode/上游模型ID`。客户端继续使用**宿主提供的 API 地址和客户端访问密钥**，不要把 CommandCode 上游 key 当成宿主客户端密钥。首次安装尚未添加账户时，模型列表可能为空；不保证所有上游模型都具备相同能力。

## 可选配置

默认配置即可连接官方服务。下面的字段均直接放在 `plugins.configs.commandcode-cliproxyapi-provider` 下，与该项的 `enabled` 同级；按需添加，不要另建嵌套 `config` 节。

| 字段 | 默认值 | 用途与限制 |
|---|---|---|
| `base-url` | `https://api.commandcode.ai/provider/v1` | 上游推理接口基础地址；账户核验与额度接口使用该地址的同源账户服务 |
| `catalog-url` | `base-url` 去掉末尾 `/` 后追加 `/models` | 模型目录地址 |
| `model-prefix.enabled` | `true` | 是否为公开模型 ID 添加前缀 |
| `model-prefix.value` | `commandcode` | 前缀须以字母或数字开头，仅含字母、数字、`.`、`_`、`-` |
| `catalog.refresh-interval` | `15m` | 模型目录刷新间隔，至少 `1m` |
| `catalog.stale-while-unavailable` | `true` | 刷新失败时继续使用已有目录；首次尚无成功目录时仍为空 |
| `protocols.chat-completions` | `true` | 允许使用 Chat Completions 上游路由 |
| `protocols.messages` | `true` | 允许使用 Messages 上游路由 |
| `protocols.responses` | `true` | 允许使用 Responses 上游路由 |
| `request-timeout` | `5m` | 推理请求超时，必须为正时长；交互式账户核验和额度查询另有限时，最长 30 秒 |
| `max-response-bytes` | `67108864` | 响应大小限制，单位字节，必须为正数 |
| `allow-http` | `false` | 默认仅允许 HTTPS；仅可信环境确需 HTTP 时开启 |
| `route-overrides` | 无 | 按上游模型 ID 覆盖路由，见下例 |

`base-url` 和 `catalog-url` 不允许带用户名、密码、查询参数或片段。自定义地址可能接收账户密钥，务必仅使用可信服务。

仅在明确知道上游支持的协议及路径时配置路由覆盖。例如以下片段需合并到现有插件配置中，并将占位模型 ID 替换为模型目录中的真实上游 ID：

```yaml
plugins:
  configs:
    commandcode-cliproxyapi-provider:
      enabled: true
      route-overrides:
        "实际上游模型ID":
          protocol: "chat-completions"
          endpoint: "/v1/chat/completions"
```

`protocol` 仅接受 `chat-completions`、`messages`、`responses`，`endpoint` 必须以 `/` 开头。覆盖配置不代表上游自动支持该接口，对应协议开关也必须开启。

**不要在插件 YAML 中配置 `api-keys`。** 此旧字段已停用，内容会被忽略，不能用于认证；请通过上述账户入口保存认证文件。

## 常见问题

| 现象 | 处理方式 |
|---|---|
| 商店没有插件或无法下载 | 检查自定义商店 URL 是否可访问、索引是否已上线，以及 Releases 是否有当前平台的 ZIP；当前仅支持 `linux/amd64` |
| 已复制 `.so` 但未加载 | 检查宿主是否支持动态插件、文件是否位于正确目录，以及总开关与插件自身开关是否均已开启 |
| 页面提示需要记住密码或管理登录失效 | 在同一地址的管理中心重新登录并勾选“记住密码”，再打开插件页面 |
| 添加 key 失败 | 确认 key 无空白字符、账户可访问身份与订阅接口、套餐属于有效 GOAT / Pro / Max；遇限流或网络异常稍后再试 |
| 模型列表为空 | 先确认存在已启用的有效认证文件，再检查目录请求、协议开关与宿主日志；未获可路由目录时不会凭空生成模型 |
| 提示重复 key 或保存状态不确定 | 先检查认证文件，包括已禁用账户，不要重复提交 |

如果是插件重命名前的安装，请将原配置迁移到 `plugins.configs.commandcode-cliproxyapi-provider`，使用新名称的 `.so` 和管理页面，并移除不再使用的旧插件文件，避免重复加载。旧页面的浏览器缓存不会自动迁移；认证与模型归属使用的 ProviderID 仍为 `commandcode`，无需改名。

## 安全提示

API key 和认证文件均属于敏感凭据。不要提交到仓库、截图分享或公开管理页面；浏览器授权后也应妥善处理剪贴板和官方页面的临时存储。备份认证文件时限制访问权限。仅对可信用户启用“记住密码”。
