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
    :title="getTitle"
    :showOkBtn="false"
    :showCancelBtn="false"
    class="pro-modal"
  >
    <div class="add-form-container" :style="{ 'min-height': `${windHeight}px` }">
      <a-tabs v-model:active-key="activeTab" type="text" class="pro-tabs">
        <!-- Tab 1: 基础房源信息 -->
        <a-tab-pane key="basic">
          <template #title>
             <span class="tab-label"><icon-home /> 房源信息</span>
          </template>
          <div class="wizard-layout">
            <div class="form-shell">
              <!-- Left: 快速导航（顶部与右侧重复，统一放左侧） -->
              <div class="form-nav">
                <div class="side-stack">
                  <div class="side-card nav-card">
                    <div class="nav-head">
                      <div class="side-title">快速导航</div>
                      <a-tag v-if="formData.id" color="blue">#{{ formData.id }}</a-tag>
                      <a-tag v-else color="gray">未生成ID</a-tag>
                    </div>
                    <a-progress
                      :percent="(wizardStep + 1) / 4"
                      size="mini"
                      :show-text="false"
                      class="nav-progress"
                    />
                    <a-steps direction="vertical" size="small" :current="wizardStep + 1" class="side-steps">
                      <a-step @click="handleWizardJump(0)" title="必填基础" description="必填信息 + 生成草稿ID" />
                      <a-step @click="handleWizardJump(1)" title="位置&配套" description="户型/楼层/物业等可选信息" />
                      <a-step @click="handleWizardJump(2)" title="图片&标签" description="房源图集与特色标签" />
                      <a-step @click="handleWizardJump(3)" title="销售&发布" description="佣金/推荐/启用等发布信息" />
                    </a-steps>
                    <div class="nav-tip">提示：点击第 2-4 步会先创建草稿。</div>
                  </div>
                </div>
              </div>

              <div class="form-main">
                <div class="form-content">
                  <a-alert
                    v-if="formError"
                    type="error"
                    show-icon
                    closable
                    role="alert"
                    class="form-error-alert"
                    @close="formError = ''"
                  >
                    <template #title>请完善必填项</template>
                    {{ formError }}
                  </a-alert>
                  <a-form ref="formRef" :model="formData" auto-label-width layout="vertical" class="pro-form">
                 
                <!-- Section: 核心属性 -->
                <div class="form-section" v-show="wizardStep === 0">
                   <div class="section-header">
                      <div class="title">核心属性</div>
                      <div class="sub">房源的基本描述与关键参数</div>
                   </div>
                   <div class="section-body">
                      <a-row :gutter="24">
                        <a-col :span="24">
                          <a-form-item field="title" label="房源标题" :rules="[{ required: true, message: '请填写房源标题' }]">
                            <a-input v-model="formData.title" placeholder="如：阳光花园精装三居室，近地铁..." allow-clear :max-length="100" show-word-limit size="large">
                               <template #prefix><icon-tag /></template>
                            </a-input>
                          </a-form-item>
                        </a-col>
                        <a-col :span="12">
                           <a-row :gutter="12">
                              <a-col :span="16">
                                <a-form-item field="price" label="售价" :rules="[positiveNumberRule('请填写价格')]">
                                  <a-input-number v-model="formData.price" :min="0" :precision="2" placeholder="0.00" size="large" class="price-input">
                                    <template #prefix>¥</template>
                                  </a-input-number>
                                </a-form-item>
                              </a-col>
                              <a-col :span="8">
                                <a-form-item field="price_unit" label="单位">
                                  <a-select v-model="formData.price_unit" size="large">
                                    <a-option v-for="opt in priceUnitOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</a-option>
                                  </a-select>
                                </a-form-item>
                              </a-col>
                           </a-row>
                        </a-col>
                        <a-col :span="6">
                          <a-form-item field="area" label="建筑面积" :rules="[positiveNumberRule('请填写建筑面积')]">
                             <a-input-number v-model="formData.area" :min="0" :precision="2" placeholder="0" size="large">
                                <template #suffix>㎡</template>
                             </a-input-number>
                          </a-form-item>
                        </a-col>
                        <a-col :span="6">
                          <a-form-item field="build_year" label="建成年份" :rules="[{ required: true, message: '请选择年份' }]">
                            <a-date-picker v-model="formData.build_year" mode="year" value-format="YYYY" placeholder="请选择" style="width: 100%" size="large"/>
                          </a-form-item>
                        </a-col>
                       </a-row>
                   </div>
                </div>

                <!-- Section: 房主与收房（更直观，放在第一步） -->
                <div class="form-section" v-show="wizardStep === 0">
                  <div class="section-header">
                    <div class="title">房主与收房</div>
                    <div class="sub">业务关键信息（可选）</div>
                  </div>
                  <div class="section-body">
                    <a-row :gutter="24">
                      <a-col :span="6">
                        <a-form-item field="owner_name" label="房主姓名">
                          <a-input v-model="formData.owner_name" placeholder="可选" size="large" allow-clear>
                            <template #prefix><icon-user /></template>
                          </a-input>
                        </a-form-item>
                      </a-col>
                      <a-col :span="6">
                        <a-form-item field="owner_phone" label="房主电话">
                          <a-input v-model="formData.owner_phone" placeholder="可选" size="large" allow-clear>
                            <template #prefix><icon-phone /></template>
                          </a-input>
                        </a-form-item>
                      </a-col>
                      <a-col :span="6">
                        <a-form-item field="receiver_name" label="收房人姓名">
                          <a-input v-model="formData.receiver_name" placeholder="可选" size="large" allow-clear>
                            <template #prefix><icon-user-group /></template>
                          </a-input>
                        </a-form-item>
                      </a-col>
                      <a-col :span="6">
                        <a-form-item field="receiver_phone" label="收房人电话">
                          <a-input v-model="formData.receiver_phone" placeholder="可选" size="large" allow-clear>
                            <template #prefix><icon-phone /></template>
                          </a-input>
                        </a-form-item>
                      </a-col>
                      <a-col :span="8">
                        <a-form-item field="receiver_price" label="收房价格(支付业主)">
                          <a-input-number
                            v-model="formData.receiver_price"
                            :min="0"
                            :precision="2"
                            size="large"
                            allow-clear
                            hide-button
                            placeholder="未填写"
                          >
                            <template #prefix>¥</template>
                          </a-input-number>
                        </a-form-item>
                      </a-col>
                    </a-row>
                  </div>
                </div>

                <!-- Section: 户型与规格 -->
                <div class="form-section" v-show="wizardStep === 1">
                   <div class="section-header">
                      <div class="title">户型与规格</div>
                   </div>
                   <div class="section-body">
                      <a-row :gutter="24">
                         <!-- 户型胶囊输入 -->
                         <a-col :span="10">
                            <a-form-item label="户型结构" required>
                               <div class="capsule-input-group">
                                  <a-input-number v-model="formData.rooms" :min="0" :max="9" hide-button size="large" class="cap-input">
                                     <template #suffix>室</template>
                                  </a-input-number>
                                  <div class="cap-divider"></div>
                                  <a-input-number v-model="formData.halls" :min="0" :max="9" hide-button size="large" class="cap-input">
                                     <template #suffix>厅</template>
                                  </a-input-number>
                                  <div class="cap-divider"></div>
                                  <a-input-number v-model="formData.bathrooms" :min="0" :max="9" hide-button size="large" class="cap-input">
                                     <template #suffix>卫</template>
                                  </a-input-number>
                               </div>
                            </a-form-item>
                         </a-col>
                         <a-col :span="7">
                            <a-form-item field="floor_level" label="楼层位置">
                               <a-select v-model="formData.floor_level" placeholder="请选择" allow-clear size="large">
                                   <a-option v-for="opt in floorLevelOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</a-option>
                               </a-select>
                            </a-form-item>
                         </a-col>
                         <a-col :span="7">
                            <a-form-item field="total_floors" label="总楼层">
                               <a-input-number v-model="formData.total_floors" :min="1" placeholder="共多少层" size="large">
                                  <template #suffix>层</template>
                               </a-input-number>
                            </a-form-item>
                         </a-col>
                         
                         <a-col :span="8">
                            <a-form-item field="orientation" label="朝向">
                               <a-select v-model="formData.orientation" placeholder="请选择" allow-clear size="large">
                                   <a-option v-for="opt in orientationOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</a-option>
                               </a-select>
                            </a-form-item>
                         </a-col>
                         <a-col :span="8">
                            <a-form-item field="property_type" label="物业类型">
                               <a-select v-model="formData.property_type" placeholder="请选择" allow-clear size="large">
                                   <a-option v-for="opt in propertyTypeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</a-option>
                               </a-select>
                            </a-form-item>
                         </a-col>
                         <a-col :span="8">
                            <a-form-item field="decoration_type" label="装修标准">
                               <a-select v-model="formData.decoration_type" placeholder="请选择" allow-clear size="large">
                                   <a-option v-for="opt in decorationTypeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</a-option>
                               </a-select>
                            </a-form-item>
                         </a-col>
                      </a-row>
                   </div>
                </div>

                <!-- Section: 地址与定位 -->
                 <div class="form-section" v-show="wizardStep === 0">
                   <div class="section-header">
                      <div class="title">地址与定位</div>
                   </div>
                   <div class="section-body">
                      <a-row :gutter="24">
                         <a-col :span="12">
                            <a-form-item field="community_name" label="小区名称" :rules="[requiredRule('请填写小区名称')]">
                               <a-input v-model="formData.community_name" placeholder="输入小区名称" size="large">
                                  <template #prefix><icon-community /></template>
                               </a-input>
                            </a-form-item>
                         </a-col>
                         <a-col :span="12">
                            <a-form-item field="area_ids" label="行政区域" :rules="[cascaderRule('请选择省/市/区')]">
                               <a-cascader
                                  :key="areaCascaderKey"
                                  v-model="formData.area_ids"
                                  :options="areaOptions"
                                  :load-more="loadAreaMore"
                                  :field-names="{ value: 'value', label: 'label', children: 'children', isLeaf: 'isLeaf' }"
                                  path-mode
                                  placeholder="省 / 市 / 区"
                                  allow-clear
                                  size="large"
                               />
                            </a-form-item>
                         </a-col>
                         <a-col :span="24">
                            <a-form-item field="address_detail" label="详细地址" :rules="[requiredRule('请填写详细地址')]">
                               <a-input v-model="formData.address_detail" placeholder="街道、门牌号、楼栋单元信息等" size="large">
                                  <template #prefix><icon-location /></template>
                               </a-input>
                            </a-form-item>
                         </a-col>
                      </a-row>
                   </div>
                 </div>
                  <!-- Section: 地图坐标（可选） -->
                  <div class="form-section" v-show="wizardStep === 1">
                     <div class="section-header">
                       <div class="title">地图坐标</div>
                      <div class="sub">非必填，可用于地图展示/附近推荐</div>
                    </div>
                    <div class="section-body">
                      <div class="coord-actions">
                        <a-space>
                          <a-button size="small" @click="openMapPicker">
                            <template #icon><icon-location /></template>
                            地图选点
                          </a-button>
                          <a-button size="small" @click="clearLatLng">清空</a-button>
                          <span class="coord-tip">选点后自动回填经纬度（保留 6 位小数）</span>
                        </a-space>
                      </div>
                      <a-row :gutter="24">
                        <a-col :span="12">
                          <a-form-item field="longitude" label="经度 (Lng)">
                            <a-input v-model="formData.longitude" placeholder="可从地图拾取" size="large">
                              <template #prefix><icon-compass /></template>
                            </a-input>
                          </a-form-item>
                        </a-col>
                        <a-col :span="12">
                          <a-form-item field="latitude" label="纬度 (Lat)">
                            <a-input v-model="formData.latitude" placeholder="可从地图拾取" size="large">
                              <template #prefix><icon-compass /></template>
                            </a-input>
                          </a-form-item>
                        </a-col>
                      </a-row>
                    </div>
                 </div>

                 <LatLngPickerModal
                   v-model:visible="mapPickerVisible"
                   :latitude="formData.latitude"
                   :longitude="formData.longitude"
                   @confirm="handleMapPicked"
                   @clear="clearLatLng"
                 />

                  <!-- Section: 图片与标签 -->
                  <div class="form-section" v-show="wizardStep === 2">
                     <div class="section-header">
                       <div class="title">图片与标签</div>
                       <div class="sub">用于列表展示与客户筛选</div>
                     </div>
                     <div class="section-body">
                       <a-row :gutter="24">
                         <a-col :span="24">
                           <a-form-item label="房源特色标签">
                             <a-input-tag v-model="formData.tags" placeholder="输入后回车，如：随时看房、满五唯一" size="large" allow-clear />
                           </a-form-item>
                         </a-col>
                         <a-col :span="24">
                           <div class="inline-section-head">
                             <div class="inline-title">发布信息</div>
                             <div class="inline-sub">第一步会创建草稿ID（用于装修/门锁等绑定）</div>
                           </div>
                         </a-col>
                         <a-col :span="8">
                            <a-form-item field="sale_status" label="销售状态" :rules="[requiredRule('请选择销售状态')]">
                              <a-select v-model="formData.sale_status" size="large">
                                <a-option v-for="opt in saleStatusOptions" :key="opt.value" :value="opt.value">
                                  <a-badge :status="getSaleStatusBadge(opt.value)" :text="opt.label" />
                                </a-option>
                              </a-select>
                            </a-form-item>
                          </a-col>
                         <a-col :span="16">
                           <a-form-item field="cover_image" label="封面主图" :rules="[requiredRule('请上传封面主图')]">
                             <FormImageBox v-model="formData.cover_image" />
                           </a-form-item>
                         </a-col>
                         <a-col :span="24">
                           <a-form-item field="video_url" label="房源视频">
                             <FormVideoBox v-model="formData.video_url" />
                           </a-form-item>
                         </a-col>
                         <a-col :span="12">
                           <a-form-item field="allow_image_download" label="图片允许下载" class="mb-0">
                             <a-switch v-model="formData.allow_image_download" :checked-value="1" :unchecked-value="0" type="round">
                               <template #checked>允许</template>
                               <template #unchecked>不允许</template>
                             </a-switch>
                             <div class="switch-tip">封面 + 图集统一开关</div>
                           </a-form-item>
                         </a-col>
                         <a-col :span="12">
                           <a-form-item field="allow_video_download" label="视频允许下载" class="mb-0">
                             <a-switch v-model="formData.allow_video_download" :checked-value="1" :unchecked-value="0" type="round">
                               <template #checked>允许</template>
                               <template #unchecked>不允许</template>
                             </a-switch>
                           </a-form-item>
                         </a-col>
                         <a-col :span="24">
                           <a-form-item field="images" label="详细图集">
                             <FormImagesBox v-model="formData.images" />
                           </a-form-item>
                         </a-col>
                       </a-row>
                    </div>
                 </div>

                 <!-- Section: 销售与发布 -->
                 <div class="form-section" v-show="wizardStep === 3">
                    <div class="section-header">
                      <div class="title">销售与发布</div>
                      <div class="sub">佣金、推荐、启用状态等</div>
                    </div>
                    <div class="section-body">
                      <a-row :gutter="24">
                        <a-col :span="6">
                          <a-form-item field="commission_rate" label="佣金比例 (%)">
                            <a-input-number v-model="formData.commission_rate" :min="0" :max="100" :precision="1" size="large">
                              <template #suffix>%</template>
                            </a-input-number>
                          </a-form-item>
                        </a-col>
                        <a-col :span="6">
                          <a-form-item field="commission_reward" label="成交奖励金">
                            <a-input-number v-model="formData.commission_reward" :min="0" :precision="2" size="large">
                              <template #prefix>¥</template>
                            </a-input-number>
                          </a-form-item>
                        </a-col>
                        <a-col :span="6">
                          <a-form-item field="weigh" label="排序权重">
                            <a-input-number v-model="formData.weigh" :min="0" size="large" />
                          </a-form-item>
                        </a-col>
                        <a-col :span="6">
                          <a-form-item field="hot_status" label="是否推荐" class="mb-0">
                            <a-switch v-model="formData.hot_status" :checked-value="1" :unchecked-value="0" type="round">
                              <template #checked>推荐</template>
                              <template #unchecked>普通</template>
                            </a-switch>
                          </a-form-item>
                        </a-col>
                        <a-col :span="6">
                          <a-form-item field="has_smart_lock" label="智能门锁" class="mb-0">
                            <a-switch v-model="formData.has_smart_lock" :checked-value="1" :unchecked-value="0" type="round" disabled>
                              <template #checked><icon-check /> 已绑定</template>
                              <template #unchecked><icon-close /> 未绑定</template>
                            </a-switch>
                            <div class="switch-tip">请在详情页绑定</div>
                          </a-form-item>
                        </a-col>
                        <a-col :span="6">
                          <a-form-item field="status" label="数据状态" class="mb-0">
                            <a-switch v-model="formData.status" :checked-value="0" :unchecked-value="1" type="round">
                              <template #checked>启用</template>
                              <template #unchecked>禁用</template>
                            </a-switch>
                          </a-form-item>
                        </a-col>

                        <a-col :span="24">
                          <div class="final-summary">
                            <a-alert type="info" show-icon>
                              <template #title>提交前确认</template>
                              <template #default>
                                 <span style="margin-right: 6px;">标题：</span><span class="strong">{{ formData.title || '-' }}</span>
                                 <span class="summary-sep">|</span>
                                 <span style="margin-right: 6px;">小区：</span><span class="strong">{{ formData.community_name || '-' }}</span>
                                 <span class="summary-sep">|</span>
                                 <span style="margin-right: 6px;">状态：</span><span class="strong">{{ getSaleStatusLabel(formData.sale_status) }}</span>
                               </template>
                             </a-alert>
                           </div>
                         </a-col>
                       </a-row>
                    </div>
                 </div>

                   </a-form>
                 </div>
              </div>

              <div class="form-side">
                <div class="side-stack">
                  <div class="side-card">
                    <div class="side-title">摘要</div>
                    <div class="kv">
                      <div class="kv-row">
                        <span class="k">标题</span>
                        <span class="v">{{ formData.title || '-' }}</span>
                      </div>
                      <div class="kv-row">
                        <span class="k">小区</span>
                        <span class="v">{{ formData.community_name || '-' }}</span>
                      </div>
                      <div class="kv-row">
                        <span class="k">售价</span>
                        <span class="v">¥{{ Number(formData.price || 0).toFixed(2) }} {{ formData.price_unit || '万' }}</span>
                      </div>
                    </div>
                  </div>

                  <div class="side-card">
                    <div class="side-title">标签预览</div>
                    <div class="tag-wrap" v-if="(formData.tags || []).length">
                      <a-tag v-for="(t, i) in (formData.tags || []).slice(0, 8)" :key="i" size="small" bordered class="pill-tag">
                        {{ t }}
                      </a-tag>
                    </div>
                    <a-empty v-else description="暂无标签" />
                  </div>
                </div>
              </div>
            </div>

          </div>
        </a-tab-pane>

        <!-- Tab 2: 装修进度管理 -->
        <a-tab-pane key="renovation" :disabled="!formData.id">
           <template #title>
             <span class="tab-label"><icon-tool /> 装修进度</span>
          </template>
          <div class="form-content">
            <a-alert v-if="!formData.id" type="warning" banner>请先保存基本信息后再填写装修记录。</a-alert>

            <a-form ref="renovationFormRef" :model="renovationData" auto-label-width layout="vertical" class="pro-form">
                   <div class="form-section">
                      <div class="section-header">
                         <div class="title">总体状态</div>
                      </div>
                      <div class="section-body">
                         <a-row :gutter="24">
                            <a-col :span="8">
                               <a-form-item field="renovation_status" label="装修阶段" required>
                                  <a-select v-model="renovationData.renovation_status" size="large" placeholder="选择状态">
                                     <a-option value="none">未装修</a-option>
                                     <a-option value="in_progress">装修中</a-option>
                                     <a-option value="done">已完工</a-option>
                                  </a-select>
                               </a-form-item>
                            </a-col>
                            <a-col :span="16" v-if="showRenovationProgress">
                               <a-form-item field="progress_percentage" label="整体进度">
                                  <a-slider v-model="renovationData.progress_percentage" :max="100" :step="5" show-input />
                               </a-form-item>
                            </a-col>
                         </a-row>
                      </div>
                   </div>

                   <div class="form-section" v-if="showRenovationProgress">
                      <div class="section-header">
                         <div class="title">施工详情</div>
                      </div>
                      <div class="section-body">
                         <a-row :gutter="24">
                            <a-col :span="8">
                               <a-form-item field="current_stage" label="当前工序">
                                  <a-select v-model="renovationData.current_stage" placeholder="如：水电改造" allow-clear size="large">
                                     <a-option v-for="o in renovationStageOptions" :key="o.value" :value="o.value">{{ o.label }}</a-option>
                                  </a-select>
                               </a-form-item>
                            </a-col>
                            <a-col :span="8">
                               <a-form-item field="start_date" label="开工日期" required>
                                  <a-date-picker v-model="renovationData.start_date" style="width: 100%" size="large"/>
                               </a-form-item>
                            </a-col>
                            <a-col :span="8">
                               <a-form-item field="estimated_finish_date" label="预计完工">
                                  <a-date-picker v-model="renovationData.estimated_finish_date" style="width: 100%" size="large"/>
                               </a-form-item>
                            </a-col>
                            <a-col :span="24">
                               <a-form-item field="notes" label="施工日志 / 备注">
                                  <a-textarea v-model="renovationData.notes" :auto-size="{minRows:3}" placeholder="记录当前施工重点..." />
                               </a-form-item>
                            </a-col>
                            <a-col :span="24">
                               <a-form-item label="现场照片">
                                  <FormImagesBox v-model="renovationData.images" />
                               </a-form-item>
                            </a-col>
                         </a-row>
                      </div>
                   </div>
            </a-form>
          </div>
        </a-tab-pane>
      </a-tabs>

      <!-- 底部操作栏：居中 + 圆角（每个页面都一致） -->
      <div class="pro-sticky-footer">
        <div class="footer-card">
          <a-space size="medium" class="footer-actions" align="center">
            <a-button size="large" shape="round" type="outline" @click="closeModal">关闭</a-button>

            <template v-if="activeTab === 'basic'">
              <a-button
                v-if="formData.id"
                size="large"
                shape="round"
                type="secondary"
                :loading="basicLoading"
                @click="handleQuickSave"
              >
                <template #icon><icon-save /></template>
                保存
              </a-button>
              <a-button v-if="wizardStep > 0" size="large" shape="round" type="secondary" @click="handleWizardPrev">上一步</a-button>
              <a-button
                v-if="wizardStep < 3"
                type="primary"
                size="large"
                shape="round"
                :loading="basicLoading"
                @click="handleWizardNext"
              >
                <template #icon><icon-right /></template>
                {{ wizardStep === 0 ? (formData.id ? '下一步' : '保存草稿并继续') : '下一步' }}
              </a-button>
              <a-button
                v-else
                type="primary"
                size="large"
                shape="round"
                :loading="basicLoading"
                @click="handleWizardFinish"
              >
                <template #icon><icon-save /></template>
                完成保存
              </a-button>
            </template>

            <template v-else>
              <a-button size="large" shape="round" type="secondary" @click="goToBasicTab">返回房源信息</a-button>
              <a-button
                size="large"
                shape="round"
                type="primary"
                :loading="renovationLoading"
                :disabled="!formData.id"
                @click="handleSaveRenovation"
              >
                <template #icon><icon-save /></template>
                保存装修记录
              </a-button>
            </template>
          </a-space>
        </div>
      </div>
    </div>
  </BasicModal>
