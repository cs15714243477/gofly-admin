<template>
  <div class="page-container">
    <div class="content-wrapper">
      <!-- 头部状态区：采用横向卡片布局，突出关键指标 -->
      <div class="dashboard-header">
        <a-row :gutter="24">
          <a-col :span="8">
            <div class="stat-card" :class="{ 'is-active': status.config_ok }">
              <div class="stat-icon-wrapper config-bg">
                <icon-settings class="stat-icon" />
              </div>
              <div class="stat-info">
                <div class="stat-label">配置状态</div>
                <div class="stat-value-row">
                  <span class="stat-value">{{ status.config_ok ? '已连接' : '未配置' }}</span>
                  <div class="status-dot" :class="status.config_ok ? 'bg-success' : 'bg-danger'"></div>
                </div>
                <div class="stat-sub">API: {{ status.api_base || '未设置' }}</div>
              </div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="stat-card" :class="{ 'is-active': status.has_token }">
              <div class="stat-icon-wrapper token-bg">
                <icon-safe class="stat-icon" />
              </div>
              <div class="stat-info">
                <div class="stat-label">Token 授权</div>
                <div class="stat-value-row">
                  <span class="stat-value">{{ status.has_token ? '有效' : '失效/未获取' }}</span>
                  <div class="status-dot" :class="status.has_token ? 'bg-success' : 'bg-warning'"></div>
                </div>
                <div class="stat-sub">过期: {{ formatExpire(status.expire_at) }}</div>
              </div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="stat-card info-card">
              <div class="stat-icon-wrapper user-bg">
                <icon-user class="stat-icon" />
              </div>
              <div class="stat-info">
                <div class="stat-label">当前商户信息</div>
                <div class="stat-value-row">
                  <span class="stat-value merchant-id">ID: {{ userInfo.business_id || '-' }}</span>
                </div>
                <div class="stat-sub">操作员 ID: {{ userInfo.id || '-' }}</div>
              </div>
            </div>
          </a-col>
        </a-row>
      </div>

      <!-- 主体内容卡片 -->
      <a-card class="main-card" :bordered="false">
        <!-- 工具栏与搜索 -->
        <div class="card-toolbar">
          <div class="left-actions">
            <a-input-search
              v-model="formModel.keyword"
              placeholder="搜索锁名、ID或MAC地址..."
              style="width: 280px"
              allow-clear
              search-button
              @search="handleSearch"
              @press-enter="handleSearch"
            >
              <template #button-icon>
                <icon-search />
              </template>
            </a-input-search>
          </div>
          <div class="right-actions">
            <a-space size="medium">
               <a-tooltip content="从云端同步最新锁列表">
                <a-button type="primary" status="success" @click="handleSync" :loading="syncLoading" class="action-btn">
                  <template #icon><icon-sync /></template>
                  同步数据
                </a-button>
              </a-tooltip>
              <a-divider direction="vertical" />
              <div class="icon-btn-group">
                 <a-tooltip content="刷新列表">
                  <a-button type="text" shape="circle" @click="fetchData">
                    <icon-refresh />
                  </a-button>
                </a-tooltip>
                <a-tooltip content="密度设置" v-if="!isFullscreen">
                   <a-dropdown @select="(val: any) => size = val">
                    <a-button type="text" shape="circle"><icon-menu /></a-button>
                    <template #content>
                      <a-doption value="mini">迷你</a-doption>
                      <a-doption value="medium">中等</a-doption>
                      <a-doption value="large">宽松</a-doption>
                    </template>
                  </a-dropdown>
                </a-tooltip>
              </div>
            </a-space>
          </div>
        </div>

        <!-- 数据表格 -->
        <div class="table-wrapper">
          <a-table
            row-key="ttlock_lock_id"
            :loading="loading"
            :pagination="pagination"
            :data="renderData"
            :columns="columns"
            :bordered="false"
            :size="size"
            :stripe="true"
            :hoverable="true"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #empty>
              <div class="empty-state">
                <icon-lock style="font-size: 48px; color: var(--color-neutral-4); margin-bottom: 16px;" />
                <div class="empty-text">暂无设备数据</div>
                <div class="empty-sub">请尝试同步设备或检查商户配置</div>
                <a-button type="primary" class="mt-4" @click="handleSync">立即同步</a-button>
              </div>
            </template>

            <!-- 自定义列渲染 -->
            <template #lockName="{ record }">
              <div class="cell-lock-main">
                <div class="lock-name">{{ record.lock_name }}</div>
                <div class="lock-id">ID: {{ record.ttlock_lock_id }}</div>
              </div>
            </template>

            <template #bind="{ record }">
              <div class="cell-bind" v-if="record.bind_property_id">
                <a-tag color="green" size="small" class="bind-tag">已绑定</a-tag>
                <div class="bind-info">
                  <div class="bind-title">
                    {{ record.bind_property_title || '房源' }}（{{ record.bind_property_id }}）
                  </div>
                  <div class="bind-sub" v-if="record.bind_property_community_name">
                    {{ record.bind_property_community_name }}
                  </div>
                </div>
              </div>
              <a-tag v-else color="gray" size="small" class="bind-tag">未绑定</a-tag>
            </template>

             <template #mac="{ record }">
                <span class="mac-text">{{ record.lock_mac }}</span>
             </template>

            <template #battery="{ record }">
              <div class="cell-battery">
                <a-progress
                  type="circle"
                  size="mini"
                  :percent="(record.battery || 0) / 100"
                  :status="getBatteryStatus(record.battery)"
                  :stroke-width="4"
                  class="battery-circle"
                />
                <span class="battery-val" :class="getBatteryStatus(record.battery)">{{ record.battery ?? 0 }}%</span>
              </div>
            </template>

            <template #sync="{ record }">
              <div class="cell-sync">
                 <div class="sync-time">{{ formatTimeAgo(record.last_sync_at) }}</div>
                 <div class="sync-full">{{ formatUnix(record.last_sync_at) }}</div>
              </div>
            </template>

            <template #action="{ record }">
              <div class="action-cell">
                <a-button type="text" size="small" class="table-btn-primary" @click="openBind(record)">
                  <template #icon><icon-link /></template>
                  {{ record.bind_property_id ? '更换绑定' : '绑定' }}
                </a-button>
                <a-button type="text" size="small" class="table-btn-default" @click="openCloudDetail(record)">
                  <template #icon><icon-info-circle /></template>
                  详情
                </a-button>
                <a-button type="text" size="small" class="table-btn-default" @click="openUnlockRecords(record)">
                  <template #icon><icon-history /></template>
                  开锁记录
                </a-button>
              </div>
            </template>
          </a-table>
        </div>
      </a-card>
    </div>

    <!-- 绑定弹窗组件 -->
    <BindPropertyModal @register="registerBindModal" @success="fetchData" />

    <!-- 详情抽屉 -->
    <a-drawer
      v-model:visible="detailVisible"
      width="400px"
      :footer="false"
      class="detail-drawer-glass"
    >
      <template #title>
        <span class="drawer-title">设备详情</span>
      </template>
      
      <a-spin :loading="detailLoading" tip="加载中..." style="width: 100%; min-height: 200px;">
        <div v-if="detailLock" class="detail-panel">
          <div class="detail-hero">
            <div class="hero-icon"><icon-lock /></div>
            <div class="hero-title">{{ cloudDetail?.lockName || detailLock.lock_name || '未知设备' }}</div>
            <div class="hero-chip">{{ cloudDetail?.modelNum || detailLock.model_num || '通用型号' }}</div>
          </div>

          <div class="detail-grid">
            <div class="detail-item full-width">
              <label>绑定房源</label>
              <div class="val">
                <template v-if="detailLock.bind_property_id">
                  {{ detailLock.bind_property_title || '房源' }}（{{ detailLock.bind_property_id }}）
                  <span class="val-sub" v-if="detailLock.bind_time"> · {{ formatUnix(detailLock.bind_time) }}</span>
                </template>
                <template v-else>未绑定</template>
              </div>
            </div>
            <div class="detail-item">
              <label>电量状态</label>
              <div class="val highlight" :class="getBatteryStatus(cloudDetail?.battery ?? detailLock.battery)">
                {{ (cloudDetail?.battery ?? detailLock.battery ?? '-') }}%
              </div>
            </div>
            <div class="detail-item">
              <label>MAC 地址</label>
              <div class="val money-font">{{ cloudDetail?.lockMac || detailLock.lock_mac || '-' }}</div>
            </div>
            <div class="detail-item full-width">
              <label>TTLock ID</label>
              <div class="val copy-row">
                 {{ cloudDetail?.lockId || detailLock.ttlock_lock_id || '-' }}
                 <icon-copy class="copy-icon" @click="copyText(String(cloudDetail?.lockId || detailLock.ttlock_lock_id || ''))" />
              </div>
            </div>

            <div class="detail-item full-width">
              <label>lockData（离线插件用）</label>
              <div class="lock-data-actions" v-if="lockDataText">
                <a-space size="mini">
                  <a-button type="text" size="mini" @click="lockDataExpanded = !lockDataExpanded">
                    {{ lockDataExpanded ? '收起' : '展开' }}
                  </a-button>
                  <a-button type="text" size="mini" @click="copyText(lockDataText)">复制</a-button>
                </a-space>
              </div>
              <div class="code-box">
                {{ lockDataDisplay }}
              </div>
            </div>

             <div class="detail-item full-width">
              <label>原始字段解析（raw_json）</label>
              <div v-if="rawParsedItems.length" class="raw-fields">
                <div v-for="it in rawParsedItems" :key="it.key" class="raw-row">
                  <div class="raw-label">{{ it.label }}</div>
                  <div class="raw-value" :class="{ mono: it.mono }">
                    <span class="raw-text">{{ it.value }}</span>
                    <icon-copy class="copy-icon" @click="copyText(it.copyText)" />
                  </div>
                </div>
              </div>
              <div v-else class="empty-raw">无数据</div>
            </div>

            <div class="detail-item full-width">
              <label>raw_json</label>
              <div class="code-box">
                {{ rawJsonPretty }}
              </div>
            </div>
          </div>
        </div>
        <div v-else-if="!detailLoading" class="empty-detail">
          <icon-empty /> 暂无详情数据
        </div>
      </a-spin>
    </a-drawer>

    <!-- 开锁记录抽屉（单独展示） -->
    <a-drawer v-model:visible="unlockVisible" width="980px" :footer="false" class="detail-drawer-glass">
      <template #title>
        <span class="drawer-title">开锁记录</span>
      </template>

      <a-spin :loading="unlockLoading" tip="加载中..." style="width: 100%; min-height: 200px">
        <div v-if="unlockLock" class="detail-panel">
          <div class="detail-hero">
            <div class="hero-icon"><icon-history /></div>
            <div class="hero-title">{{ unlockLock.lock_name || '设备' }}</div>
            <div class="hero-chip">ID: {{ unlockLock.ttlock_lock_id || '-' }}</div>
          </div>

          <div class="detail-grid">
            <div class="detail-item full-width">
              <label>绑定房源</label>
              <div class="val">
                <template v-if="unlockLock.bind_property_id">
                  {{ unlockLock.bind_property_title || '房源' }}（{{ unlockLock.bind_property_id }}）
                  <span class="val-sub" v-if="unlockLock.bind_time"> · {{ formatUnix(unlockLock.bind_time) }}</span>
                </template>
                <template v-else>未绑定</template>
              </div>
            </div>

            <div class="detail-item full-width">
              <label>重要提示</label>
              <a-alert type="warning" show-icon>
                <template #title>云端记录需要同步</template>
                <template #default>
                  云端开锁记录/操作记录需要手机通过蓝牙连接门锁时同步上传；如果未蓝牙连接（或未执行同步），云端列表可能没有数据。
                </template>
              </a-alert>
            </div>

            <div class="detail-item full-width">
              <div class="unlock-toolbar">
                <a-space size="medium" wrap>
                  <a-range-picker
                    v-model="unlockRange"
                    format="YYYY-MM-DD HH:mm"
                    :show-time="true"
                    allow-clear
                    :shortcuts="shortcuts"
                    shortcuts-position="left"
                    :style="{ width: '320px' }"
                    @change="handleUnlockRangeChange"
                  />
                  <a-button type="primary" @click="reloadUnlockData" :loading="unlockLoading">刷新</a-button>
                  <a-button @click="verifyPasscodeUsage" :loading="verifyLoading">核验使用状态</a-button>
                  <span class="val-sub" v-if="verifyStat.scannedPages">
                    已扫描 {{ verifyStat.scannedPages }} 页
                    <template v-if="verifyStat.limited">（已达到上限，建议缩小时间范围）</template>
                  </span>
                </a-space>
              </div>
            </div>

            <div class="detail-item">
              <label>密码申请记录（本地）</label>
              <a-table
                :loading="passcodeLoading"
                :data="passcodeRows"
                :pagination="passcodePagination"
                row-key="id"
                :bordered="{ wrapper: true }"
                size="mini"
                :scroll="{ y: '420px', x: 'max-content' }"
                @page-change="handlePasscodePageChange"
                @page-size-change="handlePasscodePageSizeChange"
              >
                <template #columns>
                  <a-table-column title="申请时间" :width="170">
                    <template #cell="{ record }">
                      <div>{{ formatUnix(record.createtime) }}</div>
                      <div class="val-sub">{{ formatTimeAgo(record.createtime) }}</div>
                    </template>
                  </a-table-column>
                  <a-table-column title="申请人" :width="140">
                    <template #cell="{ record }">
                      <div class="bind-title">{{ record.apply_user_name || '-' }}</div>
                      <div class="bind-sub">#{{ record.user_id || '-' }}</div>
                    </template>
                  </a-table-column>
                  <a-table-column title="房源" :width="220">
                    <template #cell="{ record }">
                      <div class="bind-title">
                        {{ record.property_title || '房源' }}（{{ record.property_id || '-' }}）
                      </div>
                      <div class="bind-sub" v-if="record.property_community_name">
                        {{ record.property_community_name }}
                      </div>
                    </template>
                  </a-table-column>
                  <a-table-column title="密码" :width="180">
                    <template #cell="{ record }">
                      <div class="val copy-row">
                        <span class="mac-text">{{ record._pwd || '-' }}</span>
                        <icon-copy class="copy-icon" @click="copyText(record._pwd)" />
                      </div>
                      <a-link v-if="record._start && record._end" @click="focusPasscode(record)" class="val-sub">
                        查云端
                      </a-link>
                    </template>
                  </a-table-column>
                  <a-table-column title="有效期" :width="240">
                    <template #cell="{ record }">
                      <div>{{ record._start ? formatMs(record._start) : '-' }}</div>
                      <div class="val-sub">{{ record._end ? formatMs(record._end) : '-' }}</div>
                    </template>
                  </a-table-column>
                  <a-table-column title="状态" :width="190">
                    <template #cell="{ record }">
                      <template v-if="record._status === 'used'">
                        <a-tag color="green" size="small">已使用</a-tag>
                        <div class="val-sub">使用：{{ formatMs(record._usedAt) }}</div>
                      </template>
                      <template v-else-if="record._status === 'expired'">
                        <a-tag color="gray" size="small">已过期</a-tag>
                      </template>
                      <template v-else-if="record._status === 'unknown'">
                        <a-tag color="orange" size="small">未确认</a-tag>
                        <div class="val-sub">建议缩小时间范围后核验</div>
                      </template>
                      <template v-else>
                        <a-tag color="orange" size="small">未发现使用记录</a-tag>
                      </template>
                      <div class="val-sub" v-if="Number(record.result) === 0 && record.err_msg">
                        失败：{{ record.err_msg }}
                      </div>
                    </template>
                  </a-table-column>
                </template>
              </a-table>
              <a-empty v-if="!passcodeRows.length" description="暂无申请记录" />
            </div>

            <div class="detail-item">
              <label>开锁记录（云端）</label>
              <a-table
                :loading="lockRecordsLoading"
                :data="lockRecordRows"
                :pagination="lockRecordPagination"
                row-key="_rowkey"
                :bordered="{ wrapper: true }"
                size="mini"
                :scroll="{ y: '420px', x: 'max-content' }"
                @page-change="handleLockRecordPageChange"
                @page-size-change="handleLockRecordPageSizeChange"
              >
                <template #columns>
                  <a-table-column title="时间" :width="180">
                    <template #cell="{ record }">
                      <div>{{ formatMs(record._time) }}</div>
                      <div class="val-sub">{{ formatMsAgo(record._time) }}</div>
                    </template>
                  </a-table-column>
                  <a-table-column title="类型" :width="210">
                    <template #cell="{ record }">
                      <a-tag :color="Number(record.recordType) === 4 ? 'arcoblue' : 'gray'" size="small">
                        {{ getRecordTypeLabel(record.recordType) }}
                      </a-tag>
                    </template>
                  </a-table-column>
                  <a-table-column title="结果" :width="90">
                    <template #cell="{ record }">
                      <a-tag :color="Number(record.success) === 1 ? 'green' : 'red'" size="small">
                        {{ Number(record.success) === 1 ? '成功' : '失败' }}
                      </a-tag>
                    </template>
                  </a-table-column>
                  <a-table-column title="密码/钥匙" :width="200">
                    <template #cell="{ record }">
                      <div class="val copy-row">
                        <span class="mac-text">{{ record.keyboardPwd || '-' }}</span>
                        <icon-copy class="copy-icon" @click="copyText(String(record.keyboardPwd || ''))" />
                      </div>
                    </template>
                  </a-table-column>
                  <a-table-column title="操作者" :width="160">
                    <template #cell="{ record }">
                      {{ record.username || record.operator || '-' }}
                    </template>
                  </a-table-column>
                </template>
              </a-table>
              <a-empty v-if="!lockRecordRows.length" description="暂无开锁记录" />
            </div>
          </div>
        </div>

        <div v-else-if="!unlockLoading" class="empty-detail">
          <icon-empty /> 暂无记录数据
        </div>
      </a-spin>
    </a-drawer>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue';
