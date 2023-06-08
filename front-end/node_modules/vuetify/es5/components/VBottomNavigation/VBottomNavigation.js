"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

require("../../../src/components/VBottomNavigation/VBottomNavigation.sass");

var _applicationable = _interopRequireDefault(require("../../mixins/applicationable"));

var _buttonGroup = _interopRequireDefault(require("../../mixins/button-group"));

var _colorable = _interopRequireDefault(require("../../mixins/colorable"));

var _measurable = _interopRequireDefault(require("../../mixins/measurable"));

var _proxyable = _interopRequireDefault(require("../../mixins/proxyable"));

var _scrollable = _interopRequireDefault(require("../../mixins/scrollable"));

var _themeable = _interopRequireDefault(require("../../mixins/themeable"));

var _toggleable = require("../../mixins/toggleable");

var _mixins = _interopRequireDefault(require("../../util/mixins"));

var _console = require("../../util/console");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

var _default = (0, _mixins.default)((0, _applicationable.default)('bottom', ['height', 'inputValue']), _colorable.default, _measurable.default, (0, _toggleable.factory)('inputValue'), _proxyable.default, _scrollable.default, _themeable.default
/* @vue/component */
).extend({
  name: 'v-bottom-navigation',
  props: {
    activeClass: {
      type: String,
      default: 'v-btn--active'
    },
    backgroundColor: String,
    grow: Boolean,
    height: {
      type: [Number, String],
      default: 56
    },
    hideOnScroll: Boolean,
    horizontal: Boolean,
    inputValue: {
      type: Boolean,
      default: true
    },
    mandatory: Boolean,
    shift: Boolean,
    tag: {
      type: String,
      default: 'div'
    }
  },
  data: function data() {
    return {
      isActive: this.inputValue
    };
  },
  computed: {
    canScroll: function canScroll() {
      return _scrollable.default.options.computed.canScroll.call(this) && (this.hideOnScroll || !this.inputValue);
    },
    classes: function classes() {
      return {
        'v-bottom-navigation--absolute': this.absolute,
        'v-bottom-navigation--grow': this.grow,
        'v-bottom-navigation--fixed': !this.absolute && (this.app || this.fixed),
        'v-bottom-navigation--horizontal': this.horizontal,
        'v-bottom-navigation--shift': this.shift
      };
    },
    styles: function styles() {
      return _objectSpread(_objectSpread({}, this.measurableStyles), {}, {
        transform: this.isActive ? 'none' : 'translateY(100%)'
      });
    }
  },
  watch: {
    canScroll: 'onScroll'
  },
  created: function created() {
    /* istanbul ignore next */
    if (this.$attrs.hasOwnProperty('active')) {
      (0, _console.breaking)('active.sync', 'value or v-model', this);
    }
  },
  methods: {
    thresholdMet: function thresholdMet() {
      if (this.hideOnScroll) {
        this.isActive = !this.isScrollingUp || this.currentScroll > this.computedScrollThreshold;
        this.$emit('update:input-value', this.isActive);
      }

      if (this.currentThreshold < this.computedScrollThreshold) return;
      this.savedScroll = this.currentScroll;
    },
    updateApplication: function updateApplication() {
      return this.$el ? this.$el.clientHeight : 0;
    },
    updateValue: function updateValue(val) {
      this.$emit('change', val);
    }
  },
  render: function render(h) {
    var data = this.setBackgroundColor(this.backgroundColor, {
      staticClass: 'v-bottom-navigation',
      class: this.classes,
      style: this.styles,
      props: {
        activeClass: this.activeClass,
        mandatory: Boolean(this.mandatory || this.value !== undefined),
        tag: this.tag,
        value: this.internalValue
      },
      on: {
        change: this.updateValue
      }
    });

    if (this.canScroll) {
      data.directives = data.directives || [];
      data.directives.push({
        arg: this.scrollTarget,
        name: 'scroll',
        value: this.onScroll
      });
    }

    return h(_buttonGroup.default, this.setTextColor(this.color, data), this.$slots.default);
  }
});

exports.default = _default;
//# sourceMappingURL=VBottomNavigation.js.map