</template>

<script lang="ts">
import { defineComponent, ref, computed, unref, watch } from 'vue';
import { BasicModal, useModalInner } from '/@/components/Modal';
import type { FormInstance } from '@arco-design/web-vue';
import useLoading from '@/hooks/loading';
import { cloneDeep } from 'lodash-es';
import { save, getContent, getRenovation, saveRenovation, getAreaMoreList, getFormOptions } from './api';
import { Message } from '@arco-design/web-vue';
import FormImageBox from '@/components/autoPlugin/Form/FormImageBox.vue';
import FormImagesBox from '@/components/autoPlugin/Form/FormImagesBox.vue';
import FormVideoBox from '@/components/autoPlugin/Form/FormVideoBox.vue';
import LatLngPickerModal from '@/components/Map/LatLngPickerModal.vue';

export default defineComponent({
  name: 'HousesAddForm',
  components: { BasicModal, FormImageBox, FormImagesBox, FormVideoBox, LatLngPickerModal },
  emits: ['success'],
  setup(_, { emit }) {
    const isUpdate = ref(false);
    // 不默认全屏：通过加大弹窗高度减少滚动
    const modelHeight = ref(900);
    const windHeight = ref(900);
    const formRef = ref<FormInstance>();
    const renovationFormRef = ref<FormInstance>();
    const formError = ref('');
    const activeTab = ref('basic');
    // 录入向导步骤：0必填基础 -> 1位置&配套 -> 2图片&标签 -> 3销售&发布
    const wizardStep = ref(0);
    // 用于强制刷新 Cascader（解决编辑回显时显示为 ID 的问题）
    const areaCascaderKey = ref(1);

    // 基本信息数据
    const baseData: any = {
      id: 0,
      title: '',
      price: 0,
      price_unit: '万',
      area: 0,
      rooms: 0,
      halls: 0,
      bathrooms: 0,
      floor_level: '',
      total_floors: 0,
      orientation: '',
      build_year: null,
      property_type: '',
      decoration_type: '',
      community_name: '',
      // address 为最终提交字段；area_ids + address_detail 用于分段输入
      address: '',
      area_ids: [] as (number | string)[],
      address_detail: '',
      latitude: '',
      longitude: '',
      tags: [] as string[],
      images: '',
      cover_image: '',
      video_url: '',
      allow_image_download: 1,
      allow_video_download: 1,
      has_smart_lock: 0,
      commission_rate: 0,
      commission_reward: 0,
      owner_name: '',
      owner_phone: '',
      receiver_name: '',
      receiver_phone: '',
      receiver_price: 0,
      sale_status: 'on_sale',
      hot_status: 0,
      status: 0,
      weigh: 0,
    };

    // 装修信息数据
    const baseRenovationData: any = {
      renovation_status: 'none',
      progress_percentage: 0,
      current_stage: '',
      start_date: '',
      estimated_finish_date: '',
      actual_finish_date: '',
      materials: [] as string[],
      images: '',
      notes: '',
      status: 0,
    };

    const formData = ref<any>(cloneDeep(baseData));
    const renovationData = ref<any>(cloneDeep(baseRenovationData));

    type SelectOption = { label: string; value: string };
    const saleStatusOptions = ref<SelectOption[]>([
      { label: '在售', value: 'on_sale' },
      { label: '已售', value: 'sold' },
      { label: '下架', value: 'off_market' },
    ]);
    const priceUnitOptions = ref<SelectOption[]>([
      { label: '万', value: '万' },
      { label: '元', value: '元' },
    ]);
    const floorLevelOptions = ref<SelectOption[]>([
      { label: '低层 (1-6)', value: '低层' },
      { label: '中层 (7-15)', value: '中层' },
      { label: '高层 (16+)', value: '高层' },
      { label: '地下', value: '地下' },
    ]);
    const orientationOptions = ref<SelectOption[]>(
      ['东', '南', '西', '北', '东南', '东北', '西南', '西北', '南北', '东西'].map((it) => ({
        label: it,
        value: it,
      }))
    );
    const propertyTypeOptions = ref<SelectOption[]>(['住宅', '公寓', '别墅', '商铺', '写字楼'].map((it) => ({ label: it, value: it })));
    const decorationTypeOptions = ref<SelectOption[]>(['毛坯', '简装', '精装', '豪装'].map((it) => ({ label: it, value: it })));
    const renovationStageOptions = ref<SelectOption[]>(
      ['设计', '拆改', '水电', '泥瓦', '木工', '油漆', '安装', '软装', '验收'].map((it) => ({ label: it, value: it }))
    );

    const formOptionsLoaded = ref(false);
    const normalizeOptions = (arr: any): SelectOption[] =>
      (Array.isArray(arr) ? arr : [])
        .map((it: any) => ({
          label: String(it?.label ?? '').trim(),
          value: String(it?.value ?? '').trim(),
        }))
        .filter((it) => it.label && it.value);

    const applyFormOptions = (data: any) => {
      const d = data || {};

      const sale = normalizeOptions(d.sale_status);
      if (sale.length) saleStatusOptions.value = sale;

      const unit = normalizeOptions(d.price_unit);
      if (unit.length) priceUnitOptions.value = unit;

      const floor = normalizeOptions(d.floor_level);
      if (floor.length) floorLevelOptions.value = floor;

      const ori = normalizeOptions(d.orientation);
      if (ori.length) orientationOptions.value = ori;

      const pt = normalizeOptions(d.property_type);
      if (pt.length) propertyTypeOptions.value = pt;

      const deco = normalizeOptions(d.decoration_type);
      if (deco.length) decorationTypeOptions.value = deco;

      const stages = normalizeOptions(d.renovation_stage);
      if (stages.length) renovationStageOptions.value = stages;
    };

    const ensureFormOptionsLoaded = async () => {
      if (formOptionsLoaded.value) return;
      try {
        const data = await getFormOptions({});
        applyFormOptions(data);
        formOptionsLoaded.value = true;
      } catch (e) {
        // ignore
      }
    };

    const ensureFormSelectValue = (key: string, optionsRef: any, required = false, fallback = '') => {
      const v = String(formData.value?.[key] ?? '').trim();
      const options: SelectOption[] = (optionsRef?.value || []) as SelectOption[];
      if (v) {
        if (!options.some((o) => o.value === v)) {
          optionsRef.value = [...options, { label: v, value: v }];
        }
        return;
      }
      if (required) {
        formData.value[key] = (options[0]?.value || fallback || '').toString();
      }
    };

    const sanitizeFormSelectValues = () => {
      ensureFormSelectValue('sale_status', saleStatusOptions, true, 'on_sale');
      ensureFormSelectValue('price_unit', priceUnitOptions, true, '万');
      ensureFormSelectValue('floor_level', floorLevelOptions);
      ensureFormSelectValue('orientation', orientationOptions);
      ensureFormSelectValue('property_type', propertyTypeOptions);
      ensureFormSelectValue('decoration_type', decorationTypeOptions);

      // 装修工序：允许历史数据不在选项内时追加显示
      const rv = String(renovationData.value?.current_stage ?? '').trim();
      const ro: SelectOption[] = renovationStageOptions.value || [];
      if (rv && !ro.some((o) => o.value === rv)) {
        renovationStageOptions.value = [...ro, { label: rv, value: rv }];
      }
    };

    const isRenovationNone = computed(() => renovationData.value.renovation_status === 'none');

    // 地图选点：回填经纬度
    const mapPickerVisible = ref(false);
    const openMapPicker = () => {
      mapPickerVisible.value = true;
    };
    const clearLatLng = () => {
      formData.value.latitude = '';
      formData.value.longitude = '';
    };
    const handleMapPicked = (payload: any) => {
      if (!payload) return;
      formData.value.longitude = payload.longitude || '';
      formData.value.latitude = payload.latitude || '';
    };

    const goToBasicTab = () => {
      activeTab.value = 'basic';
    };
    const isRenovationInProgress = computed(() => renovationData.value.renovation_status === 'in_progress');
    const isRenovationDone = computed(() => renovationData.value.renovation_status === 'done');
    const showRenovationProgress = computed(() => !isRenovationNone.value);
    const showRenovationStage = computed(() => !isRenovationNone.value);
    const showRenovationDates = computed(() => !isRenovationNone.value);

    watch(
      () => renovationData.value.renovation_status,
      (val) => {
        // 按状态自动处理字段，避免用户误填
        if (val === 'none') {
          renovationData.value.progress_percentage = 0;
          renovationData.value.current_stage = '';
          renovationData.value.start_date = '';
          renovationData.value.estimated_finish_date = '';
          renovationData.value.actual_finish_date = '';
          return;
        }
        if (val === 'done') {
          renovationData.value.progress_percentage = 100;
          renovationData.value.estimated_finish_date = '';
          return;
        }
        // in_progress
        renovationData.value.actual_finish_date = '';
        if (typeof renovationData.value.progress_percentage !== 'number') {
          renovationData.value.progress_percentage = parseInt(String(renovationData.value.progress_percentage)) || 0;
        }
      },
      { immediate: true }
    );

    const { loading, setLoading } = useLoading();
    const { loading: basicLoading, setLoading: setBasicLoading } = useLoading();
    const { loading: renovationLoading, setLoading: setRenovationLoading } = useLoading();

    const toNumber = (v: any, fallback = 0) => {
      if (v === null || v === undefined || v === '') return fallback;
      const n = Number(v);
      return Number.isFinite(n) ? n : fallback;
    };

    const requiredRule = (message: string) => ({ required: true, message });
    const positiveNumberRule = (message: string) => ({
      validator: (value: any, callback: (error?: string) => void) => {
        const n = Number(value);
        if (!Number.isFinite(n) || n <= 0) {
          callback(message);
          return;
        }
        callback();
      },
    });
    const cascaderRule = (message: string) => ({
      validator: (value: any, callback: (error?: string) => void) => {
        if (Array.isArray(value) && value.length >= 3) {
          callback();
          return;
        }
        callback(message);
      },
    });

    const getFirstValidateMessage = (err: any) => {
      if (!err || typeof err !== 'object') return '请先完善必填项';
      const firstField = Object.keys(err)[0];
      const fieldErrors = (err as any)[firstField];
      if (Array.isArray(fieldErrors) && fieldErrors[0]?.message) return fieldErrors[0].message;
      if (fieldErrors?.message) return fieldErrors.message;
      return '请先完善必填项';
    };

    const getSaleStatusLabel = (s: string) => {
      const hit = saleStatusOptions.value.find((it) => it.value === s);
      return hit?.label || s || '-';
    };

    const getSaleStatusBadge = (s: string) => {
      if (s === 'on_sale') return 'success';
      if (s === 'sold') return 'normal';
      if (s === 'off_market') return 'warning';
      return 'processing';
    };

    const basicStepFields = [
      'title',
      'price',
      'area',
      'build_year',
      'community_name',
      'area_ids',
      'address_detail',
    ] as const;

    const mediaStepFields = ['sale_status', 'cover_image'] as const;

    const areaOptions = ref<any[]>([]);

    const ensureAreaRootLoaded = async () => {
      if (areaOptions.value.length > 0) return;
      const list = await getAreaMoreList({ pid: 0 });
      areaOptions.value = (list || []).map((it: any) => ({
        ...it,
        value: it.id,
        label: it.name,
        // Arco Cascader 懒加载需要明确 isLeaf，否则会被当成叶子节点无法展开
        isLeaf: it.isLeaf === true,
        // 关键：非叶子节点不要预置 children: []，否则不会触发 load-more
        children: undefined,
      }));
    };

    const loadAreaMore = async (option: any, done: (children?: any[]) => void) => {
      try {
        const pid = option?.id ?? option?.value ?? option;
        const list = await getAreaMoreList({ pid });
        const children = (list || []).map((it: any) => ({
          ...it,
          value: it.id,
          label: it.name,
          isLeaf: it.isLeaf === true,
          children: undefined,
        }));
        // 同步写回 option.children，便于后续按 ID 回显/拼接地址时能拿到完整树
        if (option && typeof option === 'object') option.children = children;
        done(children);
      } catch (e) {
        done([]);
      }
    };

    const ensureChildrenLoaded = async (parent: any) => {
      if (!parent) return;
      if (Array.isArray(parent.children) && parent.children.length > 0) return;
      const list = await getAreaMoreList({ pid: parent.id });
      parent.children = (list || []).map((it: any) => ({
        ...it,
        value: it.id,
        label: it.name,
        isLeaf: it.isLeaf === true,
        children: undefined,
      }));
    };

    const findByName = (list: any[], name: string) => {
      if (!Array.isArray(list)) return null;
      return list.find((it) => it?.name === name || it?.shortname === name || it?.mergename === name) || null;
    };

    const findByNameFuzzy = (list: any[], name: string) => {
      if (!Array.isArray(list)) return null;
      const key = (name || '').toString().trim();
      if (!key) return null;
      // 优先精确匹配
      const exact = list.find((it) => it?.name === key || it?.shortname === key || it?.mergename === key);
      if (exact) return exact;
      // 兼容 “天津” vs “天津市”
      const fuzzy = list.find((it) => {
        const n = (it?.name || '').toString();
        const sn = (it?.shortname || '').toString();
        const mn = (it?.mergename || '').toString();
        return n.includes(key) || sn.includes(key) || mn.includes(key) || key.includes(n) || key.includes(sn);
      });
      return fuzzy || null;
    };

    const initAreaFromAddress = async (address: string) => {
      const raw = (address || '').toString().trim();
      if (!raw) return;
      const parts = raw.split(/\s+/).filter(Boolean);
      // 期望格式：省 市 区 详细地址...
      if (parts.length < 4) {
        // 兼容旧数据：如 “天津 双榆树北里2号楼”
        // 尝试用第一个词匹配省份/直辖市，并自动补齐市（同名）
        const first = parts[0] || '';
        await ensureAreaRootLoaded();
        const p = findByNameFuzzy(areaOptions.value, first);
        if (!p) {
          formData.value.address_detail = raw;
          return;
        }
        await ensureChildrenLoaded(p);
        const citySame = findByNameFuzzy(p.children || [], first) || findByNameFuzzy(p.children || [], `${first}市`);
        if (citySame) {
          formData.value.area_ids = [p.id, citySame.id];
        } else {
          formData.value.area_ids = [p.id];
        }
        await ensureAreaPathLoaded(formData.value.area_ids);
        formData.value.address_detail = parts.slice(1).join(' ') || '';
        return;
      }
      const [pName, cName, dName] = parts;
      const detail = parts.slice(3).join(' ');
      await ensureAreaRootLoaded();
      const p = findByNameFuzzy(areaOptions.value, pName);
      if (!p) {
        formData.value.address_detail = raw;
        return;
      }
      await ensureChildrenLoaded(p);
      const c = findByNameFuzzy(p.children || [], cName);
      if (!c) {
        formData.value.address_detail = raw;
        return;
      }
      await ensureChildrenLoaded(c);
      const d = findByNameFuzzy(c.children || [], dName);
      if (!d) {
        // 区没匹配到时：先回显到市级，剩余当详细地址
        formData.value.area_ids = [p.id, c.id];
        await ensureAreaPathLoaded(formData.value.area_ids);
        formData.value.address_detail = [dName, detail].filter(Boolean).join(' ');
        return;
      }
      formData.value.area_ids = [p.id, c.id, d.id];
      await ensureAreaPathLoaded(formData.value.area_ids);
      formData.value.address_detail = detail;
    };

    const getAreaLabelsByIds = (ids: any[]): string[] => {
      const out: string[] = [];
      let list: any[] = areaOptions.value;
      for (const id of ids || []) {
        const node = Array.isArray(list)
          ? list.find((it) => String(it.value ?? it.id) === String(id))
          : null;
        if (!node) break;
        out.push(node.label ?? node.name);
        list = node.children || [];
      }
      return out;
    };

    const ensureAreaPathLoaded = async (ids: any[]) => {
      let list: any[] = areaOptions.value;
      for (let i = 0; i < (ids || []).length; i++) {
        const id = ids[i];
        const node = Array.isArray(list)
          ? list.find((it) => String(it.value ?? it.id) === String(id))
          : null;
        if (!node) return;
        const hasNext = i < ids.length - 1;
        if (hasNext) {
          // 如果需要下一级但 children 未加载，则加载
          if (!Array.isArray(node.children) || node.children.length === 0) {
            const children = await getAreaMoreList({ pid: node.id });
            node.children = (children || []).map((it: any) => ({
              ...it,
              value: it.id,
              label: it.name,
              isLeaf: it.isLeaf === true,
              children: undefined,
            }));
          }
          list = node.children || [];
        }
      }
    };

    const buildFullAddress = () => {
      const ids = formData.value.area_ids || [];
      const names = getAreaLabelsByIds(ids);
      const detail = (formData.value.address_detail || '').toString().trim();
      // 用空格分隔，便于编辑回显解析
      const head = names.join(' ');
      return [head, detail].filter(Boolean).join(' ').trim();
    };

    const normalizeTagsArray = (v: any): string[] => {
      if (!v) return [];
      if (Array.isArray(v)) return v.map((i) => String(i).trim()).filter(Boolean);
      if (typeof v === 'string') {
        return v
          .split(',')
          .map((i) => i.trim())
          .filter(Boolean);
      }
      return [];
    };

    const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
      formRef.value?.resetFields();
      renovationFormRef.value?.resetFields();
      setModalProps({ confirmLoading: false });
      isUpdate.value = !!data?.isUpdate;
      activeTab.value = 'basic';
      wizardStep.value = 0;
      await Promise.all([ensureAreaRootLoaded(), ensureFormOptionsLoaded()]);

      if (unref(isUpdate)) {
        formData.value = cloneDeep(baseData);
        setLoading(true);
        try {
          // 1. 获取基本详情
          const contentData = data?.record ?? (await getContent({ id: data?.id }));
          Object.assign(formData.value, contentData);

          // 2. 数值转换
          formData.value.price = toNumber(contentData.price);
          formData.value.area = toNumber(contentData.area);
          formData.value.rooms = toNumber(contentData.rooms);
          formData.value.halls = toNumber(contentData.halls);
          formData.value.bathrooms = toNumber(contentData.bathrooms);
          formData.value.total_floors = toNumber(contentData.total_floors);
          formData.value.commission_rate = toNumber(contentData.commission_rate);
          formData.value.commission_reward = toNumber(contentData.commission_reward);
          formData.value.weigh = toNumber(contentData.weigh);
          formData.value.hot_status = toNumber(contentData.hot_status);
          formData.value.status = toNumber(contentData.status);
          formData.value.has_smart_lock = toNumber(contentData.has_smart_lock);
          formData.value.receiver_price = toNumber(contentData.receiver_price);
          formData.value.build_year = contentData.build_year ? String(contentData.build_year) : '';
          // 3. tags 兼容（列表接口可能返回逗号分隔字符串）
          formData.value.tags = normalizeTagsArray((contentData as any)?.tags);
          
          // 4. 回显地址到 Cascader
          if (contentData.address) {
             await initAreaFromAddress(contentData.address);
             // 强制刷新 Cascader 显示（避免偶发“显示为 ID”）
             areaCascaderKey.value += 1;
          }
          
          // 5. 获取装修信息
          try {
             const rData = await getRenovation({ property_id: contentData.id });
             if (rData) {
               Object.assign(renovationData.value, rData);
               // 修正进度类型
               if (typeof renovationData.value.progress_percentage !== 'number') {
                  renovationData.value.progress_percentage = parseInt(String(renovationData.value.progress_percentage || 0));
               }
             } else {
               renovationData.value = cloneDeep(baseRenovationData);
             }
          } catch(e) {
             // 没数据
             renovationData.value = cloneDeep(baseRenovationData);
          }

        } catch (e) {
          Message.error('加载房源详情失败');
        }
        setLoading(false);
      } else {
        formData.value = cloneDeep(baseData);
        renovationData.value = cloneDeep(baseRenovationData);
      }

      sanitizeFormSelectValues();
    });

    const getTitle = computed(() => (unref(isUpdate) ? '编辑房源' : '新增房源'));

    const onHeightChange = (val: any) => {
      // 跟随弹窗可用高度，尽量做到“一个窗口内完整展示”
      windHeight.value = val;
      if (typeof val === 'number' && val > 0) modelHeight.value = val;
    };

    // 保存基本信息
    const buildPostData = () => {
      const postData = cloneDeep(formData.value);
      postData.weigh = toNumber(postData.weigh);
      postData.build_year = postData.build_year ? parseInt(postData.build_year) : 0;
      // 合成地址：省 市 区 详细地址
      postData.address = buildFullAddress();
      // 可选金额：undefined/空字符串时按 null 处理，避免默认为 0
      if (postData.receiver_price === undefined || postData.receiver_price === '') {
        postData.receiver_price = null;
      }
      return postData;
    };

    const validateBasicStep = async () => {
      formError.value = '';
      // 只校验“必填基础”字段，避免后续扩展时其它步骤必填项提前阻塞
      const inst: any = formRef.value as any;
      const err = inst?.validateField ? await inst.validateField(basicStepFields as unknown as string[]) : await inst?.validate?.();
      if (err) {
        const msg = getFirstValidateMessage(err);
        formError.value = msg;
        const firstField = Object.keys(err || {})[0];
        if (firstField && typeof inst?.scrollToField === 'function') {
          try {
            inst.scrollToField(firstField);
          } catch {}
        }
        Message.error(msg);
        return false;
      }
      return true;
    };

    const validateMediaStep = async () => {
      formError.value = '';
      const inst: any = formRef.value as any;
      const err = inst?.validateField ? await inst.validateField(mediaStepFields as unknown as string[]) : await inst?.validate?.();
      if (err) {
        const msg = getFirstValidateMessage(err);
        formError.value = msg;
        const firstField = Object.keys(err || {})[0];
        if (firstField && typeof inst?.scrollToField === 'function') {
          try {
            inst.scrollToField(firstField);
          } catch {}
        }
        Message.error(msg);
        return false;
      }
      return true;
    };

    // 第一步：创建草稿（新增模式拿到ID），不关闭弹窗
    const handleSaveDraft = async () => {
      const ok = await validateBasicStep();
      if (!ok) return false;

      // 编辑模式已有ID：仅通过校验即可进入下一步
      if (formData.value.id) return true;

      try {
        setBasicLoading(true);
        const postData = buildPostData();
        const addId: any = await save(postData);
        if (addId) {
          formData.value.id = Number(addId);
          // 生成ID后即进入“编辑”语义，便于后续功能（如装修进度）解锁
          isUpdate.value = true;
        }
        Message.success('草稿已创建，可继续完善信息');
        return true;
      } catch (e: any) {
        Message.error(e?.message || '保存失败');
        return false;
      } finally {
        setBasicLoading(false);
      }
    };

    const handleWizardPrev = () => {
      if (wizardStep.value <= 0) return;
      formError.value = '';
      wizardStep.value -= 1;
    };

    const handleWizardNext = async () => {
      if (wizardStep.value === 0) {
        const ok = await handleSaveDraft();
        if (!ok) return;
      }
      if (wizardStep.value === 2) {
        const ok = await validateMediaStep();
        if (!ok) return;
      }
      formError.value = '';
      if (wizardStep.value < 3) wizardStep.value += 1;
    };

    const handleWizardFinish = async () => {
      // 兜底：如果用户未生成草稿，先按第一步校验/创建
      if (!formData.value.id) {
        const ok = await handleSaveDraft();
        if (!ok) return;
      }
      // 完成前兜底校验：基础必填 + 发布信息
      {
        const inst: any = formRef.value as any;
        const err = inst?.validateField
          ? await inst.validateField([...basicStepFields, ...mediaStepFields] as unknown as string[])
          : await inst?.validate?.();
        if (err) {
          const msg = getFirstValidateMessage(err);
          formError.value = msg;
          const firstField = Object.keys(err || {})[0];
          if (firstField && typeof inst?.scrollToField === 'function') {
            try {
              inst.scrollToField(firstField);
            } catch {}
          }
          Message.error(msg);
          return;
        }
      }
      try {
        setBasicLoading(true);
        formError.value = '';
        const postData = buildPostData();
        const ret: any = await save(postData);
        if (!formData.value.id && ret) formData.value.id = Number(ret);
        Message.success('保存成功');
        emit('success');
        closeModal();
      } catch (e: any) {
        Message.error(e?.message || '保存失败');
      } finally {
        setBasicLoading(false);
      }
    };

    // 编辑模式：任意步骤都可保存（不跳转、不关闭）
    const handleQuickSave = async () => {
      if (basicLoading.value) return;
      // 兜底：没有ID则先创建草稿
      if (!formData.value.id) {
        const ok = await handleSaveDraft();
        if (!ok) return;
      }
      const inst: any = formRef.value as any;
      const err = inst?.validateField
        ? await inst.validateField([...basicStepFields, ...mediaStepFields] as unknown as string[])
        : await inst?.validate?.();
      if (err) {
        const msg = getFirstValidateMessage(err);
        formError.value = msg;
        const firstField = Object.keys(err || {})[0];
        if (firstField && typeof inst?.scrollToField === 'function') {
          try {
            inst.scrollToField(firstField);
          } catch {}
        }
        Message.error(msg);
        return;
      }
      try {
        setBasicLoading(true);
        const postData = buildPostData();
        await save(postData);
        Message.success('保存成功');
        emit('success');
      } catch (e: any) {
        Message.error(e?.message || '保存失败');
      } finally {
        setBasicLoading(false);
      }
    };

    const handleWizardJump = async (target: number) => {
      if (basicLoading.value) return;
      if (target === wizardStep.value) return;
      // 进入非第 0 步前，确保已创建草稿（用于图片/装修/门锁等绑定）
      if (target > 0 && !formData.value.id) {
        const ok = await handleSaveDraft();
        if (!ok) {
          wizardStep.value = 0;
          return;
        }
      }
      formError.value = '';
      wizardStep.value = target;
    };

    watch(
      () => wizardStep.value,
      () => {
        formError.value = '';
      }
    );

    // 保存装修信息
    const handleSaveRenovation = async () => {
      if (!formData.value.id) {
         Message.warning('请先保存基本信息');
         return;
      }
      try {
        const err = await renovationFormRef.value?.validate();
        if (err) {
           Message.error(getFirstValidateMessage(err));
           return;
        }

        setRenovationLoading(true);
        const postData = cloneDeep(renovationData.value);
        postData.property_id = formData.value.id; // bind property_id
        
        // 日期处理
        if (postData.renovation_status === 'none') {
           postData.start_date = null;
           postData.estimated_finish_date = null;
           postData.actual_finish_date = null;
           postData.current_stage = '';
           postData.progress_percentage = 0;
        }
        
        await saveRenovation(postData);
        Message.success('装修记录已更新');
      } catch(e: any) {
         Message.error(e?.message || '保存失败');
      } finally {
         setRenovationLoading(false);
      }
    };

    return {
      registerModal,
      closeModal,
      onHeightChange,
      windHeight,
      modelHeight,
      getTitle,
      loading,
      basicLoading,
      renovationLoading,
      formRef,
      renovationFormRef,
      formData,
      renovationData,
      activeTab,
      areaCascaderKey,
      areaOptions,
      loadAreaMore,
      // helpers
      requiredRule,
      cascaderRule,
      positiveNumberRule,
      getSaleStatusLabel,
      getSaleStatusBadge,
      saleStatusOptions,
      priceUnitOptions,
      floorLevelOptions,
      orientationOptions,
      propertyTypeOptions,
      decorationTypeOptions,
      renovationStageOptions,
      wizardStep,
      mapPickerVisible,
      openMapPicker,
      clearLatLng,
      handleMapPicked,
      goToBasicTab,
      handleQuickSave,
      handleWizardPrev,
      handleWizardNext,
      handleWizardFinish,
      handleWizardJump,
      handleSaveRenovation,
      formError,
      // computed
      showRenovationProgress,
    };
  },
});
</script>

