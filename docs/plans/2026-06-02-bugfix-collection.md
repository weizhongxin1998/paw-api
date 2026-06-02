# 集合右键菜单修复计划

## 问题 1：集合右键菜单失效

**现象**：在集合上右键没有弹出上下文菜单，无法创建子集合或请求。

**原因**：
- `NTree` 的 `@contextmenu` 事件绑定在 `AppSidebar.vue`，但 `NTree` 的 `contextmenu` 事件回调参数是 `(e, option)`，而当前 `handleContextMenu` 需要 `option.key`。
- 可能原因是 `NTree` 的 `contextmenu` 事件发出的对象结构与 `handleContextMenu` 预期不匹配，或有冒泡/阻止默认行为的问题。

**复现**：
1. 打开工作区，确保存在至少一个集合
2. 在集合名称上右键
3. 应该出现右键菜单（重命名集合、新建子集合、删除集合），但实际没有

**影响**：集合的右键菜单完全不可用，用户无法通过右键创建子集合或重命名/删除集合。

---

## 问题 2：三点菜单与右键菜单冗余

**现象**：集合上同时存在 `···` 三点展开菜单和右键菜单，两者功能重叠。

**当前行为**：
- 每个集合节点右侧悬浮显示 `···` 按钮，点击弹出小菜单（仅"New Request"）
- 右键集合应弹出完整菜单（重命名、新建子集合、删除集合等）
- 两套菜单并存，交互混乱

**要求**：
1. 移除集合节点上的 `···` 三点按钮（`renderLabel` 中的 `tree-node-actions` 部分）
2. 统一使用右键菜单
3. 右键菜单补充缺少的"新建接口"选项（当前只有"New Sub-collection"、"Rename Collection"、"Delete Collection"，缺少"New Request"）

---

## 涉及文件

- `frontend/src/components/AppSidebar.vue`
  - `renderLabel` 函数（移除三点按钮）
  - `handleContextMenu` 函数（排查右键不触发的原因）
  - 右键菜单模板（补充"New Request"选项）
  - `colMenuOptions` / `handleColMenuSelect`（可移除，三点菜单的产物）

## Work

小
