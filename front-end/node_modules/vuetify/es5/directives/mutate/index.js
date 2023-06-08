"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = exports.Mutate = void 0;

function _objectWithoutProperties(source, excluded) { if (source == null) return {}; var target = _objectWithoutPropertiesLoose(source, excluded); var key, i; if (Object.getOwnPropertySymbols) { var sourceSymbolKeys = Object.getOwnPropertySymbols(source); for (i = 0; i < sourceSymbolKeys.length; i++) { key = sourceSymbolKeys[i]; if (excluded.indexOf(key) >= 0) continue; if (!Object.prototype.propertyIsEnumerable.call(source, key)) continue; target[key] = source[key]; } } return target; }

function _objectWithoutPropertiesLoose(source, excluded) { if (source == null) return {}; var target = {}; var sourceKeys = Object.keys(source); var key, i; for (i = 0; i < sourceKeys.length; i++) { key = sourceKeys[i]; if (excluded.indexOf(key) >= 0) continue; target[key] = source[key]; } return target; }

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

function inserted(el, binding, vnode) {
  var modifiers = binding.modifiers || {};
  var value = binding.value;
  var callback = _typeof(value) === 'object' ? value.handler : value;

  var once = modifiers.once,
      modifierKeys = _objectWithoutProperties(modifiers, ["once"]);

  var hasModifiers = Object.keys(modifierKeys).length > 0; // Options take top priority

  var options = _typeof(value) === 'object' && value.options ? value.options : hasModifiers // If we have modifiers, use only those provided
  ? {
    attributes: modifierKeys.attr,
    childList: modifierKeys.child,
    subtree: modifierKeys.sub,
    characterData: modifierKeys.char
  } // Defaults to everything on
  : {
    attributes: true,
    childList: true,
    subtree: true,
    characterData: true
  };
  var observer = new MutationObserver(function (mutationsList, observer) {
    /* istanbul ignore if */
    if (!el._mutate) return; // Just in case, should never fire

    callback(mutationsList, observer); // If has the once modifier, unbind

    once && unbind(el, binding, vnode);
  });
  observer.observe(el, options);
  el._mutate = Object(el._mutate);
  el._mutate[vnode.context._uid] = {
    observer: observer
  };
}

function unbind(el, binding, vnode) {
  var _a;

  if (!((_a = el._mutate) === null || _a === void 0 ? void 0 : _a[vnode.context._uid])) return;

  el._mutate[vnode.context._uid].observer.disconnect();

  delete el._mutate[vnode.context._uid];
}

var Mutate = {
  inserted: inserted,
  unbind: unbind
};
exports.Mutate = Mutate;
var _default = Mutate;
exports.default = _default;
//# sourceMappingURL=index.js.map