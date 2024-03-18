# Fireboom 后端技术沉淀

## Gitea 自动化部署

`Gitea` 是一个开源的 Git 托管服务，`Gitea`提供的 CI/CD 配置借鉴了 Github Actions，可以减少很多学习成本，并且有大量的第三方 Action 可供选择。

基于 `Fireboom` 的打包分为2部分，一个是纯代码的 Docker 镜像，一个是钩子运行时镜像

- [点击查看代码镜像配置](./.gitea/workflows/build-fb-data.yaml)
- [点击查看钩子运行时镜像配置](./.gitea/workflows/build-fb-hook.yaml)

这些脚本主要包括以下步骤：

1. 克隆项目代码
2. 登录 Harbor 镜像仓库
3. 构建 Docker 镜像
4. 推送 Docker 镜像到 Harbor 镜像仓库
5. 使用 `kubectl` 命令更新 Kubernetes 集群中的 `Deployment` 资源，生产环境请勿直接这么操作