import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatDialogModule } from '@angular/material/dialog';
import { COMMA as COMMA_KEY, ENTER as ENTER_KEY } from '@angular/cdk/keycodes';
import { FormBuilder, FormGroup, Validators, AbstractControl, ReactiveFormsModule } from '@angular/forms';
import { BlogService } from '../../services/blog.service';
import { Blog } from '../../services/blog.service';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatChipsModule } from '@angular/material/chips';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { CommonModule } from '@angular/common';
import { AccessibilityService } from '../../services/accessibility.service';
import { Subscription } from 'rxjs';
import { ConfirmDialogComponent } from '../../components/confirm-dialog/confirm-dialog.component';

interface BlogForm {
  title: string;
  author: string;
  excerpt: string;
  content: string;
  tags: string[];
  meta_title: string;
  meta_description: string;
  status: 'draft' | 'published';
  featured: boolean;
  published_at: Date | null;
}

@Component({
  selector: 'app-create-blog',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule,
    MatChipsModule,
    MatCheckboxModule,
    MatSnackBarModule,
    MatDatepickerModule,
    MatNativeDateModule,
    // Added missing Material modules
    MatIconModule,
    MatExpansionModule,
    MatDialogModule
  ],
  templateUrl: './create-blog.html',
  styleUrls: ['./create-blog.scss']
})
export class CreateBlogComponent implements OnInit, OnDestroy {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private snackBar = inject(MatSnackBar);
  private dialog = inject(MatDialog);
  private accessibilityService = inject(AccessibilityService);
  private blogService = inject(BlogService);
  private subscriptions = new Subscription();

  // Tracks submission state
  submittingFlag = false;
  errorMessageText: string | null = null;
  // Keycode for semicolon separator
  readonly SEMI_COLON = 186;
  // Expose keycode constants to template
  readonly ENTER = ENTER_KEY;
  readonly COMMA = COMMA_KEY;
  minPublishDate = new Date();

  // Constants for validation
  private readonly TITLE_MIN_LENGTH = 5;
  private readonly TITLE_MAX_LENGTH = 255;
  private readonly EXCERPT_MIN_LENGTH = 20;
  private readonly EXCERPT_MAX_LENGTH = 500;
  private readonly CONTENT_MIN_LENGTH = 10;
  private readonly META_TITLE_MAX_LENGTH = 70;
  private readonly META_DESC_MAX_LENGTH = 320;
  private readonly MAX_TAGS = 10;
  private readonly TAG_MAX_LENGTH = 50;

  blogForm!: FormGroup;

  ngOnInit(): void {
    this.initializeForm();
    // Announce page load for screen readers
    this.accessibilityService.announce('Create new blog post page loaded');
  }

  private initializeForm(): void {
    this.blogForm = this.fb.group({
      title: ['', [
        Validators.required,
        Validators.minLength(this.TITLE_MIN_LENGTH),
        Validators.maxLength(this.TITLE_MAX_LENGTH)
      ]],
      author: ['', [
        Validators.required,
        Validators.maxLength(100)
      ]],
      excerpt: ['', [
        Validators.minLength(this.EXCERPT_MIN_LENGTH),
        Validators.maxLength(this.EXCERPT_MAX_LENGTH)
      ]],
      content: ['', [
        Validators.required,
        Validators.minLength(this.CONTENT_MIN_LENGTH)
      ]],
      tags: [[], [
        this.validateMaxTags.bind(this)
      ]],
      meta_title: ['', [
        Validators.maxLength(this.META_TITLE_MAX_LENGTH)
      ]],
      meta_description: ['', [
        Validators.maxLength(this.META_DESC_MAX_LENGTH)
      ]],
      status: ['published', Validators.required],
      featured: [false],
      published_at: [null]
    });

    // Announce page load for screen readers
    this.accessibilityService.announce('Create new blog post page loaded');
  }



  // Form validation
  validateMaxTags(control: AbstractControl): { [key: string]: boolean } | null {
    const value = control.value;
    if (Array.isArray(value) && value.length > this.MAX_TAGS) {
      return { 'maxTagsExceeded': true };
    }
    return null;
  }

  // Tag management
  addTag(event: any): void {
    const input = event.chipInput?.inputElement;
    const value = (event.value || '').trim();

    if (value && value.length <= this.TAG_MAX_LENGTH) {
      const tags = [...this.blogForm.get('tags')?.value || []];
      if (!tags.includes(value) && tags.length < this.MAX_TAGS) {
        this.blogForm.get('tags')?.setValue([...tags, value]);
      }
    }

    if (input) {
      input.value = '';
    }
  }

  removeTag(tag: string): void {
    const tags = this.blogForm.get('tags')?.value.filter((t: string) => t !== tag) || [];
    this.blogForm.get('tags')?.setValue(tags);
  }

