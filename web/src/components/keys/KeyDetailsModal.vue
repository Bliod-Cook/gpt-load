<script setup lang="ts">
import { computed, ref } from "vue";
import type { APIKey } from "@/types/models";
import { copy } from "@/utils/clipboard";
import { maskKey } from "@/utils/display";
import { formatDate, formatRelativeTime } from "@/utils/time";
import { Close, CopyOutline, EyeOutline, EyeOffOutline } from "@vicons/ionicons5";
import {
  NButton,
  NCard,
  NCode,
  NDescriptions,
  NDescriptionsItem,
  NDivider,
  NIcon,
  NModal,
  NTooltip,
  NSpace,
  NTag,
  NText
} from "naive-ui";

interface Props {
  show: boolean;
  keyData: APIKey | null;
}

interface Emits {
  (e: "update:show", value: boolean): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

// 关闭弹窗
function handleClose() {
  emit("update:show", false);
}

// 是否显示完整密钥
const showKeyValue = ref(false);

// 复制通用方法
async function copyField(text: string, okMsg?: string) {
  const success = await copy(text);
  if (success) {
    window.$message.success(okMsg || "已复制到剪贴板");
  } else {
    window.$message.error("复制失败");
  }
}

async function copyKeyId() {
  if (!props.keyData) return;
  await copyField(String(props.keyData.id), "密钥 ID 已复制");
}

async function copyKeyValue() {
  if (!props.keyData) return;
  await copyField(props.keyData.key_value, "密钥值已复制");
}

function toggleShowKey() {
  showKeyValue.value = !showKeyValue.value;
}

// 成功率（%），无请求时返回 null
const successRate = computed(() => {
  const total = props.keyData?.request_count ?? 0;
  if (!total) return null;
  const failures = props.keyData?.failure_count ?? 0;
  const success = Math.max(total - failures, 0);
  return Math.round((success / total) * 100);
});

async function copyErrorMessage() {
  if (!props.keyData?.last_failure_error) {
    return;
  }
  const success = await copy(props.keyData.last_failure_error);
  if (success) {
    window.$message.success("错误信息已复制到剪贴板");
  } else {
    window.$message.error("复制失败");
  }
}
</script>

<template>
  <n-modal :show="show" @update:show="handleClose" class="key-details-modal">
    <n-card
      style="width: 680px; max-width: 92vw"
      title="密钥详情"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
    >
      <template #header-extra>
        <n-button quaternary circle @click="handleClose">
          <template #icon>
            <n-icon :component="Close" />
          </template>
        </n-button>
      </template>

      <div v-if="keyData" class="key-details">
        <!-- 基本信息 -->
        <n-descriptions title="基本信息" :column="2" label-placement="left" bordered>
          <n-descriptions-item label="密钥ID">
            <n-space align="center" size="small">
              <n-text code>{{ keyData.id }}</n-text>
              <n-tooltip trigger="hover" :show-arrow="false">
                <template #trigger>
                  <n-button text size="tiny" @click="copyKeyId" title="复制密钥ID">
                    <template #icon>
                      <n-icon :component="CopyOutline" />
                    </template>
                  </n-button>
                </template>
                复制
              </n-tooltip>
            </n-space>
          </n-descriptions-item>
          <n-descriptions-item label="密钥值" :span="2">
            <n-space align="center" justify="space-between" style="width: 100%">
              <n-text code>
                {{ showKeyValue ? keyData.key_value : maskKey(keyData.key_value) }}
              </n-text>
              <div class="inline-actions">
                <n-tooltip trigger="hover" :show-arrow="false">
                  <template #trigger>
                    <n-button text size="tiny" @click="toggleShowKey" :title="showKeyValue ? '隐藏' : '显示'">
                      <template #icon>
                        <n-icon :component="showKeyValue ? EyeOffOutline : EyeOutline" />
                      </template>
                      {{ showKeyValue ? '隐藏' : '显示' }}
                    </n-button>
                  </template>
                  {{ showKeyValue ? '隐藏密钥' : '显示密钥' }}
                </n-tooltip>
                <n-tooltip trigger="hover" :show-arrow="false">
                  <template #trigger>
                    <n-button text size="tiny" @click="copyKeyValue" title="复制密钥值">
                      <template #icon>
                        <n-icon :component="CopyOutline" />
                      </template>
                      复制
                    </n-button>
                  </template>
                  复制
                </n-tooltip>
              </div>
            </n-space>
          </n-descriptions-item>
          <n-descriptions-item label="状态">
            <n-tag
              v-if="keyData.status === 'active'"
              type="success"
              :bordered="false"
              round
            >
              有效
            </n-tag>
            <n-tag v-else type="error" :bordered="false" round>
              无效
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="请求次数">{{ keyData.request_count }}</n-descriptions-item>
          <n-descriptions-item label="失败次数">{{ keyData.failure_count }}</n-descriptions-item>
          <n-descriptions-item label="成功率">
            {{ successRate !== null ? successRate + '%' : '—' }}
          </n-descriptions-item>
          <n-descriptions-item label="最后使用">
            {{ formatDate(keyData.last_used_at || "") }}
            <n-text depth="3" style="margin-left: 8px">
              ({{ formatRelativeTime(keyData.last_used_at || "") }})
            </n-text>
          </n-descriptions-item>
        </n-descriptions>

