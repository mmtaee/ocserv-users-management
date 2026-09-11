export const telegramFeatureScenarios = {
  enabled: {
    telegramBotEnabled: true,
    navigationVisible: true,
    routeAvailable: true,
  },
  disabled: {
    telegramBotEnabled: false,
    navigationVisible: false,
    routeAvailable: false,
  },
} as const;
