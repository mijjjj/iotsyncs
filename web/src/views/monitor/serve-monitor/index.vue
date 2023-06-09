<template>
  <div>
    <div class="main-container">
      <n-spin :show="loading" description="正在从服务器查询信息...">
        <n-grid responsive="screen" cols="1 s:2 m:4 l:4 xl:4 2xl:4" x-gap="5" y-gap="5">
          <n-grid-item v-for="(item, index) of dataSource.head" :key="index" class="item-wrapper">
            <DataItem :data-model="item" :loading="loading">
              <template v-if="index === 0" #extra="{ extra }">
                <div class="margin-top-lg">
                  <div> {{ extra.data }}</div>
                  <div class="margin-top-sm"> {{ extra.data1 }}</div>
                </div>
              </template>
              <template v-else-if="index === 1" #extra="{ extra }">
                <div class="margin-top" style="position: relative">
                  <div> 已用内存：{{ extra.data }}</div>
                  <div class="margin-top-sm"> 剩余内存：{{ extra.data1 }}</div>
                  <div class="stack-avatar-wrapper"></div>
                </div>
              </template>
              <template v-else-if="index === 2" #extra="{ extra }">
                <n-progress type="line" :percentage="extra.data" />
              </template>
              <template v-else-if="index === 3" #extra>
                <LoadChart ref="mLoadChart" :data-model="dataSource.load" />
              </template>
            </DataItem>
          </n-grid-item>
        </n-grid>

        <n-grid class="mt-2">
          <n-grid-item :span="24">
            <FullYearSalesChart
              ref="fullYearSalesChart"
              :data-model="dataSource.net"
              :loading="loading"
            />
          </n-grid-item>
        </n-grid>

        <n-space vertical style="padding-top: 10px">
          <n-card title="服务器信息">
            <n-descriptions
              label-placement="top"
              bordered
              cols="1 s:1 m:2 l:3 xl:4 2xl:4"
              :label-style="{ 'font-weight': 'bold', 'font-size': '16px' }"
            >
              <n-descriptions-item label="服务器名称">
                {{ dataRunInfo.hostname }}
              </n-descriptions-item>
              <n-descriptions-item label="操作系统"> {{ dataRunInfo.os }}</n-descriptions-item>
              <n-descriptions-item label="服务器IP">
                {{ dataRunInfo.intranet_ip }} /
                {{ dataRunInfo.public_ip }}
              </n-descriptions-item>
              <n-descriptions-item label="系统架构"> {{ dataRunInfo.arch }}</n-descriptions-item>
            </n-descriptions>
          </n-card>
          <n-card title="GO运行信息">
            <n-descriptions
              label-placement="top"
              bordered
              cols="1 s:1 m:2 l:3 xl:4 2xl:4"
              :label-style="{ 'font-weight': 'bold', 'font-size': '16px' }"
            >
              <n-descriptions-item label="语言环境"> {{ dataRunInfo.goName }}</n-descriptions-item>
              <n-descriptions-item label="版本号"> {{ dataRunInfo.version }}</n-descriptions-item>
              <n-descriptions-item label="启动时间">
                {{ dataRunInfo.startTime }}</n-descriptions-item
              >
              <n-descriptions-item label="运行时长">
                {{ formatBefore(new Date(dataRunInfo.startTime)) }}
              </n-descriptions-item>
              <n-descriptions-item label="运行路径"> {{ dataRunInfo.pwd }}</n-descriptions-item>
              <n-descriptions-item label="goroutine数量">
                {{ dataRunInfo.goroutine }}
              </n-descriptions-item>
              <n-descriptions-item label="运行内存"> {{ dataRunInfo.goMem }}</n-descriptions-item>
              <n-descriptions-item label="磁盘占用"> {{ dataRunInfo.goSize }}</n-descriptions-item>
            </n-descriptions>
          </n-card>
        </n-space>
      </n-spin>
    </div>
  </div>
</template>

