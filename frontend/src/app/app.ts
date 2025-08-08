import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet, RouterLink, RouterLinkActive, Router } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { AccessibilityService } from './services/accessibility.service';
import { BlogService } from './services/blog.service';

@Component({
  selector: 'app-root',
  imports: [
    RouterOutlet,
    RouterLink,
    RouterLinkActive,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatMenuModule,
    MatTooltipModule,
    MatFormFieldModule,
    MatInputModule,
    HttpClientModule
  ],
  templateUrl: './app.html',
  styleUrls: ['./app.scss']
})
export class AppComponent {
  title = 'TechnoPrise Global Blog';
  
  // Search functionality
  searchQuery = signal('');
  private router = inject(Router);
  private blogService = inject(BlogService);

  constructor(public accessibilityService: AccessibilityService) {
    // Announce app initialization
    this.accessibilityService.announce('TechnoPrise Global Blog Platform loaded successfully');
  }
  
  onSearchInput(event: Event) {
    const target = event.target as HTMLInputElement;
    const query = target?.value?.trim() || '';
    this.searchQuery.set(query);
    
    // Update the blog service search query to trigger search in blog-home component
    this.blogService.setSearchQuery(query);
    
    // Navigate to home page if not already there to show search results
    if (this.router.url !== '/') {
      this.router.navigate(['/']);
    }
    
    // Announce search action for screen readers
    if (query) {
      this.accessibilityService.announce(`Searching for: ${query}`);
    } else {
      this.accessibilityService.announce('Search cleared, showing all articles');
    }
  }
  
  navigateToCreateBlog() {
    this.router.navigate(['/create-blog']);
    this.accessibilityService.announce('Navigating to create new blog post');
  }

  toggleHighContrast() {
    const currentSettings = this.accessibilityService.getCurrentSettings();
    this.accessibilityService.updateSettings({ 
      highContrast: !currentSettings.highContrast 
    });
    const isEnabled = !currentSettings.highContrast;
    this.accessibilityService.announce(
      `High contrast mode ${isEnabled ? 'enabled' : 'disabled'}`
    );
  }

  toggleReducedMotion() {
    const currentSettings = this.accessibilityService.getCurrentSettings();
    this.accessibilityService.updateSettings({ 
      reducedMotion: !currentSettings.reducedMotion 
    });
    const isEnabled = !currentSettings.reducedMotion;
    this.accessibilityService.announce(
      `Reduced motion ${isEnabled ? 'enabled' : 'disabled'}`
    );
  }

  increaseFontSize() {
    const currentSettings = this.accessibilityService.getCurrentSettings();
    const fontSizes: Array<'small' | 'medium' | 'large' | 'extra-large'> = ['small', 'medium', 'large', 'extra-large'];
    const currentIndex = fontSizes.indexOf(currentSettings.fontSize);
    const newIndex = Math.min(currentIndex + 1, fontSizes.length - 1);
    
    this.accessibilityService.updateSettings({ 
      fontSize: fontSizes[newIndex] 
    });
    this.accessibilityService.announce(
      `Font size increased to ${fontSizes[newIndex]}`
    );
  }

  decreaseFontSize() {
    const currentSettings = this.accessibilityService.getCurrentSettings();
    const fontSizes: Array<'small' | 'medium' | 'large' | 'extra-large'> = ['small', 'medium', 'large', 'extra-large'];
    const currentIndex = fontSizes.indexOf(currentSettings.fontSize);
    const newIndex = Math.max(currentIndex - 1, 0);
    
    this.accessibilityService.updateSettings({ 
      fontSize: fontSizes[newIndex] 
    });
    this.accessibilityService.announce(
      `Font size decreased to ${fontSizes[newIndex]}`
    );
  }
}
