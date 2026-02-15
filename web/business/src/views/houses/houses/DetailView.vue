<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :isPadding="false"
    :loading="loading"
    width="min(1380px, 98vw)"
    @height-change="onHeightChange"
    :minHeight="modelHeight"
    :height="modelHeight"
    :title="''"
    :showOkBtn="false"
    :footer="null"
    class="house-detail-modal"
  >
    <div class="detail-container" :style="{ 'min-height': `${windHeight}px` }">
      <!-- Detail Header -->
      <div class="detail-header">
        <a-card class="header-card" :bordered="false">
          <div class="header-inner">
            <div class="header-left">
              <div class="cover-preview">
                <a-image
                  :src="getImageUrl(detailData.cover_image)"
                  alt="房源封面"
                  width="96"
                  height="72"
                  fit="cover"
                  class="preview-img"
                />
              </div>
              <div class="header-info">
                <div class="title-row">
                  <div class="detail-title">{{ detailData.title || '-' }}</div>
                  <a-space size="mini" wrap>
                    <a-tag
                      :color="detailData.sale_status === 'on_sale' ? 'green' : detailData.sale_status === 'sold' ? 'gray' : 'orange'"
                    >
                      {{ getSaleStatusLabel(detailData.sale_status) }}
                    </a-tag>
                    <a-tag v-if="detailData.hot_status === 1" color="red">
                      <icon-fire />
                      推荐
                    </a-tag>
                    <a-tag v-if="detailData.status === 1" color="gray">禁用</a-tag>
                  </a-space>
                </div>

                <div class="detail-meta">
                  <span class="meta-item">
                    <icon-location />
                    {{ detailData.community_name || '-' }}
                  </span>
                  <span class="meta-item" v-if="detailData.floor_level">
                    {{ detailData.floor_level }}
                  </span>
                  <span class="meta-item addr" v-if="detailData.address">
                    {{ detailData.address }}
                  </span>
                </div>

                <div class="price-display">
                  <span class="price-amount">{{ priceAmountText }}</span>
                  <span class="price-unit">{{ detailData.price_unit || '万' }}</span>
                  <span class="price-per" v-if="pricePerSquareText">
                    ({{ pricePerSquareText }})
                  </span>
                </div>
              </div>
            </div>

            <a-button type="text" shape="circle" class="close-btn" aria-label="关闭" @click="closeModal">
              <icon-close />
            </a-button>
          </div>
        </a-card>
      </div>

      <!-- Detail Body Content -->
      <div class="detail-body">
        <a-tabs v-model:active-key="mainTabKey" type="line" size="medium">
           <a-tab-pane key="1" title="房源档案">
            <div class="tab-content">
               <div class="detail-shell">
                 <!-- Top: 基础信息 + 销售维护（左右并排） -->
                 <div class="detail-top-grid">
                   <!-- Basic Info Card -->
                   <div class="detail-card">
                     <div class="card-header">
                       <icon-file class="card-icon" />
                       <span class="card-title">基础信息</span>
                     </div>
                     <div class="card-body">
                       <div class="info-grid">
                         <div class="info-item">
                           <label>小区名称</label>
                           <div class="value">{{ detailData.community_name || '-' }}</div>
                         </div>
                         <div class="info-item">
                           <label>所在区域</label>
                           <div class="value">{{ districtText || '-' }}</div>
                         </div>
                         <div class="info-item full-width">
                           <label>详细地址</label>
                           <div class="value">{{ addressDetailText || '-' }}</div>
                         </div>
                         <div class="info-item">
                           <label>物业类型</label>
                           <div class="value">{{ detailData.property_type || '-' }}</div>
                         </div>
                         <div class="info-item">
                           <label>建筑年代</label>
                           <div class="value">{{ detailData.build_year ? `${detailData.build_year}年` : '-' }}</div>
                         </div>
                         <div class="info-item">
                           <label>装修情况</label>
                           <div class="value">{{ detailData.decoration_type || '-' }}</div>
                         </div>
                       </div>
                     </div>
                   </div>

                   <!-- Sales Info Card -->
                   <div class="detail-card">
                     <div class="card-header">
                       <icon-user-group class="card-icon" />
                       <span class="card-title">销售维护</span>
                       <div class="card-actions">
                         <a-button size="mini" type="text" @click="openStatusLogs">
                           <template #icon><icon-history /></template>
                           状态记录
                         </a-button>
                       </div>
                     </div>
                     <div class="card-body">
                       <div class="info-grid">
                         <div class="info-item">
                           <label>房主姓名</label>
                           <div class="value">{{ detailData.owner_name || '-' }}</div>
                         </div>
                         <div class="info-item">
                           <label>房主电话</label>
                           <div class="value phone-value">{{ detailData.owner_phone || '-' }}</div>
                         </div>
                         <div class="info-item">
                           <label>收房人姓名</label>
                           <div class="value">{{ detailData.receiver_name || '-' }}</div>
                         </div>
                         <div class="info-item">
                           <label>收房人电话</label>
                           <div class="value phone-value">{{ detailData.receiver_phone || '-' }}</div>
                         </div>
                         <div class="info-item full-width">
                           <label>收房价格(支付业主)</label>
                           <div class="value">{{ formatMoney(detailData.receiver_price) }}</div>
                         </div>
                         <div class="info-item">
                           <label>浏览热度</label>
                           <div class="value highlight">{{ detailData.view_count || 0 }}</div>
                         </div>
                         <div class="info-item">
                           <label>带看次数</label>
                           <div class="value highlight">{{ detailData.showing_count || 0 }}</div>
                         </div>
                       </div>
                     </div>
                   </div>
                 </div>

                 <!-- Bottom: 参数/标签/门锁 -->
                 <div class="detail-bottom-grid">
                   <div class="bottom-main">
                     <!-- Stats Card -->
                     <div class="detail-card">
                       <div class="card-header">
                         <icon-bar-chart class="card-icon" />
                         <span class="card-title">房源参数</span>
                       </div>
                       <div class="card-body">
                         <div class="stats-grid">
                           <div class="stat-item">
                             <div class="stat-value">{{ detailData.rooms }}室{{ detailData.halls }}厅</div>
                             <div class="stat-label">户型</div>
                           </div>
                           <div class="stat-item">
                             <div class="stat-value">{{ detailData.area }}㎡</div>
                             <div class="stat-label">面积</div>
                           </div>
                           <div class="stat-item">
                             <div class="stat-value">{{ detailData.floor_level }}</div>
                             <div class="stat-label">楼层</div>
                           </div>
                           <div class="stat-item">
                             <div class="stat-value">{{ detailData.orientation || '暂无' }}</div>
                             <div class="stat-label">朝向</div>
                           </div>
                       </div>
                      </div>
                     </div>
                     <!-- Custom Description -->
                     <div class="detail-card">
                       <div class="card-header">
                         <icon-file class="card-icon" />
                         <span class="card-title">房源描述</span>
                       </div>
                       <div class="card-body">
                         <div v-if="customDescText" class="desc-text">{{ customDescText }}</div>
                         <a-empty v-else description="暂无描述" />
                       </div>
                     </div>
                   </div>

                   <div class="bottom-side">
                     <!-- Tags -->
                     <div class="detail-card">
                       <div class="card-header">
                         <icon-tag class="card-icon" />
                         <span class="card-title">标签</span>
                       </div>
                       <div class="card-body">
                         <div class="tags-wrap" v-if="tagList.length">
                           <a-tag v-for="(t, i) in tagList.slice(0, 12)" :key="i" size="small" bordered class="pill-tag">
                             {{ t }}
                           </a-tag>
                         </div>
                         <a-empty v-else description="暂无标签" />
                       </div>
                     </div>

                     <!-- Smart Lock Widget -->
                     <div class="pro-card lock-widget" :class="{ 'active': lockInfo?.lock }">
                       <div class="widget-header">
                         <div class="title"><icon-lock /> 智能门锁</div>
                         <div class="status" v-if="lockInfo?.lock">
                           <div class="battery-indicator" :class="getBatteryClass(lockInfo.lock.battery)">
                             <icon-thunderbolt /> {{ lockInfo.lock.battery }}%
                           </div>
                         </div>
                         <a-button v-else type="outline" size="mini" class="btn-bind" @click="openBindLock">
                           <icon-plus /> 绑定设备
                         </a-button>
                       </div>

                       <div class="widget-body" v-if="lockInfo?.lock">
                         <div class="device-info">
                           <div class="device-icon">
                             <icon-robot />
                           </div>
                           <div class="device-text">
                             <div class="d-name">{{ lockInfo.lock.lock_alias || lockInfo.lock.lock_name }}</div>
                             <div class="d-mac">{{ lockInfo.lock.lock_mac }}</div>
                           </div>
                         </div>
                         <div class="control-grid">
                           <a-button class="control-btn primary" @click="handleRemoteUnlock">
                             <template #icon><icon-unlock /></template>
                             获取密码
                           </a-button>
                           <a-button class="control-btn" @click="openCloudDetail">
                             <template #icon><icon-settings /></template>
                             详情
                           </a-button>
                           <a-popconfirm content="解绑后无法远程控制，确认?" @ok="handleUnbind">
                             <a-button class="control-btn danger">
                               <template #icon><icon-unlink /></template>
                               解绑
                             </a-button>
                           </a-popconfirm>
                         </div>
                         <div class="sync-time">
                           上次同步: {{ formatTime(lockInfo.lock.last_sync_at) }}
                         </div>
                       </div>
                       <div class="widget-empty" v-else>
                         <div class="empty-icon"><icon-lock /></div>
                         <p>未绑定智能门锁设备</p>
                       </div>
                     </div>
                   </div>
                 </div>
               </div>
            </div>
          </a-tab-pane>

           <a-tab-pane key="2" title="图片相册">
              <div class="gallery-wrapper">
                <a-image-preview-group infinite>
                   <div class="masonry-grid">
                      <div class="grid-item" v-for="(img, idx) in galleryImages" :key="idx">
                         <a-image
                            :src="getImageUrl(img)"
                            :alt="`房源图片 ${idx + 1}`"
                            width="100%"
                            height="100%"
                            fit="cover"
                         />
                        <div class="img-actions" v-if="allowImageDownload">
                          <a-tooltip content="下载图片" position="tr">
                            <a-button
                              class="download-btn"
                              type="primary"
                              size="mini"
                              shape="circle"
                              @click.stop="downloadUrl(getImageUrl(img), `${detailData.title || '房源'}-图片-${idx + 1}.jpg`)"
                            >
                              <template #icon><icon-download /></template>
                            </a-button>
                          </a-tooltip>
                        </div>
                      </div>
                   </div>
                </a-image-preview-group>
                <a-empty v-if="galleryImages.length === 0" description="暂无图片" />
              </div>
           </a-tab-pane>

          <a-tab-pane key="3" title="房源视频">
            <div class="video-wrapper">
              <div class="video-card" v-if="videoSrc">
                <video class="video-player" :src="videoSrc" controls :poster="getImageUrl(detailData.cover_image)" />
                <div class="video-actions">
                  <a-space size="mini">
                    <a-link :href="videoSrc" target="_blank">新窗口打开</a-link>
                    <a-button v-if="allowVideoDownload" size="mini" @click="downloadUrl(videoSrc, `${detailData.title || '房源'}-视频.mp4`)">
                      <template #icon><icon-download /></template>
                      下载
                    </a-button>
                  </a-space>
                </div>
              </div>
              <a-empty v-else description="暂无视频" />
            </div>
          </a-tab-pane>

          <a-tab-pane key="4" title="装修进度">
             <div class="renovation-timeline" v-if="renovationData?.renovation_status">
                <div class="timeline-header">
                   <div class="status-badge" :class="renovationData.renovation_status">
                      {{ getRenovationStatusLabel(renovationData.renovation_status) }}
                   </div>
                   <div class="progress-info" v-if="renovationData.renovation_status !== 'none'">
                      <span>装修进度</span>
                      <a-progress 
                        :percent="(renovationData.progress_percentage || 0) / 100" 
                        size="large"
                        :color="{
                            '0%': 'rgb(var(--primary-6))',
                            '100%': 'rgb(var(--success-6))',
                        }"
                      />
                   </div>
                </div>
                
                <div class="timeline-grid">
                   <div class="t-card">
                      <div class="label">当前阶段</div>
                      <div class="val highlight">{{ renovationData.current_stage || '未开始' }}</div>
                   </div>
                   <div class="t-card">
                      <div class="label">预计完工</div>
                      <div class="val">{{ renovationData.estimated_finish_date || '-' }}</div>
                   </div>
                   <div class="t-card wide" v-if="renovationData.notes">
                      <div class="label">施工备注</div>
                      <div class="val text-area">{{ renovationData.notes }}</div>
                   </div>
                </div>

                <div class="stage-timeline" v-if="renovationStageTimeline.length">
                  <div class="stage-item" v-for="(it, idx) in renovationStageTimeline" :key="it.stage + '_' + idx">
                    <div class="stage-rail">
                      <div class="stage-dot" :class="it.status"></div>
                      <div class="stage-line" v-if="idx !== renovationStageTimeline.length - 1"></div>
                    </div>
                    <div class="stage-body">
                      <div class="stage-head">
                        <div class="stage-title">{{ it.stage }}</div>
                        <a-tag size="small" :color="getStageStatusColor(it.status)">{{ getStageStatusText(it.status) }}</a-tag>
                        <span v-if="it.date" class="stage-date">{{ it.date }}</span>
                      </div>
                      <div v-if="it.note" class="stage-note">{{ it.note }}</div>
                      <div class="stage-media" v-if="it.images && it.images.length">
                        <a-image-preview-group infinite>
                          <div class="stage-grid">
                            <div class="stage-img" v-for="(img, sidx) in it.images" :key="sidx">
                              <a-image :src="img" :alt="`装修-${it.stage}-${sidx + 1}`" width="100%" height="100%" fit="cover" />
                              <div class="img-actions" v-if="allowImageDownload">
                                <a-tooltip content="下载图片" position="tr">
                                  <a-button
                                    class="download-btn"
                                    type="primary"
                                    size="mini"
                                    shape="circle"
                                    @click.stop="downloadUrl(img, `${detailData.title || '房源'}-装修-${it.stage}-${sidx + 1}.jpg`)"
                                  >
                                    <template #icon><icon-download /></template>
                                  </a-button>
                                </a-tooltip>
                              </div>
                            </div>
                          </div>
                        </a-image-preview-group>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="renovation-media" v-if="renovationImages.length">
                   <div class="media-header">
                      <span class="media-title">装修图片</span>
                      <a-tag size="small" bordered>{{ renovationImages.length }}张</a-tag>
                   </div>
                   <a-image-preview-group infinite>
                      <div class="renovation-grid">
                         <div class="renovation-item" v-for="(img, idx) in renovationImages" :key="idx">
                            <a-image :src="img" :alt="`装修图片 ${idx + 1}`" width="100%" height="100%" fit="cover" />
                            <div class="img-actions" v-if="allowImageDownload">
                              <a-tooltip content="下载图片" position="tr">
                                <a-button
                                  class="download-btn"
                                  type="primary"
                                  size="mini"
                                  shape="circle"
                                  @click.stop="downloadUrl(img, `${detailData.title || '房源'}-装修-${idx + 1}.jpg`)"
                                >
                                  <template #icon><icon-download /></template>
                                </a-button>
                              </a-tooltip>
                            </div>
                         </div>
                      </div>
                   </a-image-preview-group>
                </div>
             </div>
             <a-empty v-else description="暂无装修信息" />
          </a-tab-pane>

          <a-tab-pane key="5" title="经纪人记录">
            <div class="records-wrapper">
              <a-tabs v-model:active-key="recordTabKey" type="card-gutter" size="small" class="records-tabs">
                <a-tab-pane key="view">
                  <template #title>
                    <icon-eye /> 浏览记录
                  </template>

                  <div class="records-toolbar">
                    <a-space size="mini" wrap>
                      <a-button size="mini" type="primary" :loading="viewLogLoading" @click="reloadViewLogs">
                        <template #icon><icon-refresh /></template>
                        刷新
                      </a-button>
                      <a-tag bordered size="small">共 {{ viewLogPagination.total || 0 }} 条</a-tag>
                    </a-space>
                    <div class="records-tip">说明：每次进入房源详情记录 1 条（不去重），数据来自小程序/业务端的行为日志。</div>
                  </div>

                  <a-spin :loading="viewLogLoading" style="width: 100%">
                    <a-table
                      row-key="id"
                      :data="viewLogList"
                      :pagination="viewLogPagination"
                      :bordered="false"
                      size="mini"
                      @page-change="handleViewLogPageChange"
                      @page-size-change="handleViewLogPageSizeChange"
                    >
                      <template #columns>
                        <a-table-column title="时间" data-index="createtime" :width="170" />
                        <a-table-column title="经纪人" :width="220">
                          <template #cell="{ record }">
                            <div class="user-cell">
                              <div class="user-main">
                                <span class="name">{{ record.user_name || record.user_username || '-' }}</span>
                                <a-tag v-if="record.user_title" size="small" bordered>{{ record.user_title }}</a-tag>
                              </div>
                              <div class="user-sub">{{ record.store_name || '-' }} · {{ record.user_mobile || '-' }}</div>
                            </div>
                          </template>
                        </a-table-column>
                        <a-table-column title="次数" :width="90">
                          <template #cell="{ record }">{{ record._count || 1 }}</template>
                        </a-table-column>
                        <a-table-column title="页面/备注">
                          <template #cell="{ record }">
                            <div>{{ record._page || '-' }}</div>
                            <div class="records-sub" v-if="record._channel">渠道：{{ record._channel }}</div>
                          </template>
                        </a-table-column>
                      </template>
                    </a-table>
                  </a-spin>
                </a-tab-pane>

                <a-tab-pane key="showing">
                  <template #title>
                    <icon-user-group /> 带看记录
                  </template>

                  <div class="records-toolbar">
                    <a-space size="mini" wrap>
                      <a-button size="mini" type="primary" :loading="showingLogLoading" @click="reloadShowingLogs">
                        <template #icon><icon-refresh /></template>
                        刷新
                      </a-button>
                      <a-tag bordered size="small">共 {{ showingLogPagination.total || 0 }} 条</a-tag>
                    </a-space>
                    <div class="records-tip">说明：带看记录来源于小程序/业务端的行为日志。</div>
                  </div>

                  <a-spin :loading="showingLogLoading" style="width: 100%">
                    <a-table
                      row-key="id"
                      :data="showingLogList"
                      :pagination="showingLogPagination"
                      :bordered="false"
                      size="mini"
                      @page-change="handleShowingLogPageChange"
                      @page-size-change="handleShowingLogPageSizeChange"
                    >
                      <template #columns>
                        <a-table-column title="时间" data-index="createtime" :width="170" />
                        <a-table-column title="经纪人" :width="220">
                          <template #cell="{ record }">
                            <div class="user-cell">
                              <div class="user-main">
                                <span class="name">{{ record.user_name || record.user_username || '-' }}</span>
                                <a-tag v-if="record.user_title" size="small" bordered>{{ record.user_title }}</a-tag>
                              </div>
                              <div class="user-sub">{{ record.store_name || '-' }} · {{ record.user_mobile || '-' }}</div>
                            </div>
                          </template>
                        </a-table-column>
                        <a-table-column title="客户" :width="140">
                          <template #cell="{ record }">{{ record._client || '-' }}</template>
                        </a-table-column>
                        <a-table-column title="电话" :width="150">
                          <template #cell="{ record }">{{ record._phone || '-' }}</template>
                        </a-table-column>
                        <a-table-column title="页面/备注">
                          <template #cell="{ record }">{{ record._page || '-' }}</template>
                        </a-table-column>
                      </template>
                    </a-table>
                  </a-spin>
                </a-tab-pane>

                <a-tab-pane key="unlock">
                  <template #title>
                    <icon-unlock /> 开锁记录
                  </template>

                  <a-tabs v-model:active-key="unlockSubTabKey" type="card-gutter" size="small">
                    <a-tab-pane key="requests">
                      <template #title>申请记录</template>

                      <div class="records-toolbar">
                        <a-space size="mini" wrap>
                          <a-button size="mini" type="primary" :loading="unlockReqLoading" @click="reloadUnlockReqs">
                            <template #icon><icon-refresh /></template>
                            刷新
                          </a-button>
                          <a-tag bordered size="small">共 {{ unlockReqPagination.total || 0 }} 条</a-tag>
                        </a-space>
                        <div class="records-tip">说明：开锁申请记录来自 business_unlock_requests（蓝牙/密码）。如暂无申请记录，可切换到“过程日志”查看历史数据。</div>
                      </div>

                      <a-spin :loading="unlockReqLoading" style="width: 100%">
                        <a-table
                          v-if="unlockReqList && unlockReqList.length"
                          row-key="id"
                          :data="unlockReqList"
                          :pagination="unlockReqPagination"
                          :bordered="false"
                          size="mini"
                          @page-change="handleUnlockReqPageChange"
                          @page-size-change="handleUnlockReqPageSizeChange"
                        >
                          <template #columns>
                            <a-table-column title="申请时间" :width="170">
                              <template #cell="{ record }">{{ record.request_time || record.createtime || '-' }}</template>
                            </a-table-column>
                            <a-table-column title="经纪人" :width="220">
                              <template #cell="{ record }">
                                <div class="user-cell">
                                  <div class="user-main">
                                    <span class="name">{{ record.user_name || record.user_username || '-' }}</span>
                                    <a-tag v-if="record.user_title" size="small" bordered>{{ record.user_title }}</a-tag>
                                  </div>
                                  <div class="user-sub">{{ record.store_name || '-' }} · {{ record.user_mobile || '-' }}</div>
                                </div>
                              </template>
                            </a-table-column>
                            <a-table-column title="方式" :width="90">
                              <template #cell="{ record }">{{ getUnlockRequestTypeLabel(record.request_type) }}</template>
                            </a-table-column>
                            <a-table-column title="状态" :width="110">
                              <template #cell="{ record }">
                                <a-tag :color="getUnlockRequestStatusColor(record.request_status)" size="small">
                                  {{ getUnlockRequestStatusLabel(record.request_status) }}
                                </a-tag>
                              </template>
                            </a-table-column>
                            <a-table-column title="完成时间" :width="170">
                              <template #cell="{ record }">{{ getUnlockRequestFinishTime(record) }}</template>
                            </a-table-column>
                            <a-table-column title="备注/过期">
                              <template #cell="{ record }">
                                <div>{{ record.approval_remark || '-' }}</div>
                                <div class="records-sub" v-if="record.expires_at">过期：{{ record.expires_at }}</div>
                              </template>
                            </a-table-column>
                          </template>
                        </a-table>

                        <a-empty
                          v-else
                          description="暂无申请记录，可切换到“过程日志”查看历史记录"
                          style="padding: 20px 0"
                        >
                          <template #extra>
                            <a-button size="mini" type="primary" @click="switchUnlockSubTab('logs')">查看过程日志</a-button>
                          </template>
                        </a-empty>
                      </a-spin>
                    </a-tab-pane>

                    <a-tab-pane key="logs">
                      <template #title>过程日志</template>

                      <div class="records-toolbar">
                        <a-space size="mini" wrap>
                          <a-button size="mini" type="primary" :loading="unlockLogLoading" @click="reloadUnlockLogs">
                            <template #icon><icon-refresh /></template>
                            刷新
                          </a-button>
                          <a-tag bordered size="small">共 {{ unlockLogPagination.total || 0 }} 条</a-tag>
                        </a-space>
                        <div class="records-tip">说明：过程日志来自小程序开锁流程日志（business_user_activity_logs，activity_type=unlock）。</div>
                      </div>

                      <a-spin :loading="unlockLogLoading" style="width: 100%">
                        <a-table
                          row-key="id"
                          :data="unlockLogList"
                          :pagination="unlockLogPagination"
                          :bordered="false"
                          size="mini"
                          @page-change="handleUnlockLogPageChange"
                          @page-size-change="handleUnlockLogPageSizeChange"
                        >
                          <template #columns>
                            <a-table-column title="时间" data-index="createtime" :width="170" />
                            <a-table-column title="经纪人" :width="220">
                              <template #cell="{ record }">
                                <div class="user-cell">
                                  <div class="user-main">
                                    <span class="name">{{ record.user_name || record.user_username || '-' }}</span>
                                    <a-tag v-if="record.user_title" size="small" bordered>{{ record.user_title }}</a-tag>
                                  </div>
                                  <div class="user-sub">{{ record.store_name || '-' }} · {{ record.user_mobile || '-' }}</div>
                                </div>
                              </template>
                            </a-table-column>
                            <a-table-column title="阶段" :width="200">
                              <template #cell="{ record }">{{ getUnlockStageLabel(record._stage) }}</template>
                            </a-table-column>
                            <a-table-column title="结果" :width="90">
                              <template #cell="{ record }">
                                <a-tag v-if="record._success === true" color="green" size="small">成功</a-tag>
                                <a-tag v-else-if="record._success === false" color="red" size="small">失败</a-tag>
                                <span v-else>-</span>
                              </template>
                            </a-table-column>
                            <a-table-column title="页面/备注">
                              <template #cell="{ record }">
                                <div>{{ record._page || '-' }}</div>
                                <div class="records-sub" v-if="record._err_msg">{{ record._err_msg }}</div>
                              </template>
                            </a-table-column>
                          </template>
                        </a-table>
                      </a-spin>
                    </a-tab-pane>
                  </a-tabs>
                </a-tab-pane>
              </a-tabs>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </div>
  </BasicModal>

  <BindLockModal @register="registerBindLockModal" @success="reloadLockInfo" />

   <a-drawer v-model:visible="cloudVisible" width="400px" :footer="false" class="cloud-drawer">
      <template #title>云端设备详情</template>
     <a-spin :loading="cloudLoading" style="width: 100%">
       <div v-if="cloudDetail" class="cloud-detail-box">
         <div class="cd-hero">
            <div class="icon-circle"><icon-cloud /></div>
            <div class="cd-title">{{ cloudDetail.lockAlias || cloudDetail.lockName }}</div>
            <div class="cd-sub">{{ cloudDetail.lockMac }}</div>
         </div>
         <div class="cd-list">
            <div class="cd-item">
               <span>电量</span>
               <span class="cd-val">{{ cloudDetail.battery }}%</span>
            </div>
            <div class="cd-item">
               <span>设备ID</span>
               <span class="cd-val">{{ cloudDetail.lockId }}</span>
            </div>
            <div class="cd-item">
               <span>型号</span>
               <span class="cd-val">{{ cloudDetail.modelNum || '标准版' }}</span>
            </div>
         </div>
       </div>
     </a-spin>
   </a-drawer>

   <a-drawer v-model:visible="statusLogVisible" width="560px" :footer="false" class="statuslog-drawer">
     <template #title>状态变更记录</template>
     <a-spin :loading="statusLogLoading" style="width: 100%">
       <a-table
         row-key="id"
         :data="statusLogList"
         :pagination="statusLogPagination"
         :bordered="false"
         size="mini"
         @page-change="handleStatusLogPageChange"
         @page-size-change="handleStatusLogPageSizeChange"
       >
         <template #columns>
           <a-table-column title="时间" data-index="createtime" :width="160" />
           <a-table-column title="字段" :width="110">
             <template #cell="{ record }">{{ getStatusFieldLabel(record.field) }}</template>
           </a-table-column>
           <a-table-column title="变更" :width="180">
             <template #cell="{ record }">
               <span class="log-before">{{ record.before_value || '-' }}</span>
               <span class="log-arrow">→</span>
               <span class="log-after">{{ record.after_value || '-' }}</span>
             </template>
           </a-table-column>
           <a-table-column title="操作人" :width="120">
             <template #cell="{ record }">{{ record.user_name || record.user_username || '-' }}</template>
           </a-table-column>
           <a-table-column title="备注">
             <template #cell="{ record }">{{ record.remark || '-' }}</template>
           </a-table-column>
         </template>
       </a-table>
     </a-spin>
   </a-drawer>
 </template>

