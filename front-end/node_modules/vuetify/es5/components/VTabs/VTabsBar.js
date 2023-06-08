"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

var _VSlideGroup = require("../VSlideGroup/VSlideGroup");

var _themeable = _interopRequireDefault(require("../../mixins/themeable"));

var _ssrBootable = _interopRequireDefault(require("../../mixins/ssr-bootable"));

var _mixins = _interopRequireDefault(require("../../util/mixins"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _createForOfIteratorHelper(o, allowArrayLike) { var it; if (typeof Symbol === "undefined" || o[Symbol.iterator] == null) { if (Array.isArray(o) || (it = _unsupportedIterableToArray(o)) || allowArrayLike && o && typeof o.length === "number") { if (it) o = it; var i = 0; var F = function F() {}; return { s: F, n: function n() { if (i >= o.length) return { done: true }; return { done: false, value: o[i++] }; }, e: function e(_e) { throw _e; }, f: F }; } throw new TypeError("Invalid attempt to iterate non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); } var normalCompletion = true, didErr = false, err; return { s: function s() { it = o[Symbol.iterator](); }, n: function n() { var step = it.next(); normalCompletion = step.done; return step; }, e: function e(_e2) { didErr = true; err = _e2; }, f: function f() { try { if (!normalCompletion && it.return != null) it.return(); } finally { if (didErr) throw err; } } }; }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

var _default = (0, _mixins.default)(_VSlideGroup.BaseSlideGroup, _ssrBootable.default, _themeable.default
/* @vue/component */
).extend({
  name: 'v-tabs-bar',
  provide: function provide() {
    return {
      tabsBar: this
    };
  },
  computed: {
    classes: function classes() {
      return _objectSpread(_objectSpread({}, _VSlideGroup.BaseSlideGroup.options.computed.classes.call(this)), {}, {
        'v-tabs-bar': true,
        'v-tabs-bar--is-mobile': this.isMobile,
        // TODO: Remove this and move to v-slide-group
        'v-tabs-bar--show-arrows': this.showArrows
      }, this.themeClasses);
    }
  },
  watch: {
    items: 'callSlider',
    internalValue: 'callSlider',
    $route: 'onRouteChange'
  },
  methods: {
    callSlider: function callSlider() {
      if (!this.isBooted) return;
      this.$emit('call:slider');
    },
    genContent: function genContent() {
      var render = _VSlideGroup.BaseSlideGroup.options.methods.genContent.call(this);

      render.data = render.data || {};
      render.data.staticClass += ' v-tabs-bar__content';
      return render;
    },
    onRouteChange: function onRouteChange(val, oldVal) {
      /* istanbul ignore next */
      if (this.mandatory) return;
      var items = this.items;
      var newPath = val.path;
      var oldPath = oldVal.path;
      var hasNew = false;
      var hasOld = false;

      var _iterator = _createForOfIteratorHelper(items),
          _step;

      try {
        for (_iterator.s(); !(_step = _iterator.n()).done;) {
          var item = _step.value;
          if (item.to === oldPath) hasOld = true;else if (item.to === newPath) hasNew = true;
          if (hasNew && hasOld) break;
        } // If we have an old item and not a new one
        // it's assumed that the user navigated to
        // a path that is not present in the items

      } catch (err) {
        _iterator.e(err);
      } finally {
        _iterator.f();
      }

      if (!hasNew && hasOld) this.internalValue = undefined;
    }
  },
  render: function render(h) {
    var render = _VSlideGroup.BaseSlideGroup.options.render.call(this, h);

    render.data.attrs = {
      role: 'tablist'
    };
    return render;
  }
});

exports.default = _default;
//# sourceMappingURL=VTabsBar.js.map