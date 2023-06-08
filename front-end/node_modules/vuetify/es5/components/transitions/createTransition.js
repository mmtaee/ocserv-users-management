"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.createSimpleTransition = createSimpleTransition;
exports.createJavascriptTransition = createJavascriptTransition;

var _mergeData = _interopRequireDefault(require("../../util/mergeData"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function mergeTransitions() {
  var _Array;

  var dest = arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : [];

  for (var _len = arguments.length, transitions = new Array(_len > 1 ? _len - 1 : 0), _key = 1; _key < _len; _key++) {
    transitions[_key - 1] = arguments[_key];
  }

  /* eslint-disable-next-line no-array-constructor */
  return (_Array = Array()).concat.apply(_Array, [dest].concat(transitions));
}

function createSimpleTransition(name) {
  var origin = arguments.length > 1 && arguments[1] !== undefined ? arguments[1] : 'top center 0';
  var mode = arguments.length > 2 ? arguments[2] : undefined;
  return {
    name: name,
    functional: true,
    props: {
      group: {
        type: Boolean,
        default: false
      },
      hideOnLeave: {
        type: Boolean,
        default: false
      },
      leaveAbsolute: {
        type: Boolean,
        default: false
      },
      mode: {
        type: String,
        default: mode
      },
      origin: {
        type: String,
        default: origin
      }
    },
    render: function render(h, context) {
      var tag = "transition".concat(context.props.group ? '-group' : '');
      var data = {
        props: {
          name: name,
          mode: context.props.mode
        },
        on: {
          beforeEnter: function beforeEnter(el) {
            el.style.transformOrigin = context.props.origin;
            el.style.webkitTransformOrigin = context.props.origin;
          }
        }
      };

      if (context.props.leaveAbsolute) {
        data.on.leave = mergeTransitions(data.on.leave, function (el) {
          var offsetTop = el.offsetTop,
              offsetLeft = el.offsetLeft,
              offsetWidth = el.offsetWidth,
              offsetHeight = el.offsetHeight;
          el._transitionInitialStyles = {
            position: el.style.position,
            top: el.style.top,
            left: el.style.left,
            width: el.style.width,
            height: el.style.height
          };
          el.style.position = 'absolute';
          el.style.top = offsetTop + 'px';
          el.style.left = offsetLeft + 'px';
          el.style.width = offsetWidth + 'px';
          el.style.height = offsetHeight + 'px';
        });
        data.on.afterLeave = mergeTransitions(data.on.afterLeave, function (el) {
          if (el && el._transitionInitialStyles) {
            var _el$_transitionInitia = el._transitionInitialStyles,
                position = _el$_transitionInitia.position,
                top = _el$_transitionInitia.top,
                left = _el$_transitionInitia.left,
                width = _el$_transitionInitia.width,
                height = _el$_transitionInitia.height;
            delete el._transitionInitialStyles;
            el.style.position = position || '';
            el.style.top = top || '';
            el.style.left = left || '';
            el.style.width = width || '';
            el.style.height = height || '';
          }
        });
      }

      if (context.props.hideOnLeave) {
        data.on.leave = mergeTransitions(data.on.leave, function (el) {
          el.style.setProperty('display', 'none', 'important');
        });
      }

      return h(tag, (0, _mergeData.default)(context.data, data), context.children);
    }
  };
}

function createJavascriptTransition(name, functions) {
  var mode = arguments.length > 2 && arguments[2] !== undefined ? arguments[2] : 'in-out';
  return {
    name: name,
    functional: true,
    props: {
      mode: {
        type: String,
        default: mode
      }
    },
    render: function render(h, context) {
      return h('transition', (0, _mergeData.default)(context.data, {
        props: {
          name: name
        },
        on: functions
      }), context.children);
    }
  };
}
//# sourceMappingURL=createTransition.js.map