<script lang="ts">
import { defineComponent, h, ref, computed, watch } from 'vue';
import { BasicModal, useModal, useModalInner } from '/@/components/Modal';
import useLoading from '@/hooks/loading';
 import { Message, Modal } from '@arco-design/web-vue';
 import { getContent, getRenovation, getStatusLogs, getBehaviorLogs, getUnlockRequests } from './api';
 import { GetFullPath } from '@/utils/tool';
 import BindLockModal from '../ttlock/modal/BindLockModal.vue';
 import { getLockDetail, getPropertyLock, remoteUnlock, unbindProperty } from '../ttlock/api';
 import dayjs from 'dayjs';
 import { useClipboard } from '@vueuse/core';

export default defineComponent({
  name: 'HousesDetailView',
  components: { BasicModal, BindLockModal },
  setup() {
    const modelHeight = ref(900);
    const windHeight = ref(900);
    const detailData = ref<any>({});
    const renovationData = ref<any>({});
    const lockLoading = ref(false);
    const lockInfo = ref<any>(null);
    const cloudVisible = ref(false);
    const cloudLoading = ref(false);
    const cloudDetail = ref<any>(null);

    const statusLogVisible = ref(false);
    const statusLogLoading = ref(false);
    const statusLogList = ref<any[]>([]);
    const statusLogPagination = ref<any>({
      current: 1,
      pageSize: 10,
      total: 0,
      showTotal: true,
      showPageSize: true,
    });

    const mainTabKey = ref('1');
    const recordTabKey = ref<'view' | 'showing' | 'unlock'>('view');
    const unlockSubTabKey = ref<'requests' | 'logs'>('requests');

    const viewLogLoading = ref(false);
    const viewLogLoaded = ref(false);
    const viewLogList = ref<any[]>([]);
    const viewLogPagination = ref<any>({
      current: 1,
      pageSize: 10,
      total: 0,
      showTotal: true,
      showPageSize: true,
    });

    const showingLogLoading = ref(false);
    const showingLogLoaded = ref(false);
    const showingLogList = ref<any[]>([]);
    const showingLogPagination = ref<any>({
      current: 1,
      pageSize: 10,
      total: 0,
      showTotal: true,
      showPageSize: true,
    });

    const unlockReqLoading = ref(false);
    const unlockReqLoaded = ref(false);
    const unlockReqList = ref<any[]>([]);
    const unlockReqPagination = ref<any>({
      current: 1,
      pageSize: 10,
      total: 0,
      showTotal: true,
      showPageSize: true,
    });

    const unlockLogLoading = ref(false);
    const unlockLogLoaded = ref(false);
    const unlockLogList = ref<any[]>([]);
    const unlockLogPagination = ref<any>({
      current: 1,
      pageSize: 10,
      total: 0,
      showTotal: true,
      showPageSize: true,
    });
    const { loading, setLoading } = useLoading();
    const { copy } = useClipboard();
    
    // 注入 Modal 控制及其方法
    const [registerModal, { setModalProps, closeModal: closeRawModal }] = useModalInner(async (data) => {
      setModalProps({ confirmLoading: false });
      detailData.value = {};
      renovationData.value = {};
       lockInfo.value = null;
       mainTabKey.value = '1';
       recordTabKey.value = 'view';
       unlockSubTabKey.value = 'requests';
       resetBehaviorLogs();

      if (data?.record?.id) {
        setLoading(true);
        try {
          const detail = await getContent({ id: data.record.id });
          detailData.value = detail || {};
          
          try {
            const renovation = await getRenovation({ property_id: data.record.id });
            if (renovation && Object.keys(renovation).length > 0) {
              renovationData.value = renovation;
            }
          } catch (e) { /* ignore */ }
          
          await reloadLockInfo();
          await ensureBehaviorLoaded();
        } catch (e) {
          console.error(e);
        }
        setLoading(false);
      }
    });

    const closeModal = () => {
       closeRawModal();
    };

    const reloadLockInfo = async () => {
      if (!detailData.value?.id) return;
      lockLoading.value = true;
      try {
        const resp: any = await getPropertyLock({ property_id: Number(detailData.value.id) });
        lockInfo.value = resp?.data ?? resp ?? null;
      } catch {
        lockInfo.value = null;
      } finally {
        lockLoading.value = false;
      }
    };

    const [registerBindLockModal, { openModal: openBindLockModal }] = useModal();
    const openBindLock = () => {
      openBindLockModal(true, { property: { id: Number(detailData.value.id), title: detailData.value.title } });
    };

    const handleUnbind = async () => {
      if (!detailData.value?.id) return;
      await unbindProperty({ property_id: Number(detailData.value.id) });
      Message.success('解绑成功');
      detailData.value.has_smart_lock = 0;
      await reloadLockInfo();
    };

    const handleRemoteUnlock = async () => {
      if (!detailData.value?.id) return;
      try {
        Message.loading({ content: '获取开锁密码中...', id: 'ttlock-passcode', duration: 0 });
        const resp: any = await remoteUnlock({ property_id: Number(detailData.value.id) });
        const data = (resp?.data ?? resp) || null;
        const pwd = String(data?.keyboardPwd || data?.keyboard_pwd || '').trim();
        if (!pwd) {
          Message.error({ content: '获取密码失败：响应缺少密码字段', id: 'ttlock-passcode', duration: 2000 });
          return;
        }
        const startMs = Number(data?.startDate || data?.start_date || 0);
        const endMs = Number(data?.endDate || data?.end_date || 0);
        const startText = startMs > 0 ? dayjs(startMs).format('YYYY-MM-DD HH:mm') : '';
        const endText = endMs > 0 ? dayjs(endMs).format('YYYY-MM-DD HH:mm') : '';
        copy(pwd);
        Message.success({ content: '开锁密码已生成（已复制）', id: 'ttlock-passcode', duration: 2000 });

        Modal.success({
          title: '开锁密码',
          titleAlign: 'start',
          hideCancel: false,
          okText: '再次复制',
          cancelText: '关闭',
          content: () => h('div', { style: { paddingTop: '4px' } }, [
            h(
              'div',
              {
                style: {
                  fontSize: '26px',
                  fontWeight: '800',
                  letterSpacing: '0.12em',
                  fontFamily: 'monospace',
                  padding: '10px 12px',
                  border: '1px solid var(--color-border-2)',
                  borderRadius: '10px',
                  background: 'var(--color-fill-1)',
                  color: 'var(--color-text-1)',
                  display: 'inline-block',
                },
              },
              pwd,
            ),
            h(
              'div',
              {
                style: {
                  marginTop: '10px',
                  color: 'var(--color-text-3)',
                  fontSize: '12px',
                  lineHeight: '18px',
                },
              },
              startText || endText
                ? `生效：${startText || '-'} · 失效：${endText || '-'}（本地时间）`
                : '一次性密码：6 小时内可使用 1 次（已复制到剪贴板）',
            ),
          ]),
          onOk: () => {
            copy(pwd);
            Message.success('已复制');
          },
        });
      } catch (e: any) {
        Message.error({ content: e?.message || '获取密码失败', id: 'ttlock-passcode', duration: 2000 });
      }
    };

    const openCloudDetail = async () => {
      if (!detailData.value?.id) return;
      cloudVisible.value = true;
      cloudLoading.value = true;
      cloudDetail.value = null;
      try {
        const resp: any = await getLockDetail({ property_id: Number(detailData.value.id) });
        cloudDetail.value = resp?.data ?? resp ?? null;
      } finally {
        cloudLoading.value = false;
      }
    };

     const formatTime = (sec: any) => {
        if(!sec) return '-';
        return dayjs(Number(sec) * 1000).format('MM-DD HH:mm');
      }
 
     const formatMoney = (v: any) => {
       if (v === null || v === undefined || v === '') return '-';
       const n = Number(v);
      if (Number.isNaN(n)) return '-';
      return `¥${n.toFixed(2)}`;
    };

    const getStatusFieldLabel = (f: string) => {
      const map: any = {
        status: '上架状态',
        sale_status: '销售状态',
        hot_status: '推荐状态',
        weigh: '排序权重',
      };
      return map[f] || f || '-';
    };

    const fetchStatusLogs = async () => {
      if (!detailData.value?.id) return;
      statusLogLoading.value = true;
      try {
        const resp: any = await getStatusLogs({
          property_id: Number(detailData.value.id),
          page: statusLogPagination.value.current,
          pageSize: statusLogPagination.value.pageSize,
        });
        const data = resp?.items ? resp : (resp?.data?.items ? resp.data : resp?.data ?? resp);
        statusLogList.value = (data?.items || []) as any[];
        statusLogPagination.value.current = Number(data?.page || statusLogPagination.value.current);
        statusLogPagination.value.total = Number(data?.total || 0);
      } catch (e: any) {
        statusLogList.value = [];
        statusLogPagination.value.total = 0;
        Message.error(e?.message || '加载状态记录失败');
      } finally {
        statusLogLoading.value = false;
      }
    };

    const openStatusLogs = async () => {
      statusLogVisible.value = true;
      statusLogPagination.value.current = 1;
      await fetchStatusLogs();
    };

    const handleStatusLogPageChange = (p: number) => {
      statusLogPagination.value.current = p;
      fetchStatusLogs();
    };

    const handleStatusLogPageSizeChange = (ps: number) => {
      statusLogPagination.value.pageSize = ps;
      statusLogPagination.value.current = 1;
      fetchStatusLogs();
    };

    const safeParseMeta = (raw: any) => {
      const s = (raw ?? '').toString().trim();
      if (!s) return {};
      try {
        const out = JSON.parse(s);
        return out && typeof out === 'object' ? out : {};
      } catch (e) {
        return {};
      }
    };

    const metaGetValue = (meta: any, key: string) => {
      if (!meta || typeof meta !== 'object') return undefined;
      if (Object.prototype.hasOwnProperty.call(meta, key)) return meta[key];
      const nested = (meta as any).meta;
      if (!nested) return undefined;
      if (typeof nested === 'object' && Object.prototype.hasOwnProperty.call(nested, key)) return nested[key];
      if (typeof nested === 'string') {
        const parsed = safeParseMeta(nested);
        if (parsed && typeof parsed === 'object' && Object.prototype.hasOwnProperty.call(parsed, key)) return (parsed as any)[key];
      }
      return undefined;
    };

    const metaString = (meta: any, ...keys: string[]) => {
      for (const k of keys) {
        const v = metaGetValue(meta, k);
        const s = (v ?? '').toString().trim();
        if (s) return s;
      }
      return '';
    };

    const metaInt = (meta: any, ...keys: string[]) => {
      for (const k of keys) {
        const v = metaGetValue(meta, k);
        const n = Number(v);
        if (Number.isFinite(n) && n > 0) return Math.trunc(n);
      }
      return 0;
    };

    const decorateActivityRows = (rows: any[], recordType: 'view' | 'showing') => {
      for (const r of rows) {
        if (!r || typeof r !== 'object') continue;
        const meta = safeParseMeta(r?.meta_data);
        r._page = metaString(meta, 'page') || metaString(meta, 'path', 'route', 'url');
        r._channel = metaString(meta, 'channel', 'source', 'from', 'scene');
        if (recordType === 'view') {
          const cnt = metaInt(meta, 'count', 'view_count', 'times');
          r._count = cnt > 0 ? cnt : 1;
        } else {
          r._client = metaString(meta, 'client_name', 'client', 'customer', 'name');
          r._phone = metaString(meta, 'client_phone', 'phone', 'mobile');
        }
      }
    };

    const getUnlockStageLabel = (raw: any) => {
      const s = String(raw || '').trim();
      const map: any = {
        ble_init_start: '蓝牙初始化',
        ble_init_ok: '蓝牙初始化成功',
        ble_init_fail: '蓝牙初始化失败',
        ble_unlock_start: '蓝牙开锁开始',
        ble_unlock_finish: '蓝牙开锁结束',
      };
      return map[s] || s || '-';
    };

    const decorateUnlockRows = (rows: any[]) => {
      for (const r of rows) {
        if (!r || typeof r !== 'object') continue;
        const meta = safeParseMeta(r?.meta_data);
        r._page = metaString(meta, 'page') || metaString(meta, 'path', 'route', 'url');
        r._stage = metaString(meta, 'stage');

        const successRaw = metaGetValue(meta, 'success');
        if (successRaw === true || successRaw === 1 || successRaw === '1' || successRaw === 'true') {
          r._success = true;
        } else if (successRaw === false || successRaw === 0 || successRaw === '0' || successRaw === 'false') {
          r._success = false;
        } else if (String(r._stage || '').includes('fail')) {
          r._success = false;
        } else {
          r._success = undefined;
        }

        r._request_id = metaInt(meta, 'request_id');
        r._err_msg = metaString(meta, 'err_msg', 'error', 'message');
      }
    };

    const getUnlockRequestTypeLabel = (raw: any) => {
      const s = String(raw || '').trim();
      const map: any = { bluetooth: '蓝牙', password: '密码' };
      return map[s] || s || '-';
    };

    const getUnlockRequestStatusLabel = (raw: any) => {
      const s = String(raw || '').trim();
      const map: any = {
        pending: '待审核',
        approved: '已通过',
        rejected: '已拒绝',
        completed: '已完成',
        cancelled: '已取消',
      };
      return map[s] || s || '-';
    };

    const getUnlockRequestStatusColor = (raw: any) => {
      const s = String(raw || '').trim();
      const map: any = {
        pending: 'orange',
        approved: 'green',
        completed: 'green',
        rejected: 'red',
        cancelled: 'gray',
      };
      return map[s] || '';
    };

    const getUnlockRequestFinishTime = (record: any) => {
      const status = String(record?.request_status || '').trim();
      if (!status || status === 'pending') return '-';
      return record?.updatetime || '-';
    };

    const normalizeListResp = (resp: any) => {
      if (!resp) return resp;
      // defHttp 默认会把 {code,data,message} 解包成 data；但为了兼容部分场景，这里做一次兜底
      if (typeof resp === 'object' && resp !== null && 'data' in resp && (resp as any).data !== undefined) {
        return (resp as any).data ?? resp;
      }
      return resp;
    };

    const normalizeListItems = (raw: any) => {
      if (Array.isArray(raw)) return raw;
      if (!raw || typeof raw !== 'object') return [];
      if (Array.isArray((raw as any).items)) return (raw as any).items;
      if (Array.isArray((raw as any).list)) return (raw as any).list;
      if (Array.isArray((raw as any).rows)) return (raw as any).rows;
      if (Array.isArray((raw as any).data)) return (raw as any).data;
      return [];
    };

    const normalizeListTotal = (raw: any, items: any[]) => {
      const n = Number((raw as any)?.total ?? (raw as any)?.count ?? (raw as any)?.totalCount ?? (raw as any)?.total_count ?? 0);
      if (Number.isFinite(n) && n >= 0) return Math.trunc(n);
      return items.length;
    };

    const normalizeListPage = (raw: any, fallback: number) => {
      const n = Number((raw as any)?.page ?? (raw as any)?.current ?? (raw as any)?.pageNo ?? (raw as any)?.page_no ?? fallback);
      if (Number.isFinite(n) && n > 0) return Math.trunc(n);
      return fallback;
    };

    const normalizeListPageSize = (raw: any, fallback: number) => {
      const n = Number((raw as any)?.pageSize ?? (raw as any)?.page_size ?? (raw as any)?.limit ?? fallback);
      if (Number.isFinite(n) && n > 0) return Math.trunc(n);
      return fallback;
    };

    const fetchBehaviorLogs = async (
      recordType: 'view' | 'showing' | 'unlock',
      listRef: any,
      paginationRef: any,
      loadingRef: any,
      loadedRef: any,
    ) => {
      if (!detailData.value?.id) return;
      if (loadingRef.value) return;
      loadingRef.value = true;
      try {
        const resp: any = await getBehaviorLogs({
          property_id: Number(detailData.value.id),
          record_type: recordType,
          page: paginationRef.value.current,
          pageSize: paginationRef.value.pageSize,
        });
        const data = normalizeListResp(resp);
        const rows = (normalizeListItems(data) as any[]).filter((it) => it && typeof it === 'object') as any[];
        if (recordType === 'view' || recordType === 'showing') {
          decorateActivityRows(rows, recordType);
        } else if (recordType === 'unlock') {
          decorateUnlockRows(rows);
        }
        listRef.value = rows;
        paginationRef.value.current = normalizeListPage(data, paginationRef.value.current);
        paginationRef.value.pageSize = normalizeListPageSize(data, paginationRef.value.pageSize);
        paginationRef.value.total = normalizeListTotal(data, rows);
        loadedRef.value = true;
      } catch (e: any) {
        listRef.value = [];
        paginationRef.value.total = 0;
        loadedRef.value = false;
        const tip: any = { view: '加载浏览记录失败', showing: '加载带看记录失败', unlock: '加载开锁记录失败' };
        Message.error(e?.message || tip[recordType] || '加载记录失败');
      } finally {
        loadingRef.value = false;
      }
    };

    const fetchUnlockReqs = async (listRef: any, paginationRef: any, loadingRef: any, loadedRef: any) => {
      if (!detailData.value?.id) return;
      if (loadingRef.value) return;
      loadingRef.value = true;
      try {
        const resp: any = await getUnlockRequests({
          property_id: Number(detailData.value.id),
          page: paginationRef.value.current,
          pageSize: paginationRef.value.pageSize,
        });
        const data = normalizeListResp(resp);
        const rows = (normalizeListItems(data) as any[]).filter((it) => it && typeof it === 'object') as any[];
        listRef.value = rows;
        paginationRef.value.current = normalizeListPage(data, paginationRef.value.current);
        paginationRef.value.pageSize = normalizeListPageSize(data, paginationRef.value.pageSize);
        paginationRef.value.total = normalizeListTotal(data, rows);
        loadedRef.value = true;
      } catch (e: any) {
        listRef.value = [];
        paginationRef.value.total = 0;
        loadedRef.value = false;
        Message.error(e?.message || '加载开锁申请记录失败');
      } finally {
        loadingRef.value = false;
      }
    };

    const resetBehaviorLogs = () => {
      viewLogLoaded.value = false;
      viewLogLoading.value = false;
      viewLogList.value = [];
      viewLogPagination.value.current = 1;
      viewLogPagination.value.total = 0;

      showingLogLoaded.value = false;
      showingLogLoading.value = false;
      showingLogList.value = [];
      showingLogPagination.value.current = 1;
      showingLogPagination.value.total = 0;

      unlockReqLoaded.value = false;
      unlockReqLoading.value = false;
      unlockReqList.value = [];
      unlockReqPagination.value.current = 1;
      unlockReqPagination.value.total = 0;

      unlockLogLoaded.value = false;
      unlockLogLoading.value = false;
      unlockLogList.value = [];
      unlockLogPagination.value.current = 1;
      unlockLogPagination.value.total = 0;
    };

    const ensureBehaviorLoaded = async () => {
      if (String(mainTabKey.value) !== '5') return;
      switch (recordTabKey.value) {
        case 'view':
          if (!viewLogLoaded.value) {
            viewLogPagination.value.current = 1;
            await fetchBehaviorLogs('view', viewLogList, viewLogPagination, viewLogLoading, viewLogLoaded);
          }
          break;
        case 'showing':
          if (!showingLogLoaded.value) {
            showingLogPagination.value.current = 1;
            await fetchBehaviorLogs('showing', showingLogList, showingLogPagination, showingLogLoading, showingLogLoaded);
          }
          break;
        case 'unlock':
          if (unlockSubTabKey.value === 'requests') {
            if (!unlockReqLoaded.value) {
              unlockReqPagination.value.current = 1;
              await fetchUnlockReqs(unlockReqList, unlockReqPagination, unlockReqLoading, unlockReqLoaded);
            }
          } else {
            if (!unlockLogLoaded.value) {
              unlockLogPagination.value.current = 1;
              await fetchBehaviorLogs('unlock', unlockLogList, unlockLogPagination, unlockLogLoading, unlockLogLoaded);
            }
          }
          break;
      }
    };

    watch(mainTabKey, () => {
      ensureBehaviorLoaded();
    });
    watch(recordTabKey, () => {
      ensureBehaviorLoaded();
    });
    watch(unlockSubTabKey, () => {
      ensureBehaviorLoaded();
    });

    const switchUnlockSubTab = (key: 'requests' | 'logs') => {
      unlockSubTabKey.value = key;
    };

    const reloadViewLogs = async () => {
      await fetchBehaviorLogs('view', viewLogList, viewLogPagination, viewLogLoading, viewLogLoaded);
    };
    const handleViewLogPageChange = (p: number) => {
      viewLogPagination.value.current = p;
      reloadViewLogs();
    };
    const handleViewLogPageSizeChange = (ps: number) => {
      viewLogPagination.value.pageSize = ps;
      viewLogPagination.value.current = 1;
      reloadViewLogs();
    };

    const reloadShowingLogs = async () => {
      await fetchBehaviorLogs('showing', showingLogList, showingLogPagination, showingLogLoading, showingLogLoaded);
    };
    const handleShowingLogPageChange = (p: number) => {
      showingLogPagination.value.current = p;
      reloadShowingLogs();
    };
    const handleShowingLogPageSizeChange = (ps: number) => {
      showingLogPagination.value.pageSize = ps;
      showingLogPagination.value.current = 1;
      reloadShowingLogs();
    };

    const reloadUnlockReqs = async () => {
      await fetchUnlockReqs(unlockReqList, unlockReqPagination, unlockReqLoading, unlockReqLoaded);
    };
    const handleUnlockReqPageChange = (p: number) => {
      unlockReqPagination.value.current = p;
      reloadUnlockReqs();
    };
    const handleUnlockReqPageSizeChange = (ps: number) => {
      unlockReqPagination.value.pageSize = ps;
      unlockReqPagination.value.current = 1;
      reloadUnlockReqs();
    };

    const reloadUnlockLogs = async () => {
      await fetchBehaviorLogs('unlock', unlockLogList, unlockLogPagination, unlockLogLoading, unlockLogLoaded);
    };
    const handleUnlockLogPageChange = (p: number) => {
      unlockLogPagination.value.current = p;
      reloadUnlockLogs();
    };
    const handleUnlockLogPageSizeChange = (ps: number) => {
      unlockLogPagination.value.pageSize = ps;
      unlockLogPagination.value.current = 1;
      reloadUnlockLogs();
    };

    const getImageUrl = (url: string) => {
      if (!url) return 'https://p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/a8c8cdb109cb051163646151a4a5083b.png~tplv-uwbnlip3yd-webp.webp';
      return GetFullPath(url) || url;
    };

    const parseTags = (v: any) => {
      if (!v) return [];
      if (Array.isArray(v)) return v.map((it) => String(it).trim()).filter(Boolean);
      if (typeof v === 'string') return v.split(',').map((it) => it.trim()).filter(Boolean);
      return [];
    };

    const tagList = computed(() => parseTags(detailData.value?.tags));
    const customDescText = computed(() => (detailData.value?.custom_desc ?? '').toString().trim());

    const districtText = computed(() => {
      const d = (detailData.value?.district || '').toString().trim();
      if (d) return d;
      const raw = (detailData.value?.address || '').toString().trim();
      const parts = raw.split(/\\s+/).filter(Boolean);
      if (parts.length >= 3) return parts.slice(0, 3).join(' / ');
      return '';
    });

    const addressDetailText = computed(() => {
      const d = (detailData.value?.address_detail || '').toString().trim();
      if (d) return d;
      const raw = (detailData.value?.address || '').toString().trim();
      if (!raw) return '';
      const parts = raw.split(/\\s+/).filter(Boolean);
      if (parts.length >= 4) return parts.slice(3).join(' ');
      return raw;
    });

    const allowImageDownload = computed(() => {
      const v = detailData.value?.allow_image_download;
      if (v === null || v === undefined || v === '') return true;
      return Number(v) === 1;
    });

    const allowVideoDownload = computed(() => {
      const v = detailData.value?.allow_video_download;
      if (v === null || v === undefined || v === '') return true;
      return Number(v) === 1;
    });

    const galleryImages = computed(() => {
      const list: string[] = [];
      const cover = (detailData.value?.cover_image || '').toString().trim();
      const imgs = parseTags(detailData.value?.images);
      if (cover) list.push(cover);
      list.push(...imgs);
      return [...new Set(list)].filter(Boolean);
    });

    const videoSrc = computed(() => {
      const v = (detailData.value?.video_url || '').toString().trim();
      if (!v) return '';
      return getImageUrl(v);
    });

    const renovationImages = computed(() => {
      const list = parseTags(renovationData.value?.images);
      return [...new Set(list.map((it) => getImageUrl(String(it))))].filter(Boolean);
    });

    type StageStatus = 'todo' | 'doing' | 'done';
    type StageTimelineItem = {
      stage: string;
      status: StageStatus;
      date: string;
      note: string;
      images: string[];
    };

    const normalizeStageStatus = (v: any): StageStatus => {
      const s = String(v ?? '').trim().toLowerCase();
      if (s === 'done' || s === 'finished' || s === 'completed') return 'done';
      if (s === 'doing' || s === 'in_progress' || s === 'progress') return 'doing';
      return 'todo';
    };

    const parseStageLogs = (raw: any): StageTimelineItem[] => {
      let arr: any = raw;
      if (typeof raw === 'string') {
        const s = raw.trim();
        if (s && s.startsWith('[')) {
          try {
            arr = JSON.parse(s);
          } catch {
            arr = [];
          }
        } else {
          arr = [];
        }
      }
      if (!Array.isArray(arr)) return [];
      return arr
        .map((it: any) => {
          const stage = String(it?.stage ?? it?.stage_name ?? it?.name ?? '').trim();
          if (!stage) return null;
          const images = parseTags(it?.images).map((u) => getImageUrl(String(u)));
          return {
            stage,
            status: normalizeStageStatus(it?.status),
            date: String(it?.date ?? '').trim(),
            note: String(it?.note ?? it?.notes ?? '').trim(),
            images: [...new Set(images)].filter(Boolean),
          } as StageTimelineItem;
        })
        .filter(Boolean) as StageTimelineItem[];
    };

    const fallbackStageOrder = ['设计', '拆改', '水电', '泥瓦', '木工', '油漆', '安装', '软装', '验收'];

    const renovationStageTimeline = computed((): StageTimelineItem[] => {
      const r: any = renovationData.value || {};
      const overall = String(r.renovation_status || '').trim();
      const currentStage = String(r.current_stage || '').trim();
      const logs = parseStageLogs(r.stage_logs);
      const order = logs.length ? logs.map((x) => x.stage) : fallbackStageOrder;
      const currentIdx = currentStage ? order.findIndex((s) => s === currentStage) : -1;

      return order.map((stage, idx) => {
        const found = logs.find((x) => x.stage === stage);
        let status: StageStatus = found?.status || 'todo';
        let date = found?.date || '';
        let note = found?.note || '';
        let images = found?.images || [];

        if (!found) {
          if (overall === 'done') status = 'done';
          else if (overall === 'none') status = 'todo';
          else if (overall === 'in_progress' && currentIdx >= 0) {
            if (idx < currentIdx) status = 'done';
            else if (idx === currentIdx) status = 'doing';
            else status = 'todo';
          }
        } else {
          // 总状态兜底覆盖（避免出现“已完工但某阶段仍未开始”）
          if (overall === 'done') status = 'done';
          if (overall === 'none') status = 'todo';
        }

        return { stage, status, date, note, images };
      });
    });

    const getStageStatusText = (s: StageStatus) => {
      if (s === 'done') return '已完成';
      if (s === 'doing') return '进行中';
      return '未开始';
    };

    const getStageStatusColor = (s: StageStatus) => {
      if (s === 'done') return 'green';
      if (s === 'doing') return 'arcoblue';
      return 'gray';
    };

    const downloadUrl = async (url: string, filename: string) => {
      const href = (url || '').toString().trim();
      if (!href) return;
      try {
        const resp = await fetch(href, { credentials: 'include' });
        if (!resp.ok) throw new Error('下载失败');
        const blob = await resp.blob();
        const blobUrl = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = blobUrl;
        a.download = filename || 'download';
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.setTimeout(() => window.URL.revokeObjectURL(blobUrl), 800);
      } catch (e) {
        window.open(href, '_blank');
      }
    };

    const priceAmountText = computed(() => {
      const v = detailData.value?.price;
      if (v === null || v === undefined || v === '') return '-';
      const n = Number(v);
      if (!Number.isFinite(n)) return '-';
      return n.toFixed(2).replace(/\\.00$/, '');
    });

    const pricePerSquareText = computed(() => {
      const area = Number(detailData.value?.area);
      if (!Number.isFinite(area) || area <= 0) return '';
      const price = Number(detailData.value?.price);
      if (!Number.isFinite(price) || price <= 0) return '';
      const unit = (detailData.value?.price_unit || '万').toString();
      const yuan = price * (unit === '元' ? 1 : 10000);
      const per = Math.round(yuan / area);
      if (!Number.isFinite(per) || per <= 0) return '';
      return `≈${per}元/㎡`;
    });

    const getSaleStatusLabel = (s: string) => {
       const map:any = { on_sale: '在售', in_sale: '预售', sold: '已售', off_market: '下架' };
       return map[s] || s;
    };
    
    const getRenovationStatusLabel = (s:string) => {
       const map:any = { none: '未装修', in_progress: '施工中', done: '已完工' };
       return map[s] || s;
    };

    const getBatteryClass = (bat: any) => {
       const b = Number(bat);
       if(b <= 20) return 'color-danger';
       if(b <= 50) return 'color-warning';
       return 'color-success';
    };

    const onHeightChange = (val: any) => {
      windHeight.value = val;
      if (typeof val === 'number' && val > 0) modelHeight.value = val;
    };

    return {
      registerModal, loading, detailData, renovationData, modelHeight, onHeightChange,
      windHeight, getImageUrl, parseTags, tagList, customDescText, galleryImages, allowImageDownload, allowVideoDownload, videoSrc, downloadUrl,
      priceAmountText, pricePerSquareText, getSaleStatusLabel,
      getRenovationStatusLabel, getBatteryClass, lockLoading, lockInfo, openBindLock,
      handleUnbind, handleRemoteUnlock, openCloudDetail, reloadLockInfo, formatTime,
      registerBindLockModal, cloudVisible, cloudLoading, cloudDetail, closeModal,
      formatMoney,
      statusLogVisible, statusLogLoading, statusLogList, statusLogPagination,
      openStatusLogs, handleStatusLogPageChange, handleStatusLogPageSizeChange, getStatusFieldLabel,
      districtText, addressDetailText, renovationImages,
      renovationStageTimeline, getStageStatusText, getStageStatusColor,

      mainTabKey, recordTabKey, unlockSubTabKey, switchUnlockSubTab,
      viewLogLoading, viewLogList, viewLogPagination, reloadViewLogs, handleViewLogPageChange, handleViewLogPageSizeChange,
      showingLogLoading, showingLogList, showingLogPagination, reloadShowingLogs, handleShowingLogPageChange, handleShowingLogPageSizeChange,
      unlockReqLoading, unlockReqList, unlockReqPagination, reloadUnlockReqs, handleUnlockReqPageChange, handleUnlockReqPageSizeChange,
      unlockLogLoading, unlockLogList, unlockLogPagination, reloadUnlockLogs, handleUnlockLogPageChange, handleUnlockLogPageSizeChange,
      getUnlockStageLabel,
      getUnlockRequestTypeLabel, getUnlockRequestStatusLabel, getUnlockRequestStatusColor, getUnlockRequestFinishTime,
    };
  },
});
</script>

