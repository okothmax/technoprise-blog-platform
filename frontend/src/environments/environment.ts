export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api/v1',
  appName: 'TechnoPrise Blog',
  version: '1.0.0',
  features: {
    voiceNavigation: true,
    highContrast: true,
    reducedMotion: true,
    screenReaderOptimization: true,
    keyboardNavigation: true,
    realTimeAccessibilityScoring: true
  },
  accessibility: {
    wcagLevel: 'AA',
    supportedLanguages: ['en', 'es', 'fr', 'de'],
    announcements: true,
    focusManagement: true
  }
};
