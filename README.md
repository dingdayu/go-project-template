# Go Template (Gin)

一键起步的 Go Web 模板，内置：Gin 路由、CI、Lint、Docker、多环境配置，以及**模块名自动重命名**。

## 初始化

> 通过 GitHub 的 **Use this template** 创建新仓库后：

**方式 A：本地脚本**

```bash
./scripts/rename-module.sh github.com/<your>/<repo>
```

**方式 B：GitHub Actions**

- 打开 GitHub → Actions → **Rename Go module (one-time)** → Run workflow
- 可留空使用默认 `github.com/<owner>/<repo>`

## 开发

```bash
make tidy
make run
```

访问：`http://localhost:8080/healthz`、`/api/v1/hello`

## 测试 & Lint

```bash
make test
make lint
```

## 构建镜像

```bash
make docker-build
```

## 配置

- `PORT`：监听端口（默认 `8080`）
- `APP_ENV`：`dev|prod`（影响 Gin 日志级别）
- `APP_NAME`：应用名

## 版本信息

`model/entity/version.go` 会在构建时被 `-ldflags` 注入，CI/Makefile 已示例。
