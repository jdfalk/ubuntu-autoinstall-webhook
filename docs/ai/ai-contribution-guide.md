<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [AI Contribution Guide](#ai-contribution-guide)
  - [How to Approach the Project](#how-to-approach-the-project)
  - [Code Implementation Guidelines](#code-implementation-guidelines)
    - [General Coding Standards](#general-coding-standards)
    - [Component Implementation](#component-implementation)
    - [Testing Approach](#testing-approach)
  - [Security Considerations](#security-considerations)
  - [Common Implementation Patterns](#common-implementation-patterns)
    - [Repository Pattern](#repository-pattern)
    - [Service Pattern](#service-pattern)
    - [Handler Pattern](#handler-pattern)
  - [Documentation Guidelines](#documentation-guidelines)
  - [Suggesting Implementation](#suggesting-implementation)
  - [Project-Specific Tips](#project-specific-tips)
  - [Working with Existing Documentation](#working-with-existing-documentation)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# AI Contribution Guide

This guide provides recommendations and best practices for AI assistants contributing to the Ubuntu Autoinstall Webhook project. It outlines how to interpret existing documentation, implement components according to project standards, and ensure code quality.

## How to Approach the Project

1. **Understand the Big Picture First**
   - Review the architecture overview (`docs/architecture/overview.md`)
   - Understand system components and their interactions
   - Review the technical design document (`docs/technical/technical-design.md`)

2. **Focus on Current Priorities**
   - Check `docs/ai/implementation-status.md` for current focus areas
   - Prioritize components marked as "In-Progress" 🔄
   - Understand dependencies before implementing new components

3. **Respect the Design Decisions**
   - Follow established architecture patterns
   - Maintain separation of concerns between components
   - Adhere to the interface-based design approach

## Code Implementation Guidelines

### General Coding Standards

The project follows Go best practices and conventions:

1. **Code Organization**
   - Place all internal code in appropriate subdirectories of `internal/`
   - Public libraries that may be consumed by external applications go in `pkg/`
   - Maintain clear separation between components

2. **Code Style**
   - Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
   - Use `gofmt` for consistent formatting
   - Prefer clarity over cleverness

3. **Naming Conventions**
   - Use descriptive names for variables, functions, and packages
   - Follow Go naming conventions:
     - CamelCase for exported names
     - lowerCamelCase for non-exported names
     - Acronyms should be consistently cased (e.g., `httpServer` or `HTTPServer`)

4. **Error Handling**
   - Check all errors and handle them appropriately
   - Provide context with errors using `fmt.Errorf("doing something: %w", err)`
   - Use custom error types for domain-specific errors

### Component Implementation

When implementing a component:

1. **Start with Interfaces**
   - Define the interface first
   - Implement against the interface, not concrete types
   - This allows for easy mocking during tests

2. **Follow Dependency Injection**
   - Components should receive their dependencies, not create them
   - This improves testability and flexibility

3. **Keep Components Focused**
   - Each component should have a single responsibility
   - Avoid feature creep within components

Example interface pattern:

```go
// Service interface
type TemplateService interface {
    GetTemplate(id string) (*Template, error)
    CreateTemplate(template *Template) error
    UpdateTemplate(template *Template) error
    DeleteTemplate(id string) error
}

// Implementation
type templateService struct {
    repo TemplateRepository
    fileEditor FileEditor
}

// Constructor with dependency injection
func NewTemplateService(repo TemplateRepository, fileEditor FileEditor) TemplateService {
    return &templateService{
        repo: repo,
        fileEditor: fileEditor,
    }
}
```

### Testing Approach

1. **Test-Driven Development**
   - Write tests before or alongside implementation
   - Focus on behavior, not implementation details

2. **Test Types**
   - Unit tests for individual functions and methods
   - Integration tests for component interactions
   - End-to-end tests for critical workflows

3. **Mocking**
   - Use interfaces to allow for easy mocking
   - Focus on testing behavior, not implementation details

Example test pattern:

```go
func TestTemplateService_GetTemplate(t *testing.T) {
    // Arrange
    mockRepo := &MockTemplateRepository{}
    mockFileEditor := &MockFileEditor{}

    expectedTemplate := &Template{ID: "123", Name: "test"}
    mockRepo.On("GetByID", "123").Return(expectedTemplate, nil)

    service := NewTemplateService(mockRepo, mockFileEditor)

    // Act
    result, err := service.GetTemplate("123")

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedTemplate, result)
    mockRepo.AssertExpectations(t)
}
```

## Security Considerations

When implementing security-sensitive components:

1. **Authentication & Authorization**
   - Always validate user permissions before actions
   - Use secure password handling practices
   - Implement proper session management

2. **Data Protection**
   - Encrypt sensitive data at rest
   - Use TLS for all network communications
   - Validate all user input

3. **Secure Development Practices**
   - Avoid hardcoded credentials
   - Use parameterized queries for database operations
   - Follow the principle of least privilege

4. **Error Handling**
   - Don't expose sensitive information in error messages
   - Log security events appropriately
   - Fail securely (deny by default)

## Common Implementation Patterns

### Repository Pattern

Used for database access:

```go
// Repository interface
type SystemRepository interface {
    GetByID(id string) (*System, error)
    GetByMacAddress(mac string) (*System, error)
    Create(system *System) error
    Update(system *System) error
    Delete(id string) error
    List(filter SystemFilter) ([]*System, error)
}

// Implementation for SQLite
type sqliteSystemRepository struct {
    db *sql.DB
}
```

### Service Pattern

Used for business logic:

```go
// Service interface
type InstallationService interface {
    StartInstallation(systemID string, templateID string) (*Installation, error)
    GetInstallation(id string) (*Installation, error)
    CancelInstallation(id string) error
    GetInstallationLogs(id string) ([]LogEntry, error)
}

// Implementation
type installationService struct {
    repo InstallationRepository
    systemRepo SystemRepository
    templateRepo TemplateRepository
    fileEditor FileEditor
}
```

### Handler Pattern

Used for API endpoints:

```go
// Handler
type SystemsHandler struct {
    systemService SystemService
    logger Logger
}

func (h *SystemsHandler) GetSystem(w http.ResponseWriter, r *http.Request) {
    // Extract ID from request
    // Call service
    // Handle response
}
```

## Documentation Guidelines

When updating or adding code, ensure documentation is kept in sync:

1. **Code Comments**
   - Document non-obvious behavior
   - Explain complex algorithms
   - Document public API functions thoroughly

2. **Package Documentation**
   - Add a package comment to each package
   - Explain the purpose and usage of the package

3. **Architecture Documentation**
   - Update component docs when changing behavior
   - Keep diagrams in sync with code changes

## Suggesting Implementation

When suggesting implementation for a component:

1. **Start with the Interface**
   - Define the public API first
   - Consider how it will be used by other components

2. **Provide Complete Examples**
   - Include both interface definition and implementation
   - Show how the component would be used by others

3. **Consider Error Cases**
   - Include proper error handling
   - Think about edge cases

4. **Include Tests**
   - Provide example tests for critical functionality
   - Demonstrate both success and failure scenarios

## Project-Specific Tips

1. **Template Processing**
   - Use Go's `text/template` package for template rendering
   - Support variable substitution and conditional logic
   - Consider hierarchical templates with inheritance

2. **Installation Process**
   - Design for idempotence
   - Implement proper status tracking
   - Consider recovery from failures

3. **API Design**
   - Follow RESTful principles
   - Use consistent response formats
   - Implement proper validation and error handling

4. **Web Interface**
   - Use standard Go templates
   - Keep JavaScript minimal and focused
   - Consider progressive enhancement

## Working with Existing Documentation

The project has comprehensive documentation that should be treated as the source of truth:

1. **Architecture Documents**
   - Define the overall system design and component interactions
   - Consult these before suggesting structural changes

2. **Technical Design**
   - Provides detailed specifications for implementation
   - Follow the patterns and approaches defined here

3. **Test Design**
   - Outlines the testing strategy and approach
   - Implement tests according to these guidelines

4. **User and Admin Guides**
   - Describe how the system should behave from a user perspective
   - Ensure implementations align with documented behavior
