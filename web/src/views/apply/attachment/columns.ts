import { h, ref } from 'vue';
import { NAvatar, NImage, NTag } from 'naive-ui';
import { getFileExt } from '@/utils/urlUtils';
import { Dicts } from '@/api/dict/dict';
import { getOptionLabel, getOptionTag, Options } from '@/utils/hotgo';
import { FormSchema } from '@/components/Form';
import { isNullOrUnDef } from '@/utils/is';
import { errorImg } from '@/utils/hotgo';
export const options = ref<Options>({
  sys_normal_disable: [],
  config_upload_drive: [],
});

export const schemas = ref<FormSchema[]>([
  {
    field: 'drive',
    component: 'NSelect',
    label: '上传驱动',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择上传驱动',
      options: [],
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
  },
  {
    field: 'member_id',
    component: 'NInput',
    label: '用户ID',
    componentProps: {
      placeholder: '请输入用户ID',
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
    rules: [{ message: '请输入用户ID', trigger: ['blur'] }],
  },
  {
    field: 'status',
    component: 'NSelect',
    label: '状态',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择类型',
      options: [],
      onUpdateValue: (e: any) => {
        console.log(e);
      },
    },
  },
]);

export const columns = [
  {
    title: '附件ID',
    key: 'id',
    width: 80,
  },
  {
    title: '应用',
    key: 'appId',
    width: 100,
  },
  {
    title: '用户ID',
    key: 'memberId',
    width: 100,
  },
  {
    title: '驱动',
    key: 'drive',
    render(row) {
      return row.drive;
    },
    width: 100,
  },
  {
    title: '上传类型',
    key: 'kind',
    render(row) {
      return h(
        NTag,
        {
          style: {
            marginRight: '6px',
          },
          type: row.kind == 'images' ? 'success' : 'warning',
          bordered: false,
        },
        {
          default: () => row.kind,
        }
      );
    },
    width: 120,
  },
  {
    title: '文件',
    key: 'fileUrl',
    width: 80,
    render(row) {
      if (row.fileUrl === '') {
        return ``;
      }
      if (row.kind !== 'images') {
        return h(
          NAvatar,
          {
            width: '40px',
            height: '40px',
            'max-width': '100%',
            'max-height': '100%',
          },
          {
            default: () => getFileExt(row.fileUrl),
          }
        );
      }
      return h(NImage, {
        width: 40,
        height: 40,
        src: row.fileUrl,
        onError: errorImg,
        style: {
          width: '40px',
          height: '40px',
          'max-width': '100%',
          'max-height': '100%',
        },
      });
    },
  },
  // {
  //   title: '本地路径',
  //   key: 'path',
  // },
  {
    title: '扩展名',
    key: 'ext',
    width: 80,
  },
  {
    title: '文件大小',
    key: 'sizeFormat',
    width: 100,
  },
  {
    title: '状态',
    key: 'status',
    render(row) {
      if (isNullOrUnDef(row.status)) {
        return ``;
      }
      return h(
        NTag,
        {
          style: {
            marginRight: '6px',
          },
          type: getOptionTag(options.value.sys_normal_disable, row.status),
          bordered: false,
        },
        {
          default: () => getOptionLabel(options.value.sys_normal_disable, row.status),
        }
      );
    },
    width: 100,
  },
  {
    title: '上传时间',
    key: 'createdAt',
    width: 180,
  },
];

async function loadOptions() {
  options.value = await Dicts({
    types: ['sys_normal_disable', 'config_upload_drive'],
  });
  for (const item of schemas.value) {
    switch (item.field) {
      case 'status':
        item.componentProps.options = options.value.sys_normal_disable;
        break;
      case 'drive':
        item.componentProps.options = options.value.config_upload_drive;
        break;
    }
  }
}

await loadOptions();
