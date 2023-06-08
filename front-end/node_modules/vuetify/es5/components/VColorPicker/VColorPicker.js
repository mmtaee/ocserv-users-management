"use strict";

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

require("../../../src/components/VColorPicker/VColorPicker.sass");

var _VSheet = _interopRequireDefault(require("../VSheet/VSheet"));

var _VColorPickerPreview = _interopRequireDefault(require("./VColorPickerPreview"));

var _VColorPickerCanvas = _interopRequireDefault(require("./VColorPickerCanvas"));

var _VColorPickerEdit = _interopRequireWildcard(require("./VColorPickerEdit"));

var _VColorPickerSwatches = _interopRequireDefault(require("./VColorPickerSwatches"));

var _util = require("./util");

var _mixins = _interopRequireDefault(require("../../util/mixins"));

var _helpers = require("../../util/helpers");

var _elevatable = _interopRequireDefault(require("../../mixins/elevatable"));

var _themeable = _interopRequireDefault(require("../../mixins/themeable"));

function _getRequireWildcardCache() { if (typeof WeakMap !== "function") return null; var cache = new WeakMap(); _getRequireWildcardCache = function _getRequireWildcardCache() { return cache; }; return cache; }

function _interopRequireWildcard(obj) { if (obj && obj.__esModule) { return obj; } if (obj === null || _typeof(obj) !== "object" && typeof obj !== "function") { return { default: obj }; } var cache = _getRequireWildcardCache(); if (cache && cache.has(obj)) { return cache.get(obj); } var newObj = {}; var hasPropertyDescriptor = Object.defineProperty && Object.getOwnPropertyDescriptor; for (var key in obj) { if (Object.prototype.hasOwnProperty.call(obj, key)) { var desc = hasPropertyDescriptor ? Object.getOwnPropertyDescriptor(obj, key) : null; if (desc && (desc.get || desc.set)) { Object.defineProperty(newObj, key, desc); } else { newObj[key] = obj[key]; } } } newObj.default = obj; if (cache) { cache.set(obj, newObj); } return newObj; }

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

var _default = (0, _mixins.default)(_elevatable.default, _themeable.default).extend({
  name: 'v-color-picker',
  props: {
    canvasHeight: {
      type: [String, Number],
      default: 150
    },
    disabled: Boolean,
    dotSize: {
      type: [Number, String],
      default: 10
    },
    flat: Boolean,
    hideCanvas: Boolean,
    hideSliders: Boolean,
    hideInputs: Boolean,
    hideModeSwitch: Boolean,
    mode: {
      type: String,
      default: 'rgba',
      validator: function validator(v) {
        return Object.keys(_VColorPickerEdit.modes).includes(v);
      }
    },
    showSwatches: Boolean,
    swatches: Array,
    swatchesMaxHeight: {
      type: [Number, String],
      default: 150
    },
    value: {
      type: [Object, String]
    },
    width: {
      type: [Number, String],
      default: 300
    }
  },
  data: function data() {
    return {
      internalValue: (0, _util.fromRGBA)({
        r: 255,
        g: 0,
        b: 0,
        a: 1
      })
    };
  },
  computed: {
    hideAlpha: function hideAlpha() {
      if (!this.value) return false;
      return !(0, _util.hasAlpha)(this.value);
    }
  },
  watch: {
    value: {
      handler: function handler(color) {
        this.updateColor((0, _util.parseColor)(color, this.internalValue));
      },
      immediate: true
    }
  },
  methods: {
    updateColor: function updateColor(color) {
      this.internalValue = color;
      var value = (0, _util.extractColor)(this.internalValue, this.value);

      if (!(0, _helpers.deepEqual)(value, this.value)) {
        this.$emit('input', value);
        this.$emit('update:color', this.internalValue);
      }
    },
    genCanvas: function genCanvas() {
      return this.$createElement(_VColorPickerCanvas.default, {
        props: {
          color: this.internalValue,
          disabled: this.disabled,
          dotSize: this.dotSize,
          width: this.width,
          height: this.canvasHeight
        },
        on: {
          'update:color': this.updateColor
        }
      });
    },
    genControls: function genControls() {
      return this.$createElement('div', {
        staticClass: 'v-color-picker__controls'
      }, [!this.hideSliders && this.genPreview(), !this.hideInputs && this.genEdit()]);
    },
    genEdit: function genEdit() {
      var _this = this;

      return this.$createElement(_VColorPickerEdit.default, {
        props: {
          color: this.internalValue,
          disabled: this.disabled,
          hideAlpha: this.hideAlpha,
          hideModeSwitch: this.hideModeSwitch,
          mode: this.mode
        },
        on: {
          'update:color': this.updateColor,
          'update:mode': function updateMode(v) {
            return _this.$emit('update:mode', v);
          }
        }
      });
    },
    genPreview: function genPreview() {
      return this.$createElement(_VColorPickerPreview.default, {
        props: {
          color: this.internalValue,
          disabled: this.disabled,
          hideAlpha: this.hideAlpha
        },
        on: {
          'update:color': this.updateColor
        }
      });
    },
    genSwatches: function genSwatches() {
      return this.$createElement(_VColorPickerSwatches.default, {
        props: {
          dark: this.dark,
          light: this.light,
          disabled: this.disabled,
          swatches: this.swatches,
          color: this.internalValue,
          maxHeight: this.swatchesMaxHeight
        },
        on: {
          'update:color': this.updateColor
        }
      });
    }
  },
  render: function render(h) {
    return h(_VSheet.default, {
      staticClass: 'v-color-picker',
      class: _objectSpread(_objectSpread({
        'v-color-picker--flat': this.flat
      }, this.themeClasses), this.elevationClasses),
      props: {
        maxWidth: this.width
      }
    }, [!this.hideCanvas && this.genCanvas(), (!this.hideSliders || !this.hideInputs) && this.genControls(), this.showSwatches && this.genSwatches()]);
  }
});

exports.default = _default;
//# sourceMappingURL=VColorPicker.js.map