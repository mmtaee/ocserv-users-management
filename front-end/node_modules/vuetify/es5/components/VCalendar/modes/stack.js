"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.stack = void 0;

var _common = require("./common");

var _timestamp = require("../util/timestamp");

function _slicedToArray(arr, i) { return _arrayWithHoles(arr) || _iterableToArrayLimit(arr, i) || _unsupportedIterableToArray(arr, i) || _nonIterableRest(); }

function _nonIterableRest() { throw new TypeError("Invalid attempt to destructure non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); }

function _iterableToArrayLimit(arr, i) { if (typeof Symbol === "undefined" || !(Symbol.iterator in Object(arr))) return; var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"] != null) _i["return"](); } finally { if (_d) throw _e; } } return _arr; }

function _arrayWithHoles(arr) { if (Array.isArray(arr)) return arr; }

function _createForOfIteratorHelper(o, allowArrayLike) { var it; if (typeof Symbol === "undefined" || o[Symbol.iterator] == null) { if (Array.isArray(o) || (it = _unsupportedIterableToArray(o)) || allowArrayLike && o && typeof o.length === "number") { if (it) o = it; var i = 0; var F = function F() {}; return { s: F, n: function n() { if (i >= o.length) return { done: true }; return { done: false, value: o[i++] }; }, e: function e(_e2) { throw _e2; }, f: F }; } throw new TypeError("Invalid attempt to iterate non-iterable instance.\nIn order to be iterable, non-array objects must have a [Symbol.iterator]() method."); } var normalCompletion = true, didErr = false, err; return { s: function s() { it = o[Symbol.iterator](); }, n: function n() { var step = it.next(); normalCompletion = step.done; return step; }, e: function e(_e3) { didErr = true; err = _e3; }, f: function f() { try { if (!normalCompletion && it.return != null) it.return(); } finally { if (didErr) throw err; } } }; }

function _unsupportedIterableToArray(o, minLen) { if (!o) return; if (typeof o === "string") return _arrayLikeToArray(o, minLen); var n = Object.prototype.toString.call(o).slice(8, -1); if (n === "Object" && o.constructor) n = o.constructor.name; if (n === "Map" || n === "Set") return Array.from(o); if (n === "Arguments" || /^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)) return _arrayLikeToArray(o, minLen); }

function _arrayLikeToArray(arr, len) { if (len == null || len > arr.length) len = arr.length; for (var i = 0, arr2 = new Array(len); i < len; i++) { arr2[i] = arr[i]; } return arr2; }

var FULL_WIDTH = 100;
var DEFAULT_OFFSET = 5;
var WIDTH_MULTIPLIER = 1.7;
/**
 * Variation of column mode where events can be stacked. The priority of this
 * mode is to stack events together taking up the least amount of space while
 * trying to ensure the content of the event is always visible as well as its
 * start and end. A sibling column has intersecting event content and must be
 * placed beside each other. Non-sibling columns are offset by 5% from the
 * previous column. The width is scaled by 1.7 so the events overlap and
 * whitespace is reduced. If there is a hole in columns the event width is
 * scaled up so it intersects with the next column. The columns have equal
 * width in the space they are given. If the event doesn't have any to the
 * right of it that intersect with it's content it's right side is extended
 * to the right side.
 */

var stack = function stack(events, firstWeekday, overlapThreshold) {
  var handler = (0, _common.getOverlapGroupHandler)(firstWeekday); // eslint-disable-next-line max-statements

  return function (day, dayEvents, timed, reset) {
    if (!timed) {
      return handler.getVisuals(day, dayEvents, timed, reset);
    }

    var dayStart = (0, _timestamp.getTimestampIdentifier)(day);
    var visuals = (0, _common.getVisuals)(dayEvents, dayStart);
    var groups = getGroups(visuals, dayStart);

    var _iterator = _createForOfIteratorHelper(groups),
        _step;

    try {
      for (_iterator.s(); !(_step = _iterator.n()).done;) {
        var group = _step.value;
        var nodes = [];

        var _iterator2 = _createForOfIteratorHelper(group.visuals),
            _step2;

        try {
          for (_iterator2.s(); !(_step2 = _iterator2.n()).done;) {
            var visual = _step2.value;
            var child = getNode(visual, dayStart);
            var index = getNextIndex(child, nodes);

            if (index === false) {
              var parent = getParent(child, nodes);

              if (parent) {
                child.parent = parent;
                child.sibling = (0, _common.hasOverlap)(child.start, child.end, parent.start, addTime(parent.start, overlapThreshold));
                child.index = parent.index + 1;
                parent.children.push(child);
              }
            } else {
              var _getOverlappingRange = getOverlappingRange(child, nodes, index - 1, index - 1),
                  _getOverlappingRange2 = _slicedToArray(_getOverlappingRange, 1),
                  _parent = _getOverlappingRange2[0];

              var children = getOverlappingRange(child, nodes, index + 1, index + nodes.length, true);
              child.children = children;
              child.index = index;

              if (_parent) {
                child.parent = _parent;
                child.sibling = (0, _common.hasOverlap)(child.start, child.end, _parent.start, addTime(_parent.start, overlapThreshold));

                _parent.children.push(child);
              }

              var _iterator3 = _createForOfIteratorHelper(children),
                  _step3;

              try {
                for (_iterator3.s(); !(_step3 = _iterator3.n()).done;) {
                  var grand = _step3.value;

                  if (grand.parent === _parent) {
                    grand.parent = child;
                  }

                  var grandNext = grand.index - child.index <= 1;

                  if (grandNext && child.sibling && (0, _common.hasOverlap)(child.start, addTime(child.start, overlapThreshold), grand.start, grand.end)) {
                    grand.sibling = true;
                  }
                }
              } catch (err) {
                _iterator3.e(err);
              } finally {
                _iterator3.f();
              }
            }

            nodes.push(child);
          }
        } catch (err) {
          _iterator2.e(err);
        } finally {
          _iterator2.f();
        }

        calculateBounds(nodes, overlapThreshold);
      }
    } catch (err) {
      _iterator.e(err);
    } finally {
      _iterator.f();
    }

    visuals.sort(function (a, b) {
      return a.left - b.left || a.event.startTimestampIdentifier - b.event.startTimestampIdentifier;
    });
    return visuals;
  };
};

exports.stack = stack;

function calculateBounds(nodes, overlapThreshold) {
  var _iterator4 = _createForOfIteratorHelper(nodes),
      _step4;

  try {
    for (_iterator4.s(); !(_step4 = _iterator4.n()).done;) {
      var node = _step4.value;
      var visual = node.visual,
          parent = node.parent;
      var columns = getMaxChildIndex(node) + 1;
      var spaceLeft = parent ? parent.visual.left : 0;
      var spaceWidth = FULL_WIDTH - spaceLeft;
      var offset = Math.min(DEFAULT_OFFSET, FULL_WIDTH / columns);
      var columnWidthMultiplier = getColumnWidthMultiplier(node, nodes);
      var columnOffset = spaceWidth / (columns - node.index + 1);
      var columnWidth = spaceWidth / (columns - node.index + (node.sibling ? 1 : 0)) * columnWidthMultiplier;

      if (parent) {
        visual.left = node.sibling ? spaceLeft + columnOffset : spaceLeft + offset;
      }

      visual.width = hasFullWidth(node, nodes, overlapThreshold) ? FULL_WIDTH - visual.left : Math.min(FULL_WIDTH - visual.left, columnWidth * WIDTH_MULTIPLIER);
    }
  } catch (err) {
    _iterator4.e(err);
  } finally {
    _iterator4.f();
  }
}

function getColumnWidthMultiplier(node, nodes) {
  if (!node.children.length) {
    return 1;
  }

  var maxColumn = node.index + nodes.length;
  var minColumn = node.children.reduce(function (min, c) {
    return Math.min(min, c.index);
  }, maxColumn);
  return minColumn - node.index;
}

function getOverlappingIndices(node, nodes) {
  var indices = [];

  var _iterator5 = _createForOfIteratorHelper(nodes),
      _step5;

  try {
    for (_iterator5.s(); !(_step5 = _iterator5.n()).done;) {
      var other = _step5.value;

      if ((0, _common.hasOverlap)(node.start, node.end, other.start, other.end)) {
        indices.push(other.index);
      }
    }
  } catch (err) {
    _iterator5.e(err);
  } finally {
    _iterator5.f();
  }

  return indices;
}

function getNextIndex(node, nodes) {
  var indices = getOverlappingIndices(node, nodes);
  indices.sort();

  for (var i = 0; i < indices.length; i++) {
    if (i < indices[i]) {
      return i;
    }
  }

  return false;
}

function getOverlappingRange(node, nodes, indexMin, indexMax) {
  var returnFirstColumn = arguments.length > 4 && arguments[4] !== undefined ? arguments[4] : false;
  var overlapping = [];

  var _iterator6 = _createForOfIteratorHelper(nodes),
      _step6;

  try {
    for (_iterator6.s(); !(_step6 = _iterator6.n()).done;) {
      var other = _step6.value;

      if (other.index >= indexMin && other.index <= indexMax && (0, _common.hasOverlap)(node.start, node.end, other.start, other.end)) {
        overlapping.push(other);
      }
    }
  } catch (err) {
    _iterator6.e(err);
  } finally {
    _iterator6.f();
  }

  if (returnFirstColumn && overlapping.length > 0) {
    var first = overlapping.reduce(function (min, n) {
      return Math.min(min, n.index);
    }, overlapping[0].index);
    return overlapping.filter(function (n) {
      return n.index === first;
    });
  }

  return overlapping;
}

function getParent(node, nodes) {
  var parent = null;

  var _iterator7 = _createForOfIteratorHelper(nodes),
      _step7;

  try {
    for (_iterator7.s(); !(_step7 = _iterator7.n()).done;) {
      var other = _step7.value;

      if ((0, _common.hasOverlap)(node.start, node.end, other.start, other.end) && (parent === null || other.index > parent.index)) {
        parent = other;
      }
    }
  } catch (err) {
    _iterator7.e(err);
  } finally {
    _iterator7.f();
  }

  return parent;
}

function hasFullWidth(node, nodes, overlapThreshold) {
  var _iterator8 = _createForOfIteratorHelper(nodes),
      _step8;

  try {
    for (_iterator8.s(); !(_step8 = _iterator8.n()).done;) {
      var other = _step8.value;

      if (other !== node && other.index > node.index && (0, _common.hasOverlap)(node.start, addTime(node.start, overlapThreshold), other.start, other.end)) {
        return false;
      }
    }
  } catch (err) {
    _iterator8.e(err);
  } finally {
    _iterator8.f();
  }

  return true;
}

function getGroups(visuals, dayStart) {
  var groups = [];

  var _iterator9 = _createForOfIteratorHelper(visuals),
      _step9;

  try {
    for (_iterator9.s(); !(_step9 = _iterator9.n()).done;) {
      var visual = _step9.value;

      var _getNormalizedRange = (0, _common.getNormalizedRange)(visual.event, dayStart),
          _getNormalizedRange2 = _slicedToArray(_getNormalizedRange, 2),
          start = _getNormalizedRange2[0],
          end = _getNormalizedRange2[1];

      var added = false;

      var _iterator10 = _createForOfIteratorHelper(groups),
          _step10;

      try {
        for (_iterator10.s(); !(_step10 = _iterator10.n()).done;) {
          var group = _step10.value;

          if ((0, _common.hasOverlap)(start, end, group.start, group.end)) {
            group.visuals.push(visual);
            group.end = Math.max(group.end, end);
            added = true;
            break;
          }
        }
      } catch (err) {
        _iterator10.e(err);
      } finally {
        _iterator10.f();
      }

      if (!added) {
        groups.push({
          start: start,
          end: end,
          visuals: [visual]
        });
      }
    }
  } catch (err) {
    _iterator9.e(err);
  } finally {
    _iterator9.f();
  }

  return groups;
}

function getNode(visual, dayStart) {
  var _getNormalizedRange3 = (0, _common.getNormalizedRange)(visual.event, dayStart),
      _getNormalizedRange4 = _slicedToArray(_getNormalizedRange3, 2),
      start = _getNormalizedRange4[0],
      end = _getNormalizedRange4[1];

  return {
    parent: null,
    sibling: true,
    index: 0,
    visual: visual,
    start: start,
    end: end,
    children: []
  };
}

function getMaxChildIndex(node) {
  var max = node.index;

  var _iterator11 = _createForOfIteratorHelper(node.children),
      _step11;

  try {
    for (_iterator11.s(); !(_step11 = _iterator11.n()).done;) {
      var child = _step11.value;
      var childMax = getMaxChildIndex(child);

      if (childMax > max) {
        max = childMax;
      }
    }
  } catch (err) {
    _iterator11.e(err);
  } finally {
    _iterator11.f();
  }

  return max;
}

function addTime(identifier, minutes) {
  var removeMinutes = identifier % 100;
  var totalMinutes = removeMinutes + minutes;
  var addHours = Math.floor(totalMinutes / 60);
  var addMinutes = totalMinutes % 60;
  return identifier - removeMinutes + addHours * 100 + addMinutes;
}
//# sourceMappingURL=stack.js.map