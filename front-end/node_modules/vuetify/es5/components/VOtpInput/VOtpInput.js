"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

require("../../../src/components/VTextField/VTextField.sass");

require("../../../src/components/VOtpInput/VOtpInput.sass");

var _VInput = _interopRequireDefault(require("../VInput"));

var _VTextField2 = _interopRequireDefault(require("../VTextField/VTextField"));

var _ripple = _interopRequireDefault(require("../../directives/ripple"));

var _helpers = require("../../util/helpers");

var _console = require("../../util/console");

var _mixins = _interopRequireDefault(require("../../util/mixins"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _toConsumableArray(arr) { return _arrayWithoutHoles(arr) || _iterableToArray(arr) || _unsupportedIterableToArray(arr) || _nonIterableSpread(); }

function _nonIterableSpread() { throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _iterableToArray(iter) { if (typeof Symbol !== "undefined" && Symbol.iterator in Object(iter)) return Array.from(iter); }

function _arrayWithoutHoles(arr) { if (Array.isArray(arr)) return _arrayLikeToArray(arr); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

var baseMixins = (0, _mixins.default)(_VInput.default);
/* @vue/component */

var _default = baseMixins.extend().extend({
  name: 'v-otp-input',
  directives: {
    ripple: _ripple.default
  },
  inheritAttrs: false,
  props: {
    length: {
      type: [Number, String],
      default: 6
    },
    type: {
      type: String,
      default: 'text'
    },
    plain: Boolean
  },
  data: function data() {
    return {
      initialValue: null,
      isBooted: false,
      otp: []
    };
  },
  computed: {
    outlined: function outlined() {
      return !this.plain;
    },
    classes: function classes() {
      return _objectSpread(_objectSpread(_objectSpread({}, _VInput.default.options.computed.classes.call(this)), _VTextField2.default.options.computed.classes.call(this)), {}, {
        'v-otp-input--plain': this.plain
      });
    }
  },
  watch: {
    isFocused: 'updateValue',
    value: function value(val) {
      this.lazyValue = val;
      this.otp = (val === null || val === void 0 ? void 0 : val.split('')) || [];
    }
  },
  created: function created() {
    var _a;
    /* istanbul ignore next */


    if (this.$attrs.hasOwnProperty('browser-autocomplete')) {
      (0, _console.breaking)('browser-autocomplete', 'autocomplete', this);
    }

    this.otp = ((_a = this.internalValue) === null || _a === void 0 ? void 0 : _a.split('')) || [];
  },
  mounted: function mounted() {
    var _this = this;

    requestAnimationFrame(function () {
      return _this.isBooted = true;
    });
  },
  methods: {
    /** @public */
    focus: function focus(e, otpIdx) {
      this.onFocus(e, otpIdx || 0);
    },
    genInputSlot: function genInputSlot(otpIdx) {
      var _this2 = this;

      return this.$createElement('div', this.setBackgroundColor(this.backgroundColor, {
        staticClass: 'v-input__slot',
        style: {
          height: (0, _helpers.convertToUnit)(this.height)
        },
        on: {
          click: function click() {
            return _this2.onClick(otpIdx);
          },
          mousedown: function mousedown(e) {
            return _this2.onMouseDown(e, otpIdx);
          },
          mouseup: function mouseup(e) {
            return _this2.onMouseUp(e, otpIdx);
          }
        }
      }), [this.genDefaultSlot(otpIdx)]);
    },
    genControl: function genControl(otpIdx) {
      return this.$createElement('div', {
        staticClass: 'v-input__control'
      }, [this.genInputSlot(otpIdx)]);
    },
    genDefaultSlot: function genDefaultSlot(otpIdx) {
      return [this.genFieldset(), this.genTextFieldSlot(otpIdx)];
    },
    genContent: function genContent() {
      var _this3 = this;

      return Array.from({
        length: +this.length
      }, function (_, i) {
        return _this3.$createElement('div', _this3.setTextColor(_this3.validationState, {
          staticClass: 'v-input',
          class: _this3.classes
        }), [_this3.genControl(i)]);
      });
    },
    genFieldset: function genFieldset() {
      return this.$createElement('fieldset', {
        attrs: {
          'aria-hidden': true
        }
      }, [this.genLegend()]);
    },
    genLegend: function genLegend() {
      var span = this.$createElement('span', {
        domProps: {
          innerHTML: '&#8203;'
        }
      });
      return this.$createElement('legend', {
        style: {
          width: '0px'
        }
      }, [span]);
    },
    genInput: function genInput(otpIdx) {
      var _this4 = this;

      var listeners = Object.assign({}, this.listeners$);
      delete listeners.change; // Change should not be bound externally

      return this.$createElement('input', {
        style: {},
        domProps: {
          value: this.otp[otpIdx],
          min: this.type === 'number' ? 0 : null
        },
        attrs: _objectSpread(_objectSpread({}, this.attrs$), {}, {
          autocomplete: 'one-time-code',
          disabled: this.isDisabled,
          readonly: this.isReadonly,
          type: this.type,
          id: "".concat(this.computedId, "--").concat(otpIdx),
          class: "otp-field-box--".concat(otpIdx)
        }),
        on: Object.assign(listeners, {
          blur: this.onBlur,
          input: function input(e) {
            return _this4.onInput(e, otpIdx);
          },
          focus: function focus(e) {
            return _this4.onFocus(e, otpIdx);
          },
          keydown: this.onKeyDown,
          keyup: function keyup(e) {
            return _this4.onKeyUp(e, otpIdx);
          }
        }),
        ref: 'input',
        refInFor: true
      });
    },
    genTextFieldSlot: function genTextFieldSlot(otpIdx) {
      return this.$createElement('div', {
        staticClass: 'v-text-field__slot'
      }, [this.genInput(otpIdx)]);
    },
    onBlur: function onBlur(e) {
      var _this5 = this;

      this.isFocused = false;
      e && this.$nextTick(function () {
        return _this5.$emit('blur', e);
      });
    },
    onClick: function onClick(otpIdx) {
      if (this.isFocused || this.isDisabled || !this.$refs.input[otpIdx]) return;
      this.onFocus(undefined, otpIdx);
    },
    onFocus: function onFocus(e, otpIdx) {
      e === null || e === void 0 ? void 0 : e.preventDefault();
      e === null || e === void 0 ? void 0 : e.stopPropagation();
      var elements = this.$refs.input;
      var ref = this.$refs.input && elements[otpIdx || 0];
      if (!ref) return;

      if (document.activeElement !== ref) {
        ref.focus();
        return ref.select();
      }

      if (!this.isFocused) {
        this.isFocused = true;
        ref.select();
        e && this.$emit('focus', e);
      }
    },
    onInput: function onInput(e, index) {
      var maxCursor = +this.length - 1;
      var target = e.target;
      var value = target.value;
      var inputDataArray = (value === null || value === void 0 ? void 0 : value.split('')) || [];

      var newOtp = _toConsumableArray(this.otp);

      for (var i = 0; i < inputDataArray.length; i++) {
        var appIdx = index + i;
        if (appIdx > maxCursor) break;
        newOtp[appIdx] = inputDataArray[i].toString();
      }

      if (!inputDataArray.length) {
        newOtp.splice(index, 1);
      }

      this.otp = newOtp;
      this.internalValue = this.otp.join('');

      if (index + inputDataArray.length >= +this.length) {
        this.onCompleted();
        this.clearFocus(index);
      } else if (inputDataArray.length) {
        this.changeFocus(index + inputDataArray.length);
      }
    },
    clearFocus: function clearFocus(index) {
      var input = this.$refs.input[index];
      input.blur();
    },
    onKeyDown: function onKeyDown(e) {
      if (e.keyCode === _helpers.keyCodes.enter) {
        this.$emit('change', this.internalValue);
      }

      this.$emit('keydown', e);
    },
    onMouseDown: function onMouseDown(e, otpIdx) {
      // Prevent input from being blurred
      if (e.target !== this.$refs.input[otpIdx]) {
        e.preventDefault();
        e.stopPropagation();
      }

      _VInput.default.options.methods.onMouseDown.call(this, e);
    },
    onMouseUp: function onMouseUp(e, otpIdx) {
      if (this.hasMouseDown) this.focus(e, otpIdx);

      _VInput.default.options.methods.onMouseUp.call(this, e);
    },
    changeFocus: function changeFocus(index) {
      this.onFocus(undefined, index || 0);
    },
    updateValue: function updateValue(val) {
      // Sets validationState from validatable
      this.hasColor = val;

      if (val) {
        this.initialValue = this.lazyValue;
      } else if (this.initialValue !== this.lazyValue) {
        this.$emit('change', this.lazyValue);
      }
    },
    onKeyUp: function onKeyUp(event, index) {
      event.preventDefault();
      var eventKey = event.key;

      if (['Tab', 'Shift', 'Meta', 'Control', 'Alt'].includes(eventKey)) {
        return;
      }

      if (['Delete'].includes(eventKey)) {
        return;
      }

      if (eventKey === 'ArrowLeft' || eventKey === 'Backspace' && !this.otp[index]) {
        return index > 0 && this.changeFocus(index - 1);
      }

      if (eventKey === 'ArrowRight') {
        return index + 1 < +this.length && this.changeFocus(index + 1);
      }
    },
    onCompleted: function onCompleted() {
      var rsp = this.otp.join('');

      if (rsp.length === +this.length) {
        this.$emit('finish', rsp);
      }
    }
  },
  render: function render(h) {
    return h('div', {
      staticClass: 'v-otp-input',
      class: this.themeClasses
    }, this.genContent());
  }
});

exports.default = _default;
//# sourceMappingURL=VOtpInput.js.map