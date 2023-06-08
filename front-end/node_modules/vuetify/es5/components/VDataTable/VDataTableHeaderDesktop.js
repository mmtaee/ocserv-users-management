"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

var _mixins = _interopRequireDefault(require("../../util/mixins"));

var _header = _interopRequireDefault(require("./mixins/header"));

var _helpers = require("../../util/helpers");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _toConsumableArray(arr) { return _arrayWithoutHoles(arr) || _iterableToArray(arr) || _unsupportedIterableToArray(arr) || _nonIterableSpread(); }

function _nonIterableSpread() { throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _iterableToArray(iter) { if (typeof Symbol !== "undefined" && Symbol.iterator in Object(iter)) return Array.from(iter); }

function _arrayWithoutHoles(arr) { if (Array.isArray(arr)) return _arrayLikeToArray(arr); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

var _default = (0, _mixins.default)(_header.default).extend({
  name: 'v-data-table-header-desktop',
  methods: {
    genGroupByToggle: function genGroupByToggle(header) {
      var _this = this;

      return this.$createElement('span', {
        on: {
          click: function click(e) {
            e.stopPropagation();

            _this.$emit('group', header.value);
          }
        }
      }, ['group']);
    },
    getAria: function getAria(beingSorted, isDesc) {
      var _this2 = this;

      var $t = function $t(key) {
        return _this2.$vuetify.lang.t("$vuetify.dataTable.ariaLabel.".concat(key));
      };

      var ariaSort = 'none';
      var ariaLabel = [$t('sortNone'), $t('activateAscending')];

      if (!beingSorted) {
        return {
          ariaSort: ariaSort,
          ariaLabel: ariaLabel.join(' ')
        };
      }

      if (isDesc) {
        ariaSort = 'descending';
        ariaLabel = [$t('sortDescending'), $t(this.options.mustSort ? 'activateAscending' : 'activateNone')];
      } else {
        ariaSort = 'ascending';
        ariaLabel = [$t('sortAscending'), $t('activateDescending')];
      }

      return {
        ariaSort: ariaSort,
        ariaLabel: ariaLabel.join(' ')
      };
    },
    genHeader: function genHeader(header) {
      var _this3 = this;

      var data = {
        attrs: {
          role: 'columnheader',
          scope: 'col',
          'aria-label': header.text || ''
        },
        style: {
          width: (0, _helpers.convertToUnit)(header.width),
          minWidth: (0, _helpers.convertToUnit)(header.width)
        },
        class: ["text-".concat(header.align || 'start')].concat(_toConsumableArray((0, _helpers.wrapInArray)(header.class)), [header.divider && 'v-data-table__divider']),
        on: {}
      };
      var children = [];

      if (header.value === 'data-table-select' && !this.singleSelect) {
        return this.$createElement('th', data, [this.genSelectAll()]);
      }

      children.push(this.$scopedSlots.hasOwnProperty(header.value) ? this.$scopedSlots[header.value]({
        header: header
      }) : this.$createElement('span', [header.text]));

      if (!this.disableSort && (header.sortable || !header.hasOwnProperty('sortable'))) {
        data.on.click = function () {
          return _this3.$emit('sort', header.value);
        };

        var sortIndex = this.options.sortBy.findIndex(function (k) {
          return k === header.value;
        });
        var beingSorted = sortIndex >= 0;
        var isDesc = this.options.sortDesc[sortIndex];
        data.class.push('sortable');

        var _this$getAria = this.getAria(beingSorted, isDesc),
            ariaLabel = _this$getAria.ariaLabel,
            ariaSort = _this$getAria.ariaSort;

        data.attrs['aria-label'] += "".concat(header.text ? ': ' : '').concat(ariaLabel);
        data.attrs['aria-sort'] = ariaSort;

        if (beingSorted) {
          data.class.push('active');
          data.class.push(isDesc ? 'desc' : 'asc');
        }

        if (header.align === 'end') children.unshift(this.genSortIcon());else children.push(this.genSortIcon());

        if (this.options.multiSort && beingSorted) {
          children.push(this.$createElement('span', {
            class: 'v-data-table-header__sort-badge'
          }, [String(sortIndex + 1)]));
        }
      }

      if (this.showGroupBy && header.groupable !== false) children.push(this.genGroupByToggle(header));
      return this.$createElement('th', data, children);
    }
  },
  render: function render() {
    var _this4 = this;

    return this.$createElement('thead', {
      staticClass: 'v-data-table-header'
    }, [this.$createElement('tr', this.headers.map(function (header) {
      return _this4.genHeader(header);
    }))]);
  }
});

exports.default = _default;
//# sourceMappingURL=VDataTableHeaderDesktop.js.map