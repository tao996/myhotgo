import { h } from 'vue';
import { NTag } from 'naive-ui';

export const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 60,
  },
  {
    title: '岗位',
    key: 'name',
    width: 120,
    render(row) {
      return h(
        NTag,
        {
          type: 'info',
        },
        {
          default: () => row.name,
        }
      );
    },
  },
  {
    title: '岗位编码',
    key: 'code',
    width: 120,
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      return h(
        NTag,
        {
          style: {
            marginRight: '6px',
          },
          type: row.status == 1 ? 'info' : 'error',
          bordered: false,
        },
        {
          default: () => (row.status == 1 ? '正常' : '已禁用'),
        }
      );
    },
  },
  {
    title: '备注',
    key: 'remark',
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 180,
  },
];
