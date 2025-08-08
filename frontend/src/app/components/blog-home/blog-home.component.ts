import { Component, OnInit, OnDestroy, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { Subject, takeUntil, debounceTime, distinctUntilChanged } from 'rxjs';

import { BlogService, Blog, BlogListResponse } from '../../services/blog.service';
import { AccessibilityService } from '../../services/accessibility.service';

@Component({
  selector: 'app-blog-home',
  standalone: true,

  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatPaginatorModule
  ],
  templateUrl: './blog-home.component.html',
  styleUrl: './blog-home.component.scss'
})
export class BlogHomeComponent implements OnInit, OnDestroy {
  private blogService = inject(BlogService);
  private accessibilityService = inject(AccessibilityService);
  private destroy$ = new Subject<void>();

  // Signals for reactive state management
  blogResponse = signal<BlogListResponse | null>(null);
  featuredBlogs = signal<Blog[]>([]);
  loading = signal<boolean>(false);
  error = signal<string | null>(null);
  searchQuery = signal<string>('');
  currentPage = signal<number>(1);
  pageSize = signal<number>(10);
  accessibilityScore = signal<any>(null);

  private searchSubject = new Subject<string>();

  ngOnInit(): void {
    this.initializeComponent();
    this.loadInitialData();
    this.subscribeToAccessibilityScore();
    this.subscribeToSearchQuery();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private initializeComponent(): void {
    // Announce page load to screen readers
    this.accessibilityService.announce('Blog home page loaded');
    
    // Calculate accessibility score
    setTimeout(() => {
      this.accessibilityService.calculateAccessibilityScore();
    }, 1000);
  }



  private loadInitialData(): void {
    this.loadBlogs();
    this.loadFeaturedBlogs();
  }

  private subscribeToAccessibilityScore(): void {
    this.accessibilityService.score$.pipe(
      takeUntil(this.destroy$)
    ).subscribe(score => {
      this.accessibilityScore.set(score);
    });
  }

  private subscribeToSearchQuery(): void {
    this.blogService.searchQuery$.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      takeUntil(this.destroy$)
    ).subscribe(query => {
      this.searchQuery.set(query);
      this.currentPage.set(1);
      this.loadBlogs();
    });
  }







  private loadBlogs(): void {
    this.loading.set(true);
    
    const filters = {
      page: this.currentPage(),
      limit: this.pageSize(),
      search: this.searchQuery() || undefined,
      published: true
    };

    this.blogService.getBlogs(filters).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (response) => {
        this.blogResponse.set(response);
        this.loading.set(false);
      },
      error: (error) => {
        console.error('Error loading blogs:', error);
        this.error.set('Failed to load blog posts. Please try again.');
        this.loading.set(false);
      }
    });
  }



  private loadFeaturedBlogs(): void {
    this.blogService.getFeaturedBlogs(3).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (blogs) => {
        this.featuredBlogs.set(blogs);
      },
      error: (error) => {
        console.error('Error loading featured blogs:', error);
      }
    });
  }

  onSearchInput(event: Event): void {
    const target = event.target as HTMLInputElement;
    if (target) {
      const query = target.value.trim();
      this.searchQuery.set(query);
      this.searchSubject.next(query);
    }
  }

  onPageChange(event: PageEvent): void {
    this.currentPage.set(event.pageIndex + 1);
    this.pageSize.set(event.pageSize);
    this.loadBlogs();
    
    // Scroll to top of results
    document.getElementById('posts-title')?.scrollIntoView({ 
      behavior: 'smooth', 
      block: 'start' 
    });
    
    // Announce page change
    this.accessibilityService.announce(
      `Navigated to page ${this.currentPage()} of ${Math.ceil((this.blogResponse()?.total || 0) / this.pageSize())}`
    );
  }

  retryLoad(): void {
    this.loadBlogs();
  }

  trackByBlogId(index: number, blog: Blog): number {
    return blog.id;
  }

  formatDate(dateString: string): string {
    return this.blogService.formatDate(dateString);
  }

  formatReadingTime(minutes: number): string {
    return this.blogService.formatReadingTime(minutes);
  }

  searchResultsAnnouncement(): string {
    const response = this.blogResponse();
    if (!response) return '';
    
    if (this.searchQuery()) {
      return `Found ${response.total} articles matching "${this.searchQuery()}"`;
    }
    
    return `Showing ${response.blogs.length} of ${response.total} articles`;
  }
}

