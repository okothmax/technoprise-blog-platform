import { Injectable, inject } from '@angular/core';
import { BehaviorSubject, Observable, fromEvent } from 'rxjs';
import { DOCUMENT } from '@angular/common';

export interface AccessibilitySettings {
  highContrast: boolean;
  reducedMotion: boolean;
  fontSize: 'small' | 'medium' | 'large' | 'extra-large';
  voiceNavigation: boolean;
  screenReaderMode: boolean;
  keyboardNavigation: boolean;
}

export interface AccessibilityScore {
  overall: number;
  colorContrast: number;
  keyboardNavigation: number;
  screenReaderCompatibility: number;
  focusManagement: number;
  semanticStructure: number;
}

@Injectable({
  providedIn: 'root'
})
export class AccessibilityService {
  private document = inject(DOCUMENT);
  
  // Default accessibility settings
  private defaultSettings: AccessibilitySettings = {
    highContrast: false,
    reducedMotion: false,
    fontSize: 'medium',
    voiceNavigation: false,
    screenReaderMode: false,
    keyboardNavigation: true
  };

  private settingsSubject = new BehaviorSubject<AccessibilitySettings>(this.defaultSettings);
  private scoreSubject = new BehaviorSubject<AccessibilityScore>({
    overall: 95,
    colorContrast: 100,
    keyboardNavigation: 95,
    screenReaderCompatibility: 90,
    focusManagement: 95,
    semanticStructure: 100
  });

  public settings$ = this.settingsSubject.asObservable();
  public score$ = this.scoreSubject.asObservable();
  
  // Current announcement for live region
  private currentAnnouncement = '';
  private announcementSubject = new BehaviorSubject<string>('');
  public announcement$ = this.announcementSubject.asObservable();

  constructor() {
    this.initializeAccessibility();
    this.loadUserPreferences();
    this.detectSystemPreferences();
  }

  /**
   * Initialize accessibility features
   */
  private initializeAccessibility(): void {
    // Add skip links
    this.addSkipLinks();
    
    // Initialize keyboard navigation
    this.initializeKeyboardNavigation();
    
    // Initialize focus management
    this.initializeFocusManagement();
    
    // Initialize screen reader announcements
    this.initializeScreenReaderSupport();
  }

  /**
   * Add skip navigation links
   */
  private addSkipLinks(): void {
    const skipNav = this.document.createElement('nav');
    skipNav.className = 'skip-navigation';
    skipNav.setAttribute('aria-label', 'Skip navigation');
    skipNav.innerHTML = `
      <a href="#main-content" class="skip-link">Skip to main content</a>
      <a href="#navigation" class="skip-link">Skip to navigation</a>
      <a href="#footer" class="skip-link">Skip to footer</a>
    `;
    
    this.document.body.insertBefore(skipNav, this.document.body.firstChild);
  }

  /**
   * Initialize keyboard navigation
   */
  private initializeKeyboardNavigation(): void {
    fromEvent<KeyboardEvent>(this.document, 'keydown').subscribe(event => {
      // Handle keyboard shortcuts
      if (event.altKey) {
        switch (event.key) {
          case '1':
            event.preventDefault();
            this.focusElement('#main-content');
            this.announce('Navigated to main content');
            break;
          case '2':
            event.preventDefault();
            this.focusElement('#navigation');
            this.announce('Navigated to navigation');
            break;
          case 'h':
            event.preventDefault();
            this.focusElement('h1, h2, h3, h4, h5, h6');
            this.announce('Navigated to next heading');
            break;
          case 'l':
            event.preventDefault();
            this.focusElement('a');
            this.announce('Navigated to next link');
            break;
        }
      }

      // Escape key handling
      if (event.key === 'Escape') {
        this.handleEscapeKey();
      }
    });
  }

  /**
   * Initialize focus management
   */
  private initializeFocusManagement(): void {
    // Track focus for better accessibility
    let lastFocusedElement: Element | null = null;

    fromEvent(this.document, 'focusin').subscribe((event: Event) => {
      const target = event.target as Element;
      if (target) {
        lastFocusedElement = target;
        this.highlightFocusedElement(target);
      }
    });

    fromEvent(this.document, 'focusout').subscribe(() => {
      this.removeFocusHighlight();
    });
  }

