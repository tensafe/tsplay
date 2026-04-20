# Lesson 139 Large Flow Layout

这个目录不是为了演示“某个单独 action 怎么用”，而是为了演示大型 Flow 应该怎么拆。

推荐把一个较大的教程包拆成：

- `flows/collect.flow.yaml`
- `flows/verify.flow.yaml`
- `flows/publish.flow.yaml`
- `checklists/review.md`
- `data/`

这样的目的不是形式主义，而是为了让：

- 采集
- 校验
- 交付

这三件事能被分别 review、分别复跑、分别定位问题。
