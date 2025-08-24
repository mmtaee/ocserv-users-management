import "vuetify/lib/styles/main.css"
import "@mdi/font/css/materialdesignicons.css"

import type {DisplayThresholds} from "vuetify"
import {createVuetify} from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import {createVueI18nAdapter} from "vuetify/locale/adapters/vue-i18n"
import {useI18n} from "vue-i18n";
import i18n from "@/plugins/i18n.ts";


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
        adapter: createVueI18nAdapter({i18n: i18n as any, useI18n}),
    },
    display: {mobileBreakpoint: 'sm', thresholds: breakpoints},
    icons: {defaultSet: 'mdi'},
    theme: {
        defaultTheme: 'light',
        themes: {
            light: {
                dark: false,
                colors: {
                    primary: '#2C7BE5',      // Bright blue
                    secondary: '#6C757D',    // Cool gray
                    background: '#F8F9FA',   // Very light gray
                    surface: '#FFFFFF',      // White
                    error: '#D32F2F',        // Strong red
                    success: '#388E3C',      // Green
                    warning: '#FFA000',      // Amber
                    info: '#1976D2',         // Deep blue
                    onPrimary: '#FFFFFF',    // Text on primary
                    onSecondary: '#FFFFFF',  // Text on secondary
                    onBackground: '#212529', // Main text
                    onSurface: '#212529',    // Secondary text
                    odd: '#E6F0FF',   // Very light blue (soft, related to primary)
                    even: '#F8F9FA',  // Very light gray (clean, neutral)
                },
            },
            dark: {
                dark: true,
                colors: {
                    primary: '#2C7BE5',       // Same blue
                    secondary: '#ADB5BD',     // Soft gray
                    background: '#121212',    // Dark background
                    surface: '#1E1E1E',       // Slightly lighter than background
                    error: '#EF5350',         // Softer red
                    success: '#81C784',       // Light green
                    warning: '#FFB74D',       // Warm amber
                    info: '#64B5F6',          // Sky blue
                    onPrimary: '#FFFFFF',     // White text
                    onSecondary: '#000000',   // Black text
                    onBackground: '#E0E0E0',  // Light text
                    onSurface: '#E0E0E0',     // Light secondary
                    odd: '#1A1A1A',     // Dark gray
                    even: '#121212',    // Match background
                },
            }
        },
    },
})