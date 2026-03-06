# GitHub Secrets 手动配置指南

## 配置步骤

1. 打开你的GitHub仓库页面
2. 点击 **Settings** (设置)
3. 左侧菜单找到 **Secrets and variables** → **Actions**
4. 点击 **New repository secret** 按钮
5. 按照下面的表格逐个添加

---

## 需要添加的 Secrets

### 私有Docker Registry配置

| Name | Value |
|------|-------|
| `PRIVATE_REGISTRY_URL` | `docker.topren.top` |
| `PRIVATE_REGISTRY_USERNAME` | `admin` |
| `PRIVATE_REGISTRY_PASSWORD` | `Password123` |

### 生产环境服务器配置

| Name | Value |
|------|-------|
| `DEPLOY_HOST` | `154.12.55.183` |
| `DEPLOY_USER` | `root` |
| `DEPLOY_PASSWORD` | `ZXGNaunfomFD25PD` |
| `DEPLOY_PORT` | `22` |
| `DEPLOY_DIR` | `/root/CLIProxyAPI` |

---

## 配置完成后

### 验证配置
1. 去 **Settings → Secrets and variables → Actions**
2. 确认所有8个Secrets都已添加
3. Secret名称必须完全匹配（区分大小写）

### 触发首次部署
```bash
# 方式1：推送代码到main分支
git add .
git commit -m "feat: setup auto deploy workflow"
git push origin main

# 方式2：手动触发workflow
# 去 GitHub Actions 页面 → Dev Private Deploy → Run workflow
```

### 查看部署日志
1. 去仓库的 **Actions** 页面
2. 点击最新的 workflow run
3. 查看 build-and-push 和 deploy-prod 的日志

---

## 常见问题

### Q: Secret添加后看不到值？
A: 正常的，GitHub出于安全考虑会隐藏Secret的值，只能看到名称和最后更新时间。

### Q: Workflow运行失败？
A: 检查以下几点：
1. Secret名称是否完全匹配（区分大小写）
2. 服务器SSH连接是否正常
3. 私有仓库登录凭证是否正确
4. 查看Actions日志的详细错误信息

### Q: 如何修改Secret？
A: 点击Secret名称 → Update secret → 输入新值 → Update secret

---

## 安全建议

⚠️ **重要提示**：
- 不要在代码中硬编码密码
- 不要在commit message中包含敏感信息
- 定期更换服务器密码和Docker Registry密码
- 生产环境建议使用SSH密钥而不是密码认证
