<template>
  <div class="container">
    <page-card breadcrumb fullTable :isFullscreen="isFullscreen" @onheight="handlePageHeight">
      <template #searchLeft>
        <a-space>
          <a-button type="primary" @click="fetchTables" v-permission="['view']">
            <template #icon>
              <icon-refresh />
            </template>
            刷新
          </a-button>
          <a-input style="width: 240px" v-model="selectedTable" placeholder="选择/输入表名" allow-clear />
          <a-button @click="handleColumns" v-permission="['view']">字段</a-button>
          <a-button @click="openAddColumnModal" v-permission="['create']">新增字段</a-button>
          <a-button @click="handleCreateSQL" v-permission="['view']">建表SQL</a-button>
          <a-popconfirm content="您确定要删除该表吗？（仅开发环境允许）" @ok="handleDrop">
            <a-button status="danger" v-permission="['drop']">删除表</a-button>
          </a-popconfirm>
        </a-space>
      </template>
      <template #searchRight>
        <tabletool :showbtn="['fullscreen']" @fullscreen="(v)=>isFullscreen=v"></tabletool>
      </template>

      <template #table>
        <a-row :gutter="12" style="height: 100%" class="dbtablemgr-row">
          <a-col :xs="24" :md="10" style="height: 100%">
            <a-card title="数据表" :bordered="false" class="dbtablemgr-left-card">
              <a-table
                row-key="name"
                :data="tables"
                :pagination="false"
                :loading="loadingTables"
                :scroll="{ y: '100%' }"
                @row-click="(r:any)=>{ selectedTable = r.name }"
                :bordered="{ wrapper: true, cell: true }"
              >
                <template #columns>
                  <a-table-column title="表名" data-index="name" :width="220" />
                  <a-table-column title="备注" data-index="title" />
                </template>
              </a-table>
            </a-card>
          </a-col>

          <a-col :xs="24" :md="14" style="height: 100%">
            <a-scrollbar class="dbtablemgr-right-scroll" style="overflow: auto" :style="{ height: `${rightScrollHeight}px` }">
              <div class="dbtablemgr-right-inner">
                <a-card title="建表" :bordered="false" style="margin-bottom: 12px">
                  <a-tabs>
                <a-tab-pane key="visual" title="可视化建表">
                  <a-alert type="warning" style="margin-bottom: 10px">
                    建议仅开发环境使用。点击“预览SQL”确认无误后再创建。
                  </a-alert>

                  <a-form :model="tableForm" layout="inline" style="margin-bottom: 10px">
                    <a-form-item field="tableName" label="表名" required>
                      <a-input v-model="tableForm.tableName" style="width: 200px" placeholder="例如: test_user" />
                    </a-form-item>
                    <a-form-item field="tableComment" label="表备注">
                      <a-input v-model="tableForm.tableComment" style="width: 240px" placeholder="可选" />
                    </a-form-item>
                    <a-form-item field="engine" label="引擎">
                      <a-select v-model="tableForm.engine" style="width: 140px" allow-clear>
                        <a-option value="InnoDB">InnoDB</a-option>
                        <a-option value="MyISAM">MyISAM</a-option>
                      </a-select>
                    </a-form-item>
                    <a-form-item field="charset" label="字符集">
                      <a-select v-model="tableForm.charset" style="width: 160px" allow-clear>
                        <a-option value="utf8mb4">utf8mb4</a-option>
                        <a-option value="utf8">utf8</a-option>
                      </a-select>
                    </a-form-item>
                  </a-form>

                  <div style="margin-bottom: 8px">
                    <a-space>
                      <a-button type="primary" @click="addFieldRow">添加字段</a-button>
                      <a-button @click="addDefaultIdField">一键添加ID主键</a-button>
                      <a-button :loading="loadingPreview" @click="handlePreviewSQL">预览SQL</a-button>
                      <a-button type="primary" status="success" :loading="loadingCreateVisual" @click="handleCreateTable" v-permission="['create']">
                        创建表
                      </a-button>
                    </a-space>
                  </div>

                  <a-table :data="fields" :pagination="false" row-key="_rowid" :bordered="{ wrapper: true, cell: true }" :scroll="{ y: '220px' }">
                    <template #columns>
                      <a-table-column title="字段名" :width="140">
                        <template #cell="{ record }">
                          <a-input v-model="record.name" placeholder="例如: id" />
                        </template>
                      </a-table-column>
                      <a-table-column title="类型" :width="140">
                        <template #cell="{ record }">
                          <a-select v-model="record.type" allow-clear>
                            <a-option v-for="t in typeOptions" :key="t.value" :value="t.value">{{ t.label }}</a-option>
                          </a-select>
                        </template>
                      </a-table-column>
                      <a-table-column title="长度" :width="100">
                        <template #cell="{ record }">
                          <a-input v-model="record.length" :placeholder="getLengthPlaceholder(record.type)" />
                        </template>
                      </a-table-column>
                      <a-table-column title="可空" :width="70">
                        <template #cell="{ record }">
                          <a-switch v-model="record.nullable" />
                        </template>
                      </a-table-column>
                      <a-table-column title="无符号" :width="80">
                        <template #cell="{ record }">
                          <a-switch v-model="record.unsigned" />
                        </template>
                      </a-table-column>
                      <a-table-column title="默认值" :width="140">
                        <template #cell="{ record }">
                          <a-input v-model="record.default" placeholder="可选" />
                        </template>
                      </a-table-column>
                      <a-table-column title="备注">
                        <template #cell="{ record }">
                          <a-input v-model="record.comment" placeholder="写入 COLUMN_COMMENT" />
                        </template>
                      </a-table-column>
                      <a-table-column title="PK" :width="60">
                        <template #cell="{ record }">
                          <a-switch v-model="record.primaryKey" />
                        </template>
                      </a-table-column>
                      <a-table-column title="AI" :width="60">
                        <template #cell="{ record }">
                          <a-switch v-model="record.autoIncrement" />
                        </template>
                      </a-table-column>
                      <a-table-column title="操作" :width="80">
                        <template #cell="{ rowIndex }">
                          <a-button status="danger" size="mini" @click="removeFieldRow(rowIndex)">删除</a-button>
                        </template>
                      </a-table-column>
                    </template>
                  </a-table>

                  <a-divider />
                  <a-textarea :model-value="previewSQL" :auto-size="{ minRows: 6, maxRows: 10 }" readonly placeholder="预览SQL会显示在这里" />
                </a-tab-pane>

                <a-tab-pane key="sql" title="手写SQL建表">
                  <a-alert type="warning" style="margin-bottom: 10px">
                    仅允许单条 CREATE TABLE 语句（禁止分号 / 多语句）。生产环境会拒绝执行创建/删除。
                  </a-alert>
                  <a-textarea v-model="createSQL" :auto-size="{ minRows: 6, maxRows: 12 }" placeholder="CREATE TABLE ..." />
                  <div style="margin-top: 10px">
                    <a-button type="primary" :loading="loadingCreate" @click="handleCreateBySQL" v-permission="['create']">
                      执行创建
                    </a-button>
                  </div>
                </a-tab-pane>
                  </a-tabs>
                </a-card>

                <a-card title="结果" :bordered="false">
                  <a-tabs>
                <a-tab-pane key="columns" title="字段">
                  <a-table
                    row-key="COLUMN_NAME"
                    :data="columns"
                    :pagination="false"
                    :loading="loadingColumns"
                    :scroll="{ y: '320px' }"
                    :bordered="{ wrapper: true, cell: true }"
                  >
                    <template #columns>
                      <a-table-column title="字段" data-index="COLUMN_NAME" :width="140" />
                      <a-table-column title="类型" data-index="COLUMN_TYPE" :width="180" />
                      <a-table-column title="可空" data-index="IS_NULLABLE" :width="80" />
                      <a-table-column title="默认" data-index="COLUMN_DEFAULT" :width="120" />
                      <a-table-column title="备注" data-index="COLUMN_COMMENT" />
                      <a-table-column title="EXTRA" data-index="EXTRA" :width="140" />
                      <a-table-column title="操作" :width="120">
                        <template #cell="{ record }">
                          <a-space>
                            <a-button size="mini" @click="openEditColumnModal(record)" v-permission="['update']">编辑</a-button>
                            <a-popconfirm content="确定删除该字段吗？" @ok="handleDropColumn(record)" v-permission="['drop']">
                              <a-button size="mini" status="danger">删除</a-button>
                            </a-popconfirm>
                          </a-space>
                        </template>
                      </a-table-column>
                    </template>
                  </a-table>
                </a-tab-pane>
                <a-tab-pane key="ddl" title="建表SQL">
                  <a-textarea :model-value="createTableSQLText" :auto-size="{ minRows: 14, maxRows: 18 }" readonly />
                </a-tab-pane>
                  </a-tabs>
                </a-card>
              </div>
            </a-scrollbar>
          </a-col>
        </a-row>
      </template>
    </page-card>

    <a-modal
      v-model:visible="columnModalVisible"
      :title="columnModalMode === 'add' ? '新增字段' : '编辑字段'"
      :ok-text="columnModalMode === 'add' ? '新增' : '保存'"
      cancel-text="取消"
      @ok="submitColumnModal"
    >
      <a-form :model="columnModalForm" layout="vertical">
        <a-form-item field="name" label="字段名" required>
          <a-input v-model="columnModalForm.name" placeholder="例如: username" />
        </a-form-item>
        <a-form-item v-if="columnModalMode === 'edit'" field="newName" label="重命名为（可选）">
          <a-input v-model="columnModalForm.newName" placeholder="留空表示不改名" />
        </a-form-item>
        <a-form-item field="type" label="类型" required>
          <a-select v-model="columnModalForm.type" allow-clear>
            <a-option v-for="t in typeOptions" :key="t.value" :value="t.value">{{ t.label }}</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="length" label="长度（可选）">
          <a-input v-model="columnModalForm.length" placeholder="例如: 64 或 10,2" />
        </a-form-item>
        <a-form-item field="nullable" label="可空">
          <a-switch v-model="columnModalForm.nullable" />
        </a-form-item>
        <a-form-item field="unsigned" label="无符号">
          <a-switch v-model="columnModalForm.unsigned" />
        </a-form-item>
        <a-form-item field="default" label="默认值（可选）">
          <a-input v-model="columnModalForm.default" placeholder="留空表示不设置默认值" />
        </a-form-item>
        <a-form-item field="comment" label="备注（可选）">
          <a-input v-model="columnModalForm.comment" />
        </a-form-item>
        <a-form-item field="autoIncrement" label="自增（谨慎）">
          <a-switch v-model="columnModalForm.autoIncrement" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { tabletool } from '/@/components/tabletool';
