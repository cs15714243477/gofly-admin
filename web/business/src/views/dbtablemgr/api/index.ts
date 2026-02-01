import { defHttp } from '@/utils/http';

enum Api {
  getTables = '/dbtablemgr/getTables',
  getColumns = '/dbtablemgr/getColumns',
  getCreateSQL = '/dbtablemgr/getCreateSQL',
  previewCreateSQL = '/dbtablemgr/previewCreateSQL',
  createTable = '/dbtablemgr/createTable',
  addColumn = '/dbtablemgr/addColumn',
  modifyColumn = '/dbtablemgr/modifyColumn',
  dropColumn = '/dbtablemgr/dropColumn',
  createTableBySQL = '/dbtablemgr/createTableBySQL',
  dropTable = '/dbtablemgr/dropTable',
}

export interface TableItem {
  name: string;
  title: string;
}

export interface CreateTableField {
  name: string;
  type: string;
  length?: string;
  unsigned?: boolean;
  nullable?: boolean;
  default?: any;
  comment?: string;
  primaryKey?: boolean;
  autoIncrement?: boolean;
}

export interface CreateTableReq {
  tableName: string;
  tableComment?: string;
  engine?: string;
  charset?: string;
  collate?: string;
  fields: CreateTableField[];
}

export interface AlterColumnReq {
  tableName: string;
  oldName?: string;
  newName?: string;
  field: CreateTableField;
  position?: 'FIRST' | 'AFTER' | '';
  after?: string;
}

export function getTables(params: object) {
  return defHttp.get({ url: Api.getTables, params }, { errorMessageMode: 'message' }) as unknown as Promise<TableItem[]>;
}

export function getColumns(params: { tablename: string }) {
  return defHttp.get({ url: Api.getColumns, params }, { errorMessageMode: 'message' }) as unknown as Promise<any[]>;
}

export function getCreateSQL(params: { tablename: string }) {
  return defHttp.get({ url: Api.getCreateSQL, params }, { errorMessageMode: 'message' }) as unknown as Promise<any>;
}

export function previewCreateSQL(params: CreateTableReq) {
  return defHttp.post({ url: Api.previewCreateSQL, params }, { errorMessageMode: 'message' }) as unknown as Promise<{ sql: string }>;
}

export function createTable(params: CreateTableReq) {
  return defHttp.post({ url: Api.createTable, params }, { errorMessageMode: 'message' }) as unknown as Promise<{ sql: string }>;
}

export function addColumn(params: AlterColumnReq) {
  return defHttp.post({ url: Api.addColumn, params }, { errorMessageMode: 'message' }) as unknown as Promise<{ sql: string }>;
}

export function modifyColumn(params: AlterColumnReq) {
  return defHttp.post({ url: Api.modifyColumn, params }, { errorMessageMode: 'message' }) as unknown as Promise<{ sql: string }>;
}

export function dropColumn(params: { tableName: string; name: string }) {
  return defHttp.delete({ url: Api.dropColumn, params }, { errorMessageMode: 'message' }) as unknown as Promise<{ sql: string }>;
}

export function createTableBySQL(params: { sql: string }) {
  return defHttp.post({ url: Api.createTableBySQL, params }, { errorMessageMode: 'message' });
}

export function dropTable(params: { tablename: string }) {
  return defHttp.delete({ url: Api.dropTable, params }, { errorMessageMode: 'message' });
}

