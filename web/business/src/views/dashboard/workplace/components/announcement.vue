<template>
  <a-spin :loading="loading" style="width: 100%">
    <a-card
      class="general-card"
      :title="$t('workplace.announcement')"
      :header-style="{ paddingBottom: '0' }"
      :body-style="{ padding: '15px 20px 13px 20px' }"
    >
      <div v-if="list.length">
        <div v-for="item in list" :key="item.id" class="item">
          <a-tag :color="tagColor(item.type)" size="small">{{ tagText(item.type) }}</a-tag>
          <span class="item-content" :class="{ unread: item.isread === 0 }">
            {{ item.title || item.content }}
          </span>
        </div>
      </div>
      <a-empty v-else description="暂无消息" />
    </a-card>
  </a-spin>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import useLoading from '@/hooks/loading';
  import { queryMessageList } from '@/api/message';

  type MessageItem = {
    id: number;
    type: string | number;
    title?: string;
    content?: string;
    isread?: number;
    createtime?: number | string;
    path?: string;
  };

  const { loading, setLoading } = useLoading(true);
  const list = ref<MessageItem[]>([]);

  const tagText = (type: MessageItem['type']) => {
    if (type === 'notice' || type === 2) return '通知';
    if (type === 'todo' || type === 3) return '待办';
    return '消息';
  };
  const tagColor = (type: MessageItem['type']) => {
    if (type === 'notice' || type === 2) return 'blue';
    if (type === 'todo' || type === 3) return 'orangered';
    return 'cyan';
  };

  const fetchData = async () => {
    setLoading(true);
    try {
      const data = await queryMessageList();
      list.value = (data?.items || []).slice(0, 5) as MessageItem[];
    } catch (err) {
      // 你可以在这里统一处理错误提示
    } finally {
      setLoading(false);
    }
  };

  fetchData();
</script>

<style scoped lang="less">
  .item {
    display: flex;
    align-items: center;
    width: 100%;
    height: 24px;
    margin-bottom: 4px;
    .item-content {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      margin-left: 4px;
      color: var(--color-text-2);
      text-decoration: none;
      font-size: 13px;
      cursor: pointer;
    }
    .unread {
      color: var(--color-text-1);
      font-weight: 500;
    }
  }
</style>