        <n-divider />

        <!-- 失败信息 -->
        <n-descriptions title="最后一次失败信息" :column="1" label-placement="left" bordered>
          <n-descriptions-item v-if="keyData.last_failure_at" label="失败时间">
            {{ formatDate(keyData.last_failure_at || "") }}
            <n-text depth="3" style="margin-left: 8px">
              ({{ formatRelativeTime(keyData.last_failure_at || "") }})
            </n-text>
          </n-descriptions-item>
          <n-descriptions-item v-if="keyData.last_failure_error"  label="错误信息" :span="2">
            <div class="error-section">
              <n-space align="center" justify="space-between" style="margin-bottom: 8px">
                <span>错误详情</span>
                <n-button text size="tiny" @click="copyErrorMessage" title="复制错误信息">
                  <template #icon>
                    <n-icon :component="CopyOutline" />
                  </template>
                  复制
                </n-button>
              </n-space>
              <n-code
                language="text"
                :code="keyData.last_failure_error"
                word-wrap
                style="max-height: 40vh; max-width: 100%; overflow: auto;"
              />
            </div>
          </n-descriptions-item>
          <n-descriptions-item v-if="!keyData.last_failure_error && !keyData.last_failure_at" label="状态" :span="2">
            <n-text depth="3">该密钥尚未记录失败信息</n-text>
          </n-descriptions-item>
        </n-descriptions>

        <n-divider />

        <!-- 时间信息 -->
        <n-descriptions title="时间信息" :column="2" label-placement="left" bordered>
          <n-descriptions-item label="创建时间">
            {{ formatDate(keyData.created_at) }}
            <n-text depth="3" style="margin-left: 8px">
              ({{ formatRelativeTime(keyData.created_at || "") }})
            </n-text>
          </n-descriptions-item>
          <n-descriptions-item label="更新时间">
            {{ formatDate(keyData.updated_at) }}
            <n-text depth="3" style="margin-left: 8px">
              ({{ formatRelativeTime(keyData.updated_at || "") }})
            </n-text>
          </n-descriptions-item>
        </n-descriptions>
      </div>

      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 12px">
          <n-button @click="handleClose">关闭</n-button>
        </div>
      </template>
    </n-card>
  </n-modal>
</template>

<style scoped>
.key-details-modal {
  --n-color: rgba(255, 255, 255, 0.95);
}

.key-details {
  max-height: 70vh;
  overflow-y: auto;
}

.error-section {
  width: 100%;
}

.inline-actions {
  display: flex;
  gap: 4px;
}

/* 确保错误信息在任意尺寸下都不会溢出屏幕 */
.error-section :deep(.n-code) {
  max-width: 100%;
  overflow: auto;
}
.error-section :deep(pre) {
  white-space: pre-wrap;         /* 允许换行 */
  overflow-wrap: anywhere;       /* 极长单词/连续字符也可换行 */
}

:deep(.n-card-header) {
  border-bottom: 1px solid rgba(239, 239, 245, 0.8);
  padding: 10px 20px;
}

:deep(.n-card__content) {
  max-height: calc(100vh - 68px - 61px - 50px);
  overflow-y: auto;
}

:deep(.n-card__footer) {
  border-top: 1px solid rgba(239, 239, 245, 0.8);
  padding: 10px 15px;
}

:deep(.n-descriptions) {
  margin-bottom: 0;
}

:deep(.n-descriptions .n-descriptions-table-wrapper) {
  border-radius: 6px;
}

@media (max-width: 768px) {
  :deep(.n-descriptions) {
    --n-th-padding: 8px;
    --n-td-padding: 8px;
  }
}
</style>
