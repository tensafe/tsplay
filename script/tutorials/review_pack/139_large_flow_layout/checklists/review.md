# Lesson 139 Package Review Checklist

review 一个大型 Flow 包时，优先确认：

1. `collect / verify / publish` 是否真的拆开了责任。
2. 每个 stage 能不能单独复跑。
3. 每个 stage 的 artifact 是否在自己的子目录里。
4. 输出文件名能不能直接看懂用途。
5. 新同学能不能只看目录就大概猜到执行顺序。
