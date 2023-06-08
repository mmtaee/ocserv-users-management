module.exports = (api, options, rootOptions) => {
    api.extendPackage({
      devDependencies: {
        axios: "^0.18.0"
      }
    });
  
    // Render vuetify plugin file
    api.render(
      {
        "./src/plugins/axios.js": "./templates/plugins/axios.js"
      },
      options
    );
  
    const fs = require("fs");
    const helpers = require('./helpers')(api)
  
    // adapted from https://github.com/Akryum/vue-cli-plugin-apollo/blob/master/generator/index.js#L68-L91
    api.onCreateComplete(() => {
      // Modify main.js
      helpers.updateMain(src => {
        let vueImportIndex = src.findIndex(line => line.match(/^import Vue/));
  
        let axiosImportIndex = src.findIndex(line => line.match(/\/plugins\/axios/));
        if(axiosImportIndex < 0){
          src.splice(++vueImportIndex, 0, "import './plugins/axios'");
        }      
        return src;
      });
    });
  };
  