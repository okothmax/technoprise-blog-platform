import { Component, OnInit, OnDestroy, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatCardModule } from '@angular/material/card';
import { MatTooltipModule } from '@angular/material/tooltip';
import { Subject, takeUntil } from 'rxjs';
import { Title, Meta } from '@angular/platform-browser';

import { BlogService, Blog } from '../../services/blog.service';
import { AccessibilityService } from '../../services/accessibility.service';

@Component({
  selector: 'app-blog-post',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatCardModule,
    MatTooltipModule
  ],
  templateUrl: './blog-post.component.html',
  styleUrl: './blog-post.component.scss'
})
export class BlogPostComponent implements OnInit, OnDestroy {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private blogService = inject(BlogService);
  private accessibilityService = inject(AccessibilityService);
  private titleService = inject(Title);
  private metaService = inject(Meta);
  private destroy$ = new Subject<void>();

  // Signals for reactive state management
  blog = signal<Blog | null>(null);
  loading = signal<boolean>(false);
  error = signal<string | null>(null);
  highContrast = signal<boolean>(false);
  fontSize = signal<'small' | 'medium' | 'large' | 'extra-large'>('medium');
  isReading = signal<boolean>(false);
  accessibilityScore = signal<any>(null);

  private speechSynthesis: SpeechSynthesis | null = null;
  private currentUtterance: SpeechSynthesisUtterance | null = null;

  ngOnInit(): void {
    this.initializeComponent();
    this.loadBlogPost();
    this.subscribeToAccessibilitySettings();
    this.subscribeToAccessibilityScore();
  }

  ngOnDestroy(): void {
    this.stopReading();
    this.destroy$.next();
    this.destroy$.complete();
  }

  private initializeComponent(): void {
    // Initialize speech synthesis if available
    if ('speechSynthesis' in window) {
      this.speechSynthesis = window.speechSynthesis;
    }

    // Announce page load to screen readers
    this.accessibilityService.announce('Blog post page loaded');
  }

