# PC web 平台

## Gitea 自动化部署

`Gitea` 是一个开源的 Git 托管服务，`Gitea`提供的 CI/CD 配置借鉴了 Github Actions，可以减少很多学习成本，并且有大量的第三方 Action 可供选择。

[点击查看前端 Gitea Action打包部署脚本](./.gitea/workflows/build-docker.yaml)。该脚本主要包括以下步骤：

1. 克隆项目代码
2. 登录 Harbor 镜像仓库
3. 构建 Docker 镜像
4. 推送 Docker 镜像到 Harbor 镜像仓库
5. 使用 `kubectl` 命令更新 Kubernetes 集群中的 `Deployment` 资源，生产环境请勿直接这么操作