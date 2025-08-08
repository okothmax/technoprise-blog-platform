import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';

export interface Blog {
  id: number;
  title: string;
  slug: string;
  content?: string;
  excerpt: string;
  author: string;
  published: boolean;
  featured: boolean;
  tags: string[];
  meta_title?: string;
  meta_description?: string;
  reading_time: number;
  view_count: number;
  created_at: string;
  updated_at: string;
  published_at?: string;
}

export interface BlogListResponse {
  blogs: Blog[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface BlogFilters {
  page?: number;
  limit?: number;
  search?: string;
  featured?: boolean;
  published?: boolean;
  tags?: string;
  categories?: string;
  startDate?: string;
  endDate?: string;
  sort?: string;
}

@Injectable({
  providedIn: 'root'
})
export class BlogService {
  private http = inject(HttpClient);
  private readonly apiUrl = `${environment.apiUrl}/blogs`;
  
  // State management for accessibility and search
  private loadingSubject = new BehaviorSubject<boolean>(false);
  private errorSubject = new BehaviorSubject<string | null>(null);
  private searchQuerySubject = new BehaviorSubject<string>('');
  
  public loading$ = this.loadingSubject.asObservable();
  public error$ = this.errorSubject.asObservable();
  public searchQuery$ = this.searchQuerySubject.asObservable();

  /**
   * Get paginated list of blog posts with accessibility features
   */
  getBlogs(filters: BlogFilters = {}): Observable<BlogListResponse> {
    this.setLoading(true);
    this.clearError();

    let params = new HttpParams();
    
    if (filters.page) params = params.set('page', filters.page.toString());
    if (filters.limit) params = params.set('limit', filters.limit.toString());
    if (filters.search) params = params.set('search', filters.search);
    if (filters.featured !== undefined) params = params.set('featured', filters.featured.toString());
    if (filters.published !== undefined) params = params.set('published', filters.published.toString());
    if (filters.tags) params = params.set('tags', filters.tags);
    if (filters.categories) params = params.set('categories', filters.categories);
    if (filters.startDate) params = params.set('startDate', filters.startDate);
    if (filters.endDate) params = params.set('endDate', filters.endDate);
    if (filters.sort) params = params.set('sort', filters.sort);

    return this.http.get<BlogListResponse>(this.apiUrl, { params }).pipe(
      map(response => {
        this.setLoading(false);
        // Announce to screen readers
        this.announceToScreenReader(`Loaded ${response.blogs.length} blog posts`);
        return response;
      }),
      catchError(error => {
        this.setLoading(false);
        this.setError('Failed to load blog posts. Please try again.');
        this.announceToScreenReader('Error loading blog posts');
        throw error;
      })
    );
  }

  /**
   * Get a single blog post by slug
   */
  getBlogBySlug(slug: string): Observable<Blog> {
    this.setLoading(true);
    this.clearError();

    return this.http.get<Blog>(`${this.apiUrl}/${slug}`).pipe(
      map(blog => {
        this.setLoading(false);
        // Announce to screen readers
        this.announceToScreenReader(`Loaded article: ${blog.title}`);
        return blog;
      }),
      catchError(error => {
        this.setLoading(false);
        if (error.status === 404) {
          this.setError('Blog post not found.');
          this.announceToScreenReader('Blog post not found');
        } else {
          this.setError('Failed to load blog post. Please try again.');
          this.announceToScreenReader('Error loading blog post');
        }
        throw error;
      })
    );
  }

  /**
   * Search blogs with debouncing for accessibility
   */
  searchBlogs(query: string, page: number = 1, limit: number = 10): Observable<BlogListResponse> {
    return this.getBlogs({ search: query, page, limit });
  }

  /**
   * Get featured blogs for homepage
   */
  getFeaturedBlogs(limit: number = 3): Observable<Blog[]> {
    return this.getBlogs({ featured: true, limit }).pipe(
      map(response => response.blogs)
    );
  }

  /**
   * Get recent blogs
   */
  getRecentBlogs(limit: number = 5): Observable<Blog[]> {
    return this.getBlogs({ limit }).pipe(
      map(response => response.blogs)
    );
  }

  // Accessibility helper methods
  private setLoading(loading: boolean): void {
    this.loadingSubject.next(loading);
  }

  private setError(error: string | null): void {
    this.errorSubject.next(error);
  }

  private clearError(): void {
    this.errorSubject.next(null);
  }

  /**
   * Update the search query and notify subscribers
   */
  setSearchQuery(query: string): void {
    this.searchQuerySubject.next(query);
  }

  /**
   * Get the current search query
   */
  getCurrentSearchQuery(): string {
    return this.searchQuerySubject.value;
  }

  /**
   * Announce messages to screen readers
   */
  private announceToScreenReader(message: string): void {
    const announcement = document.createElement('div');
    announcement.setAttribute('aria-live', 'polite');
    announcement.setAttribute('aria-atomic', 'true');
    announcement.className = 'sr-only';
    announcement.textContent = message;
    
    document.body.appendChild(announcement);
    
    // Remove after announcement
    setTimeout(() => {
      document.body.removeChild(announcement);
    }, 1000);
  }

  /**
   * Format reading time for accessibility
   */
  formatReadingTime(minutes: number): string {
    if (minutes < 1) return 'Less than 1 minute read';
    if (minutes === 1) return '1 minute read';
    return `${minutes} minutes read`;
  }

  /**
   * Format date for accessibility
   */
  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  /**
   * Create a new blog post
   * @param blogData The blog post data to create
   * @returns An Observable of the created blog post
   */
  createBlog(blogData: {
    title: string;
    content: string;
    excerpt?: string;
    author: string;
    published: boolean;
    featured: boolean;
    tags: string;
    meta_title?: string;
    meta_description?: string;
  }): Observable<Blog> {
    this.setLoading(true);
    this.clearError();

    // Ensure required fields are provided
    if (!blogData.title || !blogData.content || !blogData.author) {
      const error = 'Title, content, and author are required';
      this.setError(error);
      this.announceToScreenReader('Error: ' + error);
      return new Observable(subscriber => {
        subscriber.error(new Error(error));
      });
    }

    // Prepare the request body
    const requestBody = {
      title: blogData.title,
      content: blogData.content,
      excerpt: blogData.excerpt || '',
      author: blogData.author,
      published: blogData.published || false,
      featured: blogData.featured || false,
      tags: blogData.tags || '',
      meta_title: blogData.meta_title || '',
      meta_description: blogData.meta_description || ''
    };

    return this.http.post<Blog>(this.apiUrl, requestBody).pipe(
      map(response => {
        this.setLoading(false);
        this.announceToScreenReader('Blog post created successfully');
        return response;
      }),
      catchError(error => {
        this.setLoading(false);
        let errorMessage = 'Failed to create blog post. Please try again.';
        
        if (error.error && error.error.error) {
          errorMessage = error.error.error;
        } else if (error.status === 400) {
          errorMessage = 'Invalid data. Please check your input and try again.';
        }
        
        this.setError(errorMessage);
        this.announceToScreenReader('Error: ' + errorMessage);
        throw error;
      })
    );
  }

  /**
   * Generate accessible excerpt
   */
  generateAccessibleExcerpt(content: string, maxLength: number = 150): string {
    if (!content) return '';
    
    // Strip HTML tags
    const stripped = content.replace(/<[^>]*>/g, ' ');
    
    if (stripped.length <= maxLength) return stripped;
    
    // Find last complete word
    const truncated = stripped.substring(0, maxLength);
    const lastSpace = truncated.lastIndexOf(' ');
    
    return lastSpace > 0 ? truncated.substring(0, lastSpace) + '...' : truncated + '...';
  }
}