<style lang="less" scoped>
/* Pro Modal Styling */
.add-form-container {
   background: var(--color-bg-2);
   padding: 0;
   display: flex;
   flex-direction: column;
}

.coord-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}
.coord-tip {
  font-size: 12px;
  color: var(--color-text-3);
}

.pro-tabs {
  flex: 1;
  min-height: 0;
}

.pro-modal {
   // Global overrides for modal if needed, usually handled by attrs
}

.pro-tabs {
   :deep(.arco-tabs-nav-tab) {
      justify-content: center; 
      padding-top: 12px;
   }
   :deep(.arco-tabs-tab) {
      font-size: 15px;
      padding: 8px 24px;
   }
   
   .tab-label {
      display: flex;
      align-items: center;
      gap: 6px;
   }
}

.form-content {
   padding: 6px 24px 14px;
}

.form-error-alert {
  margin: 0 0 14px;
}

.form-shell .form-content {
  padding: 0 0 12px;
}

.wizard-layout {
  display: flex;
  flex-direction: column;
}

.form-shell {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr) 340px;
  gap: 20px;
  padding: 6px 24px 0;
  align-items: start;
}

.form-nav {
  position: sticky;
  top: 12px;
  align-self: start;
}

.form-main {
  min-width: 0;
}

.form-side {
  position: sticky;
  top: 12px;
  align-self: start;
}

