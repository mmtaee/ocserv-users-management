"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

var _mixins = _interopRequireDefault(require("../../util/mixins"));

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _toConsumableArray(arr) { return _arrayWithoutHoles(arr) || _iterableToArray(arr) || _unsupportedIterableToArray(arr) || _nonIterableSpread(); }

function _nonIterableSpread() { throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _iterableToArray(iter) { if (typeof Symbol !== "undefined" && Symbol.iterator in Object(iter)) return Array.from(iter); }

function _arrayWithoutHoles(arr) { if (Array.isArray(arr)) return _arrayLikeToArray(arr); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

function searchChildren(children) {
  var results = [];

  for (var index = 0; index < children.length; index++) {
    var child = children[index];

    if (child.isActive && child.isDependent) {
      results.push(child);
    } else {
      results.push.apply(results, _toConsumableArray(searchChildren(child.$children)));
    }
  }

  return results;
}
/* @vue/component */


var _default = (0, _mixins.default)().extend({
  name: 'dependent',
  data: function data() {
    return {
      closeDependents: true,
      isActive: false,
      isDependent: true
    };
  },
  watch: {
    isActive: function isActive(val) {
      if (val) return;
      var openDependents = this.getOpenDependents();

      for (var index = 0; index < openDependents.length; index++) {
        openDependents[index].isActive = false;
      }
    }
  },
  methods: {
    getOpenDependents: function getOpenDependents() {
      if (this.closeDependents) return searchChildren(this.$children);
      return [];
    },
    getOpenDependentElements: function getOpenDependentElements() {
      var result = [];
      var openDependents = this.getOpenDependents();

      for (var index = 0; index < openDependents.length; index++) {
        result.push.apply(result, _toConsumableArray(openDependents[index].getClickableDependentElements()));
      }

      return result;
    },
    getClickableDependentElements: function getClickableDependentElements() {
      var result = [this.$el];
      if (this.$refs.content) result.push(this.$refs.content);
      if (this.overlay) result.push(this.overlay.$el);
      result.push.apply(result, _toConsumableArray(this.getOpenDependentElements()));
      return result;
    }
  }
});

exports.default = _default;
//# sourceMappingURL=index.js.map