import useLoading from '@/hooks/loading';
import { Message } from '@arco-design/web-vue';
import { Pagination } from '@/types/global';
import { SizeProps } from '/@/components/tabletool';
import { useModal } from '/@/components/Modal';
import useUserStore from '@/store/modules/user';
import BindPropertyModal from './modal/BindPropertyModal.vue';
import { getLockDetail, getLockList, getStatus, syncLocks, getLockRecords, getPasscodeEvents } from './api';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';
import { useClipboard } from '@vueuse/core';
import { shortcuts } from '@/utils/dayjs';

dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

const { copy } = useClipboard();
const { loading, setLoading } = useLoading(true);
const syncLoading = ref(false);
const isFullscreen = ref(false);
const size = ref<SizeProps>('large');
const userStore = useUserStore();
const userInfo = computed(() => userStore.userInfo);

const basePagination: Pagination = { current: 1, pageSize: 20 };
const pagination = reactive({
  ...basePagination,
  total: 0,
  showTotal: true,
  showPageSize: true,
});

const columns = [
  { title: '锁设备信息', slotName: 'lockName', width: 220 },
  { title: '绑定房源', slotName: 'bind', width: 260 },
  { title: 'MAC 地址', slotName: 'mac', width: 160 },
  { title: '电量', slotName: 'battery', width: 120, align: 'center' },
  { title: '最近同步', slotName: 'sync', width: 180 },
  { title: '操作', slotName: 'action', width: 240, fixed: 'right', align: 'center' },
];

