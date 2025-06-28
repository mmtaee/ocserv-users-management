import "vuetify/lib/styles/main.css"
import "@mdi/font/css/materialdesignicons.css"

import type {DisplayThresholds} from "vuetify"
import {createVuetify} from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import i18n from "@/plugins/i18n"
import {createVueI18nAdapter} from "vuetify/locale/adapters/vue-i18n"
import {useI18n} from "vue-i18n"


const breakpoints: DisplayThresholds = {
    xs: 0,
    sm: 600,
    md: 960,
    lg: 1280,
    xl: 1920,
    xxl: 2560
}

export default createVuetify({
    components,
    directives,
    locale: {
        adapter: createVueI18nAdapter({i18n, useI18n}),
    },
    display: {
        mobileBreakpoint: 'sm',
        thresholds: breakpoints
    },
    icons: {
        defaultSet: 'mdi', // This is already the default value - only for display purposes
    },
    // icons: {
    //     defaultSet: 'mdi',
    //     aliases,
    //     sets: {mdi},
    // },
    theme: {
        defaultTheme: 'light', // or 'dark' depending on your needs
        themes: {
            light: {
                dark: false,
                colors: {
                    background: '#FDFDFD',
                    surface: '#FFFFFF',
                    primary: '#7C3AED',          // Vivid Violet
                    'primary-darken-1': '#6D28D9',
                    secondary: '#F97316',        // Vibrant Orange
                    accent: '#06B6D4',           // Electric Cyan
                    error: '#EF4444',            // Bright Red
                    warning: '#FACC15',          // Punchy Yellow
                    success: '#10B981',          // Lush Green
                    info: '#3B82F6',             // Azure Blue
                    odd: '#8f8585',
                    'on-primary': '#FFFFFF',
                    'on-surface': '#1E1E2F',
                    'on-background': '#111827',
                },
            },
            dark: {
                dark: true,
                colors: {
                    background: '#0F172A',         // Deep Navy
                    surface: '#1E293B',            // Dark Slate
                    primary: '#C084FC',            // Electric Purple
                    'primary-darken-1': '#A855F7',
                    secondary: '#FB923C',          // Lively Orange
                    accent: '#22D3EE',             // Aqua Neon
                    error: '#F87171',              // Bright Coral
                    warning: '#FCD34D',            // Sunny Yellow
                    success: '#34D399',            // Green Pulse
                    info: '#60A5FA',               // Sky Blue
                    odd: '#8f8585',
                    'on-primary': '#000000',
                    'on-surface': '#F3F4F6',
                    'on-background': '#E5E7EB',
                },
            }
        },
    },
})