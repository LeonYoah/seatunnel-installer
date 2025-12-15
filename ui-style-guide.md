# UI 全局样式简要规范

> 供开发人员在实现页面时统一风格使用（基于 Element Plus + 自定义主题）。

---

## 1. 颜色（Color）

- **主色（Primary）**
  - CSS 变量：`--primary`, `--primary-2`
  - 用途：
    - 主要按钮、激活状态、关键链接
- **中性色（Neutral）**
  - 背景：`--bg`, `--surface`, `--surface-2`
  - 边框：`--border`
  - 文本：`--text`, `--muted`

> 规则：想“突出可点击 / 当前状态”用主色，作为背景或文字；想“做背景 / 边线”用中性色。

---

## 2. 圆角（Radius）

- 全局按钮、输入框等小组件：`4px ~ 6px`
- 卡片、弹框：`6px ~ 8px`
- 当前按钮圆角变量：`--btn-radius: 6px`

> 要求：同一页面中，按钮和输入框的圆角保持一致，不混用尖角/大圆角。

---

## 3. 按钮（Button）

### 3.1 尺寸（Size）

- 使用 CSS 变量：
  - 默认（Comfortable）：
    - 高度：`--btn-height: 36px`
    - 水平内边距：`--btn-padding-x: 18px`
  - 紧凑（Compact）：
    - 高度：`--btn-height-compact: 32px`
    - 水平内边距：`--btn-padding-x-compact: 14px`
  - 字体：`--btn-font-size: 14px`, `--btn-font-weight: 500`

- 映射到 Element Plus：
  - 默认按钮（无 `size`） → **Comfortable**
  - `size="small"` → **Compact**

### 3.2 形态（Shape）

- 全局覆盖：`.el-button:not(.is-circle):not(.is-text)`
  - 使用上面的高度、圆角、内边距、字体变量
  - 相邻按钮自动保持 `--btn-gap: 12px` 的水平间距

### 3.3 类型（Type）

- **主操作按钮（Hero Primary）**
  - 使用类：`.btn-primary-action` + Element `type="primary"`
  - 效果：
    - 填充主色背景 `var(--primary)`，文字白色
    - Hover / Active 使用 `var(--primary-2)`
  - 场景：
    - 空状态 CTA（如「立即巡检」）
    - 创建/启动/部署等主操作

- **次要按钮（Utility / Secondary）**
  - 使用 Element 默认 `plain` 或在后续通过 `.btn-secondary` 扩展
  - 场景：
    - 同一操作组中的辅助操作（如「定期巡检」）
    - 工具栏中的一般操作（刷新、筛选）

> 规范：一个操作区域内尽量只保留 **一个** Hero Primary，其余使用次要样式。

---

## 4. 间距（Spacing）

- 统一使用 8/12/16/24px 作为主要间距节奏：
  - 8px：组件内部小间距（图标与文字）
  - 12px：按钮组、表格工具栏内的间距
  - 16px：卡片内容内的段落间距
  - 24px：页面主体与卡片之间间距

---

## 5. 字体（Typography）

- 基础：浏览器默认 sans-serif（已在全局设置）
- 建议层级：
  - 页面主标题：18–20px，`font-weight: 600`
  - 区块标题 / 卡片标题：16px，`font-weight: 600`
  - 正文：14px
  - 次要说明：12–13px，颜色使用 `--muted`

> 页面开发时，尽量复用这四个层级，不再额外创造很多字号。

---

## 6. 使用示例

- 主操作按钮（Hero Primary）：

  ```vue
  <el-button type="primary">
    立即巡检
  </el-button>
  ```

- 同一组中的次要按钮（Utility）：

  ```vue
  <el-button plain>
    定期巡检
  </el-button>
  ```

开发新页面时，优先遵循以上变量和约定，再根据具体业务进行微调。