<style lang="less" scoped>
/* Modal and Container Styling */
.house-detail-modal {
  :deep(.arco-modal-body) {
    padding: 0;
  }
}

.detail-container {
  background: var(--color-bg-1);
  border-radius: 8px;
  overflow: hidden;
}

.gallery-wrapper {
  .grid-item {
    position: relative;
    overflow: hidden;
    border-radius: 10px;
  }

  .img-actions {
    position: absolute;
    top: 8px;
    right: 8px;
    z-index: 2;
  }

  .download-btn {
    box-shadow: 0 8px 22px rgba(0, 0, 0, 0.12);
  }
}

.renovation-media {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px dashed var(--color-border-2);
  .media-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;
    .media-title {
      font-weight: 600;
      color: var(--color-text-1);
    }
  }
  .renovation-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 10px;
    @media (max-width: 1400px) {
      grid-template-columns: repeat(3, 1fr);
    }
    @media (max-width: 1200px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }
  .renovation-item {
    position: relative;
    width: 100%;
    height: 120px;
    border-radius: 10px;
    overflow: hidden;
    background: var(--color-fill-2);
    border: 1px solid var(--color-border-2);
    transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
    &:hover {
      transform: translateY(-1px);
      border-color: rgb(var(--primary-6));
      box-shadow: 0 10px 24px rgba(0, 0, 0, 0.08);
    }
    :deep(.arco-image) {
      width: 100%;
      height: 100%;
    }

    .img-actions {
      position: absolute;
      top: 8px;
      right: 8px;
      z-index: 2;
    }
    .download-btn {
      box-shadow: 0 8px 22px rgba(0, 0, 0, 0.12);
    }
  }
}

