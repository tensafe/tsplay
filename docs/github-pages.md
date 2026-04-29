# GitHub Pages 发布说明

这套文档站已经按 `MkDocs Material + GitHub Actions + GitHub Pages` 的方式接好了发布骨架。

## 已经加好的文件

- [mkdocs.yml](https://github.com/tensafe/tsplay/blob/main/mkdocs.yml)
- [requirements-docs.txt](https://github.com/tensafe/tsplay/blob/main/requirements-docs.txt)
- [site-src/index.md](https://github.com/tensafe/tsplay/blob/main/site-src/index.md)
- [site-src/getting-started.md](https://github.com/tensafe/tsplay/blob/main/site-src/getting-started.md)
- [site-src/assets/stylesheets/extra.css](https://github.com/tensafe/tsplay/blob/main/site-src/assets/stylesheets/extra.css)
- [tools/prepare_docs_site.py](https://github.com/tensafe/tsplay/blob/main/tools/prepare_docs_site.py)
- [theme-overrides/main.html](https://github.com/tensafe/tsplay/blob/main/theme-overrides/main.html)
- [docs-site.yml](https://github.com/tensafe/tsplay/blob/main/.github/workflows/docs-site.yml)

## 首次启用怎么做

1. 把这些改动推到 GitHub 仓库。
2. 打开仓库 `Settings -> Pages`。
3. 在 `Build and deployment` 里把 `Source` 设为 `GitHub Actions`。
4. 确认默认发布分支会触发 [docs-site.yml](https://github.com/tensafe/tsplay/blob/main/.github/workflows/docs-site.yml)。
5. 等待 workflow 跑完，站点就会发布到默认地址。

默认项目站点地址通常是：

```text
https://tensafe.github.io/tsplay/
```

如果你改了仓库名或组织名，记得同步更新 [mkdocs.yml](https://github.com/tensafe/tsplay/blob/main/mkdocs.yml) 里的 `site_url` 和 `repo_url`。

## 本地预览

```bash
python3 -m venv .venv-docs
source .venv-docs/bin/activate
pip install -r requirements-docs.txt
python3 tools/prepare_docs_site.py
mkdocs serve
```

默认预览地址通常是：

```text
http://127.0.0.1:8000
```

## 这版为什么这样组织

- 构建前会先运行 [tools/prepare_docs_site.py](https://github.com/tensafe/tsplay/blob/main/tools/prepare_docs_site.py)
- 它会把 `README.zh-CN.md`、`ReadMe.md`、`docs/`、`script/`、`demo/` 复制到临时站点目录 `site-docs/`
- 这样现有教程里大量指向 `../README.zh-CN.md`、`../../script/...`、`../../demo/...` 的链接都还能继续工作
- 不需要为了发站点，去大改几百个教程链接
- `site-docs/` 和最终生成的 `site/` 都属于构建产物，不需要提交

## 这版站点额外做了什么

- 用 `Material for MkDocs` 替换了默认风格，导航、搜索、代码复制和移动端体验会更友好
- 首页改成了面向角色的 landing page，而不是只放一段普通索引
- 增加了 [site-src/getting-started.md](https://github.com/tensafe/tsplay/blob/main/site-src/getting-started.md) 作为独立快速开始页
- 用 [site-src/assets/stylesheets/extra.css](https://github.com/tensafe/tsplay/blob/main/site-src/assets/stylesheets/extra.css) 定义了品牌色、hero 区、卡片和表格样式
- 用 [theme-overrides/main.html](https://github.com/tensafe/tsplay/blob/main/theme-overrides/main.html) 加了公告条入口

## 后续建议

- 如果你想把首页做得更像产品官网，可以继续补截图、路线图和“我属于哪类用户”的入口卡片
- 如果你想给中文用户更顺手的导航，可以逐步把更多中文页面放进主导航
- 如果你后面要接自定义域名，可以在仓库 Pages 设置里继续加 `Custom domain`
