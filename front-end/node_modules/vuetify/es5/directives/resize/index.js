"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = exports.Resize = void 0;

function inserted(el, binding, vnode) {
  var callback = binding.value;
  var options = binding.options || {
    passive: true
  };
  window.addEventListener('resize', callback, options);
  el._onResize = Object(el._onResize);
  el._onResize[vnode.context._uid] = {
    callback: callback,
    options: options
  };

  if (!binding.modifiers || !binding.modifiers.quiet) {
    callback();
  }
}

function unbind(el, binding, vnode) {
  var _a;

  if (!((_a = el._onResize) === null || _a === void 0 ? void 0 : _a[vnode.context._uid])) return;
  var _el$_onResize$vnode$c = el._onResize[vnode.context._uid],
      callback = _el$_onResize$vnode$c.callback,
      options = _el$_onResize$vnode$c.options;
  window.removeEventListener('resize', callback, options);
  delete el._onResize[vnode.context._uid];
}

var Resize = {
  inserted: inserted,
  unbind: unbind
};
exports.Resize = Resize;
var _default = Resize;
exports.default = _default;
//# sourceMappingURL=index.js.map