# Clash Auto Updater

一个用 Go 编写的自动化工具，用于获取多个Clash订阅，根据关键词进行过滤，并生成一个合并后的 `config.yaml` 文件。

## ✨ 功能

-   **多订阅合并**: 从多个URL获取订阅源，并将所有代理节点合并。
-   **智能过滤**: 根据您设置的关键词（如地区、类型等）筛选出您需要的节点。
-   **模板化配置**: 使用一个基础的 `template.yaml`，程序会自动将筛选后的节点填充进去。
-   **简单易用**: 只需一个配置文件和一条命令即可完成所有操作。

## 🚀 如何开始

### 1. 准备

-   确保您已经安装了 Go (版本 1.18+)。

### 2. 配置

1.  将 `config/` 目录下的 `config.yaml.example` 文件复制一份，并重命名为 `config.yaml`。
2.  打开 `config.yaml` 并填入您的信息：
    -   `subscriptions`: 您的Clash订阅链接列表。**重要提示**: 订阅源必须提供Clash兼容的YAML格式的代理列表。
    -   `filter_rules.include_keywords`: 您想要保留的节点的关键词。节点的名称（不区分大小写）中只要包含任意一个关键词，就会被保留。
    -   `template_path`: 您的Clash基础模板文件的路径，默认为 `./config/template.yaml`。
    -   `output_path`:最终生成的配置文件的存放路径，默认为 `./dist/config.yaml`。
3.  (可选) 您可以根据自己的需求修改 `config/template.yaml` 文件，例如更改端口、模式或添加自己的规则。

### 3. 运行

在项目根目录下执行以下命令：

```bash
go run cmd/clash-auto/main.go
```

或者，您也可以构建一个二进制文件来运行：

```bash
# 构建
go build -o clash-auto cmd/clash-auto/main.go

# 运行
./clash-auto
```

程序执行成功后，您将在 `dist/` 目录下找到最终生成的 `config.yaml` 文件。

## 📂 项目结构

```
clash-auto/
├── cmd/clash-auto/     # 程序主入口
├── internal/           # 内部模块 (配置、下载、解析、过滤、生成)
├── config/             # 用户配置文件和模板
│   ├── config.yaml.example
│   └── template.yaml
├── dist/               # 生成的配置文件输出目录
├── go.mod
└── README.md
```
