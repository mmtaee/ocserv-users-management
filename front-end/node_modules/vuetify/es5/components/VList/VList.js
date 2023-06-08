"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

require("../../../src/components/VList/VList.sass");

var _VSheet = _interopRequireDefault(require("../VSheet/VSheet"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _createForOfIteratorHelper(o, allowArrayLike) { var it; if (typeof Symbol === "undefined" || o[Symbol.iterator] == null) { if (Array.isArray(o) || (it = _unsupportedIterableToArray(o)) || allowArrayLike && o && typeof o.length === "number") { if (it) o = it; var i = 0; var F = function F() {}; return { s: F, n: function n() { if (i >= o.length) return { done: true }; return { done: false, value: o[i++] }; }, e: function e(_e) { throw _e; }, f: F }; } throw new TypeError("Invalid attempt to iterate non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); } var normalCompletion = true, didErr = false, err; return { s: function s() { it = o[Symbol.iterator](); }, n: function n() { var step = it.next(); normalCompletion = step.done; return step; }, e: function e(_e2) { didErr = true; err = _e2; }, f: function f() { try { if (!normalCompletion && it.return != null) it.return(); } finally { if (didErr) throw err; } } }; }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

/* @vue/component */
var _default = _VSheet.default.extend().extend({
  name: 'v-list',
  provide: function provide() {
    return {
      isInList: true,
      list: this
    };
  },
  inject: {
    isInMenu: {
      default: false
    },
    isInNav: {
      default: false
    }
  },
  props: {
    dense: Boolean,
    disabled: Boolean,
    expand: Boolean,
    flat: Boolean,
    nav: Boolean,
    rounded: Boolean,
    subheader: Boolean,
    threeLine: Boolean,
    twoLine: Boolean
  },
  data: function data() {
    return {
      groups: []
    };
  },
  computed: {
    classes: function classes() {
      return _objectSpread(_objectSpread({}, _VSheet.default.options.computed.classes.call(this)), {}, {
        'v-list--dense': this.dense,
        'v-list--disabled': this.disabled,
        'v-list--flat': this.flat,
        'v-list--nav': this.nav,
        'v-list--rounded': this.rounded,
        'v-list--subheader': this.subheader,
        'v-list--two-line': this.twoLine,
        'v-list--three-line': this.threeLine
      });
    }
  },
  methods: {
    register: function register(content) {
      this.groups.push(content);
    },
    unregister: function unregister(content) {
      var index = this.groups.findIndex(function (g) {
        return g._uid === content._uid;
      });
      if (index > -1) this.groups.splice(index, 1);
    },
    listClick: function listClick(uid) {
      if (this.expand) return;

      var _iterator = _createForOfIteratorHelper(this.groups),
          _step;

      try {
        for (_iterator.s(); !(_step = _iterator.n()).done;) {
          var group = _step.value;
          group.toggle(uid);
        }
      } catch (err) {
        _iterator.e(err);
      } finally {
        _iterator.f();
      }
    }
  },
  render: function render(h) {
    var data = {
      staticClass: 'v-list',
      class: this.classes,
      style: this.styles,
      attrs: _objectSpread({
        role: this.isInNav || this.isInMenu ? undefined : 'list'
      }, this.attrs$)
    };
    return h(this.tag, this.setBackgroundColor(this.color, data), [this.$slots.default]);
  }
});

exports.default = _default;
//# sourceMappingURL=VList.js.map