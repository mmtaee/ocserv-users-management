function enabled(value: string | undefined): boolean {
  if (value === undefined) return true;
  return value.trim().toLowerCase() === "true";
}

export const telegramBotEnabled = enabled(
  import.meta.env.TELEGRAM_BOT_ENABLED,
);