<script lang="ts">
  import DataItem from './components/DataItem.vue';
  import LoadChart from './components/chart/LoadChart.vue';
  import FullYearSalesChart from './components/chart/FullYearSalesChart.vue';
  import { defineComponent, inject, onMounted, ref, onUpdated } from 'vue';
  import { SocketEnum } from '@/enums/socketEnum';
  import { addOnMessage, sendMsg } from '@/utils/websocket';
  import { formatBefore } from '@/utils/dateUtil';
  import { useDialog, useMessage } from 'naive-ui';

  export default defineComponent({
    name: 'Home',
    components: {
      DataItem,
      LoadChart,
      FullYearSalesChart,
    },
    setup() {
      const dataRunInfo = ref({
        arch: '',
        goMem: '0MB',
        goName: 'Golang',
        goSize: '0MB',
        goroutine: 0,
        hostname: '',
        intranet_ip: '127.0.0.1',
        os: '',
        public_ip: '0.0.0.0',
        pwd: '/',
        rootPath: '/',
        runTime: 0,
        startTime: '',
        version: '',
      });
      const dataSource = ref({
        head: [
          {
            title: 'CPU',
            data: '0%',
            bottomTitle: 'CPU数量',
            totalSum: '',
            iconClass: 'HardwareChip',
            extra: {
              data: '',
              data1: '',
            },
          },
          {
            title: '内存',
            data: '0%',
            bottomTitle: '总内存',
            totalSum: '0GB',
            iconClass: 'AppsSharp',
            extra: {
              data: '0GB',
              data1: '0GB',
            },
          },
          {
            title: '磁盘',
            data: '已用 0GB',
            bottomTitle: '总容量',
            totalSum: '0GB',
            iconClass: 'PieChart',
            extra: {
              data: 0,
            },
          },
          {
            title: '负载',
            data: '0%',
            bottomTitle: '总进程数',
            totalSum: '0个',
            iconClass: 'Analytics',
            extra: {
              data: 80,
            },
          },
        ],
        load: {},
        net: {},
      });

      const message = useMessage();
      const dialog = useDialog();
      const loading = ref(true);
      const onMessageList = inject('onMessageList');

      const onAdminMonitor = (res) => {
        const data = JSON.parse(res.data);
        if (data.event === SocketEnum.EventAdminMonitorRunInfo) {
          loading.value = false;
          if (data.code == SocketEnum.CodeErr) {
            message.error('查询出错:' + data.event);
            return;
          }

          dataRunInfo.value = data.data;
          return;
        }

        if (data.event === SocketEnum.EventAdminMonitorTrends) {
          loading.value = false;
          if (data.code == SocketEnum.CodeErr) {
            message.error('查询出错:' + data.event);
            return;
          }
          dataSource.value = data.data;
          return;
        }
      };

      addOnMessage(onMessageList, onAdminMonitor);

      onMounted(() => {
        getInfo();

        setTimeout(() => {
          if (loading.value) {
            loading.value = false;
            dialog.error({
              title: '错误',
              content: '连接超时，请刷新重试。如仍未解决请检查websocket连接是否正确！',
              positiveText: '确定',
            });
          }
        }, 5000);
      });

      onUpdated(() => {
        // 切换页面后直接出发一次
        if (loading.value === false) {
          sendMsg(SocketEnum.EventAdminMonitorTrends);
          sendMsg(SocketEnum.EventAdminMonitorRunInfo);
        }
      });

      function getInfo() {
        loading.value = true;
        sendMsg(SocketEnum.EventAdminMonitorTrends);
        sendMsg(SocketEnum.EventAdminMonitorRunInfo);

        setInterval(function () {
          sendMsg(SocketEnum.EventAdminMonitorTrends);
        }, 1000 * 2);

        setInterval(function () {
          sendMsg(SocketEnum.EventAdminMonitorRunInfo);
        }, 1000 * 10);
      }

      const mLoadChart = ref<InstanceType<typeof LoadChart>>();
      const fullYearSalesChart = ref<InstanceType<typeof FullYearSalesChart>>();
      const onResize = () => {
        setTimeout(() => {
          mLoadChart.value?.updateChart();
        }, 500);
      };
      const collapse = true;
      onResize();
      return {
        loading,
        collapse,
        mLoadChart,
        fullYearSalesChart,
        dataSource,
        dataRunInfo,
        formatBefore,
      };
    },
  });
</script>

<style lang="less" scoped>
  @media screen and (max-width: 992px) {
    .item-wrapper {
      margin-bottom: 5px;
    }

    .map-margin-tb {
      margin: 5px 0;
    }
  }

  .light {
    .chart-item {
      background-color: #fff;
    }
  }

  .stack-avatar-wrapper {
    position: absolute;
    right: -2%;
    top: 10%;
  }
</style>