import {
  getTables,
  getColumns,
  getCreateSQL,
  createTableBySQL,
  dropTable,
  previewCreateSQL,
  createTable,
  addColumn,
  modifyColumn,
  dropColumn,
  type TableItem,
  type CreateTableField,
  type CreateTableReq,
  type AlterColumnReq,
} from './api';

const isFullscreen = ref(false);
const pageHeight = ref(document.documentElement.clientHeight);
const handlePageHeight = (h: number) => {
  pageHeight.value = h;
};
const rightScrollHeight = computed(() => Math.max(pageHeight.value - 60, 320));

const tables = ref([] as TableItem[]);
const columns = ref([] as any[]);
const selectedTable = ref('');
const createSQL = ref('');
const createTableSQLText = ref('');
const previewSQL = ref('');

const loadingTables = ref(false);
const loadingColumns = ref(false);
const loadingCreate = ref(false);
const loadingPreview = ref(false);
const loadingCreateVisual = ref(false);
const columnModalVisible = ref(false);
const columnModalMode = ref('add' as 'add' | 'edit');
const columnModalOldName = ref('');
const columnModalForm = ref({
  name: '',
  newName: '',
  type: 'varchar',
  length: '64',
  unsigned: false,
  nullable: true,
  default: '',
  comment: '',
  autoIncrement: false,
  primaryKey: false,
} as any);

