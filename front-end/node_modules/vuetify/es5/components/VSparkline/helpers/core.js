"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.genPoints = genPoints;
exports.genBars = genBars;

function _toConsumableArray(arr) { return _arrayWithoutHoles(arr) || _iterableToArray(arr) || _unsupportedIterableToArray(arr) || _nonIterableSpread(); }

function _nonIterableSpread() { throw new TypeError("Invalid attempt to spread non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _iterableToArray(iter) { if (typeof Symbol !== "undefined" && Symbol.iterator in Object(iter)) return Array.from(iter); }

function _arrayWithoutHoles(arr) { if (Array.isArray(arr)) return _arrayLikeToArray(arr); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

function genPoints(values, boundary) {
  var minX = boundary.minX,
      maxX = boundary.maxX,
      minY = boundary.minY,
      maxY = boundary.maxY;
  var totalValues = values.length;
  var maxValue = Math.max.apply(Math, _toConsumableArray(values));
  var minValue = Math.min.apply(Math, _toConsumableArray(values));
  var gridX = (maxX - minX) / (totalValues - 1);
  var gridY = (maxY - minY) / (maxValue - minValue || 1);
  return values.map(function (value, index) {
    return {
      x: minX + index * gridX,
      y: maxY - (value - minValue) * gridY,
      value: value
    };
  });
}

function genBars(values, boundary) {
  var minX = boundary.minX,
      maxX = boundary.maxX,
      minY = boundary.minY,
      maxY = boundary.maxY;
  var totalValues = values.length;
  var maxValue = Math.max.apply(Math, _toConsumableArray(values));
  var minValue = Math.min.apply(Math, _toConsumableArray(values));
  if (minValue > 0) minValue = 0;
  if (maxValue < 0) maxValue = 0;
  var gridX = maxX / totalValues;
  var gridY = (maxY - minY) / (maxValue - minValue || 1);
  var horizonY = maxY - Math.abs(minValue * gridY);
  return values.map(function (value, index) {
    var height = Math.abs(gridY * value);
    return {
      x: minX + index * gridX,
      y: horizonY - height + +(value < 0) * height,
      height: height,
      value: value
    };
  });
}
//# sourceMappingURL=core.js.map