.video-wrapper {
  padding: 16px 20px 20px;
}

.video-card {
  border: 1px solid var(--color-border-2);
  background: var(--color-bg-2);
  border-radius: 12px;
  padding: 12px;
}

.video-player {
  width: 100%;
  height: 420px;
  border-radius: 10px;
  background: #000;
  object-fit: contain;
}

.video-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.records-wrapper {
  padding: 16px 20px 20px;
}

.records-toolbar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}

.records-tip {
  font-size: 12px;
  color: var(--color-text-3);
  line-height: 20px;
  padding-top: 2px;
}

.records-sub {
  margin-top: 2px;
  font-size: 12px;
  color: var(--color-text-3);
  line-height: 18px;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-cell {
  .user-main {
    display: flex;
    align-items: center;
    gap: 6px;
    min-width: 0;
  }
  .name {
    font-weight: 600;
    color: var(--color-text-1);
    max-width: 160px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .user-sub {
    margin-top: 2px;
    font-size: 12px;
    color: var(--color-text-3);
    max-width: 260px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

/* Header Section */
.detail-header {
  padding: 12px 20px 0;
  background: var(--color-bg-1);

  .header-card {
    border-radius: 12px;
    border: 1px solid var(--color-border-2);
    background: var(--color-bg-2);

    :deep(.arco-card-body) {
      padding: 12px 14px;
    }
  }

  .header-inner {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;
  }

  .header-left {
    display: flex;
    align-items: flex-start;
    gap: 16px;
    min-width: 0;
  }

  .cover-preview {
    flex-shrink: 0;

    .preview-img {
      border-radius: 10px;
      border: 1px solid var(--color-border-2);
      box-shadow: 0 6px 18px rgba(0, 0, 0, 0.06);
    }
  }

  .header-info {
    min-width: 0;

    .title-row {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 6px;
    }

    .detail-title {
      font-size: 16px;
      font-weight: 700;
      color: var(--color-text-1);
      line-height: 1.35;
      margin: 0;
      max-width: 520px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .detail-meta {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 6px 12px;
      font-size: 13px;
      color: var(--color-text-2);
      margin-bottom: 6px;

      .meta-item {
        display: inline-flex;
        align-items: center;
        gap: 4px;
      }

      .meta-item.addr {
        max-width: 420px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .price-display {
      display: flex;
      align-items: baseline;
      gap: 4px;

      .price-amount {
        font-size: 22px;
        font-weight: 800;
        color: rgb(var(--primary-6));
      }

      .price-unit {
        font-size: 13px;
        font-weight: 600;
        color: var(--color-text-2);
        margin-right: 10px;
      }

      .price-per {
        font-size: 12px;
        color: var(--color-text-3);
      }
    }
  }

  .close-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s ease;

    &:hover {
      background: var(--color-fill-2);
      transform: rotate(90deg);
    }
  }
}

.detail-body {
  padding: 14px 20px 16px;
  
  :deep(.arco-tabs-nav) {
    margin-bottom: 12px;
  }
  
  .tab-content {
    min-height: 0;
  }
}

/* Detail Grid Layout */
.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  
  .detail-card {
    background: var(--color-bg-2);
    border: 1px solid var(--color-border-2);
    border-radius: 12px;
    padding: 20px;
    transition: all 0.2s ease;
    
    &:hover {
      border-color: var(--color-border-3);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
    }
    
    .card-header {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 16px;
      
      .card-icon {
        color: rgb(var(--primary-6));
        font-size: 16px;
      }
      
      .card-title {
        font-size: 15px;
        font-weight: 600;
        color: var(--color-text-1);
      }

      .card-actions {
        margin-left: auto;
        display: inline-flex;
        align-items: center;
        gap: 8px;

        :deep(.arco-btn) {
          padding-left: 8px;
          padding-right: 8px;
        }
      }
    }
    
    .card-body {
      .info-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 16px 20px;
        
        .info-item {
          &.full-width {
            grid-column: 1 / -1;
          }
          
          label {
            font-size: 12px;
            color: var(--color-text-3);
            margin-bottom: 4px;
            display: block;
          }
          
          .value {
            font-size: 14px;
            color: var(--color-text-1);
            font-weight: 500;
            
            &.phone-value {
              font-family: monospace;
            }
            
            &.highlight {
              color: rgb(var(--primary-6));
              font-weight: 600;
            }
          }
        }
      }
      
      .stats-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 16px;
        
        .stat-item {
          text-align: center;
          padding: 12px;
          background: var(--color-fill-1);
          border-radius: 8px;
          
          .stat-value {
            font-size: 18px;
            font-weight: 700;
            color: var(--color-text-1);
            margin-bottom: 4px;
          }
          
          .stat-label {
            font-size: 12px;
            color: var(--color-text-3);
          }
        }
      }
    }
  }
  
  /* Smart Lock Widget - Full Width */
  .lock-widget {
    grid-column: 1 / -1;
    
    &.active {
      background: linear-gradient(135deg, var(--color-bg-2) 0%, var(--color-fill-1) 100%);
      border-color: rgba(var(--primary-6), 0.3);
    }
    
    .widget-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      
      .title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 15px;
        font-weight: 600;
        color: var(--color-text-1);
      }
      
      .status {
        .battery-indicator {
          display: flex;
          align-items: center;
          gap: 4px;
          font-size: 13px;
          font-weight: 500;
          
          &.color-success { color: rgb(var(--success-6)); }
          &.color-warning { color: rgb(var(--warning-6)); }
          &.color-danger { color: rgb(var(--danger-6)); }
        }
      }
      
      .btn-bind {
        border-radius: 6px;
      }
    }
    
    .widget-body {
      display: flex;
      align-items: center;
      gap: 14px;
      
      .device-info {
        flex: 1;
        display: flex;
        align-items: center;
        gap: 16px;
        
        .device-icon {
          width: 40px;
          height: 40px;
          background: var(--color-bg-1);
          border: 1px solid var(--color-border-2);
          border-radius: 10px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 20px;
          color: rgb(var(--primary-6));
        }
        
        .device-text {
          .d-name {
            font-weight: 600;
            font-size: 14px;
            color: var(--color-text-1);
            margin-bottom: 2px;
          }
          
          .d-mac {
            font-family: monospace;
            color: var(--color-text-3);
            font-size: 12px;
          }
        }
      }
      
      .control-grid {
        display: flex;
        gap: 10px;
        
        .control-btn {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          min-width: 88px;
          height: 36px;
          border-radius: 10px;
          background: var(--color-bg-1);
          border: 1px solid var(--color-border-2);
          cursor: pointer;
          transition: all 0.2s ease;
          font-size: 12px;
          color: var(--color-text-2);
          
          &:hover {
            border-color: var(--color-border-3);
            box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
          }
          
          &.primary {
            background: rgb(var(--primary-6));
            color: #fff;
            border: none;
            
            &:hover {
              background: rgb(var(--primary-5));
              box-shadow: 0 4px 12px rgba(var(--primary-6), 0.3);
            }
          }
          
          &.danger {
            color: rgb(var(--danger-6));
            
            &:hover {
              background: rgba(var(--danger-6), 0.05);
            }
          }
        }
      }
      
      .sync-time {
        font-size: 12px;
        color: var(--color-text-3);
        margin-top: 8px;
      }
    }
    
    .widget-empty {
      text-align: center;
      padding: 18px 14px;
      color: var(--color-text-3);
      
      .empty-icon {
        font-size: 28px;
        margin-bottom: 8px;
        opacity: 0.5;
      }
      
      p {
        margin: 0;
        font-size: 14px;
      }
    }
  }
}

/* Detail Shell Layout (主内容 + 右侧信息栏) */
.detail-shell {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.detail-top-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  align-items: start;
}

.detail-bottom-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
  align-items: start;
}

