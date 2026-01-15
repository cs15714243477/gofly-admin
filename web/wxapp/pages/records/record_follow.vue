<template>
	<view class="record-page">
		<TopHeader title="关注记录" />
		<scroll-view scroll-y class="list">
			<view class="card" v-for="(it, idx) in items" :key="idx">
				<view class="row">
					<text class="title">{{ it.title }}</text>
					<text class="tag">关注</text>
				</view>
				<view class="sub">{{ it.time }}</view>
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
					{ title: '阳光花园三期 2室1厅', time: '2026-01-07 14:20' },
					{ title: '金地名津 精装大三房', time: '2026-01-06 19:02' },
					{ title: '万科四季花城', time: '2026-01-05 10:11' }
				]
			}
		}
	}
</script>

<style lang="scss">
	.record-page {
		height: 100vh;
		display: flex;
		flex-direction: column;
		background: #f6f7f8;
	}
	.list {
		flex: 1;
		overflow: hidden;
		padding: 24rpx 16rpx;
		box-sizing: border-box;
	}
	.card {
		background: #fff;
		border-radius: 24rpx;
		padding: 20rpx;
		border: 1px solid #f1f5f9;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.02);
		margin-bottom: 16rpx;
	}
	.row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16rpx;
	}
	.title {
		font-size: 30rpx;
		font-weight: 800;
		color: #0f172a;
		flex: 1;
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	.tag {
		font-size: 22rpx;
		font-weight: 800;
		color: #2563eb;
		background: rgba(37, 99, 235, 0.12);
		border: 1rpx solid rgba(37, 99, 235, 0.18);
		padding: 8rpx 14rpx;
		border-radius: 999rpx;
	}
	.sub {
		margin-top: 10rpx;
		font-size: 24rpx;
		color: #64748b;
	}
	.bottom-spacer {
		height: 40rpx;
	}
</style>

