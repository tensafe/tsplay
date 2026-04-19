# Lesson 41: 查看和导出命名会话信息

这一节不新建状态文件，  
而是学习怎么检查上一节注册好的命名会话。

目标：

- 列出当前所有已保存会话
- 查看一个会话的详情
- 导出可直接复用的 snippet

## 开始前

建议先完成：

- [Lesson 40](40-save-named-session.md)

## Step 1: 列出当前会话

```bash
./tsplay -action list-sessions
```

预期结果：

- 能看到 `session_lab_demo`

## Step 2: 查看单个会话详情

```bash
./tsplay -action get-session -session-name session_lab_demo
```

这一步适合确认：

- 会话类型
- 底层 storage state 路径
- 推荐的 `browser.use_session` 写法

## Step 3: 导出复用 snippet

导出推荐的 Flow YAML 片段：

```bash
./tsplay -action export-session \
  -session-name session_lab_demo \
  -session-format flow_yaml
```

如果你想看浏览器块本身，也可以：

```bash
./tsplay -action export-session \
  -session-name session_lab_demo \
  -session-format browser_yaml
```

## Step 4: 这一节在课程中的作用

这一节是在帮你建立“会话不是黑盒”的意识。  
你不仅能用它，还能：

- 看见它
- 检查它
- 导出它

## 下一节

下一节正式进入最常用的复用方式：`use_session`。
[Lesson 42](42-use-named-session.md)
