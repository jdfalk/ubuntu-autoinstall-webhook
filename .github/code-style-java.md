<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Java Code Style Guide](#java-code-style-guide)
  - [Source File Structure](#source-file-structure)
  - [Naming Conventions](#naming-conventions)
  - [Formatting](#formatting)
  - [Import Statements](#import-statements)
  - [Class Declaration](#class-declaration)
  - [Javadoc](#javadoc)
  - [Programming Practices](#programming-practices)
  - [Exception Handling](#exception-handling)
  - [Concurrency](#concurrency)
  - [Tools](#tools)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Java Code Style Guide

This guide follows Google's Java Style Guide.

## Source File Structure

- Each source file contains a single public class or interface
- Files are encoded in UTF-8
- Use Unix-style line endings (LF)
- Files must end with a newline character
- Maximum line length is 100 characters

## Naming Conventions

- `UpperCamelCase` for class and interface names
- `lowerCamelCase` for method and variable names
- `UPPER_SNAKE_CASE` for constants
- Package names in all lowercase, no underscores: `com.example.deepspace`
- Test classes should be named ending with `Test`: `HashTest` or `HashIntegrationTest`

## Formatting

- Use 2 space indentation (no tabs)
- Use one statement per line
- Column limit: 100 characters
- No line break before opening brace
- Line break after opening brace and before closing brace
- Empty blocks may be concise: `{}` on same line

```java
// Correct
void doNothing() {}

// Correct for non-empty blocks
void doSomething() {
  performAction();
}
```

## Import Statements

- No wildcard imports
- No static imports (except for common test libraries)
- Ordering: static imports first, then non-static imports
- Within each group, imports appear in ASCII sort order
- No line-wrapping in import statements
- No unused imports

```java
import java.util.List;
import java.util.Map;
import javax.annotation.Nullable;
import javax.inject.Inject;
```

## Class Declaration

- Each top-level class resides in a source file of its own
- Class members appear in a logical order (not by scope or accessibility)
- Recommended order:
  1. Static fields
  2. Instance fields
  3. Constructors
  4. Methods
- Overloads appear sequentially

## Javadoc

- Write Javadoc for every public class and method
- Use complete sentences
- Use third person, not second person
- Use `@param`, `@return`, and `@throws` tags as appropriate

```java
/**
 * Returns the element at the specified position in this list.
 *
 * @param index index of the element to return
 * @return the element at the specified position
 * @throws IndexOutOfBoundsException if the index is out of range
 */
public E get(int index) {
  // Implementation
}
```

## Programming Practices

- Always use braces with `if`, `else`, `for`, `do` and `while` statements
- Use `@Override` annotation whenever applicable
- Caught exceptions should not be ignored
- Static members should be accessed via class name: `Foo.aStaticMethod()`
- Numeric literals should use underscores for readability: `1_000_000`
- Use diamond operator with generics when possible: `Map<String, List<String>> map = new HashMap<>()`

## Exception Handling

- Throw specific exceptions rather than generic ones
- Document exceptions with `@throws` Javadoc tags
- Use try-with-resources for AutoCloseable resources
- Don't catch exceptions you can't handle properly

## Concurrency

- Prefer higher-level concurrency utilities (from `java.util.concurrent`)
- Prefer `synchronized` methods over synchronized blocks
- Document thread-safety for each class

## Tools

- Use Google Java Format tool
- Configure IDE to match Google Style Guide
- Use CheckStyle or similar tools with Google's configuration
