"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

var _vue = _interopRequireDefault(require("vue"));

var _helpers = require("../../util/helpers");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _toConsumableArray(arr) { return _arrayWithoutHoles(arr) || _iterableToArray(arr) || _unsupportedIterableToArray(arr) || _nonIterableSpread(); }

function _nonIterableSpread() { throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _iterableToArray(iter) { if (typeof Symbol !== "undefined" && Symbol.iterator in Object(iter)) return Array.from(iter); }

function _arrayWithoutHoles(arr) { if (Array.isArray(arr)) return _arrayLikeToArray(arr); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

function needsTd(slot) {
  var _a;

  return slot.length !== 1 || !['td', 'th'].includes((_a = slot[0]) === null || _a === void 0 ? void 0 : _a.tag);
}

var _default = _vue.default.extend({
  name: 'row',
  functional: true,
  props: {
    headers: Array,
    index: Number,
    item: Object,
    rtl: Boolean
  },
  render: function render(h, _ref) {
    var props = _ref.props,
        slots = _ref.slots,
        data = _ref.data;
    var computedSlots = slots();
    var columns = props.headers.map(function (header) {
      var children = [];
      var value = (0, _helpers.getObjectValueByPath)(props.item, header.value);
      var slotName = header.value;
      var scopedSlot = data.scopedSlots && data.scopedSlots.hasOwnProperty(slotName) && data.scopedSlots[slotName];
      var regularSlot = computedSlots.hasOwnProperty(slotName) && computedSlots[slotName];

      if (scopedSlot) {
        children.push.apply(children, _toConsumableArray((0, _helpers.wrapInArray)(scopedSlot({
          item: props.item,
          isMobile: false,
          header: header,
          index: props.index,
          value: value
        }))));
      } else if (regularSlot) {
        children.push.apply(children, _toConsumableArray((0, _helpers.wrapInArray)(regularSlot)));
      } else {
        children.push(value == null ? value : String(value));
      }

      var textAlign = "text-".concat(header.align || 'start');
      return needsTd(children) ? h('td', {
        class: [textAlign, header.cellClass, {
          'v-data-table__divider': header.divider
        }]
      }, children) : children;
    });
    return h('tr', data, columns);
  }
});

exports.default = _default;
//# sourceMappingURL=Row.js.map