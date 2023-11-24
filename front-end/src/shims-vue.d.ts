declare module '*.vue' {
  import Vue from 'vue'
  export default Vue
}
// src/shims-vue.d.ts

declare module '*.vue' {
  import Vue from 'vue';
  export default Vue;
}

declare module '*.json' {
  const value: any;
  export default value;
}

declare module '@/*' {
  import { Component } from 'vue';
  const component: Component;
  export default component;
}

// Add your process.env declarations here
declare const process: {
  env: {
    VUE_APP_DOCKERIZED: string;
    NODE_ENV: string,
    BASE_URL: string
    // Add other environment variables as needed
  };
};