.bottom-main,
.bottom-side {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

@media (max-width: 1200px) {
  .detail-top-grid,
  .detail-bottom-grid {
    grid-template-columns: minmax(0, 1fr);
  }
}

.detail-shell .detail-card {
  background: var(--color-bg-2);
  border: 1px solid var(--color-border-2);
  border-radius: 12px;
  padding: 14px;
  transition: all 0.2s ease;

  &:hover {
    border-color: var(--color-border-3);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
  }
}

.detail-shell .detail-card .card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.detail-shell .detail-card .card-header .card-icon {
  color: rgb(var(--primary-6));
  font-size: 16px;
}

.detail-shell .detail-card .card-header .card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-1);
}

.detail-shell .detail-card .card-header .card-actions {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 8px;

  :deep(.arco-btn) {
    padding-left: 8px;
    padding-right: 8px;
  }
}

.detail-shell .detail-card .card-body .info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px 14px;
}

.detail-shell .detail-card .card-body .info-grid .info-item.full-width {
  grid-column: 1 / -1;
}

.detail-shell .detail-card .card-body .info-grid .info-item label {
  font-size: 12px;
  color: var(--color-text-3);
  margin-bottom: 4px;
  display: block;
}

.detail-shell .detail-card .card-body .info-grid .info-item .value {
  font-size: 14px;
  color: var(--color-text-1);
  font-weight: 500;
}

