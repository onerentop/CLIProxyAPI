#!/bin/bash
# GitHub Secrets 自动配置脚本
# 使用方法：
#   1. 安装 GitHub CLI: https://cli.github.com/
#   2. 登录: gh auth login
#   3. 运行: bash setup-github-secrets.sh

set -e

echo "=========================================="
echo "🔧 配置 GitHub Secrets for CLIProxyAPI"
echo "=========================================="

# 检查 gh CLI 是否安装
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) 未安装"
    echo "请访问 https://cli.github.com/ 安装"
    exit 1
fi

# 检查是否已登录
if ! gh auth status &> /dev/null; then
    echo "❌ 未登录 GitHub CLI"
    echo "请运行: gh auth login"
    exit 1
fi

echo ""
echo "📝 开始配置 Secrets..."
echo ""

# 私有Docker Registry配置
echo "🐳 配置私有Docker Registry..."
gh secret set PRIVATE_REGISTRY_URL --body "docker.topren.top"
gh secret set PRIVATE_REGISTRY_USERNAME --body "admin"
gh secret set PRIVATE_REGISTRY_PASSWORD --body "Password123"

# 生产环境服务器配置
echo "🚀 配置生产环境服务器..."
gh secret set DEPLOY_HOST --body "154.12.55.183"
gh secret set DEPLOY_USER --body "root"
gh secret set DEPLOY_PASSWORD --body "ZXGNaunfomFD25PD"
gh secret set DEPLOY_PORT --body "22"
gh secret set DEPLOY_DIR --body "/root/CLIProxyAPI"

echo ""
echo "=========================================="
echo "✅ 所有 Secrets 配置完成！"
echo "=========================================="
echo ""
echo "📋 已配置的 Secrets:"
echo "  - PRIVATE_REGISTRY_URL"
echo "  - PRIVATE_REGISTRY_USERNAME"
echo "  - PRIVATE_REGISTRY_PASSWORD"
echo "  - DEPLOY_HOST"
echo "  - DEPLOY_USER"
echo "  - DEPLOY_PASSWORD"
echo "  - DEPLOY_PORT"
echo "  - DEPLOY_DIR"
echo ""
echo "🎯 下一步："
echo "  1. 推送代码到 main 分支触发自动部署"
echo "  2. 或者去 GitHub Actions 页面手动触发 workflow"
echo ""