const tableForm = ref({
  tableName: '',
  tableComment: '',
  engine: 'InnoDB',
  charset: 'utf8mb4',
} as { tableName: string; tableComment: string; engine: string; charset: string });

type FieldRow = CreateTableField & { _rowid: string };
const fields = ref([] as FieldRow[]);
const typeOptions = [
  { value: 'bigint', label: 'bigint（大整数）' },
  { value: 'int', label: 'int（整数）' },
  { value: 'smallint', label: 'smallint（小整数）' },
  { value: 'tinyint', label: 'tinyint（微整数）' },
  { value: 'decimal', label: 'decimal（定点数）' },
  { value: 'float', label: 'float（浮点数）' },
  { value: 'double', label: 'double（双精度）' },
  { value: 'varchar', label: 'varchar（可变字符串）' },
  { value: 'char', label: 'char（定长字符串）' },
  { value: 'text', label: 'text（文本）' },
  { value: 'datetime', label: 'datetime（日期时间）' },
  { value: 'date', label: 'date（日期）' },
  { value: 'timestamp', label: 'timestamp（时间戳）' },
  { value: 'json', label: 'json（JSON）' },
  { value: 'blob', label: 'blob（二进制）' },
  { value: 'bool', label: 'bool（布尔，等价 tinyint(1)）' },
];

