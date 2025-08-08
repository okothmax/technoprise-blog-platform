package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"technoprise-blog-backend/internal/models"
)

// Initialize sets up the database connection and runs migrations
func Initialize() (*gorm.DB, error) {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "technoprise_blog")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Construct database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Connect to database
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		// Fallback to SQLite for development
		log.Println("PostgreSQL connection failed, falling back to SQLite...")
		db, err = gorm.Open("sqlite3", "technoprise_blog.db")
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %v", err)
		}
	}

	// Configure connection pool
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Enable logging in development
	if os.Getenv("GIN_MODE") != "release" {
		db.LogMode(true)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	// Seed database with sample data
	if err := seedDatabase(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	log.Println("✅ Database initialized successfully")
	return db, nil
}

// runMigrations creates or updates database tables
func runMigrations(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")
	
	// Auto-migrate models
	if err := db.AutoMigrate(&models.Blog{}).Error; err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}

// seedDatabase populates the database with sample blog posts
func seedDatabase(db *gorm.DB) error {
	// Check if blogs already exist
	var count int64
	if err := db.Model(&models.Blog{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("📚 Database already contains blog posts, skipping seed")
		return nil
	}

	log.Println("🌱 Seeding database with sample blog posts...")

	sampleBlogs := []models.Blog{
		{
			Title:       "The Future of Web Accessibility: AI-Powered Inclusive Design",
			Slug:        "future-web-accessibility-ai-inclusive-design",
			Content:     `<h2>Introduction</h2><p>As we advance into the digital age, web accessibility has become more crucial than ever. At TechnoPrise Global, we believe that the future of web accessibility lies in AI-powered inclusive design that automatically adapts to users' needs.</p><h2>AI-Driven Accessibility Features</h2><p>Modern AI technologies are revolutionizing how we approach accessibility:</p><ul><li><strong>Automatic Alt Text Generation:</strong> AI can analyze images and generate descriptive alt text for screen readers.</li><li><strong>Real-time Caption Generation:</strong> Speech-to-text AI provides instant captions for video content.</li><li><strong>Adaptive UI:</strong> Interfaces that automatically adjust based on user preferences and disabilities.</li><li><strong>Voice Navigation:</strong> Natural language processing enables hands-free website navigation.</li></ul><h2>Implementation Best Practices</h2><p>When implementing AI-powered accessibility features, consider:</p><ol><li>User privacy and data protection</li><li>Fallback mechanisms for AI failures</li><li>Continuous learning and improvement</li><li>User control and customization options</li></ol><p>The future is bright for inclusive web experiences that truly serve everyone.</p>`,
			Excerpt:     "Explore how AI-powered technologies are revolutionizing web accessibility and creating truly inclusive digital experiences for all users.",
			Author:      "Dr. Sarah Chen",
			Published:   true,
			Featured:    true,
			Tags:        "accessibility, AI, inclusive design, web development, WCAG",
			MetaTitle:   "AI-Powered Web Accessibility: The Future of Inclusive Design",
			MetaDesc:    "Discover how artificial intelligence is transforming web accessibility with automatic alt text, real-time captions, and adaptive interfaces.",
		},
		{
			Title:       "Mobile-First Accessibility: Designing for Touch and Voice",
			Slug:        "mobile-first-accessibility-touch-voice-design",
			Content:     `<h2>The Mobile Accessibility Revolution</h2><p>With over 6.8 billion smartphone users worldwide, mobile accessibility has become paramount. TechnoPrise Global leads the charge in creating mobile experiences that work for everyone, regardless of ability.</p><h2>Touch Accessibility Fundamentals</h2><p>Mobile touch interfaces present unique challenges and opportunities:</p><ul><li><strong>Target Size:</strong> Minimum 44px touch targets for easy interaction</li><li><strong>Gesture Alternatives:</strong> Provide button alternatives for complex gestures</li><li><strong>Haptic Feedback:</strong> Use vibration to confirm actions for users with visual impairments</li><li><strong>One-Handed Operation:</strong> Design for thumb-friendly navigation</li></ul><h2>Voice Interface Integration</h2><p>Voice commands are transforming mobile accessibility:</p><pre><code>// Voice command integration
class VoiceAccessibility {
  initializeVoiceCommands() {
    this.recognition = new SpeechRecognition();
    this.setupCommands([
      'navigate to home',
      'read article',
      'increase text size',
      'activate high contrast'
    ]);
  }
}</code></pre><h2>Testing Mobile Accessibility</h2><p>Essential testing strategies include:</p><ol><li>Screen reader testing on iOS VoiceOver and Android TalkBack</li><li>Switch control navigation testing</li><li>Voice control testing</li><li>High contrast and zoom testing</li><li>One-handed operation validation</li></ol><p>Mobile accessibility isn't just compliance—it's about creating delightful experiences for all users.</p>`,
			Excerpt:     "Master mobile accessibility design with touch-friendly interfaces, voice commands, and comprehensive testing strategies for inclusive mobile experiences.",
			Author:      "Alex Rivera",
			Published:   true,
			Featured:    true,
			Tags:        "mobile accessibility, touch interfaces, voice commands, responsive design, mobile UX",
			MetaTitle:   "Mobile-First Accessibility: Touch and Voice Design Guide",
			MetaDesc:    "Learn to design accessible mobile interfaces with touch-friendly controls, voice commands, and comprehensive mobile accessibility testing.",
		},
		{
			Title:       "Dark Mode Accessibility: Beyond Just Inverting Colors",
			Slug:        "dark-mode-accessibility-design-principles",
			Content:     `<h2>The Science Behind Dark Mode</h2><p>Dark mode isn't just a trendy design choice—it's a crucial accessibility feature that can reduce eye strain, improve battery life, and enhance usability for users with certain visual conditions.</p><h2>Accessibility Benefits of Dark Mode</h2><ul><li><strong>Reduced Eye Strain:</strong> Lower blue light emission in low-light environments</li><li><strong>Better for Light Sensitivity:</strong> Helps users with photophobia and migraines</li><li><strong>Improved Focus:</strong> Reduces visual distractions for users with ADHD</li><li><strong>Battery Conservation:</strong> Extends device battery life on OLED screens</li></ul><h2>Design Principles for Accessible Dark Mode</h2><h3>Color Contrast Considerations</h3><p>Maintaining WCAG contrast ratios in dark mode requires careful planning:</p><pre><code>/* Dark mode color palette */
:root[data-theme='dark'] {
  --bg-primary: #0d1117;
  --bg-secondary: #161b22;
  --text-primary: #f0f6fc;
  --text-secondary: #8b949e;
  --accent: #58a6ff;
  --border: #30363d;
}</code></pre><h3>Avoiding Pure Black and White</h3><p>Pure black (#000000) can cause halation effects. Instead, use:</p><ul><li>Dark grays (#0d1117) for backgrounds</li><li>Off-whites (#f0f6fc) for text</li><li>Sufficient color differentiation for interactive elements</li></ul><h2>Implementation Best Practices</h2><ol><li><strong>System Preference Detection:</strong> Respect user's OS-level dark mode preference</li><li><strong>Manual Toggle:</strong> Provide user control with persistent settings</li><li><strong>Gradual Transitions:</strong> Smooth animations between light and dark modes</li><li><strong>Image Adaptations:</strong> Adjust images and illustrations for dark backgrounds</li></ol><h2>Testing Dark Mode Accessibility</h2><p>Comprehensive testing should include:</p><ul><li>Contrast ratio validation in both modes</li><li>Screen reader testing in dark mode</li><li>Color blindness simulation</li><li>Low vision user testing</li><li>Performance impact assessment</li></ul><p>Dark mode accessibility is about creating comfortable, inclusive experiences that adapt to users' needs and preferences.</p>`,
			Excerpt:     "Discover the accessibility benefits of dark mode design and learn how to implement inclusive dark themes that go beyond simple color inversion.",
			Author:      "Jordan Kim",
			Published:   true,
			Featured:    false,
			Tags:        "dark mode, accessibility, color contrast, visual design, UX, eye strain",
			MetaTitle:   "Dark Mode Accessibility: Inclusive Design Beyond Color Inversion",
			MetaDesc:    "Learn to design accessible dark mode interfaces with proper contrast ratios, user preferences, and comprehensive accessibility testing.",
		},
		{
			Title:       "WCAG 2.2: What's New and How to Implement the Latest Guidelines",
			Slug:        "wcag-2-2-new-guidelines-implementation",
			Content:     `<h2>Overview of WCAG 2.2</h2><p>The Web Content Accessibility Guidelines (WCAG) 2.2 introduces several new success criteria that further enhance web accessibility. This update focuses on improving the experience for users with cognitive disabilities and mobile device users.</p><h2>New Success Criteria in WCAG 2.2</h2><h3>2.4.11 Focus Not Obscured (Minimum) - AA</h3><p>When a user interface component receives keyboard focus, the component is not entirely hidden due to author-created content.</p><h3>2.4.12 Focus Not Obscured (Enhanced) - AAA</h3><p>When a user interface component receives keyboard focus, no part of the component is hidden by author-created content.</p><h3>2.5.7 Dragging Movements - AA</h3><p>All functionality that uses a dragging movement can be achieved by a single pointer without dragging.</p><h3>2.5.8 Target Size (Minimum) - AA</h3><p>The size of the target for pointer inputs is at least 24 by 24 CSS pixels.</p><h3>3.2.6 Consistent Help - A</h3><p>If a web page contains help mechanisms, they are provided in a consistent order relative to other page content.</p><h3>3.3.7 Redundant Entry - A</h3><p>Information previously entered by or provided to the user that is required to be entered again is either auto-populated or available for selection.</p><h3>3.3.8 Accessible Authentication (Minimum) - AA</h3><p>A cognitive function test is not required for any step in an authentication process.</p><h3>3.3.9 Accessible Authentication (Enhanced) - AAA</h3><p>A cognitive function test or a test that requires the user to remember or transcribe information is not required.</p><h2>Implementation Strategies</h2><p>To successfully implement WCAG 2.2:</p><ul><li>Audit your current accessibility compliance</li><li>Prioritize the new AA-level criteria</li><li>Update your design system and components</li><li>Train your development team</li><li>Implement automated testing</li><li>Conduct user testing with people with disabilities</li></ul><p>Remember, accessibility is not a one-time task but an ongoing commitment to inclusive design.</p>`,
			Excerpt:     "Learn about the new success criteria in WCAG 2.2 and discover practical strategies for implementing these latest accessibility guidelines.",
			Author:      "Michael Rodriguez",
			Published:   true,
			Featured:    false,
			Tags:        "WCAG, accessibility guidelines, compliance, web standards",
			MetaTitle:   "WCAG 2.2 Implementation Guide: New Guidelines Explained",
			MetaDesc:    "Complete guide to WCAG 2.2 new success criteria including focus management, dragging movements, and accessible authentication.",
		},
		{
			Title:       "Building Accessible React Components: A Developer's Guide",
			Slug:        "building-accessible-react-components-guide",
			Content:     `<h2>Why Accessibility Matters in React Development</h2><p>React's component-based architecture provides an excellent foundation for building accessible web applications. However, developers must be intentional about implementing accessibility features from the ground up.</p><h2>Essential Accessibility Patterns</h2><h3>1. Semantic HTML Foundation</h3><pre><code>// Good: Using semantic HTML
function BlogPost({ title, content, author }) {
  return (
    &lt;article&gt;
      &lt;header&gt;
        &lt;h1&gt;{title}&lt;/h1&gt;
        &lt;p&gt;By &lt;span className="author"&gt;{author}&lt;/span&gt;&lt;/p&gt;
      &lt;/header&gt;
      &lt;main dangerouslySetInnerHTML={{ __html: content }} /&gt;
    &lt;/article&gt;
  );
}</code></pre><h3>2. ARIA Labels and Roles</h3><pre><code>// Accessible button component
function AccessibleButton({ children, onClick, disabled, ariaLabel }) {
  return (
    &lt;button
      onClick={onClick}
      disabled={disabled}
      aria-label={ariaLabel}
      className="btn"
    &gt;
      {children}
    &lt;/button&gt;
  );
}</code></pre><h3>3. Focus Management</h3><pre><code>import { useRef, useEffect } from 'react';

function Modal({ isOpen, onClose, children }) {
  const modalRef = useRef(null);
  const previousFocus = useRef(null);

  useEffect(() => {
    if (isOpen) {
      previousFocus.current = document.activeElement;
      modalRef.current?.focus();
    } else {
      previousFocus.current?.focus();
    }
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    &lt;div
      ref={modalRef}
      role="dialog"
      aria-modal="true"
      tabIndex={-1}
      className="modal"
    &gt;
      {children}
    &lt;/div&gt;
  );
}</code></pre><h2>Testing Your Components</h2><p>Use these tools to ensure your React components are accessible:</p><ul><li><strong>eslint-plugin-jsx-a11y:</strong> Catch accessibility issues during development</li><li><strong>@testing-library/jest-dom:</strong> Test accessibility in your unit tests</li><li><strong>axe-core:</strong> Automated accessibility testing</li><li><strong>Screen readers:</strong> Manual testing with actual assistive technology</li></ul><h2>Best Practices Summary</h2><ol><li>Start with semantic HTML</li><li>Use ARIA attributes appropriately</li><li>Manage focus properly</li><li>Ensure keyboard navigation works</li><li>Test with real users and assistive technology</li><li>Maintain color contrast ratios</li><li>Provide alternative text for images</li></ol><p>Building accessible React components requires attention to detail, but the result is a more inclusive web for everyone.</p>`,
			Excerpt:     "Master the art of building accessible React components with practical examples, best practices, and testing strategies for inclusive web development.",
			Author:      "Emma Thompson",
			Published:   true,
			Featured:    false,
			Tags:        "React, accessibility, components, development, ARIA, semantic HTML",
			MetaTitle:   "Accessible React Components: Developer's Complete Guide",
			MetaDesc:    "Learn to build accessible React components with semantic HTML, ARIA labels, focus management, and comprehensive testing strategies.",
		},
		{
			Title:       "The Business Case for Digital Accessibility: ROI and Beyond",
			Slug:        "business-case-digital-accessibility-roi",
			Content:     `<h2>Understanding the Business Impact</h2><p>Digital accessibility isn't just about compliance—it's a strategic business decision that drives growth, innovation, and market expansion. Organizations that prioritize accessibility see measurable returns on their investment.</p><h2>Financial Benefits</h2><h3>Market Expansion</h3><p>The global disability market represents over 1.3 billion people with a combined spending power of $13 trillion annually. By making your digital products accessible, you tap into this significant market segment.</p><h3>Cost Savings</h3><ul><li><strong>Reduced Legal Risk:</strong> Proactive accessibility reduces the risk of costly lawsuits</li><li><strong>Lower Maintenance Costs:</strong> Accessible code is typically cleaner and more maintainable</li><li><strong>Decreased Support Tickets:</strong> Better usability reduces customer support burden</li></ul><h3>Revenue Growth</h3><p>Studies show that accessible websites experience:</p><ul><li>Increased conversion rates (up to 23%)</li><li>Higher customer satisfaction scores</li><li>Better search engine rankings</li><li>Improved mobile experience</li></ul><h2>Brand and Reputation Benefits</h2><h3>Corporate Social Responsibility</h3><p>Accessibility demonstrates your commitment to inclusion and social responsibility, enhancing brand reputation and employee satisfaction.</p><h3>Innovation Driver</h3><p>Designing for accessibility often leads to innovative solutions that benefit all users. Features like voice controls, captions, and simplified interfaces improve the experience for everyone.</p><h2>Measuring ROI</h2><h3>Key Performance Indicators</h3><ul><li>Website traffic and user engagement</li><li>Conversion rates and sales</li><li>Customer satisfaction scores</li><li>Support ticket volume</li><li>Legal compliance costs</li><li>Employee retention and satisfaction</li></ul><h3>Calculating Return on Investment</h3><p>To calculate accessibility ROI:</p><ol><li>Baseline current performance metrics</li><li>Implement accessibility improvements</li><li>Measure changes in key metrics</li><li>Calculate cost savings and revenue increases</li><li>Factor in risk mitigation value</li></ol><h2>Implementation Strategy</h2><h3>Phase 1: Foundation (Months 1-3)</h3><ul><li>Accessibility audit and assessment</li><li>Team training and awareness</li><li>Policy development</li><li>Quick wins implementation</li></ul><h3>Phase 2: Integration (Months 4-9)</h3><ul><li>Design system updates</li><li>Development process integration</li><li>Testing automation</li><li>Content strategy alignment</li></ul><h3>Phase 3: Optimization (Months 10-12)</h3><ul><li>Advanced features implementation</li><li>User feedback integration</li><li>Performance monitoring</li><li>Continuous improvement</li></ul><h2>Success Stories</h2><p>Companies like Microsoft, Apple, and Target have demonstrated that accessibility investments yield significant returns through increased market share, improved customer loyalty, and reduced operational costs.</p><p>The business case for digital accessibility is clear: it's not just the right thing to do—it's the smart thing to do.</p>`,
			Excerpt:     "Discover the compelling business case for digital accessibility, including ROI calculations, market opportunities, and strategic implementation approaches.",
			Author:      "David Park",
			Published:   true,
			Featured:    true,
			Tags:        "business case, ROI, digital accessibility, market expansion, compliance",
			MetaTitle:   "Digital Accessibility ROI: The Complete Business Case",
			MetaDesc:    "Learn how digital accessibility drives business growth with market expansion, cost savings, and measurable ROI across industries.",
		},
		{
			Title:       "Voice User Interfaces: Designing for the Next Generation of Accessibility",
			Slug:        "voice-user-interfaces-accessibility-design",
			Content:     `<h2>The Rise of Voice Technology</h2><p>Voice User Interfaces (VUIs) are transforming how we interact with technology, offering unprecedented opportunities for accessible design. From smart speakers to voice-controlled web applications, VUIs are breaking down barriers for users with various disabilities.</p><h2>Accessibility Benefits of VUIs</h2><h3>Motor Disabilities</h3><p>Voice interfaces eliminate the need for precise motor control, enabling users with limited mobility to navigate and interact with digital content effortlessly.</p><h3>Visual Impairments</h3><p>VUIs provide an alternative to visual interfaces, allowing users with blindness or low vision to access information and complete tasks through audio feedback.</p><h3>Cognitive Disabilities</h3><p>Natural language processing makes interfaces more intuitive, reducing cognitive load and making technology more accessible to users with learning disabilities.</p><h2>Design Principles for Accessible VUIs</h2><h3>1. Clear and Consistent Commands</h3><pre><code>// Good: Clear, consistent voice commands
const voiceCommands = {
  navigation: [
    "Go to home page",
    "Navigate to blog",
    "Open contact page"
  ],
  actions: [
    "Read article",
    "Play next",
    "Save bookmark"
  ]
};</code></pre><h3>2. Contextual Help and Guidance</h3><p>Provide users with discoverable help commands and contextual guidance:</p><ul><li>"What can I say?" - Lists available commands</li><li>"Help with navigation" - Context-specific assistance</li><li>"Repeat that" - Replay last response</li></ul><h3>3. Error Handling and Recovery</h3><pre><code>// Graceful error handling
function handleVoiceError(error) {
  switch(error.type) {
    case 'no-speech':
      return "I didn't hear anything. Please try again.";
    case 'audio-capture':
      return "Microphone not available. Please check your settings.";
    case 'not-allowed':
      return "Microphone access denied. Please enable microphone permissions.";
    default:
      return "I didn't understand. Could you please rephrase that?";
  }
}</code></pre><h2>Implementation Best Practices</h2><h3>Progressive Enhancement</h3><p>Implement VUI as an enhancement to existing interfaces, not a replacement:</p><pre><code>// Progressive VUI enhancement
class AccessibleInterface {
  constructor() {
    this.hasVoiceSupport = 'webkitSpeechRecognition' in window;
    this.initializeVoice();
  }

  initializeVoice() {
    if (this.hasVoiceSupport) {
      this.setupSpeechRecognition();
      this.addVoiceControls();
    }
  }

  setupSpeechRecognition() {
    this.recognition = new webkitSpeechRecognition();
    this.recognition.continuous = false;
    this.recognition.interimResults = false;
    this.recognition.lang = 'en-US';
    
    this.recognition.onresult = (event) => {
      const command = event.results[0][0].transcript;
      this.processVoiceCommand(command);
    };
  }
}</code></pre><h3>Multimodal Feedback</h3><p>Combine voice with visual and haptic feedback for comprehensive accessibility:</p><ul><li>Visual confirmation of voice commands</li><li>Audio feedback for successful actions</li><li>Haptic feedback on mobile devices</li><li>Text alternatives for all voice content</li></ul><h2>Testing VUI Accessibility</h2><h3>User Testing</h3><ul><li>Test with users who have different types of disabilities</li><li>Evaluate in various environments (quiet, noisy, etc.)</li><li>Test with different accents and speech patterns</li><li>Assess cognitive load and learning curve</li></ul><h3>Technical Testing</h3><ul><li>Speech recognition accuracy</li><li>Response time and latency</li><li>Error handling effectiveness</li><li>Fallback mechanism reliability</li></ul><h2>Future Considerations</h2><h3>Privacy and Security</h3><p>Voice interfaces raise important privacy concerns:</p><ul><li>Local processing vs. cloud-based recognition</li><li>Data retention and deletion policies</li><li>User consent and control</li><li>Secure transmission of voice data</li></ul><h3>Emerging Technologies</h3><ul><li>AI-powered natural language understanding</li><li>Emotion recognition in voice</li><li>Multilingual voice interfaces</li><li>Integration with IoT and smart environments</li></ul><p>Voice User Interfaces represent a significant step forward in accessible design, offering new ways to interact with technology that can benefit users of all abilities.</p>`,
			Excerpt:     "Explore how Voice User Interfaces are revolutionizing accessibility, with design principles, implementation strategies, and best practices for inclusive VUI development.",
			Author:      "Lisa Wang",
			Published:   false,
			Featured:    false,
			Tags:        "VUI, voice interfaces, accessibility, design, speech recognition",
			MetaTitle:   "Voice UI Accessibility: Designing Inclusive Voice Interfaces",
			MetaDesc:    "Learn to design accessible Voice User Interfaces with best practices for inclusive VUI development and implementation strategies.",
		},
		{
			Title:       "Cognitive Accessibility: Designing for Neurodiversity",
			Slug:        "cognitive-accessibility-neurodiversity-design",
			Content:     `<h2>Understanding Cognitive Accessibility</h2><p>Cognitive accessibility focuses on making digital experiences usable for people with cognitive and neurological differences, including ADHD, autism, dyslexia, and learning disabilities. At TechnoPrise Global, we champion neurodiversity in design.</p><h2>Key Principles for Cognitive Accessibility</h2><ul><li><strong>Clear Navigation:</strong> Consistent, predictable navigation patterns</li><li><strong>Simple Language:</strong> Plain language and clear instructions</li><li><strong>Reduced Cognitive Load:</strong> Minimize distractions and unnecessary complexity</li><li><strong>Flexible Timing:</strong> Allow users to control time limits and auto-playing content</li><li><strong>Error Prevention:</strong> Clear error messages and recovery paths</li></ul><h2>Design Patterns for Neurodiversity</h2><h3>1. Progressive Disclosure</h3><p>Present information in digestible chunks:</p><pre><code>// Progressive disclosure component
function StepByStep({ steps, currentStep }) {
  return (
    <div className="step-container">
      <ProgressIndicator current={currentStep} total={steps.length} />
      <StepContent step={steps[currentStep]} />
      <NavigationControls 
        onNext={() => setCurrentStep(currentStep + 1)}
        onPrevious={() => setCurrentStep(currentStep - 1)}
      />
    </div>
  );
}</code></pre><h3>2. Customizable Interface</h3><p>Allow users to adjust the interface to their needs:</p><ul><li>Font size and spacing controls</li><li>Color and contrast adjustments</li><li>Animation and motion preferences</li><li>Reading speed controls</li></ul><h2>Testing with Neurodiverse Users</h2><p>Include people with cognitive differences in your testing process:</p><ol><li>Recruit diverse participants</li><li>Create comfortable testing environments</li><li>Allow extra time for tasks</li><li>Focus on task completion rather than speed</li><li>Gather qualitative feedback about cognitive load</li></ol><p>Cognitive accessibility benefits everyone by creating clearer, more intuitive digital experiences.</p>`,
			Excerpt:     "Learn how to design inclusive digital experiences for neurodiversity, focusing on cognitive accessibility principles and testing strategies.",
			Author:      "Dr. Maya Patel",
			Published:   true,
			Featured:    false,
			Tags:        "cognitive accessibility, neurodiversity, ADHD, autism, dyslexia, inclusive design",
			MetaTitle:   "Cognitive Accessibility: Designing for Neurodiversity and Inclusion",
			MetaDesc:    "Master cognitive accessibility design principles for neurodiversity, including ADHD, autism, and dyslexia considerations.",
		},
		{
			Title:       "Accessibility Testing Automation: Tools and Strategies",
			Slug:        "accessibility-testing-automation-tools-strategies",
			Content:     `<h2>The Importance of Automated Accessibility Testing</h2><p>While manual testing remains crucial, automated accessibility testing helps catch issues early and ensures consistent compliance across large applications. TechnoPrise Global integrates accessibility testing into every stage of development.</p><h2>Essential Automated Testing Tools</h2><h3>1. axe-core Integration</h3><pre><code>// Jest + axe-core integration
import { axe, toHaveNoViolations } from 'jest-axe';

expect.extend(toHaveNoViolations);

test('should not have accessibility violations', async () => {
  const { container } = render(<BlogPost />);
  const results = await axe(container);
  expect(results).toHaveNoViolations();
});</code></pre><h3>2. Lighthouse CI</h3><p>Integrate Lighthouse accessibility audits into your CI/CD pipeline:</p><pre><code># .github/workflows/accessibility.yml
name: Accessibility Audit
on: [push, pull_request]
jobs:
  lighthouse:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run Lighthouse CI
        uses: treosh/lighthouse-ci-action@v8
        with:
          configPath: './lighthouserc.json'</code></pre><h3>3. Pa11y for Command Line Testing</h3><pre><code># Install pa11y globally
npm install -g pa11y

# Test a single page
pa11y https://example.com

# Test multiple pages
pa11y-ci --sitemap https://example.com/sitemap.xml</code></pre><h2>Building a Comprehensive Testing Strategy</h2><ol><li><strong>Unit Level:</strong> Test individual components with axe-core</li><li><strong>Integration Level:</strong> Test page interactions and workflows</li><li><strong>End-to-End:</strong> Test complete user journeys</li><li><strong>Performance:</strong> Monitor accessibility performance metrics</li></ol><h2>Setting Up Accessibility Gates</h2><p>Prevent accessibility regressions by setting up quality gates:</p><pre><code>// Accessibility threshold configuration
module.exports = {
  ci: {
    collect: {
      url: ['http://localhost:3000'],
      numberOfRuns: 3
    },
    assert: {
      assertions: {
        'categories:accessibility': ['error', {minScore: 0.9}]
      }
    }
  }
};</code></pre><h2>Monitoring and Reporting</h2><ul><li>Set up accessibility dashboards</li><li>Track accessibility metrics over time</li><li>Generate regular accessibility reports</li><li>Alert teams to new violations</li></ul><p>Automated testing accelerates accessibility compliance while maintaining high standards across your entire application.</p>`,
			Excerpt:     "Master automated accessibility testing with comprehensive tools, CI/CD integration, and strategies for maintaining accessibility compliance at scale.",
			Author:      "Carlos Rodriguez",
			Published:   true,
			Featured:    true,
			Tags:        "accessibility testing, automation, axe-core, lighthouse, CI/CD, quality assurance",
			MetaTitle:   "Accessibility Testing Automation: Tools and CI/CD Integration",
			MetaDesc:    "Learn to implement automated accessibility testing with axe-core, Lighthouse CI, and comprehensive testing strategies.",
		},
		{
			Title:       "Creating Accessible Data Visualizations",
			Slug:        "accessible-data-visualizations-design-guide",
			Content:     `<h2>The Challenge of Accessible Data Visualization</h2><p>Data visualizations are powerful tools for communication, but they often exclude users with visual impairments or cognitive differences. Creating accessible charts and graphs requires thoughtful design and alternative representations.</p><h2>Fundamental Principles</h2><ul><li><strong>Multiple Formats:</strong> Provide data in visual, textual, and tabular formats</li><li><strong>Color Independence:</strong> Never rely solely on color to convey information</li><li><strong>Clear Labeling:</strong> Use descriptive titles, axis labels, and legends</li><li><strong>Logical Structure:</strong> Organize data in a meaningful sequence</li></ul><h2>Implementation Techniques</h2><h3>1. Alternative Text for Charts</h3><pre><code>// Comprehensive alt text for charts
function AccessibleChart({ data, title, description }) {
  const altText = title + '. ' + description + '. ' +
    'Data shows: ' + data.map(item => 
      item.label + ': ' + item.value
    ).join(', ') + '.';
  
  return (
    <div role="img" aria-label={altText}>
      <Chart data={data} />
      <DataTable data={data} className="sr-only" />
    </div>
  );
}</code></pre><h3>2. Sonification for Data</h3><p>Convert data to audio representations:</p><pre><code>// Audio representation of data trends
function SonifyData(data) {
  const audioContext = new AudioContext();
  data.forEach((point, index) => {
    const frequency = 200 + (point.value * 10);
    const startTime = audioContext.currentTime + (index * 0.5);
    playTone(frequency, startTime, 0.4);
  });
}</code></pre><h3>3. Tactile Patterns</h3><p>Use different patterns and textures for distinguishing data series:</p><pre><code>/* CSS patterns for accessibility */
.data-series-1 { fill: url(#diagonal-lines); }
.data-series-2 { fill: url(#dots); }
.data-series-3 { fill: url(#cross-hatch); }</code></pre><h2>Interactive Accessibility</h2><ul><li><strong>Keyboard Navigation:</strong> Allow users to navigate through data points</li><li><strong>Focus Management:</strong> Clearly indicate which data point is selected</li><li><strong>Voice Announcements:</strong> Announce data values as users navigate</li><li><strong>Zoom and Pan:</strong> Enable users to explore data at different scales</li></ul><h2>Testing Data Visualizations</h2><ol><li>Test with screen readers</li><li>Verify keyboard navigation</li><li>Check color contrast ratios</li><li>Validate with color blindness simulators</li><li>Test cognitive load with complex datasets</li></ol><p>Accessible data visualization opens insights to everyone, creating more inclusive and impactful data storytelling.</p>`,
			Excerpt:     "Learn to create accessible data visualizations with alternative formats, sonification, and inclusive design principles for charts and graphs.",
			Author:      "Dr. Lisa Zhang",
			Published:   true,
			Featured:    false,
			Tags:        "data visualization, accessibility, charts, graphs, sonification, inclusive design",
			MetaTitle:   "Accessible Data Visualizations: Inclusive Charts and Graphs",
			MetaDesc:    "Master accessible data visualization design with alternative formats, sonification, and comprehensive accessibility techniques.",
		},
		{
			Title:       "Legal Compliance and Accessibility: ADA, WCAG, and Beyond",
			Slug:        "legal-compliance-accessibility-ada-wcag-standards",
			Content:     `<h2>The Legal Landscape of Digital Accessibility</h2><p>Digital accessibility isn't just good practice—it's increasingly a legal requirement. Understanding compliance standards helps organizations avoid litigation while creating inclusive experiences.</p><h2>Key Accessibility Laws and Standards</h2><h3>Americans with Disabilities Act (ADA)</h3><p>While the ADA doesn't explicitly mention websites, courts increasingly apply it to digital spaces:</p><ul><li>Title III applies to places of public accommodation</li><li>No specific technical standards, but WCAG is often referenced</li><li>Enforcement through private lawsuits</li></ul><h3>Section 508 (Federal Agencies)</h3><p>Requires federal agencies to make electronic content accessible:</p><ul><li>Applies to all federal websites and applications</li><li>References WCAG 2.0 Level AA as standard</li><li>Includes procurement requirements</li></ul><h3>European Accessibility Act</h3><p>Comprehensive accessibility legislation for EU member states:</p><ul><li>Covers websites, mobile apps, and digital services</li><li>Based on EN 301 549 standard</li><li>Includes enforcement mechanisms and penalties</li></ul><h2>WCAG Compliance Levels</h2><h3>Level A (Minimum)</h3><ul><li>Basic accessibility features</li><li>Essential for any public-facing site</li><li>Addresses major barriers</li></ul><h3>Level AA (Standard)</h3><ul><li>Recommended compliance level</li><li>Required for most government sites</li><li>Covers most user needs</li></ul><h3>Level AAA (Enhanced)</h3><ul><li>Highest level of accessibility</li><li>Not required for entire sites</li><li>Applied to specific content areas</li></ul><h2>Building a Compliance Strategy</h2><ol><li><strong>Accessibility Audit:</strong> Assess current compliance status</li><li><strong>Risk Assessment:</strong> Identify high-risk areas and user flows</li><li><strong>Remediation Plan:</strong> Prioritize fixes based on impact and effort</li><li><strong>Training Program:</strong> Educate teams on accessibility requirements</li><li><strong>Ongoing Monitoring:</strong> Implement regular testing and reviews</li></ol><h2>Documentation and Evidence</h2><p>Maintain comprehensive accessibility documentation:</p><ul><li>Accessibility statements</li><li>Testing reports and audit results</li><li>Remediation timelines and progress</li><li>User feedback and complaint resolution</li><li>Training records and certifications</li></ul><h2>Working with Legal Teams</h2><p>Collaborate effectively with legal counsel:</p><ul><li>Provide technical expertise on accessibility standards</li><li>Translate compliance requirements into actionable tasks</li><li>Document accessibility efforts and improvements</li><li>Prepare for potential accessibility audits</li></ul><p>Proactive accessibility compliance protects organizations while creating better experiences for all users.</p>`,
			Excerpt:     "Navigate the complex legal landscape of digital accessibility with comprehensive guidance on ADA, WCAG, Section 508, and compliance strategies.",
			Author:      "Jennifer Martinez, JD",
			Published:   true,
			Featured:    false,
			Tags:        "legal compliance, ADA, WCAG, Section 508, accessibility law, digital rights",
			MetaTitle:   "Legal Compliance and Digital Accessibility: ADA and WCAG Guide",
			MetaDesc:    "Understand digital accessibility legal requirements including ADA, WCAG, Section 508, and compliance strategies for organizations.",
		},
		{
			Title:       "Accessibility in E-commerce: Converting All Customers",
			Slug:        "accessibility-ecommerce-inclusive-shopping-experience",
			Content:     `<h2>The Business Case for Accessible E-commerce</h2><p>Accessible e-commerce isn't just about compliance—it's about reaching the $13 trillion disability market. TechnoPrise Global helps businesses create shopping experiences that convert for everyone.</p><h2>Critical E-commerce Accessibility Areas</h2><h3>1. Product Discovery</h3><ul><li><strong>Search Functionality:</strong> Voice search and predictive text</li><li><strong>Filtering and Sorting:</strong> Keyboard accessible controls</li><li><strong>Product Images:</strong> Comprehensive alt text and zoom functionality</li><li><strong>Product Videos:</strong> Captions and audio descriptions</li></ul><h3>2. Shopping Cart and Checkout</h3><pre><code>// Accessible form validation
function AccessibleCheckout() {
  const [errors, setErrors] = useState({});
  
  const validateField = (field, value) => {
    const newErrors = { ...errors };
    if (!value) {
      newErrors[field] = field + ' is required';
      announceError('Error: ' + newErrors[field]);
    } else {
      delete newErrors[field];
      announceSuccess(field + ' validated successfully');
    }
    setErrors(newErrors);
  };
  
  return (
    <form aria-label="Checkout form">
      <fieldset>
        <legend>Billing Information</legend>
        {/* Form fields with proper labels and error handling */}
      </fieldset>
    </form>
  );
}</code></pre><h3>3. Payment Processing</h3><ul><li>Clear error messages for payment failures</li><li>Multiple payment method options</li><li>Secure, accessible payment forms</li><li>Progress indicators for multi-step processes</li></ul><h2>Mobile Commerce Accessibility</h2><p>Mobile shopping requires special attention:</p><ul><li><strong>Touch Targets:</strong> Minimum 44px for all interactive elements</li><li><strong>Gesture Alternatives:</strong> Button alternatives for swipe actions</li><li><strong>Voice Shopping:</strong> Integration with voice assistants</li><li><strong>One-Handed Operation:</strong> Thumb-friendly navigation</li></ul><h2>Accessibility Features That Drive Sales</h2><ol><li><strong>Voice Search:</strong> "Find red dresses under $100"</li><li><strong>Smart Recommendations:</strong> AI-powered accessible product suggestions</li><li><strong>Wishlist Management:</strong> Easy saving and organization</li><li><strong>Order Tracking:</strong> Clear status updates and notifications</li></ol><h2>Testing E-commerce Accessibility</h2><h3>Critical User Journeys</h3><ul><li>Product search and discovery</li><li>Adding items to cart</li><li>Checkout process completion</li><li>Account creation and management</li><li>Order history and reordering</li></ul><h3>Assistive Technology Testing</h3><ul><li>Screen reader navigation</li><li>Voice control shopping</li><li>Switch navigation for checkout</li><li>Magnification software compatibility</li></ul><h2>Measuring Accessibility ROI</h2><p>Track the business impact of accessibility improvements:</p><ul><li>Conversion rate improvements</li><li>Reduced cart abandonment</li><li>Increased customer satisfaction scores</li><li>Expanded market reach</li><li>Reduced customer service inquiries</li></ul><p>Accessible e-commerce creates better experiences for all customers while opening new market opportunities.</p>`,
			Excerpt:     "Transform your e-commerce platform with accessibility best practices that improve conversions and reach the $13 trillion disability market.",
			Author:      "Rachel Thompson",
			Published:   false,
			Featured:    false,
			Tags:        "e-commerce, accessibility, online shopping, conversion optimization, mobile commerce",
			MetaTitle:   "E-commerce Accessibility: Inclusive Shopping Experiences That Convert",
			MetaDesc:    "Learn to create accessible e-commerce experiences that improve conversions and reach customers with disabilities.",
		},
	}

	for _, blog := range sampleBlogs {
		if err := db.Create(&blog).Error; err != nil {
			return fmt.Errorf("failed to create blog: %v", err)
		}
	}

	log.Printf("✅ Successfully seeded database with %d blog posts", len(sampleBlogs))
	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
