import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./components/blog-home/blog-home.component').then(m => m.BlogHomeComponent),
    title: 'TechnoPrise Global Blog - Accessibility-First Insights'
  },
  {
    path: 'blog/:slug',
    loadComponent: () => import('./components/blog-post/blog-post.component').then(m => m.BlogPostComponent),
    title: 'Blog Post - TechnoPrise Global',
    data: { prerender: false }
  },
  {
    path: 'create-blog',
    loadComponent: () => import('./components/create-blog/create-blog').then(m => m.CreateBlogComponent),
    title: 'Create New Blog Post - TechnoPrise Global'
  },
  {
    path: '**',
    redirectTo: ''
  }
];