const getLengthPlaceholder = (t: string) => {
  switch ((t || '').toLowerCase()) {
    case 'varchar':
    case 'char':
      return '必填，例如 64';
    case 'decimal':
      return '必填，例如 10,2';
    case 'bool':
      return '可选，默认 1';
    default:
      return '可选';
  }
};

const addFieldRow = () => {
  fields.value.push({
    _rowid: `${Date.now()}_${Math.random()}`,
    name: '',
    type: 'varchar',
    length: '64',
    nullable: true,
    unsigned: false,
    default: undefined,
    comment: '',
    primaryKey: false,
    autoIncrement: false,
  });
};

const addDefaultIdField = () => {
  // 如果已存在 id 则不重复加
  if (fields.value.some((f) => f.name === 'id')) return;
  fields.value.unshift({
    _rowid: `${Date.now()}_${Math.random()}`,
    name: 'id',
    type: 'bigint',
    length: '20',
    nullable: false,
    unsigned: true,
    default: undefined,
    comment: 'ID',
    primaryKey: true,
    autoIncrement: true,
  });
};

const removeFieldRow = (idx: number) => {
  fields.value.splice(idx, 1);
};

const fetchTables = async () => {
  loadingTables.value = true;
  try {
    tables.value = await getTables({});
  } finally {
    loadingTables.value = false;
  }
};

const handleColumns = async () => {
  if (!selectedTable.value) return Message.error('请先选择表');
  loadingColumns.value = true;
  try {
    columns.value = await getColumns({ tablename: selectedTable.value });
  } finally {
    loadingColumns.value = false;
  }
};

const parseColumnType = (columnType: string) => {
  const raw = (columnType || '').toLowerCase();
  const unsigned = raw.includes('unsigned');
  const typeMatch = raw.match(/^[a-z]+/);
  const type = typeMatch ? typeMatch[0] : 'varchar';
  const lenMatch = raw.match(/\(([^)]+)\)/);
  const length = lenMatch ? lenMatch[1] : '';
  return { type, length, unsigned };
};

const openAddColumnModal = () => {
  if (!selectedTable.value) return Message.error('请先选择表');
  columnModalMode.value = 'add';
  columnModalOldName.value = '';
  columnModalForm.value = {
    name: '',
    newName: '',
    type: 'varchar',
    length: '64',
    unsigned: false,
    nullable: true,
    default: '',
    comment: '',
    autoIncrement: false,
    primaryKey: false,
  };
  columnModalVisible.value = true;
};

const openEditColumnModal = (record: any) => {
  if (!selectedTable.value) return Message.error('请先选择表');
  columnModalMode.value = 'edit';
  const name = record?.COLUMN_NAME || '';
  columnModalOldName.value = name;
  const { type, length, unsigned } = parseColumnType(record?.COLUMN_TYPE || '');
  columnModalForm.value = {
    name,
    newName: '',
    type,
    length,
    unsigned,
    nullable: String(record?.IS_NULLABLE || '').toUpperCase() === 'YES',
    default: record?.COLUMN_DEFAULT ?? '',
    comment: record?.COLUMN_COMMENT ?? '',
    autoIncrement: String(record?.EXTRA || '').toLowerCase().includes('auto_increment'),
    primaryKey: false,
  };
  columnModalVisible.value = true;
};

