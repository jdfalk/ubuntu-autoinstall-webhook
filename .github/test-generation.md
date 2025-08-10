<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Test Generation Guidelines](#test-generation-guidelines)
  - [Test Structure](#test-structure)
  - [Test Naming Conventions](#test-naming-conventions)
  - [Test Types](#test-types)
  - [Unit Test Best Practices](#unit-test-best-practices)
  - [Test Coverage Guidelines](#test-coverage-guidelines)
  - [Mocking Guidelines](#mocking-guidelines)
  - [Assertion Best Practices](#assertion-best-practices)
  - [Data Management](#data-management)
  - [Test Organization](#test-organization)
  - [Testing Anti-Patterns to Avoid](#testing-anti-patterns-to-avoid)
  - [Special Testing Considerations](#special-testing-considerations)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Test Generation Guidelines

## Test Structure

Follow this structure when creating tests:

```markdown
[Setup] - Prepare the test environment and inputs [Exercise] - Execute the
functionality being tested [Verify] - Check that the results match expectations
[Teardown] - Clean up any resources (if needed)
```

## Test Naming Conventions

- Use descriptive names that indicate what's being tested
- Follow the pattern: `test[UnitOfWork_StateUnderTest_ExpectedBehavior]`
- Examples:
  - `testLoginService_InvalidCredentials_ReturnsError`
  - `testUserRepository_SaveDuplicateEmail_ThrowsException`
  - `testCalculator_DivideByZero_ThrowsException`

## Test Types

- **Unit Tests**: Focus on a single function, method, or class in isolation
- **Integration Tests**: Verify interactions between components
- **Functional Tests**: Test complete features from a user perspective
- **Performance Tests**: Measure response times and resource usage
- **Security Tests**: Identify vulnerabilities and validate safeguards
- **Accessibility Tests**: Ensure compliance with accessibility standards

## Unit Test Best Practices

- Test one specific behavior per test case
- Mock external dependencies
- Use setup/teardown to avoid code duplication
- Write deterministic tests (consistent results)
- Cover both happy paths and edge cases
- Test public interfaces, not internal implementation
- Keep tests independent of each other
- Use clear assertions with meaningful messages

## Test Coverage Guidelines

- Aim for high coverage of business-critical code
- Prioritize coverage of complex logic and edge cases
- Don't obsess over 100% coverage at the expense of meaningful tests
- Focus on code paths rather than simple line coverage
- Consider risk and complexity when prioritizing what to test

## Mocking Guidelines

- Mock external services and dependencies
- Use stubs for predetermined responses
- Use spies to verify interactions
- Avoid excessive mocking that reduces test value
- Mock at the appropriate abstraction level
- Document mock behavior clearly

## Assertion Best Practices

- Use specific assertions (e.g., `assertEquals` instead of `assertTrue`)
- Check only one logical assertion per test
- Write clear assertion messages explaining expected vs actual
- Verify the right things: state changes, interactions, or exceptions
- For collections, verify content regardless of order when appropriate

## Data Management

- Use test data factories or builders for complex objects
- Create minimal test data sets that focus on the test requirements
- Avoid shared mutable test data between tests
- Consider using test database fixtures for integration tests
- Clean up test data reliably after tests complete

## Test Organization

- Group tests by feature or component
- Separate slow tests from fast tests
- Use test suites to organize related tests
- Maintain parallel structure between code and tests
- Place tests in a location that mirrors the code structure

## Testing Anti-Patterns to Avoid

- Flaky tests with inconsistent results
- Tests that depend on external services without mocks
- Overly complex test setups
- Testing implementation details instead of behavior
- Excessive use of sleep/delay in asynchronous tests
- Ignoring test failures
- Testing trivial code with no logic
- Writing tests after code is already in production

## Special Testing Considerations

- **Asynchronous Code**: Use async/await patterns and avoid arbitrary delays
- **APIs**: Test request validation, response formats, and error cases
- **UIs**: Test component rendering, user interactions, and state management
- **Data Access**: Test query correctness and error handling
- **Security**: Test authorization, input validation, and error handling