.detail-shell .detail-card .card-body .info-grid .info-item .value.phone-value {
  font-family: monospace;
}

.detail-shell .detail-card .card-body .info-grid .info-item .value.highlight {
  color: rgb(var(--primary-6));
  font-weight: 600;
}

.detail-shell .detail-card .card-body .stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.detail-shell .detail-card .card-body .stats-grid .stat-item {
  text-align: center;
  padding: 10px;
  background: var(--color-fill-1);
  border-radius: 8px;
}

.detail-shell .detail-card .card-body .stats-grid .stat-item .stat-value {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-1);
  margin-bottom: 4px;
}

.detail-shell .detail-card .card-body .stats-grid .stat-item .stat-label {
  font-size: 12px;
  color: var(--color-text-3);
}

.tags-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.pill-tag {
  border-radius: 999px;
  padding: 0 10px;
  border-color: var(--color-border-2);
  background: rgba(var(--primary-6), 0.08);
  color: rgb(var(--primary-6));

  :deep(.arco-tag-content) {
    max-width: 180px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.desc-text {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.7;
  color: var(--color-text-1);
}

/* Smart Lock Widget (脱离 detail-grid，兼容右侧栏展示) */
.lock-widget {
  &.active {
    background: linear-gradient(135deg, var(--color-bg-2) 0%, var(--color-fill-1) 100%);
    border-color: rgba(var(--primary-6), 0.3);
  }

  .widget-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;

    .title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 15px;
      font-weight: 600;
      color: var(--color-text-1);
    }

    .status {
      .battery-indicator {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 13px;
        font-weight: 500;

        &.color-success {
          color: rgb(var(--success-6));
        }
        &.color-warning {
          color: rgb(var(--warning-6));
        }
        &.color-danger {
          color: rgb(var(--danger-6));
        }
      }
    }

    .btn-bind {
      border-radius: 6px;
    }
  }

  .widget-body {
    display: flex;
    align-items: center;
    gap: 14px;

    .device-info {
      flex: 1;
      display: flex;
      align-items: center;
      gap: 16px;

      .device-icon {
        width: 40px;
        height: 40px;
        background: var(--color-bg-1);
        border: 1px solid var(--color-border-2);
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 20px;
        color: rgb(var(--primary-6));
      }

      .device-text {
        min-width: 0;

        .d-name {
          font-weight: 600;
          font-size: 14px;
          color: var(--color-text-1);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .d-mac {
          font-size: 12px;
          color: var(--color-text-3);
          font-family: monospace;
          margin-top: 4px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }

    .control-grid {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;

      .control-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 88px;
        height: 36px;
        border-radius: 10px;
        background: var(--color-bg-1);
        border: 1px solid var(--color-border-2);
        cursor: pointer;
        transition: all 0.2s ease;
        font-size: 12px;
        color: var(--color-text-2);
        padding: 0 10px;
        white-space: nowrap;

        &:hover {
          border-color: var(--color-border-3);
          box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
        }

        &.primary {
          background: rgb(var(--primary-6));
          color: #fff;
          border: none;

          &:hover {
            background: rgb(var(--primary-5));
            box-shadow: 0 4px 12px rgba(var(--primary-6), 0.3);
          }
        }

        &.danger {
          color: rgb(var(--danger-6));

          &:hover {
            background: rgba(var(--danger-6), 0.05);
          }
        }
      }
    }

    .sync-time {
      font-size: 12px;
      color: var(--color-text-3);
      margin-top: 8px;
    }
  }

  .widget-empty {
    text-align: center;
    padding: 18px 14px;
    color: var(--color-text-3);

    .empty-icon {
      font-size: 28px;
      margin-bottom: 8px;
      opacity: 0.5;
    }

    p {
      margin: 0;
      font-size: 14px;
    }
  }
}

/* Gallery Section */
.gallery-wrapper {
  padding: 16px;
  
  .masonry-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
    
    .grid-item {
      border-radius: 8px;
      overflow: hidden;
      transition: all 0.2s ease;
      cursor: pointer;
      border: 1px solid var(--color-border-2);
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
        z-index: 10;
      }
    }
  }
}

