// Mixins
import Measurable from '../../mixins/measurable';
import Toggleable from '../../mixins/toggleable'; // Directives

import intersect from '../../directives/intersect'; // Utilities

import mixins from '../../util/mixins';
import { getSlot } from '../../util/helpers';
export default mixins(Measurable, Toggleable).extend({
  name: 'VLazy',
  directives: {
    intersect
  },
  props: {
    options: {
      type: Object,
      // For more information on types, navigate to:
      // https://developer.mozilla.org/en-US/docs/Web/API/Intersection_Observer_API
      default: () => ({
        root: undefined,
        rootMargin: undefined,
        threshold: undefined
      })
    },
    tag: {
      type: String,
      default: 'div'
    },
    transition: {
      type: String,
      default: 'fade-transition'
    }
  },
  computed: {
    styles() {
      return { ...this.measurableStyles
      };
    }

  },
  methods: {
    genContent() {
      const children = this.isActive && getSlot(this);
      return this.transition ? this.$createElement('transition', {
        props: {
          name: this.transition
        }
      }, children) : children;
    },

    onObserve(entries, observer, isIntersecting) {
      if (this.isActive) return;
      this.isActive = isIntersecting;
    }

  },

  render(h) {
    return h(this.tag, {
      staticClass: 'v-lazy',
      attrs: this.$attrs,
      directives: [{
        name: 'intersect',
        value: {
          handler: this.onObserve,
          options: this.options
        }
      }],
      on: this.$listeners,
      style: this.styles
    }, [this.genContent()]);
  }

});
//# sourceMappingURL=VLazy.js.map