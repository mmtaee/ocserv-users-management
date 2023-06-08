"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.calculateUpdatedOffset = calculateUpdatedOffset;
exports.calculateCenteredOffset = calculateCenteredOffset;
exports.default = exports.BaseSlideGroup = void 0;

require("../../../src/components/VSlideGroup/VSlideGroup.sass");

var _VIcon = _interopRequireDefault(require("../VIcon"));

var _transitions = require("../transitions");

var _VItemGroup = require("../VItemGroup/VItemGroup");

var _mobile = _interopRequireDefault(require("../../mixins/mobile"));

var _resize = _interopRequireDefault(require("../../directives/resize"));

var _touch = _interopRequireDefault(require("../../directives/touch"));

var _mixins = _interopRequireDefault(require("../../util/mixins"));

var _helpers = require("../../util/helpers");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _createForOfIteratorHelper(o, allowArrayLike) { var it; if (typeof Symbol === "undefined" || o[Symbol.iterator] == null) { if (Array.isArray(o) || (it = _unsupportedIterableToArray(o)) || allowArrayLike && o && typeof o.length === "number") { if (it) o = it; var i = 0; var F = function F() {}; return { s: F, n: function n() { if (i >= o.length) return { done: true }; return { done: false, value: o[i++] }; }, e: function e(_e) { throw _e; }, f: F }; } throw new TypeError("Invalid attempt to iterate non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); } var normalCompletion = true, didErr = false, err; return { s: function s() { it = o[Symbol.iterator](); }, n: function n() { var step = it.next(); normalCompletion = step.done; return step; }, e: function e(_e2) { didErr = true; err = _e2; }, f: function f() { try { if (!normalCompletion && it.return != null) it.return(); } finally { if (didErr) throw err; } } }; }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

function bias(val) {
  var c = 0.501;
  var x = Math.abs(val);
  return Math.sign(val) * (x / ((1 / c - 2) * (1 - x) + 1));
}

function calculateUpdatedOffset(selectedElement, widths, rtl, currentScrollOffset) {
  var clientWidth = selectedElement.clientWidth;
  var offsetLeft = rtl ? widths.content - selectedElement.offsetLeft - clientWidth : selectedElement.offsetLeft;

  if (rtl) {
    currentScrollOffset = -currentScrollOffset;
  }

  var totalWidth = widths.wrapper + currentScrollOffset;
  var itemOffset = clientWidth + offsetLeft;
  var additionalOffset = clientWidth * 0.4;

  if (offsetLeft <= currentScrollOffset) {
    currentScrollOffset = Math.max(offsetLeft - additionalOffset, 0);
  } else if (totalWidth <= itemOffset) {
    currentScrollOffset = Math.min(currentScrollOffset - (totalWidth - itemOffset - additionalOffset), widths.content - widths.wrapper);
  }

  return rtl ? -currentScrollOffset : currentScrollOffset;
}

function calculateCenteredOffset(selectedElement, widths, rtl) {
  var offsetLeft = selectedElement.offsetLeft,
      clientWidth = selectedElement.clientWidth;

  if (rtl) {
    var offsetCentered = widths.content - offsetLeft - clientWidth / 2 - widths.wrapper / 2;
    return -Math.min(widths.content - widths.wrapper, Math.max(0, offsetCentered));
  } else {
    var _offsetCentered = offsetLeft + clientWidth / 2 - widths.wrapper / 2;

    return Math.min(widths.content - widths.wrapper, Math.max(0, _offsetCentered));
  }
}

var BaseSlideGroup = (0, _mixins.default)(_VItemGroup.BaseItemGroup, _mobile.default).extend({
  name: 'base-slide-group',
  directives: {
    Resize: _resize.default,
    Touch: _touch.default
  },
  props: {
    activeClass: {
      type: String,
      default: 'v-slide-item--active'
    },
    centerActive: Boolean,
    nextIcon: {
      type: String,
      default: '$next'
    },
    prevIcon: {
      type: String,
      default: '$prev'
    },
    showArrows: {
      type: [Boolean, String],
      validator: function validator(v) {
        return typeof v === 'boolean' || ['always', 'desktop', 'mobile'].includes(v);
      }
    }
  },
  data: function data() {
    return {
      isOverflowing: false,
      resizeTimeout: 0,
      startX: 0,
      isSwipingHorizontal: false,
      isSwiping: false,
      scrollOffset: 0,
      widths: {
        content: 0,
        wrapper: 0
      }
    };
  },
  computed: {
    canTouch: function canTouch() {
      return typeof window !== 'undefined';
    },
    __cachedNext: function __cachedNext() {
      return this.genTransition('next');
    },
    __cachedPrev: function __cachedPrev() {
      return this.genTransition('prev');
    },
    classes: function classes() {
      return _objectSpread(_objectSpread({}, _VItemGroup.BaseItemGroup.options.computed.classes.call(this)), {}, {
        'v-slide-group': true,
        'v-slide-group--has-affixes': this.hasAffixes,
        'v-slide-group--is-overflowing': this.isOverflowing
      });
    },
    hasAffixes: function hasAffixes() {
      switch (this.showArrows) {
        // Always show arrows on desktop & mobile
        case 'always':
          return true;
        // Always show arrows on desktop

        case 'desktop':
          return !this.isMobile;
        // Show arrows on mobile when overflowing.
        // This matches the default 2.2 behavior

        case true:
          return this.isOverflowing || Math.abs(this.scrollOffset) > 0;
        // Always show on mobile

        case 'mobile':
          return this.isMobile || this.isOverflowing || Math.abs(this.scrollOffset) > 0;
        // https://material.io/components/tabs#scrollable-tabs
        // Always show arrows when
        // overflowed on desktop

        default:
          return !this.isMobile && (this.isOverflowing || Math.abs(this.scrollOffset) > 0);
      }
    },
    hasNext: function hasNext() {
      if (!this.hasAffixes) return false;
      var _this$widths = this.widths,
          content = _this$widths.content,
          wrapper = _this$widths.wrapper; // Check one scroll ahead to know the width of right-most item

      return content > Math.abs(this.scrollOffset) + wrapper;
    },
    hasPrev: function hasPrev() {
      return this.hasAffixes && this.scrollOffset !== 0;
    }
  },
  watch: {
    internalValue: 'setWidths',
    // When overflow changes, the arrows alter
    // the widths of the content and wrapper
    // and need to be recalculated
    isOverflowing: 'setWidths',
    scrollOffset: function scrollOffset(val) {
      if (this.$vuetify.rtl) val = -val;
      var scroll = val <= 0 ? bias(-val) : val > this.widths.content - this.widths.wrapper ? -(this.widths.content - this.widths.wrapper) + bias(this.widths.content - this.widths.wrapper - val) : -val;
      if (this.$vuetify.rtl) scroll = -scroll;
      this.$refs.content.style.transform = "translateX(".concat(scroll, "px)");
    }
  },
  mounted: function mounted() {
    var _this = this;

    if (typeof ResizeObserver !== 'undefined') {
      var obs = new ResizeObserver(function () {
        _this.onResize();
      });
      obs.observe(this.$el);
      obs.observe(this.$refs.content);
      this.$on('hook:destroyed', function () {
        obs.disconnect();
      });
    } else {
      var itemsLength = 0;
      this.$on('hook:beforeUpdate', function () {
        var _a;

        itemsLength = (((_a = _this.$refs.content) === null || _a === void 0 ? void 0 : _a.children) || []).length;
      });
      this.$on('hook:updated', function () {
        var _a;

        if (itemsLength === (((_a = _this.$refs.content) === null || _a === void 0 ? void 0 : _a.children) || []).length) return;

        _this.setWidths();
      });
    }
  },
  methods: {
    onScroll: function onScroll() {
      this.$refs.wrapper.scrollLeft = 0;
    },
    onFocusin: function onFocusin(e) {
      if (!this.isOverflowing) return; // Focused element is likely to be the root of an item, so a
      // breadth-first search will probably find it in the first iteration

      var _iterator = _createForOfIteratorHelper((0, _helpers.composedPath)(e)),
          _step;

      try {
        for (_iterator.s(); !(_step = _iterator.n()).done;) {
          var el = _step.value;

          var _iterator2 = _createForOfIteratorHelper(this.items),
              _step2;

          try {
            for (_iterator2.s(); !(_step2 = _iterator2.n()).done;) {
              var vm = _step2.value;

              if (vm.$el === el) {
                this.scrollOffset = calculateUpdatedOffset(vm.$el, this.widths, this.$vuetify.rtl, this.scrollOffset);
                return;
              }
            }
          } catch (err) {
            _iterator2.e(err);
          } finally {
            _iterator2.f();
          }
        }
      } catch (err) {
        _iterator.e(err);
      } finally {
        _iterator.f();
      }
    },
    // Always generate next for scrollable hint
    genNext: function genNext() {
      var _this2 = this;

      var slot = this.$scopedSlots.next ? this.$scopedSlots.next({}) : this.$slots.next || this.__cachedNext;
      return this.$createElement('div', {
        staticClass: 'v-slide-group__next',
        class: {
          'v-slide-group__next--disabled': !this.hasNext
        },
        on: {
          click: function click() {
            return _this2.onAffixClick('next');
          }
        },
        key: 'next'
      }, [slot]);
    },
    genContent: function genContent() {
      return this.$createElement('div', {
        staticClass: 'v-slide-group__content',
        ref: 'content',
        on: {
          focusin: this.onFocusin
        }
      }, this.$slots.default);
    },
    genData: function genData() {
      return {
        class: this.classes,
        directives: [{
          name: 'resize',
          value: this.onResize
        }]
      };
    },
    genIcon: function genIcon(location) {
      var icon = location;

      if (this.$vuetify.rtl && location === 'prev') {
        icon = 'next';
      } else if (this.$vuetify.rtl && location === 'next') {
        icon = 'prev';
      }

      var upperLocation = "".concat(location[0].toUpperCase()).concat(location.slice(1));
      var hasAffix = this["has".concat(upperLocation)];
      if (!this.showArrows && !hasAffix) return null;
      return this.$createElement(_VIcon.default, {
        props: {
          disabled: !hasAffix
        }
      }, this["".concat(icon, "Icon")]);
    },
    // Always generate prev for scrollable hint
    genPrev: function genPrev() {
      var _this3 = this;

      var slot = this.$scopedSlots.prev ? this.$scopedSlots.prev({}) : this.$slots.prev || this.__cachedPrev;
      return this.$createElement('div', {
        staticClass: 'v-slide-group__prev',
        class: {
          'v-slide-group__prev--disabled': !this.hasPrev
        },
        on: {
          click: function click() {
            return _this3.onAffixClick('prev');
          }
        },
        key: 'prev'
      }, [slot]);
    },
    genTransition: function genTransition(location) {
      return this.$createElement(_transitions.VFadeTransition, [this.genIcon(location)]);
    },
    genWrapper: function genWrapper() {
      var _this4 = this;

      return this.$createElement('div', {
        staticClass: 'v-slide-group__wrapper',
        directives: [{
          name: 'touch',
          value: {
            start: function start(e) {
              return _this4.overflowCheck(e, _this4.onTouchStart);
            },
            move: function move(e) {
              return _this4.overflowCheck(e, _this4.onTouchMove);
            },
            end: function end(e) {
              return _this4.overflowCheck(e, _this4.onTouchEnd);
            }
          }
        }],
        ref: 'wrapper',
        on: {
          scroll: this.onScroll
        }
      }, [this.genContent()]);
    },
    calculateNewOffset: function calculateNewOffset(direction, widths, rtl, currentScrollOffset) {
      var sign = rtl ? -1 : 1;
      var newAbosluteOffset = sign * currentScrollOffset + (direction === 'prev' ? -1 : 1) * widths.wrapper;
      return sign * Math.max(Math.min(newAbosluteOffset, widths.content - widths.wrapper), 0);
    },
    onAffixClick: function onAffixClick(location) {
      this.$emit("click:".concat(location));
      this.scrollTo(location);
    },
    onResize: function onResize() {
      /* istanbul ignore next */
      if (this._isDestroyed) return;
      this.setWidths();
    },
    onTouchStart: function onTouchStart(e) {
      var content = this.$refs.content;
      this.startX = this.scrollOffset + e.touchstartX;
      content.style.setProperty('transition', 'none');
      content.style.setProperty('willChange', 'transform');
    },
    onTouchMove: function onTouchMove(e) {
      if (!this.canTouch) return;

      if (!this.isSwiping) {
        // only calculate disableSwipeHorizontal during the first onTouchMove invoke
        // in order to ensure disableSwipeHorizontal value is consistent between onTouchStart and onTouchEnd
        var diffX = e.touchmoveX - e.touchstartX;
        var diffY = e.touchmoveY - e.touchstartY;
        this.isSwipingHorizontal = Math.abs(diffX) > Math.abs(diffY);
        this.isSwiping = true;
      }

      if (this.isSwipingHorizontal) {
        // sliding horizontally
        this.scrollOffset = this.startX - e.touchmoveX; // temporarily disable window vertical scrolling

        document.documentElement.style.overflowY = 'hidden';
      }
    },
    onTouchEnd: function onTouchEnd() {
      if (!this.canTouch) return;
      var _this$$refs = this.$refs,
          content = _this$$refs.content,
          wrapper = _this$$refs.wrapper;
      var maxScrollOffset = content.clientWidth - wrapper.clientWidth;
      content.style.setProperty('transition', null);
      content.style.setProperty('willChange', null);

      if (this.$vuetify.rtl) {
        /* istanbul ignore else */
        if (this.scrollOffset > 0 || !this.isOverflowing) {
          this.scrollOffset = 0;
        } else if (this.scrollOffset <= -maxScrollOffset) {
          this.scrollOffset = -maxScrollOffset;
        }
      } else {
        /* istanbul ignore else */
        if (this.scrollOffset < 0 || !this.isOverflowing) {
          this.scrollOffset = 0;
        } else if (this.scrollOffset >= maxScrollOffset) {
          this.scrollOffset = maxScrollOffset;
        }
      }

      this.isSwiping = false; // rollback whole page scrolling to default

      document.documentElement.style.removeProperty('overflow-y');
    },
    overflowCheck: function overflowCheck(e, fn) {
      e.stopPropagation();
      this.isOverflowing && fn(e);
    },
    scrollIntoView
    /* istanbul ignore next */
    : function scrollIntoView() {
      if (!this.selectedItem && this.items.length) {
        var lastItemPosition = this.items[this.items.length - 1].$el.getBoundingClientRect();
        var wrapperPosition = this.$refs.wrapper.getBoundingClientRect();

        if (this.$vuetify.rtl && wrapperPosition.right < lastItemPosition.right || !this.$vuetify.rtl && wrapperPosition.left > lastItemPosition.left) {
          this.scrollTo('prev');
        }
      }

      if (!this.selectedItem) {
        return;
      }

      if (this.selectedIndex === 0 || !this.centerActive && !this.isOverflowing) {
        this.scrollOffset = 0;
      } else if (this.centerActive) {
        this.scrollOffset = calculateCenteredOffset(this.selectedItem.$el, this.widths, this.$vuetify.rtl);
      } else if (this.isOverflowing) {
        this.scrollOffset = calculateUpdatedOffset(this.selectedItem.$el, this.widths, this.$vuetify.rtl, this.scrollOffset);
      }
    },
    scrollTo
    /* istanbul ignore next */
    : function scrollTo(location) {
      this.scrollOffset = this.calculateNewOffset(location, {
        // Force reflow
        content: this.$refs.content ? this.$refs.content.clientWidth : 0,
        wrapper: this.$refs.wrapper ? this.$refs.wrapper.clientWidth : 0
      }, this.$vuetify.rtl, this.scrollOffset);
    },
    setWidths: function setWidths() {
      var _this5 = this;

      window.requestAnimationFrame(function () {
        if (_this5._isDestroyed) return;
        var _this5$$refs = _this5.$refs,
            content = _this5$$refs.content,
            wrapper = _this5$$refs.wrapper;
        _this5.widths = {
          content: content ? content.clientWidth : 0,
          wrapper: wrapper ? wrapper.clientWidth : 0
        }; // https://github.com/vuetifyjs/vuetify/issues/13212
        // We add +1 to the wrappers width to prevent an issue where the `clientWidth`
        // gets calculated wrongly by the browser if using a different zoom-level.

        _this5.isOverflowing = _this5.widths.wrapper + 1 < _this5.widths.content;

        _this5.scrollIntoView();
      });
    }
  },
  render: function render(h) {
    return h('div', this.genData(), [this.genPrev(), this.genWrapper(), this.genNext()]);
  }
});
exports.BaseSlideGroup = BaseSlideGroup;

var _default = BaseSlideGroup.extend({
  name: 'v-slide-group',
  provide: function provide() {
    return {
      slideGroup: this
    };
  }
});

exports.default = _default;
//# sourceMappingURL=VSlideGroup.js.map