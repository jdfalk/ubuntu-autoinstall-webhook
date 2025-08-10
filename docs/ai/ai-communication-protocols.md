<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [AI Communication Protocols](#ai-communication-protocols)
  - [General Communication Guidelines](#general-communication-guidelines)
    - [Response Structure](#response-structure)
    - [Code Suggestions](#code-suggestions)
  - [Project-Specific Protocols](#project-specific-protocols)
    - [Architecture Discussions](#architecture-discussions)
    - [Implementation Suggestions](#implementation-suggestions)
    - [Security-Related Communications](#security-related-communications)
  - [Response Formats for Specific Tasks](#response-formats-for-specific-tasks)
    - [For Component Implementation](#for-component-implementation)
  - [Implementation](#implementation)
  - [Usage Example](#usage-example)
  - [Test Cases](#test-cases)
  - [Security Considerations](#security-considerations)
  - [Performance Considerations](#performance-considerations)
- [Bug Fix: [Brief Description]](#bug-fix-brief-description)
  - [Issue](#issue)
  - [Root Cause](#root-cause)
  - [Proposed Fix](#proposed-fix)
  - [Testing](#testing)
- [Code Review: [File or Component Name]](#code-review-file-or-component-name)
  - [Overview](#overview)
  - [Key Observations](#key-observations)
    - [Strengths](#strengths)
    - [Areas for Improvement](#areas-for-improvement)
  - [Security Review](#security-review)
  - [Performance Review](#performance-review)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# AI Communication Protocols

This document establishes protocols for how AI assistants should communicate
with developers working on the Ubuntu Autoinstall Webhook project. Following
these guidelines will ensure consistent, productive interactions that align with
project goals.

## General Communication Guidelines

### Response Structure

1. **For Implementation Suggestions**
   - Begin with a concise summary of what you're proposing
   - Explain the reasoning behind your approach
   - Provide code examples with clear comments
   - Discuss potential alternatives
   - Highlight security or performance implications

2. **For Code Reviews**
   - Start with a summary assessment
   - Group feedback by category (bugs, optimizations, style)
   - Provide specific recommendations for improvements
   - Include code examples for complex changes
   - Reference relevant project standards or patterns

3. **For Technical Questions**
   - Directly address the question first
   - Provide context and background information
   - Include relevant code examples or documentation references
   - Suggest related considerations if applicable

### Code Suggestions

When suggesting code:

1. **Use Clear File Path Headers**

   ```
   // File: internal/api/systems.go
   ```

2. **Include Context for Modifications**

   ```
   // Existing function to modify:
   func GetSystem(id string) (*System, error) {
     // ... existing code
   }

   // Modified version:
   func GetSystem(id string) (*System, error) {
     // ... modified code
   }
   ```

3. **Explain Non-Obvious Code**
   ```go
   // Using a custom error type to allow clients to distinguish between
   // "not found" errors and other types of failures
   if system == nil {
     return nil, ErrSystemNotFound
   }
   ```

## Project-Specific Protocols

### Architecture Discussions

When discussing architectural changes:

1. **Refer to Existing Architecture Documents**
   - Reference relevant sections in the architecture overview
   - Explain how your suggestion aligns with established patterns
   - Identify components affected by proposed changes

2. **Focus on Interfaces First**
   - Begin with interface definitions rather than implementations
   - Explain interaction patterns between components
   - Consider backward compatibility concerns

3. **Consider Deployment Impact**
   - Discuss how changes affect deployment requirements
   - Address migration concerns for existing data
   - Consider operational complexity implications

### Implementation Suggestions

When providing implementation advice:

1. **Follow the Technical Requirements**
   - Ensure suggestions align with documented requirements
   - Reference specific sections of the technical requirements
   - Highlight where trade-offs might be necessary

2. **Prioritize Project Standards**
   - Follow Go best practices and project coding standards
   - Use established project patterns for consistency
   - Maintain separation of concerns between components

3. **Include Tests with Implementations**
   - Provide unit test examples alongside implementation code
   - Demonstrate both success and error cases
   - Show how to test edge conditions

### Security-Related Communications

For security-related topics:

1. **Explicitly Identify Security Implications**
   - Highlight security considerations with clear markers
   - Explain the nature of potential vulnerabilities
   - Reference relevant security best practices

2. **Provide Secure Alternatives**
   - When identifying security issues, always suggest secure alternatives
   - Explain the security benefits of recommended approaches
   - Reference industry standards where applicable

3. **Consider Security Throughout**
   - Integrate security considerations into all discussions
   - Don't treat security as a separate topic
   - Consider authentication, authorization, and data protection

## Response Formats for Specific Tasks

### For Component Implementation

````
# [Component Name] Implementation

## Overview
[Brief description of the component and its purpose]

## Interface Definition
```go
// [Interface code with documentation comments]
````

## Implementation

```go
// [Implementation code with documentation comments]
```

## Usage Example

```go
// [Example showing how to use this component]
```

## Test Cases

```go
// [Example test cases]
```

## Security Considerations

[Security-related notes, if applicable]

## Performance Considerations

[Performance-related notes, if applicable]

```

### For Bug Fixes

```

# Bug Fix: [Brief Description]

## Issue

[Description of the bug and its impact]

## Root Cause

[Analysis of what's causing the issue]

## Proposed Fix

```go
// [Code changes to resolve the issue]
```

## Testing

[How to verify the fix works correctly]

```

### For Code Review

```

# Code Review: [File or Component Name]

## Overview

[General assessment of the code]

## Key Observations

### Strengths

- [Positive aspects of the code]

### Areas for Improvement

1. [Issue #1]

   ```go
   // Current code

   // Suggested improvement
   ```

2. [Issue #2]

   ```go
   // Current code

   // Suggested improvement
   ```

## Security Review

[Security-related observations]

## Performance Review

[Performance-related observations]

```

## Special Considerations

### Handling Uncertainties

When you're uncertain about aspects of the implementation:

1. **Clearly State Assumptions**
   - Identify what you're assuming and why
   - Explain how different assumptions would change your recommendation

2. **Provide Multiple Options**
   - Present alternative approaches with pros and cons
   - Explain the conditions under which each would be preferred

3. **Request Clarification**
   - Clearly articulate what additional information would help
   - Suggest how that information might be obtained

### Prioritizing Tasks

When helping prioritize development work:

1. **Consider Dependencies**
   - Identify component dependencies
   - Suggest logical implementation ordering

2. **Balance Technical Debt**
   - Consider when to address technical debt vs. new features
   - Explain the long-term implications of postponing refactoring

3. **Focus on Core Functionality First**
   - Prioritize components needed for basic functionality
   - Defer advanced features until core features are stable

### Collaborative Development

When working alongside human developers:

1. **Build on Their Ideas**
   - Acknowledge their approach before suggesting modifications
   - Frame suggestions as enhancements rather than replacements

2. **Respect Design Decisions**
   - Work within established design patterns
   - Seek to understand rationale before suggesting alternatives

3. **Provide Educational Context**
   - Explain the reasoning behind suggestions
   - Include references to Go best practices or design patterns
   - Share knowledge that helps developers grow their skills
```