  private loadBlogPost(): void {
    const slug = this.route.snapshot.paramMap.get('slug');
    if (!slug) {
      this.error.set('Invalid blog post URL');
      return;
    }

    this.loading.set(true);
    this.error.set(null);

    this.blogService.getBlogBySlug(slug).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (blog) => {
        this.blog.set(blog);
        this.loading.set(false);
        this.updateSEOMetadata(blog);
        this.calculateAccessibilityScore();
        
        // Announce article loaded
        this.accessibilityService.announce(`Article loaded: ${blog.title}`);
      },
      error: (error) => {
        this.loading.set(false);
        if (error.status === 404) {
          this.error.set('Blog post not found. It may have been moved or deleted.');
        } else {
          this.error.set('Failed to load blog post. Please check your connection and try again.');
        }
        console.error('Error loading blog post:', error);
      }
    });
  }

  private subscribeToAccessibilitySettings(): void {
    this.accessibilityService.settings$.pipe(
      takeUntil(this.destroy$)
    ).subscribe(settings => {
      this.highContrast.set(settings.highContrast);
      this.fontSize.set(settings.fontSize);
    });
  }

  private subscribeToAccessibilityScore(): void {
    this.accessibilityService.score$.pipe(
      takeUntil(this.destroy$)
    ).subscribe(score => {
      this.accessibilityScore.set(score);
    });
  }

  private updateSEOMetadata(blog: Blog): void {
    // Update page title
    this.titleService.setTitle(`${blog.title} | TechnoPrise Blog`);

    // Update meta tags
    this.metaService.updateTag({ name: 'description', content: blog.meta_description || blog.excerpt });
    this.metaService.updateTag({ name: 'author', content: blog.author });
    this.metaService.updateTag({ name: 'keywords', content: blog.tags.join(', ') });
    
    // Open Graph tags
    this.metaService.updateTag({ property: 'og:title', content: blog.meta_title || blog.title });
    this.metaService.updateTag({ property: 'og:description', content: blog.meta_description || blog.excerpt });
    this.metaService.updateTag({ property: 'og:type', content: 'article' });
    this.metaService.updateTag({ property: 'og:url', content: window.location.href });
    
    // Article specific tags
    this.metaService.updateTag({ property: 'article:author', content: blog.author });
    this.metaService.updateTag({ property: 'article:published_time', content: blog.published_at || blog.created_at });
    this.metaService.updateTag({ property: 'article:modified_time', content: blog.updated_at });
    this.metaService.updateTag({ property: 'article:tag', content: blog.tags.join(', ') });

    // Twitter Card tags
    this.metaService.updateTag({ name: 'twitter:card', content: 'summary_large_image' });
    this.metaService.updateTag({ name: 'twitter:title', content: blog.meta_title || blog.title });
    this.metaService.updateTag({ name: 'twitter:description', content: blog.meta_description || blog.excerpt });
  }

  private calculateAccessibilityScore(): void {
    setTimeout(() => {
      this.accessibilityService.calculateAccessibilityScore();
    }, 1000);
  }

  goBack(): void {
    this.router.navigate(['/']);
    this.accessibilityService.announce('Navigated back to blog home');
  }

  retryLoad(): void {
    this.loadBlogPost();
  }

  toggleHighContrast(): void {
    const newValue = !this.highContrast();
    this.accessibilityService.updateSettings({ highContrast: newValue });
    this.accessibilityService.announce(
      newValue ? 'High contrast mode enabled' : 'High contrast mode disabled'
    );
  }

  increaseFontSize(): void {
    const sizes: Array<'small' | 'medium' | 'large' | 'extra-large'> = ['small', 'medium', 'large', 'extra-large'];
    const currentIndex = sizes.indexOf(this.fontSize());
    if (currentIndex < sizes.length - 1) {
      const newSize = sizes[currentIndex + 1];
      this.accessibilityService.updateSettings({ fontSize: newSize });
      this.accessibilityService.announce(`Font size increased to ${newSize}`);
    }
  }

  decreaseFontSize(): void {
    const sizes: Array<'small' | 'medium' | 'large' | 'extra-large'> = ['small', 'medium', 'large', 'extra-large'];
    const currentIndex = sizes.indexOf(this.fontSize());
    if (currentIndex > 0) {
      const newSize = sizes[currentIndex - 1];
      this.accessibilityService.updateSettings({ fontSize: newSize });
      this.accessibilityService.announce(`Font size decreased to ${newSize}`);
    }
  }

  readAloud(): void {
    if (!this.speechSynthesis || !this.blog()) return;

    if (this.isReading()) {
      this.stopReading();
      return;
    }

    const blog = this.blog()!;
    const textToRead = `${blog.title}. ${blog.excerpt}. ${this.stripHTML(blog.content || '')}`;
    
    this.currentUtterance = new SpeechSynthesisUtterance(textToRead);
    this.currentUtterance.rate = 0.8;
    this.currentUtterance.pitch = 1;
    this.currentUtterance.volume = 1;

    this.currentUtterance.onstart = () => {
      this.isReading.set(true);
      this.accessibilityService.announce('Started reading article aloud');
    };

    this.currentUtterance.onend = () => {
      this.isReading.set(false);
      this.accessibilityService.announce('Finished reading article');
    };

    this.currentUtterance.onerror = () => {
      this.isReading.set(false);
      this.accessibilityService.announce('Error occurred while reading');
    };

    this.speechSynthesis.speak(this.currentUtterance);
  }

  private stopReading(): void {
    if (this.speechSynthesis && this.isReading()) {
      this.speechSynthesis.cancel();
      this.isReading.set(false);
      this.accessibilityService.announce('Stopped reading article');
    }
  }

  shareArticle(): void {
    const blog = this.blog();
    if (!blog) return;

    if (navigator.share) {
      navigator.share({
        title: blog.title,
        text: blog.excerpt,
        url: window.location.href
      }).then(() => {
        this.accessibilityService.announce('Article shared successfully');
      }).catch(() => {
        this.fallbackShare();
      });
    } else {
      this.fallbackShare();
    }
  }

  private fallbackShare(): void {
    navigator.clipboard.writeText(window.location.href).then(() => {
      this.accessibilityService.announce('Article URL copied to clipboard');
    }).catch(() => {
      this.accessibilityService.announce('Unable to share article');
    });
  }

  private stripHTML(html: string): string {
    const tmp = document.createElement('div');
    tmp.innerHTML = html;
    return tmp.textContent || tmp.innerText || '';
  }

  formatDate(dateString: string): string {
    return this.blogService.formatDate(dateString);
  }

  formatReadingTime(minutes: number): string {
    return this.blogService.formatReadingTime(minutes);
  }
}