const renderData = ref<any[]>([]);
const status = reactive<any>({ config_ok: false, has_token: false, api_base: '', expire_at: 0 });
const formModel = reactive({ keyword: '' });

// API Actions
const fetchStatus = async () => {
  try {
    const resp: any = await getStatus();
    Object.assign(status, (resp?.data ?? resp) || {});
  } catch (e) {
    console.error('Status Error', e);
  }
};

const fetchData = async () => {
  setLoading(true);
  try {
    const resp: any = await getLockList({ page: pagination.current, pageSize: pagination.pageSize });
    const data = resp?.items ? resp : (resp?.data?.items ? resp.data : resp?.data ?? resp);
    const items = (data?.items || []) as any[];
    
    // 前端关键词过滤
    const kw = String(formModel.keyword || '').trim().toLowerCase();
    renderData.value = kw 
      ? items.filter(it =>
          `${it.lock_name} ${it.ttlock_lock_id} ${it.lock_mac} ${it.bind_property_title || ''} ${it.bind_property_id || ''}`
            .toLowerCase()
            .includes(kw)
        )
      : items;

    pagination.current = Number(data?.page || pagination.current);
    pagination.total = Number(data?.total || 0);
  } catch (e: any) {
    Message.error(e?.message || '加载列表失败');
    renderData.value = [];
  } finally {
    setLoading(false);
  }
};