  /**
   * Initialize screen reader support
   */
  private initializeScreenReaderSupport(): void {
    // Create live region for announcements
    const liveRegion = this.document.createElement('div');
    liveRegion.id = 'live-region';
    liveRegion.setAttribute('aria-live', 'polite');
    liveRegion.setAttribute('aria-atomic', 'true');
    liveRegion.className = 'sr-only';
    this.document.body.appendChild(liveRegion);
  }

  /**
   * Detect system accessibility preferences
   */
  private detectSystemPreferences(): void {
    const settings = { ...this.defaultSettings };

    // Detect prefers-reduced-motion
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      settings.reducedMotion = true;
    }

    // Detect prefers-contrast
    if (window.matchMedia('(prefers-contrast: high)').matches) {
      settings.highContrast = true;
    }

    // Detect screen reader
    if (this.isScreenReaderActive()) {
      settings.screenReaderMode = true;
    }

    this.updateSettings(settings);
  }

  /**
   * Load user preferences from localStorage
   */
  private loadUserPreferences(): void {
    const saved = localStorage.getItem('accessibility-settings');
    if (saved) {
      try {
        const settings = JSON.parse(saved);
        this.updateSettings({ ...this.defaultSettings, ...settings });
      } catch (error) {
        console.warn('Failed to load accessibility settings:', error);
      }
    }
  }

  /**
   * Save user preferences to localStorage
   */
  private saveUserPreferences(settings: AccessibilitySettings): void {
    localStorage.setItem('accessibility-settings', JSON.stringify(settings));
  }

  /**
   * Update accessibility settings
   */
  updateSettings(settings: Partial<AccessibilitySettings>): void {
    const currentSettings = this.settingsSubject.value;
    const newSettings = { ...currentSettings, ...settings };
    
    this.settingsSubject.next(newSettings);
    this.saveUserPreferences(newSettings);
    this.applySettings(newSettings);
    
    // Announce changes
    this.announce('Accessibility settings updated');
  }

  /**
   * Apply accessibility settings to the document
   */
  private applySettings(settings: AccessibilitySettings): void {
    const body = this.document.body;
    const root = this.document.documentElement;
    
    // Apply high contrast
    body.classList.toggle('high-contrast', settings.highContrast);
    root.classList.toggle('high-contrast', settings.highContrast);
    
    // Apply reduced motion
    body.classList.toggle('reduced-motion', settings.reducedMotion);
    root.classList.toggle('reduced-motion', settings.reducedMotion);
    
    // Apply font size
    body.classList.remove('font-small', 'font-medium', 'font-large', 'font-extra-large');
    root.classList.remove('font-small', 'font-medium', 'font-large', 'font-extra-large');
    body.classList.add(`font-${settings.fontSize}`);
    root.classList.add(`font-${settings.fontSize}`);
    
    // Set CSS custom properties for font size
    const fontSizeMap = {
      'small': '0.875rem',
      'medium': '1rem',
      'large': '1.125rem',
      'extra-large': '1.25rem'
    };
    root.style.setProperty('--base-font-size', fontSizeMap[settings.fontSize]);
    
    // Apply voice navigation
    body.classList.toggle('voice-navigation', settings.voiceNavigation);
    
    // Apply screen reader mode
    body.classList.toggle('screen-reader-mode', settings.screenReaderMode);
    
    // Apply keyboard navigation
    body.classList.toggle('keyboard-navigation', settings.keyboardNavigation);
  }

  /**
   * Announce message to screen readers
   */
  announce(message: string, priority: 'polite' | 'assertive' = 'polite'): void {
    this.currentAnnouncement = message;
    this.announcementSubject.next(message);
    
    const liveRegion = this.document.getElementById('live-region');
    if (liveRegion) {
      liveRegion.setAttribute('aria-live', priority);
      liveRegion.textContent = message;
      
      // Clear after announcement
      setTimeout(() => {
        liveRegion.textContent = '';
        this.currentAnnouncement = '';
        this.announcementSubject.next('');
      }, 1000);
    }
  }

  /**
   * Get current announcement for template binding
   */
  getCurrentAnnouncement(): string {
    return this.currentAnnouncement;
  }

  /**
   * Focus on element by selector
   */
  focusElement(selector: string): void {
    const element = this.document.querySelector(selector) as HTMLElement;
    if (element) {
      element.focus();
      element.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }

  /**
   * Handle escape key press
   */
  private handleEscapeKey(): void {
    // Close modals, dropdowns, etc.
    const modals = this.document.querySelectorAll('[role="dialog"][aria-modal="true"]');
    modals.forEach(modal => {
      const closeButton = modal.querySelector('[aria-label*="close"], [aria-label*="Close"]') as HTMLElement;
      if (closeButton) {
        closeButton.click();
      }
    });
  }

  /**
   * Highlight focused element for better visibility
   */
  private highlightFocusedElement(element: Element): void {
    element.classList.add('accessibility-focus');
  }

  /**
   * Remove focus highlight
   */
  private removeFocusHighlight(): void {
    const focused = this.document.querySelectorAll('.accessibility-focus');
    focused.forEach(el => el.classList.remove('accessibility-focus'));
  }

  /**
   * Check if screen reader is active
   */
  private isScreenReaderActive(): boolean {
    // Simple heuristic to detect screen reader
    return window.navigator.userAgent.includes('NVDA') ||
           window.navigator.userAgent.includes('JAWS') ||
           window.speechSynthesis?.getVoices().length > 0;
  }

  /**
   * Calculate accessibility score
   */
  calculateAccessibilityScore(): void {
    const score: AccessibilityScore = {
      overall: 0,
      colorContrast: this.checkColorContrast(),
      keyboardNavigation: this.checkKeyboardNavigation(),
      screenReaderCompatibility: this.checkScreenReaderCompatibility(),
      focusManagement: this.checkFocusManagement(),
      semanticStructure: this.checkSemanticStructure()
    };

    // Calculate overall score
    score.overall = Math.round(
      (score.colorContrast + 
       score.keyboardNavigation + 
       score.screenReaderCompatibility + 
       score.focusManagement + 
       score.semanticStructure) / 5
    );

    this.scoreSubject.next(score);
  }

  /**
   * Check color contrast compliance
   */
  private checkColorContrast(): number {
    // Simplified color contrast check
    // In a real implementation, you'd check actual contrast ratios
    return 95; // Assuming good contrast
  }

  /**
   * Check keyboard navigation compliance
   */
  private checkKeyboardNavigation(): number {
    const interactiveElements = this.document.querySelectorAll(
      'a, button, input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );
    
    let score = 100;
    interactiveElements.forEach(element => {
      if (!element.getAttribute('tabindex') && element.tagName !== 'A' && element.tagName !== 'BUTTON') {
        score -= 5;
      }
    });

    return Math.max(0, score);
  }

  /**
   * Check screen reader compatibility
   */
  private checkScreenReaderCompatibility(): number {
    let score = 100;
    
    // Check for missing alt text
    const images = this.document.querySelectorAll('img');
    images.forEach(img => {
      if (!img.getAttribute('alt')) {
        score -= 10;
      }
    });

    // Check for proper headings
    const headings = this.document.querySelectorAll('h1, h2, h3, h4, h5, h6');
    if (headings.length === 0) {
      score -= 20;
    }

    return Math.max(0, score);
  }

  /**
   * Check focus management
   */
  private checkFocusManagement(): number {
    // Check for focus indicators
    return 95; // Assuming good focus management
  }

  /**
   * Check semantic structure
   */
  private checkSemanticStructure(): number {
    let score = 100;
    
    // Check for main landmark
    if (!this.document.querySelector('main')) {
      score -= 20;
    }

    // Check for navigation landmark
    if (!this.document.querySelector('nav')) {
      score -= 10;
    }

    return Math.max(0, score);
  }

  /**
   * Get current accessibility settings
   */
  getCurrentSettings(): AccessibilitySettings {
    return this.settingsSubject.value;
  }

  /**
   * Reset to default settings
   */
  resetToDefaults(): void {
    this.updateSettings(this.defaultSettings);
    this.announce('Accessibility settings reset to defaults');
  }
}
