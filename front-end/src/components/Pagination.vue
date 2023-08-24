<template>
  <v-row align="center" justify="start" class="my-1 mx-1 pa-0">
    <v-col class="ma-0 pa-0" cols="auto">
      <span class="font-weight-medium me-2 text-subtitle-2 text--secondary">
        Item per page
      </span>
    </v-col>
    <v-col class="ma-0 pa-0" cols="1">
      <v-select
        v-model="itemPerPage"
        :items="[10, 20, 30, 50, 100]"
        class="small-size my-0 py-0"
        item-text="text"
        item-value="value"
        hide-details
        dense
        outlined
        @change="$emit('changePage', 1, itemPerPage)"
      >
        <template v-slot:selection="{ item }">
          <span
            class="font-weight-medium text-caption text--secondary"
            v-text="item"
          />
        </template>

        <template v-slot:item="{ item }">
          <span
            class="font-weight-medium text-caption text--secondary"
            v-text="item"
          />
        </template>
      </v-select>
    </v-col>
    <v-col cols="5" class="ma-0 pa-0 text-start" v-if="pages > 1">
      <v-pagination
        v-model="currentPage"
        class="ma-0 pa-0 small-size-page"
        :length="pages"
        :total-visible="7"
        @next="$emit('changePage', currentPage, itemPerPage)"
        @previous="$emit('changePage', currentPage, itemPerPage)"
        @input="$emit('changePage', currentPage, itemPerPage)"
      />
    </v-col>
    <v-spacer />
    <v-col v-if="!condition && Boolean(count)" class="ma-0 pa-0" cols="auto">
      <v-chip class="ma-0" color="primary" label outlined small>
        <span class="font-weight-medium"> Total item: {{ this.count }} </span>
      </v-chip>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "Pagination",

  props: {
    page: Number,
    pages: Number,
    perPage: {
      type: Number,
      default: 0,
    },
    count: {
      type: Number,
      default: 0,
    },
  },

  data(): {
    itemPerPage: number;
    currentPage: number;
  } {
    return {
      itemPerPage: 10,
      currentPage: 1,
    };
  },

  watch: {
    page: {
      immediate: true,
      handler() {
        this.currentPage = this.page;
      },
    },
    perPage: {
      immediate: true,
      handler() {
        this.itemPerPage = this.perPage;
      },
    },
  },

  computed: {
    condition() {
      if (this.itemPerPage == 10 && this.pages == 1) {
        return true;
      }
      return false;
    },
  },
});
</script>

<style>
.deactivateClassDefault {
  color: #0000005e !important;
  cursor: default;
}
.updatedClassDefault {
  background-color: #2cbe6020;
  cursor: default;
}
.notAvailableClassDefault {
  text-decoration: line-through !important;
  text-decoration-color: #ef3e36de !important;
  background-color: #f5bc0020;
  cursor: default;
}

.classDefault {
  color: #000000de !important;
  cursor: default;
}
.selectedActiveClassDefault {
  background-color: #2e4a7620 !important;
  color: #000000de !important;
  cursor: default;
}
.pointer {
  cursor: pointer !important;
}
/* v-icon.v-icon--dense */
.block-display {
  display: block !important;
}
.span-chip {
  display: block !important;
  max-height: 24px !important;
  min-height: 24px !important;
  min-width: 90px !important;
  max-width: 90px !important;
  border: 1px solid #9e9e9e !important;
  border-radius: 4px !important;
  padding: auto !important;
  /* text-justify: auto; */
}
.button-icon {
  border: 1px solid #2cbe60 !important;
  border-radius: 4px !important;
  padding: 2px !important;
}
.activeClass {
  color: #000000de !important;
  cursor: pointer !important;
}
.cursorDefault {
  cursor: default !important;
}
.pointer {
  cursor: pointer !important;
}
.small-size .v-input__slot {
  max-height: 28px !important;
  min-height: 28px !important;
  font-size: 0.75rem !important;
}
.small-size .v-input {
  font-size: 0.75rem !important;
}
.small-size .v-input__append-inner {
  margin-top: 2.8px !important;
}
.small-size .v-input__prepend-inner {
  margin-top: 2.8px !important;
}
.small-size .v-icon.v-icon {
  font-size: 18px !important;
}
.small-size-page .v-pagination .v-pagination__item {
  /* max-width: 27px !important; */
  /* min-width: 27px !important; */
  max-height: 27px !important;
  min-height: 27px !important;
  font-size: 0.75rem !important;
  box-shadow: none !important;
  /* color: #0099FF !important; */
  font-weight: 500 !important;
  border: solid 0.5px #0099ff !important;
}
.small-size-page .v-pagination .v-pagination__navigation {
  /* max-width: 27px !important;
  min-width: 27px !important; */
  max-height: 27px !important;
  min-height: 27px !important;
  font-size: 0.75rem !important;
  box-shadow: none !important;
  /* color: #0099FF !important; */
  border: solid 0.5px #0099ff !important;
  padding-top: 1px !important;
  padding-top: 1px !important;
}
.small-size-page .v-pagination .v-pagination__navigation .v-icon {
  color: #0099ff !important;
}
.small-size-page .v-pagination {
  width: auto !important;
}
.goldClass {
  background-color: #d4af3720;
  cursor: pointer !important;
}
.silverClass {
  background-color: #b5b5bd20;
  cursor: pointer !important;
}
.bronzClass {
  background-color: #cd7f3220;
  cursor: pointer !important;
}

.small-list-item {
  min-height: 32px !important;
  max-height: 32px !important;
}

/* .test .v-tooltip__content {
  padding: 5px 0px !important;
} */
</style>