// Handlers
const handleSearch = () => { pagination.current = 1; fetchData(); };
const handlePageChange = (p: number) => { pagination.current = p; fetchData(); };
const handlePageSizeChange = (ps: number) => { pagination.pageSize = ps; handleSearch(); };

const handleSync = async () => {
  syncLoading.value = true;
  try {
    await syncLocks();
    Message.success('同步指令已下发');
    setTimeout(() => { fetchStatus(); fetchData(); syncLoading.value = false; }, 1500);
  } catch (e: any) {
    syncLoading.value = false;
    Message.error('同步失败: ' + (e?.message || '未知错误'));
  }
};

const [registerBindModal, { openModal: openBindModal }] = useModal();
const openBind = (row: any) => openBindModal(true, { lock: row });

const detailVisible = ref(false);
const detailLoading = ref(false);
const cloudDetail = ref<any>(null);
const detailLock = ref<any>(null);
const lockDataExpanded = ref(false);

// 开锁记录（独立抽屉）
const unlockVisible = ref(false);
const unlockLock = ref<any>(null);
const unlockRange = ref<any>(undefined);
const unlockLoading = computed(() => lockRecordsLoading.value || passcodeLoading.value);

const lockRecordsLoading = ref(false);
const lockRecords = ref<any[]>([]);
const lockRecordPagination = reactive({
  current: 1,
  pageSize: 50,
  total: 0,
  showTotal: true,
  showPageSize: true,
});
const passcodeLoading = ref(false);
const passcodeEvents = ref<any[]>([]);
const passcodePagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showTotal: true,
  showPageSize: true,
});

const verifyLoading = ref(false);
const verifyStat = reactive<{ scannedPages: number; limited: boolean }>({ scannedPages: 0, limited: false });
const passcodeUsedAtByPwd = ref<Record<string, number>>({});

const resetVerifyState = () => {
  verifyStat.scannedPages = 0;
  verifyStat.limited = false;
  passcodeUsedAtByPwd.value = {};
};

const unwrapRespData = (resp: any) => {
  const r1 = resp?.data ?? resp;
  return r1?.data ?? r1;
};

const getUnlockRangeMs = () => {
  const r = unlockRange.value;
  if (!Array.isArray(r) || r.length !== 2 || !r[0] || !r[1]) return { start: 0, end: 0 };
  const start = dayjs(r[0]).valueOf();
  const end = dayjs(r[1]).valueOf();
  if (!Number.isFinite(start) || !Number.isFinite(end) || start <= 0 || end <= 0) return { start: 0, end: 0 };
  return start <= end ? { start, end } : { start: end, end: start };
};

