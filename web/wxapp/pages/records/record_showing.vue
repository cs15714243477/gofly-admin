<template>
	<view class="record-page">
		<TopHeader title="带看记录" />
		<scroll-view scroll-y class="list">
			<view class="card" v-for="(it, idx) in items" :key="idx">
				<view class="row">
					<text class="title">{{ it.title }}</text>
					<text class="tag">{{ it.type }}</text>
				</view>
				<view class="sub">{{ it.time }} · {{ it.client }}</view>
			</view>
			<view class="bottom-spacer"></view>
		</scroll-view>
		<BottomTabBar active="me" />
	</view>
</template>

<script>
	import TopHeader from '@/components/TopHeader.vue'
	import BottomTabBar from '@/components/BottomTabBar.vue'
	export default {
		components: { TopHeader, BottomTabBar },
		data() {
			return {
				items: [
					{ title: '滨江别墅 独栋花园', time: '2026-01-07 13:40', client: '张先生', type: '线下' },
					{ title: '阳光花园三期 2室', time: '2026-01-06 16:20', client: '王女士', type: '自助开锁' },
					{ title: '金城中心 单身公寓', time: '2026-01-05 10:30', client: '李先生', type: '线上' }
				]
			}
		}
	}
</script>

<style lang="scss">
	.record-page { height: 100vh; display: flex; flex-direction: column; background: #f6f7f8; }
	.list { flex: 1; overflow: hidden; padding: 24rpx 16rpx; box-sizing: border-box; }
	.card { background: #fff; border-radius: 24rpx; padding: 20rpx; border: 1px solid #f1f5f9; box-shadow: 0 2rpx 10rpx rgba(0,0,0,.02); margin-bottom: 16rpx; }
	.row { display: flex; align-items: center; justify-content: space-between; gap: 16rpx; }
	.title { font-size: 30rpx; font-weight: 800; color: #0f172a; flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	.tag { font-size: 22rpx; font-weight: 800; color: #0f172a; background: rgba(148,163,184,.16); border: 1rpx solid rgba(148,163,184,.18); padding: 8rpx 14rpx; border-radius: 999rpx; }
	.sub { margin-top: 10rpx; font-size: 24rpx; color: #64748b; }
	.bottom-spacer { height: 40rpx; }
</style>

