"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

require("../../../src/components/VStepper/VStepper.sass");

var _VSheet = _interopRequireDefault(require("../VSheet"));

var _registrable = require("../../mixins/registrable");

var _proxyable = _interopRequireDefault(require("../../mixins/proxyable"));

var _mixins = _interopRequireDefault(require("../../util/mixins"));

var _console = require("../../util/console");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

var baseMixins = (0, _mixins.default)(_VSheet.default, (0, _registrable.provide)('stepper'), _proxyable.default);
/* @vue/component */

var _default = baseMixins.extend({
  name: 'v-stepper',
  provide: function provide() {
    return {
      stepClick: this.stepClick,
      isVertical: this.vertical
    };
  },
  props: {
    altLabels: Boolean,
    nonLinear: Boolean,
    flat: Boolean,
    vertical: Boolean
  },
  data: function data() {
    var data = {
      isBooted: false,
      steps: [],
      content: [],
      isReverse: false
    };
    data.internalLazyValue = this.value != null ? this.value : (data[0] || {}).step || 1;
    return data;
  },
  computed: {
    classes: function classes() {
      return _objectSpread({
        'v-stepper--flat': this.flat,
        'v-stepper--is-booted': this.isBooted,
        'v-stepper--vertical': this.vertical,
        'v-stepper--alt-labels': this.altLabels,
        'v-stepper--non-linear': this.nonLinear
      }, _VSheet.default.options.computed.classes.call(this));
    },
    styles: function styles() {
      return _objectSpread({}, _VSheet.default.options.computed.styles.call(this));
    }
  },
  watch: {
    internalValue: function internalValue(val, oldVal) {
      this.isReverse = Number(val) < Number(oldVal);
      oldVal && (this.isBooted = true);
      this.updateView();
    }
  },
  created: function created() {
    /* istanbul ignore next */
    if (this.$listeners.input) {
      (0, _console.breaking)('@input', '@change', this);
    }
  },
  mounted: function mounted() {
    this.updateView();
  },
  methods: {
    register: function register(item) {
      if (item.$options.name === 'v-stepper-step') {
        this.steps.push(item);
      } else if (item.$options.name === 'v-stepper-content') {
        item.isVertical = this.vertical;
        this.content.push(item);
      }
    },
    unregister: function unregister(item) {
      if (item.$options.name === 'v-stepper-step') {
        this.steps = this.steps.filter(function (i) {
          return i !== item;
        });
      } else if (item.$options.name === 'v-stepper-content') {
        item.isVertical = this.vertical;
        this.content = this.content.filter(function (i) {
          return i !== item;
        });
      }
    },
    stepClick: function stepClick(step) {
      var _this = this;

      this.$nextTick(function () {
        return _this.internalValue = step;
      });
    },
    updateView: function updateView() {
      for (var index = this.steps.length; --index >= 0;) {
        this.steps[index].toggle(this.internalValue);
      }

      for (var _index = this.content.length; --_index >= 0;) {
        this.content[_index].toggle(this.internalValue, this.isReverse);
      }
    }
  },
  render: function render(h) {
    return h(this.tag, {
      staticClass: 'v-stepper',
      class: this.classes,
      style: this.styles
    }, this.$slots.default);
  }
});

exports.default = _default;
//# sourceMappingURL=VStepper.js.map