const fetchPasscodeList = async () => {
  const lockId = Number(unlockLock.value?.ttlock_lock_id || 0);
  if (!lockId) {
    passcodeEvents.value = [];
    passcodePagination.total = 0;
    return;
  }

  passcodeLoading.value = true;
  try {
    const resp: any = await getPasscodeEvents({
      ttlock_lock_id: lockId,
      page: passcodePagination.current,
      pageSize: passcodePagination.pageSize,
    });
    const data = unwrapRespData(resp) || {};
    passcodeEvents.value = (data?.items || []) as any[];
    passcodePagination.total = Number(data?.total || 0);
    passcodePagination.current = Number(data?.page || passcodePagination.current);
    passcodePagination.pageSize = Number(data?.pageSize || passcodePagination.pageSize);
  } catch (e: any) {
    passcodeEvents.value = [];
    passcodePagination.total = 0;
    Message.error(e?.message || '加载密码申请记录失败');
  } finally {
    passcodeLoading.value = false;
  }
};

const fetchLockRecordPage = async () => {
  const lockId = Number(unlockLock.value?.ttlock_lock_id || 0);
  if (!lockId) {
    lockRecords.value = [];
    lockRecordPagination.total = 0;
    return;
  }

  lockRecordsLoading.value = true;
  try {
    const { start, end } = getUnlockRangeMs();
    const resp: any = await getLockRecords({
      ttlock_lock_id: lockId,
      ...(start > 0 ? { start_date: start } : {}),
      ...(end > 0 ? { end_date: end } : {}),
      page: lockRecordPagination.current,
      pageSize: lockRecordPagination.pageSize,
    });
    const data = unwrapRespData(resp) || {};
    lockRecords.value = (data?.list || []) as any[];
    lockRecordPagination.total = Number(data?.total || 0);
    lockRecordPagination.current = Number(data?.pageNo || data?.page || lockRecordPagination.current);
    lockRecordPagination.pageSize = Number(data?.pageSize || lockRecordPagination.pageSize);
  } catch (e: any) {
    lockRecords.value = [];
    lockRecordPagination.total = 0;
    Message.error(e?.message || '加载开锁记录失败');
  } finally {
    lockRecordsLoading.value = false;
  }
};

const verifyPasscodeUsage = async () => {
  const lockId = Number(unlockLock.value?.ttlock_lock_id || 0);
  if (!lockId) return;

  const pwdSet = new Set<string>();
  for (const it of passcodeEvents.value || []) {
    const pwd = String((it as any)?.keyboard_pwd || (it as any)?.keyboardPwd || '').trim();
    if (pwd) pwdSet.add(pwd);
  }

  resetVerifyState();
  if (pwdSet.size === 0) return;

  verifyLoading.value = true;
  try {
    const { start, end } = getUnlockRangeMs();
    const scanPageSize = 100;
    const maxPages = 20;
    let totalPages = 0;
    let lastLen = 0;
    let allMatched = false;

    for (let pageNo = 1; pageNo <= maxPages; pageNo++) {
      const resp: any = await getLockRecords({
        ttlock_lock_id: lockId,
        ...(start > 0 ? { start_date: start } : {}),
        ...(end > 0 ? { end_date: end } : {}),
        page: pageNo,
        pageSize: scanPageSize,
      });
      const data = unwrapRespData(resp) || {};
      const list = (data?.list || []) as any[];
      lastLen = list.length;

      if (pageNo === 1) {
        const total = Number(data?.total || 0);
        if (total > 0) totalPages = Math.ceil(total / scanPageSize);
      }

      for (const r of list) {
        const pwd = String(r?.keyboardPwd || '').trim();
        if (!pwd || !pwdSet.has(pwd)) continue;

        const usedAt = Number(r?.lockDate || r?.serverDate || 0);
        if (!Number.isFinite(usedAt) || usedAt <= 0) continue;

        const recordType = Number(r?.recordType || 0);
        const success = Number(r?.success || 0) === 1;
        if (!success || recordType !== 4) continue;

        const exist = passcodeUsedAtByPwd.value[pwd];
        if (!exist || usedAt > exist) passcodeUsedAtByPwd.value[pwd] = usedAt;
      }

      verifyStat.scannedPages = pageNo;

      if (Object.keys(passcodeUsedAtByPwd.value).length >= pwdSet.size) {
        allMatched = true;
        break;
      }
      if (totalPages > 0 && pageNo >= totalPages) break;
      if (totalPages === 0 && list.length < scanPageSize) break;
    }

    const scanned = verifyStat.scannedPages;
    verifyStat.limited =
      !allMatched &&
      ((totalPages > 0 && scanned < totalPages) || (totalPages === 0 && scanned >= maxPages && lastLen >= scanPageSize));
  } catch (e: any) {
    passcodeUsedAtByPwd.value = {};
    verifyStat.limited = true;
    Message.error(e?.message || '核验使用状态失败');
  } finally {
    verifyLoading.value = false;
  }
};

const reloadUnlockData = async () => {
  resetVerifyState();
  await Promise.all([fetchPasscodeList(), fetchLockRecordPage()]);
  await verifyPasscodeUsage();
};

const handleUnlockRangeChange = async () => {
  lockRecordPagination.current = 1;
  resetVerifyState();
  await fetchLockRecordPage();
  await verifyPasscodeUsage();
};

const handlePasscodePageChange = async (p: number) => {
  passcodePagination.current = p;
  await fetchPasscodeList();
  await verifyPasscodeUsage();
};

const handlePasscodePageSizeChange = async (ps: number) => {
  passcodePagination.pageSize = ps;
  passcodePagination.current = 1;
  await fetchPasscodeList();
  await verifyPasscodeUsage();
};

const handleLockRecordPageChange = async (p: number) => {
  lockRecordPagination.current = p;
  await fetchLockRecordPage();
};

const handleLockRecordPageSizeChange = async (ps: number) => {
  lockRecordPagination.pageSize = ps;
  lockRecordPagination.current = 1;
  await fetchLockRecordPage();
};