@media (max-width: 1200px) {
  .form-shell {
    grid-template-columns: minmax(0, 1fr);
    padding: 8px 20px 0;
  }

  .form-nav {
    position: relative;
    top: auto;
  }

  .form-side {
    position: relative;
    top: auto;
  }

}

.pro-sticky-footer {
  position: sticky;
  bottom: 0;
  z-index: 6;
  padding: 14px 0 18px;
  display: flex;
  justify-content: center;
  background: linear-gradient(
    to bottom,
    rgba(255, 255, 255, 0),
    rgba(255, 255, 255, 0.88) 34%,
    rgba(255, 255, 255, 0.98)
  );
  border-top: 1px solid rgba(0, 0, 0, 0.02);
  backdrop-filter: blur(8px);
}

.footer-card {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid var(--color-border-2);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.06);
}

.footer-actions :deep(.arco-btn) {
  min-width: 108px;
}

@media (max-width: 768px) {
  .form-content {
    padding: 6px 16px 14px;
  }
}

@media (prefers-reduced-motion: reduce) {
  * {
    transition: none !important;
    animation: none !important;
  }
}

.side-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.side-card {
  background: var(--color-bg-1);
  border: 1px solid var(--color-border-2);
  border-radius: 12px;
  padding: 14px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.03);
}

