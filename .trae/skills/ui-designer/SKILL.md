---
name: "ui-designer"
description: "Designs modern, minimalist UI for mobile apps. Invoke when user requests UI redesign for pages like login, ensuring clean aesthetics while preserving all business logic."
---

# UI Designer - 移动端简约大气风格

## 设计原则

### 1. 简约大气风格特征
- **留白艺术**: 大量留白，呼吸感强
- **色彩克制**: 主色调不超过3种，使用渐变色增加层次感
- **无边框设计**: 去除包裹边框，使用阴影和背景色区分层次
- **大字体标题**: 突出核心信息
- **圆角元素**: 按钮、输入框使用统一圆角
- **图标精简**: 使用简洁的图标，避免复杂装饰

### 2. 配色方案（房产/商务类App）
```scss
// 主色调 - 专业蓝绿
$primary: #0d9488;      // 主色 -  teal-600
$primary-dark: #0f766e; // 深色 - teal-700
$primary-light: #14b8a6;// 浅色 - teal-500

// 辅助色
$accent: #3b82f6;       // 强调色 - blue-500
$success: #10b981;      // 成功色

// 中性色
$text-primary: #1f2937;   // 主要文字 - gray-800
$text-secondary: #6b7280; // 次要文字 - gray-500
$text-tertiary: #9ca3af;  // 辅助文字 - gray-400
$bg-page: #f9fafb;        // 页面背景 - gray-50
$bg-card: #ffffff;        // 卡片背景
$border: #e5e7eb;         // 边框色 - gray-200
```

### 3. 布局规范
- 页面边距: 40rpx-48rpx
- 元素间距: 24rpx-32rpx
- 大标题: 48rpx-56rpx, font-weight: 700
- 小标题: 28rpx-32rpx, font-weight: 600
- 正文: 26rpx-30rpx
- 按钮高度: 96rpx-104rpx
- 圆角: 16rpx-24rpx

### 4. 阴影规范
```scss
// 轻微阴影
box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);

// 中等阴影
box-shadow: 0 8rpx 24rpx rgba(13, 148, 136, 0.12);

// 强调阴影
box-shadow: 0 12rpx 32rpx rgba(13, 148, 136, 0.18);
```

### 5. 登录页专用设计模式

#### 顶部区域
- Logo居中，大尺寸（120rpx-160rpx）
- 品牌名称使用粗体大字号
- Slogan使用次要色小字号

#### 中部表单区域
- 去除包裹边框，直接放置在页面上
- 输入框使用底部边框或轻微背景色
- 按钮使用渐变背景，增加视觉重量
- 社交登录按钮使用品牌色

#### 底部区域
- 协议链接居中排列
- 使用分割线或间距区分
- 字体小一号，颜色更淡

### 6. 动画效果
```scss
// 按钮点击效果
&:active {
  transform: scale(0.98);
  opacity: 0.9;
}

// 页面进入动画
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
```

## 重构步骤

1. **分析现有代码**: 识别所有数据绑定、事件处理、条件渲染
2. **保留逻辑**: 确保所有methods、data、生命周期函数不变
3. **重构模板**: 重新组织HTML结构，应用新的CSS类
4. **重写样式**: 使用新的设计规范替换原有样式
5. **验证功能**: 确保所有交互功能正常工作