const submitColumnModal = async () => {
  if (!selectedTable.value) return Message.error('请先选择表');
  const f = columnModalForm.value;
  if (!f.name) return Message.error('请输入字段名');
  const field: CreateTableField = {
    name: f.name,
    type: f.type,
    length: f.length,
    unsigned: !!f.unsigned,
    nullable: !!f.nullable,
    default: f.default === '' ? undefined : f.default,
    comment: f.comment,
    autoIncrement: !!f.autoIncrement,
    primaryKey: !!f.primaryKey,
  };
  if (columnModalMode.value === 'add') {
    const req: AlterColumnReq = {
      tableName: selectedTable.value,
      field,
    };
    await addColumn(req);
    Message.success('新增字段成功');
  } else {
    const req: AlterColumnReq = {
      tableName: selectedTable.value,
      oldName: columnModalOldName.value,
      newName: f.newName || undefined,
      field,
    };
    await modifyColumn(req);
    Message.success('修改字段成功');
  }
  columnModalVisible.value = false;
  await handleColumns();
};

const handleDropColumn = async (record: any) => {
  if (!selectedTable.value) return Message.error('请先选择表');
  const name = record?.COLUMN_NAME;
  if (!name) return;
  await dropColumn({ tableName: selectedTable.value, name });
  Message.success('删除字段成功');
  await handleColumns();
};

const handleCreateSQL = async () => {
  if (!selectedTable.value) return Message.error('请先选择表');
  const res = await getCreateSQL({ tablename: selectedTable.value });
  createTableSQLText.value = JSON.stringify(res, null, 2);
};

const handleCreateBySQL = async () => {
  const sql = createSQL.value.trim();
  if (!sql) return Message.error('请输入SQL');
  loadingCreate.value = true;
  try {
    await createTableBySQL({ sql });
    Message.success('创建成功');
    await fetchTables();
  } finally {
    loadingCreate.value = false;
  }
};

const buildVisualReq = (): CreateTableReq | null => {
  const tableName = tableForm.value.tableName.trim();
  if (!tableName) {
    Message.error('请输入表名');
    return null;
  }
  if (!fields.value.length) {
    Message.error('请至少添加一个字段');
    return null;
  }
  const req: CreateTableReq = {
    tableName,
    tableComment: tableForm.value.tableComment?.trim() || '',
    engine: tableForm.value.engine || 'InnoDB',
    charset: tableForm.value.charset || 'utf8mb4',
    fields: fields.value.map(({ _rowid, ...rest }) => rest),
  };
  return req;
};

const handlePreviewSQL = async () => {
  const req = buildVisualReq();
  if (!req) return;
  loadingPreview.value = true;
  try {
    const res = await previewCreateSQL(req);
    previewSQL.value = res.sql || '';
  } finally {
    loadingPreview.value = false;
  }
};

const handleCreateTable = async () => {
  const req = buildVisualReq();
  if (!req) return;
  loadingCreateVisual.value = true;
  try {
    const res = await createTable(req);
    previewSQL.value = res.sql || previewSQL.value;
    Message.success('创建成功');
    await fetchTables();
  } finally {
    loadingCreateVisual.value = false;
  }
};

const handleDrop = async () => {
  if (!selectedTable.value) return Message.error('请先选择表');
  await dropTable({ tablename: selectedTable.value });
  Message.success('删除成功');
  selectedTable.value = '';
  await fetchTables();
};

fetchTables();
addDefaultIdField();
</script>

<style scoped lang="less">
.dbtablemgr-row {
  height: 100%;
}

.dbtablemgr-left-card {
  height: 100%;
  display: flex;
  flex-direction: column;

  :deep(.arco-card-body) {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  :deep(.arco-table-container) {
    flex: 1;
    height: 100%;
  }
}

.dbtablemgr-right-scroll {
  height: 100%;
}

.dbtablemgr-right-inner {
  padding-bottom: 12px;
}
</style>
