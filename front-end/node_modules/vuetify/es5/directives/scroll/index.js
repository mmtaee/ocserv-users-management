"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = exports.Scroll = void 0;

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

function inserted(el, binding, vnode) {
  var _ref = binding.modifiers || {},
      _ref$self = _ref.self,
      self = _ref$self === void 0 ? false : _ref$self;

  var value = binding.value;
  var options = _typeof(value) === 'object' && value.options || {
    passive: true
  };
  var handler = typeof value === 'function' || 'handleEvent' in value ? value : value.handler;
  var target = self ? el : binding.arg ? document.querySelector(binding.arg) : window;
  if (!target) return;
  target.addEventListener('scroll', handler, options);
  el._onScroll = Object(el._onScroll);
  el._onScroll[vnode.context._uid] = {
    handler: handler,
    options: options,
    // Don't reference self
    target: self ? undefined : target
  };
}

function unbind(el, binding, vnode) {
  var _a;

  if (!((_a = el._onScroll) === null || _a === void 0 ? void 0 : _a[vnode.context._uid])) return;
  var _el$_onScroll$vnode$c = el._onScroll[vnode.context._uid],
      handler = _el$_onScroll$vnode$c.handler,
      options = _el$_onScroll$vnode$c.options,
      _el$_onScroll$vnode$c2 = _el$_onScroll$vnode$c.target,
      target = _el$_onScroll$vnode$c2 === void 0 ? el : _el$_onScroll$vnode$c2;
  target.removeEventListener('scroll', handler, options);
  delete el._onScroll[vnode.context._uid];
}

var Scroll = {
  inserted: inserted,
  unbind: unbind
};
exports.Scroll = Scroll;
var _default = Scroll;
exports.default = _default;
//# sourceMappingURL=index.js.map