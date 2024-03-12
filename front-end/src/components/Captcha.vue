<template>
  <div
    ref="recaptcha"
    class="g-recaptcha"
    data-callback="callbackSuccess"
    data-error-callback="callbackError"
    data-expired-callback="callbackExpired"
    :data-sitekey="sitekey"
    :data-size="compact ? 'compact' : 'normal'"
    :data-theme="$vuetify.theme.dark ? 'dark' : 'light'"
  />
</template>
<script>
export default {
  name: "CaptchaV2",

  props: {
    compact: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      sitekey: this.$store.getters.getGoogleSiteKey,
    };
  },

  methods: {
    callbackSuccess(token) {
      this.$emit("token", token);
      this.$emit("validForm", true);
    },

    callbackError() {
      window.grecaptcha.reset();
      this.$emit("validForm", false);
    },

    callbackExpired() {
      this.$emit("validForm", false);
    },

    destroyRecaptcha() {
      try {
        document.getElementById("captcha").remove();
        this.$refs.recaptcha.remove();
      } catch {}
    },
  },

  mounted() {
    this.recaptchaScript = document.createElement("script");
    this.recaptchaScript.setAttribute(
      "src",
      "https://www.google.com/recaptcha/api.js"
    );
    this.recaptchaScript.setAttribute("async", "");
    this.recaptchaScript.setAttribute("defer", "");
    this.recaptchaScript.setAttribute("id", "captcha");
    document.head.appendChild(this.recaptchaScript);
    window.callbackSuccess = this.callbackSuccess;
    window.callbackError = this.callbackError;
    window.callbackExpired = this.callbackExpired;
  },
};
</script>
<!-- <Captcha ref="recaptcha" @token="function" @valid="function" compact /> -->
<!-- <button @click="$refs.recaptcha.destroyRecaptcha()">click</button> -->
