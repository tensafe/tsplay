# TSPlay 文档健康检查与断层清单

> 目标：把“链接有没有坏”和“人会不会在这里断掉”分开检查。

## 当前状态

### 链接健康

当前仓库已经有可重复执行的相对链接检查：

- 一条统一命令：`tools/check_docs_suite.py`
- 链接检查脚本：`tools/check_markdown_links.py`
- 连续性检查脚本：`tools/check_doc_continuity.py`
- CI 接入：`.github/workflows/docs-site.yml`

最近一次本地检查结论：

```text
## Repository Docs
OK

## Generated Site Docs
OK
```

也就是说：

- 仓库源码视图里的相对 Markdown 链接是通的
- 生成到 `site-docs/` 和 `site/` 后的站点链接也是通的

### 已经补掉的硬断层

- 教程入口页原先跳到不存在的根目录 `getting-started.md`
- 根目录快速开始入口只在站点层存在，仓库源码视图里不够顺手
- 文档站 CI 以前不检查 Markdown 相对链接

这些问题现在都已经补上。

## 仍要持续注意的内容断层

下面这些不属于“坏链”，但会影响阅读连续性。

### 1. 首跑成功后的下一步分叉容易过宽

风险：

- 用户跑通 `Lesson 01` 后，会同时看到教程、训练、MCP、单二进制、培训体系
- 如果没有一条默认下一步，第一次上手的人容易重新掉回“大地图焦虑”

当前处理：

- `site-src/getting-started.md` 已补“默认下一步”和“常见断层提醒”
- [../getting-started.md](../getting-started.md) 也同步补了源码视图可读版

### 2. `site-src/` 看起来像坏链，但发布路径是对的

原因：

- `site-src/` 不是最终站点根目录
- 它会先被 `tools/prepare_docs_site.py` 复制到 `site-docs/`

结论：

- 不应直接按 `site-src/` 原地检查所有相对链接
- 应该以生成后的 `site-docs/` 为准

### 3. 英文教程地图目前仍承担“导航层”多于“完整英文课程层”

现状：

- [tutorials/README.md](tutorials/README.md) 已能承担英文入口和路线说明
- 但大量单节 lesson 仍主要是中文内容

这不是坏链问题，但属于体验断层：

- 英文用户能找到路
- 不一定能在每一节里得到完整英文讲解

### 4. 单二进制入口仍需长期和 release 包保持同步

现状：

- `list-assets / extract-assets / file-srv / getting-started`
- 这四个入口已经有教程和路线图承接

风险：

- 以后如果 release 产物结构变化，文档入口和 release 资产说明容易再次分叉

建议：

- 每次改 release workflow 时，顺手检查 [product/core-feature-roadmap.md](product/core-feature-roadmap.md) 和 `getting-started` 相关说明

### 5. 入口页文案风格需要持续保持同一语气

风险：

- 入口页如果有的偏“命令式”，有的偏“说明式”，用户会觉得页面之间像是不同人写的
- 即使链接都通，语气突兀也会让阅读节奏断掉

建议：

- 先说明这页帮用户解决什么，再给默认下一步
- 多用“先、再、接着、如果你更关心”这类承接语气
- 少用过硬的命令感表达，尤其是在首页、快速开始、文档总入口这几类页面

## 建议的固定检查动作

每次文档或站点入口有明显变更时，至少跑一次：

```bash
python3 tools/check_docs_suite.py
```

如果你只是缩小范围排查，也可以单独跑：

```bash
python3 tools/check_markdown_links.py
python3 tools/check_doc_continuity.py
```

每次下列文件变更时，也建议顺手检查“内容断层”：

- `README.zh-CN.md`
- `getting-started.md`
- `site-src/getting-started.md`
- `docs/tutorials/README.zh-CN.md`
- `docs/README.md`
- `mkdocs.yml`

## 文档健康判断标准

可以用下面这 4 条快速判断：

- 用户第一次进仓库，能在 1 次跳转内找到最小起步命令
- 用户跑通后，能在 1 次跳转内找到下一条默认学习路径
- 用户切换到 `Agent / MCP / 单二进制 / 培训` 任一方向时，不需要先读完整地图
- 文档站和 GitHub 源码视图里的入口指向保持一致