@media (max-width: 768px) {
  .detail-header {
    padding: 16px 16px 0;
  }

  .detail-body {
    padding: 12px 16px 20px;
  }

  .gallery-wrapper {
    padding: 12px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .gallery-wrapper .masonry-grid .grid-item {
    transition: none !important;
  }
}

/* Renovation Timeline */
.renovation-timeline {
  padding: 20px;
  
  .timeline-header {
    display: flex;
    align-items: center;
    gap: 24px;
    margin-bottom: 32px;
    
    .status-badge {
      font-size: 18px;
      font-weight: 600;
      padding: 8px 16px;
      border-radius: 20px;
      
      &.none {
        color: var(--color-text-3);
        background: var(--color-fill-2);
      }
      &.in_progress {
        color: rgb(var(--primary-6));
        background: rgba(var(--primary-6), 0.1);
      }
      &.done {
        color: rgb(var(--success-6));
        background: rgba(var(--success-6), 0.1);
      }
    }
    
    .progress-info {
      flex: 1;
      display: flex;
      align-items: center;
      gap: 12px;
      
      span {
        font-size: 13px;
        color: var(--color-text-2);
        white-space: nowrap;
      }
      
      :deep(.arco-progress) {
        width: 100%;
        max-width: 400px;
      }
    }
  }
  
  .timeline-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
    
    .t-card {
      background: var(--color-fill-1);
      padding: 16px;
      border-radius: 8px;
      border: 1px solid var(--color-border-2);
      
      &.wide {
        grid-column: 1 / -1;
      }
      
      .label {
        color: var(--color-text-3);
        font-size: 12px;
        margin-bottom: 8px;
        display: block;
      }
      
      .val {
        font-size: 15px;
        color: var(--color-text-1);
        font-weight: 500;
        
        &.highlight {
          color: rgb(var(--primary-6));
          font-weight: 600;
          font-size: 16px;
        }
        
        &.text-area {
          font-size: 14px;
          line-height: 1.6;
          font-weight: 400;
        }
      }
    }
  }

  .stage-timeline {
    margin-top: 18px;
    padding-top: 14px;
    border-top: 1px dashed var(--color-border-2);
  }

  .stage-item {
    display: flex;
    align-items: stretch;
    gap: 14px;
    padding: 8px 0;
  }

  .stage-rail {
    width: 18px;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding-top: 10px;
    flex-shrink: 0;
  }

  .stage-dot {
    width: 10px;
    height: 10px;
    border-radius: 999px;
    background: var(--color-fill-3);
    border: 2px solid var(--color-border-3);
    box-sizing: border-box;

    &.doing {
      background: rgb(var(--primary-6));
      border-color: rgba(var(--primary-6), 0.35);
      box-shadow: 0 0 0 4px rgba(var(--primary-6), 0.12);
    }
    &.done {
      background: rgb(var(--success-6));
      border-color: rgba(var(--success-6), 0.35);
      box-shadow: 0 0 0 4px rgba(var(--success-6), 0.12);
    }
  }

  .stage-line {
    width: 2px;
    flex: 1;
    background: var(--color-border-2);
    margin-top: 6px;
    border-radius: 2px;
  }

  .stage-body {
    flex: 1;
    background: var(--color-fill-1);
    border: 1px solid var(--color-border-2);
    border-radius: 10px;
    padding: 12px 14px;
    min-width: 0;
    transition: border-color 0.18s ease, box-shadow 0.18s ease, transform 0.18s ease;

    &:hover {
      border-color: rgba(var(--primary-6), 0.35);
      box-shadow: 0 10px 26px rgba(0, 0, 0, 0.06);
      transform: translateY(-1px);
    }
  }

  .stage-head {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
  }

  .stage-title {
    font-weight: 700;
    color: var(--color-text-1);
  }

  .stage-date {
    font-size: 12px;
    color: var(--color-text-3);
  }

  .stage-note {
    margin-top: 6px;
    font-size: 13px;
    color: var(--color-text-2);
    line-height: 1.65;
    white-space: pre-wrap;
    word-break: break-word;
  }

  .stage-media {
    margin-top: 10px;
  }

  .stage-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 10px;

    @media (max-width: 1400px) {
      grid-template-columns: repeat(3, 1fr);
    }
    @media (max-width: 1200px) {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  .stage-img {
    position: relative;
    width: 100%;
    height: 120px;
    border-radius: 10px;
    overflow: hidden;
    background: var(--color-fill-2);
    border: 1px solid var(--color-border-2);
    transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;

    &:hover {
      transform: translateY(-1px);
      border-color: rgba(var(--primary-6), 0.55);
      box-shadow: 0 10px 24px rgba(0, 0, 0, 0.08);
    }

    :deep(.arco-image) {
      width: 100%;
      height: 100%;
    }

    .img-actions {
      position: absolute;
      top: 8px;
      right: 8px;
      z-index: 2;
    }
    .download-btn {
      box-shadow: 0 8px 22px rgba(0, 0, 0, 0.12);
    }
  }
}

/* Cloud Detail Drawer */
.cloud-drawer {
  :deep(.arco-drawer-body) {
    padding: 20px;
  }
}

.statuslog-drawer {
  :deep(.arco-drawer-body) {
    padding: 16px 16px 20px;
  }

  .log-before {
    color: var(--color-text-3);
  }
  .log-arrow {
    margin: 0 8px;
    color: var(--color-text-4);
  }
  .log-after {
    color: var(--color-text-1);
    font-weight: 600;
  }
}

.cloud-detail-box {
  .cd-hero {
    text-align: center;
    margin-bottom: 24px;
    
    .icon-circle {
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
    
    .cd-title {
      font-size: 18px;
      font-weight: 700;
      color: var(--color-text-1);
      margin-bottom: 4px;
    }
    
    .cd-sub {
      font-size: 12px;
      color: var(--color-text-3);
      font-family: monospace;
    }
  }
  
  .cd-list {
    .cd-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid var(--color-border-2);
      
      &:last-child {
        border-bottom: none;
      }
      
      span:first-child {
        font-size: 13px;
        color: var(--color-text-3);
      }
      
      .cd-val {
        font-size: 14px;
        color: var(--color-text-1);
        font-weight: 500;
        font-family: monospace;
      }
    }
  }
}
</style>

