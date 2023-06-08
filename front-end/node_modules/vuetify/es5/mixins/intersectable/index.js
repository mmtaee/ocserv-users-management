"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = intersectable;

var _intersect = _interopRequireDefault(require("../../directives/intersect"));

var _console = require("../../util/console");

var _vue = _interopRequireDefault(require("vue"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

// Directives
// Utilities
// Types
function intersectable(options) {
  return _vue.default.extend({
    name: 'intersectable',
    data: function data() {
      return {
        isIntersecting: false
      };
    },
    mounted: function mounted() {
      _intersect.default.inserted(this.$el, {
        name: 'intersect',
        value: this.onObserve
      }, this.$vnode);
    },
    destroyed: function destroyed() {
      _intersect.default.unbind(this.$el, {
        name: 'intersect',
        value: this.onObserve
      }, this.$vnode);
    },
    methods: {
      onObserve: function onObserve(entries, observer, isIntersecting) {
        this.isIntersecting = isIntersecting;
        if (!isIntersecting) return;

        for (var i = 0, length = options.onVisible.length; i < length; i++) {
          var callback = this[options.onVisible[i]];

          if (typeof callback === 'function') {
            callback();
            continue;
          }

          (0, _console.consoleWarn)(options.onVisible[i] + ' method is not available on the instance but referenced in intersectable mixin options');
        }
      }
    }
  });
}
//# sourceMappingURL=index.js.map