.side-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--color-text-1);
  margin-bottom: 8px;
  letter-spacing: 0.02em;
}

.nav-card {
  .nav-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    margin-bottom: 8px;
  }

  .nav-progress {
    margin: 4px 0 10px;
  }

  .nav-tip {
    margin-top: 10px;
    font-size: 12px;
    color: var(--color-text-3);
    line-height: 1.5;
  }

  :deep(.arco-steps-item) {
    cursor: pointer;
  }
  :deep(.arco-steps-item-title) {
    font-weight: 600;
  }
  :deep(.arco-steps-item-description) {
    font-size: 12px;
  }
}

.kv {
  display: grid;
  gap: 8px;

  .kv-row {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 10px;
    min-width: 0;

    .k {
      color: var(--color-text-3);
      font-size: 12px;
      white-space: nowrap;
    }
    .v {
      color: var(--color-text-1);
      font-size: 12px;
      font-weight: 600;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      text-align: right;
    }
  }
}

.tag-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.pill-tag {
  border-radius: 999px;
  padding: 0 10px;
  border-color: var(--color-border-2);
  background: rgba(var(--primary-6), 0.08);
  color: rgb(var(--primary-6));

  :deep(.arco-tag-content) {
    max-width: 160px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}


.inline-section-head {
  margin: 2px 0 10px;
  padding: 10px 12px;
  border: 1px dashed var(--color-border-2);
  border-radius: 10px;
  background: var(--color-fill-1);

  .inline-title {
    font-size: 13px;
    font-weight: 700;
    color: var(--color-text-1);
    line-height: 1.2;
    margin-bottom: 2px;
    letter-spacing: 0.02em;
  }

  .inline-sub {
    font-size: 12px;
    color: var(--color-text-3);
    line-height: 1.4;
  }
}

.form-section {
  margin-bottom: 16px;
  background: var(--color-bg-1);
  border-radius: 8px;
  
  .section-header {
    margin-bottom: 10px;
    padding-left: 8px;
    border-left: 4px solid rgb(var(--primary-6));
    
    .title {
      font-size: 16px;
      font-weight: 600;
      color: var(--color-text-1);
      margin-bottom: 2px;
      line-height: 1.2;
    }
    .sub {
      font-size: 12px;
      color: var(--color-text-3);
    }
  }
  
  .section-body {
     // padding: 12px;
     // opacity: background/border if wanted
  }
}

.pro-form {
  :deep(.arco-form-item) {
    margin-bottom: 12px;
  }

  :deep(.arco-form-item-label-col) {
    padding-bottom: 4px;
  }
}

/* Capsule Input Group */
.capsule-input-group {
   display: flex;
   align-items: center;
   border: 1px solid var(--color-border-2);
   border-radius: 4px; /* match arco input radius */
   overflow: hidden;
   
   .cap-input {
      flex: 1;
      :deep(.arco-input-wrapper) {
         border: none;
         border-radius: 0;
         padding-left: 8px; 
         padding-right: 8px;
         background: var(--color-fill-1);
         &:focus-within { background: #fff; box-shadow: none; }
      }
   }
   
   .cap-divider {
      width: 1px;
      height: 20px;
      background: var(--color-border-2);
   }
}

.mb-0 { :deep(.arco-form-item) { margin-bottom: 0; } }

.switch-tip {
   font-size: 12px;
   color: var(--color-text-3);
   margin-left: 8px;
}

.final-summary {
  margin-top: 6px;

  .strong {
    font-weight: 600;
    color: var(--color-text-1);
  }

  .summary-sep {
    margin: 0 10px;
    color: var(--color-text-3);
  }
}
</style>