  // Form submission
  onSubmit(): void {
    if (this.blogForm.invalid) {
      this.markFormGroupTouched(this.blogForm);
      return;
    }

    this.submittingFlag = true;
    this.errorMessageText = null;

    const formValue = this.blogForm.getRawValue();
    
    // Prepare the blog post data
    const blogPost = {
      title: formValue.title,
      content: formValue.content,
      excerpt: formValue.excerpt,
      author: formValue.author,
      published: formValue.status === 'published',
      featured: formValue.featured,
      tags: formValue.tags?.join(','),
      meta_title: formValue.meta_title,
      meta_description: formValue.meta_description,
      published_at: formValue.status === 'published' && !formValue.published_at 
        ? new Date().toISOString() 
        : formValue.published_at
    };

    // Call the blog service to create the post
    const createSub = this.blogService.createBlog(blogPost).subscribe({
      next: (response: Blog) => {
        this.submittingFlag = false;
        
        // Announce success to screen readers
        this.accessibilityService.announce(`Blog post ${response.title} ${response.published ? 'published' : 'saved as draft'} successfully`);
        
        // Show success message
        this.snackBar.open(
          `Blog post ${formValue.status === 'published' ? 'published' : 'saved as draft'} successfully!`,
          'Close',
          {
            duration: 5000,
            panelClass: ['success-snackbar']
          }
        );
        
        // Navigate to the new blog post or blog list
        if (response && response.slug) {
          this.router.navigate(['/blog', response.slug]);
        } else {
          this.router.navigate(['/blog']);
        }
      },
      error: (err: any) => {
        console.error('Error creating blog post:', err);
        this.submittingFlag = false;
        
        // Get error message from response or use default
        let errorMessage = 'An error occurred while saving the blog post. Please try again.';
        if (err.error?.message) {
          errorMessage = err.error.message;
        } else if (err.status === 400) {
          errorMessage = 'Invalid blog post data. Please check your inputs.';
        } else if (err.status === 500) {
          errorMessage = 'Server error occurred. Please try again later.';
        }
        
        this.errorMessageText = errorMessage;
        
        // Show error message
        this.snackBar.open(errorMessage, 'Close', {
          duration: 10000,
          panelClass: ['error-snackbar']
        });
        
        // Announce error to screen readers
        this.accessibilityService.announce(errorMessage);
        
        // Focus on the first invalid control
        const invalidControl = this.getFirstInvalidControl();
        if (invalidControl) {
          invalidControl.focus();
        }
      }
    });

    this.subscriptions.add(createSub);
  }

  // Handle cancel button click
  onCancel(): void {
    if (this.blogForm.dirty) {
      const dialogRef = this.dialog.open(ConfirmDialogComponent, {
        width: '400px',
        data: {
          title: 'Discard Changes',
          message: 'You have unsaved changes. Are you sure you want to discard them?',
          confirmText: 'Discard',
          cancelText: 'Keep Editing',
          confirmColor: 'warn'
        }
      });

      const dialogSub = dialogRef.afterClosed().subscribe(result => {
        if (result) {
          this.router.navigate(['/blog']);
        }
      });

      this.subscriptions.add(dialogSub);
    } else {
      this.router.navigate(['/blog']);
    }
  }

  // Clean up subscriptions
  ngOnDestroy() {
    this.subscriptions.unsubscribe();
  }

  private markFormGroupTouched(group: FormGroup): void {
    Object.values(group.controls).forEach(control => {
      control.markAsTouched();
      if (control instanceof FormGroup) {
        this.markFormGroupTouched(control);
      }
    });
  }

  private getFirstInvalidControl(): HTMLElement | null {
    const formElement = document.querySelector('form');
    if (!formElement) return null;

    const invalidControls = formElement.querySelectorAll('.ng-invalid');
    return invalidControls.length > 0 ? invalidControls[0] as HTMLElement : null;
  }

  private calculateReadingTime(content: string): number {
    const wordsPerMinute = 200;
    const wordCount = content.split(/\s+/).length;
    return Math.ceil(wordCount / wordsPerMinute);
  }

  // Helper for template to indicate submitting state
  isSubmitting(): boolean {
    return this.submittingFlag;
  }

  // Getter for error message used in template
  errorMessage(): string | null {
    return this.errorMessageText;
  }

  getFieldError(fieldName: string): string {
    const field = this.blogForm.get(fieldName);
    if (field?.errors && field.touched) {
      if (field.errors['required']) {
        return `${fieldName.charAt(0).toUpperCase() + fieldName.slice(1)} is required`;
      }
      if (field.errors['minlength']) {
        return `${fieldName.charAt(0).toUpperCase() + fieldName.slice(1)} must be at least ${field.errors['minlength'].requiredLength} characters`;
      }
      if (field.errors['maxlength']) {
        return `${fieldName.charAt(0).toUpperCase() + fieldName.slice(1)} must not exceed ${field.errors['maxlength'].requiredLength} characters`;
      }
    }
    return '';
  }
}