const focusPasscode = async (record: any) => {
  const start = Number(record?._start || 0);
  const end = Number(record?._end || 0);
  if (!Number.isFinite(start) || !Number.isFinite(end) || start <= 0 || end <= 0) return;

  const pad = 60 * 60 * 1000;
  unlockRange.value = [new Date(Math.max(start - pad, 0)), new Date(end + pad)];
  await handleUnlockRangeChange();
};

type RawItem = {
  key: string;
  label: string;
  value: string;
  copyText: string;
  mono?: boolean;
};

const rawJsonObj = computed<Record<string, any> | null>(() => {
  const raw = detailLock.value?.raw_json;
  if (!raw) return null;
  if (typeof raw === 'object') return raw as Record<string, any>;
  try {
    return JSON.parse(String(raw)) as Record<string, any>;
  } catch {
    return null;
  }
});

const formatRawValue = (v: any) => {
  if (v == null) return '-';
  if (typeof v === 'string') return v.trim() || '-';
  if (typeof v === 'number' || typeof v === 'boolean') return String(v);
  try {
    return JSON.stringify(v);
  } catch {
    return String(v);
  }
};

const lockDataText = computed(() => {
  const v = rawJsonObj.value?.lockData ?? cloudDetail.value?.lockData;
  const s = typeof v === 'string' ? v.trim() : String(v || '').trim();
  return s || '';
});

const lockDataDisplay = computed(() => {
  if (!lockDataText.value) return '无数据';
  if (lockDataExpanded.value) return lockDataText.value;
  const s = lockDataText.value;
  return s.length > 220 ? s.slice(0, 220) + '…' : s;
});

const rawFieldLabelMap: Record<string, string> = {
  protocolType: '协议类型（protocolType）',
  lockVersion: '锁版本（lockVersion）',
  specialValue: 'specialValue',
  timezoneRawOffset: '时区偏移（timezoneRawOffset）',
  featureValue: 'featureValue',
  lockFlagPos: 'lockFlagPos',
  lockType: 'lockType',
  keyboardPwdVersion: '键盘密码版本（keyboardPwdVersion）',
  remoteEnable: '远程功能（remoteEnable）',
  hasGateway: '是否有网关（hasGateway）',
};

const rawParsedItems = computed<RawItem[]>(() => {
  const obj = rawJsonObj.value;
  if (!obj) return [];

  const ignore = new Set(['lockId', 'lockName', 'lockMac', 'battery', 'modelNum', 'lockData']);
  const priority = [
    'protocolType',
    'lockVersion',
    'specialValue',
    'timezoneRawOffset',
    'featureValue',
    'lockFlagPos',
    'lockType',
    'keyboardPwdVersion',
  ];

  const out: RawItem[] = [];
  const add = (key: string) => {
    if (!key || ignore.has(key) || !(key in obj)) return;
    const val = formatRawValue(obj[key]);
    out.push({
      key,
      label: rawFieldLabelMap[key] || key,
      value: val,
      copyText: val === '-' ? '' : val,
      mono: typeof obj[key] === 'number' || typeof obj[key] === 'boolean',
    });
  };

  for (const k of priority) add(k);
  for (const k of Object.keys(obj).sort()) {
    if (out.length >= 12) break;
    if (priority.includes(k)) continue;
    add(k);
  }
  return out;
});

const rawJsonPretty = computed(() => {
  const raw = detailLock.value?.raw_json;
  if (!raw) return '无数据';
  try {
    const obj = typeof raw === 'string' ? JSON.parse(raw) : raw;
    return JSON.stringify(obj, null, 2);
  } catch {
    return String(raw);
  }
});

const formatMs = (ms: any) => {
  const v = Number(ms || 0);
  if (!Number.isFinite(v) || v <= 0) return '-';
  return dayjs(v).format('YYYY-MM-DD HH:mm:ss');
};
const formatMsAgo = (ms: any) => {
  const v = Number(ms || 0);
  if (!Number.isFinite(v) || v <= 0) return '-';
  return dayjs(v).fromNow();
};
// 云端锁操作记录类型（recordType）说明：来自 TTLock 官方文档（云端获取 lockRecord/list）
const cloudRecordTypeLabelMap: Record<number, string> = {
  1: '蓝牙开锁',
  4: '密码开锁',
  5: '车位锁升',
  6: '车位锁降',
  7: 'IC卡开锁',
  8: '指纹开锁',
  9: '手环开锁',
  10: '机械钥匙开锁',
  11: '蓝牙闭锁',
  12: '网关开锁',
  29: '非法开锁',
  30: '门磁合上（关门）',
  31: '门磁打开（开门）',
  32: '从内部打开门',
  33: '指纹关锁',
  34: '密码关锁',
  35: 'IC卡关锁',
  36: '机械钥匙关锁',
  37: 'APP按键控制锁',
  42: '邮局本地邮件',
  43: '邮局外地邮件',
  44: '防撬报警',
  45: '超时自动闭锁',
  46: '开锁按键开锁',
  47: '闭锁按键闭锁',
  48: '系统被锁定',
  49: '酒店卡开锁',
  50: '高温开锁',
  51: '用已删除的卡开锁',
  52: '用APP锁定锁',
  53: '用密码锁定锁',
  54: '车离开（车位锁）',
  55: '遥控开/闭锁',
  57: '二维码开锁成功',
  58: '二维码开锁失败（过期）',
  59: '开启反锁',
  60: '关闭反锁',
  61: '二维码闭锁成功',
  62: '二维码开锁失败（门已反锁）',
  63: '常开时间段自动开锁',
  64: '门未关报警',
  65: '开锁超时',
  66: '闭锁超时',
  67: '3D人脸开锁成功',
  68: '3D人脸开锁失败（门反锁）',
  69: '3D人脸闭锁',
  71: '3D人脸开门失败（过期/未生效）',
  75: 'App授权按键开锁成功',
  76: '网关授权按键开锁成功',
  77: '双重认证蓝牙开锁验证成功（等待第二用户）',
  78: '双重认证密码开锁验证成功（等待第二用户）',
  79: '双重认证指纹开锁验证成功（等待第二用户）',
  80: '双重认证IC卡开锁验证成功（等待第二用户）',
  81: '双重认证人脸卡开锁验证成功（等待第二用户）',
  82: '双重认证遥控开锁验证成功（等待第二用户）',
  83: '双重认证掌静脉开锁验证成功（等待第二用户）',
  84: '掌静脉开锁成功',
  85: '掌静脉开锁失败（门反锁）',
  86: '掌静脉闭锁',
  88: '掌静脉开门失败（过期/未生效）',
  92: '管理员密码开锁',
};

