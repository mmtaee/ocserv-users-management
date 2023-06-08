"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = void 0;

require("../../../src/components/VAutocomplete/VAutocomplete.sass");

var _VSelect = _interopRequireDefault(require("../VSelect/VSelect"));

var _VAutocomplete2 = _interopRequireDefault(require("../VAutocomplete/VAutocomplete"));

var _helpers = require("../../util/helpers");

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function _typeof(obj) { "@babel/helpers - typeof"; if (typeof Symbol === "function" && typeof Symbol.iterator === "symbol") { _typeof = function _typeof(obj) { return typeof obj; }; } else { _typeof = function _typeof(obj) { return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj; }; } return _typeof(obj); }

function ownKeys(object, enumerableOnly) { var keys = Object.keys(object); if (Object.getOwnPropertySymbols) { var symbols = Object.getOwnPropertySymbols(object); if (enumerableOnly) symbols = symbols.filter(function (sym) { return Object.getOwnPropertyDescriptor(object, sym).enumerable; }); keys.push.apply(keys, symbols); } return keys; }

function _objectSpread(target) { for (var i = 1; i < arguments.length; i++) { var source = arguments[i] != null ? arguments[i] : {}; if (i % 2) { ownKeys(Object(source), true).forEach(function (key) { _defineProperty(target, key, source[key]); }); } else if (Object.getOwnPropertyDescriptors) { Object.defineProperties(target, Object.getOwnPropertyDescriptors(source)); } else { ownKeys(Object(source)).forEach(function (key) { Object.defineProperty(target, key, Object.getOwnPropertyDescriptor(source, key)); }); } } return target; }

function _defineProperty(obj, key, value) { if (key in obj) { Object.defineProperty(obj, key, { value: value, enumerable: true, configurable: true, writable: true }); } else { obj[key] = value; } return obj; }

