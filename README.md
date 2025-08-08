# TechnoPrise Global Blog Platform

[![Angular](https://img.shields.io/badge/Angular-20-red.svg)](https://angular.io/)
[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.8-blue.svg)](https://www.typescriptlang.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.4-38B2AC.svg)](https://tailwindcss.com/)
[![Angular Material](https://img.shields.io/badge/Angular_Material-20-FF6D00.svg)](https://material.angular.io/)
[![WCAG 2.2 AA](https://img.shields.io/badge/WCAG-2.2_AA-green.svg)](https://www.w3.org/WAI/WCAG21/quickref/)

A modern, accessibility-first blog platform built with Angular and Golang, featuring a beautiful UI, comprehensive accessibility support, and enterprise-grade code quality.

## Features

### Core Features
- **Blog Management**: Create, read, and manage blog posts with rich content
- **Real-time Search**: Instant search with debouncing and accessibility announcements
- **Pagination**: Smooth pagination with keyboard navigation support
- **Responsive Design**: Mobile-first design that works on all devices
- **SEO Optimized**: Meta tags, Open Graph, and Twitter Card support
- **Accessibility First**: WCAG 2.2 AA compliant with screen reader support

### Advanced Features
- **Modern UI**: Clean, futuristic design with smooth animations
- **Dark Mode Ready**: High contrast and reduced motion support
- **Accessibility Scoring**: Real-time accessibility score calculation
- **Screen Reader Support**: Live regions and ARIA attributes throughout
- **Keyboard Navigation**: Full keyboard accessibility
- **Font Size Controls**: Dynamic font scaling for better readability

### Technical Excellence
- **Angular 20**: Latest Angular with standalone components
- **Golang Backend**: High-performance REST API with Gin framework
- **PostgreSQL/SQLite Database**: PostgreSQL for production, SQLite fallback with GORM ORM
- **Material Design**: Angular Material components with custom theming
- **Tailwind CSS**: Utility-first CSS framework for rapid development
- **TypeScript**: Full type safety and modern JavaScript features

## Prerequisites

Before you begin, ensure you have the following installed on your system:

### Required Software
- **Node.js** (v18.0.0 or higher) - [Download here](https://nodejs.org/)
- **Go** (v1.21 or higher) - [Download here](https://golang.org/dl/)
- **PostgreSQL** (v12 or higher) - [Download here](https://www.postgresql.org/download/) *(Primary database)*
- **Git** - [Download here](https://git-scm.com/)

### Recommended Tools
- **VS Code** with Angular and Go extensions
- **Postman** or similar API testing tool
- **Chrome DevTools** for debugging

## Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/okothmax/technoprise-blog-platform.git
cd technoprise-blog-platform
```

### 2. Database Setup

The application supports both PostgreSQL (primary) and SQLite (fallback):

#### Option A: PostgreSQL Setup (Recommended for Production)
1. Install PostgreSQL and create a database:
```bash
# Create database
createdb technoprise_blog

# Or using psql
psql -U postgres
CREATE DATABASE technoprise_blog;
\q
```

#### Option B: SQLite Setup (Development/Testing)
SQLite requires no additional setup - the database file will be created automatically.

### 3. Backend Setup

Navigate to the backend directory:
```bash
cd backend
```

Install Go dependencies:
```bash
go mod download
```

Create environment file:
```bash
cp .env.example .env
```

Edit `.env` file with your configuration:
```env
# Database Configuration
DB_TYPE=sqlite
DB_NAME=technoprise_blog.db

# Server Configuration
PORT=8080
GIN_MODE=release

# CORS Configuration
FRONTEND_URL=http://localhost:4200
```

Run the backend server:
```bash
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080` with the following endpoints:
- `GET /api/v1/blogs` - Get paginated blog posts
- `GET /api/v1/blogs/:slug` - Get single blog post
- `POST /api/v1/blogs` - Create new blog post
- `GET /api/v1/health` - Health check

### 4. Frontend Setup

Open a new terminal and navigate to the frontend directory:
```bash
cd frontend
```

Install dependencies:
```bash
npm install
```

Start the development server:
```bash
npm start
```

The frontend will start on `http://localhost:4200`

### 5. Access the Application

Open your browser and navigate to:
- **Frontend**: http://localhost:4200
- **Backend API**: http://localhost:8080/api/v1/health

## Project Structure

```
technoprise-blog-platform/
├── backend/                    # Golang backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # Application entry point
│   ├── internal/
│   │   ├── config/            # Configuration management
│   │   ├── database/          # Database setup and migrations
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── middleware/        # HTTP middleware
│   │   ├── models/            # Data models and structs
│   │   └── routes/            # Route definitions
│   ├── pkg/                   # Shared packages
│   ├── .env.example          # Environment variables template
│   ├── go.mod                # Go module dependencies
│   └── go.sum                # Go module checksums
├── frontend/                  # Angular frontend
│   ├── src/
│   │   ├── app/
│   │   │   ├── components/    # Angular components
│   │   │   │   ├── blog-home/ # Homepage component
│   │   │   │   ├── blog-post/ # Blog post detail component
│   │   │   │   └── create-blog/ # Blog creation component
│   │   │   ├── services/      # Angular services
│   │   │   │   ├── accessibility.service.ts
│   │   │   │   └── blog.service.ts
│   │   │   ├── app.routes.ts  # Route configuration
│   │   │   ├── app.component.ts # Root component
│   │   │   └── app.config.ts  # App configuration
│   │   ├── environments/      # Environment configurations
│   │   ├── styles.scss       # Global styles
│   │   └── index.html        # Main HTML file
│   ├── angular.json          # Angular CLI configuration
│   ├── package.json          # Node.js dependencies
│   ├── tailwind.config.js    # Tailwind CSS configuration
│   └── tsconfig.json         # TypeScript configuration
└── README.md                 # This file
```

## Key Components

### Frontend Components

#### Blog Home Component
- Displays paginated list of blog posts
- Real-time search functionality
- Responsive card-based layout
- Loading states and error handling

#### Blog Post Component
- Dynamic routing with slug-based URLs
- Full blog post content display
- SEO metadata management
- Accessibility score calculation
- Social sharing capabilities

#### Create Blog Component
- Rich form with validation
- Tag management system
- SEO settings panel
- Draft/publish functionality
- Accessibility-compliant form fields

### Backend Services

#### Database Layer
- **Primary**: PostgreSQL database with GORM ORM
- **Fallback**: SQLite for development/testing
- Automatic migrations for both databases
- Seeded sample data
- Optimized queries with indexing

#### API Layer
- RESTful API design
- JSON request/response
- CORS configuration
- Error handling middleware

#### Business Logic
- Blog CRUD operations
- Search and filtering
- Slug generation
- Reading time calculation

## Styling and Design

### Design System
- **Colors**: Modern blue and gray palette
- **Typography**: Inter font family for readability
- **Spacing**: Consistent 8px grid system
- **Components**: Material Design with custom theming

### Responsive Breakpoints
```scss
// Mobile First Approach
@media (min-width: 640px)  { /* sm */ }
@media (min-width: 768px)  { /* md */ }
@media (min-width: 1024px) { /* lg */ }
@media (min-width: 1280px) { /* xl */ }
```

### Accessibility Features
- **WCAG 2.2 AA Compliance**: Full accessibility standard compliance
- **Screen Reader Support**: ARIA labels and live regions
- **Keyboard Navigation**: Tab order and focus management
- **High Contrast Mode**: Enhanced visibility options
- **Font Size Controls**: Dynamic text scaling
- **Reduced Motion**: Respects user motion preferences

## Development

### Available Scripts

#### Frontend Scripts
```bash
npm start          # Start development server
npm run build      # Build for production
npm run test       # Run unit tests
npm run lint       # Lint TypeScript code
npm run e2e        # Run end-to-end tests
```

#### Backend Scripts
```bash
go run cmd/server/main.go    # Start development server
go build                     # Build binary
go test ./...               # Run tests
go mod tidy                 # Clean up dependencies
```

### Environment Variables

#### Backend (.env)
```env
# Database Configuration (PostgreSQL - Primary)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=technoprise_blog
DB_SSLMODE=disable

# For SQLite fallback, comment above and use:
# DB_TYPE=sqlite
# DB_NAME=technoprise_blog.db

# Server Configuration
PORT=8080
GIN_MODE=debug

# CORS Configuration
FRONTEND_URL=http://localhost:4200

# Security
JWT_SECRET=your-jwt-secret-key-here
API_KEY=your-api-key-here

# Features
ENABLE_SWAGGER=true
ENABLE_METRICS=true
ENABLE_LOGGING=true
```

#### Frontend (environment.ts)
```typescript
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api/v1'
};
```

## Testing

### Frontend Testing
```bash
# Unit tests with Jasmine/Karma
npm run test

# E2E tests with Protractor
npm run e2e

# Test coverage
npm run test -- --code-coverage
```

### Backend Testing
```bash
# Run all tests
go test ./...

# Test with coverage
go test -cover ./...

# Verbose test output
go test -v ./...
```

### Manual Testing Checklist
- [ ] Homepage loads and displays blog posts
- [ ] Search functionality works correctly
- [ ] Pagination navigates properly
- [ ] Blog post pages load with correct content
- [ ] Create blog form validates and submits
- [ ] Responsive design works on mobile/tablet
- [ ] Accessibility features function properly
- [ ] Dark mode and high contrast work
- [ ] Keyboard navigation is smooth

## Deployment

### Frontend Deployment

#### Build for Production
```bash
cd frontend
npm run build
```

#### Deploy to Netlify
1. Connect your GitHub repository to Netlify
2. Set build command: `npm run build`
3. Set publish directory: `dist/frontend`
4. Add environment variables in Netlify dashboard

#### Deploy to Vercel
```bash
npm install -g vercel
vercel --prod
```

### Backend Deployment

#### Build Binary
```bash
cd backend
go build -o technoprise-blog cmd/server/main.go
```

#### Deploy to Railway
1. Connect GitHub repository
2. Set root directory to `backend`
3. Railway will auto-detect Go application

#### Deploy to Heroku
```bash
# Create Procfile in backend directory
echo "web: ./technoprise-blog" > Procfile

# Deploy
git subtree push --prefix backend heroku main
```

## Security Considerations

### Backend Security
- **CORS Configuration**: Properly configured for frontend domain
- **Input Validation**: All user inputs are validated
- **SQL Injection Prevention**: GORM ORM provides protection
- **Environment Variables**: Sensitive data stored in .env files

### Frontend Security
- **XSS Prevention**: Angular's built-in sanitization
- **CSRF Protection**: Angular's HTTP interceptors
- **Content Security Policy**: Configured in index.html
- **Secure Headers**: Set in production builds

## Performance Optimization

### Frontend Optimizations
- **Lazy Loading**: Routes are lazy-loaded
- **Tree Shaking**: Unused code is eliminated
- **Minification**: Production builds are minified
- **Caching**: Browser caching strategies implemented
- **Bundle Analysis**: Use `npm run build -- --stats-json`

### Backend Optimizations
- **Database Indexing**: Proper indexes on search fields
- **Query Optimization**: Efficient database queries
- **Caching**: In-memory caching for frequent requests
- **Compression**: Gzip compression enabled

## Troubleshooting

### Common Issues

#### Frontend Issues

**Issue**: `npm install` fails
```bash
# Solution: Clear cache and reinstall
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

**Issue**: Angular Material styles not loading
```bash
# Solution: Check if styles are imported in angular.json
"styles": [
  "@angular/material/prebuilt-themes/indigo-pink.css",
  "src/styles.scss"
]
```

**Issue**: Tailwind classes not working
```bash
# Solution: Rebuild with Tailwind
npm run build
```

#### Backend Issues

**Issue**: Database connection fails
```bash
# Solution: Check if .env file exists and has correct values
cp .env.example .env
# Edit .env with correct database path
```

**Issue**: CORS errors in browser
```bash
# Solution: Update FRONTEND_URL in .env
FRONTEND_URL=http://localhost:4200
```

**Issue**: Go modules not found
```bash
# Solution: Download dependencies
go mod download
go mod tidy
```

### Debug Mode

#### Enable Debug Logging
```bash
# Backend
export GIN_MODE=debug

# Frontend
ng serve --configuration=development
```

#### Browser DevTools
- Open Chrome DevTools (F12)
- Check Console for JavaScript errors
- Use Network tab to debug API calls
- Use Lighthouse for performance auditing

## Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `npm test` and `go test ./...`
5. Commit changes: `git commit -m 'Add amazing feature'`
6. Push to branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Style Guidelines

#### TypeScript/Angular
- Use Angular CLI for generating components
- Follow Angular style guide
- Use TypeScript strict mode
- Write unit tests for components and services

#### Go
- Follow Go formatting standards (`go fmt`)
- Use meaningful variable names
- Write unit tests for handlers and models
- Document public functions

### Commit Message Format
```
type(scope): description

feat(blog): add search functionality
fix(api): resolve CORS issue
docs(readme): update installation guide
style(ui): improve button styling
test(blog): add unit tests for blog service
```

## Additional Resources

### Documentation
- [Angular Documentation](https://angular.io/docs)
- [Go Documentation](https://golang.org/doc/)
- [Material Design Guidelines](https://material.io/design)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [WCAG 2.2 Guidelines](https://www.w3.org/WAI/WCAG21/quickref/)

### Learning Resources
- [Angular University](https://angular-university.io/)
- [Go by Example](https://gobyexample.com/)
- [Web Accessibility Guidelines](https://webaim.org/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)

### Tools and Extensions
- **VS Code Extensions**:
  - Angular Language Service
  - Go for Visual Studio Code
  - Tailwind CSS IntelliSense
  - ESLint
  - Prettier

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Authors

- **Max Okoth** - *Initial work* - [okothmax](https://github.com/okothmax)

## Acknowledgments

- Angular team for the amazing framework
- Go team for the excellent language
- Material Design team for the design system
- Tailwind CSS team for the utility framework
- Web accessibility community for the guidelines