const getRecordTypeLabel = (t: any) => {
  const v = Number(t || 0);
  if (!v) return '-';
  const label = cloudRecordTypeLabelMap[v];
  return label ? `${label}（${v}）` : `类型 ${v}`;
};

const passcodeRows = computed(() => {
  const now = Date.now();
  const list = (passcodeEvents.value || []) as any[];
  return list.map((it: any) => {
    const pwd = String(it?.keyboard_pwd || it?.keyboardPwd || '').trim();
    const start = Number(it?.start_date || it?.startDate || 0);
    const end = Number(it?.end_date || it?.endDate || 0);
    const usedAt = pwd ? Number(passcodeUsedAtByPwd.value[pwd] || 0) : 0;
    let status: 'used' | 'expired' | 'unused' | 'unknown' = 'unused';
    if (usedAt > 0) {
      status = 'used';
    } else if (verifyStat.limited) {
      status = 'unknown';
    } else if (end > 0 && end < now) {
      status = 'expired';
    } else {
      status = 'unused';
    }
    return { ...it, _pwd: pwd, _start: start, _end: end, _usedAt: usedAt, _status: status };
  });
});

const lockRecordRows = computed(() => {
  const list = (lockRecords.value || []) as any[];
  return list
    .map((it: any, idx: number) => {
      const t = Number(it?.lockDate || it?.serverDate || 0);
      return { ...it, _time: t, _rowkey: String(it?.recordId || t || idx) };
    })
    .sort((a: any, b: any) => Number(b._time || 0) - Number(a._time || 0));
});

const openCloudDetail = async (row: any) => {
  detailVisible.value = true;
  detailLoading.value = true;
  cloudDetail.value = null;
  detailLock.value = row || null;
  lockDataExpanded.value = false;
  const lockId = Number(row?.ttlock_lock_id || 0);
  try {
    const resp: any = await getLockDetail({ ttlock_lock_id: lockId });
    cloudDetail.value = unwrapRespData(resp) || null;
  } catch (e: any) {
    cloudDetail.value = null;
    Message.error(e?.message || '加载云端详情失败');
  } finally {
    detailLoading.value = false;
  }
};

const openUnlockRecords = async (row: any) => {
  unlockVisible.value = true;
  unlockLock.value = row || null;
  lockRecords.value = [];
  passcodeEvents.value = [];

  resetVerifyState();
  passcodePagination.current = 1;
  lockRecordPagination.current = 1;
  unlockRange.value = [dayjs().subtract(7, 'day').startOf('day').toDate(), dayjs().toDate()];

  await reloadUnlockData();
};

const copyText = (text: string) => {
  if (!text) return;
  copy(text);
  Message.success('已复制');
};

// Formatters
const formatUnix = (s: any) => !s ? '-' : dayjs(Number(s) * 1000).format('YYYY-MM-DD HH:mm');
const formatTimeAgo = (s: any) => !s ? '-' : dayjs(Number(s) * 1000).fromNow();
const formatExpire = (s: any) => {
  if (!s) return '-';
  const isExp = Date.now() / 1000 > Number(s);
  return formatUnix(s) + (isExp ? ' (已过期)' : '');
};
const getBatteryStatus = (b: any) => {
  const val = Number(b || 0);
  return val <= 20 ? 'danger' : val <= 50 ? 'warning' : 'success';
};

onMounted(() => { fetchStatus(); fetchData(); });
</script>

<style scoped lang="less">
// page layout
.page-container {
  padding: 20px;
  background-color: var(--color-fill-1);
  min-height: 100vh;
}

.content-wrapper {
  max-width: 1600px; // Pro Max: constraint layout on ultra-wide screens
  margin: 0 auto;
}

// Header Stats Cards
.dashboard-header {
  margin-bottom: 24px;
}

.stat-card {
  background: var(--color-bg-2);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;
  border: 1px solid transparent;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);
  }

  &.is-active {
    border-color: rgba(var(--success-6), 0.3);
    background: linear-gradient(135deg, var(--color-bg-2) 0%, rgba(var(--success-6), 0.05) 100%);
  }

  &.info-card {
     border-color: rgba(var(--arcoblue-6), 0.3);
      background: linear-gradient(135deg, var(--color-bg-2) 0%, rgba(var(--arcoblue-6), 0.05) 100%);
  }

  .stat-icon-wrapper {
    width: 56px;
    height: 56px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    
    &.config-bg { background-color: rgba(var(--primary-6), 0.1); color: rgb(var(--primary-6)); }
    &.token-bg { background-color: rgba(var(--success-6), 0.1); color: rgb(var(--success-6)); }
    &.user-bg { background-color: rgba(var(--arcoblue-6), 0.1); color: rgb(var(--arcoblue-6)); }
    
    .stat-icon {
      font-size: 28px;
    }
  }

  .stat-info {
    flex: 1;
    overflow: hidden;
    
    .stat-label {
      font-size: 13px;
      color: var(--color-text-3);
      margin-bottom: 4px;
    }
    
    .stat-value-row {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 4px;
      
      .stat-value {
        font-size: 18px;
        font-weight: 700;
        color: var(--color-text-1);
        line-height: 1.2;
      }
      
      .status-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        &.bg-success { background-color: rgb(var(--success-6)); box-shadow: 0 0 4px rgb(var(--success-6)); }
        &.bg-danger { background-color: rgb(var(--danger-6)); }
        &.bg-warning { background-color: rgb(var(--warning-6)); }
      }
    }
    
    .stat-sub {
      font-size: 12px;
      color: var(--color-text-4);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      opacity: 0.8;
    }
  }
}

