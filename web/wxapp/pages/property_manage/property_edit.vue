<template>
  <view class="pe">
    <TopHeader :title="id ? '编辑房源' : '新增房源'">
      <template #left>
        <view class="icon-btn" @tap="goBack">
          <text class="material-symbols-outlined">arrow_back</text>
        </view>
      </template>
    </TopHeader>

    <scroll-view scroll-y="true" class="content">
      <!-- 基础信息 -->
      <view class="card">
        <view class="card-title">基础信息</view>

        <view class="form-row">
          <text class="label required">房源标题</text>
          <input
            v-model="form.title"
            class="input"
            type="text"
            maxlength="100"
            placeholder="请输入房源标题"
            placeholder-class="placeholder"
          />
        </view>

        <view class="grid">
          <view class="form-row">
            <text class="label required">销售状态</text>
            <picker
              mode="selector"
              :range="saleStatusOptions.map((o) => o.label)"
              :value="saleStatusIndex"
              @change="onSaleStatusChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  saleStatusOptions[saleStatusIndex]?.label || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >expand_more</text
                >
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">价格</text>
            <view class="price-row">
              <input
                v-model="form.price"
                class="input"
                type="digit"
                placeholder="0"
                placeholder-class="placeholder"
              />
              <picker
                mode="selector"
                :range="priceUnitOptions.map((o) => o.label)"
                :value="priceUnitIndex"
                @change="onPriceUnitChange"
              >
                <view class="unit-picker">
                  <text class="unit-text">{{
                    priceUnitOptions[priceUnitIndex]?.label ||
                    form.price_unit ||
                    "万"
                  }}</text>
                  <text class="material-symbols-outlined unit-ic"
                    >expand_more</text
                  >
                </view>
              </picker>
            </view>
          </view>

          <view class="form-row">
            <text class="label">面积(㎡)</text>
            <input
              v-model="form.area"
              class="input"
              type="digit"
              placeholder="0"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">户型</text>
            <view class="room-row">
              <input
                v-model="form.rooms"
                class="mini"
                type="number"
                placeholder="室"
                placeholder-class="placeholder"
              />
              <text class="mini-suffix">室</text>
              <input
                v-model="form.halls"
                class="mini"
                type="number"
                placeholder="厅"
                placeholder-class="placeholder"
              />
              <text class="mini-suffix">厅</text>
              <input
                v-model="form.bathrooms"
                class="mini"
                type="number"
                placeholder="卫"
                placeholder-class="placeholder"
              />
              <text class="mini-suffix">卫</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 房源参数 -->
      <view class="card">
        <view class="card-title">房源参数</view>

        <view class="grid">
          <view class="form-row">
            <text class="label">建成年份</text>
            <picker
              mode="selector"
              :range="buildYearOptions.map((o) => o.label)"
              :value="buildYearIndex"
              @change="onBuildYearChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  buildYearOptions[buildYearIndex]?.label || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >expand_more</text
                >
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">楼层位置</text>
            <picker
              mode="selector"
              :range="floorLevelOptions.map((o) => o.label)"
              :value="floorLevelIndex"
              @change="onFloorLevelChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  floorLevelOptions[floorLevelIndex]?.label || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >expand_more</text
                >
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">总楼层</text>
            <input
              v-model="form.total_floors"
              class="input"
              type="number"
              placeholder="共多少层"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">朝向</text>
            <picker
              mode="selector"
              :range="orientationOptions.map((o) => o.label)"
              :value="orientationIndex"
              @change="onOrientationChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  orientationOptions[orientationIndex]?.label || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >expand_more</text
                >
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">物业类型</text>
            <picker
              mode="selector"
              :range="propertyTypeOptions.map((o) => o.label)"
              :value="propertyTypeIndex"
              @change="onPropertyTypeChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  propertyTypeOptions[propertyTypeIndex]?.label || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >expand_more</text
                >
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">装修标准</text>
            <picker
              mode="selector"
              :range="decorationTypeOptions.map((o) => o.label)"
              :value="decorationTypeIndex"
              @change="onDecorationTypeChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  decorationTypeOptions[decorationTypeIndex]?.label || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >expand_more</text
                >
              </view>
            </picker>
          </view>
        </view>
      </view>

      <!-- 房主与收房 -->
      <view class="card">
        <view class="card-title">房主与收房</view>
        <view class="grid">
          <view class="form-row">
            <text class="label">房主姓名</text>
            <input
              v-model="form.owner_name"
              class="input"
              type="text"
              placeholder="可选"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">房主电话</text>
            <input
              v-model="form.owner_phone"
              class="input"
              type="number"
              placeholder="可选"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">收房人姓名</text>
            <input
              v-model="form.receiver_name"
              class="input"
              type="text"
              placeholder="可选"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">收房人电话</text>
            <input
              v-model="form.receiver_phone"
              class="input"
              type="number"
              placeholder="可选"
              placeholder-class="placeholder"
            />
          </view>
        </view>
        <view class="form-row">
          <text class="label">收房价格(支付业主)</text>
          <input
            v-model="form.receiver_price"
            class="input"
            type="digit"
            placeholder="未填写"
            placeholder-class="placeholder"
          />
        </view>
      </view>

      <!-- 图片 -->
      <view class="card">
        <view class="card-title">图片</view>
        <view class="hint">建议上传清晰横图，点击图片可设置封面/删除</view>

        <view class="media-grid">
          <view
            v-for="(img, idx) in images"
            :key="'img-' + idx"
            class="media-item"
            @tap="openImageActions(idx)"
          >
            <image
              class="media-img"
              :src="toFullMediaUrl(img)"
              mode="aspectFill"
            ></image>
            <view v-if="img === form.cover_image" class="badge">封面</view>
          </view>
          <view class="media-item add" @tap="pickImages">
            <text class="material-symbols-outlined add-ic"
              >add_photo_alternate</text
            >
            <text class="add-text">添加</text>
          </view>
        </view>

        <view class="switch-row">
          <view class="switch-left">
            <text class="material-symbols-outlined sw-ic">download</text>
            <text class="sw-text">允许下载图片</text>
          </view>
          <switch
            :checked="Number(form.allow_image_download) === 1"
            @change="
              (e) => (form.allow_image_download = e.detail.value ? 1 : 0)
            "
            color="#2563eb"
          />
        </view>
      </view>

      <!-- 视频 -->
      <view class="card">
        <view class="card-title">视频</view>
        <view class="hint">支持单个房源视频，建议 10-30 秒</view>

        <view class="video-box" @tap="openVideoActions">
          <view v-if="!form.video_url" class="video-empty">
            <text class="material-symbols-outlined video-ic">movie</text>
            <view class="video-texts">
              <text class="video-title">未上传视频</text>
              <text class="video-sub">点击上传/更换</text>
            </view>
          </view>
          <video
            v-else
            class="video"
            :src="toFullMediaUrl(form.video_url)"
            controls
            object-fit="cover"
          ></video>
        </view>

        <view class="switch-row">
          <view class="switch-left">
            <text class="material-symbols-outlined sw-ic">download</text>
            <text class="sw-text">允许下载视频</text>
          </view>
          <switch
            :checked="Number(form.allow_video_download) === 1"
            @change="
              (e) => (form.allow_video_download = e.detail.value ? 1 : 0)
            "
            color="#2563eb"
          />
        </view>
      </view>

      <!-- 标签 -->
      <view class="card">
        <view class="card-title">标签</view>
        <view class="tag-row">
          <view class="tag" v-for="(t, idx) in tags" :key="'t-' + idx">
            <text class="tag-text">{{ t }}</text>
            <text
              class="material-symbols-outlined tag-x"
              @tap.stop="removeTag(idx)"
              >close</text
            >
          </view>
          <view v-if="tags.length === 0" class="hint small">暂无标签</view>
        </view>
        <view class="tag-add">
          <input
            v-model="newTag"
            class="input"
            type="text"
            maxlength="12"
            placeholder="输入标签（最多6个）"
            placeholder-class="placeholder"
          />
          <button class="btn add" @tap="addTag">添加</button>
        </view>
      </view>

      <!-- 位置 -->
      <view class="card">
        <view class="card-title">位置</view>
        <view class="form-row">
          <text class="label">小区名称</text>
          <input
            v-model="form.community_name"
            class="input"
            type="text"
            placeholder="请输入小区/楼盘"
            placeholder-class="placeholder"
          />
        </view>

        <view class="form-row">
          <text class="label">省市区</text>
          <picker
            mode="region"
            :value="addressRegion"
            @change="onAddressRegionChange"
          >
            <view class="picker">
              <text class="picker-text">{{ addressRegionText }}</text>
              <text class="material-symbols-outlined picker-ic"
                >expand_more</text
              >
            </view>
          </picker>
        </view>

        <view class="form-row">
          <text class="label">详细地址</text>
          <input
            v-model="addressDetail"
            class="input"
            type="text"
            placeholder="街道、门牌号、楼栋单元等"
            placeholder-class="placeholder"
          />
        </view>

        <view class="coord-row">
          <view class="coord-left">
            <text class="material-symbols-outlined coord-ic">my_location</text>
            <text class="coord-text">{{
              form.latitude && form.longitude
                ? Number(form.latitude).toFixed(6) +
                  ", " +
                  Number(form.longitude).toFixed(6)
                : "未选择坐标"
            }}</text>
          </view>
          <button class="btn ghost" @tap="pickLocation">地图选点</button>
        </view>
      </view>

      <!-- 发布与佣金 -->
      <view class="card">
        <view class="card-title">发布与佣金</view>
        <view class="grid">
          <view class="form-row">
            <text class="label">佣金比例(%)</text>
            <input
              v-model="form.commission_rate"
              class="input"
              type="digit"
              placeholder="0-100"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">成交奖励金</text>
            <input
              v-model="form.commission_reward"
              class="input"
              type="digit"
              placeholder="0"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">排序权重</text>
            <input
              v-model="form.weigh"
              class="input"
              type="number"
              placeholder="0"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">智能门锁</text>
            <view class="picker disabled">
              <text class="picker-text">{{
                Number(form.has_smart_lock) === 1 ? "已绑定" : "未绑定"
              }}</text>
              <text class="material-symbols-outlined picker-ic">lock</text>
            </view>
            <view class="hint small tip">门锁请在后台/详情页绑定</view>
          </view>
        </view>

        <view class="switch-row">
          <view class="switch-left">
            <text class="material-symbols-outlined sw-ic">stars</text>
            <text class="sw-text">推荐房源</text>
          </view>
          <switch
            :checked="Number(form.hot_status) === 1"
            @change="(e) => (form.hot_status = e.detail.value ? 1 : 0)"
            color="#2563eb"
          />
        </view>

        <view class="switch-row">
          <view class="switch-left">
            <text class="material-symbols-outlined sw-ic">visibility</text>
            <text class="sw-text">房源启用</text>
          </view>
          <switch
            :checked="Number(form.status) === 0"
            @change="(e) => (form.status = e.detail.value ? 0 : 1)"
            color="#2563eb"
          />
        </view>
      </view>

      <!-- 装修进度（放在最后） -->
      <view class="card">
        <view class="card-title">装修进度</view>
        <view class="hint">装修信息无需单独保存，点击底部“保存”会一起提交</view>

        <view class="form-row">
          <text class="label">装修状态</text>
          <picker
            mode="selector"
            :range="renovationStatusOptions.map((o) => o.label)"
            :value="renovationStatusIndex"
            @change="onRenovationStatusChange"
          >
            <view class="picker">
              <text class="picker-text">{{
                renovationStatusOptions[renovationStatusIndex]?.label ||
                "请选择"
              }}</text>
              <text class="material-symbols-outlined picker-ic"
                >expand_more</text
              >
            </view>
          </picker>
          <view v-if="!id" class="hint small tip"
            >请先保存房源生成ID后再维护装修进度</view
          >
        </view>

        <view v-if="id && renovation.renovation_status !== 'none'" class="grid">
          <view class="form-row">
            <text class="label">进度(%)</text>
            <input
              v-model="renovation.progress_percentage"
              class="input"
              type="number"
              placeholder="0-100"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">当前阶段</text>
            <picker
              mode="selector"
              :range="renovationStageOptions.map((o) => o.label)"
              :value="renovationStageIndex"
              @change="onRenovationStageChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  renovationStageOptions[renovationStageIndex]?.label ||
                  "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >expand_more</text
                >
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">开始日期</text>
            <picker
              mode="date"
              :value="renovation.start_date"
              @change="(e) => (renovation.start_date = e.detail.value)"
            >
              <view class="picker">
                <text class="picker-text">{{
                  renovation.start_date || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >calendar_month</text
                >
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">预计完工</text>
            <picker
              mode="date"
              :value="renovation.estimated_finish_date"
              @change="
                (e) => (renovation.estimated_finish_date = e.detail.value)
              "
            >
              <view class="picker">
                <text class="picker-text">{{
                  renovation.estimated_finish_date || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic">event</text>
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">实际完工</text>
            <picker
              mode="date"
              :value="renovation.actual_finish_date"
              @change="(e) => (renovation.actual_finish_date = e.detail.value)"
            >
              <view class="picker">
                <text class="picker-text">{{
                  renovation.actual_finish_date || "请选择"
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >verified</text
                >
              </view>
            </picker>
          </view>
        </view>

        <view v-if="id && renovation.renovation_status !== 'none'">
          <view class="form-row">
            <text class="label">材料</text>
            <view class="tag-row">
              <view
                class="tag"
                v-for="(t, idx) in renovationMaterials"
                :key="'m-' + idx"
              >
                <text class="tag-text">{{ t }}</text>
                <text
                  class="material-symbols-outlined tag-x"
                  @tap.stop="removeMaterial(idx)"
                  >close</text
                >
              </view>
              <view v-if="renovationMaterials.length === 0" class="hint small"
                >暂无材料</view
              >
            </view>
            <view class="tag-add">
              <input
                v-model="newMaterial"
                class="input"
                type="text"
                maxlength="12"
                placeholder="输入材料（最多10个）"
                placeholder-class="placeholder"
              />
              <button class="btn add" @tap="addMaterial">添加</button>
            </view>
          </view>

          <view class="form-row">
            <text class="label">装修图片</text>
            <view class="media-grid">
              <view
                v-for="(img, idx) in renovationImages"
                :key="'rimg-' + idx"
                class="media-item"
                @tap="openRenovationImageActions(idx)"
              >
                <image
                  class="media-img"
                  :src="toFullMediaUrl(img)"
                  mode="aspectFill"
                ></image>
              </view>
              <view class="media-item add" @tap="pickRenovationImages">
                <text class="material-symbols-outlined add-ic"
                  >add_photo_alternate</text
                >
                <text class="add-text">添加</text>
              </view>
            </view>
          </view>

          <!-- 工序时间线（按工序分组图片） -->
          <view class="stage-logs">
            <view class="stage-logs-head">
              <text class="label">工序时间线</text>
              <view class="stage-logs-actions">
                <button
                  class="mini-btn"
                  :disabled="savingRenovation || uploading"
                  @tap="confirmGenerateStageLogs"
                >
                  生成默认工序
                </button>
                <button
                  class="mini-btn ghost"
                  :disabled="savingRenovation || uploading"
                  @tap="syncStageLogsByCurrentStage"
                >
                  按当前阶段同步
                </button>
                <button
                  class="mini-btn ghost"
                  :disabled="savingRenovation || uploading"
                  @tap="addStageLog"
                >
                  新增工序
                </button>
                <button
                  class="mini-btn danger"
                  :disabled="savingRenovation || uploading"
                  @tap="clearStageLogs"
                >
                  清空
                </button>
              </view>
            </view>

            <view
              v-if="!(renovation.stage_logs || []).length"
              class="hint small"
            >
              暂无工序记录，可点击“生成默认工序”快速创建；每个工序可上传一组图片（如水电/泥瓦）。
            </view>

            <view v-else class="stage-log-list">
              <view
                class="stage-log-item"
                v-for="(log, idx) in renovation.stage_logs"
                :key="'sl-' + idx"
              >
                <view class="stage-log-top">
                  <view class="stage-log-title">
                    <text class="name">{{ log.stage || "未选择工序" }}</text>
                    <text class="tag" :class="log.status || 'todo'">{{
                      getStageStatusText(log.status)
                    }}</text>
                  </view>
                  <view class="stage-log-tools">
                    <button
                      class="mini-icon-btn"
                      :disabled="idx === 0"
                      @tap.stop="moveStageLog(idx, -1)"
                    >
                      上移
                    </button>
                    <button
                      class="mini-icon-btn"
                      :disabled="idx === (renovation.stage_logs || []).length - 1"
                      @tap.stop="moveStageLog(idx, 1)"
                    >
                      下移
                    </button>
                    <button
                      class="mini-icon-btn danger"
                      @tap.stop="removeStageLog(idx)"
                    >
                      删除
                    </button>
                  </view>
                </view>

                <view class="grid stage-log-grid">
                  <view class="form-row">
                    <text class="label">工序</text>
                    <picker
                      mode="selector"
                      :range="renovationStageOptions.map((o) => o.label)"
                      :value="getStageLogStageIndex(log)"
                      @change="(e) => onStageLogStageChange(idx, e)"
                    >
                      <view class="picker">
                        <text class="picker-text">{{
                          log.stage || "请选择"
                        }}</text>
                        <text class="material-symbols-outlined picker-ic"
                          >expand_more</text
                        >
                      </view>
                    </picker>
                  </view>

                  <view class="form-row">
                    <text class="label">状态</text>
                    <picker
                      mode="selector"
                      :range="stageStatusOptions.map((o) => o.label)"
                      :value="getStageLogStatusIndex(log)"
                      @change="(e) => onStageLogStatusChange(idx, e)"
                    >
                      <view class="picker">
                        <text class="picker-text">{{
                          stageStatusOptions[getStageLogStatusIndex(log)]
                            ?.label || "请选择"
                        }}</text>
                        <text class="material-symbols-outlined picker-ic"
                          >expand_more</text
                        >
                      </view>
                    </picker>
                  </view>

                  <view class="form-row">
                    <text class="label">日期</text>
                    <picker
                      mode="date"
                      :value="log.date"
                      @change="(e) => onStageLogDateChange(idx, e)"
                    >
                      <view class="picker">
                        <text class="picker-text">{{
                          log.date || "请选择"
                        }}</text>
                        <text class="material-symbols-outlined picker-ic"
                          >calendar_month</text
                        >
                      </view>
                    </picker>
                  </view>

                  <view class="form-row stage-log-note">
                    <text class="label">备注</text>
                    <textarea
                      v-model="log.note"
                      class="textarea small"
                      maxlength="200"
                      placeholder="可选：该工序说明"
                      placeholder-class="placeholder"
                    ></textarea>
                  </view>
                </view>

                <view class="form-row">
                  <text class="label">工序图片</text>
                  <view class="media-grid stage-media-grid">
                    <view
                      v-for="(img, iidx) in ensureStageLogImages(log)"
                      :key="'slimg-' + idx + '-' + iidx"
                      class="media-item stage-media-item"
                      @tap="openStageLogImageActions(idx, iidx)"
                    >
                      <image
                        class="media-img"
                        :src="toFullMediaUrl(img)"
                        mode="aspectFill"
                      ></image>
                    </view>
                    <view
                      class="media-item add stage-media-add"
                      @tap="pickStageLogImages(idx)"
                    >
                      <text class="material-symbols-outlined add-ic"
                        >add_photo_alternate</text
                      >
                      <text class="add-text">添加</text>
                    </view>
                  </view>
                </view>
              </view>
            </view>
          </view>

          <view class="form-row">
            <text class="label">施工说明</text>
            <textarea
              v-model="renovation.notes"
              class="textarea"
              maxlength="500"
              placeholder="可选：记录进度说明、施工要点等"
              placeholder-class="placeholder"
            ></textarea>
          </view>

          <view class="renovation-actions">
            <button
              class="btn ghost"
              :disabled="savingRenovation || uploading"
              @tap="loadRenovation"
            >
              刷新
            </button>
          </view>
        </view>
      </view>

      <view class="bottom-spacer"></view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="footer">
      <button
        class="footer-btn ghost"
        :disabled="!id || saving || savingRenovation || uploading"
        @tap="preview"
      >
        预览
      </button>
      <button
        class="footer-btn primary"
        :disabled="saving || savingRenovation || uploading"
        @tap="save"
      >
        {{ saving || savingRenovation ? "保存中..." : "保存" }}
      </button>
    </view>
  </view>
</template>

<script>
import TopHeader from "@/components/TopHeader.vue";
import propertyApi from "@/api/property";
import $store from "@/store";
import md5 from "js-md5";
import { baseUrl } from "@/utils/config";

// 默认省市区（新增房源、或地址无法解析时使用；带省/市/区后缀）
const DEFAULT_ADDRESS_REGION = ["辽宁省", "沈阳市", "沈河区"];

// base64 编码（兼容小程序端无 window.btoa）
function base64Encode(str = "") {
  try {
    // #ifdef H5
    return window.btoa(str);
    // #endif
  } catch (e) {}

  // #ifndef H5
  const utf8ToBytes = (s) => {
    const bytes = [];
    for (let i = 0; i < s.length; i++) {
      let codePoint = s.charCodeAt(i);
      if (codePoint < 0x80) {
        bytes.push(codePoint);
      } else if (codePoint < 0x800) {
        bytes.push(0xc0 | (codePoint >> 6));
        bytes.push(0x80 | (codePoint & 0x3f));
      } else {
        bytes.push(0xe0 | (codePoint >> 12));
        bytes.push(0x80 | ((codePoint >> 6) & 0x3f));
        bytes.push(0x80 | (codePoint & 0x3f));
      }
    }
    return new Uint8Array(bytes);
  };
  const buf = utf8ToBytes(str).buffer;
  // eslint-disable-next-line no-undef
  if (typeof wx !== "undefined" && wx.arrayBufferToBase64)
    return wx.arrayBufferToBase64(buf);
  if (typeof uni !== "undefined" && uni.arrayBufferToBase64)
    return uni.arrayBufferToBase64(buf);
  return "";
  // #endif
}

export default {
  components: { TopHeader },
  data() {
    const curYear = new Date().getFullYear();
    const buildYearOptions = [{ label: "请选择", value: "" }].concat(
      Array.from({ length: curYear - 1950 + 1 }).map((_, i) => {
        const y = String(curYear - i);
        return { label: y, value: y };
      })
    );

    return {
      id: 0,
      saving: false,
      uploading: false,
      form: {
        title: "",
        price: "",
        price_unit: "万",
        area: "",
        rooms: 0,
        halls: 0,
        bathrooms: 0,
        build_year: "",
        floor_level: "",
        total_floors: "",
        orientation: "",
        property_type: "",
        decoration_type: "",
        community_name: "",
        address: "",
        latitude: "",
        longitude: "",
        tags: [],
        images: [],
        cover_image: "",
        video_url: "",
        allow_image_download: 1,
        allow_video_download: 1,
        sale_status: "on_sale",
        owner_name: "",
        owner_phone: "",
        receiver_name: "",
        receiver_phone: "",
        receiver_price: "",
        commission_rate: "",
        commission_reward: "",
        weigh: "",
        hot_status: 0,
        has_smart_lock: 0,
        status: 0,
      },
      tags: [],
      images: [],
      newTag: "",
      // 地址分段输入：省市区 + 详细地址（最终保存时合成到 form.address）
      addressRegion: [...DEFAULT_ADDRESS_REGION],
      addressDetail: "",
      saleStatusOptions: [
        { label: "在售", value: "on_sale" },
        { label: "已售", value: "sold" },
        { label: "下架", value: "off_market" },
      ],
      buildYearOptions,
      floorLevelOptions: [
        { label: "低层 (1-6)", value: "低层" },
        { label: "中层 (7-15)", value: "中层" },
        { label: "高层 (16+)", value: "高层" },
        { label: "地下", value: "地下" },
      ],
      orientationOptions: [
        "东",
        "南",
        "西",
        "北",
        "东南",
        "东北",
        "西南",
        "西北",
        "南北",
        "东西",
      ].map((it) => ({ label: it, value: it })),
      propertyTypeOptions: ["住宅", "公寓", "别墅", "商铺", "写字楼"].map(
        (it) => ({
          label: it,
          value: it,
        })
      ),
      decorationTypeOptions: ["毛坯", "简装", "精装", "豪装"].map((it) => ({
        label: it,
        value: it,
      })),
      priceUnitOptions: [
        { label: "万", value: "万" },
        { label: "元", value: "元" },
      ],
      // 装修进度（可维护房源）
      savingRenovation: false,
      renovation: {
        renovation_status: "none",
        progress_percentage: 0,
        current_stage: "",
        start_date: "",
        estimated_finish_date: "",
        actual_finish_date: "",
        stage_logs: [],
        notes: "",
        status: 0,
      },
      stageStatusOptions: [
        { label: "未开始", value: "todo" },
        { label: "进行中", value: "doing" },
        { label: "已完成", value: "done" },
      ],
      renovationStatusOptions: [
        { label: "未装修", value: "none" },
        { label: "装修中", value: "in_progress" },
        { label: "已完成", value: "done" },
      ],
      renovationStageOptions: [{ label: "请选择", value: "" }].concat(
        [
          "设计",
          "拆改",
          "水电",
          "泥瓦",
          "木工",
          "油漆",
          "安装",
          "软装",
          "验收",
        ].map((it) => ({ label: it, value: it }))
      ),
      renovationMaterials: [],
      renovationImages: [],
      newMaterial: "",
    };
  },
  computed: {
    saleStatusIndex() {
      const v = String(this.form.sale_status || "on_sale");
      const idx = this.saleStatusOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
    priceUnitIndex() {
      const v = String(this.form.price_unit || "万");
      const idx = this.priceUnitOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
    buildYearIndex() {
      const v = String(this.form.build_year || "");
      const idx = this.buildYearOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
    floorLevelIndex() {
      const v = String(this.form.floor_level || "");
      const idx = this.floorLevelOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
    orientationIndex() {
      const v = String(this.form.orientation || "");
      const idx = this.orientationOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
    propertyTypeIndex() {
      const v = String(this.form.property_type || "");
      const idx = this.propertyTypeOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
    decorationTypeIndex() {
      const v = String(this.form.decoration_type || "");
      const idx = this.decorationTypeOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
    addressRegionText() {
      const arr = Array.isArray(this.addressRegion) ? this.addressRegion : [];
      const text = arr.filter((it) => String(it || "").trim()).join(" / ");
      return text || "请选择省市区";
    },
    renovationStatusIndex() {
      const v = String(this.renovation.renovation_status || "none");
      const idx = this.renovationStatusOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
    renovationStageIndex() {
      const v = String(this.renovation.current_stage || "");
      const idx = this.renovationStageOptions.findIndex((o) => o.value === v);
      return idx >= 0 ? idx : 0;
    },
  },
  onLoad(options) {
    const id = Number((options && (options.id || options.ID)) || 0) || 0;
    this.id = id;
    this.bootstrap();
  },
  methods: {
    // 兼容：上传接口可能返回相对路径（/resource/uploads/...），这里统一转为可展示的完整地址
    toFullMediaUrl(url) {
      const s = String(url || "").trim();
      if (!s) return "";
      if (
        /^https?:\/\//i.test(s) ||
        s.startsWith("wxfile://") ||
        s.startsWith("file://")
      )
        return s;
      const root = String(baseUrl || "").replace(/\/+$/, "");
      const path = s.startsWith("/") ? s : `/${s}`;
      return root ? `${root}${path}` : path;
    },
    goBack() {
      uni.navigateBack();
    },
    async ensurePermission() {
      const userStore = $store("user");
      const token = uni.getStorageSync("token");
      if (!token && !userStore.isLogin) {
        uni.reLaunch({ url: "/pages/login/login" });
        return false;
      }
      try {
        const ui = userStore.userInfo || {};
        // 兼容：老版本缓存的 userInfo 可能没有 can_manage_properties，需要强制刷新一次
        if (!ui.id || typeof ui.can_manage_properties === "undefined") {
          await userStore.getInfo();
        }
      } catch (e) {}
      const canManage =
        Number((userStore.userInfo || {}).can_manage_properties) === 1;
      if (!canManage) {
        uni.showModal({
          title: "提示",
          content: "暂无房源维护权限，请联系后台管理员开启。",
          showCancel: false,
          success: () => {
            uni.navigateBack();
          },
        });
        return false;
      }
      return true;
    },
    onSaleStatusChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt = this.saleStatusOptions[idx] || this.saleStatusOptions[0];
      this.form.sale_status = opt.value;
    },
    onPriceUnitChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt = this.priceUnitOptions[idx] || this.priceUnitOptions[0];
      this.form.price_unit = (opt && opt.value) || "万";
    },
    // 页面初始化：权限校验 + 下拉选项缓存 + 回显
    async bootstrap() {
      const ok = await this.ensurePermission();
      if (!ok) return;
      await this.loadFormOptions();
      if (this.id) {
        await this.loadContent();
        await this.loadRenovation();
      }
    },
    // 首次加载表单下拉项时写入缓存；有缓存直接使用
    async loadFormOptions() {
      const userStore = $store("user");
      const businessId = Number((userStore.userInfo || {}).business_id) || 1;
      const cacheKey = `wx_property_form_options_v2_${businessId}`;
      const cached = uni.getStorageSync(cacheKey);
      if (cached && typeof cached === "object") {
        this.applyFormOptions(cached);
        return;
      }

      const res = await propertyApi.getFormOptions({});
      if (!res || res.code !== 0) return;
      const data = res.data || {};
      this.applyFormOptions(data);
      try {
        uni.setStorageSync(cacheKey, data);
      } catch (e) {}
    },
    applyFormOptions(data) {
      const normalize = (arr) =>
        (Array.isArray(arr) ? arr : [])
          .map((it) => ({
            label: String((it && it.label) || "").trim(),
            value: String((it && it.value) || "").trim(),
          }))
          .filter((it) => it.label && it.value);

      const sale = normalize(data.sale_status);
      if (sale.length) this.saleStatusOptions = sale;

      const unit = normalize(data.price_unit);
      if (unit.length) this.priceUnitOptions = unit;

      const floor = normalize(data.floor_level);
      if (floor.length) this.floorLevelOptions = floor;

      const ori = normalize(data.orientation);
      if (ori.length) this.orientationOptions = ori;

      const pt = normalize(data.property_type);
      if (pt.length) this.propertyTypeOptions = pt;

      const deco = normalize(data.decoration_type);
      if (deco.length) this.decorationTypeOptions = deco;

      const stages = normalize(data.renovation_stage);
      if (stages.length)
        this.renovationStageOptions = [{ label: "请选择", value: "" }].concat(
          stages
        );

      // 兜底：当前值不在选项内时，回退到第一个选项
      const fixValue = (key, options, fallback) => {
        const v = String(this.form[key] || "").trim();
        const hit = (options || []).some((o) => String(o.value) === v);
        if (!hit)
          this.form[key] =
            (options && options[0] && options[0].value) || fallback;
      };
      fixValue("sale_status", this.saleStatusOptions, "on_sale");
      fixValue("price_unit", this.priceUnitOptions, "万");
      fixValue("floor_level", this.floorLevelOptions, "");
      fixValue("orientation", this.orientationOptions, "");
      fixValue("property_type", this.propertyTypeOptions, "");
      fixValue("decoration_type", this.decorationTypeOptions, "");

      // 装修工序：允许历史数据不在选项内时追加显示，避免 picker 显示异常
      const stageVal = String(this.renovation.current_stage || "").trim();
      if (
        stageVal &&
        !(this.renovationStageOptions || []).some(
          (o) => String(o.value) === stageVal
        )
      ) {
        this.renovationStageOptions = (
          this.renovationStageOptions || []
        ).concat([{ label: stageVal, value: stageVal }]);
      }

      // 兜底：工序时间线里的工序不在下拉里时，追加显示，避免 picker 显示异常
      const logStages = Array.isArray(this.renovation.stage_logs)
        ? this.renovation.stage_logs
            .map((it) => String((it && it.stage) || "").trim())
            .filter(Boolean)
        : [];
      if (logStages.length) {
        const exist = new Set(
          (this.renovationStageOptions || []).map((o) =>
            String((o && o.value) || "").trim()
          )
        );
        const missing = [];
        logStages.forEach((s) => {
          if (!exist.has(s)) {
            exist.add(s);
            missing.push({ label: s, value: s });
          }
        });
        if (missing.length) {
          this.renovationStageOptions = (this.renovationStageOptions || []).concat(
            missing
          );
        }
      }
    },
    onBuildYearChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt = this.buildYearOptions[idx] || this.buildYearOptions[0];
      this.form.build_year = opt.value;
    },
    onFloorLevelChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt = this.floorLevelOptions[idx] || this.floorLevelOptions[0];
      this.form.floor_level = opt.value;
    },
    onOrientationChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt = this.orientationOptions[idx] || this.orientationOptions[0];
      this.form.orientation = opt.value;
    },
    onPropertyTypeChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt = this.propertyTypeOptions[idx] || this.propertyTypeOptions[0];
      this.form.property_type = opt.value;
    },
    onDecorationTypeChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt =
        this.decorationTypeOptions[idx] || this.decorationTypeOptions[0];
      this.form.decoration_type = opt.value;
    },
    onAddressRegionChange(e) {
      const v = (e && e.detail && e.detail.value) || [];
      if (Array.isArray(v)) this.addressRegion = v;
    },
    buildFullAddress() {
      const region = Array.isArray(this.addressRegion)
        ? this.addressRegion
            .map((it) => String(it || "").trim())
            .filter(Boolean)
        : [];
      let detail = String(this.addressDetail || "").trim();

      // 防止重复：当“详细地址”里本身包含省市区（比如地图选点返回了完整地址/用户粘贴了完整地址），
      // 保存时再拼接一次会出现“省市区省市区省市区…”的累积。
      if (region.length >= 3 && detail) {
        const esc = (s) => String(s || "").replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
        const p0 = esc(region[0]);
        const p1 = esc(region[1]);
        const p2 = esc(region[2]);
        // 支持空格或“/”作为分隔符；移除开头连续重复的“省 市 区”片段
        const re = new RegExp(
          `^(?:${p0}\\s*(?:/|\\s)+${p1}\\s*(?:/|\\s)+${p2}\\s*)+`,
        );
        let next = detail.replace(re, "").trim();
        // 兼容：某些来源地址可能没有任何分隔符（如“辽宁省沈阳市沈河区…”）
        if (next === detail) {
          const re2 = new RegExp(`^(?:${p0}${p1}${p2}\\s*)+`);
          next = detail.replace(re2, "").trim();
        }
        if (next !== detail) {
          detail = next;
          // 同步回写，避免用户看到的输入框内容也越来越长
          this.addressDetail = detail;
        }
      }

      return [...region, detail].filter(Boolean).join(" ");
    },
    initAddressUIFromAddress(address) {
      const raw = String(address || "").trim();
      // 兜底：保证“省市区”始终有默认值（避免进入页面为空）
      this.addressRegion = [...DEFAULT_ADDRESS_REGION];
      this.addressDetail = "";
      if (!raw) return;
      const parts = raw.split(/\\s+/).filter(Boolean);
      if (parts.length >= 3) {
        let p0 = parts[0] || "";
        const p1 = parts[1] || "";
        const p2 = parts[2] || "";

        // 兼容：直辖市历史数据可能是“北京 北京市 海淀区 ...”，region 期望“北京市”
        if (
          ["北京", "天津", "上海", "重庆"].includes(p0) &&
          p1 &&
          p1.endsWith("市") &&
          p1.includes(p0)
        ) {
          p0 = p1;
        }

        this.addressRegion = [p0, p1, p2];
        this.addressDetail = parts.slice(3).join(" ");
        return;
      }
      // 只有详细地址/或字段内容不满足 3 段：保留默认省市区，把 raw 作为“详细地址”
      this.addressDetail = raw;
    },
    onRenovationStageChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt =
        this.renovationStageOptions[idx] || this.renovationStageOptions[0];
      this.renovation.current_stage = (opt && opt.value) || "";
    },
    async loadContent() {
      const res = await propertyApi.getManageContent({ id: this.id });
      if (!res || res.code !== 0) return;
      const data = res.data || {};
      this.form = {
        ...this.form,
        ...data,
      };
      this.tags = Array.isArray(data.tags) ? data.tags : [];
      this.images = Array.isArray(data.images) ? data.images : [];
      if (this.form.build_year)
        this.form.build_year = String(this.form.build_year);
      if (this.form.total_floors)
        this.form.total_floors = String(this.form.total_floors);
      if (this.form.weigh || this.form.weigh === 0)
        this.form.weigh = String(this.form.weigh);
      if (
        this.form.receiver_price === null ||
        typeof this.form.receiver_price === "undefined"
      ) {
        this.form.receiver_price = "";
      }
      if (
        this.form.commission_rate === null ||
        typeof this.form.commission_rate === "undefined"
      ) {
        this.form.commission_rate = "";
      }
      if (
        this.form.commission_reward === null ||
        typeof this.form.commission_reward === "undefined"
      ) {
        this.form.commission_reward = "";
      }
      // 地址拆分：用于省市区选择 + 详细地址输入
      this.initAddressUIFromAddress(this.form.address);
      // 兜底：封面未设置时，用第一张图
      if (!this.form.cover_image && this.images.length > 0) {
        this.form.cover_image = this.images[0];
      }
    },
    async loadRenovation() {
      if (!this.id) return;
      const res = await propertyApi.getManageRenovation({
        property_id: this.id,
      });
      if (!res || res.code !== 0) return;
      const data = res.data || {};
      this.renovation = {
        ...this.renovation,
        ...data,
      };
      this.renovationMaterials = Array.isArray(data.materials)
        ? data.materials
        : [];
      this.renovationImages = Array.isArray(data.images) ? data.images : [];
      this.renovation.stage_logs = this.normalizeStageLogs(
        data.stage_logs || this.renovation.stage_logs
      );

      // 兜底类型修正
      if (
        this.renovation.progress_percentage === null ||
        typeof this.renovation.progress_percentage === "undefined"
      ) {
        this.renovation.progress_percentage = 0;
      }
      if (
        this.renovation.renovation_status !== "none" &&
        this.renovation.renovation_status !== "in_progress" &&
        this.renovation.renovation_status !== "done"
      ) {
        this.renovation.renovation_status = "none";
      }

      // 兜底：阶段值不在下拉里时，追加一项用于回显
      const stageVal = String(this.renovation.current_stage || "").trim();
      if (
        stageVal &&
        !(this.renovationStageOptions || []).some(
          (o) => String(o.value) === stageVal
        )
      ) {
        this.renovationStageOptions = (
          this.renovationStageOptions || []
        ).concat([{ label: stageVal, value: stageVal }]);
      }
    },
    addTag() {
      const t = String(this.newTag || "").trim();
      if (!t) {
        uni.showToast({ title: "请输入标签", icon: "none" });
        return;
      }
      if (this.tags.includes(t)) {
        this.newTag = "";
        return;
      }
      if (this.tags.length >= 6) {
        uni.showToast({ title: "最多6个标签", icon: "none" });
        return;
      }
      this.tags = this.tags.concat([t]);
      this.newTag = "";
    },
    removeTag(idx) {
      this.tags = this.tags.filter((_, i) => i !== idx);
    },
    openImageActions(idx) {
      const img = this.images[idx];
      if (!img) return;
      const isCover = img === this.form.cover_image;
      const list = isCover ? ["删除"] : ["设为封面", "删除"];
      uni.showActionSheet({
        itemList: list,
        success: (res) => {
          const tap = Number(res.tapIndex);
          if (!isCover && tap === 0) {
            this.form.cover_image = img;
            return;
          }
          const delIndex = isCover ? 0 : 1;
          if (tap === delIndex) {
            this.images = this.images.filter((_, i) => i !== idx);
            if (this.form.cover_image === img) {
              this.form.cover_image = this.images[0] || "";
            }
          }
        },
      });
    },
    openVideoActions() {
      const has = !!this.form.video_url;
      const list = has ? ["更换视频", "移除视频"] : ["上传视频"];
      uni.showActionSheet({
        itemList: list,
        success: async (res) => {
          const tap = Number(res.tapIndex);
          if (!has) {
            await this.pickVideo();
            return;
          }
          if (tap === 0) await this.pickVideo();
          if (tap === 1) this.form.video_url = "";
        },
      });
    },
    onRenovationStatusChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt =
        this.renovationStatusOptions[idx] || this.renovationStatusOptions[0];
      this.renovation.renovation_status = opt.value;
      if (opt.value === "none") {
        this.renovation.progress_percentage = 0;
        this.renovation.current_stage = "";
        this.renovation.start_date = "";
        this.renovation.estimated_finish_date = "";
        this.renovation.actual_finish_date = "";
        this.renovation.stage_logs = [];
        this.renovation.notes = "";
        this.renovationMaterials = [];
        this.renovationImages = [];
      }
      if (opt.value === "done") {
        this.renovation.progress_percentage = 100;
        if (
          !Array.isArray(this.renovation.stage_logs) ||
          this.renovation.stage_logs.length === 0
        ) {
          this.generateDefaultStageLogs({ overwrite: true });
        }
        this.syncStageLogsByCurrentStage({ silent: true });
      }
      if (opt.value === "in_progress") {
        if (
          !Array.isArray(this.renovation.stage_logs) ||
          this.renovation.stage_logs.length === 0
        ) {
          this.generateDefaultStageLogs({ overwrite: true });
          this.syncStageLogsByCurrentStage({ silent: true });
        }
      }
    },
    normalizeRenovationStageStatus(v) {
      const s = String(v || "").trim().toLowerCase();
      if (s === "done" || s === "finished" || s === "completed") return "done";
      if (s === "doing" || s === "in_progress" || s === "progress")
        return "doing";
      return "todo";
    },
    normalizeStageLogs(raw) {
      const arr = Array.isArray(raw) ? raw : [];
      const out = [];
      arr.forEach((it) => {
        if (!it) return;
        const stage = String(it.stage || it.stage_name || it.name || "").trim();
        if (!stage) return;
        const status = this.normalizeRenovationStageStatus(it.status);
        const date = String(it.date || "").trim();
        const note = String(it.note || it.notes || "").trim();

        let imgs = [];
        if (Array.isArray(it.images)) imgs = it.images;
        else if (typeof it.images === "string") {
          imgs = String(it.images || "")
            .split(",")
            .map((x) => String(x || "").trim())
            .filter(Boolean);
        }
        const normalizedImgs = imgs
          .map((u) => String(u || "").trim())
          .filter(Boolean);

        out.push({
          stage,
          status,
          date,
          note,
          images: Array.from(new Set(normalizedImgs)),
        });
      });
      return out;
    },
    ensureStageLogImages(log) {
      if (!log) return [];
      if (Array.isArray(log.images)) return log.images;
      if (typeof log.images === "string") {
        return String(log.images || "")
          .split(",")
          .map((x) => String(x || "").trim())
          .filter(Boolean);
      }
      return [];
    },
    getStageStatusText(v) {
      const s = this.normalizeRenovationStageStatus(v);
      if (s === "done") return "已完成";
      if (s === "doing") return "进行中";
      return "未开始";
    },
    getDefaultStageOrder() {
      const list = (this.renovationStageOptions || [])
        .map((o) => String((o && o.value) || "").trim())
        .filter(Boolean);
      if (list.length) return list;
      return ["设计", "拆改", "水电", "泥瓦", "木工", "油漆", "安装", "软装", "验收"];
    },
    confirmGenerateStageLogs() {
      uni.showModal({
        title: "生成默认工序",
        content: "将按默认工序生成时间线（会覆盖现有工序记录），是否继续？",
        confirmText: "继续",
        success: (res) => {
          if (!res.confirm) return;
          this.generateDefaultStageLogs({ overwrite: true });
          this.syncStageLogsByCurrentStage({ silent: true });
        },
      });
    },
    generateDefaultStageLogs({ overwrite = false } = {}) {
      const existing = Array.isArray(this.renovation.stage_logs)
        ? this.renovation.stage_logs
        : [];
      if (!overwrite && existing.length) return;
      const order = this.getDefaultStageOrder();
      this.renovation.stage_logs = order.map((stage) => ({
        stage,
        status: "todo",
        date: "",
        note: "",
        images: [],
      }));
    },
    syncStageLogsByCurrentStage(options = {}) {
      const silent = !!(options && options.silent);
      const overall = String(this.renovation.renovation_status || "none").trim();
      if (overall === "none") {
        this.renovation.stage_logs = [];
        return;
      }
      let logs = Array.isArray(this.renovation.stage_logs)
        ? this.normalizeStageLogs(this.renovation.stage_logs)
        : [];
      if (!logs.length) {
        this.generateDefaultStageLogs({ overwrite: true });
        logs = Array.isArray(this.renovation.stage_logs)
          ? this.normalizeStageLogs(this.renovation.stage_logs)
          : [];
      }

      if (overall === "done") {
        this.renovation.stage_logs = logs.map((it) => ({
          ...it,
          status: "done",
        }));
        return;
      }

      const cur = String(this.renovation.current_stage || "").trim();
      const order = logs.map((it) => it.stage);
      const curIdx = cur ? order.findIndex((s) => s === cur) : -1;
      if (curIdx < 0) {
        if (!silent)
          uni.showToast({ title: "当前阶段未选择/不在工序内", icon: "none" });
        this.renovation.stage_logs = logs;
        return;
      }
      this.renovation.stage_logs = logs.map((it, idx) => {
        const next = { ...it };
        if (idx < curIdx) next.status = "done";
        else if (idx === curIdx) next.status = "doing";
        else next.status = "todo";
        return next;
      });
    },
    addStageLog() {
      const logs = Array.isArray(this.renovation.stage_logs)
        ? this.normalizeStageLogs(this.renovation.stage_logs)
        : [];
      logs.push({
        stage: "",
        status: "todo",
        date: "",
        note: "",
        images: [],
      });
      this.renovation.stage_logs = logs;
    },
    clearStageLogs() {
      uni.showModal({
        title: "清空工序时间线",
        content: "清空后将删除所有工序记录（含工序图片），且无法恢复。",
        confirmText: "清空",
        confirmColor: "#ef4444",
        success: (res) => {
          if (!res.confirm) return;
          this.renovation.stage_logs = [];
        },
      });
    },
    removeStageLog(idx) {
      uni.showModal({
        title: "删除工序",
        content: "确定要删除该工序记录吗？",
        confirmText: "删除",
        confirmColor: "#ef4444",
        success: (res) => {
          if (!res.confirm) return;
          const logs = Array.isArray(this.renovation.stage_logs)
            ? this.renovation.stage_logs
            : [];
          this.renovation.stage_logs = logs.filter((_, i) => i !== idx);
        },
      });
    },
    moveStageLog(idx, step) {
      const logs = Array.isArray(this.renovation.stage_logs)
        ? this.renovation.stage_logs.slice(0)
        : [];
      const next = idx + step;
      if (next < 0 || next >= logs.length) return;
      const tmp = logs[idx];
      logs[idx] = logs[next];
      logs[next] = tmp;
      this.renovation.stage_logs = logs;
    },
    getStageLogStageIndex(log) {
      const v = String((log && log.stage) || "").trim();
      const idx = (this.renovationStageOptions || []).findIndex(
        (o) => String(o.value) === v
      );
      return idx >= 0 ? idx : 0;
    },
    onStageLogStageChange(logIdx, e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt =
        this.renovationStageOptions[idx] || this.renovationStageOptions[0];
      const stage = String((opt && opt.value) || "").trim();
      const logs = Array.isArray(this.renovation.stage_logs)
        ? this.renovation.stage_logs.slice(0)
        : [];
      if (!logs[logIdx]) return;
      logs[logIdx] = { ...logs[logIdx], stage };
      this.renovation.stage_logs = logs;
    },
    getStageLogStatusIndex(log) {
      const v = this.normalizeRenovationStageStatus(
        (log && log.status) || "todo"
      );
      const idx = (this.stageStatusOptions || []).findIndex(
        (o) => String(o.value) === v
      );
      return idx >= 0 ? idx : 0;
    },
    onStageLogStatusChange(logIdx, e) {
      const idx = Number(e && e.detail && e.detail.value) || 0;
      const opt = this.stageStatusOptions[idx] || this.stageStatusOptions[0];
      const status = this.normalizeRenovationStageStatus(opt && opt.value);
      const logs = Array.isArray(this.renovation.stage_logs)
        ? this.renovation.stage_logs.slice(0)
        : [];
      if (!logs[logIdx]) return;
      logs[logIdx] = { ...logs[logIdx], status };
      this.renovation.stage_logs = logs;
    },
    onStageLogDateChange(logIdx, e) {
      const date =
        e && e.detail && typeof e.detail.value !== "undefined"
          ? String(e.detail.value || "").trim()
          : "";
      const logs = Array.isArray(this.renovation.stage_logs)
        ? this.renovation.stage_logs.slice(0)
        : [];
      if (!logs[logIdx]) return;
      logs[logIdx] = { ...logs[logIdx], date };
      this.renovation.stage_logs = logs;
    },
    openStageLogImageActions(logIdx, imgIdx) {
      const logs = Array.isArray(this.renovation.stage_logs)
        ? this.renovation.stage_logs
        : [];
      const log = logs[logIdx];
      if (!log) return;
      const imgs = this.ensureStageLogImages(log);
      const img = imgs[imgIdx];
      if (!img) return;
      uni.showActionSheet({
        itemList: ["预览", "删除"],
        success: (res) => {
          const tap = Number(res.tapIndex);
          if (tap === 0) {
            const urls = imgs.map((u) => this.toFullMediaUrl(u)).filter(Boolean);
            const current = this.toFullMediaUrl(img);
            if (!urls.length || !current) return;
            uni.previewImage({ current, urls });
            return;
          }
          if (tap === 1) {
            const nextImgs = imgs.filter((_, i) => i !== imgIdx);
            const nextLogs = logs.slice(0);
            nextLogs[logIdx] = { ...log, images: nextImgs };
            this.renovation.stage_logs = nextLogs;
          }
        },
      });
    },
    pickStageLogImages(logIdx) {
      if (this.uploading) return;
      if (!this.id) {
        uni.showToast({ title: "请先保存房源", icon: "none" });
        return;
      }
      uni.chooseImage({
        count: 9,
        sizeType: ["compressed"],
        sourceType: ["album", "camera"],
        success: async (res) => {
          const files = (res && res.tempFilePaths) || [];
          if (!files.length) return;
          await this.uploadStageLogImages(files, logIdx);
        },
      });
    },
    uploadStageLogImages(filePaths = [], logIdx) {
      if (!Array.isArray(filePaths) || filePaths.length === 0)
        return Promise.resolve();
      return new Promise(async (resolve) => {
        this.uploading = true;
        uni.showLoading({ title: "上传中", mask: true });
        try {
          for (let i = 0; i < filePaths.length; i++) {
            const fp = filePaths[i];
            const out = await this.uploadSingle(fp, "image").catch(() => null);
            const url = out && out.code === 0 && out.data ? out.data.url : "";
            if (!url) continue;
            const logs = Array.isArray(this.renovation.stage_logs)
              ? this.renovation.stage_logs.slice(0)
              : [];
            const log = logs[logIdx];
            if (!log) continue;
            const imgs = this.ensureStageLogImages(log);
            logs[logIdx] = { ...log, images: imgs.concat([url]) };
            this.renovation.stage_logs = logs;
          }
        } finally {
          uni.hideLoading();
          this.uploading = false;
          resolve();
        }
      });
    },
    addMaterial() {
      const t = String(this.newMaterial || "").trim();
      if (!t) {
        uni.showToast({ title: "请输入材料", icon: "none" });
        return;
      }
      if (this.renovationMaterials.includes(t)) {
        this.newMaterial = "";
        return;
      }
      if (this.renovationMaterials.length >= 10) {
        uni.showToast({ title: "最多10个材料", icon: "none" });
        return;
      }
      this.renovationMaterials = this.renovationMaterials.concat([t]);
      this.newMaterial = "";
    },
    removeMaterial(idx) {
      this.renovationMaterials = this.renovationMaterials.filter(
        (_, i) => i !== idx
      );
    },
    openRenovationImageActions(idx) {
      const img = this.renovationImages[idx];
      if (!img) return;
      uni.showActionSheet({
        itemList: ["删除"],
        success: () => {
          this.renovationImages = this.renovationImages.filter(
            (_, i) => i !== idx
          );
        },
      });
    },
    pickRenovationImages() {
      if (this.uploading) return;
      if (!this.id) {
        uni.showToast({ title: "请先保存房源", icon: "none" });
        return;
      }
      uni.chooseImage({
        count: 9,
        sizeType: ["compressed"],
        sourceType: ["album", "camera"],
        success: async (res) => {
          const files = (res && res.tempFilePaths) || [];
          if (!files.length) return;
          await this.uploadFiles(files, "image", "renovation");
        },
      });
    },
    pickImages() {
      if (this.uploading) return;
      uni.chooseImage({
        count: 9,
        sizeType: ["compressed"],
        sourceType: ["album", "camera"],
        success: async (res) => {
          const files = (res && res.tempFilePaths) || [];
          if (!files.length) return;
          await this.uploadFiles(files, "image", "property");
        },
      });
    },
    pickVideo() {
      if (this.uploading) return Promise.resolve();
      return new Promise((resolve) => {
        uni.chooseVideo({
          sourceType: ["album", "camera"],
          maxDuration: 60,
          compressed: true,
          success: async (res) => {
            const filePath =
              res && (res.tempFilePath || res.tempFilePaths?.[0]);
            if (!filePath) return resolve();
            await this.uploadFiles([filePath], "video", "property");
            resolve();
          },
          fail: () => resolve(),
        });
      });
    },
    uploadFiles(filePaths = [], filetype = "image", target = "property") {
      if (!Array.isArray(filePaths) || filePaths.length === 0)
        return Promise.resolve();
      return new Promise(async (resolve) => {
        this.uploading = true;
        uni.showLoading({ title: "上传中", mask: true });
        try {
          for (let i = 0; i < filePaths.length; i++) {
            const fp = filePaths[i];
            const out = await this.uploadSingle(fp, filetype).catch(() => null);
            const url = out && out.code === 0 && out.data ? out.data.url : "";
            if (!url) continue;
            if (filetype === "video") {
              this.form.video_url = url;
            } else {
              if (target === "renovation") {
                this.renovationImages = this.renovationImages.concat([url]);
              } else {
                this.images = this.images.concat([url]);
                if (!this.form.cover_image) this.form.cover_image = url;
              }
            }
          }
        } finally {
          uni.hideLoading();
          this.uploading = false;
          resolve();
        }
      });
    },
    uploadSingle(filePath, filetype) {
      return new Promise((resolve, reject) => {
        const token = uni.getStorageSync("token");
        const timestamp = Math.floor(Date.now() / 1000);
        const passstr = md5(import.meta.env.GF_API_SECRET + timestamp);
        const header = {
          Accept: "text/json",
          Businessid: import.meta.env.GF_BUSINESUSSID,
          apiverify: base64Encode(passstr + "#" + timestamp),
        };
        if (token) header.Authorization = `${token}`;

        uni.uploadFile({
          url: `${baseUrl}/common/upload/upFile`,
          filePath,
          name: "file",
          formData: { filetype },
          header,
          success: (upRes) => {
            try {
              const raw = upRes && upRes.data ? upRes.data : "{}";
              const out = typeof raw === "string" ? JSON.parse(raw) : raw;
              if (!out || out.code !== 0) {
                uni.showToast({
                  title: (out && out.message) || "上传失败",
                  icon: "none",
                });
                reject(out);
                return;
              }
              resolve(out);
            } catch (e) {
              uni.showToast({ title: "解析上传结果失败", icon: "none" });
              reject(e);
            }
          },
          fail: (err) => {
            uni.showToast({ title: "上传失败", icon: "none" });
            reject(err);
          },
        });
      });
    },
    pickLocation() {
      uni.chooseLocation({
        success: (res) => {
          if (!res) return;
          const lat = Number(res.latitude);
          const lng = Number(res.longitude);
          if (isFinite(lat) && isFinite(lng)) {
            this.form.latitude = lat;
            this.form.longitude = lng;
          }
          // 兼容：有 name 则给到小区名称
          if (res.name && !this.form.community_name)
            this.form.community_name = String(res.name);
          if (res.address && !String(this.addressDetail || "").trim())
            this.addressDetail = String(res.address);
        },
        fail: () => {
          uni.showToast({ title: "无法打开地图", icon: "none" });
        },
      });
    },
    async save() {
      if (this.saving || this.uploading) return;
      const title = String(this.form.title || "").trim();
      if (!title) {
        uni.showToast({ title: "请填写房源标题", icon: "none" });
        return;
      }
      this.saving = true;
      try {
        const payload = {
          ...(this.id ? { id: this.id } : {}),
          title,
          sale_status: this.form.sale_status,
          price: this.form.price,
          price_unit: this.form.price_unit,
          area: this.form.area,
          rooms: Number(this.form.rooms || 0),
          halls: Number(this.form.halls || 0),
          bathrooms: Number(this.form.bathrooms || 0),
          community_name: String(this.form.community_name || "").trim(),
          address: this.buildFullAddress(),
          latitude: this.form.latitude || "",
          longitude: this.form.longitude || "",
          tags: this.tags,
          images: this.images,
          cover_image: this.form.cover_image || "",
          video_url: this.form.video_url || "",
          allow_image_download:
            Number(this.form.allow_image_download) === 1 ? 1 : 0,
          allow_video_download:
            Number(this.form.allow_video_download) === 1 ? 1 : 0,
          build_year: this.form.build_year ? Number(this.form.build_year) : 0,
          floor_level: String(this.form.floor_level || "").trim(),
          total_floors: Number(this.form.total_floors || 0),
          orientation: String(this.form.orientation || "").trim(),
          property_type: String(this.form.property_type || "").trim(),
          decoration_type: String(this.form.decoration_type || "").trim(),
          owner_name: String(this.form.owner_name || "").trim(),
          owner_phone: String(this.form.owner_phone || "").trim(),
          receiver_name: String(this.form.receiver_name || "").trim(),
          receiver_phone: String(this.form.receiver_phone || "").trim(),
          receiver_price: this.form.receiver_price,
          commission_rate: this.form.commission_rate,
          commission_reward: this.form.commission_reward,
          weigh: Number(this.form.weigh || 0),
          hot_status: Number(this.form.hot_status) === 1 ? 1 : 0,
          status: Number(this.form.status) === 0 ? 0 : 1,
        };
        const res = await propertyApi.saveManage(payload);
        if (!res || res.code !== 0) return;

        const createdNew = !this.id;
        const nid = Number(res && res.data && res.data.id) || 0;
        if (!this.id && nid) {
          this.id = nid;
        }

        // 装修信息跟随底部保存统一提交
        let renovationSaved = true;
        if (this.id) {
          renovationSaved = await this.saveRenovation({
            silent: true,
            reload: false,
          });
        }

        // 新增后刷新一次装修数据，保证页面状态与服务端一致
        if (createdNew && this.id) {
          await this.loadRenovation();
        }

        if (renovationSaved) {
          uni.showToast({ title: "保存成功", icon: "none" });
        } else {
          uni.showToast({ title: "房源已保存，装修保存失败", icon: "none" });
        }
      } finally {
        this.saving = false;
      }
    },
    async saveRenovation(options = {}) {
      const silent = !!(options && options.silent);
      const reload = !options || options.reload !== false;

      if (!this.id) {
        if (!silent) uni.showToast({ title: "请先保存房源", icon: "none" });
        return false;
      }
      if (this.savingRenovation || this.uploading) return false;

      const status = String(this.renovation.renovation_status || "none");
      const progress = Number(this.renovation.progress_percentage || 0);
      if (!isFinite(progress) || progress < 0 || progress > 100) {
        if (!silent)
          uni.showToast({ title: "进度需在0-100之间", icon: "none" });
        return false;
      }

      this.savingRenovation = true;
      try {
        let stageLogs = this.normalizeStageLogs(this.renovation.stage_logs);
        if (status === "none") {
          stageLogs = [];
        }
        const payload = {
          property_id: this.id,
          renovation_status: status,
          progress_percentage: Math.floor(progress),
          current_stage: String(this.renovation.current_stage || "").trim(),
          start_date: this.renovation.start_date || "",
          estimated_finish_date: this.renovation.estimated_finish_date || "",
          actual_finish_date: this.renovation.actual_finish_date || "",
          materials: this.renovationMaterials,
          images: this.renovationImages,
          stage_logs: stageLogs,
          notes: String(this.renovation.notes || "").trim(),
          status: Number(this.renovation.status) === 1 ? 1 : 0,
        };
        const res = await propertyApi.saveManageRenovation(payload);
        if (!res || res.code !== 0) return false;
        if (!silent) {
          uni.showToast({ title: "保存成功", icon: "none" });
        }
        if (reload) {
          await this.loadRenovation();
        }
        return true;
      } finally {
        this.savingRenovation = false;
      }
    },
    preview() {
      if (!this.id) return;
      uni.navigateTo({
        url: `/pages/property_detail/property_detail?id=${this.id}`,
      });
    },
  },
};
</script>

<style>
.pe {
  --pe-bg-a: #edf3fa;
  --pe-bg-b: #f9fbff;
  --pe-card-bg: rgba(255, 255, 255, 0.9);
  --pe-line: #dbe6f2;
  --pe-text-1: #0f172a;
  --pe-text-2: #334155;
  --pe-text-3: #64748b;
  --pe-accent-1: #2d9cf0;
  --pe-accent-2: #2563eb;

  height: 100vh;
  background: linear-gradient(
    160deg,
    var(--pe-bg-a) 0%,
    var(--pe-bg-b) 52%,
    #f1f6fd 100%
  );
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

.pe::before {
  content: "";
  position: absolute;
  inset: 0;
  pointer-events: none;
  background: radial-gradient(
      560rpx 360rpx at -12% -6%,
      rgba(45, 156, 240, 0.16),
      transparent 62%
    ),
    radial-gradient(
      680rpx 400rpx at 112% 4%,
      rgba(37, 99, 235, 0.12),
      transparent 66%
    );
  z-index: 0;
}

.icon-btn {
  width: 72rpx;
  height: 72rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(255, 255, 255, 0.95);
  background: linear-gradient(145deg, #ffffff, #edf3fa);
  box-shadow: 8rpx 8rpx 16rpx rgba(148, 163, 184, 0.18),
    -8rpx -8rpx 16rpx rgba(255, 255, 255, 0.94);
  color: var(--pe-text-1);
}

.icon-btn.primary {
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  border: none;
  color: #ffffff;
  box-shadow: 0 12rpx 28rpx rgba(37, 99, 235, 0.22);
}

.content {
  flex: 1;
  overflow: hidden;
  position: relative;
  z-index: 1;
}

.card {
  margin: 18rpx 24rpx 0;
  padding: 22rpx 20rpx 16rpx;
  background: var(--pe-card-bg);
  border-radius: 30rpx;
  border: 1px solid rgba(255, 255, 255, 0.95);
  box-shadow: 14rpx 14rpx 28rpx rgba(148, 163, 184, 0.2),
    -14rpx -14rpx 28rpx rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  position: relative;
  z-index: 1;
}

.card::before {
  content: "";
  position: absolute;
  top: 0;
  left: 22rpx;
  width: 86rpx;
  height: 8rpx;
  border-radius: 0 0 999rpx 999rpx;
  background: linear-gradient(135deg, var(--pe-accent-1), var(--pe-accent-2));
  opacity: 0.75;
}

.card-title {
  font-size: 31rpx;
  font-weight: 800;
  color: var(--pe-text-1);
  margin-bottom: 16rpx;
  letter-spacing: 0.2rpx;
}

.hint {
  font-size: 24rpx;
  color: #76879d;
  margin-top: -4rpx;
  margin-bottom: 14rpx;
  line-height: 1.5;
}
.hint.small {
  margin: 0;
}
.hint.small.tip {
  margin-top: 10rpx;
  color: #64748b;
}

.form-row {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  margin-bottom: 14rpx;
}

.label {
  font-size: 24rpx;
  color: var(--pe-text-2);
  font-weight: 700;
  letter-spacing: 0.2rpx;
}
.label.required::after {
  content: " *";
  color: #ef4444;
}

.input {
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e5edf7;
  background: linear-gradient(145deg, #f9fbfe, #eef3f9);
  box-shadow: inset 6rpx 6rpx 12rpx rgba(208, 220, 234, 0.36),
    inset -6rpx -6rpx 12rpx rgba(255, 255, 255, 0.95);
  padding: 0 16rpx;
  font-size: 28rpx;
  color: var(--pe-text-1);
  box-sizing: border-box;
}

.textarea {
  width: 100%;
  min-height: 160rpx;
  border-radius: 20rpx;
  border: 1px solid #e5edf7;
  background: linear-gradient(145deg, #f9fbfe, #eef3f9);
  box-shadow: inset 6rpx 6rpx 12rpx rgba(208, 220, 234, 0.34),
    inset -6rpx -6rpx 12rpx rgba(255, 255, 255, 0.92);
  padding: 14rpx 16rpx;
  font-size: 28rpx;
  color: var(--pe-text-1);
  box-sizing: border-box;
}
.textarea.small {
  min-height: 120rpx;
}

.placeholder {
  color: #94a3b8;
}

/* 横屏/宽屏：表单占满宽度，不保留左右留白 */
@media screen and (orientation: landscape), screen and (min-width: 750px) {
  .card {
    margin-left: 0;
    margin-right: 0;
    border-radius: 0;
  }

  .footer {
    padding-left: 0;
    padding-right: 0;
    border-radius: 0;
  }
}

.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14rpx;
}

.stage-logs {
  margin-top: 14rpx;
  padding-top: 14rpx;
  border-top: 1px dashed #dce6f2;
}

.stage-logs-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12rpx;
}

.stage-logs-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 10rpx;
}

.mini-btn {
  padding: 10rpx 14rpx;
  border-radius: 999rpx;
  border: 1px solid rgba(37, 99, 235, 0.28);
  background: rgba(37, 99, 235, 0.12);
  color: #2563eb;
  font-size: 22rpx;
  font-weight: 800;
  line-height: 1;
}

.mini-btn.ghost {
  border-color: rgba(100, 116, 139, 0.22);
  background: rgba(148, 163, 184, 0.12);
  color: #334155;
}

.mini-btn.danger {
  border-color: rgba(239, 68, 68, 0.25);
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
}

.mini-btn[disabled] {
  opacity: 0.6;
}

.stage-log-list {
  margin-top: 12rpx;
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.stage-log-item {
  border-radius: 24rpx;
  border: 1px solid #dde8f4;
  background: linear-gradient(145deg, #fbfdff, #edf4fb);
  box-shadow: 8rpx 8rpx 16rpx rgba(169, 182, 199, 0.18),
    -6rpx -6rpx 14rpx rgba(255, 255, 255, 0.78);
  padding: 14rpx;
}

.stage-log-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12rpx;
  margin-bottom: 10rpx;
}

.stage-log-title {
  display: flex;
  align-items: center;
  gap: 10rpx;
  min-width: 0;
}

.stage-log-title .name {
  font-size: 26rpx;
  font-weight: 900;
  color: var(--pe-text-1);
  max-width: 260rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stage-log-title .tag {
  padding: 6rpx 12rpx;
  border-radius: 999rpx;
  font-size: 20rpx;
  font-weight: 900;
  background: rgba(148, 163, 184, 0.18);
  color: #475569;
}

.stage-log-title .tag.doing {
  background: rgba(37, 99, 235, 0.14);
  color: #1d4ed8;
}
.stage-log-title .tag.done {
  background: rgba(34, 197, 94, 0.14);
  color: #15803d;
}

.stage-log-tools {
  display: flex;
  align-items: center;
  gap: 10rpx;
  flex-shrink: 0;
}

.mini-icon-btn {
  padding: 8rpx 12rpx;
  border-radius: 999rpx;
  border: 1px solid rgba(100, 116, 139, 0.22);
  background: rgba(148, 163, 184, 0.12);
  color: #334155;
  font-size: 22rpx;
  font-weight: 800;
  line-height: 1;
}
.mini-icon-btn.danger {
  border-color: rgba(239, 68, 68, 0.25);
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
}
.mini-icon-btn[disabled] {
  opacity: 0.5;
}

.stage-log-grid {
  margin-top: 6rpx;
}

.stage-log-note {
  grid-column: 1 / -1;
}

.media-grid.stage-media-grid {
  grid-template-columns: repeat(4, 1fr);
}

.media-item.stage-media-item {
  height: 160rpx;
}

.media-item.add.stage-media-add {
  height: 160rpx;
}

.picker {
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e5edf7;
  background: linear-gradient(145deg, #f9fbfe, #eef3f9);
  box-shadow: inset 6rpx 6rpx 12rpx rgba(208, 220, 234, 0.32),
    inset -6rpx -6rpx 12rpx rgba(255, 255, 255, 0.9);
  padding: 0 14rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.picker-text {
  font-size: 28rpx;
  color: var(--pe-text-1);
}

.picker-ic {
  font-size: 36rpx !important;
  color: #7f93aa;
}
.picker.disabled {
  background: linear-gradient(145deg, #f0f4fa, #e9f0f8);
  border-color: #d8e2ef;
  opacity: 0.9;
}

.price-row {
  display: flex;
  gap: 10rpx;
  align-items: center;
}

.unit-picker {
  width: 140rpx;
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e5edf7;
  background: linear-gradient(145deg, #f9fbfe, #eef3f9);
  box-shadow: inset 6rpx 6rpx 12rpx rgba(208, 220, 234, 0.32),
    inset -6rpx -6rpx 12rpx rgba(255, 255, 255, 0.9);
  padding: 0 12rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
}

.unit-text {
  font-size: 28rpx;
  color: var(--pe-text-1);
}
.unit-ic {
  font-size: 34rpx !important;
  color: #7f93aa;
}

.room-row {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.mini {
  width: 120rpx;
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e5edf7;
  background: linear-gradient(145deg, #f9fbfe, #eef3f9);
  box-shadow: inset 6rpx 6rpx 12rpx rgba(208, 220, 234, 0.32),
    inset -6rpx -6rpx 12rpx rgba(255, 255, 255, 0.9);
  padding: 0 12rpx;
  font-size: 28rpx;
  color: var(--pe-text-1);
  box-sizing: border-box;
}

.mini-suffix {
  font-size: 26rpx;
  color: #64748b;
  margin-right: 8rpx;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12rpx;
}

.media-item {
  height: 200rpx;
  border-radius: 20rpx;
  overflow: hidden;
  position: relative;
  background: linear-gradient(145deg, #f4f8fd, #e9f0f8);
  border: 1px solid #dde8f4;
  box-shadow: 8rpx 8rpx 16rpx rgba(169, 182, 199, 0.2),
    -6rpx -6rpx 14rpx rgba(255, 255, 255, 0.78);
}

.media-item.add {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 8rpx;
  background: linear-gradient(145deg, #fbfdff, #edf4fb);
}

.media-img {
  width: 100%;
  height: 100%;
}

.badge {
  position: absolute;
  left: 10rpx;
  top: 10rpx;
  padding: 6rpx 12rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  background: rgba(37, 99, 235, 0.85);
  color: #ffffff;
  backdrop-filter: blur(8px);
}

.add-ic {
  font-size: 54rpx !important;
  color: #2563eb;
}

.add-text {
  font-size: 24rpx;
  color: #2563eb;
  font-weight: 700;
}

.switch-row {
  margin-top: 14rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14rpx 4rpx 0;
  border-top: 1px dashed #dce6f2;
}

.switch-left {
  display: flex;
  align-items: center;
  gap: 10rpx;
  color: var(--pe-text-2);
}

.sw-ic {
  font-size: 34rpx !important;
  color: #5f7a98;
}

.sw-text {
  font-size: 26rpx;
  font-weight: 700;
}

.video-box {
  height: 280rpx;
  border-radius: 24rpx;
  overflow: hidden;
  border: 1px solid #dbe6f3;
  background: linear-gradient(160deg, #0f2237, #132f4b 56%, #0f1d31);
  box-shadow: inset 0 0 0 1rpx rgba(255, 255, 255, 0.04);
}

.video-empty {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
  color: rgba(255, 255, 255, 0.9);
}

.video-ic {
  font-size: 64rpx !important;
  opacity: 0.85;
}

.video-texts {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.video-title {
  font-size: 30rpx;
  font-weight: 800;
}

.video-sub {
  font-size: 24rpx;
  opacity: 0.8;
}

.video {
  width: 100%;
  height: 100%;
}

.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10rpx;
}

.tag {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 14rpx;
  border-radius: 999rpx;
  background: linear-gradient(145deg, #f8fbff, #edf3fb);
  border: 1px solid #dce6f2;
  box-shadow: 6rpx 6rpx 12rpx rgba(177, 190, 206, 0.14),
    -4rpx -4rpx 10rpx rgba(255, 255, 255, 0.86);
}

.tag-text {
  font-size: 24rpx;
  color: #0f172a;
  font-weight: 700;
}

.tag-x {
  font-size: 26rpx !important;
  color: #64748b;
}

.tag-add {
  margin-top: 12rpx;
  display: flex;
  gap: 12rpx;
  align-items: center;
}

.btn {
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #dce7f3;
  background: linear-gradient(145deg, #ffffff, #eff5fc);
  font-size: 26rpx;
  color: var(--pe-text-1);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 18rpx;
  box-shadow: 8rpx 8rpx 16rpx rgba(169, 183, 201, 0.18),
    -6rpx -6rpx 14rpx rgba(255, 255, 255, 0.9);
}
.btn::after {
  border: none;
}
.btn.add {
  background: linear-gradient(135deg, var(--pe-accent-1), var(--pe-accent-2));
  border: none;
  color: #ffffff;
  font-weight: 800;
  min-width: 140rpx;
  box-shadow: 0 12rpx 24rpx rgba(37, 99, 235, 0.24);
}

.btn.ghost {
  background: linear-gradient(145deg, #f9fbff, #edf4fc);
  color: #1f6ed7;
  border-color: rgba(37, 99, 235, 0.2);
  min-width: 160rpx;
}

.renovation-actions {
  margin-top: 12rpx;
  display: flex;
  justify-content: flex-end;
}

.coord-row {
  margin-top: 12rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
  padding-top: 10rpx;
  border-top: 1px dashed #dce6f2;
}

.coord-left {
  display: flex;
  align-items: center;
  gap: 10rpx;
  min-width: 0;
}

.coord-ic {
  font-size: 34rpx !important;
  color: #5f7a98;
}

.coord-text {
  font-size: 24rpx;
  color: var(--pe-text-2);
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.bottom-spacer {
  height: calc(env(safe-area-inset-bottom) + 120rpx);
}

.footer {
  padding: 14rpx 24rpx calc(env(safe-area-inset-bottom) + 14rpx);
  background: rgba(242, 248, 255, 0.88);
  backdrop-filter: blur(12px);
  border-top: 1rpx solid rgba(214, 227, 242, 0.9);
  box-shadow: 0 -10rpx 24rpx rgba(113, 136, 162, 0.12);
  display: flex;
  gap: 14rpx;
  position: relative;
  z-index: 2;
}

.footer-btn {
  flex: 1;
  height: 88rpx;
  border-radius: 24rpx;
  font-size: 28rpx;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #dce7f3;
  background: linear-gradient(145deg, #ffffff, #eff5fc);
  color: var(--pe-text-1);
  box-shadow: 8rpx 8rpx 16rpx rgba(169, 183, 201, 0.2),
    -6rpx -6rpx 14rpx rgba(255, 255, 255, 0.88);
}
.footer-btn::after {
  border: none;
}

.footer-btn.primary {
  border: none;
  color: #ffffff;
  background: linear-gradient(135deg, var(--pe-accent-1), var(--pe-accent-2));
  box-shadow: 0 14rpx 28rpx rgba(37, 99, 235, 0.28);
}

.footer-btn.ghost {
  background: linear-gradient(145deg, #f9fbff, #edf4fc);
  color: #1f6ed7;
  border-color: rgba(37, 99, 235, 0.2);
}

.footer-btn:disabled,
.btn:disabled {
  opacity: 0.58;
  box-shadow: none;
}
</style>
