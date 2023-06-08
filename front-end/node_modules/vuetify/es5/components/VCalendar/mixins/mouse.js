"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

var _vue = _interopRequireDefault(require("vue"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

var _default = _vue.default.extend({
  name: 'mouse',
  methods: {
    getDefaultMouseEventHandlers: function getDefaultMouseEventHandlers(suffix, getEvent) {
      var listeners = Object.keys(this.$listeners).filter(function (key) {
        return key.endsWith(suffix);
      }).reduce(function (acc, key) {
        acc[key] = {
          event: key.slice(0, -suffix.length)
        };
        return acc;
      }, {});
      return this.getMouseEventHandlers(_objectSpread(_objectSpread({}, listeners), {}, _defineProperty({}, 'contextmenu' + suffix, {
        event: 'contextmenu',
        prevent: true,
        result: false
      })), getEvent);
    },
    getMouseEventHandlers: function getMouseEventHandlers(events, getEvent) {
      var _this = this;

      var on = {};

      var _loop = function _loop(event) {
        var eventOptions = events[event];
        if (!_this.$listeners[event]) return "continue"; // TODO somehow pull in modifiers

        var prefix = eventOptions.passive ? '&' : (eventOptions.once ? '~' : '') + (eventOptions.capture ? '!' : '');
        var key = prefix + eventOptions.event;

        var handler = function handler(e) {
          var _a, _b;

          var mouseEvent = e;

          if (eventOptions.button === undefined || mouseEvent.buttons > 0 && mouseEvent.button === eventOptions.button) {
            if (eventOptions.prevent) {
              e.preventDefault();
            }

            if (eventOptions.stop) {
              e.stopPropagation();
            } // Due to TouchEvent target always returns the element that is first placed
            // Even if touch point has since moved outside the interactive area of that element
            // Ref: https://developer.mozilla.org/en-US/docs/Web/API/Touch/target
            // This block of code aims to make sure touchEvent is always dispatched from the element that is being pointed at


            if (e && 'touches' in e) {
              var classSeparator = ' ';
              var eventTargetClasses = (_a = e.currentTarget) === null || _a === void 0 ? void 0 : _a.className.split(classSeparator);
              var currentTargets = document.elementsFromPoint(e.changedTouches[0].clientX, e.changedTouches[0].clientY); // Get "the same kind" current hovering target by checking
              // If element has the same class of initial touch start element (which has touch event listener registered)

              var currentTarget = currentTargets.find(function (t) {
                return t.className.split(classSeparator).some(function (c) {
                  return eventTargetClasses.includes(c);
                });
              });

              if (currentTarget && !((_b = e.target) === null || _b === void 0 ? void 0 : _b.isSameNode(currentTarget))) {
                currentTarget.dispatchEvent(new TouchEvent(e.type, {
                  changedTouches: e.changedTouches,
                  targetTouches: e.targetTouches,
                  touches: e.touches
                }));
                return;
              }
            }

            _this.$emit(event, getEvent(e), e);
          }

          return eventOptions.result;
        };

        if (key in on) {
          /* istanbul ignore next */
          if (Array.isArray(on[key])) {
            on[key].push(handler);
          } else {
            on[key] = [on[key], handler];
          }
        } else {
          on[key] = handler;
        }
      };

      for (var event in events) {
        var _ret = _loop(event);

        if (_ret === "continue") continue;
      }

      return on;
    }
  }
});

exports.default = _default;
//# sourceMappingURL=mouse.js.map