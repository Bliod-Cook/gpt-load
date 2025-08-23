/**
 * 格式化相对时间
 * @param date 日期字符串
 * @returns 格式化后的相对时间字符串
 */
export function formatRelativeTime(date: string): string {
  if (!date) {
    return "未记录";
  }

  const now = new Date();
  const target = new Date(date);
  const diffSeconds = Math.floor((now.getTime() - target.getTime()) / 1000);
  const diffMinutes = Math.floor(diffSeconds / 60);
  const diffHours = Math.floor(diffMinutes / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffDays > 0) {
    return `${diffDays}天前`;
  }
  if (diffHours > 0) {
    return `${diffHours}小时前`;
  }
  if (diffMinutes > 0) {
    return `${diffMinutes}分钟前`;
  }
  if (diffSeconds > 0) {
    return `${diffSeconds}秒前`;
  }
  return "刚刚";
}

/**
 * 格式化完整日期时间
 * @param date 日期字符串
 * @returns 格式化后的日期时间字符串
 */
export function formatDate(date: string): string {
  if (!date) {
    return "未记录";
  }

  return new Date(date).toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

/**
 * 格式化相对时间，针对"最后使用"字段的特殊处理
 * @param date 日期字符串
 * @returns 格式化后的相对时间字符串，空值时返回"从未"
 */
export function formatLastUsedTime(date: string): string {
  if (!date) {
    return "从未";
  }

  return formatRelativeTime(date);
}
