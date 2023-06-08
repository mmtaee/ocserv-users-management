import VIcon from '../../VIcon';
import VSimpleCheckbox from '../../VCheckbox/VSimpleCheckbox';
import ripple from '../../../directives/ripple';
import mixins from '../../../util/mixins';
export default mixins().extend({
  // https://github.com/vuejs/vue/issues/6872
  directives: {
    ripple
  },
  props: {
    headers: {
      type: Array,
      default: () => []
    },
    options: {
      type: Object,
      default: () => ({
        page: 1,
        itemsPerPage: 10,
        sortBy: [],
        sortDesc: [],
        groupBy: [],
        groupDesc: [],
        multiSort: false,
        mustSort: false
      })
    },
    checkboxColor: String,
    sortIcon: {
      type: String,
      default: '$sort'
    },
    everyItem: Boolean,
    someItems: Boolean,
    showGroupBy: Boolean,
    singleSelect: Boolean,
    disableSort: Boolean
  },
  methods: {
    genSelectAll() {
      var _a;

      const data = {
        props: {
          value: this.everyItem,
          indeterminate: !this.everyItem && this.someItems,
          color: (_a = this.checkboxColor) !== null && _a !== void 0 ? _a : ''
        },
        on: {
          input: v => this.$emit('toggle-select-all', v)
        }
      };

      if (this.$scopedSlots['data-table-select']) {
        return this.$scopedSlots['data-table-select'](data);
      }

      return this.$createElement(VSimpleCheckbox, {
        staticClass: 'v-data-table__checkbox',
        ...data
      });
    },

    genSortIcon() {
      return this.$createElement(VIcon, {
        staticClass: 'v-data-table-header__icon',
        props: {
          size: 18
        }
      }, [this.sortIcon]);
    }

  }
});
//# sourceMappingURL=header.js.map