/* @vue/component */
var _default2 = _VAutocomplete2.default.extend({
  name: 'v-combobox',
  props: {
    delimiters: {
      type: Array,
      default: function _default() {
        return [];
      }
    },
    returnObject: {
      type: Boolean,
      default: true
    }
  },
  data: function data() {
    return {
      editingIndex: -1
    };
  },
  computed: {
    computedCounterValue: function computedCounterValue() {
      return this.multiple ? this.selectedItems.length : (this.internalSearch || '').toString().length;
    },
    hasSlot: function hasSlot() {
      return _VSelect.default.options.computed.hasSlot.call(this) || this.multiple;
    },
    isAnyValueAllowed: function isAnyValueAllowed() {
      return true;
    },
    menuCanShow: function menuCanShow() {
      if (!this.isFocused) return false;
      return this.hasDisplayedItems || !!this.$slots['no-data'] && !this.hideNoData;
    },
    searchIsDirty: function searchIsDirty() {
      return this.internalSearch != null;
    }
  },
  methods: {
    onInternalSearchChanged: function onInternalSearchChanged(val) {
      if (val && this.multiple && this.delimiters.length) {
        var delimiter = this.delimiters.find(function (d) {
          return val.endsWith(d);
        });

        if (delimiter != null) {
          this.internalSearch = val.slice(0, val.length - delimiter.length);
          this.updateTags();
        }
      }

      this.updateMenuDimensions();
    },
    genInput: function genInput() {
      var input = _VAutocomplete2.default.options.methods.genInput.call(this);

      delete input.data.attrs.name;
      input.data.on.paste = this.onPaste;
      return input;
    },
    genChipSelection: function genChipSelection(item, index) {
      var _this = this;

      var chip = _VSelect.default.options.methods.genChipSelection.call(this, item, index); // Allow user to update an existing value


      if (this.multiple) {
        chip.componentOptions.listeners = _objectSpread(_objectSpread({}, chip.componentOptions.listeners), {}, {
          dblclick: function dblclick() {
            _this.editingIndex = index;
            _this.internalSearch = _this.getText(item);
            _this.selectedIndex = -1;
          }
        });
      }

      return chip;
    },
    onChipInput: function onChipInput(item) {
      _VSelect.default.options.methods.onChipInput.call(this, item);

      this.editingIndex = -1;
    },
    // Requires a manual definition
    // to overwrite removal in v-autocomplete
    onEnterDown: function onEnterDown(e) {
      e.preventDefault(); // If has menu index, let v-select-list handle

      if (this.getMenuIndex() > -1) return;
      this.$nextTick(this.updateSelf);
    },
    onKeyDown: function onKeyDown(e) {
      var keyCode = e.keyCode;

      if (e.ctrlKey || ![_helpers.keyCodes.home, _helpers.keyCodes.end].includes(keyCode)) {
        _VSelect.default.options.methods.onKeyDown.call(this, e);
      } // If user is at selection index of 0
      // create a new tag


      if (this.multiple && keyCode === _helpers.keyCodes.left && this.$refs.input.selectionStart === 0) {
        this.updateSelf();
      } else if (keyCode === _helpers.keyCodes.enter) {
        this.onEnterDown(e);
      } // The ordering is important here
      // allows new value to be updated
      // and then moves the index to the
      // proper location


      this.changeSelectedIndex(keyCode);
    },
    onTabDown: function onTabDown(e) {
      // When adding tags, if searching and
      // there is not a filtered options,
      // add the value to the tags list
      if (this.multiple && this.internalSearch && this.getMenuIndex() === -1) {
        e.preventDefault();
        e.stopPropagation();
        return this.updateTags();
      }

      _VAutocomplete2.default.options.methods.onTabDown.call(this, e);
    },
    selectItem: function selectItem(item) {
      // Currently only supports items:<string[]>
      if (this.editingIndex > -1) {
        this.updateEditing();
      } else {
        _VAutocomplete2.default.options.methods.selectItem.call(this, item); // if selected item contains search value,
        // remove the search string


        if (this.internalSearch && this.multiple && this.getText(item).toLocaleLowerCase().includes(this.internalSearch.toLocaleLowerCase())) {
          this.internalSearch = null;
        }
      }
    },
    setSelectedItems: function setSelectedItems() {
      if (this.internalValue == null || this.internalValue === '') {
        this.selectedItems = [];
      } else {
        this.selectedItems = this.multiple ? this.internalValue : [this.internalValue];
      }
    },
    setValue: function setValue(value) {
      _VSelect.default.options.methods.setValue.call(this, value === undefined ? this.internalSearch : value);
    },
    updateEditing: function updateEditing() {
      var _this2 = this;

      var value = this.internalValue.slice();
      var index = this.selectedItems.findIndex(function (item) {
        return _this2.getText(item) === _this2.internalSearch;
      }); // If user enters a duplicate text on chip edit,
      // don't add it, move it to the end of the list

      if (index > -1) {
        var item = _typeof(value[index]) === 'object' ? Object.assign({}, value[index]) : value[index];
        value.splice(index, 1);
        value.push(item);
      } else {
        value[this.editingIndex] = this.internalSearch;
      }

      this.setValue(value);
      this.editingIndex = -1;
      this.internalSearch = null;
    },
    updateCombobox: function updateCombobox() {
      // If search is not dirty, do nothing
      if (!this.searchIsDirty) return; // The internal search is not matching
      // the internal value, update the input

      if (this.internalSearch !== this.getText(this.internalValue)) this.setValue(); // Reset search if using slot to avoid a double input

      var isUsingSlot = Boolean(this.$scopedSlots.selection) || this.hasChips;
      if (isUsingSlot) this.internalSearch = null;
    },
    updateSelf: function updateSelf() {
      this.multiple ? this.updateTags() : this.updateCombobox();
    },
    updateTags: function updateTags() {
      var _this3 = this;

      var menuIndex = this.getMenuIndex(); // If the user is not searching
      // and no menu item is selected
      // or if the search is empty
      // do nothing

      if (menuIndex < 0 && !this.searchIsDirty || !this.internalSearch) return;

      if (this.editingIndex > -1) {
        return this.updateEditing();
      }

      var index = this.selectedItems.findIndex(function (item) {
        return _this3.internalSearch === _this3.getText(item);
      }); // If the duplicate item is an object,
      // copy it, so that it can be added again later

      var itemToSelect = index > -1 && _typeof(this.selectedItems[index]) === 'object' ? Object.assign({}, this.selectedItems[index]) : this.internalSearch; // If it already exists, do nothing
      // this might need to change to bring
      // the duplicated item to the last entered

      if (index > -1) {
        var internalValue = this.internalValue.slice();
        internalValue.splice(index, 1);
        this.setValue(internalValue);
      } // If menu index is greater than 1
      // the selection is handled elsewhere
      // TODO: find out where


      if (menuIndex > -1) return this.internalSearch = null;
      this.selectItem(itemToSelect);
      this.internalSearch = null;
    },
    onPaste: function onPaste(event) {
      var _a;

      this.$emit('paste', event);
      if (!this.multiple || this.searchIsDirty) return;
      var pastedItemText = (_a = event.clipboardData) === null || _a === void 0 ? void 0 : _a.getData('text/vnd.vuetify.autocomplete.item+plain');

      if (pastedItemText && this.findExistingIndex(pastedItemText) === -1) {
        event.preventDefault();

        _VSelect.default.options.methods.selectItem.call(this, pastedItemText);
      }
    },
    clearableCallback: function clearableCallback() {
      this.editingIndex = -1;

      _VAutocomplete2.default.options.methods.clearableCallback.call(this);
    }
  }
});

exports.default = _default2;
//# sourceMappingURL=VCombobox.js.map