// Main Card Styling
.main-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.02);
  transition: all 0.3s;
  
  :deep(.arco-card-body) {
    padding: 24px;
  }
}

.card-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  
  .left-actions {
    flex: 1;
  }
  .right-actions {
    margin-left: 16px;
  }
  
  .action-btn {
    border-radius: 6px;
    font-weight: 500;
  }
  
  .icon-btn-group {
    display: flex;
    gap: 8px;
    color: var(--color-text-2);
  }
}

// Table Styling Overrides
.table-wrapper {
  // Add sleek styling to table
  :deep(.arco-table-container) {
    border-radius: 8px;
  }
  :deep(.arco-table-th) {
    background-color: var(--color-fill-1);
    font-weight: 600;
  }
}

// Cell Renderers
.cell-lock-main {
  .lock-name {
    font-size: 15px;
    font-weight: 600;
    color: var(--color-text-1);
  }
  .lock-id {
    font-size: 12px;
    color: var(--color-text-3);
    font-family: 'Menlo', 'Monaco', monospace; // Tech feel
  }
}

.cell-bind {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  min-width: 0;

  .bind-tag {
    flex-shrink: 0;
    margin-top: 2px;
  }

  .bind-info {
    min-width: 0;
  }

  .bind-title {
    font-size: 13px;
    font-weight: 600;
    color: var(--color-text-1);
    line-height: 18px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .bind-sub {
    font-size: 12px;
    color: var(--color-text-3);
    line-height: 16px;
    margin-top: 2px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.mac-text {
  font-family: monospace;
  background: var(--color-fill-2);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  color: var(--color-text-2);
}

.cell-battery {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  
  .battery-val {
    font-size: 12px;
    font-weight: 700;
    
    &.danger { color: rgb(var(--danger-6)); }
    &.warning { color: rgb(var(--warning-6)); }
    &.success { color: rgb(var(--success-6)); }
  }
}

.cell-sync {
  line-height: 1.4;
  .sync-time {
    color: var(--color-text-1);
    font-weight: 500;
  }
  .sync-full {
    font-size: 12px;
    color: var(--color-text-4);
  }
}

.action-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  
  .table-btn-primary {
    color: rgb(var(--primary-6));
    &:hover { background-color: var(--color-fill-2); }
  }
  .table-btn-default {
    color: var(--color-text-2);
    &:hover { color: var(--color-text-1); background-color: var(--color-fill-2); }
  }
}

// Empty State
.empty-state {
  padding: 40px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  
  .empty-text {
    font-size: 16px;
    color: var(--color-text-1);
    font-weight: 500;
    margin-bottom: 4px;
  }
  .empty-sub {
    font-size: 13px;
    color: var(--color-text-3);
  }
  .mt-4 {
    margin-top: 16px;
  }
}

// Detail Drawer Styling
.detail-drawer-glass {
  // Optional glass effect for drawer
}

.detail-panel {
  padding: 16px 4px;
}

.detail-hero {
  text-align: center;
  margin-bottom: 24px;
  
  .hero-icon {
    width: 64px;
    height: 64px;
    margin: 0 auto 12px;
    background: linear-gradient(135deg, rgb(var(--primary-6)), rgb(var(--primary-5)));
    border-radius: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-size: 32px;
    box-shadow: 0 8px 16px rgba(var(--primary-6), 0.25);
  }
  
  .hero-title {
    font-size: 18px;
    font-weight: 700;
    color: var(--color-text-1);
    margin-bottom: 4px;
  }
  
  .hero-chip {
    display: inline-block;
    padding: 2px 10px;
    background: var(--color-fill-2);
    border-radius: 12px;
    font-size: 12px;
    color: var(--color-text-3);
  }
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.detail-item {
  background: var(--color-fill-1);
  padding: 12px;
  border-radius: 8px;
  
  &.full-width {
    grid-column: span 2;
  }
  
  label {
    display: block;
    font-size: 12px;
    color: var(--color-text-3);
    margin-bottom: 6px;
  }
  
  .val {
    font-size: 14px;
    font-weight: 500;
    color: var(--color-text-1);

    .val-sub {
      font-size: 12px;
      font-weight: 400;
      color: var(--color-text-3);
    }
    
    &.highlight.success { color: rgb(var(--success-6)); }
    &.highlight.warning { color: rgb(var(--warning-6)); }
    &.highlight.danger { color: rgb(var(--danger-6)); }
    
    &.money-font {
      font-family: monospace;
    }
    
    &.copy-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      
      .copy-icon {
        cursor: pointer;
        color: var(--color-text-4);
        &:hover { color: rgb(var(--primary-6)); }
      }
    }
  }
  
  .code-box {
    font-family: monospace;
    font-size: 11px;
    color: var(--color-text-3);
    white-space: pre-wrap;
    word-break: break-word;
    max-height: 180px;
    overflow-y: auto;
  }
}

.lock-data-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: -2px;
  margin-bottom: 8px;
}

.raw-fields {
  display: grid;
  grid-template-columns: 1fr;
  gap: 10px;
}

.raw-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.raw-label {
  width: 160px;
  flex-shrink: 0;
  font-size: 12px;
  color: var(--color-text-3);
  line-height: 18px;
}

.raw-value {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
  font-size: 13px;
  color: var(--color-text-1);
  line-height: 18px;

  &.mono {
    font-family: monospace;
  }

  .raw-text {
    word-break: break-word;
  }

  .copy-icon {
    flex-shrink: 0;
    cursor: pointer;
    color: var(--color-text-4);
    margin-top: 1px;
    &:hover { color: rgb(var(--primary-6)); }
  }
}

.empty-raw {
  font-size: 12px;
  color: var(--color-text-4);
}

.empty-detail {
  padding: 40px;
  text-align: center;
  color: var(--color-text-3);
}
</style>
