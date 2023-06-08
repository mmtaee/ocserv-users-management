"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = exports.ClickOutside = void 0;

var _dom = require("../../util/dom");

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

function defaultConditional() {
  return true;
}

function checkEvent(e, el, binding) {
  // The include element callbacks below can be expensive
  // so we should avoid calling them when we're not active.
  // Explicitly check for false to allow fallback compatibility
  // with non-toggleable components
  if (!e || checkIsActive(e, binding) === false) return false; // If we're clicking inside the shadowroot, then the app root doesn't get the same
  // level of introspection as to _what_ we're clicking. We want to check to see if
  // our target is the shadowroot parent container, and if it is, ignore.

  var root = (0, _dom.attachedRoot)(el);
  if (typeof ShadowRoot !== 'undefined' && root instanceof ShadowRoot && root.host === e.target) return false; // Check if additional elements were passed to be included in check
  // (click must be outside all included elements, if any)

  var elements = (_typeof(binding.value) === 'object' && binding.value.include || function () {
    return [];
  })(); // Add the root element for the component this directive was defined on


  elements.push(el); // Check if it's a click outside our elements, and then if our callback returns true.
  // Non-toggleable components should take action in their callback and return falsy.
  // Toggleable can return true if it wants to deactivate.
  // Note that, because we're in the capture phase, this callback will occur before
  // the bubbling click event on any outside elements.

  return !elements.some(function (el) {
    return el.contains(e.target);
  });
}

function checkIsActive(e, binding) {
  var isActive = _typeof(binding.value) === 'object' && binding.value.closeConditional || defaultConditional;
  return isActive(e);
}

function directive(e, el, binding) {
  var handler = typeof binding.value === 'function' ? binding.value : binding.value.handler;
  el._clickOutside.lastMousedownWasOutside && checkEvent(e, el, binding) && setTimeout(function () {
    checkIsActive(e, binding) && handler && handler(e);
  }, 0);
}

function handleShadow(el, callback) {
  var root = (0, _dom.attachedRoot)(el);
  callback(document);

  if (typeof ShadowRoot !== 'undefined' && root instanceof ShadowRoot) {
    callback(root);
  }
}

var ClickOutside = {
  // [data-app] may not be found
  // if using bind, inserted makes
  // sure that the root element is
  // available, iOS does not support
  // clicks on body
  inserted: function inserted(el, binding, vnode) {
    var onClick = function onClick(e) {
      return directive(e, el, binding);
    };

    var onMousedown = function onMousedown(e) {
      el._clickOutside.lastMousedownWasOutside = checkEvent(e, el, binding);
    };

    handleShadow(el, function (app) {
      app.addEventListener('click', onClick, true);
      app.addEventListener('mousedown', onMousedown, true);
    });

    if (!el._clickOutside) {
      el._clickOutside = {
        lastMousedownWasOutside: true
      };
    }

    el._clickOutside[vnode.context._uid] = {
      onClick: onClick,
      onMousedown: onMousedown
    };
  },
  unbind: function unbind(el, binding, vnode) {
    if (!el._clickOutside) return;
    handleShadow(el, function (app) {
      var _a;

      if (!app || !((_a = el._clickOutside) === null || _a === void 0 ? void 0 : _a[vnode.context._uid])) return;
      var _el$_clickOutside$vno = el._clickOutside[vnode.context._uid],
          onClick = _el$_clickOutside$vno.onClick,
          onMousedown = _el$_clickOutside$vno.onMousedown;
      app.removeEventListener('click', onClick, true);
      app.removeEventListener('mousedown', onMousedown, true);
    });
    delete el._clickOutside[vnode.context._uid];
  }
};
exports.ClickOutside = ClickOutside;
var _default = ClickOutside;
exports.default = _default;
//# sourceMappingURL=index.js.map