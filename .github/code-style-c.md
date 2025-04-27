<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [C Code Style Guide](#c-code-style-guide)
  - [Header Files](#header-files)
  - [Naming Conventions](#naming-conventions)
  - [Functions](#functions)
  - [Comments](#comments)
  - [Formatting](#formatting)
  - [Variables](#variables)
  - [Memory Management](#memory-management)
  - [Error Handling](#error-handling)
  - [Best Practices](#best-practices)
  - [Tools](#tools)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# C Code Style Guide

This guide follows Google's C Style Guide and industry standards.

## Header Files

- Every `.c` file should have an associated `.h` file
- Use header guards with unique names
- Include headers in the following order:
  1. Main module header (matching the `.c` file name)
  2. C standard library headers
  3. Other library headers
  4. Your project's headers
- Declare function prototypes in headers, implement in source files

```c
/* In file.h */
#ifndef PROJECT_FILE_H_
#define PROJECT_FILE_H_

/* Function declarations */

#endif  /* PROJECT_FILE_H_ */
```

## Naming Conventions

- File names: `snake_case.c` and `snake_case.h`
- Functions and variables: `snake_case`
- Constants and macros: `ALL_CAPS_WITH_UNDERSCORES`
- Type names (typedef, structs, enums): `snake_case_t` or `PascalCase`
- Global variables: `g_prefixed_name`
- Struct members: `snake_case` without special prefix

## Functions

- Keep functions short and focused (aim for under 40 lines)
- Declare function parameters on a single line if possible
- If parameters don't fit on one line, align parameter names
- Use consistent parameter ordering across related functions
- Return error codes (typically integers) rather than using exceptions

```c
/* Good: parameters fit on one line */
void process_data(const char* input, size_t length, int options);

/* Good: parameters aligned when split across multiple lines */
int complex_function(
    const char* input_file_path,
    char* output_buffer,
    size_t buffer_size,
    int processing_options);
```

## Comments

- Use `/* */` style comments for multi-line comments
- Use `//` for single-line comments if C99 is allowed
- Every function, struct, typedef, and file should have a comment
- Comments should explain why, not what
- Use doxygen-style comments for API documentation

```c
/**
 * Processes the input data according to specified options.
 *
 * @param input Pointer to input data buffer
 * @param length Size of input data in bytes
 * @param options Processing option flags
 * @return 0 on success, error code otherwise
 */
int process_data(const char* input, size_t length, int options);
```

## Formatting

- Line length: 80 characters maximum
- Use 2 spaces for indentation (no tabs)
- Place opening braces on the same line as control statements
- Place opening braces for functions on a new line
- Use braces even for single-line blocks
- Switch statements: indent case statements, use braces for blocks

```c
/* Function definition style */
int calculate_sum(int a, int b)
{
  return a + b;
}

/* Control flow style */
if (condition) {
  do_something();
} else {
  do_alternative();
}

/* Switch statement style */
switch (value) {
  case 0:
    handle_zero();
    break;
  case 1: {
    int temp = calculate_something();
    handle_one(temp);
    break;
  }
  default:
    handle_default();
    break;
}
```

## Variables

- Declare one variable per line
- Initialize variables at declaration when possible
- Keep variable declarations at the start of blocks
- Use const whenever a variable should not be modified
- Avoid global variables when possible

## Memory Management

- Always check return values of memory allocation functions
- Free dynamically allocated memory in the reverse order of allocation
- Use a consistent pattern for allocating and freeing resources
- Consider using a cleanup function pattern for complex resources

## Error Handling

- Always check return values and handle errors
- Use a consistent error handling approach (return codes, goto cleanup, etc.)
- Document error codes and their meanings
- Avoid silent failures

## Best Practices

- Prefer C99 or C11 features when available
- Use enums instead of #define for related constants
- Use inline functions (C99) instead of function-like macros
- Use static for file-scope functions
- Use const for pointers that should not modify their targets
- Avoid using non-standard or platform-specific features
- Validate all input parameters

## Tools

- Use a static analyzer like clang-analyze or Coverity
- Use a style checker like clang-format
- Configure tools to match this style guide
